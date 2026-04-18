package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"ci4-visualizer/pkg/analyzer"
	"ci4-visualizer/pkg/astparser"
	"ci4-visualizer/pkg/graph"
	"ci4-visualizer/pkg/security" // Import security agent
	"ci4-visualizer/pkg/tester"
)

//go:embed frontend/*
var embeddedFrontend embed.FS

//go:embed parser/*
var embeddedParser embed.FS

func main() {
	// Parse CLI flags
	projectPathFlag := flag.String("path", "", "Path to the CodeIgniter 4 project to analyze on startup")
	phpExePathFlag := flag.String("php", "C:/xampp_v8_1_25/php/php.exe", "Path to the PHP executable (e.g., C:/xampp/php/php.exe)")
	runTestsFlag := flag.Bool("test", false, "Run basic health checks on startup")
	flag.Parse()

	// Initial security check
	security.CheckRepoSafety()

	// Get a sub-filesystem for the parser to avoid "parser/" prefix in paths
	parserFS, err := fs.Sub(embeddedParser, "parser")
	if err != nil {
		log.Fatalf("Failed to create sub-filesystem for parser: %v", err)
	}

	// Initialize the AST parser with the PHP executable path and the scoped parser FS
	err = astparser.Init(*phpExePathFlag, parserFS)
	if err != nil {
		log.Fatalf("Failed to initialize AST parser: %v", err)
	}
	defer astparser.Cleanup()

	// If a project path is provided via CLI, trigger analysis immediately
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
	mux.HandleFunc("/", serveFrontend)
	mux.HandleFunc("/api/analyze", handleAnalyze)
	mux.HandleFunc("/api/graph", handleGraph)

	if *runTestsFlag {
		tester.RunBasicTests(mux)
	}

	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, mux))
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
