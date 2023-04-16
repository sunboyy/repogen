package mongo

import (
	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

func (g RepositoryGenerator) generateInsertBody(
	operation spec.InsertOperation) codegen.FunctionBody {

	if operation.Mode == spec.QueryModeOne {
		return g.generateInsertOneBody()
	}
	return g.generateInsertManyBody()
}

func (g RepositoryGenerator) generateInsertOneBody() codegen.FunctionBody {
	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"result", "err"},
			Values: codegen.StatementList{
				codegen.ChainStatement{
					codegen.Identifier("r"),
					codegen.Identifier("collection"),
					codegen.CallStatement{
						FuncName: "InsertOne",
						Params: []codegen.Statement{
							codegen.Identifier("arg0"),
							codegen.Identifier("arg1"),
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
		codegen.ReturnStatement{
			codegen.ChainStatement{
				codegen.Identifier("result"),
				codegen.Identifier("InsertedID"),
			},
			codegen.Identifier("nil"),
		},
	}
}

func (g RepositoryGenerator) generateInsertManyBody() codegen.FunctionBody {
	return codegen.FunctionBody{
		codegen.DeclStatement{
			Name: "entities",
			Type: code.ArrayType{
				ContainedType: code.InterfaceType{},
			},
		},
		codegen.RawBlock{
			Header: []string{"for _, model := range arg1"},
			Statements: []codegen.Statement{
				codegen.AssignStatement{
					Vars: []string{"entities"},
					Values: codegen.StatementList{
						codegen.CallStatement{
							FuncName: "append",
							Params: codegen.StatementList{
								codegen.Identifier("entities"),
								codegen.Identifier("model"),
							},
						},
					},
				},
			},
		},
		codegen.DeclAssignStatement{
			Vars: []string{"result", "err"},
			Values: codegen.StatementList{
				codegen.ChainStatement{
					codegen.Identifier("r"),
					codegen.Identifier("collection"),
					codegen.CallStatement{
						FuncName: "InsertMany",
						Params: codegen.StatementList{
							codegen.Identifier("arg0"),
							codegen.Identifier("entities"),
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
		codegen.ReturnStatement{
			codegen.ChainStatement{
				codegen.Identifier("result"),
				codegen.Identifier("InsertedIDs"),
			},
			codegen.Identifier("nil"),
		},
	}
}
