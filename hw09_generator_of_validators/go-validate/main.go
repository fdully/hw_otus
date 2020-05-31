package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

func main() {
	goFile := os.Getenv("GOFILE")
	if goFile == "" {
		log.Fatal("please set GOFILE environment variable or use go generate.")
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, goFile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create(strings.TrimSuffix(goFile, ".go") + "_validation_generated.go")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	validationData, err := GenerateValidationData(node)
	if err != nil {
		log.Fatal(err)
	}
	if validationData == nil {
		fmt.Println("nothing to generate.")
		os.Exit(0)
	}

	_, err = fmt.Fprint(outFile, string(validationData))
	if err != nil {
		log.Fatal(err)
	}
}
