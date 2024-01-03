package cmd

import (
	"flag"
	"fmt"
	"strings"
)

func FindCmdVerb(prefix string) []string {
	commands := []string{"init", "install", "call", "export", "quit"}
	var matches []string

	for _, cmd := range commands {
		if strings.HasPrefix(cmd, prefix) {
			matches = append(matches, cmd)
		}
	}

	return matches
}
func ValidateCommand(command string) error {
	numArgsRecv := flag.NArg()
	fmt.Printf("Validating command, number of args from flag: %v\n", numArgsRecv)
	switch command {
	case "init":
		if flag.NArg() != 3 { //this takes only dir to init env and nameOfEnv init . myfirstEnv
			fmt.Printf("Command: %v takes %v arguments, received %v\n", command, 3, numArgsRecv)
			return fmt.Errorf("wrong number of arguments")
		}
	case "install":
		if flag.NArg() != 2 { //this takes just the toolName e.g bioenv install samtools
			return fmt.Errorf("wrong number of arguments")
		}
	case "call", "export":
		if flag.NArg() < 1 {
			return fmt.Errorf("wrong number of arguments")
		}
	case "quit":
		if flag.NArg() != 1 {
			return fmt.Errorf("wrong number of arguments")
		}
	default:
		return fmt.Errorf("unknown command - %s", flag.Arg(0))
	}
	return nil
}
