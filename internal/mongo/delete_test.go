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

func TestGenerateMethod_Delete(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "simple delete one method",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeBool),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorEqual,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
										Tag: `bson:"_id,omitempty"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteOne(arg0, bson.M{
		"_id": arg1,
	})
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil`,
		},
		{
			Name: "simple delete many method",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByGender",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorEqual,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
										Tag: `bson:"gender"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteMany(arg0, bson.M{
		"gender": arg1,
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`,
		},
		{
			Name: "delete with And operator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByGenderAndAge",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Operator: spec.OperatorAnd,
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorEqual,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
										Tag: `bson:"gender"`,
									},
								},
								ParamIndex: 1,
							},
							{
								Comparator: spec.ComparatorEqual,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteMany(arg0, bson.M{
		"$and": []bson.M{
			{
				"gender": arg1,
			},
			{
				"age": arg2,
			},
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`,
		},
		{
			Name: "delete with Or operator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByGenderOrAge",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Operator: spec.OperatorOr,
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorEqual,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
										Tag: `bson:"gender"`,
									},
								},
								ParamIndex: 1,
							},
							{
								Comparator: spec.ComparatorEqual,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteMany(arg0, bson.M{
		"$or": []bson.M{
			{
				"gender": arg1,
			},
			{
				"age": arg2,
			},
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`,
		},
		{
			Name: "delete with Not comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByGenderNot",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorNot,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
										Tag: `bson:"gender"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteMany(arg0, bson.M{
		"gender": bson.M{
			"$ne": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`,
		},
		{
			Name: "delete with LessThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByAgeLessThan",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorLessThan,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteMany(arg0, bson.M{
		"age": bson.M{
			"$lt": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`,
		},
		{
			Name: "delete with LessThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByAgeLessThanEqual",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorLessThanEqual,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteMany(arg0, bson.M{
		"age": bson.M{
			"$lte": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`,
		},
		{
			Name: "delete with GreaterThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByAgeGreaterThan",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorGreaterThan,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteMany(arg0, bson.M{
		"age": bson.M{
			"$gt": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`,
		},
		{
			Name: "delete with GreaterThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByAgeGreaterThanEqual",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorGreaterThanEqual,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteMany(arg0, bson.M{
		"age": bson.M{
			"$gte": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`,
		},
		{
			Name: "delete with Between comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByAgeBetween",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorBetween,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteMany(arg0, bson.M{
		"age": bson.M{
			"$gte": arg1,
			"$lte": arg2,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`,
		},
		{
			Name: "delete with In comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByGenderIn",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorIn,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
										Tag: `bson:"gender"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	result, err := r.collection.DeleteMany(arg0, bson.M{
		"gender": bson.M{
			"$in": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(testutils.Pkg, "User", "UserRepository")
			expectedReceiver := codegen.MethodReceiver{
				Name:    "r",
				Type:    "UserRepositoryMongo",
				Pointer: true,
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
