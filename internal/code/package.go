package code

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	errNoPackagesFound = errors.New("no packages found")
)

type dirparseFn func(dir string) (pkgs map[string]*ast.Package, err error)

type pkgPathParseFn func(dir string) (path string, err error)

func dirParser(dir string) (pkgs map[string]*ast.Package, err error) {
	return parser.ParseDir(token.NewFileSet(), dir, nil, parser.ParseComments)
}
func parsePackagePath(dir string) (string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName,
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

type PackageParser struct {
	DirParser     dirparseFn
	PkgPathParser pkgPathParseFn
}

func NewPackageParser() *PackageParser {
	return &PackageParser{
		DirParser:     dirParser,
		PkgPathParser: parsePackagePath,
	}
}

// ParsePackage extracts package name, struct and interface implementations from pkgDir.
// Test files will be ignored.
func (p *PackageParser) ParsePackage(pkgDir string) (Package, error) {
	pkg := NewPackage()
	var err error
	pkg.Path, err = p.PkgPathParser(pkgDir)
	if err != nil {
		return Package{}, err
	}

	pkgs, err := p.DirParser(pkgDir)
	if err != nil {
		return Package{}, err
	}

	for _, astPkg := range pkgs {
		for fileName, file := range astPkg.Files {
			if strings.HasSuffix(fileName, "_test.go") {
				continue
			}

			if err := pkg.addFile(ExtractComponents(file)); err != nil {
				return Package{}, err
			}
		}
	}
	return pkg, nil
}

// Package stores package name, struct and interface implementations as a result
// from ParsePackage.
type Package struct {
	Name       string
	Path       string // the path of package itself
	Structs    map[string]Struct
	Interfaces map[string]InterfaceType
}

// NewPackage is a constructor function for Package.
func NewPackage() Package {
	return Package{
		Structs:    map[string]Struct{},
		Interfaces: map[string]InterfaceType{},
	}
}

// addFile alters the Package by adding struct and interface implementations in
// the extracted file. If the package name conflicts, it will return error.
func (pkg *Package) addFile(file File) error {
	if pkg.Name == "" {
		pkg.Name = file.PackageName
	} else if pkg.Name != file.PackageName {
		return ErrAmbiguousPackageName
	}

	for _, structImpl := range file.Structs {
		if _, ok := pkg.Structs[structImpl.Name]; ok {
			return DuplicateStructError(structImpl.Name)
		}
		pkg.Structs[structImpl.Name] = structImpl
	}

	for _, interfaceImpl := range file.Interfaces {
		if _, ok := pkg.Interfaces[interfaceImpl.Name]; ok {
			return DuplicateInterfaceError(interfaceImpl.Name)
		}
		pkg.Interfaces[interfaceImpl.Name] = interfaceImpl
	}

	return nil
}
