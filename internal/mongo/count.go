package mongo

import (
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

func (g RepositoryGenerator) generateCountBody(
	operation spec.CountOperation) (codegen.FunctionBody, error) {

	querySpec, err := g.mongoQuerySpec(operation.Query)
	if err != nil {
		return nil, err
	}

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"count", "err"},
			Values: codegen.StatementList{
				codegen.ChainStatement{
					codegen.Identifier("r"),
					codegen.Identifier("collection"),
					codegen.CallStatement{
						FuncName: "CountDocuments",
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
					codegen.Identifier("count"),
				},
			},
			codegen.Identifier("nil"),
		},
	}, nil
}
