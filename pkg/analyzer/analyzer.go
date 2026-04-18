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
	"ci4-visualizer/pkg/security"
)

var (
	CurrentProject   string
	CurrentGraphData graph.GraphData
	GraphDataMutex   sync.Mutex
)

// PerformAnalysis encapsulates all analysis logic.
func PerformAnalysis(projectPath string) ([]graph.Node, []graph.Edge, error) {
	projectPath = strings.TrimSpace(projectPath)
	
	nodes := []graph.Node{}
	edges := []graph.Edge{}

	graph.Reset()
	graph.AddNode(&nodes, graph.Node{ID: "project", Label: filepath.Base(projectPath), Type: "project"})

	files, err := walkProject(projectPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error walking project %s: %w", projectPath, err)
	}

	log.Printf("Found %d PHP files in %s", len(files), projectPath)

	for i, file := range files {
		fileNodeID := fmt.Sprintf("file_%d", i)
		fileName := filepath.Base(file)
		graph.AddNode(&nodes, graph.Node{ID: fileNodeID, Label: fileName, Type: "file"})
		graph.AddEdge(&edges, graph.Edge{Source: "project", Target: fileNodeID})

		astJSON, err := astparser.ExecutePhpParser(file)
		if err != nil {
			log.Printf("Non-fatal: Skipping %s due to parse error: %v", file, err)
			errorNodeID := fmt.Sprintf("file_error_%d", i)
			graph.AddNode(&nodes, graph.Node{ID: errorNodeID, Label: fmt.Sprintf("Syntax Error: %s", fileName), Type: "error", ParentID: fileNodeID})
			graph.AddEdge(&edges, graph.Edge{Source: fileNodeID, Target: errorNodeID, Label: "has_error"})
			continue
		}

		// Handle flexible AST root
		var statements []extractor.PhpAstNode
		if err := json.Unmarshal(astJSON, &statements); err != nil {
			var singleStmt extractor.PhpAstNode
			if err2 := json.Unmarshal(astJSON, &singleStmt); err2 != nil {
				log.Printf("Failed to unmarshal AST for %s: %v", file, err)
				continue
			}
			statements = []extractor.PhpAstNode{singleStmt}
		}

		// Process statements for the current file
		for _, stmt := range statements {
			if stmt == nil { continue }
			
			stmtType, ok := stmt["__type"].(string)
			if !ok { continue }

			switch stmtType {
			case "Stmt_Class":
				processClass(stmt, fileNodeID, fileName, &nodes, &edges)
			case "Stmt_Function":
				processGlobalFunction(stmt, fileNodeID, fileName, &nodes, &edges)
			case "Stmt_Expression":
				// Check for $routes calls specifically in the routes file
				if strings.HasSuffix(file, filepath.Join("Config", "Routes.php")) {
					processRouteExpression(stmt, fileNodeID, &nodes, &edges)
				}
			}
		}
	}

	return nodes, edges, nil
}

func processClass(stmt extractor.PhpAstNode, fileNodeID, fileName string, nodes *[]graph.Node, edges *[]graph.Edge) {
	className, ok := extractor.GetNameFromNode(stmt)
	if !ok { return }

	classID := fmt.Sprintf("class_%s", className)
	classType := "class"

	if extendsNode, extendsOk := stmt["extends"].(extractor.PhpAstNode); extendsOk && extendsNode != nil {
		if parentClassFQN, parentOk := extractor.GetFullyQualifiedName(extendsNode); parentOk {
			if strings.Contains(parentClassFQN, `CodeIgniter\Controller`) || strings.Contains(parentClassFQN, `BaseController`) {
				classType = "controller"
			} else if strings.Contains(parentClassFQN, `CodeIgniter\Model`) || strings.Contains(parentClassFQN, `BaseModel`) {
				classType = "model"
			}
			externalClassID := "external_class_" + strings.ReplaceAll(parentClassFQN, "", "_")
			graph.AddNode(nodes, graph.Node{ID: externalClassID, Label: security.PrivacyScrub(parentClassFQN), Type: "base_class"})
			graph.AddEdge(edges, graph.Edge{Source: classID, Target: externalClassID, Label: "extends"})
		}
	}
	
	log.Printf("  [Found %s] %s (in %s)", strings.ToUpper(classType), className, fileName)

	graph.AddNode(nodes, graph.Node{ID: classID, Label: security.PrivacyScrub(className), Type: classType, ParentID: fileNodeID})
	graph.AddEdge(edges, graph.Edge{Source: fileNodeID, Target: classID, Label: "defines"})

	if classStmts, csOk := stmt["stmts"].([]interface{}); csOk {
		for _, classStmt := range classStmts {
			if methodStmt, isMethod := classStmt.(extractor.PhpAstNode); isMethod {
				if methodType, mtOk := methodStmt["__type"].(string); mtOk && methodType == "Stmt_ClassMethod" {
					processClassMethod(methodStmt, className, classID, nodes, edges)
				}
			}
		}
	}
}

func processClassMethod(methodStmt extractor.PhpAstNode, className, classID string, nodes *[]graph.Node, edges *[]graph.Edge) {
	methodName, ok := extractor.GetNameFromNode(methodStmt)
	if !ok { return }

	if flags, fOk := methodStmt["flags"].(float64); fOk && int(flags)&1 != 1 {
		return // Not public
	}

	log.Printf("    -> Method: %s", methodName)
	methodID := fmt.Sprintf("method_%s_%s", className, methodName)
	graph.AddNode(nodes, graph.Node{ID: methodID, Label: security.PrivacyScrub(methodName), Type: "method", ParentID: classID})
	graph.AddEdge(edges, graph.Edge{Source: classID, Target: methodID, Label: "has_method"})

	if methodStmts, mbOk := methodStmt["stmts"].([]interface{}); mbOk {
		for _, mStmt := range methodStmts {
			if bodyAstNode, isBodyAst := mStmt.(extractor.PhpAstNode); isBodyAst {
				viewCalls := extractor.ExtractViewCalls(bodyAstNode)
				for _, viewName := range viewCalls {
					viewID := fmt.Sprintf("view_%s", viewName)
					graph.AddNode(nodes, graph.Node{ID: viewID, Label: security.PrivacyScrub(viewName), Type: "view"})
					graph.AddEdge(edges, graph.Edge{Source: methodID, Target: viewID, Label: "renders"})
				}
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

func processGlobalFunction(stmt extractor.PhpAstNode, fileNodeID, fileName string, nodes *[]graph.Node, edges *[]graph.Edge) {
	funcName, ok := extractor.GetNameFromNode(stmt)
	if !ok { return }

	log.Printf("  [Found GLOBAL FUNC] %s (in %s)", funcName, fileName)
	funcID := fmt.Sprintf("func_%s_%s", fileNodeID, funcName)
	graph.AddNode(nodes, graph.Node{ID: funcID, Label: security.PrivacyScrub(funcName), Type: "function", ParentID: fileNodeID})
	graph.AddEdge(edges, graph.Edge{Source: fileNodeID, Target: funcID, Label: "contains_func"})
}

func processRouteExpression(stmt extractor.PhpAstNode, fileNodeID string, nodes *[]graph.Node, edges *[]graph.Edge) {
	if expr, exprOk := stmt["expr"].(extractor.PhpAstNode); exprOk {
		if exprType, etOk := expr["__type"].(string); etOk && exprType == "Expr_MethodCall" {
			if varNode, varOk := expr["var"].(extractor.PhpAstNode); varOk {
				if varName, nameOk := extractor.GetNameFromNode(varNode); nameOk && varName == "routes" {
					if methodNode, methOk := expr["name"].(extractor.PhpAstNode); methOk {
						if routeMethod, rmOk := extractor.GetNameFromNode(methodNode); rmOk {
							if args, argsOk := expr["args"].([]interface{}); argsOk && len(args) >= 2 {
								if patternArg, patOk := args[0].(extractor.PhpAstNode); patOk {
									if routePattern, rpOk := extractor.GetNameFromNode(patternArg["value"].(extractor.PhpAstNode)); rpOk {
										if handlerArg, handOk := args[1].(extractor.PhpAstNode); handOk {
											if handlerString, hsOk := extractor.GetNameFromNode(handlerArg["value"].(extractor.PhpAstNode)); hsOk {
												processRoute(routeMethod, routePattern, handlerString, fileNodeID, nodes, edges)
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

func processRoute(method, pattern, handler, fileNodeID string, nodes *[]graph.Node, edges *[]graph.Edge) {
	if !strings.Contains(handler, "::") { return }

	parts := strings.Split(handler, "::")
	if len(parts) != 2 { return }

	controllerName, methodName := parts[0], parts[1]
	routeID := fmt.Sprintf("route_%s_%s", method, strings.ReplaceAll(pattern, "/", "_"))
	routeLabel := security.PrivacyScrub(fmt.Sprintf("%s %s", strings.ToUpper(method), pattern))
	
	graph.AddNode(nodes, graph.Node{ID: routeID, Label: routeLabel, Type: "route", ParentID: fileNodeID})
	graph.AddEdge(edges, graph.Edge{Source: fileNodeID, Target: routeID, Label: "defines_route"})

	controllerClassID := fmt.Sprintf("class_%s", controllerName)
	controllerMethodID := fmt.Sprintf("method_%s_%s", controllerName, methodName)
	
	graph.AddNode(nodes, graph.Node{ID: controllerClassID, Label: security.PrivacyScrub(controllerName), Type: "controller"})
	graph.AddNode(nodes, graph.Node{ID: controllerMethodID, Label: security.PrivacyScrub(methodName), Type: "method", ParentID: controllerClassID})
	graph.AddEdge(edges, graph.Edge{Source: routeID, Target: controllerMethodID, Label: "calls"})
}

func walkProject(projectPath string) ([]string, error) {
	var phpFiles []string
	appPath := filepath.Join(projectPath, "app")
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("CodeIgniter 'app' directory not found at: %s", appPath)
	}
	
	err := filepath.WalkDir(appPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".php") {
			phpFiles = append(phpFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking project directory %s: %w", projectPath, err)
	}
	return phpFiles, nil
}
