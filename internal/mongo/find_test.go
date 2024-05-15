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

func TestGenerateMethod_Find(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "simple find one method",
			MethodSpec: spec.MethodSpec{
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
			ExpectedBody: `	var entity User
	if err := r.collection.FindOne(arg0, bson.M{
		"_id": arg1,
	}, options.FindOne().SetSort(bson.M{
	})).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil`,
		},
		{
			Name: "simple find many method",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGender",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": arg1,
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with deep field reference",
			MethodSpec: spec.MethodSpec{
				Name: "FindByNameFirst",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeString),
					},
					[]*types.Var{
						createTypeVar(types.NewPointer(testutils.TypeUserStruct)),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorEqual,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Name"),
										Tag: `bson:"name"`,
									},
									{
										Var: testutils.FindStructFieldByName(testutils.TypeNameStruct, "First"),
										Tag: `bson:"first"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	var entity User
	if err := r.collection.FindOne(arg0, bson.M{
		"name.first": arg1,
	}, options.FindOne().SetSort(bson.M{
	})).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil`,
		},
		{
			Name: "find with And operator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderAndAge",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"$and": []bson.M{
			{
				"gender": arg1,
			},
			{
				"age": arg2,
			},
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with Or operator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderOrAge",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"$or": []bson.M{
			{
				"gender": arg1,
			},
			{
				"age": arg2,
			},
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with Not comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderNot",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": bson.M{
			"$ne": arg1,
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with LessThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeLessThan",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{
			"$lt": arg1,
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with LessThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeLessThanEqual",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{
			"$lte": arg1,
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with GreaterThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeGreaterThan",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{
			"$gt": arg1,
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with GreaterThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeGreaterThanEqual",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{
			"$gte": arg1,
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with Between comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeBetween",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
						createTypeVar(code.TypeInt),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{
			"$gte": arg1,
			"$lte": arg2,
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with In comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderIn",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": bson.M{
			"$in": arg1,
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with NotIn comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderNotIn",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeGenderNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorNotIn,
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": bson.M{
			"$nin": arg1,
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with True comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByEnabledTrue",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorTrue,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Enabled"),
										Tag: `bson:"enabled"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"enabled": true,
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with False comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByEnabledFalse",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorFalse,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Enabled"),
										Tag: `bson:"enabled"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"enabled": false,
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with Exists comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByReferrerExists",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorExists,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Referrer"),
										Tag: `bson:"referrer"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"referrer": bson.M{
			"$exists": 1,
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with NotExists comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByReferrerNotExists",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator: spec.ComparatorNotExists,
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Referrer"),
										Tag: `bson:"referrer"`,
									},
								},
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
		"referrer": bson.M{
			"$exists": 0,
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with sort ascending",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByAge",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
	}, options.Find().SetSort(bson.M{
		"age": 1,
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with sort descending",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByAgeDesc",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
	}, options.Find().SetSort(bson.M{
		"age": -1,
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with deep sort ascending",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByNameFirst",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Name"),
									Tag: `bson:"name"`,
								},
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "First"),
									Tag: `bson:"first"`,
								},
							},
							Ordering: spec.OrderingAscending,
						},
					},
				},
			},
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
	}, options.Find().SetSort(bson.M{
		"name.first": 1,
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with multiple sorts",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByGenderAndAgeDesc",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "Gender"),
									Tag: `bson:"gender"`,
								},
							},
							Ordering: spec.OrderingAscending,
						},
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
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
	}, options.Find().SetSort(bson.M{
		"gender": 1,
		"age": -1,
	}))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with limit",
			MethodSpec: spec.MethodSpec{
				Name: "FindTop5AllOrderByAgeDesc",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserStruct))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
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
					Limit: 5,
				},
			},
			ExpectedBody: `	cursor, err := r.collection.Find(arg0, bson.M{
	}, options.Find().SetSort(bson.M{
		"age": -1,
	}).SetLimit(5))
	if err != nil {
		return nil, err
	}
	entities := []*User{
	}
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
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
