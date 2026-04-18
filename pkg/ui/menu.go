package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// MenuOption represents a single item in the interactive menu.
type MenuOption struct {
	Label   string
	Handler func()
}

// ShowMainMenu displays an interactive CLI menu to the user.
func ShowMainMenu(options []MenuOption) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- CodeIgniter 4 AST Visualizer ---")
		for i, opt := range options {
			fmt.Printf("%d. %s\n", i+1, opt.Label)
		}
		fmt.Printf("%d. Exit\n", len(options)+1)
		fmt.Print("\nSelect an option: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == fmt.Sprintf("%d", len(options)+1) || strings.ToLower(input) == "exit" || strings.ToLower(input) == "q" {
			fmt.Println("Exiting...")
			return
		}

		var choice int
		fmt.Sscanf(input, "%d", &choice)

		if choice > 0 && choice <= len(options) {
			options[choice-1].Handler()
		} else {
			fmt.Println("Invalid selection, please try again.")
		}
	}
}

// GetInput prompts the user for string input with a specific message.
func GetInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
