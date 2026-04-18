package netutil

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// CheckPortUsage finds and displays the process using a specific port.
func CheckPortUsage(port int) {
	var cmd *exec.Cmd
	portStr := fmt.Sprintf(":%d", port)

	fmt.Printf("Checking for processes using port %d...\\n", port)

	if runtime.GOOS == "windows" {
		// Using double quotes for the findstr pattern
		cmdStr := fmt.Sprintf("netstat -ano | findstr \"%s\"", portStr)
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("lsof -i %s", portStr))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Could not execute command: %v\\n", err)
		return
	}

	outStr := strings.TrimSpace(string(output))
	if len(outStr) == 0 {
		fmt.Printf("No process found using port %d.\\n", port)
	} else {
		fmt.Printf("\\n--- Process(es) using port %d ---\\n", port)
		fmt.Println(outStr)
		fmt.Println("---------------------------------")
		if runtime.GOOS == "windows" {
			fmt.Println("Find the PID in the last column. Use 'taskkill /F /PID <PID_NUMBER>' to terminate it.")
		}
	}
}
