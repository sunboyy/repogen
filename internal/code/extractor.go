package code

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"
)

// ExtractComponents converts ast file into code components model
func ExtractComponents(f *ast.File) File {
	var file File
	file.PackageName = f.Name.Name

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if ok {
				var imp Import
				if importSpec.Name != nil {
					imp.Name = importSpec.Name.Name
				}
				importPath, err := strconv.Unquote(importSpec.Path.Value)
				if err != nil {
					fmt.Printf("cannot unquote import %s : %s \n", importSpec.Path.Value, err)
					continue
				}
				imp.Path = importPath

				file.Imports = append(file.Imports, imp)
			}

			typeSpec, ok := spec.(*ast.TypeSpec)
			if ok {
				structType, ok := typeSpec.Type.(*ast.StructType)
				if ok {
					str := Struct{
						Name: typeSpec.Name.Name,
					}

					for _, field := range structType.Fields.List {
						var strField StructField
						for _, name := range field.Names {
							strField.Name = name.Name
							break
						}
						strField.Type = getType(field.Type)
						if field.Tag != nil {
							strField.Tags = extractStructTag(field.Tag.Value)
						}

						str.Fields = append(str.Fields, strField)
					}

					file.Structs = append(file.Structs, str)
				}

				interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
				if ok {
					intf := InterfaceType{
						Name: typeSpec.Name.Name,
					}

					for _, method := range interfaceType.Methods.List {
						funcType, ok := method.Type.(*ast.FuncType)
						if !ok {
							continue
						}

						var name string
						for _, n := range method.Names {
							name = n.Name
							break
						}

						meth := extractFunction(name, funcType)

						intf.Methods = append(intf.Methods, meth)
					}

					file.Interfaces = append(file.Interfaces, intf)
				}
			}
		}
	}
	return file
}

func extractStructTag(tagValue string) map[string][]string {
	tagTokens := strings.Fields(tagValue[1 : len(tagValue)-1])

	tags := make(map[string][]string)
	for _, tagToken := range tagTokens {
		colonIndex := strings.Index(tagToken, ":")
		if colonIndex == -1 {
			continue
		}
		tagKey := tagToken[:colonIndex]
		tagValue, err := strconv.Unquote(tagToken[colonIndex+1:])
		if err != nil {
			fmt.Printf("cannot unquote struct tag %s : %s\n", tagToken[colonIndex+1:], err)
			continue
		}
		tagValues := strings.Split(tagValue, ",")
		tags[tagKey] = tagValues
	}

	return tags
}

func extractFunction(name string, funcType *ast.FuncType) Method {
	meth := Method{
		Name: name,
	}
	for _, param := range funcType.Params.List {
		paramType := getType(param.Type)

		if len(param.Names) == 0 {
			meth.Params = append(meth.Params, Param{Type: paramType})
			continue
		}

		for _, name := range param.Names {
			meth.Params = append(meth.Params, Param{
				Name: name.Name,
				Type: paramType,
			})
		}
	}

	if funcType.Results != nil {
		for _, result := range funcType.Results.List {
			meth.Returns = append(meth.Returns, getType(result.Type))
		}
	}

	return meth
}

func getType(expr ast.Expr) Type {
	identExpr, ok := expr.(*ast.Ident)
	if ok {
		return SimpleType(identExpr.Name)
	}

	selectorExpr, ok := expr.(*ast.SelectorExpr)
	if ok {
		xExpr, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			return ExternalType{Name: selectorExpr.Sel.Name}
		}
		return ExternalType{PackageAlias: xExpr.Name, Name: selectorExpr.Sel.Name}
	}

	starExpr, ok := expr.(*ast.StarExpr)
	if ok {
		containedType := getType(starExpr.X)
		return PointerType{ContainedType: containedType}
	}

	arrayType, ok := expr.(*ast.ArrayType)
	if ok {
		containedType := getType(arrayType.Elt)
		return ArrayType{containedType}
	}

	intfType, ok := expr.(*ast.InterfaceType)
	if ok {
		var methods []Method
		for _, method := range intfType.Methods.List {
			funcType, ok := method.Type.(*ast.FuncType)
			if !ok {
				continue
			}

			var name string
			for _, n := range method.Names {
				name = n.Name
				break
			}

			methods = append(methods, extractFunction(name, funcType))
		}

		return InterfaceType{
			Methods: methods,
		}
	}

	return nil
}
