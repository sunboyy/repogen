package mongo

import (
	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

func (g RepositoryGenerator) generateFindBody(
	operation spec.FindOperation) (codegen.FunctionBody, error) {

	querySpec, err := g.mongoQuerySpec(operation.Query)
	if err != nil {
		return nil, err
	}

	sortsCode, err := g.mongoSorts(operation.Sorts)
	if err != nil {
		return nil, err
	}

	if operation.Mode == spec.QueryModeOne {
		return g.generateFindOneBody(querySpec, sortsCode), nil
	}

	return g.generateFindManyBody(querySpec, sortsCode), nil
}

func (g RepositoryGenerator) generateFindOneBody(querySpec querySpec,
	sortsCode codegen.MapStatement) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclStatement{
			Name: "entity",
			Type: code.SimpleType(g.StructModel.Name),
		},
		codegen.IfBlock{
			Condition: []codegen.Statement{
				codegen.DeclAssignStatement{
					Vars: []string{"err"},
					Values: codegen.StatementList{
						codegen.ChainStatement{
							codegen.Identifier("r"),
							codegen.Identifier("collection"),
							codegen.CallStatement{
								FuncName: "FindOne",
								Params: codegen.StatementList{
									codegen.Identifier("arg0"),
									querySpec.Code(),
									codegen.ChainStatement{
										codegen.Identifier("options"),
										codegen.CallStatement{
											FuncName: "FindOne",
										},
										codegen.CallStatement{
											FuncName: "SetSort",
											Params: codegen.StatementList{
												sortsCode,
											},
										},
									},
								},
							},
							codegen.CallStatement{
								FuncName: "Decode",
								Params: codegen.StatementList{
									codegen.RawStatement("&entity"),
								},
							},
						},
					},
				},
				codegen.RawStatement("err != nil"),
			},
			Statements: []codegen.Statement{
				codegen.ReturnStatement{
					codegen.Identifier("nil"),
					codegen.Identifier("err"),
				},
			},
		},
		codegen.ReturnStatement{
			codegen.RawStatement("&entity"),
			codegen.Identifier("nil"),
		},
	}
}

func (g RepositoryGenerator) generateFindManyBody(querySpec querySpec,
	sortsCode codegen.MapStatement) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"cursor", "err"},
			Values: codegen.StatementList{
				codegen.ChainStatement{
					codegen.Identifier("r"),
					codegen.Identifier("collection"),
					codegen.CallStatement{
						FuncName: "Find",
						Params: codegen.StatementList{
							codegen.Identifier("arg0"),
							querySpec.Code(),
							codegen.ChainStatement{
								codegen.Identifier("options"),
								codegen.CallStatement{
									FuncName: "Find",
								},
								codegen.CallStatement{
									FuncName: "SetSort",
									Params: codegen.StatementList{
										sortsCode,
									},
								},
							},
						},
					},
				},
			},
		},
		codegen.IfBlock{
			Condition: []codegen.Statement{
				codegen.RawStatement("err != nil"),
			},
			Statements: []codegen.Statement{
				codegen.ReturnStatement{
					codegen.Identifier("nil"),
					codegen.Identifier("err"),
				},
			},
		},
		codegen.DeclStatement{
			Name: "entities",
			Type: code.ArrayType{
				ContainedType: code.PointerType{
					ContainedType: code.SimpleType(g.StructModel.Name),
				},
			},
		},
		codegen.IfBlock{
			Condition: []codegen.Statement{
				codegen.DeclAssignStatement{
					Vars: []string{"err"},
					Values: codegen.StatementList{
						codegen.ChainStatement{
							codegen.Identifier("cursor"),
							codegen.CallStatement{
								FuncName: "All",
								Params: codegen.StatementList{
									codegen.Identifier("arg0"),
									codegen.RawStatement("&entities"),
								},
							},
						},
					},
				},
				codegen.RawStatement("err != nil"),
			},
			Statements: []codegen.Statement{
				codegen.ReturnStatement{
					codegen.Identifier("nil"),
					codegen.Identifier("err"),
				},
			},
		},
		codegen.ReturnStatement{
			codegen.Identifier("entities"),
			codegen.Identifier("nil"),
		},
	}
}

func (g RepositoryGenerator) mongoSorts(sortSpec []spec.Sort) (
	codegen.MapStatement, error) {

	sortsCode := codegen.MapStatement{
		Type: "bson.M",
	}

	for _, s := range sortSpec {
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
