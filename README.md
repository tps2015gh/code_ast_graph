# Code AST Visualizer

A standalone, modular web application that parses PHP CodeIgniter 4 projects to generate an interactive, node-based knowledge graph. It leverages the power of Go for the backend engine and `nikic/php-parser` (PHP) for high-fidelity Abstract Syntax Tree (AST) analysis.

## Features
- **Obsidian-Style Dark Theme:** A modern, clean interface designed for developers.
- **Hierarchical Visualization:** Defaults to a top-down `dagre` layout for clear architectural mapping.
- **Deep Inspection:** Interactive "Node Details" panel showing all methods, global functions, and routes within a file or class.
- **HMVC Orchestration:** A modular backend with specialized analyzers for Routes, Controllers, Models, and Functions.
- **Security & Privacy:** Integrated "Security Agent" with auto-generated `.gitignore` and privacy scrubbing for sensitive labels.
- **Interactive CLI:** A user-friendly menu system for analysis, server management, and health checks.
- **Single Binary:** Everything (frontend, parser, logic) is bundled into a single Go executable.

## The Team

| Name | Role | Job Tasks |
| :--- | :--- | :--- |
| **The Architect** (User) | Lead Visionary, PM & QA Lead | Defined core requirements, established HMVC modular architecture, provided expert CI4 guidance, identified critical bugs, and drove UX/performance refinements. |
| **The Builder** (Gemini CLI) | Full-Stack Engineer & Automation Specialist | Implemented Go backend, integrated PHP-Parser v5, developed interactive CLI, created Security & Tester agents, optimized graph rendering, and managed modularization. |

## AI's Opinion
This project has evolved into a robust, developer-centric tool. The shift to a modular HMVC-inspired backend proves that complex code analysis can be both fast and maintainable. By implementing security and privacy layers early, we've ensured that "Code AST Visualizer" is safe for professional use. The addition of the deep details panel and global function scanning truly bridges the gap between a high-level architectural overview and low-level code inspection. It's a testament to the power of human-AI collaboration in solving real-world engineering challenges.

## License
Distributed under the MIT License. See `LICENSE` for more information.
