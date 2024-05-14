package generator_test

import (
	"go/token"
	"go/types"
	"os"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/generator"
	"github.com/sunboyy/repogen/internal/spec"
	"github.com/sunboyy/repogen/internal/testutils"
)

func createSignature(params []*types.Var, results []*types.Var) *types.Signature {
	return types.NewSignatureType(nil, nil, nil, types.NewTuple(params...), types.NewTuple(results...), false)
}

func createTypeVar(t types.Type) *types.Var {
	return types.NewVar(token.NoPos, nil, "", t)
}

func TestGenerateMongoRepository(t *testing.T) {
	methods := []spec.MethodSpec{
		// test find: One mode
		{
			Name: "FindByID",
			Signature: createSignature(
				[]*types.Var{
					createTypeVar(testutils.TypeContextNamed),
					createTypeVar(testutils.TypeObjectIDNamed),
				},
				[]*types.Var{
					createTypeVar(types.NewPointer(testutils.TypeUserNamed)),
					createTypeVar(code.TypeError),
				},
			),
			Operation: spec.FindOperation{
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
							ParamIndex: 1,
						},
					},
				},
			},
		},
		// test find: Many mode, And operator, NOT and LessThan comparator
		{
			Name: "FindByGenderNotAndAgeLessThan",
			Signature: createSignature(
				[]*types.Var{
					createTypeVar(testutils.TypeContextNamed),
					createTypeVar(testutils.TypeGenderNamed),
					createTypeVar(code.TypeInt),
				},
				[]*types.Var{
					createTypeVar(types.NewPointer(testutils.TypeUserNamed)),
					createTypeVar(code.TypeError),
				},
			),
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
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
							Comparator: spec.ComparatorNot,
							ParamIndex: 1,
						},
						{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
									Tag: `bson:"age"`,
								},
							},
							Comparator: spec.ComparatorLessThan,
							ParamIndex: 2,
						},
					},
				},
			},
		},
		{
			Name: "FindByAgeLessThanEqualOrderByAge",
			Signature: createSignature(
				[]*types.Var{
					createTypeVar(testutils.TypeContextNamed),
					createTypeVar(code.TypeInt),
				},
				[]*types.Var{
					createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserNamed))),
					createTypeVar(code.TypeError),
				},
			),
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
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
				Sorts: []spec.Sort{
					{
						FieldReference: spec.FieldReference{
							{
								Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
								Tag: `bson:"age"`,
							},
						},
						Ordering: spec.OrderingAscending,
					},
				},
			},
		},
		{
			Name: "FindByAgeGreaterThanOrderByAgeAsc",
			Signature: createSignature(
				[]*types.Var{
					createTypeVar(testutils.TypeContextNamed),
					createTypeVar(code.TypeInt),
				},
				[]*types.Var{
					createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserNamed))),
					createTypeVar(code.TypeError),
				},
			),
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
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
				Sorts: []spec.Sort{
					{
						FieldReference: spec.FieldReference{
							{
								Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
								Tag: `bson:"age"`,
							},
						},
						Ordering: spec.OrderingAscending,
					},
				},
			},
		},
		{
			Name: "FindByAgeGreaterThanEqualOrderByAgeDesc",
			Signature: createSignature(
				[]*types.Var{
					createTypeVar(testutils.TypeContextNamed),
					createTypeVar(code.TypeInt),
				},
				[]*types.Var{
					createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserNamed))),
					createTypeVar(code.TypeError),
				},
			),
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
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
				Sorts: []spec.Sort{
					{
						FieldReference: spec.FieldReference{
							{
								Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Age"),
								Tag: `bson:"age"`,
							},
						},
						Ordering: spec.OrderingDescending,
					},
				},
			},
		},
		{
			Name: "FindByAgeBetween",
			Signature: createSignature(
				[]*types.Var{
					createTypeVar(testutils.TypeContextNamed),
					createTypeVar(code.TypeInt),
					createTypeVar(code.TypeInt),
				},
				[]*types.Var{
					createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserNamed))),
					createTypeVar(code.TypeError),
				},
			),
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
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
		{
			Name: "FindByGenderOrAge",
			Signature: createSignature(
				[]*types.Var{
					createTypeVar(testutils.TypeContextNamed),
					createTypeVar(testutils.TypeGenderNamed),
					createTypeVar(code.TypeInt),
				},
				[]*types.Var{
					createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserNamed))),
					createTypeVar(code.TypeError),
				},
			),
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
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
	}
	expectedBytes, err := os.ReadFile("../../test/generator_test_expected.txt")
	if err != nil {
		t.Fatal(err)
	}
	expectedCode := string(expectedBytes)

	code, err := generator.GenerateRepository(testutils.Pkg, "User", "UserRepository", methods)

	if err != nil {
		t.Fatal(err)
	}
	if err := testutils.ExpectMultiLineString(expectedCode, code); err != nil {
		t.Error(err)
	}
}
