package main

import (
	"fmt"
	"gitlab-automation/modules/gitlabpermissions"
	"os"
)

func main() {
	// All commands are formatted:
	//  gitlab-automation <command> <args for the command...>

	// Collect arguments to the program
	args := os.Args

	// Args should have at least 2 args... the bin and the command
	if len(args) == 1 {
		commandHelp()
		os.Exit(0)
	}

	if len(args) > 1 {
		switch args[1] {
		case "help":
			commandHelp()
		case "permissions":
			gitlabpermissions.CommandGitlabPermissions()
		default:
			unknownCommand(args[1])
		}
	}
}

func commandHelp() {
	fmt.Println("Usage: gitlab-automation <command>")
	fmt.Println("Available command:")
	fmt.Println(" - help : Show this help")
	fmt.Println(" - permissions : List the users who have access to a group")
}

func unknownCommand(command string) {
	fmt.Printf("Unknown command: %s", command)
	fmt.Println()
	commandHelp()
}
