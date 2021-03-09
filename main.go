package main

import (
	"errors"
	"flag"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/generator"
	"github.com/sunboyy/repogen/internal/spec"
)

func main() {
	sourcePtr := flag.String("src", "", "source file")
	destPtr := flag.String("dest", "", "destination file")
	modelPtr := flag.String("model", "", "model struct name")
	repoPtr := flag.String("repo", "", "repository interface name")

	flag.Parse()

	if *sourcePtr == "" {
		panic("-source flag required")
	}
	if *modelPtr == "" {
		panic("-model flag required")
	}
	if *repoPtr == "" {
		panic("-repo flag required")
	}

	code, err := generateFromRequest(*sourcePtr, *modelPtr, *repoPtr)
	if err != nil {
		panic(err)
	}

	dest := os.Stdout
	if *destPtr != "" {
		if err := os.MkdirAll(filepath.Dir(*destPtr), os.ModePerm); err != nil {
			panic(err)
		}
		file, err := os.Create(*destPtr)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		dest = file
	}

	if _, err := dest.WriteString(code); err != nil {
		panic(err)
	}
}

func generateFromRequest(fileName, structModelName, repositoryInterfaceName string) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	file := code.ExtractComponents(f)

	structModel, ok := file.Structs.ByName(structModelName)
	if !ok {
		return "", errors.New("struct model not found")
	}

	intf, ok := file.Interfaces.ByName(repositoryInterfaceName)
	if !ok {
		return "", errors.New("interface model not found")
	}

	var methodSpecs []spec.MethodSpec
	for _, method := range intf.Methods {
		methodSpec, err := spec.ParseInterfaceMethod(file.Structs, structModel, method)
		if err != nil {
			return "", err
		}
		methodSpecs = append(methodSpecs, methodSpec)
	}

	return generator.GenerateRepository(file.PackageName, structModel, intf.Name, methodSpecs)
}
