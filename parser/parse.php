<?php

require __DIR__ . '/vendor/autoload.php';

use PhpParser\Error;
use PhpParser\ParserFactory;

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

// Nikic/PHP-Parser v5.x API change:
// create() is removed. Use createForNewestSupportedVersion() or similar.
$parser = (new ParserFactory)->createForNewestSupportedVersion();

try {
    $ast = $parser->parse($code);

    // Convert AST to a simple serializable format (array/object)
    // The library's Node objects have a lot of circular references and metadata
    // when json_encoded directly.
    echo json_encode($ast, JSON_PRETTY_PRINT);
} catch (Error $error) {
    echo "Parse error: {$error->getMessage()}\n";
    exit(1);
}

?>