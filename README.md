# Code AST Visualizer

A professional-grade, standalone web application that transforms PHP CodeIgniter 4 projects into interactive, AST-driven architectural maps. Built with Go and powered by `nikic/php-parser` v5, it provides high-fidelity analysis of code relationships and folder hierarchies.

## Features
- **Hierarchical Tree-View:** Auto-detects CI4 folder structures and groups nodes into logical "Buckets" (Controllers, Models, Views, Config).
- **Greedy AST Analysis:** Scans the entire project using inheritance-based detection to identify Controllers and Models anywhere in the source, not just specific folders.
- **Obsidion-Style Dark Theme:** A premium developer experience with high-contrast, color-coded visualization.
- **Deep Inspection Sidebar:** Click any node to "drill down" and see all nested classes, methods, global functions, and routes.
- **Integrated Security Agent:** Automated `.gitignore` management and privacy scrubbing to protect sensitive source code data.
- **Interactive CLI Orchestration:** A menu-driven interface for analysis, network/port management, and automated health checks.

## The Team

| Name | Role | Job Tasks |
| :--- | :--- | :--- |
| **The Architect** (User) | Visionary Lead, QA Director & UX Designer | Established the HMVC architectural foundation, defined the "Greedy Analysis" requirement, and meticulously identified rendering/performance bottlenecks. |
| **The Builder** (Gemini CLI) | Principal Engineer & Security Guardian | Engineered the Go modular backend, implemented the AST-to-Graph transformation, developed the Security/Tester agents, and optimized the Cytoscape.js implementation. |

## AI's Opinion on the Project
This project represents a **paradigm shift** in how developers interact with legacy and modern CI4 codebases. By moving from text-based navigation to a visual, hierarchical map, we've reduced the "cognitive load" required to understand complex architectures. The implementation of "Greedy Analysis" was a turning point; it transformed the tool from a directory viewer into a true **static analysis engine** that understands code intent through inheritance rather than just file paths.

## Layout Analysis: Why Dagre?
For CodeIgniter 4, the **Dagre (Top-Down)** layout is the "Gold Standard." CI4’s strict adherence to Namespaces and MVC hierarchies means that a vertical flow most naturally represents the relationship between a Request (Route), the Logic (Controller), and the Result (View). It transforms a flat list of files into a navigable organizational chart.

### Future Layout Considerations:
- **Breadth-First (Dependency Layers):** To be added for tracking "Service Layers" and dependency injection chains.
- **Klay (Advanced Hierarchy):** A candidate for even tighter, more professional node packing in very large projects.

## License
Distributed under the MIT License. See `LICENSE` for more information.
