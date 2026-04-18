# TODO for AI: CodeIgniter 4 Visualizer

## 1. Project Status
- [x] Initial implementation of Go backend and Cytoscape.js frontend.
- [x] Integration with `nikic/php-parser` for PHP AST extraction.
- [x] **Fixed:** `ParserFactory::create()` error by updating `parse.php`.
- [x] **Improved:** Error reporting now captures both `stdout` and `stderr`.
- [x] Project modularized into `pkg/` structure.
- [x] **Refactored:** Backend follow HMVC (Modular) principles.
- [x] **Enhanced grouping:** Added CI4 Buckets and physical folder hierarchy.
- [x] **UI/UX:** Obsidian-style dark theme with "Swimlane" layered view.
- [x] **Deep Inspection:** Recursive functional units list in sidebar.

## 2. Immediate Tasks (The "Next Steps")

### AI Static Analysis Agent
- [ ] **Implement Call-Chain Tracking:**
    - Trace the execution path from a Route -> Controller -> Method -> Model.
    - Highlight the "Full Path" in the graph when a node is selected.
- [ ] **Unused Code Detection:**
    - Use the graph to identify "Orphan Nodes" (classes or functions with no incoming edges).
    - Flag these as potential candidates for removal/refactoring.

### Multi-Layer Improvements
- [ ] **Cross-Layer Highlighting:**
    - When clicking a Controller, automatically highlight the Models it uses and the Views it renders in the other swimlanes.

## 3. Future Enhancements

### Advanced Analysis Features
- [ ] **Detect `use` statements:** Correct mapping of model class names even if namespaced.
- [ ] **Analyze CI4 Filters:** Link routes to applied filters.

### Automated Test Suite
- [ ] Add unit tests for `pkg/extractor` and `pkg/graph`.
- [ ] Implement an integration test for a sample CI4 project.

---
*Last Updated: 2026-04-19*
