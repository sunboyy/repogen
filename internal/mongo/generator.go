package mongo

import (
	"bytes"
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
	for _, param := range methodSpec.Params {
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
	case spec.InsertOperation:
		return g.generateInsertImplementation(operation)
	case spec.FindOperation:
		return g.generateFindImplementation(operation)
	case spec.UpdateOperation:
		return g.generateUpdateImplementation(operation)
	case spec.DeleteOperation:
		return g.generateDeleteImplementation(operation)
	case spec.CountOperation:
		return g.generateCountImplementation(operation)
	}

	return "", OperationNotSupportedError
}

func (g RepositoryGenerator) generateInsertImplementation(operation spec.InsertOperation) (string, error) {
	if operation.Mode == spec.QueryModeOne {
		return insertOneTemplate, nil
	}
	return insertManyTemplate, nil
}

func (g RepositoryGenerator) generateFindImplementation(operation spec.FindOperation) (string, error) {
	querySpec, err := g.mongoQuerySpec(operation.Query)
	if err != nil {
		return "", err
	}

	tmplData := mongoFindTemplateData{
		EntityType: g.StructModel.Name,
		QuerySpec:  querySpec,
	}

	if operation.Mode == spec.QueryModeOne {
		return generateFromTemplate("mongo_repository_findone", findOneTemplate, tmplData)
	}
	return generateFromTemplate("mongo_repository_findmany", findManyTemplate, tmplData)
}

func (g RepositoryGenerator) generateUpdateImplementation(operation spec.UpdateOperation) (string, error) {
	var fields []updateField
	for _, field := range operation.Fields {
		bsonTag, err := g.bsonTagFromFieldName(field.Name)
		if err != nil {
			return "", err
		}
		fields = append(fields, updateField{BsonTag: bsonTag, ParamIndex: field.ParamIndex})
	}

	querySpec, err := g.mongoQuerySpec(operation.Query)
	if err != nil {
		return "", err
	}

	tmplData := mongoUpdateTemplateData{
		UpdateFields: fields,
		QuerySpec:    querySpec,
	}

	if operation.Mode == spec.QueryModeOne {
		return generateFromTemplate("mongo_repository_updateone", updateOneTemplate, tmplData)
	}
	return generateFromTemplate("mongo_repository_updatemany", updateManyTemplate, tmplData)
}

func (g RepositoryGenerator) generateDeleteImplementation(operation spec.DeleteOperation) (string, error) {
	querySpec, err := g.mongoQuerySpec(operation.Query)
	if err != nil {
		return "", err
	}

	tmplData := mongoDeleteTemplateData{
		QuerySpec: querySpec,
	}

	if operation.Mode == spec.QueryModeOne {
		return generateFromTemplate("mongo_repository_deleteone", deleteOneTemplate, tmplData)
	}
	return generateFromTemplate("mongo_repository_deletemany", deleteManyTemplate, tmplData)
}

func (g RepositoryGenerator) generateCountImplementation(operation spec.CountOperation) (string, error) {
	querySpec, err := g.mongoQuerySpec(operation.Query)
	if err != nil {
		return "", err
	}

	tmplData := mongoCountTemplateData{
		QuerySpec: querySpec,
	}

	return generateFromTemplate("mongo_repository_count", countTemplate, tmplData)
}

func (g RepositoryGenerator) mongoQuerySpec(query spec.QuerySpec) (querySpec, error) {
	var predicates []predicate

	for _, predicateSpec := range query.Predicates {
		bsonTag, err := g.bsonTagFromFieldName(predicateSpec.Field)
		if err != nil {
			return querySpec{}, err
		}

		predicates = append(predicates, predicate{
			Field:      bsonTag,
			Comparator: predicateSpec.Comparator,
			ParamIndex: predicateSpec.ParamIndex,
		})
	}

	return querySpec{
		Operator:   query.Operator,
		Predicates: predicates,
	}, nil
}

func (g RepositoryGenerator) bsonTagFromFieldName(fieldName string) (string, error) {
	structField, ok := g.StructModel.Fields.ByName(fieldName)
	if !ok {
		return "", fmt.Errorf("struct field %s not found", fieldName)
	}

	bsonTag, ok := structField.Tags["bson"]
	if !ok {
		return "", BsonTagNotFoundError
	}

	return bsonTag[0], nil
}

func (g RepositoryGenerator) structName() string {
	return g.InterfaceName + "Mongo"
}

func generateFromTemplate(name string, templateString string, tmplData interface{}) (string, error) {
	tmpl, err := template.New(name).Parse(templateString)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	if err := tmpl.Execute(buffer, tmplData); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
