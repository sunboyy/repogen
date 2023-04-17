package mongo

import (
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

func (g RepositoryGenerator) generateCountBody(
	operation spec.CountOperation) (codegen.FunctionBody, error) {

	querySpec, err := g.convertQuerySpec(operation.Query)
	if err != nil {
		return nil, err
	}

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"count", "err"},
			Values: codegen.StatementList{
				codegen.NewChainBuilder("r").
					Chain("collection").
					Call("CountDocuments",
						codegen.Identifier("arg0"),
						querySpec.Code(),
					).Build(),
			},
		},
		ifErrReturn0Err,
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
