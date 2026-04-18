package analyzer

import (
	"fmt"
	"strings"

	"ci4-visualizer/pkg/extractor"
	"ci4-visualizer/pkg/graph"
	"ci4-visualizer/pkg/security" // Import security agent
)

type RouteAnalyzer struct{}

func (ra *RouteAnalyzer) Analyze(statements []extractor.PhpAstNode, fileInfo FileInfo, nodes *[]graph.Node, edges *[]graph.Edge) error {
	if !fileInfo.IsRoute {
		return nil
	}

	for _, stmt := range statements {
		// Look for $routes->method(...) calls
		if stmtType, ok := stmt["__type"].(string); ok && stmtType == "Stmt_Expression" {
			if expr, exprOk := stmt["expr"].(extractor.PhpAstNode); exprOk {
				if exprType, etOk := expr["__type"].(string); etOk && exprType == "Expr_MethodCall" {
					if varNode, varOk := expr["var"].(extractor.PhpAstNode); varOk {
						if varType, vtOk := varNode["__type"].(string); vtOk && varType == "Expr_Variable" {
							if varName, nameOk := extractor.GetNameFromNode(varNode); nameOk && varName == "routes" {
								if methodNode, methOk := expr["name"].(extractor.PhpAstNode); methOk {
									if routeMethod, rmOk := extractor.GetNameFromNode(methodNode); rmOk {
										if args, argsOk := expr["args"].([]interface{}); argsOk && len(args) >= 2 {
											if patternArg, patOk := args[0].(extractor.PhpAstNode); patOk {
												if routePattern, rpOk := extractor.GetNameFromNode(patternArg["value"].(extractor.PhpAstNode)); rpOk {
													if handlerArg, handOk := args[1].(extractor.PhpAstNode); handOk {
														if handlerString, hsOk := extractor.GetNameFromNode(handlerArg["value"].(extractor.PhpAstNode)); hsOk {
															ra.processRoute(routeMethod, routePattern, handlerString, fileInfo, nodes, edges)
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (ra *RouteAnalyzer) processRoute(method, pattern, handler string, fileInfo FileInfo, nodes *[]graph.Node, edges *[]graph.Edge) {
	if !strings.Contains(handler, "::") {
		return
	}

	parts := strings.Split(handler, "::")
	if len(parts) != 2 {
		return
	}

	controllerName := parts[0]
	methodName := parts[1]

	routeID := fmt.Sprintf("route_%s_%s", method, strings.ReplaceAll(pattern, "/", "_"))
	
	// Apply privacy scrub to the label
	routeLabel := security.PrivacyScrub(fmt.Sprintf("%s %s", strings.ToUpper(method), pattern))
	
	graph.AddNode(nodes, graph.Node{
		ID:    routeID,
		Label: routeLabel,
		Type:  "route",
		ParentID: fileInfo.NodeID,
	})
	graph.AddEdge(edges, graph.Edge{Source: fileInfo.NodeID, Target: routeID, Label: "defines_route"})

	controllerClassID := fmt.Sprintf("class_%s", controllerName)
	controllerMethodID := fmt.Sprintf("method_%s_%s", controllerName, methodName)
	
	graph.AddNode(nodes, graph.Node{ID: controllerClassID, Label: security.PrivacyScrub(controllerName), Type: "controller"})
	graph.AddNode(nodes, graph.Node{ID: controllerMethodID, Label: security.PrivacyScrub(methodName), Type: "method", ParentID: controllerClassID})

	graph.AddEdge(edges, graph.Edge{Source: routeID, Target: controllerMethodID, Label: "calls"})
}
