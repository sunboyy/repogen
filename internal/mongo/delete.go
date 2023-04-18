package mongo

import (
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

func (g RepositoryGenerator) generateDeleteBody(
	operation spec.DeleteOperation) (codegen.FunctionBody, error) {

	return deleteBodyGenerator{
		baseMethodGenerator: g.baseMethodGenerator,
		operation:           operation,
	}.generate()
}

type deleteBodyGenerator struct {
	baseMethodGenerator
	operation spec.DeleteOperation
}

func (g deleteBodyGenerator) generate() (codegen.FunctionBody, error) {
	querySpec, err := g.convertQuerySpec(g.operation.Query)
	if err != nil {
		return nil, err
	}

	if g.operation.Mode == spec.QueryModeOne {
		return g.generateDeleteOneBody(querySpec), nil
	}

	return g.generateDeleteManyBody(querySpec), nil
}

func (g deleteBodyGenerator) generateDeleteOneBody(
	querySpec querySpec) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"result", "err"},
			Values: codegen.StatementList{
				codegen.NewChainBuilder("r").
					Chain("collection").
					Call("DeleteOne",
						codegen.Identifier("arg0"),
						querySpec.Code(),
					).Build(),
			},
		},
		ifErrReturnFalseErr,
		codegen.ReturnStatement{
			codegen.RawStatement("result.DeletedCount > 0"),
			codegen.Identifier("nil"),
		},
	}
}

func (g deleteBodyGenerator) generateDeleteManyBody(
	querySpec querySpec) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"result", "err"},
			Values: codegen.StatementList{
				codegen.NewChainBuilder("r").
					Chain("collection").
					Call("DeleteMany",
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
					codegen.NewChainBuilder("result").Chain("DeletedCount").Build(),
				},
			},
			codegen.Identifier("nil"),
		},
	}
}
