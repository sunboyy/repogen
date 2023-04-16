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

func TestGenerateMethod_Update(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "update model method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateByID",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "model", Type: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateModel{},
					Mode:   spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{idField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     2,
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{ageField},
							ParamIndex:     1,
							Operator:       spec.UpdateOperatorSet,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{idField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     2,
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{ageField},
							ParamIndex:     1,
							Operator:       spec.UpdateOperatorSet,
						},
					},
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{genderField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     2,
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "consentHistory", Type: code.SimpleType("ConsentHistory")},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{consentHistoryField},
							ParamIndex:     1,
							Operator:       spec.UpdateOperatorPush,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{idField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     2,
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.TypeInt},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{ageField},
							ParamIndex:     1,
							Operator:       spec.UpdateOperatorInc,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{idField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     2,
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "enabled", Type: code.TypeBool},
					{Name: "consentHistory", Type: code.SimpleType("ConsentHistory")},
					{Name: "gender", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{enabledField},
							ParamIndex:     1,
							Operator:       spec.UpdateOperatorSet,
						},
						spec.UpdateField{
							FieldReference: spec.FieldReference{consentHistoryField},
							ParamIndex:     2,
							Operator:       spec.UpdateOperatorPush,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{idField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     3,
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "firstName", Type: code.TypeString},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{nameField, firstNameField},
							ParamIndex:     1,
							Operator:       spec.UpdateOperatorSet,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{idField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     2,
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
