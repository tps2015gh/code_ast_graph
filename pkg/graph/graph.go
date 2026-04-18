package graph

import "fmt"

// Define structures for the graph data
type Node struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Type     string `json:"type"` // e.g., "controller", "model", "view", "route", "method"
	ParentID string `json:"parent,omitempty"` // For hierarchical nodes (compound nodes in Cytoscape)
}

type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label,omitempty"`
}

type GraphData struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Global maps to keep track of added nodes and edges to avoid duplicates
var (
	existingNodes = make(map[string]bool)
	existingEdges = make(map[string]bool)
)

func AddNode(nodes *[]Node, node Node) {
	if _, exists := existingNodes[node.ID]; !exists {
		*nodes = append(*nodes, node)
		existingNodes[node.ID] = true
	}
}

func AddEdge(edges *[]Edge, edge Edge) {
	edgeID := fmt.Sprintf("%s-%s-%s", edge.Source, edge.Label, edge.Target)
	if _, exists := existingEdges[edgeID]; !exists {
		edge.ID = edgeID // Ensure the ID is set!
		*edges = append(*edges, edge)
		existingEdges[edgeID] = true
	}
}

// Reset clears the maps for a new analysis
func Reset() {
	existingNodes = make(map[string]bool)
	existingEdges = make(map[string]bool)
}
