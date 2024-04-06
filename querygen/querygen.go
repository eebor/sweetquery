package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eebor/sweetquery/querygen/generator"
	"github.com/eebor/sweetquery/querygen/model"
)

var (
	src string
	out string
)

func init() {
	flag.StringVar(&src, "src", "", "path to src file")
	flag.StringVar(&out, "out", "", "path to out file")
}

func main() {
	flag.Parse()

	if src == "" {
		src = os.Getenv("GOFILE")

		if src == "" {
			log.Fatal("src or GOFILE must be set")
		}
	}

	if out == "" {
		out = strings.TrimSuffix(src, ".go") + "_builder.go"
	}

	tasks, _, err := parseSource(src)
	if err != nil {
		log.Fatal(err)
	}

	if len(tasks) == 0 {
		log.Printf("no query found in file %s\n", src)
		return
	}

	imorts, err := getImports()
	if err != nil {
		log.Fatal(err)
	}

	pkg := model.Package{
		Name:    getOutPackageName(),
		Imports: imorts,
		Tasks:   tasks,
	}

	if !outIsSrc() {
		pkg.SrcPkgHandlePrefix = srcPkgName
	}

	buf := bytes.Buffer{}

	err = generator.Generate(&pkg, &buf)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(filepath.Dir(out), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()

	buf.WriteTo(outFile)
}

func outIsSrc() bool {
	return filepath.Dir(out) == filepath.Dir(src)
}

func getOutPackageName() string {
	return filepath.Base(filepath.Dir(out))
}
