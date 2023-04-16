package mongo

import (
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

func (g RepositoryGenerator) generateDeleteBody(
	operation spec.DeleteOperation) (codegen.FunctionBody, error) {

	querySpec, err := g.mongoQuerySpec(operation.Query)
	if err != nil {
		return nil, err
	}

	if operation.Mode == spec.QueryModeOne {
		return g.generateDeleteOneBody(querySpec), nil
	}

	return g.generateDeleteManyBody(querySpec), nil
}

func (g RepositoryGenerator) generateDeleteOneBody(
	querySpec querySpec) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"result", "err"},
			Values: codegen.StatementList{
				codegen.ChainStatement{
					codegen.Identifier("r"),
					codegen.Identifier("collection"),
					codegen.CallStatement{
						FuncName: "DeleteOne",
						Params: codegen.StatementList{
							codegen.Identifier("arg0"),
							querySpec.Code(),
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
					codegen.Identifier("false"),
					codegen.Identifier("err"),
				},
			},
		},
		codegen.ReturnStatement{
			codegen.RawStatement("result.DeletedCount > 0"),
			codegen.Identifier("nil"),
		},
	}
}

func (g RepositoryGenerator) generateDeleteManyBody(
	querySpec querySpec) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"result", "err"},
			Values: codegen.StatementList{
				codegen.ChainStatement{
					codegen.Identifier("r"),
					codegen.Identifier("collection"),
					codegen.CallStatement{
						FuncName: "DeleteMany",
						Params: codegen.StatementList{
							codegen.Identifier("arg0"),
							querySpec.Code(),
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
					codegen.Identifier("0"),
					codegen.Identifier("err"),
				},
			},
		},
		codegen.ReturnStatement{
			codegen.CallStatement{
				FuncName: "int",
				Params: codegen.StatementList{
					codegen.ChainStatement{
						codegen.Identifier("result"),
						codegen.Identifier("DeletedCount"),
					},
				},
			},
			codegen.Identifier("nil"),
		},
	}
}
