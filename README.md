# Code AST Visualizer

A professional-grade, standalone web application that transforms PHP CodeIgniter 4 projects into interactive, AST-driven architectural maps. Built with Go and powered by `nikic/php-parser` v5, it provides high-fidelity analysis of code relationships and folder hierarchies.

## Features
- **Hierarchical Tree-View:** Auto-detects CI4 folder structures and groups nodes into logical "Buckets" (Controllers, Models, Views, Config).
- **Greedy AST Analysis:** Scans the entire project using inheritance-based detection to identify logic anywhere in the source.
- **AI-Ready Static Analysis:** Leverages the Abstract Syntax Tree to track code dependencies, function calls, and data flow with 100% semantic precision.
- **Deep Inspection Sidebar:** Click any node to "drill down" and see all nested classes, methods, global functions, and routes.
- **Integrated Security Agent:** Automated `.gitignore` management and privacy scrubbing to protect sensitive source code data.
- **Interactive CLI Orchestration:** A menu-driven interface for analysis, network/port management, and automated health checks.

## The Team

| Name | Role | Job Tasks |
| :--- | :--- | :--- |
| **The Architect** (User) | Visionary Lead, QA Director & UX Designer | Established the HMVC architectural foundation, defined the "Greedy Analysis" requirement, and meticulously identified rendering/performance bottlenecks. |
| **The Builder** (Gemini CLI) | Principal Engineer & Security Guardian | Engineered the Go modular backend, implemented the AST-to-Graph transformation, developed the Security/Tester agents, and optimized the Cytoscape.js implementation. |

## Why This Project is a "Goldmine" for AI
For an AI like me, this tool is more than a visualizer—it's a **contextual accelerator**.
1.  **Semantic Navigation:** Unlike raw text search, the AST allows an AI to follow the "living logic" of the code. I can move from a Route to a Method to a Model instantly, knowing exactly what each node is.
2.  **Context Compression:** A 480-node graph is a highly compressed "map" of your project. It allows an AI to understand the entire architecture in seconds, which would normally require reading thousands of lines of code.
3.  **Precise Impact Prediction:** If a developer asks "What happens if I delete this function?", an AI using this tool doesn't guess—it follows the graph edges to identify every single dependency.

## Advice: How to turn this into a full "AI Tool"
To move from a "Visualizer" to a truly "Proactive AI Tool," we should implement:
- **Symbolic Data Flow:** Track variables (not just functions). This would let an AI detect "What database columns does this specific form input touch?"
- **NLQ (Natural Language Query):** Implement an interface where a user asks: *"Show me all Controllers that use the User model but don't have a login check."* 
- **Automated Refactoring Proposals:** Use the graph to identify "God Classes" (nodes with too many edges) and have the AI propose a service-oriented split.
- **Context-Aware Prompts:** Export selected parts of the graph/AST directly into a prompt format that can be used by LLMs for automated bug fixing.

## AI's Opinion on the Project
By using an AST instead of raw text, we've given this tool a "semantic brain." This project represents a shift from simple visualization to **AI-assisted static analysis**. Because I can "track" the tree, I can understand the intent of the code, find broken connections, and identify security vulnerabilities that a human developer might miss in a large codebase.

## Layout Analysis: Why Dagre?
For CodeIgniter 4, the **Dagre (Top-Down)** layout is the "Gold Standard." CI4’s strict adherence to Namespaces and MVC hierarchies means that a vertical flow most naturally represents the relationship between a Request (Route), the Logic (Controller), and the Result (View).

## License
Distributed under the MIT License. See `LICENSE` for more information.
