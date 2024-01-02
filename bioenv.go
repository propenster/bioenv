package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bioenv/cmd"
	"bioenv/vienv"
)

const (
	usage = `Usage:
	bioenv init <directory> <name>    Initialize a new project in the specified directory with the given name
	bioenv install <toolName>         Install the specified bio tool
	bioenv call <toolName> [args...]  Call the specified bio tool with optional arguments
	bioenv export <argument>          Export data using the specified argument
	bioenv quit                       Quit the application
	bioenv --version                  Print the application version
  
  Commands:
	init      Initialize a new project
	install   Install a bio tool
	call      Call a bio tool
	export    Export data
	quit      Quit the application
	--version Print the application version
  
  Arguments:
	<directory> Directory path for init command
	<name>      Project name for init command
	<toolName>  Name of the bio tool for install and call commands
	[args...]   Optional arguments for the call command
	<argument>  Argument for the export command
  
  Examples:
	bioenv init /path/to/project MyFirstBioVirtualEnvironment
	bioenv install GATK
	bioenv call GATK --arg1 value1 --arg2 value2
	bioenv export firstenv.json
	bioenv quit
	bioenv --version
  `
	version = "0.0.1"
)

var venv vienv.VirtualEnvironment

func main() {
	fmt.Println("Welcome to BioEnv")
	args := os.Args
	fmt.Printf("Arguments received: %v\n", args)

	log.SetFlags(0) // Don't prefix with time
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "show version and exit")
	flag.Usage = func() {
		// name := path.Base(os.Args[0])
		// fmt.Printf(usage, name)
		fmt.Println(usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Println("Invalid number of arguments NArg() == zero")
		log.Fatalf("error: wrong number of arguments")
	}

	matches := cmd.FindCmdVerb(flag.Arg(0))
	switch len(matches) {
	case 0:
		log.Fatalf("error: unknown command - %q\n", flag.Arg(0))
	case 1: /* nop */
	default:
		log.Fatalf("error: too many matches to %q\n", flag.Arg(0))
	}

	command := matches[0]
	if err := cmd.ValidateCommand(command); err != nil {
		log.Fatalf("error validating command verb: %s\n", err)
	}

	switch command {
	case "init":
		//Initialise bioenv function...
		venv, err := vienv.Init(".", "myfirstbioenv")
		if err != nil {
			log.Fatal("Error initialising a new bio virtual environment")
		}
		fmt.Printf("New virtual environment: %v\n", venv)

	case "install":
		//install tool...
		toolName := "gatk"
		if err := venv.InstallTool(toolName); err != nil {
			log.Fatalf("Error installing tool: %s\n", toolName)
		}

	default:
		log.Fatalf("could not understand given command.\n")
	}

	// call tool... bioenv call gatk -- params params2 params3

}
