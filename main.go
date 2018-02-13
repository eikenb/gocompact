package main

import (
	"bytes"
	"flag"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"

	"github.com/eikenb/gocompact/printer"
	"github.com/eikenb/tools/imports"
	// https://github.com/golang/go/issues/23782
	// "golang.org/x/tools/imports"
)

const printerMode = printer.UseSpaces | printer.TabIndent

var (
	cfg       = printer.Config{Mode: printerMode, Tabwidth: 4}
	writeFile = flag.Bool("w", false,
		"write back to source file instead of stdout")
)

func main() {
	flag.Parse()
	paths := flag.Args()

	if flag.NArg() == 0 {
		*writeFile = false
		src, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		err = processFile("<stdin>", src)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, path := range paths {
		fp, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		src, err := ioutil.ReadAll(fp)
		if err != nil {
			log.Fatal(err)
		}
		err = processFile(path, src)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func processFile(path string, src []byte) error {
	res, err := format(path, src)
	if err != nil {
		return err
	}

	if !bytes.Equal(res, src) {
		if *writeFile {
			err = ioutil.WriteFile(path, res, 0644)
			if err != nil {
				return err
			}
		} else {
			_, err := os.Stdout.Write(res)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func format(filename string, src []byte) ([]byte, error) {
	var fset *token.FileSet = token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	_, err = imports.FixImports(fset, file, filename)
	if err != nil {
		return nil, err
	}

	// cfg.Fprint(os.Stdout, fset, file)
	var buf bytes.Buffer
	cfg.Fprint(&buf, fset, file)
	res := buf.Bytes()
	return res, nil
}
