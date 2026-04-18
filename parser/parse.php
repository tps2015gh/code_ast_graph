<?php

require __DIR__ . '/vendor/autoload.php';

use PhpParser\Error;
use PhpParser\ParserFactory;
use PhpParser\Node;

if (!isset($argv[1])) {
    echo "Usage: php parse.php <file_path>\n";
    exit(1);
}

$filePath = $argv[1];

if (!file_exists($filePath)) {
    echo "Error: File not found: " . $filePath . "\n";
    exit(1);
}

$code = file_get_contents($filePath);

$parser = (new ParserFactory)->createForNewestSupportedVersion();

try {
    $ast = $parser->parse($code);

    // Transform AST nodes to include the node type for Go
    $serializableAst = transformNodes($ast);

    echo json_encode($serializableAst, JSON_PRETTY_PRINT);
} catch (Error $error) {
    echo "Parse error: {$error->getMessage()}\n";
    exit(1);
}

/**
 * Recursively transforms PHP-Parser nodes into arrays with a __type field.
 */
function transformNodes($nodes) {
    if (!is_array($nodes) && !($nodes instanceof Node)) {
        return $nodes;
    }

    if ($nodes instanceof Node) {
        $result = ['__type' => $nodes->getType()];
        foreach ($nodes->getSubNodeNames() as $name) {
            $subNode = $nodes->$name;
            $result[$name] = transformNodes($subNode);
        }
        return $result;
    }

    $result = [];
    foreach ($nodes as $key => $value) {
        $result[$key] = transformNodes($value);
    }
    return $result;
}
?>