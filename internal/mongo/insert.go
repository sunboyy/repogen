package mongo

import (
	"go/types"

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
				codegen.NewChainBuilder("r").
					Chain("collection").
					Call("InsertOne",
						codegen.Identifier("arg0"),
						codegen.Identifier("arg1"),
					).Build(),
			},
		},
		ifErrReturnNilErr,
		codegen.ReturnStatement{
			codegen.NewChainBuilder("result").Chain("InsertedID").Build(),
			codegen.Identifier("nil"),
		},
	}
}

func (g RepositoryGenerator) generateInsertManyBody() codegen.FunctionBody {
	return codegen.FunctionBody{
		codegen.NewDeclStatement(
			g.targetPkg,
			"entities",
			types.NewSlice(types.NewInterfaceType(nil, nil)),
		),
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
				codegen.NewChainBuilder("r").
					Chain("collection").
					Call("InsertMany",
						codegen.Identifier("arg0"),
						codegen.Identifier("entities"),
					).Build(),
			},
		},
		ifErrReturnNilErr,
		codegen.ReturnStatement{
			codegen.NewChainBuilder("result").Chain("InsertedIDs").Build(),
			codegen.Identifier("nil"),
		},
	}
}
