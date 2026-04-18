# TODO for AI: CodeIgniter 4 Visualizer

## 1. Project Status
- [x] Initial implementation of Go backend and Cytoscape.js frontend.
- [x] Integration with `nikic/php-parser` for PHP AST extraction.
- [x] **Fixed:** `ParserFactory::create()` error by updating `parse.php` for PHP-Parser v5 compatibility.
- [x] **Improved:** Error reporting now captures both `stdout` and `stderr` from the PHP parser.
- [x] Basic detection of Controllers, Models, Views, and Routes.
- [x] Project modularized into `pkg/` structure.
- [x] Standalone executable `ci4-visualizer.exe` successfully built.
- [x] Embedded static assets (frontend and parser script).
- [x] Refactored backend to follow HMVC (Modular) principles.
- [x] **Enhanced grouping:** Added CI4 Buckets (🕹️ Controllers, 📦 Models, 🖼️ Views) for better organization.
- [x] **Inspector:** Recursive "drill down" list for folders and classes.
- [x] **UI/UX:** Obsidian-style dark theme with high-performance rendering.
- [x] **Interactive CLI:** Menu-based orchestration with port/process check tool.

## 2. Immediate Tasks (The "Next Steps")

### "Multi-Layer" Swimlane Layout
- [ ] **Implement Swimlane View:**
    - Develop a custom layout logic that forces top-level categories (Controllers, Models, Views) into distinct vertical columns (Swimlanes).
    - This will allow the user to easily "locate" a file by looking at a specific layer of the diagram.
    - Ensure edges between columns (e.g., Controller -> Model) are clear and non-obstructive.

### Advanced Analysis Features
- [ ] **Detect `use` statements:** Correct mapping of model class names even if they are namespaced and alias-ed.
- [ ] **Analyze CI4 Filters and Middlewares:** Link routes to filters.
- [ ] **Database Schema Linkage:** Parse database migrations to link Models to their corresponding database tables.

## 3. Future Enhancements

### Tester Agent Expansion
- [ ] **Automated Test Suite:**
    - Add unit tests for `pkg/extractor` and `pkg/graph`.
    - Implement an integration test that runs the analyzer on a small, dummy CodeIgniter 4 project.

### Frontend Improvement
- [ ] **Search and Filter:** Add a search bar to the frontend to quickly find specific nodes.
- [ ] **Node Inspection:** Implement a detail panel that shows more information (e.g., the PHP code snippet).

---
*Last Updated: 2026-04-19*
