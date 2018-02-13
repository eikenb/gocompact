package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"

	"gocompact/printer"

	// "golang.org/x/tools/imports"
	"github.com/eikenb/tools/imports"
)

const printerMode = printer.UseSpaces | printer.TabIndent

var cfg = printer.Config{Mode: printerMode, Tabwidth: 4}

func main() {
	flag.Parse()
	paths := flag.Args()
	for _, path := range paths {
		err := processFile(path)
		if err != nil {
			panic(err)
		}
	}
}

func processFile(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	src, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}
	return format(path, src)
}

func format(filename string, src []byte) error {
	var fset *token.FileSet = token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}
	a, err := imports.FixImports(fset, file, filename)
	if err != nil {
		return err
	}
	fmt.Println(a, err)

	cfg.Fprint(os.Stdout, fset, file)
	return nil
}
