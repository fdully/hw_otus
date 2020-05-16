package main

import (
	"fmt"
	"log"
	"os"
)

const minimalArgs = 3

func main() {
	// Place your code here
	if len(os.Args) < minimalArgs {
		fmt.Println("Usage: go-envdir /path/to/env/dir shell-command arg1 arg2")
	}

	environmentFromFiles, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	commandEnvrionment := MakeCommandEnv(os.Environ(), environmentFromFiles)

	os.Exit(RunCmd(os.Args[2:], commandEnvrionment))
}
