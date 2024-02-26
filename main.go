package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/sunboyy/repogen/internal/code"
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
	pkgDirPtr := flag.String("pkg", ".", "package directory to scan for repository interface")
	destPtr := flag.String("dest", "", "destination file")
	modelDirPtr := flag.String("model-dir", ".", "package directory to scan for model struct")
	modelPtr := flag.String("model", "", "model struct name")
	repoPtr := flag.String("repo", "", "repository interface name")
	destPkgName := flag.String("dest-pkg", "", "destination package name")

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

	code, err := generateFromRequest(*pkgDirPtr, *modelDirPtr, *modelPtr, *repoPtr, *destPkgName)
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

var (
	errStructNotFound    = errors.New("struct not found")
	errInterfaceNotFound = errors.New("interface not found")
)

func generateFromRequest(pkgDir, modelDirPtr, structModelName,
	repositoryInterfaceName, destPkgName string) (string, error) {
	pkg, err := parsePkg(pkgDir)
	if err != nil {
		return "", err
	}
	if destPkgName == "" {
		destPkgName = pkg.Name
	}
	if pkgDir == modelDirPtr {
		structModel, ok := pkg.Structs[structModelName]
		if !ok {
			return "", errStructNotFound
		}
		return generateRepository(pkg, structModel, repositoryInterfaceName, destPkgName)
	} else {
		modelPkg, err := parsePkg(modelDirPtr)
		if err != nil {
			return "", err
		}
		structModel, ok := modelPkg.Structs[structModelName]
		if !ok {
			return "", errStructNotFound
		}
		structModel.PackageAlias = modelPkg.Name
		return generateRepository(pkg, structModel, repositoryInterfaceName, destPkgName)
	}
}

func parsePkg(pkgDir string) (code.Package, error) {
	dirParser := func(dir string) (pkgs map[string]*ast.Package, err error) {
		return parser.ParseDir(token.NewFileSet(), dir, nil, parser.ParseComments)
	}

	pkgPserser := code.NewPackageParser(dirParser, parsePackageID)
	return pkgPserser.ParsePackage(pkgDir)
}

var (
	errNoPackagesFound = errors.New("no packages found")
)

func parsePackageID(dir string) (string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes,
		Dir:  dir,
	}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		return "", err
	}
	if len(pkgs) > 0 {
		return pkgs[0].ID, nil
	}
	return "", errNoPackagesFound
}

func generateRepository(pkg code.Package, structModel code.Struct, repositoryInterfaceName,
	destPkgName string) (string, error) {
	intf, ok := pkg.Interfaces[repositoryInterfaceName]
	if !ok {
		return "", errInterfaceNotFound
	}

	var methodSpecs []spec.MethodSpec
	for _, method := range intf.Methods {
		methodSpec, err := spec.ParseInterfaceMethod(pkg.Structs, structModel, method)
		if err != nil {
			return "", err
		}
		methodSpecs = append(methodSpecs, methodSpec)
	}

	return generator.GenerateRepository(destPkgName, structModel, intf.Name, methodSpecs)
}
