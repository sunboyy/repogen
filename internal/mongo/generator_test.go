package mongo_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
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
	referrerField = code.StructField{
		Name: "Referrer",
		Type: code.PointerType{ContainedType: code.SimpleType("UserModel")},
		Tags: map[string][]string{"bson": {"referrer"}},
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
		referrerField,
		consentHistoryField,
		enabledField,
		accessTokenField,
	},
}

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
		Body: codegen.FunctionBody{
			codegen.ReturnStatement{
				codegen.StructStatement{
					Type: "&UserRepositoryMongo",
					Pairs: []codegen.StructFieldPair{{
						Key:   "collection",
						Value: codegen.Identifier("collection"),
					}},
				},
			},
		},
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
	if !reflect.DeepEqual(expected.Body, actual.Body) {
		t.Errorf("incorrect function body: expected %+v got %+v",
			expected.Body,
			actual.Body,
		)
	}
}

type GenerateMethodTestCase struct {
	Name         string
	MethodSpec   spec.MethodSpec
	ExpectedBody string
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
