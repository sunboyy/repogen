package mongo_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
	"github.com/sunboyy/repogen/internal/testutils"
)

var (
	idField = code.StructField{
		Name: "ID",
		Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
		Tags: map[string][]string{"bson": {"_id", "omitempty"}},
	}
	genderField = code.StructField{
		Name: "Gender",
		Type: code.SimpleType("Gender"),
		Tags: map[string][]string{"bson": {"gender"}},
	}
	ageField = code.StructField{
		Name: "Age",
		Type: code.TypeInt,
		Tags: map[string][]string{"bson": {"age"}},
	}
	nameField = code.StructField{
		Name: "Name",
		Type: code.SimpleType("NameModel"),
		Tags: map[string][]string{"bson": {"name"}},
	}
	consentHistoryField = code.StructField{
		Name: "ConsentHistory",
		Type: code.ArrayType{ContainedType: code.SimpleType("ConsentHistory")},
		Tags: map[string][]string{"bson": {"consent_history"}},
	}
	enabledField = code.StructField{
		Name: "Enabled",
		Type: code.TypeBool,
		Tags: map[string][]string{"bson": {"enabled"}},
	}
	accessTokenField = code.StructField{
		Name: "AccessToken",
		Type: code.TypeString,
	}

	firstNameField = code.StructField{
		Name: "First",
		Type: code.TypeString,
		Tags: map[string][]string{"bson": {"first"}},
	}
)

var userModel = code.Struct{
	Name: "UserModel",
	Fields: code.StructFields{
		idField,
		code.StructField{
			Name: "Username",
			Type: code.TypeString,
			Tags: map[string][]string{"bson": {"username"}},
		},
		genderField,
		ageField,
		nameField,
		consentHistoryField,
		enabledField,
		accessTokenField,
	},
}

const expectedConstructorBody = `	return &UserRepositoryMongo{
		collection: collection,
	}`

func TestImports(t *testing.T) {
	generator := mongo.NewGenerator(userModel, "UserRepository")
	expected := [][]code.Import{
		{
			{Path: "context"},
		},
		{
			{Path: "go.mongodb.org/mongo-driver/bson"},
			{Path: "go.mongodb.org/mongo-driver/bson/primitive"},
			{Path: "go.mongodb.org/mongo-driver/mongo"},
			{Path: "go.mongodb.org/mongo-driver/mongo/options"},
		},
	}

	actual := generator.Imports()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("incorrect imports: expected %+v, got %+v", expected, actual)
	}
}

func TestGenerateStruct(t *testing.T) {
	generator := mongo.NewGenerator(userModel, "UserRepository")
	expected := codegen.StructBuilder{
		Name: "UserRepositoryMongo",
		Fields: []code.StructField{
			{
				Name: "collection",
				Type: code.PointerType{
					ContainedType: code.ExternalType{
						PackageAlias: "mongo",
						Name:         "Collection",
					},
				},
			},
		},
	}

	actual := generator.GenerateStruct()

	if expected.Name != actual.Name {
		t.Errorf(
			"incorrect struct name: expected %s, got %s",
			expected.Name,
			actual.Name,
		)
	}
	if !reflect.DeepEqual(expected.Fields, actual.Fields) {
		t.Errorf(
			"incorrect struct fields: expected %+v, got %+v",
			expected.Fields,
			actual.Fields,
		)
	}
}

func TestGenerateConstructor(t *testing.T) {
	generator := mongo.NewGenerator(userModel, "UserRepository")
	expected := codegen.FunctionBuilder{
		Name: "NewUserRepository",
		Params: []code.Param{
			{
				Name: "collection",
				Type: code.PointerType{
					ContainedType: code.ExternalType{
						PackageAlias: "mongo",
						Name:         "Collection",
					},
				},
			},
		},
		Returns: []code.Type{
			code.SimpleType("UserRepository"),
		},
		Body: expectedConstructorBody,
	}

	actual, err := generator.GenerateConstructor()

	if err != nil {
		t.Fatal(err)
	}
	if expected.Name != actual.Name {
		t.Errorf(
			"incorrect function name: expected %s, got %s",
			expected.Name,
			actual.Name,
		)
	}
	if !reflect.DeepEqual(expected.Params, actual.Params) {
		t.Errorf(
			"incorrect struct params: expected %+v, got %+v",
			expected.Params,
			actual.Params,
		)
	}
	if err := testutils.ExpectMultiLineString(expected.Body, actual.Body); err != nil {
		t.Error(err)
	}
}

type GenerateMethodTestCase struct {
	Name         string
	MethodSpec   spec.MethodSpec
	ExpectedBody string
}

func TestGenerateMethod_Insert(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "insert one method",
			MethodSpec: spec.MethodSpec{
				Name: "InsertOne",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "userModel", Type: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				},
				Returns: []code.Type{
					code.InterfaceType{},
					code.TypeError,
				},
				Operation: spec.InsertOperation{
					Mode: spec.QueryModeOne,
				},
			},
			ExpectedBody: `	result, err := r.collection.InsertOne(arg0, arg1)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil`,
		},
		{
			Name: "insert many method",
			MethodSpec: spec.MethodSpec{
				Name: "Insert",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "userModel", Type: code.ArrayType{
						ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")},
					}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.InterfaceType{}},
					code.TypeError,
				},
				Operation: spec.InsertOperation{
					Mode: spec.QueryModeMany,
				},
			},
			ExpectedBody: `	var entities []interface{}
	for _, model := range arg1 {
		entities = append(entities, model)
	}
	result, err := r.collection.InsertMany(arg0, entities)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			var expectedParams []code.Param
			for i, param := range testCase.MethodSpec.Params {
				expectedParams = append(expectedParams, code.Param{
					Name: fmt.Sprintf("arg%d", i),
					Type: param.Type,
				})
			}
			expected := codegen.MethodBuilder{
				Receiver: codegen.MethodReceiver{
					Name:    "r",
					Type:    "UserRepositoryMongo",
					Pointer: true,
				},
				Name:    testCase.MethodSpec.Name,
				Params:  expectedParams,
				Returns: testCase.MethodSpec.Returns,
				Body:    testCase.ExpectedBody,
			}

			actual, err := generator.GenerateMethod(testCase.MethodSpec)

			if err != nil {
				t.Fatal(err)
			}
			if expected.Receiver != actual.Receiver {
				t.Errorf(
					"incorrect method receiver: expected %+v, got %+v",
					expected.Receiver,
					actual.Receiver,
				)
			}
			if expected.Name != actual.Name {
				t.Errorf(
					"incorrect method name: expected %s, got %s",
					expected.Name,
					actual.Name,
				)
			}
			if !reflect.DeepEqual(expected.Params, actual.Params) {
				t.Errorf(
					"incorrect struct params: expected %+v, got %+v",
					expected.Params,
					actual.Params,
				)
			}
			if !reflect.DeepEqual(expected.Returns, actual.Returns) {
				t.Errorf(
					"incorrect struct returns: expected %+v, got %+v",
					expected.Returns,
					actual.Returns,
				)
			}
			if err := testutils.ExpectMultiLineString(expected.Body, actual.Body); err != nil {
				t.Error(err)
			}
		})
	}
}

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
			{"gender": arg1},
			{"age": arg2},
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
			{"gender": arg1},
			{"age": arg2},
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
		"gender": bson.M{"$ne": arg1},
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
		"age": bson.M{"$lt": arg1},
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
		"age": bson.M{"$lte": arg1},
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
		"age": bson.M{"$gt": arg1},
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
		"age": bson.M{"$gte": arg1},
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
		"age": bson.M{"$gte": arg1, "$lte": arg2},
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
		"gender": bson.M{"$in": arg1},
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
		"gender": bson.M{"$nin": arg1},
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
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			var expectedParams []code.Param
			for i, param := range testCase.MethodSpec.Params {
				expectedParams = append(expectedParams, code.Param{
					Name: fmt.Sprintf("arg%d", i),
					Type: param.Type,
				})
			}
			expected := codegen.MethodBuilder{
				Receiver: codegen.MethodReceiver{
					Name:    "r",
					Type:    "UserRepositoryMongo",
					Pointer: true,
				},
				Name:    testCase.MethodSpec.Name,
				Params:  expectedParams,
				Returns: testCase.MethodSpec.Returns,
				Body:    testCase.ExpectedBody,
			}

			actual, err := generator.GenerateMethod(testCase.MethodSpec)

			if err != nil {
				t.Fatal(err)
			}
			if expected.Receiver != actual.Receiver {
				t.Errorf(
					"incorrect method receiver: expected %+v, got %+v",
					expected.Receiver,
					actual.Receiver,
				)
			}
			if expected.Name != actual.Name {
				t.Errorf(
					"incorrect method name: expected %s, got %s",
					expected.Name,
					actual.Name,
				)
			}
			if !reflect.DeepEqual(expected.Params, actual.Params) {
				t.Errorf(
					"incorrect struct params: expected %+v, got %+v",
					expected.Params,
					actual.Params,
				)
			}
			if !reflect.DeepEqual(expected.Returns, actual.Returns) {
				t.Errorf(
					"incorrect struct returns: expected %+v, got %+v",
					expected.Returns,
					actual.Returns,
				)
			}
			if err := testutils.ExpectMultiLineString(expected.Body, actual.Body); err != nil {
				t.Error(err)
			}
		})
	}
}

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
	return result.MatchedCount > 0, err`,
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
	return result.MatchedCount > 0, err`,
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
	return int(result.MatchedCount), err`,
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
	return result.MatchedCount > 0, err`,
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
	return result.MatchedCount > 0, err`,
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
	return result.MatchedCount > 0, err`,
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
	return result.MatchedCount > 0, err`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			var expectedParams []code.Param
			for i, param := range testCase.MethodSpec.Params {
				expectedParams = append(expectedParams, code.Param{
					Name: fmt.Sprintf("arg%d", i),
					Type: param.Type,
				})
			}
			expected := codegen.MethodBuilder{
				Receiver: codegen.MethodReceiver{
					Name:    "r",
					Type:    "UserRepositoryMongo",
					Pointer: true,
				},
				Name:    testCase.MethodSpec.Name,
				Params:  expectedParams,
				Returns: testCase.MethodSpec.Returns,
				Body:    testCase.ExpectedBody,
			}

			actual, err := generator.GenerateMethod(testCase.MethodSpec)

			if err != nil {
				t.Fatal(err)
			}
			if expected.Receiver != actual.Receiver {
				t.Errorf(
					"incorrect method receiver: expected %+v, got %+v",
					expected.Receiver,
					actual.Receiver,
				)
			}
			if expected.Name != actual.Name {
				t.Errorf(
					"incorrect method name: expected %s, got %s",
					expected.Name,
					actual.Name,
				)
			}
			if !reflect.DeepEqual(expected.Params, actual.Params) {
				t.Errorf(
					"incorrect struct params: expected %+v, got %+v",
					expected.Params,
					actual.Params,
				)
			}
			if !reflect.DeepEqual(expected.Returns, actual.Returns) {
				t.Errorf(
					"incorrect struct returns: expected %+v, got %+v",
					expected.Returns,
					actual.Returns,
				)
			}
			if err := testutils.ExpectMultiLineString(expected.Body, actual.Body); err != nil {
				t.Error(err)
			}
		})
	}
}

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
			{"gender": arg1},
			{"age": arg2},
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
			{"gender": arg1},
			{"age": arg2},
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
		"gender": bson.M{"$ne": arg1},
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
		"age": bson.M{"$lt": arg1},
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
		"age": bson.M{"$lte": arg1},
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
		"age": bson.M{"$gt": arg1},
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
		"age": bson.M{"$gte": arg1},
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
		"age": bson.M{"$gte": arg1, "$lte": arg2},
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
		"gender": bson.M{"$in": arg1},
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
			var expectedParams []code.Param
			for i, param := range testCase.MethodSpec.Params {
				expectedParams = append(expectedParams, code.Param{
					Name: fmt.Sprintf("arg%d", i),
					Type: param.Type,
				})
			}
			expected := codegen.MethodBuilder{
				Receiver: codegen.MethodReceiver{
					Name:    "r",
					Type:    "UserRepositoryMongo",
					Pointer: true,
				},
				Name:    testCase.MethodSpec.Name,
				Params:  expectedParams,
				Returns: testCase.MethodSpec.Returns,
				Body:    testCase.ExpectedBody,
			}

			actual, err := generator.GenerateMethod(testCase.MethodSpec)

			if err != nil {
				t.Fatal(err)
			}
			if expected.Receiver != actual.Receiver {
				t.Errorf(
					"incorrect method receiver: expected %+v, got %+v",
					expected.Receiver,
					actual.Receiver,
				)
			}
			if expected.Name != actual.Name {
				t.Errorf(
					"incorrect method name: expected %s, got %s",
					expected.Name,
					actual.Name,
				)
			}
			if !reflect.DeepEqual(expected.Params, actual.Params) {
				t.Errorf(
					"incorrect struct params: expected %+v, got %+v",
					expected.Params,
					actual.Params,
				)
			}
			if !reflect.DeepEqual(expected.Returns, actual.Returns) {
				t.Errorf(
					"incorrect struct returns: expected %+v, got %+v",
					expected.Returns,
					actual.Returns,
				)
			}
			if err := testutils.ExpectMultiLineString(expected.Body, actual.Body); err != nil {
				t.Error(err)
			}
		})
	}
}

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
			{"gender": arg1},
			{"age": arg2},
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
			{"gender": arg1},
			{"age": arg2},
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
		"gender": bson.M{"$ne": arg1},
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
		"age": bson.M{"$lt": arg1},
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
		"age": bson.M{"$lte": arg1},
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
		"age": bson.M{"$gt": arg1},
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
		"age": bson.M{"$gte": arg1},
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
		"age": bson.M{"$gte": arg1, "$lte": arg2},
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
		"age": bson.M{"$in": arg1},
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
			var expectedParams []code.Param
			for i, param := range testCase.MethodSpec.Params {
				expectedParams = append(expectedParams, code.Param{
					Name: fmt.Sprintf("arg%d", i),
					Type: param.Type,
				})
			}
			expected := codegen.MethodBuilder{
				Receiver: codegen.MethodReceiver{
					Name:    "r",
					Type:    "UserRepositoryMongo",
					Pointer: true,
				},
				Name:    testCase.MethodSpec.Name,
				Params:  expectedParams,
				Returns: testCase.MethodSpec.Returns,
				Body:    testCase.ExpectedBody,
			}

			actual, err := generator.GenerateMethod(testCase.MethodSpec)

			if err != nil {
				t.Fatal(err)
			}
			if expected.Receiver != actual.Receiver {
				t.Errorf(
					"incorrect method receiver: expected %+v, got %+v",
					expected.Receiver,
					actual.Receiver,
				)
			}
			if expected.Name != actual.Name {
				t.Errorf(
					"incorrect method name: expected %s, got %s",
					expected.Name,
					actual.Name,
				)
			}
			if !reflect.DeepEqual(expected.Params, actual.Params) {
				t.Errorf(
					"incorrect struct params: expected %+v, got %+v",
					expected.Params,
					actual.Params,
				)
			}
			if !reflect.DeepEqual(expected.Returns, actual.Returns) {
				t.Errorf(
					"incorrect struct returns: expected %+v, got %+v",
					expected.Returns,
					actual.Returns,
				)
			}
			if err := testutils.ExpectMultiLineString(expected.Body, actual.Body); err != nil {
				t.Error(err)
			}
		})
	}
}

type GenerateMethodInvalidTestCase struct {
	Name          string
	Method        spec.MethodSpec
	ExpectedError error
}

type StubOperation struct {
}

func (o StubOperation) Name() string {
	return "Stub"
}

func TestGenerateMethod_Invalid(t *testing.T) {
	testTable := []GenerateMethodInvalidTestCase{
		{
			Name: "operation not supported",
			Method: spec.MethodSpec{
				Name: "SearchByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: StubOperation{},
			},
			ExpectedError: mongo.NewOperationNotSupportedError("Stub"),
		},
		{
			Name: "bson tag not found in query",
			Method: spec.MethodSpec{
				Name: "FindByAccessToken",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{accessTokenField},
								Comparator:     spec.ComparatorEqual,
								ParamIndex:     1,
							},
						},
					},
				},
			},
			ExpectedError: mongo.NewBsonTagNotFoundError("AccessToken"),
		},
		{
			Name: "bson tag not found in sort",
			Method: spec.MethodSpec{
				Name: "FindAllOrderByAccessToken",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Sorts: []spec.Sort{
						{FieldReference: spec.FieldReference{accessTokenField}, Ordering: spec.OrderingAscending},
					},
				},
			},
			ExpectedError: mongo.NewBsonTagNotFoundError("AccessToken"),
		},
		{
			Name: "bson tag not found in update field",
			Method: spec.MethodSpec{
				Name: "UpdateAccessTokenByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{accessTokenField},
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
			ExpectedError: mongo.NewBsonTagNotFoundError("AccessToken"),
		},
		{
			Name: "update type not supported",
			Method: spec.MethodSpec{
				Name: "UpdateAgeByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
				Operation: spec.UpdateOperation{
					Update: StubUpdate{},
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
			ExpectedError: mongo.NewUpdateTypeNotSupportedError(StubUpdate{}),
		},
		{
			Name: "update operator not supported",
			Method: spec.MethodSpec{
				Name: "UpdateConsentHistoryAppendByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
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
							Operator:       "APPEND",
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
			ExpectedError: mongo.NewUpdateOperatorNotSupportedError("APPEND"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")

			_, err := generator.GenerateMethod(testCase.Method)

			if !errors.Is(err, testCase.ExpectedError) {
				t.Errorf("\nExpected = %+v\nReceived = %+v", testCase.ExpectedError, err)
			}
		})
	}
}
