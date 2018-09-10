package main

import (
	"os"
	"fmt"
	"go/parser"
	"go/token"
	"go/ast"
	"go/printer"
	)

func printUsageAndExit() {
	fmt.Println("Usage: go-ast-trace <trace-type> <input-files>")
	fmt.Println("Trace Types:")
	fmt.Println(" locks		Prints trace before and after locks of mutexes and channels")
	os.Exit(0)
}

func printErrorAndExit(err error) {
	fmt.Println("ERROR: ", err.Error())
	os.Exit(1)
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		printUsageAndExit()
	}
	traceType := args[0]
	inputFilenames := args[1:]

	for _, inputFilename := range inputFilenames {
		fset, node := parseFile(inputFilename)

		switch traceType {
		case "locks":
			traceLocks(node)
		default:
			printUsageAndExit()
		}

		writeFile(inputFilename, fset, node)
	}
}

func parseFile(filename string) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		printErrorAndExit(err)
	}
	ast.Print(fset, node)
	return fset, node
}

func writeFile(filename string, fset *token.FileSet, node *ast.File) {
	f, err := os.Create(filename)
	if err != nil {
		printErrorAndExit(err)
	}
	defer f.Close()
	if err := printer.Fprint(f, fset, node); err != nil {
		printErrorAndExit(err)
	}
}