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

func TestGenerateMethod_Count(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "simple count method",
			MethodSpec: spec.MethodSpec{
				Name: "CountByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{genderField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     1,
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
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Operator: spec.OperatorAnd,
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{genderField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     1,
							},
							{
								FieldReference: spec.FieldReference{ageField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     2,
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
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Operator: spec.OperatorOr,
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{genderField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     1,
							},
							{
								FieldReference: spec.FieldReference{ageField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     2,
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
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{genderField},
								Comparator:     spec.ComparatorNot,
								ParamIndex:     1,
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
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{ageField},
								Comparator:     spec.ComparatorLessThan,
								ParamIndex:     1,
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
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{ageField},
								Comparator:     spec.ComparatorLessThanEqual,
								ParamIndex:     1,
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
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{ageField},
								Comparator:     spec.ComparatorGreaterThan,
								ParamIndex:     1,
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
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{ageField},
								Comparator:     spec.ComparatorGreaterThanEqual,
								ParamIndex:     1,
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
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{ageField},
								Comparator:     spec.ComparatorBetween,
								ParamIndex:     1,
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
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.TypeInt}},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{ageField},
								Comparator:     spec.ComparatorIn,
								ParamIndex:     1,
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
