package extractor

import (
	"strconv"
	"strings"
)

// --- PHP AST Structures (simplified for initial parsing) ---
type PhpAstNode map[string]interface{}

// GetNameFromNode extracts a simple string name from various AST node types.
func GetNameFromNode(node PhpAstNode) (string, bool) {
	if nameData, ok := node["name"].(map[string]interface{}); ok {
		if name, nameOk := nameData["name"].(string); nameOk {
			return name, true
		}
	}
	if parts, ok := node["parts"].([]interface{}); ok && len(parts) > 0 {
		if len(parts) == 1 {
			if part0, ok := parts[0].(string); ok {
				return part0, true
			}
		}
		if lastPart, ok := parts[len(parts)-1].(string); ok {
			return lastPart, true
		}
	}
	if name, ok := node["name"].(string); ok { // Direct string name
		return name, true
	}
	if value, ok := node["value"].(string); ok { // For nodes like String_
		return value, true
	}
	if value, ok := node["value"].(float64); ok { // For nodes like LNumber
		return strconv.FormatFloat(value, 'f', -1, 64), true
	}
	return "", false
}

// GetFullyQualifiedName extracts a fully qualified name (e.g., "App\Controllers\Home") from an AST Name node.
func GetFullyQualifiedName(node PhpAstNode) (string, bool) {
	if parts, ok := node["parts"].([]interface{}); ok {
		strParts := make([]string, len(parts))
		for i, p := range parts {
			if s, ok := p.(string); ok {
				strParts[i] = s
			} else {
				return "", false
			}
		}
		return strings.Join(strParts, `\`), true // Corrected join separator
	}
	return "", false
}

// ExtractViewCalls recursively finds `view()` function calls within an AST subtree.
func ExtractViewCalls(astNode PhpAstNode) []string {
	var viewNames []string

	if astNode == nil {
		return viewNames
	}

	// Check if it's a function call (Expr_FuncCall)
	if nodeType, ok := astNode["__type"].(string); ok && nodeType == "Expr_FuncCall" {
		if funcNode, fnOk := astNode["name"].(PhpAstNode); fnOk {
			if funcName, nameOk := GetNameFromNode(funcNode); nameOk && strings.ToLower(funcName) == "view" {
				if args, argsOk := astNode["args"].([]interface{}); argsOk && len(args) > 0 {
					if firstArg, argOk := args[0].(PhpAstNode); argOk {
						if valueNode, valOk := firstArg["value"].(PhpAstNode); valOk {
							if viewName, vnOk := GetNameFromNode(valueNode); vnOk {
								viewNames = append(viewNames, viewName)
							}
						}
					}
				}
			}
		}
	}

	// Recursively check children
	for _, key := range []string{"expr", "stmts", "else", "elseifs", "cond", "if", "loop", "body", "args", "value", "items"} {
		if child, ok := astNode[key]; ok {
			if childNode, isNode := child.(PhpAstNode); isNode {
				viewNames = append(viewNames, ExtractViewCalls(childNode)...)
			} else if childNodes, isNodeArray := child.([]interface{}); isNodeArray {
				for _, cn := range childNodes {
					if childNodeItem, isNodeItem := cn.(PhpAstNode); isNodeItem {
						viewNames = append(viewNames, ExtractViewCalls(childNodeItem)...)
					}
				}
			}
		}
	}
	return viewNames
}

// ExtractModelUsage recursively finds `new ModelName()` instantiations within an AST subtree.
func ExtractModelUsage(astNode PhpAstNode) []string {
	var modelNames []string

	if astNode == nil {
		return modelNames
	}

	// Check if it's an object instantiation (Expr_New)
	if nodeType, ok := astNode["__type"].(string); ok && nodeType == "Expr_New" {
		if classNode, classOk := astNode["class"].(PhpAstNode); classOk {
			if modelClassName, nameOk := GetNameFromNode(classNode); nameOk {
				// Simple heuristic: if class name ends with "Model"
				if strings.HasSuffix(modelClassName, "Model") {
					modelNames = append(modelNames, modelClassName)
				}
			}
		}
	}

	// Recursively check children
	for _, key := range []string{"expr", "stmts", "else", "elseifs", "cond", "if", "loop", "body", "args", "value", "items"} {
		if child, ok := astNode[key]; ok {
			if childNode, isNode := child.(PhpAstNode); isNode {
				modelNames = append(modelNames, ExtractModelUsage(childNode)...)
			} else if childNodes, isNodeArray := child.([]interface{}); isNodeArray {
				for _, cn := range childNodes {
					if childNodeItem, isNodeItem := cn.(PhpAstNode); isNodeItem {
						modelNames = append(modelNames, ExtractModelUsage(childNodeItem)...)
					}
				}
			}
		}
	}
	return modelNames
}
