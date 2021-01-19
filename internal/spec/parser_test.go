package spec_test

import (
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/spec"
)

type TestCase struct {
	Name           string
	Interface      code.Interface
	ExpectedOutput spec.RepositorySpec
}

func TestParseRepositoryInterface(t *testing.T) {
	testTable := []TestCase{
		{
			Name: "interface spec",
			Interface: code.Interface{
				Name: "UserRepository",
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
			},
		},
		{
			Name: "FindOneByArg method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
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
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
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
			},
		},
		{
			Name: "FindOneByMultiWordArg method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
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
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
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
			},
		},
		{
			Name: "FindByArg method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
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
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
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
			},
		},
		{
			Name: "FindAll method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
						Name: "FindAll",
						Params: []code.Param{
							{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
						},
						Returns: []code.Type{
							code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
							code.SimpleType("error"),
						},
					},
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
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
			},
		},
		{
			Name: "FindByArgAndArg method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
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
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
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
			},
		},
		{
			Name: "FindByArgOrArg method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
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
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
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
			},
		},
		{
			Name: "FindByArgNot method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
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
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
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
			},
		},
		{
			Name: "FindByArgLessThan method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
						Name: "FindByAgeLessThan",
						Params: []code.Param{
							{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
							{Type: code.SimpleType("string")},
						},
						Returns: []code.Type{
							code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
							code.SimpleType("error"),
						},
					},
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
						Name: "FindByAgeLessThan",
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
								{Field: "Age", Comparator: spec.ComparatorLessThan},
							}},
						},
					},
				},
			},
		},
		{
			Name: "FindByArgLessThanEqual method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
						Name: "FindByAgeLessThanEqual",
						Params: []code.Param{
							{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
							{Type: code.SimpleType("string")},
						},
						Returns: []code.Type{
							code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
							code.SimpleType("error"),
						},
					},
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
						Name: "FindByAgeLessThanEqual",
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
								{Field: "Age", Comparator: spec.ComparatorLessThanEqual},
							}},
						},
					},
				},
			},
		},
		{
			Name: "FindByArgGreaterThan method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
						Name: "FindByAgeGreaterThan",
						Params: []code.Param{
							{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
							{Type: code.SimpleType("string")},
						},
						Returns: []code.Type{
							code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
							code.SimpleType("error"),
						},
					},
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
						Name: "FindByAgeGreaterThan",
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
								{Field: "Age", Comparator: spec.ComparatorGreaterThan},
							}},
						},
					},
				},
			},
		},
		{
			Name: "FindByArgGreaterThanEqual method",
			Interface: code.Interface{
				Name: "UserRepository",
				Methods: []code.Method{
					{
						Name: "FindByAgeGreaterThanEqual",
						Params: []code.Param{
							{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
							{Type: code.SimpleType("string")},
						},
						Returns: []code.Type{
							code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
							code.SimpleType("error"),
						},
					},
				},
			},
			ExpectedOutput: spec.RepositorySpec{
				InterfaceName: "UserRepository",
				Methods: []spec.MethodSpec{
					{
						Name: "FindByAgeGreaterThanEqual",
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
								{Field: "Age", Comparator: spec.ComparatorGreaterThanEqual},
							}},
						},
					},
				},
			},
		},
	}

	structModel := code.Struct{
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

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSpec, err := spec.ParseRepositoryInterface(structModel, testCase.Interface)

			if err != nil {
				t.Errorf("Error = %s", err)
			}
			if !reflect.DeepEqual(actualSpec, testCase.ExpectedOutput) {
				t.Errorf("Expected = %v\nReceived = %v", testCase.ExpectedOutput, actualSpec)
			}
		})
	}
}
