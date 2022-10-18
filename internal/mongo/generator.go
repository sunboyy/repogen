package mongo

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
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

// Imports returns necessary imports for the mongo repository implementation.
func (g RepositoryGenerator) Imports() [][]code.Import {
	return [][]code.Import{
		{
			{Path: "context"},
		},
		{
			{Path: "go.mongodb.org/mongo-driver/bson"},
			{Path: "go.mongodb.org/mongo-driver/bson/primitive"},
			{Path: "go.mongodb.org/mongo-driver/mongo"},
			{Path: "go.mongodb.org/mongo-driver/mongo/options"},
		},
	}
}

// GenerateStruct creates codegen.StructBuilder of mongo repository
// implementation struct.
func (g RepositoryGenerator) GenerateStruct() codegen.StructBuilder {
	return codegen.StructBuilder{
		Name: g.repoImplStructName(),
		Fields: code.StructFields{
			{
				Name: "collection",
				Type: code.PointerType{
					ContainedType: code.ExternalType{
						PackageAlias: "mongo",
						Name:         "Collection",
					},
				},
			},
		},
	}
}

// GenerateConstructor creates codegen.FunctionBuilder of a constructor for
// mongo repository implementation struct.
func (g RepositoryGenerator) GenerateConstructor() (codegen.FunctionBuilder, error) {
	tmpl, err := template.New("mongo_constructor_body").Parse(constructorBody)
	if err != nil {
		return codegen.FunctionBuilder{}, err
	}

	tmplData := constructorBodyData{
		ImplStructName: g.repoImplStructName(),
	}

	buffer := new(bytes.Buffer)
	if err := tmpl.Execute(buffer, tmplData); err != nil {
		return codegen.FunctionBuilder{}, err
	}

	return codegen.FunctionBuilder{
		Name: "New" + g.InterfaceName,
		Params: []code.Param{
			{
				Name: "collection",
				Type: code.PointerType{
					ContainedType: code.ExternalType{
						PackageAlias: "mongo",
						Name:         "Collection",
					},
				},
			},
		},
		Returns: []code.Type{
			code.SimpleType(g.InterfaceName),
		},
		Body: buffer.String(),
	}, nil
}

// GenerateMethod creates codegen.MethodBuilder of repository method from the
// provided method specification.
func (g RepositoryGenerator) GenerateMethod(methodSpec spec.MethodSpec) (codegen.MethodBuilder, error) {
	var params []code.Param
	for i, param := range methodSpec.Params {
		params = append(params, code.Param{
			Name: fmt.Sprintf("arg%d", i),
			Type: param.Type,
		})
	}

	implementation, err := g.generateMethodImplementation(methodSpec)
	if err != nil {
		return codegen.MethodBuilder{}, err
	}

	return codegen.MethodBuilder{
		Receiver: codegen.MethodReceiver{
			Name:    "r",
			Type:    code.SimpleType(g.repoImplStructName()),
			Pointer: true,
		},
		Name:    methodSpec.Name,
		Params:  params,
		Returns: methodSpec.Returns,
		Body:    implementation,
	}, nil
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
	default:
		return "", NewOperationNotSupportedError(operation.Name())
	}
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

	sorts, err := g.mongoSorts(operation.Sorts)
	if err != nil {
		return "", err
	}

	tmplData := mongoFindTemplateData{
		EntityType: g.StructModel.Name,
		QuerySpec:  querySpec,
		Sorts:      sorts,
	}

	if operation.Mode == spec.QueryModeOne {
		return generateFromTemplate("mongo_repository_findone", findOneTemplate, tmplData)
	}
	return generateFromTemplate("mongo_repository_findmany", findManyTemplate, tmplData)
}

func (g RepositoryGenerator) mongoSorts(sortSpec []spec.Sort) ([]findSort, error) {
	var sorts []findSort

	for _, s := range sortSpec {
		bsonFieldReference, err := g.bsonFieldReference(s.FieldReference)
		if err != nil {
			return nil, err
		}

		sorts = append(sorts, findSort{
			BsonTag:  bsonFieldReference,
			Ordering: s.Ordering,
		})
	}

	return sorts, nil
}

func (g RepositoryGenerator) generateUpdateImplementation(operation spec.UpdateOperation) (string, error) {
	update, err := g.getMongoUpdate(operation.Update)
	if err != nil {
		return "", err
	}

	querySpec, err := g.mongoQuerySpec(operation.Query)
	if err != nil {
		return "", err
	}

	tmplData := mongoUpdateTemplateData{
		Update:    update,
		QuerySpec: querySpec,
	}

	if operation.Mode == spec.QueryModeOne {
		return generateFromTemplate("mongo_repository_updateone", updateOneTemplate, tmplData)
	}
	return generateFromTemplate("mongo_repository_updatemany", updateManyTemplate, tmplData)
}

func (g RepositoryGenerator) getMongoUpdate(updateSpec spec.Update) (update, error) {
	switch updateSpec := updateSpec.(type) {
	case spec.UpdateModel:
		return updateModel{}, nil
	case spec.UpdateFields:
		update := make(updateFields)
		for _, field := range updateSpec {
			bsonFieldReference, err := g.bsonFieldReference(field.FieldReference)
			if err != nil {
				return querySpec{}, err
			}

			updateKey := getUpdateOperatorKey(field.Operator)
			if updateKey == "" {
				return querySpec{}, NewUpdateOperatorNotSupportedError(field.Operator)
			}
			updateField := updateField{
				BsonTag:    bsonFieldReference,
				ParamIndex: field.ParamIndex,
			}
			update[updateKey] = append(update[updateKey], updateField)
		}
		return update, nil
	default:
		return nil, NewUpdateTypeNotSupportedError(updateSpec)
	}
}

func getUpdateOperatorKey(operator spec.UpdateOperator) string {
	switch operator {
	case spec.UpdateOperatorSet:
		return "$set"
	case spec.UpdateOperatorPush:
		return "$push"
	case spec.UpdateOperatorInc:
		return "$inc"
	default:
		return ""
	}
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
		bsonFieldReference, err := g.bsonFieldReference(predicateSpec.FieldReference)
		if err != nil {
			return querySpec{}, err
		}

		predicates = append(predicates, predicate{
			Field:      bsonFieldReference,
			Comparator: predicateSpec.Comparator,
			ParamIndex: predicateSpec.ParamIndex,
		})
	}

	return querySpec{
		Operator:   query.Operator,
		Predicates: predicates,
	}, nil
}

func (g RepositoryGenerator) bsonFieldReference(fieldReference spec.FieldReference) (string, error) {
	var bsonTags []string
	for _, field := range fieldReference {
		tag, err := g.bsonTagFromField(field)
		if err != nil {
			return "", err
		}
		bsonTags = append(bsonTags, tag)
	}
	return strings.Join(bsonTags, "."), nil
}

func (g RepositoryGenerator) bsonTagFromField(field code.StructField) (string, error) {
	bsonTag, ok := field.Tags["bson"]
	if !ok {
		return "", NewBsonTagNotFoundError(field.Name)
	}

	return bsonTag[0], nil
}

func (g RepositoryGenerator) repoImplStructName() string {
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
