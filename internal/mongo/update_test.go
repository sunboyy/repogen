package mongo_test

import (
	"fmt"
	"go/token"
	"go/types"
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
	"github.com/sunboyy/repogen/internal/testutils"
)

func TestGenerateMethod_Update(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "update model method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeUserNamed),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeBool),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.UpdateOperation{
					Update: spec.UpdateModel{},
					Mode:   spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
										Tag: `bson:"_id,omitempty"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg2,
	}, bson.M{
		"$set": arg1,
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, nil`,
		},
		{
			Name: "simple update one method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateAgeByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeBool),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
									Tag: `bson:"age"`,
								},
							},
							ParamIndex: 1,
							Operator:   spec.UpdateOperatorSet,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
										Tag: `bson:"_id,omitempty"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg2,
	}, bson.M{
		"$set": bson.M{
			"age": arg1,
		},
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, nil`,
		},
		{
			Name: "simple update many method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateAgeByGender",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
						createTypeVar(testutils.TypeGenderNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
									Tag: `bson:"age"`,
								},
							},
							ParamIndex: 1,
							Operator:   spec.UpdateOperatorSet,
						},
					},
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
										Tag: `bson:"gender"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.UpdateMany(arg0, bson.M{
		"gender": arg2,
	}, bson.M{
		"$set": bson.M{
			"age": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.MatchedCount), nil`,
		},
		{
			Name: "simple update push method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateConsentHistoryPushByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeConsentHistoryNamed),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeBool),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ConsentHistory"),
									Tag: `bson:"consent_history"`,
								},
							},
							ParamIndex: 1,
							Operator:   spec.UpdateOperatorPush,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
										Tag: `bson:"_id,omitempty"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg2,
	}, bson.M{
		"$push": bson.M{
			"consent_history": arg1,
		},
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, nil`,
		},
		{
			Name: "simple update inc method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateAgeIncByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeBool),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
									Tag: `bson:"age"`,
								},
							},
							ParamIndex: 1,
							Operator:   spec.UpdateOperatorInc,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
										Tag: `bson:"_id,omitempty"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg2,
	}, bson.M{
		"$inc": bson.M{
			"age": arg1,
		},
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, nil`,
		},
		{
			Name: "simple update set and push method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateEnabledAndConsentHistoryPushByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeBool),
						createTypeVar(testutils.TypeConsentHistoryNamed),
						createTypeVar(testutils.TypeGenderNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeBool),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Enabled"),
									Tag: `bson:"enabled"`,
								},
							},
							ParamIndex: 1,
							Operator:   spec.UpdateOperatorSet,
						},
						spec.UpdateField{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ConsentHistory"),
									Tag: `bson:"consent_history"`,
								},
							},
							ParamIndex: 2,
							Operator:   spec.UpdateOperatorPush,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
										Tag: `bson:"_id,omitempty"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 3,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg3,
	}, bson.M{
		"$push": bson.M{
			"consent_history": arg2,
		},
		"$set": bson.M{
			"enabled": arg1,
		},
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, nil`,
		},
		{
			Name: "update with deeply referenced field",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateNameFirstByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeString),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeBool),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Name"),
									Tag: `bson:"name"`,
								},
								{
									Var: testutils.FindStructFieldByName(testutils.TypeNameStruct, "FirstName"),
									Tag: `bson:"first"`,
								},
							},
							ParamIndex: 1,
							Operator:   spec.UpdateOperatorSet,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
										Tag: `bson:"_id,omitempty"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg2,
	}, bson.M{
		"$set": bson.M{
			"name.first": arg1,
		},
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, nil`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(testutils.Pkg, testutils.TypeUserNamed, "UserRepository")
			expectedReceiver := codegen.MethodReceiver{
				Name:     "r",
				TypeName: "UserRepositoryMongo",
				Pointer:  true,
			}

			params := testCase.MethodSpec.Signature.Params()
			var expectedParamVars []*types.Var
			for i := 0; i < params.Len(); i++ {
				expectedParamVars = append(expectedParamVars, types.NewVar(token.NoPos, nil, fmt.Sprintf("arg%d", i),
					params.At(i).Type()))
			}
			expectedParams := types.NewTuple(expectedParamVars...)
			returns := testCase.MethodSpec.Signature.Results()
			var expectedReturns []types.Type
			for i := 0; i < returns.Len(); i++ {
				expectedReturns = append(expectedReturns, returns.At(i).Type())
			}

			actual, err := generator.GenerateMethod(testCase.MethodSpec)

			if err != nil {
				t.Fatal(err)
			}
			if expectedReceiver != actual.Receiver {
				t.Errorf(
					"incorrect method receiver: expected %+v, got %+v",
					expectedReceiver,
					actual.Receiver,
				)
			}
			if testCase.MethodSpec.Name != actual.Name {
				t.Errorf(
					"incorrect method name: expected %s, got %s",
					testCase.MethodSpec.Name,
					actual.Name,
				)
			}
			if !reflect.DeepEqual(expectedParams, actual.Params) {
				t.Errorf(
					"incorrect struct params: expected %+v, got %+v",
					expectedParams,
					actual.Params,
				)
			}
			if !reflect.DeepEqual(expectedReturns, actual.Returns) {
				t.Errorf(
					"incorrect struct returns: expected %+v, got %+v",
					expectedReturns,
					actual.Returns,
				)
			}
			if err := testutils.ExpectMultiLineString(testCase.ExpectedBody, actual.Body.Code()); err != nil {
				t.Error(err)
			}
		})
	}
}
