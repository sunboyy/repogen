package mongo

import (
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

func (g RepositoryGenerator) generateUpdateBody(
	operation spec.UpdateOperation) (codegen.FunctionBody, error) {

	return updateBodyGenerator{
		baseMethodGenerator: g.baseMethodGenerator,
		operation:           operation,
	}.generate()
}

type updateBodyGenerator struct {
	baseMethodGenerator
	operation spec.UpdateOperation
}

func (g updateBodyGenerator) generate() (codegen.FunctionBody, error) {
	update, err := g.convertUpdate(g.operation.Update)
	if err != nil {
		return nil, err
	}

	querySpec, err := g.convertQuerySpec(g.operation.Query)
	if err != nil {
		return nil, err
	}

	if g.operation.Mode == spec.QueryModeOne {
		return g.generateUpdateOneBody(update, querySpec), nil
	}

	return g.generateUpdateManyBody(update, querySpec), nil
}

func (g updateBodyGenerator) generateUpdateOneBody(update update,
	querySpec querySpec) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"result", "err"},
			Values: codegen.StatementList{
				codegen.NewChainBuilder("r").
					Chain("collection").
					Call("UpdateOne",
						codegen.Identifier("arg0"),
						querySpec.Code(),
						update.Code(),
					).Build(),
			},
		},
		ifErrReturnFalseErr,
		codegen.ReturnStatement{
			codegen.RawStatement("result.MatchedCount > 0"),
			codegen.Identifier("nil"),
		},
	}
}

func (g updateBodyGenerator) generateUpdateManyBody(update update,
	querySpec querySpec) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"result", "err"},
			Values: codegen.StatementList{
				codegen.NewChainBuilder("r").
					Chain("collection").
					Call("UpdateMany",
						codegen.Identifier("arg0"),
						querySpec.Code(),
						update.Code(),
					).Build(),
			},
		},
		ifErrReturn0Err,
		codegen.ReturnStatement{
			codegen.CallStatement{
				FuncName: "int",
				Params: codegen.StatementList{
					codegen.NewChainBuilder("result").
						Chain("MatchedCount").Build(),
				},
			},
			codegen.Identifier("nil"),
		},
	}
}

func (g updateBodyGenerator) convertUpdate(updateSpec spec.Update) (update, error) {
	switch updateSpec := updateSpec.(type) {
	case spec.UpdateModel:
		return updateModel{}, nil
	case spec.UpdateFields:
		update := make(updateFields)
		for _, field := range updateSpec {
			bsonFieldReference, err := g.bsonFieldReference(field.FieldReference)
			if err != nil {
				return nil, err
			}

			updateKey := getUpdateOperatorKey(field.Operator)
			if updateKey == "" {
				return nil, NewUpdateOperatorNotSupportedError(field.Operator)
			}
			updateField := updateField{
				BsonTag:    bsonFieldReference,
				ParamIndex: field.ParamIndex,
			}
			update[updateKey] = append(update[updateKey], updateField)
		}
		return update, nil
	default:
		return nil, NewUpdateTypeNotSupportedError(updateSpec)
	}
}

func getUpdateOperatorKey(operator spec.UpdateOperator) string {
	switch operator {
	case spec.UpdateOperatorSet:
		return "$set"
	case spec.UpdateOperatorPush:
		return "$push"
	case spec.UpdateOperatorInc:
		return "$inc"
	default:
		return ""
	}
}
