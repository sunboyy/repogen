package spec_test

import (
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/spec"
)

var structModel = code.Struct{
	Name: "UserModel",
	Fields: code.StructFields{
		{
			Name: "ID",
			Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
		},
		{
			Name: "PhoneNumber",
			Type: code.SimpleType("string"),
		},
		{
			Name: "Gender",
			Type: code.SimpleType("Gender"),
		},
		{
			Name: "City",
			Type: code.SimpleType("string"),
		},
		{
			Name: "Age",
			Type: code.SimpleType("int"),
		},
	},
}

type ParseInterfaceMethodTestCase struct {
	Name           string
	Method         code.Method
	ExpectedOutput spec.MethodSpec
}

func TestParseInterfaceMethod(t *testing.T) {
	testTable := []ParseInterfaceMethodTestCase{
		{
			Name: "FindOneByArg method",
			Method: code.Method{
				Name: "FindOneByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindOneByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{Predicates: []spec.Predicate{
						{Field: "ID", Comparator: spec.ComparatorEqual},
					}},
				},
			},
		},
		{
			Name: "FindOneByMultiWordArg method",
			Method: code.Method{
				Name: "FindOneByPhoneNumber",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindOneByPhoneNumber",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{Predicates: []spec.Predicate{
						{Field: "PhoneNumber", Comparator: spec.ComparatorEqual},
					}},
				},
			},
		},
		{
			Name: "FindByArg method",
			Method: code.Method{
				Name: "FindByCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindByCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{Predicates: []spec.Predicate{
						{Field: "City", Comparator: spec.ComparatorEqual},
					}},
				},
			},
		},
		{
			Name: "FindAll method",
			Method: code.Method{
				Name: "FindAll",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindAll",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
				},
			},
		},
		{
			Name: "FindByArgAndArg method",
			Method: code.Method{
				Name: "FindByCityAndGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindByCityAndGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Operator: spec.OperatorAnd,
						Predicates: []spec.Predicate{
							{Field: "City", Comparator: spec.ComparatorEqual},
							{Field: "Gender", Comparator: spec.ComparatorEqual},
						},
					},
				},
			},
		},
		{
			Name: "FindByArgOrArg method",
			Method: code.Method{
				Name: "FindByCityOrGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindByCityOrGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Operator: spec.OperatorOr,
						Predicates: []spec.Predicate{
							{Field: "City", Comparator: spec.ComparatorEqual},
							{Field: "Gender", Comparator: spec.ComparatorEqual},
						},
					},
				},
			},
		},
		{
			Name: "FindByArgNot method",
			Method: code.Method{
				Name: "FindByCityNot",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindByCityNot",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{Predicates: []spec.Predicate{
						{Field: "City", Comparator: spec.ComparatorNot},
					}},
				},
			},
		},
		{
			Name: "FindByArgLessThan method",
			Method: code.Method{
				Name: "FindByAgeLessThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindByAgeLessThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{Predicates: []spec.Predicate{
						{Field: "Age", Comparator: spec.ComparatorLessThan},
					}},
				},
			},
		},
		{
			Name: "FindByArgLessThanEqual method",
			Method: code.Method{
				Name: "FindByAgeLessThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindByAgeLessThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{Predicates: []spec.Predicate{
						{Field: "Age", Comparator: spec.ComparatorLessThanEqual},
					}},
				},
			},
		},
		{
			Name: "FindByArgGreaterThan method",
			Method: code.Method{
				Name: "FindByAgeGreaterThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindByAgeGreaterThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{Predicates: []spec.Predicate{
						{Field: "Age", Comparator: spec.ComparatorGreaterThan},
					}},
				},
			},
		},
		{
			Name: "FindByArgGreaterThanEqual method",
			Method: code.Method{
				Name: "FindByAgeGreaterThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindByAgeGreaterThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{Predicates: []spec.Predicate{
						{Field: "Age", Comparator: spec.ComparatorGreaterThanEqual},
					}},
				},
			},
		},
		{
			Name: "FindByArgIn method",
			Method: code.Method{
				Name: "FindByCityIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.SimpleType("string")}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOutput: spec.MethodSpec{
				Name: "FindByCityIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.SimpleType("string")}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{Predicates: []spec.Predicate{
						{Field: "City", Comparator: spec.ComparatorIn},
					}},
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(structModel, testCase.Method)

			if err != nil {
				t.Errorf("Error = %s", err)
			}
			if !reflect.DeepEqual(actualSpec, testCase.ExpectedOutput) {
				t.Errorf("Expected = %v\nReceived = %v", testCase.ExpectedOutput, actualSpec)
			}
		})
	}
}

type ParseInterfaceMethodInvalidTestCase struct {
	Name          string
	Method        code.Method
	ExpectedError error
}

func TestParseInterfaceMethodInvalid(t *testing.T) {
	testTable := []ParseInterfaceMethodInvalidTestCase{
		{
			Name: "unknown operation",
			Method: code.Method{
				Name: "SearchByID",
			},
			ExpectedError: spec.UnknownOperationError,
		},
		{
			Name: "unsupported find method name",
			Method: code.Method{
				Name: "Find",
			},
			ExpectedError: spec.UnsupportedNameError,
		},
		{
			Name: "invalid number of returns",
			Method: code.Method{
				Name: "FindOneByID",
				Returns: []code.Type{
					code.SimpleType("UserModel"),
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "unsupported return values from find method",
			Method: code.Method{
				Name: "FindOneByID",
				Returns: []code.Type{
					code.SimpleType("UserModel"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "error return not provided",
			Method: code.Method{
				Name: "FindOneByID",
				Returns: []code.Type{
					code.SimpleType("UserModel"),
					code.SimpleType("int"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "misplaced operator token (leftmost)",
			Method: code.Method{
				Name: "FindByAndGender",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidQueryError,
		},
		{
			Name: "misplaced operator token (rightmost)",
			Method: code.Method{
				Name: "FindByGenderAnd",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidQueryError,
		},
		{
			Name: "misplaced operator token (double operator)",
			Method: code.Method{
				Name: "FindByGenderAndAndCity",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidQueryError,
		},
		{
			Name: "ambiguous query",
			Method: code.Method{
				Name: "FindByGenderAndCityOrAge",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidQueryError,
		},
		{
			Name: "no context parameter",
			Method: code.Method{
				Name: "FindByGender",
				Params: []code.Param{
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.ContextParamRequiredError,
		},
		{
			Name: "mismatched number of parameters",
			Method: code.Method{
				Name: "FindByCountry",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidParamError,
		},
		{
			Name: "struct field not found",
			Method: code.Method{
				Name: "FindByCountry",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.StructFieldNotFoundError,
		},
		{
			Name: "mismatched method parameter type",
			Method: code.Method{
				Name: "FindByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidParamError,
		},
		{
			Name: "mismatched method parameter type for special case",
			Method: code.Method{
				Name: "FindByCityIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidParamError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(structModel, testCase.Method)

			if err != testCase.ExpectedError {
				t.Errorf("\nExpected = %v\nReceived = %v", testCase.ExpectedError, err)
			}
		})
	}
}
