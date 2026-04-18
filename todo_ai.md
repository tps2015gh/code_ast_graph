# TODO for AI: CodeIgniter 4 Visualizer

## 1. Project Status
- [x] Initial implementation of Go backend and Cytoscape.js frontend.
- [x] Integration with `nikic/php-parser` for PHP AST extraction.
- [x] **Fixed:** `ParserFactory::create()` error by updating `parse.php` for PHP-Parser v5 compatibility.
- [x] **Improved:** Error reporting now captures both `stdout` and `stderr` from the PHP parser, providing clear syntax error messages.
- [x] Basic detection of Controllers, Models, Views, and Routes.
- [x] Project modularized into `pkg/` structure for better maintainability.
- [x] Standalone executable `ci4-visualizer.exe` successfully built on Windows.
- [x] Embedded static assets (frontend and parser script) for single-binary distribution.
- [x] Refactored backend to follow HMVC (Modular) principles.
- [x] Enhanced graph data structure (hierarchical nodes).
- [x] Tester Agent Integration (basic health checks).

## 2. Immediate Tasks (The "Next Steps")

### Tester Agent Expansion
- [ ] **Expand `pkg/tester` package:**
    - Implement checks that verify the backend API endpoints return valid JSON.
    - Add a test that verifies the analysis logic correctly identifies specific edges (e.g., "renders", "calls", "uses") for a known sample code snippet.
- [ ] **Automated Test Suite:**
    - Add unit tests for `pkg/extractor` and `pkg/graph`.
    - Implement an integration test that runs the analyzer on a small, dummy CodeIgniter 4 project and verifies the output graph structure.

### Advanced Analysis Features
- [ ] **Detect `use` statements:** Correct mapping of model class names even if they are namespaced and alias-ed.
- [ ] **Analyze CI4 Filters and Middlewares:** Link routes to filters.
- [ ] **Database Schema Linkage:** Parse database migrations to link Models to their corresponding database tables.
- [ ] **Detect HMVC Modules:** Specifically detect and visualize CodeIgniter 4 modules located in `app/Modules` (if used in the target project).

## 3. Future Enhancements

### Frontend Improvement
- [ ] **Search and Filter:** Add a search bar to the frontend to quickly find specific nodes (e.g., a specific controller or view).
- [ ] **Node Inspection:** Implement a detail panel that shows more information when a node is clicked (e.g., the PHP code snippet, docblocks).
- [ ] **Improved Layouts:** Experiment with different Cytoscape.js layouts (e.g., `dagre`, `klay`) to better represent hierarchical structures.

---
*Last Updated: 2026-04-19*
