package cmd

import (
	"flag"
	"fmt"
	"testing"
)

func TestFindCmdVerb(t *testing.T) {
	tests := []struct {
		name    string
		verb    string
		matches int
	}{
		{"Valid init verb", "init", 1},
		{"Invalid init verb", "innit", 0},
		{"Valid install verb", "install", 1},
		{"Invalid install verb", "iinstall", 0},
		{"Valid call verb", "call", 1},
		{"Invalid call verb", "calI", 0},
		{"Valid export verb", "export", 1},
		{"Invalid export verb", "exprt", 0},
		{"Valid quit verb", "quit", 1},
		{"Invalid quit verb", "quitt", 0},
		{"Unknown verb", "unknown", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := tt.verb
			matches := FindCmdVerb(command)
			if len(matches) != tt.matches {
				t.Errorf("No matches for command: %v found no match", command)
			}
		})
	}

}

// func TestValidateCommandV2(t *testing.T) {
// 	args := []string{"arg0", "arg1", "arg2"}
// 	fs := flag.NewFlagSet("test", flag.ExitOnError)
// 	fs.Parse(args)
// 	flag.CommandLine = fs
// 	err := ValidateCommand("init")
// 	if err != nil {
// 		t.Errorf("Expected no error, got: %v", err)
// 	}
// }

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		numArgs int
		wantErr bool
	}{
		{"Valid init command", "init", 3, false},
		{"Invalid init command", "init", 2, true},
		{"Valid install command", "install", 2, false},
		{"Invalid install command", "install", 3, true},
		{"Valid call command", "call", 1, false},
		{"Invalid call command", "call", 0, true},
		{"Valid export command", "export", 1, false},
		{"Valid quit command", "quit", 1, false},
		{"Invalid quit command", "quit", 2, true},
		{"Unknown command", "unknown", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up flags for the test
			args := make([]string, 0)
			fs := flag.NewFlagSet(tt.command, flag.ExitOnError)
			for i := 0; i < tt.numArgs; i++ {
				args = append(args, fmt.Sprintf("arg%v", i))
			}
			fs.Parse(args)
			flag.CommandLine = fs

			err := ValidateCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
