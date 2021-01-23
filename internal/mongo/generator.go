package mongo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"text/template"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/spec"
)

// NewGenerator creates a new instance of MongoDB repository generator
func NewGenerator(structModel code.Struct, interfaceName string) RepositoryGenerator {
	return RepositoryGenerator{
		StructModel:   structModel,
		InterfaceName: interfaceName,
	}
}

// RepositoryGenerator provides repository constructor and method generation from provided specification
type RepositoryGenerator struct {
	StructModel   code.Struct
	InterfaceName string
}

// GenerateConstructor generates mongo repository struct implementation and constructor for the struct
func (g RepositoryGenerator) GenerateConstructor(buffer io.Writer) error {
	tmpl, err := template.New("mongo_repository_base").Parse(constructorTemplate)
	if err != nil {
		return err
	}

	tmplData := mongoConstructorTemplateData{
		InterfaceName:  g.InterfaceName,
		ImplStructName: g.structName(),
	}

	if err := tmpl.Execute(buffer, tmplData); err != nil {
		return err
	}

	return nil
}

// GenerateMethod generates implementation of from provided method specification
func (g RepositoryGenerator) GenerateMethod(methodSpec spec.MethodSpec, buffer io.Writer) error {
	tmpl, err := template.New("mongo_repository_method").Parse(methodTemplate)
	if err != nil {
		return err
	}

	implementation, err := g.generateMethodImplementation(methodSpec)
	if err != nil {
		return err
	}

	var paramTypes []code.Type
	for _, param := range methodSpec.Params[1:] {
		paramTypes = append(paramTypes, param.Type)
	}

	tmplData := mongoMethodTemplateData{
		StructName:     g.structName(),
		MethodName:     methodSpec.Name,
		ParamTypes:     paramTypes,
		ReturnTypes:    methodSpec.Returns,
		Implementation: implementation,
	}

	if err := tmpl.Execute(buffer, tmplData); err != nil {
		return err
	}

	return nil
}

func (g RepositoryGenerator) generateMethodImplementation(methodSpec spec.MethodSpec) (string, error) {
	switch operation := methodSpec.Operation.(type) {
	case spec.FindOperation:
		return g.generateFindImplementation(operation)
	}

	return "", errors.New("method spec not supported")
}

func (g RepositoryGenerator) generateFindImplementation(operation spec.FindOperation) (string, error) {
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

func (g RepositoryGenerator) structName() string {
	return g.InterfaceName + "Mongo"
}
