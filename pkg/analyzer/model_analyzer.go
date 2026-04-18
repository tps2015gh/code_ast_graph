package analyzer

import (
	"fmt"
	"strings"

	"ci4-visualizer/pkg/extractor"
	"ci4-visualizer/pkg/graph"
	"ci4-visualizer/pkg/security" // Import security agent
)

type ModelAnalyzer struct{}

func (ma *ModelAnalyzer) Analyze(statements []extractor.PhpAstNode, fileInfo FileInfo, nodes *[]graph.Node, edges *[]graph.Edge) error {
	if !fileInfo.IsModel {
		return nil
	}
	
	for _, stmt := range statements {
		if stmtType, ok := stmt["__type"].(string); ok && stmtType == "Stmt_Class" {
			ma.processClass(stmt, fileInfo, nodes, edges)
		}
	}
	return nil
}

func (ma *ModelAnalyzer) processClass(stmt extractor.PhpAstNode, fileInfo FileInfo, nodes *[]graph.Node, edges *[]graph.Edge) {
	className, ok := extractor.GetNameFromNode(stmt)
	if !ok {
		return
	}

	classID := fmt.Sprintf("class_%s", className)
	
	// Ensure external parent node and inheritance
	if extendsNode, extendsOk := stmt["extends"].(extractor.PhpAstNode); extendsOk && extendsNode != nil {
		if parentClassFQN, parentOk := extractor.GetFullyQualifiedName(extendsNode); parentOk {
			externalClassID := "external_class_" + strings.ReplaceAll(parentClassFQN, "\\", "_")
			graph.AddNode(nodes, graph.Node{ID: externalClassID, Label: security.PrivacyScrub(parentClassFQN), Type: "base_class"})
			graph.AddEdge(edges, graph.Edge{Source: classID, Target: externalClassID, Label: "extends"})
		}
	}

	graph.AddNode(nodes, graph.Node{
		ID:       classID,
		Label:    security.PrivacyScrub(className),
		Type:     "model",
		ParentID: fileInfo.NodeID,
	})
	graph.AddEdge(edges, graph.Edge{Source: fileInfo.NodeID, Target: classID, Label: "defines"})
}
