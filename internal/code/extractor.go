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
			switch spec := spec.(type) {
			case *ast.ImportSpec:
				var imp Import
				if spec.Name != nil {
					imp.Name = spec.Name.Name
				}
				importPath, err := strconv.Unquote(spec.Path.Value)
				if err != nil {
					fmt.Printf("cannot unquote import %s : %s \n", spec.Path.Value, err)
					continue
				}
				imp.Path = importPath

				file.Imports = append(file.Imports, imp)

			case *ast.TypeSpec:
				switch t := spec.Type.(type) {
				case *ast.StructType:
					file.Structs = append(file.Structs, extractStructType(spec.Name.Name, t))
				case *ast.InterfaceType:
					file.Interfaces = append(file.Interfaces, extractInterfaceType(spec.Name.Name, t))
				}
			}
		}
	}
	return file
}

func extractStructType(name string, structType *ast.StructType) Struct {
	str := Struct{
		Name: name,
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

	return str
}

func extractInterfaceType(name string, interfaceType *ast.InterfaceType) InterfaceType {
	intf := InterfaceType{
		Name: name,
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

		var comments []string
		if method.Doc != nil {
			for _, comment := range method.Doc.List {
				commentRunes := []rune(comment.Text)
				commentText := strings.TrimSpace(string(commentRunes[2:]))
				comments = append(comments, commentText)
			}
		}

		meth := extractFunction(name, comments, funcType)

		intf.Methods = append(intf.Methods, meth)
	}

	return intf
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

func extractFunction(name string, comments []string, funcType *ast.FuncType) Method {
	meth := Method{
		Name:     name,
		Comments: comments,
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
	switch expr := expr.(type) {
	case *ast.Ident:
		return SimpleType(expr.Name)

	case *ast.SelectorExpr:
		xExpr, ok := expr.X.(*ast.Ident)
		if !ok {
			return ExternalType{Name: expr.Sel.Name}
		}
		return ExternalType{PackageAlias: xExpr.Name, Name: expr.Sel.Name}

	case *ast.StarExpr:
		containedType := getType(expr.X)
		return PointerType{ContainedType: containedType}

	case *ast.ArrayType:
		containedType := getType(expr.Elt)
		return ArrayType{ContainedType: containedType}

	case *ast.MapType:
		keyType := getType(expr.Key)
		valueType := getType(expr.Value)
		return MapType{KeyType: keyType, ValueType: valueType}

	case *ast.InterfaceType:
		return extractInterfaceType("", expr)
	}

	return nil
}
