package mongo

import (
	"fmt"
	"strings"

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

// RepositoryGenerator is a MongoDB repository generator that provides
// necessary information required to construct an implementation.
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
		Body: codegen.FunctionBody{
			codegen.ReturnStatement{
				codegen.StructStatement{
					Type: fmt.Sprintf("&%s", g.repoImplStructName()),
					Pairs: []codegen.StructFieldPair{{
						Key:   "collection",
						Value: codegen.Identifier("collection"),
					}},
				},
			},
		},
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

func (g RepositoryGenerator) generateMethodImplementation(
	methodSpec spec.MethodSpec) (codegen.FunctionBody, error) {

	switch operation := methodSpec.Operation.(type) {
	case spec.InsertOperation:
		return g.generateInsertBody(operation), nil
	case spec.FindOperation:
		return g.generateFindBody(operation)
	case spec.UpdateOperation:
		return g.generateUpdateBody(operation)
	case spec.DeleteOperation:
		return g.generateDeleteBody(operation)
	case spec.CountOperation:
		return g.generateCountBody(operation)
	default:
		return nil, NewOperationNotSupportedError(operation.Name())
	}
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
