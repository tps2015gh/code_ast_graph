package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"ci4-visualizer/pkg/analyzer"
	"ci4-visualizer/pkg/astparser"
	"ci4-visualizer/pkg/graph"
	"ci4-visualizer/pkg/netutil"
	"ci4-visualizer/pkg/security"
	"ci4-visualizer/pkg/tester"
	"ci4-visualizer/pkg/ui"
)

//go:embed frontend/*
var embeddedFrontend embed.FS

//go:embed parser/*
var embeddedParser embed.FS

func main() {
	projectPathFlag := flag.String("path", "", "Path to the CodeIgniter 4 project to analyze on startup")
	phpExePathFlag := flag.String("php", "C:/xampp_v8_1_25/php/php.exe", "Path to the PHP executable")
	runTestsFlag := flag.Bool("test", false, "Run basic health checks on startup")
	interactiveFlag := flag.Bool("interactive", false, "Start in interactive CLI menu mode")
	flag.Parse()

	security.CheckRepoSafety()

	parserFS, err := fs.Sub(embeddedParser, "parser")
	if err != nil {
		log.Fatalf("Failed to create sub-filesystem for parser: %v", err)
	}

	err = astparser.Init(*phpExePathFlag, parserFS)
	if err != nil {
		log.Fatalf("Failed to initialize AST parser: %v", err)
	}
	defer astparser.Cleanup()

	if *interactiveFlag || (len(os.Args) == 1) {
		startInteractiveMode()
		return
	}

	if *projectPathFlag != "" {
		log.Printf("CLI path provided: %s. Starting initial analysis...", *projectPathFlag)
		nodes, edges, analyzeErr := analyzer.PerformAnalysis(*projectPathFlag)
		analyzer.GraphDataMutex.Lock()
		if analyzeErr != nil {
			log.Printf("Initial CLI analysis failed: %v", analyzeErr)
			analyzer.CurrentGraphData = graph.GraphData{Nodes: []graph.Node{{ID: "error", Label: fmt.Sprintf("Initial analysis error: %v", analyzeErr), Type: "error"}}}
		} else {
			analyzer.CurrentGraphData = graph.GraphData{Nodes: nodes, Edges: edges}
			analyzer.CurrentProject = *projectPathFlag
		}
		analyzer.GraphDataMutex.Unlock()
		log.Printf("Initial analysis complete via CLI.")
	}

	mux := http.NewServeMux()
	registerHandlers(mux)
	
	if *runTestsFlag {
		tester.RunBasicTests(mux)
	}

	startServer(mux)
}

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", serveFrontend)
	mux.HandleFunc("/api/analyze", handleAnalyze)
	mux.HandleFunc("/api/graph", handleGraph)
	mux.HandleFunc("/favicon.ico", handleFavicon)
}


func startInteractiveMode() {
	options := []ui.MenuOption{
		{
			Label: "Analyze CodeIgniter 4 Project",
			Handler: func() {
				path := ui.GetInput("Enter project path: ")
				if path != "" {
					log.Printf("Analyzing project: %s", path)
					nodes, edges, err := analyzer.PerformAnalysis(path)
					analyzer.GraphDataMutex.Lock()
					if err != nil {
						log.Printf("Analysis failed: %v", err)
						analyzer.CurrentGraphData = graph.GraphData{Nodes: []graph.Node{{ID: "error", Label: fmt.Sprintf("Error: %v", err), Type: "error"}}}
					} else {
						analyzer.CurrentGraphData = graph.GraphData{Nodes: nodes, Edges: edges}
						analyzer.CurrentProject = path
						fmt.Printf("Analysis complete! Found %d nodes and %d edges.\\n", len(nodes), len(edges))
					}
					analyzer.GraphDataMutex.Unlock()
				}
			},
		},
		{
			Label: "Start Visualization Server",
			Handler: func() {
				fmt.Println("Starting server in background... visit http://localhost:8080")
				mux := http.NewServeMux()
				registerHandlers(mux)
				go startServer(mux)
			},
		},
		{
			Label: "Check Port 8080 Usage",
			Handler: func() {
				netutil.CheckPortUsage(8080)
			},
		},
		{
			Label: "Run Security Check",
			Handler: func() {
				security.CheckRepoSafety()
			},
		},
		{
			Label: "Run Health Checks (Tests)",
			Handler: func() {
				mux := http.NewServeMux()
				registerHandlers(mux)
				tester.RunBasicTests(mux)
			},
		},
	}

	ui.ShowMainMenu(options)
}

func startServer(mux *http.ServeMux) {
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Printf("Server failed (maybe already running?): %v", err)
	}
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNoContent)
}


func serveFrontend(w http.ResponseWriter, r *http.Request) {
	fsys, err := fs.Sub(embeddedFrontend, "frontend")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error getting sub-filesystem for frontend: %v", err)
		return
	}
	http.FileServer(http.FS(fsys)).ServeHTTP(w, r)
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received request to analyze project: %s", req.Path)

	projectPath := req.Path

	go func() {
		analyzer.GraphDataMutex.Lock()
		defer analyzer.GraphDataMutex.Unlock()

		analyzer.CurrentProject = projectPath
		
		nodes, edges, err := analyzer.PerformAnalysis(projectPath)
		if err != nil {
			log.Printf("Analysis failed for %s: %v", projectPath, err)
			analyzer.CurrentGraphData = graph.GraphData{Nodes: []graph.Node{{ID: "error", Label: fmt.Sprintf("Error: %v", err), Type: "error"}}}
		} else {
			analyzer.CurrentGraphData = graph.GraphData{Nodes: nodes, Edges: edges}
			log.Printf("Analysis complete for %s. Found %d nodes and %d edges.", projectPath, len(nodes), len(edges))
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Analysis initiated. Fetch graph data via /api/graph periodically."})
}

func handleGraph(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	analyzer.GraphDataMutex.Lock()
	defer analyzer.GraphDataMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(analyzer.CurrentGraphData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
