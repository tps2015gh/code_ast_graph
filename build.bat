@echo off
echo Building CI4 AST Visualizer...
go build -o ci4-visualizer.exe
if %ERRORLEVEL% EQU 0 (
    echo.
    echo ---------------------------------------
    echo Build SUCCESSFUL: ci4-visualizer.exe
    echo ---------------------------------------
) else (
    echo.
    echo #######################################
    echo Build FAILED! Check the errors above.
    echo #######################################
)
pause
