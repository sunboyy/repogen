package mongo

import (
	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

func (g RepositoryGenerator) generateFindBody(
	operation spec.FindOperation) (codegen.FunctionBody, error) {

	return findBodyGenerator{
		baseMethodGenerator: g.baseMethodGenerator,
		operation:           operation,
	}.generate()
}

type findBodyGenerator struct {
	baseMethodGenerator
	operation spec.FindOperation
}

func (g findBodyGenerator) generate() (codegen.FunctionBody, error) {
	querySpec, err := g.convertQuerySpec(g.operation.Query)
	if err != nil {
		return nil, err
	}

	sortsCode, err := g.generateSortMap()
	if err != nil {
		return nil, err
	}

	if g.operation.Mode == spec.QueryModeOne {
		return g.generateFindOneBody(querySpec, sortsCode), nil
	}

	return g.generateFindManyBody(querySpec, sortsCode), nil
}

func (g findBodyGenerator) generateFindOneBody(querySpec querySpec,
	sortsCode codegen.MapStatement) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclStatement{
			Name: "entity",
			Type: code.SimpleType(g.structModel.Name),
		},
		codegen.IfBlock{
			Condition: []codegen.Statement{
				codegen.DeclAssignStatement{
					Vars: []string{"err"},
					Values: codegen.StatementList{
						codegen.NewChainBuilder("r").
							Chain("collection").
							Call("FindOne",
								codegen.Identifier("arg0"),
								querySpec.Code(),
								codegen.NewChainBuilder("options").
									Call("FindOne").
									Call("SetSort", sortsCode).
									Build(),
							).
							Call("Decode",
								codegen.RawStatement("&entity"),
							).Build(),
					},
				},
				codegen.RawStatement("err != nil"),
			},
			Statements: []codegen.Statement{
				returnNilErr,
			},
		},
		codegen.ReturnStatement{
			codegen.RawStatement("&entity"),
			codegen.Identifier("nil"),
		},
	}
}

func (g findBodyGenerator) generateFindManyBody(querySpec querySpec,
	sortsCode codegen.MapStatement) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"cursor", "err"},
			Values: codegen.StatementList{
				codegen.NewChainBuilder("r").
					Chain("collection").
					Call("Find",
						codegen.Identifier("arg0"),
						querySpec.Code(),
						codegen.NewChainBuilder("options").
							Call("Find").
							Call("SetSort", sortsCode).
							Build(),
					).Build(),
			},
		},
		ifErrReturnNilErr,
		codegen.DeclStatement{
			Name: "entities",
			Type: code.ArrayType{
				ContainedType: code.PointerType{
					ContainedType: code.SimpleType(g.structModel.Name),
				},
			},
		},
		codegen.IfBlock{
			Condition: []codegen.Statement{
				codegen.DeclAssignStatement{
					Vars: []string{"err"},
					Values: codegen.StatementList{
						codegen.NewChainBuilder("cursor").
							Call("All",
								codegen.Identifier("arg0"),
								codegen.RawStatement("&entities"),
							).Build(),
					},
				},
				codegen.RawStatement("err != nil"),
			},
			Statements: []codegen.Statement{
				returnNilErr,
			},
		},
		codegen.ReturnStatement{
			codegen.Identifier("entities"),
			codegen.Identifier("nil"),
		},
	}
}

func (g findBodyGenerator) generateSortMap() (
	codegen.MapStatement, error) {

	sortsCode := codegen.MapStatement{
		Type: "bson.M",
	}

	for _, s := range g.operation.Sorts {
		bsonFieldReference, err := g.bsonFieldReference(s.FieldReference)
		if err != nil {
			return codegen.MapStatement{}, err
		}

		sortValueIdentifier := codegen.Identifier("1")
		if s.Ordering == spec.OrderingDescending {
			sortValueIdentifier = codegen.Identifier("-1")
		}

		sortsCode.Pairs = append(sortsCode.Pairs, codegen.MapPair{
			Key:   bsonFieldReference,
			Value: sortValueIdentifier,
		})
	}

	return sortsCode, nil
}
