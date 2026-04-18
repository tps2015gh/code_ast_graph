# CodeIgniter 4 AST Visualizer

A standalone, modular web application that parses PHP CodeIgniter 4 projects to generate an interactive, node-based knowledge graph. It leverages the power of Go for the backend engine and `nikic/php-parser` (PHP) for high-fidelity Abstract Syntax Tree (AST) analysis.

## Features
- **Hierarchical Visualization:** Uses Cytoscape.js to show relationships between Controllers, Models, Views, and Routes.
- **AST-Based Analysis:** Goes beyond simple regex to understand the actual structure of the PHP code.
- **HMVC Architecture:** Built with a modular analyzer system (Route, Controller, and Model analyzers).
- **Standalone:** Bundles all assets (PHP parser, frontend) into a single Go executable.

## The Team

| Name | Role | Job Tasks |
| :--- | :--- | :--- |
| **The Architect** (User) | Lead Visionary & Project Manager | Defined the core requirements, established the HMVC modular architectural direction, and provided expert CI4 guidance. |
| **The Builder** (Gemini CLI) | Implementation Specialist | Developed the Go backend, integrated the PHP-based AST parser, implemented the Cytoscape.js frontend, and handled modularization. |

## AI's Opinion
This project is a perfect example of **pragmatic engineering**. By bridging the Go and PHP ecosystems, we avoid reinventing the wheel (parsing PHP is complex!) while gaining Go's speed and deployment simplicity. Using AST analysis for visualization transforms a codebase from a collection of text files into a navigable map, which is invaluable for onboarding, auditing, and refactoring complex CI4 applications. The addition of hierarchical nodes makes the "Obsidian-style" graph truly reflect the logical nesting of a professional web framework.

## License
Distributed under the MIT License. See `LICENSE` for more information.
