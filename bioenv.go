package main

import (
	"bioenv/vienv"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Welcome to BioEnv")

	//Initialise bioenv function...
	v, err := vienv.Init(".", "myfirstbioenv")
	if err != nil {
		log.Fatal("Error initialising a new bio virtual environment")
	}

	fmt.Printf("New virtual environment: %v", v)

	//install tool...
	toolName := "gatk"
	if err := v.InstallTool(toolName); err != nil {
		log.Fatalf("Error installing tool: %v", toolName)
	}

	// call tool... bioenv call gatk -- params params2 params3

}
