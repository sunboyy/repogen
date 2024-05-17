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
	pkgPtr := flag.String("pkg", ".", "package directory to scan for model struct and repository interface")
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

	request := GenerationRequest{
		Pkg:       *pkgPtr,
		ModelName: *modelPtr,
		RepoName:  *repoPtr,
		Dest:      *destPtr,
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
}

func printUsage() {
	fmt.Println(usageText)
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Println(version)
}

var (
	errNoPackageFound = errors.New("no package found")
)

func generateFromRequest(request GenerationRequest) (string, error) {
	cfg := packages.Config{
		Mode: packages.NeedName | packages.NeedTypes,
	}
	pkgs, err := packages.Load(&cfg, request.Pkg)
	if err != nil {
		return "", err
	}
	if len(pkgs) == 0 {
		return "", errNoPackageFound
	}

	pkg := pkgs[0]

	return generator.GenerateRepositoryImpl(pkg.Types, request.ModelName, request.RepoName)
}
