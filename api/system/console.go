package system

import (
	"fmt"
	"os"
)

// ConsoleMenu implements the TUI selection logic
type ConsoleMenu struct {
	Options []string
}

// ShowMenu displays the interactive console interface
func (m ConsoleMenu) ShowMenu() {
	fmt.Println("\n--- ZeroWall Console Setup ---")
	for i, opt := range m.Options {
		fmt.Printf("%d) %s\n", i+1, opt)
	}
	fmt.Printf("%d) Logout\n", len(m.Options)+1)
	fmt.Print("\nEnter selection: ")
}

// HandleSelection routes console inputs
func (m ConsoleMenu) HandleSelection(choice int) {
	switch choice {
	case 1:
		fmt.Println("Assigning Interfaces...")
	case 2:
		fmt.Println("Configuring IP Addresses...")
	case 3:
		fmt.Println("Resetting Admin Password...")
	case 4:
		fmt.Println("Rebooting System...")
		// exec.Command("reboot").Run()
	default:
		fmt.Println("Invalid Selection.")
	}
}
