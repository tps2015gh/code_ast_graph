package analyzer

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"ci4-visualizer/pkg/astparser"
	"ci4-visualizer/pkg/extractor"
	"ci4-visualizer/pkg/graph"
)

var (
	// These will store the result of the last analysis
	CurrentProject   string
	CurrentGraphData graph.GraphData
	GraphDataMutex   sync.Mutex
)

// ProjectAnalyzer orchestrates the analysis of a CI4 project.
type ProjectAnalyzer struct {
	Analyzers []Analyzer
}

func NewProjectAnalyzer() *ProjectAnalyzer {
	return &ProjectAnalyzer{
		Analyzers: []Analyzer{
			&RouteAnalyzer{},
			&ControllerAnalyzer{},
			&ModelAnalyzer{},
		},
	}
}

func PerformAnalysis(projectPath string) ([]graph.Node, []graph.Edge, error) {
	nodes := []graph.Node{}
	edges := []graph.Edge{}

	// Reset global tracking maps for a fresh analysis
	graph.Reset()

	graph.AddNode(&nodes, graph.Node{ID: "project", Label: filepath.Base(projectPath), Type: "project"})

	files, err := walkProject(projectPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error walking project %s: %w", projectPath, err)
	}

	log.Printf("Found %d PHP files in %s", len(files), projectPath)

	pa := NewProjectAnalyzer()

	for i, file := range files {
		fileNodeID := fmt.Sprintf("file_%d", i)
		fileInfo := pa.getFileInfo(file, projectPath, fileNodeID)
		
		graph.AddNode(&nodes, graph.Node{ID: fileNodeID, Label: fileInfo.Name, Type: "file"})
		graph.AddEdge(&edges, graph.Edge{Source: "project", Target: fileNodeID})

		astJSON, err := astparser.ExecutePhpParser(file)
		if err != nil {
			log.Printf("Non-fatal: Skipping %s due to parse error: %v", file, err)
			
			// Add an error node to the graph so the user sees it in the diagram
			errorNodeID := fmt.Sprintf("file_error_%d", i)
			graph.AddNode(&nodes, graph.Node{
				ID:       errorNodeID,
				Label:    fmt.Sprintf("Syntax Error: %s", fileInfo.Name),
				Type:     "error",
				ParentID: fileNodeID,
			})
			graph.AddEdge(&edges, graph.Edge{Source: fileNodeID, Target: errorNodeID, Label: "has_error"})
			continue
		}

		var statements []extractor.PhpAstNode
		if err := json.Unmarshal(astJSON, &statements); err != nil {
			log.Printf("Failed to unmarshal AST for %s: %v", file, err)
			continue
		}

		for _, a := range pa.Analyzers {
			if err := a.Analyze(statements, fileInfo, &nodes, &edges); err != nil {
				log.Printf("Analyzer error for %s: %v", file, err)
			}
		}
	}

	return nodes, edges, nil
}

func (pa *ProjectAnalyzer) getFileInfo(path, projectPath, nodeID string) FileInfo {
	name := filepath.Base(path)
	info := FileInfo{
		Path:   path,
		Name:   name,
		NodeID: nodeID,
	}

	relPath, _ := filepath.Rel(projectPath, path)
	relPath = filepath.ToSlash(relPath)

	if strings.HasSuffix(relPath, "Config/Routes.php") {
		info.IsRoute = true
	} else if strings.Contains(relPath, "app/Controllers") {
		info.IsController = true
	} else if strings.Contains(relPath, "app/Models") {
		info.IsModel = true
	}

	return info
}

func walkProject(projectPath string) ([]string, error) {
	var phpFiles []string
	appPath := filepath.Join(projectPath, "app")
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("CodeIgniter 'app' directory not found at: %s", appPath)
	}
	
	err := filepath.WalkDir(appPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v", path, err)
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".php") {
			relPath, _ := filepath.Rel(projectPath, path)
			relPath = filepath.ToSlash(relPath)
			
			if strings.HasSuffix(relPath, "Config/Routes.php") || 
			   strings.Contains(relPath, "app/Controllers") || 
			   strings.Contains(relPath, "app/Models") {
				phpFiles = append(phpFiles, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking project directory %s: %w", projectPath, err)
	}
	return phpFiles, nil
}
