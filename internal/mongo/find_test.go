package mongo_test

import (
	"fmt"
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorEqual,
								FieldReference: spec.FieldReference{idField},
								ParamIndex:     1,
							},
						},
					},
				},
			},
			ExpectedBody: `	var entity UserModel
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorEqual,
								FieldReference: spec.FieldReference{genderField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with deep field reference",
			MethodSpec: spec.MethodSpec{
				Name: "FindByNameFirst",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "firstName", Type: code.TypeString},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorEqual,
								FieldReference: spec.FieldReference{nameField, firstNameField},
								ParamIndex:     1,
							},
						},
					},
				},
			},
			ExpectedBody: `	var entity UserModel
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Operator: spec.OperatorAnd,
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorEqual,
								FieldReference: spec.FieldReference{genderField},
								ParamIndex:     1,
							},
							{
								Comparator:     spec.ComparatorEqual,
								FieldReference: spec.FieldReference{ageField},
								ParamIndex:     2,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with Or operator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderOrAge",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Operator: spec.OperatorOr,
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorEqual,
								FieldReference: spec.FieldReference{genderField},
								ParamIndex:     1,
							},
							{
								Comparator:     spec.ComparatorEqual,
								FieldReference: spec.FieldReference{ageField},
								ParamIndex:     2,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with Not comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderNot",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorNot,
								FieldReference: spec.FieldReference{genderField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with LessThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeLessThan",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorLessThan,
								FieldReference: spec.FieldReference{ageField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with LessThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeLessThanEqual",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorLessThanEqual,
								FieldReference: spec.FieldReference{ageField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with GreaterThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeGreaterThan",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorGreaterThan,
								FieldReference: spec.FieldReference{ageField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with GreaterThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeGreaterThanEqual",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorGreaterThanEqual,
								FieldReference: spec.FieldReference{ageField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with Between comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeBetween",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "fromAge", Type: code.TypeInt},
					{Name: "toAge", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorBetween,
								FieldReference: spec.FieldReference{ageField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with In comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderIn",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.ArrayType{ContainedType: code.SimpleType("Gender")}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorIn,
								FieldReference: spec.FieldReference{genderField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with NotIn comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderNotIn",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.ArrayType{ContainedType: code.SimpleType("Gender")}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorNotIn,
								FieldReference: spec.FieldReference{genderField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with True comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByEnabledTrue",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorTrue,
								FieldReference: spec.FieldReference{enabledField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with False comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByEnabledFalse",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorFalse,
								FieldReference: spec.FieldReference{enabledField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with Exists comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByReferrerExists",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorExists,
								FieldReference: spec.FieldReference{referrerField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with NotExists comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByReferrerNotExists",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorNotExists,
								FieldReference: spec.FieldReference{referrerField},
								ParamIndex:     1,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with sort ascending",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByAge",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingAscending},
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with sort descending",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByAgeDesc",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingDescending},
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with deep sort ascending",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByNameFirst",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{
							FieldReference: spec.FieldReference{nameField, firstNameField},
							Ordering:       spec.OrderingAscending,
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with multiple sorts",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByGenderAndAgeDesc",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{FieldReference: spec.FieldReference{genderField}, Ordering: spec.OrderingAscending},
						{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingDescending},
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
		{
			Name: "find with limit",
			MethodSpec: spec.MethodSpec{
				Name: "FindTop5AllOrderByAgeDesc",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingDescending},
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
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			expectedReceiver := codegen.MethodReceiver{
				Name:    "r",
				Type:    "UserRepositoryMongo",
				Pointer: true,
			}
			var expectedParams []code.Param
			for i, param := range testCase.MethodSpec.Params {
				expectedParams = append(expectedParams, code.Param{
					Name: fmt.Sprintf("arg%d", i),
					Type: param.Type,
				})
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
			if !reflect.DeepEqual(testCase.MethodSpec.Returns, actual.Returns) {
				t.Errorf(
					"incorrect struct returns: expected %+v, got %+v",
					testCase.MethodSpec.Returns,
					actual.Returns,
				)
			}
			if err := testutils.ExpectMultiLineString(testCase.ExpectedBody, actual.Body.Code()); err != nil {
				t.Error(err)
			}
		})
	}
}
