package mongo

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/spec"
	"golang.org/x/tools/imports"
)

// GenerateMongoRepository generates mongodb repository
func GenerateMongoRepository(packageName string, structModel code.Struct, intf code.Interface) (string, error) {
	repositorySpec, err := spec.ParseRepositoryInterface(structModel, intf)
	if err != nil {
		return "", err
	}

	generator := mongoRepositoryGenerator{
		PackageName:    packageName,
		StructModel:    structModel,
		RepositorySpec: repositorySpec,
	}

	output, err := generator.Generate()
	if err != nil {
		return "", err
	}

	return output, nil
}

type mongoRepositoryGenerator struct {
	PackageName    string
	StructModel    code.Struct
	RepositorySpec spec.RepositorySpec
}

func (g mongoRepositoryGenerator) Generate() (string, error) {
	buffer := new(bytes.Buffer)
	if err := g.generateBaseContent(buffer); err != nil {
		return "", err
	}

	for _, method := range g.RepositorySpec.Methods {
		if err := g.generateMethod(buffer, method); err != nil {
			return "", err
		}
	}

	newOutput, err := imports.Process("", buffer.Bytes(), nil)
	if err != nil {
		return "", err
	}

	return string(newOutput), nil
}

func (g mongoRepositoryGenerator) generateBaseContent(buffer *bytes.Buffer) error {
	tmpl, err := template.New("mongo_repository_base").Parse(baseTemplate)
	if err != nil {
		return err
	}

	tmplData := mongoBaseTemplateData{
		PackageName:   g.PackageName,
		InterfaceName: g.RepositorySpec.InterfaceName,
		StructName:    g.structName(),
	}

	if err := tmpl.Execute(buffer, tmplData); err != nil {
		return err
	}

	return nil
}

func (g mongoRepositoryGenerator) generateMethod(buffer *bytes.Buffer, method spec.MethodSpec) error {
	tmpl, err := template.New("mongo_repository_method").Parse(methodTemplate)
	if err != nil {
		return err
	}

	implementation, err := g.generateMethodImplementation(method)
	if err != nil {
		return err
	}

	var paramTypes []code.Type
	for _, param := range method.Params[1:] {
		paramTypes = append(paramTypes, param.Type)
	}

	tmplData := mongoMethodTemplateData{
		StructName:     g.structName(),
		MethodName:     method.Name,
		ParamTypes:     paramTypes,
		ReturnTypes:    method.Returns,
		Implementation: implementation,
	}

	if err := tmpl.Execute(buffer, tmplData); err != nil {
		return err
	}

	return nil
}

func (g mongoRepositoryGenerator) generateMethodImplementation(methodSpec spec.MethodSpec) (string, error) {
	switch operation := methodSpec.Operation.(type) {
	case spec.FindOperation:
		return g.generateFindImplementation(operation)
	}

	return "", errors.New("method spec not supported")
}

func (g mongoRepositoryGenerator) generateFindImplementation(operation spec.FindOperation) (string, error) {
	buffer := new(bytes.Buffer)

	var predicates []predicate
	for _, predicateSpec := range operation.Query.Predicates {
		structField, ok := g.StructModel.Fields.ByName(predicateSpec.Field)
		if !ok {
			return "", fmt.Errorf("struct field %s not found", predicateSpec.Field)
		}

		bsonTag, ok := structField.Tags["bson"]
		if !ok {
			return "", fmt.Errorf("struct field %s does not have bson tag", predicateSpec.Field)
		}

		predicates = append(predicates, predicate{Field: bsonTag[0], Comparator: predicateSpec.Comparator})
	}

	tmplData := mongoFindTemplateData{
		EntityType: g.StructModel.Name,
		QuerySpec: querySpec{
			Operator:   operation.Query.Operator,
			Predicates: predicates,
		},
	}

	if operation.Mode == spec.QueryModeOne {
		tmpl, err := template.New("mongo_repository_findone").Parse(findOneTemplate)
		if err != nil {
			return "", err
		}

		if err := tmpl.Execute(buffer, tmplData); err != nil {
			return "", err
		}
	} else {
		tmpl, err := template.New("mongo_repository_findmany").Parse(findManyTemplate)
		if err != nil {
			return "", err
		}

		if err := tmpl.Execute(buffer, tmplData); err != nil {
			return "", err
		}
	}

	return buffer.String(), nil
}

func (g mongoRepositoryGenerator) structName() string {
	return g.RepositorySpec.InterfaceName + "Mongo"
}
