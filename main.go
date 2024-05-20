package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sunboyy/repogen/internal/generator"
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
	pkgPtr := flag.String(
		"pkg",
		".",
		"package directory to scan for model struct and repository interface",
	)
	destPtr := flag.String("dest", "", "destination file")
	modelPtr := flag.String("model", "", "model struct name")
	repoPtr := flag.String("repo", "", "repository interface name")
	modelPkgPtr := flag.String(
		"model-pkg",
		"",
		"package directory to scan for model struct. If not set, will fallback to -pkg.",
	)
	destPkgPtr := flag.String(
		"dest-pkg",
		"",
		"destination package path. If not set, will consider as in the same package as repository interface.",
	)
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

	request := GenerationRequest{
		Pkg:       *pkgPtr,
		ModelName: *modelPtr,
		RepoName:  *repoPtr,
		Dest:      *destPtr,
		ModelPkg:  *modelPkgPtr,
		DestPkg:   *destPkgPtr,
	}
	code, err := generateFromRequest(request)
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

type GenerationRequest struct {
	Pkg       string
	ModelName string
	RepoName  string
	Dest      string
	ModelPkg  string
	DestPkg   string
}

func printUsage() {
	fmt.Println(usageText)
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Println(version)
}

var (
	errNoPackageFound        = errors.New("no package found")
	errUnsupportMultiplePkgs = errors.New(
		`multiple packages are not supported, 
		please specify the package ID or directory path that only contains one package`,
	)
	errMissingPackageName = errors.New("missing package name")
)

func generateFromRequest(request GenerationRequest) (string, error) {
	cfg := packages.Config{
		Mode: packages.NeedName | packages.NeedTypes,
	}
	if request.ModelPkg == "" {
		request.ModelPkg = request.Pkg
	}
	if request.DestPkg == "" {
		request.DestPkg = request.Pkg
	}
	intfPkgID, err := getPkgID(request.Pkg)
	if err != nil {
		return "", err
	}
	modelPkgID, err := getPkgID(request.ModelPkg)
	if err != nil {
		return "", err
	}
	destPkgID, err := getPkgID(request.DestPkg)
	if err != nil {
		return "", err
	}
	pkgs, err := packages.Load(&cfg, intfPkgID, modelPkgID, destPkgID)
	if err != nil {
		return "", err
	}
	pkgM := packagesToMap(pkgs)
	return generator.GenerateRepositoryImpl(
		pkgM[modelPkgID].Types,
		pkgM[intfPkgID].Types,
		pkgM[destPkgID].Types,
		request.ModelName,
		request.RepoName,
	)
}

func getPkgID(pattern string) (string, error) {
	pkgs, err := packages.Load(nil, pattern)
	if err != nil {
		return "", err
	}
	if len(pkgs) < 1 {
		return "", errNoPackageFound
	}
	if len(pkgs) > 1 {
		return "", errUnsupportMultiplePkgs
	}
	// when no go file in the package, the package name will be empty
	// this prevent the missing field upfront.
	if pkgs[0].Name == "" {
		return "", fmt.Errorf("%w on %s", errMissingPackageName, pattern)
	}
	return pkgs[0].ID, nil
}

func packagesToMap(pkgs []*packages.Package) map[string]*packages.Package {
	m := make(map[string]*packages.Package)
	for _, pkg := range pkgs {
		m[pkg.ID] = pkg
	}
	return m
}
