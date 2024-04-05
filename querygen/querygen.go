package main

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/eebor/sweetquery/querygen/generator"
)

func main() {
	path := os.Getenv("GOFILE")
	if path == "" {
		log.Fatal("GOFILE must be set")
	}

	pkg, err := parseSource(path)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.Buffer{}

	err = generator.Generate(pkg, &buf)
	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create(strings.TrimSuffix(path, ".go") + "_builder.go")
	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()

	buf.WriteTo(outFile)
}
