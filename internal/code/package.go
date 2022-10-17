package code

import (
	"go/ast"
	"strings"
)

// ParsePackage extracts package name, struct and interface implementations from
// map[string]*ast.Package. Test files will be ignored.
func ParsePackage(pkgs map[string]*ast.Package) (Package, error) {
	pkg := NewPackage()
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
// from ParsePackage
type Package struct {
	Name       string
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
