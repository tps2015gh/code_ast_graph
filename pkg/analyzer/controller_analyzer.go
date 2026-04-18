package analyzer

import (
	"fmt"
	"strings"

	"ci4-visualizer/pkg/extractor"
	"ci4-visualizer/pkg/graph"
	"ci4-visualizer/pkg/security" // Import security agent
)

type ControllerAnalyzer struct{}

func (ca *ControllerAnalyzer) Analyze(statements []extractor.PhpAstNode, fileInfo FileInfo, nodes *[]graph.Node, edges *[]graph.Edge) error {
	if !fileInfo.IsController {
		return nil
	}
	
	for _, stmt := range statements {
		if stmtType, ok := stmt["__type"].(string); ok && stmtType == "Stmt_Class" {
			ca.processClass(stmt, fileInfo, nodes, edges)
		}
	}
	return nil
}

func (ca *ControllerAnalyzer) processClass(stmt extractor.PhpAstNode, fileInfo FileInfo, nodes *[]graph.Node, edges *[]graph.Edge) {
	className, ok := extractor.GetNameFromNode(stmt)
	if !ok {
		return
	}

	classID := fmt.Sprintf("class_%s", className)

	// Check for inheritance
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
		Type:     "controller",
		ParentID: fileInfo.NodeID,
	})
	graph.AddEdge(edges, graph.Edge{Source: fileInfo.NodeID, Target: classID, Label: "defines"})

	// Find methods within the class
	if classStmts, csOk := stmt["stmts"].([]interface{}); csOk {
		for _, classStmt := range classStmts {
			if methodStmt, isMethod := classStmt.(extractor.PhpAstNode); isMethod {
				if methodType, mtOk := methodStmt["__type"].(string); mtOk && methodType == "Stmt_ClassMethod" {
					ca.processMethod(methodStmt, className, classID, nodes, edges)
				}
			}
		}
	}
}

func (ca *ControllerAnalyzer) processMethod(methodStmt extractor.PhpAstNode, className, classID string, nodes *[]graph.Node, edges *[]graph.Edge) {
	methodName, ok := extractor.GetNameFromNode(methodStmt)
	if !ok {
		return
	}

	// Check visibility (Public)
	if flags, fOk := methodStmt["flags"].(float64); fOk {
		if int(flags)&1 != 1 { // Not public
			return
		}
	}

	methodID := fmt.Sprintf("method_%s_%s", className, methodName)
	graph.AddNode(nodes, graph.Node{
		ID:       methodID,
		Label:    security.PrivacyScrub(methodName),
		Type:     "method",
		ParentID: classID,
	})
	graph.AddEdge(edges, graph.Edge{Source: classID, Target: methodID, Label: "has_method"})

	// Extract view calls and model usage from method body
	if methodStmts, mbOk := methodStmt["stmts"].([]interface{}); mbOk {
		for _, mStmt := range methodStmts {
			if bodyAstNode, isBodyAst := mStmt.(extractor.PhpAstNode); isBodyAst {
				// View calls
				viewCalls := extractor.ExtractViewCalls(bodyAstNode)
				for _, viewName := range viewCalls {
					viewID := fmt.Sprintf("view_%s", viewName)
					graph.AddNode(nodes, graph.Node{ID: viewID, Label: security.PrivacyScrub(viewName), Type: "view"})
					graph.AddEdge(edges, graph.Edge{Source: methodID, Target: viewID, Label: "renders"})
				}
				// Model usage
				modelUsages := extractor.ExtractModelUsage(bodyAstNode)
				for _, modelName := range modelUsages {
					modelInstanceID := fmt.Sprintf("model_instance_%s", modelName)
					graph.AddNode(nodes, graph.Node{ID: modelInstanceID, Label: security.PrivacyScrub(modelName), Type: "model_instance"})
					graph.AddEdge(edges, graph.Edge{Source: methodID, Target: modelInstanceID, Label: "uses"})
				}
			}
		}
	}
}
