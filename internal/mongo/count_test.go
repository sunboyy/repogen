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

func TestGenerateMethod_Count(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "simple count method",
			MethodSpec: spec.MethodSpec{
				Name: "CountByGender",
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
				Operation: spec.CountOperation{
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
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	count, err := r.collection.CountDocuments(arg0, bson.M{
		"gender": arg1,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil`,
		},
		{
			Name: "count with And operator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByGenderAndCity",
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
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Operator: spec.OperatorAnd,
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
										Tag: `bson:"gender"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 1,
							},
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedBody: `	count, err := r.collection.CountDocuments(arg0, bson.M{
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
	return int(count), nil`,
		},
		{
			Name: "count with Or operator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByGenderOrCity",
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
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Operator: spec.OperatorOr,
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
										Tag: `bson:"gender"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 1,
							},
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedBody: `	count, err := r.collection.CountDocuments(arg0, bson.M{
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
	return int(count), nil`,
		},
		{
			Name: "count with Not comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByGenderNot",
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
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
										Tag: `bson:"gender"`,
									},
								},
								Comparator: spec.ComparatorNot,
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	count, err := r.collection.CountDocuments(arg0, bson.M{
		"gender": bson.M{
			"$ne": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil`,
		},
		{
			Name: "count with LessThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeLessThan",
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
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								Comparator: spec.ComparatorLessThan,
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{
			"$lt": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil`,
		},
		{
			Name: "count with LessThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeLessThanEqual",
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
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								Comparator: spec.ComparatorLessThanEqual,
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{
			"$lte": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil`,
		},
		{
			Name: "count with GreaterThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeGreaterThan",
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
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								Comparator: spec.ComparatorGreaterThan,
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{
			"$gt": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil`,
		},
		{
			Name: "count with GreaterThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeGreaterThanEqual",
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
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								Comparator: spec.ComparatorGreaterThanEqual,
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{
			"$gte": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil`,
		},
		{
			Name: "count with Between comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeBetween",
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
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								Comparator: spec.ComparatorBetween,
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{
			"$gte": arg1,
			"$lte": arg2,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil`,
		},
		{
			Name: "count with In comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeIn",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
										Tag: `bson:"age"`,
									},
								},
								Comparator: spec.ComparatorIn,
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{
			"$in": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil`,
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
