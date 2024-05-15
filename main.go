package main

import (
	"errors"
	"flag"
	"fmt"
	"go/types"
	"log"
	"os"
	"path/filepath"

	"github.com/sunboyy/repogen/internal/generator"
	"github.com/sunboyy/repogen/internal/spec"
	"golang.org/x/tools/go/packages"
)

const usageText = `repogen generates MongoDB repository implementation from repository interface

  Find more information at: https://github.com/sunboyy/repogen

Supported options:`

// version indicates the version of repogen.
const version = "v0.4-next"

func main() {
	flag.Usage = printUsage

	versionPtr := flag.Bool("version", false, "print version of repogen")
	pkgDirPtr := flag.String("pkg", ".", "package directory to scan for model struct and repository interface")
	destPtr := flag.String("dest", "", "destination file")
	modelPtr := flag.String("model", "", "model struct name")
	repoPtr := flag.String("repo", "", "repository interface name")

	flag.Parse()

	if *versionPtr {
		printVersion()
		return
	}

	if *modelPtr == "" {
		printUsage()
		log.Fatal("-model flag required")
	}
	if *repoPtr == "" {
		printUsage()
		log.Fatal("-repo flag required")
	}

	code, err := generateFromRequest(*pkgDirPtr, *modelPtr, *repoPtr)
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

func printUsage() {
	fmt.Println(usageText)
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Println(version)
}

func generateFromRequest(pkgDir, structModelName, repositoryInterfaceName string) (string, error) {
	cfg := packages.Config{
		Mode: packages.NeedName | packages.NeedTypes,
	}
	pkgs, err := packages.Load(&cfg, pkgDir)
	if err != nil {
		return "", err
	}
	if len(pkgs) == 0 {
		return "", errNoPackageFound
	}

	pkg := pkgs[0]

	return generateRepository(pkg.Types, structModelName, repositoryInterfaceName)
}

var (
	errNoPackageFound    = errors.New("no package found")
	errStructNotFound    = errors.New("struct not found")
	errNotNamedStruct    = errors.New("not a named struct")
	errInterfaceNotFound = errors.New("interface not found")
	errNotInterface      = errors.New("not an interface")
)

func generateRepository(pkg *types.Package, structModelName, repositoryInterfaceName string) (string, error) {
	structModelObj := pkg.Scope().Lookup(structModelName)
	if structModelObj == nil {
		return "", errStructNotFound
	}
	namedStruct, ok := structModelObj.Type().(*types.Named)
	if !ok {
		return "", errNotNamedStruct
	}

	intfObj := pkg.Scope().Lookup(repositoryInterfaceName)
	if intfObj == nil {
		return "", errInterfaceNotFound
	}
	intf, ok := intfObj.Type().Underlying().(*types.Interface)
	if !ok {
		return "", errNotInterface
	}

	var methodSpecs []spec.MethodSpec
	for i := 0; i < intf.NumMethods(); i++ {
		method := intf.Method(i)
		log.Println("Generating method:", method.Name())

		methodSpec, err := spec.ParseInterfaceMethod(pkg, namedStruct, method)
		if err != nil {
			return "", err
		}
		methodSpecs = append(methodSpecs, methodSpec)
	}

	return generator.GenerateRepository(pkg, namedStruct, repositoryInterfaceName, methodSpecs)
}
