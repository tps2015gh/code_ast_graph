package analyzer

import (
	"ci4-visualizer/pkg/extractor"
	"ci4-visualizer/pkg/graph"
)

// Analyzer defines the interface for specialized project analyzers.
type Analyzer interface {
	Analyze(statements []extractor.PhpAstNode, fileInfo FileInfo, nodes *[]graph.Node, edges *[]graph.Edge) error
}

// FileInfo contains metadata about the file being analyzed.
type FileInfo struct {
	Path       string
	Name       string
	NodeID     string
	IsRoute    bool
	IsController bool
	IsModel    bool
}
