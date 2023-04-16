package mongo

import (
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

func (g RepositoryGenerator) generateUpdateBody(
	operation spec.UpdateOperation) (codegen.FunctionBody, error) {

	update, err := g.getMongoUpdate(operation.Update)
	if err != nil {
		return nil, err
	}

	querySpec, err := g.mongoQuerySpec(operation.Query)
	if err != nil {
		return nil, err
	}

	if operation.Mode == spec.QueryModeOne {
		return g.generateUpdateOneBody(update, querySpec), nil
	}

	return g.generateUpdateManyBody(update, querySpec), nil
}

func (g RepositoryGenerator) generateUpdateOneBody(update update,
	querySpec querySpec) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"result", "err"},
			Values: codegen.StatementList{
				codegen.ChainStatement{
					codegen.Identifier("r"),
					codegen.Identifier("collection"),
					codegen.CallStatement{
						FuncName: "UpdateOne",
						Params: codegen.StatementList{
							codegen.Identifier("arg0"),
							querySpec.Code(),
							update.Code(),
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
			codegen.RawStatement("result.MatchedCount > 0"),
			codegen.Identifier("nil"),
		},
	}
}

func (g RepositoryGenerator) generateUpdateManyBody(update update,
	querySpec querySpec) codegen.FunctionBody {

	return codegen.FunctionBody{
		codegen.DeclAssignStatement{
			Vars: []string{"result", "err"},
			Values: codegen.StatementList{
				codegen.ChainStatement{
					codegen.Identifier("r"),
					codegen.Identifier("collection"),
					codegen.CallStatement{
						FuncName: "UpdateMany",
						Params: codegen.StatementList{
							codegen.Identifier("arg0"),
							querySpec.Code(),
							update.Code(),
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
						codegen.Identifier("MatchedCount"),
					},
				},
			},
			codegen.Identifier("nil"),
		},
	}
}

func (g RepositoryGenerator) getMongoUpdate(updateSpec spec.Update) (update, error) {
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
