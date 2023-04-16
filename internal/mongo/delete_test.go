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

func TestGenerateMethod_Delete(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "simple delete one method",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByID",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{code.TypeBool, code.TypeError},
				Operation: spec.DeleteOperation{
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.DeleteOperation{
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.DeleteOperation{
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.DeleteOperation{
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.DeleteOperation{
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.DeleteOperation{
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.DeleteOperation{
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.DeleteOperation{
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.DeleteOperation{
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "fromAge", Type: code.TypeInt},
					{Name: "toAge", Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.DeleteOperation{
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.ArrayType{ContainedType: code.SimpleType("Gender")}},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.DeleteOperation{
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
