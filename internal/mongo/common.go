package mongo

import (
	"go/token"
	"go/types"
	"strings"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

var (
	mongoCollectionType types.Type
	bsonMType           types.Type
)

func init() {
	bareMongoPkg := types.NewPackage("go.mongodb.org/mongo-driver/mongo", "mongo")
	mongoCollectionType = types.NewNamed(types.NewTypeName(token.NoPos, bareMongoPkg, "Collection", nil), nil, nil)

	bareBsonPkg := types.NewPackage("go.mongodb.org/mongo-driver/bson", "bson")
	bsonMType = types.NewNamed(types.NewTypeName(token.NoPos, bareBsonPkg, "M", nil), nil, nil)
}

var errOccurred = codegen.RawStatement("err != nil")

var returnNilErr = codegen.ReturnStatement{
	codegen.Identifier("nil"),
	codegen.Identifier("err"),
}

var ifErrReturnNilErr = codegen.IfBlock{
	Condition: []codegen.Statement{
		errOccurred,
	},
	Statements: []codegen.Statement{
		returnNilErr,
	},
}

var ifErrReturn0Err = codegen.IfBlock{
	Condition: []codegen.Statement{
		errOccurred,
	},
	Statements: []codegen.Statement{
		codegen.ReturnStatement{
			codegen.Identifier("0"),
			codegen.Identifier("err"),
		},
	},
}

var ifErrReturnFalseErr = codegen.IfBlock{
	Condition: []codegen.Statement{
		errOccurred,
	},
	Statements: []codegen.Statement{
		codegen.ReturnStatement{
			codegen.Identifier("false"),
			codegen.Identifier("err"),
		},
	},
}

type baseMethodGenerator struct {
	targetPkg        *types.Package
	structModelNamed *types.Named
}

func (g baseMethodGenerator) bsonFieldReference(fieldReference spec.FieldReference) (string, error) {
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

func (g baseMethodGenerator) bsonTagFromField(field code.StructField) (string, error) {
	bsonTag, ok := field.Tag.Lookup("bson")
	if !ok {
		return "", NewBsonTagNotFoundError(field.Var.Name())
	}

	documentKey := strings.Split(bsonTag, ",")[0]
	return documentKey, nil
}

func (g baseMethodGenerator) convertQuerySpec(query spec.QuerySpec) (querySpec, error) {
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
		TargetPkg:  g.targetPkg,
		Operator:   query.Operator,
		Predicates: predicates,
	}, nil
}
