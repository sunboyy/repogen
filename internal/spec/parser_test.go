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
		{
			Name: "Enabled",
			Type: code.SimpleType("bool"),
		},
	},
}

type ParseInterfaceMethodTestCase struct {
	Name              string
	Method            code.Method
	ExpectedOperation spec.Operation
}

func TestParseInterfaceMethod_Insert(t *testing.T) {
	testTable := []ParseInterfaceMethodTestCase{
		{
			Name: "InsertOne method",
			Method: code.Method{
				Name: "InsertOne",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				},
				Returns: []code.Type{
					code.InterfaceType{},
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.InsertOperation{
				Mode: spec.QueryModeOne,
			},
		},
		{
			Name: "InsertMany method",
			Method: code.Method{
				Name: "InsertMany",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.InterfaceType{}},
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.InsertOperation{
				Mode: spec.QueryModeMany,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(structModel, testCase.Method)

			if err != nil {
				t.Errorf("Error = %s", err)
			}
			expectedOutput := spec.MethodSpec{
				Name:      testCase.Method.Name,
				Params:    testCase.Method.Params,
				Returns:   testCase.Method.Returns,
				Operation: testCase.ExpectedOperation,
			}
			if !reflect.DeepEqual(actualSpec, expectedOutput) {
				t.Errorf("Expected = %v\nReceived = %v", expectedOutput, actualSpec)
			}
		})
	}
}

func TestParseInterfaceMethod_Find(t *testing.T) {
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "ID", Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "PhoneNumber", Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorAnd,
					Predicates: []spec.Predicate{
						{Field: "City", Comparator: spec.ComparatorEqual, ParamIndex: 1},
						{Field: "Gender", Comparator: spec.ComparatorEqual, ParamIndex: 2},
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorOr,
					Predicates: []spec.Predicate{
						{Field: "City", Comparator: spec.ComparatorEqual, ParamIndex: 1},
						{Field: "Gender", Comparator: spec.ComparatorEqual, ParamIndex: 2},
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorNot, ParamIndex: 1},
				}},
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Age", Comparator: spec.ComparatorLessThan, ParamIndex: 1},
				}},
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Age", Comparator: spec.ComparatorLessThanEqual, ParamIndex: 1},
				}},
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Age", Comparator: spec.ComparatorGreaterThan, ParamIndex: 1},
				}},
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Age", Comparator: spec.ComparatorGreaterThanEqual, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindByArgBetween method",
			Method: code.Method{
				Name: "FindByAgeBetween",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Age", Comparator: spec.ComparatorBetween, ParamIndex: 1},
				}},
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
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorIn, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindByArgNotIn method",
			Method: code.Method{
				Name: "FindByCityNotIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.SimpleType("string")}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorNotIn, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindByArgTrue method",
			Method: code.Method{
				Name: "FindByEnabledTrue",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Enabled", Comparator: spec.ComparatorTrue, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindByArgFalse method",
			Method: code.Method{
				Name: "FindByEnabledFalse",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Enabled", Comparator: spec.ComparatorFalse, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindByArgOrderByArg method",
			Method: code.Method{
				Name: "FindByCityOrderByAge",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
				Sorts: []spec.Sort{
					{FieldName: "Age", Ordering: spec.OrderingAscending},
				},
			},
		},
		{
			Name: "FindByArgOrderByArgAsc method",
			Method: code.Method{
				Name: "FindByCityOrderByAgeAsc",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
				Sorts: []spec.Sort{
					{FieldName: "Age", Ordering: spec.OrderingAscending},
				},
			},
		},
		{
			Name: "FindByArgOrderByArgDesc method",
			Method: code.Method{
				Name: "FindByCityOrderByAgeDesc",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
				Sorts: []spec.Sort{
					{FieldName: "Age", Ordering: spec.OrderingDescending},
				},
			},
		},
		{
			Name: "FindByArgOrderByArgAndArg method",
			Method: code.Method{
				Name: "FindByCityOrderByCityAndAgeDesc",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
				Sorts: []spec.Sort{
					{FieldName: "City", Ordering: spec.OrderingAscending},
					{FieldName: "Age", Ordering: spec.OrderingDescending},
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
			expectedOutput := spec.MethodSpec{
				Name:      testCase.Method.Name,
				Params:    testCase.Method.Params,
				Returns:   testCase.Method.Returns,
				Operation: testCase.ExpectedOperation,
			}
			if !reflect.DeepEqual(actualSpec, expectedOutput) {
				t.Errorf("Expected = %v\nReceived = %v", expectedOutput, actualSpec)
			}
		})
	}
}

func TestParseInterfaceMethod_Update(t *testing.T) {
	testTable := []ParseInterfaceMethodTestCase{
		{
			Name: "UpdateByArg",
			Method: code.Method{
				Name: "UpdateByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateModel{},
				Mode:   spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "ID", Comparator: spec.ComparatorEqual, ParamIndex: 2},
				}},
			},
		},
		{
			Name: "UpdateArgByArg one method",
			Method: code.Method{
				Name: "UpdateGenderByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateFields{
					{Name: "Gender", ParamIndex: 1},
				},
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "ID", Comparator: spec.ComparatorEqual, ParamIndex: 2},
				}},
			},
		},
		{
			Name: "UpdateArgByArg many method",
			Method: code.Method{
				Name: "UpdateGenderByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateFields{
					{Name: "Gender", ParamIndex: 1},
				},
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "ID", Comparator: spec.ComparatorEqual, ParamIndex: 2},
				}},
			},
		},
		{
			Name: "UpdateArgAndArgByArg method",
			Method: code.Method{
				Name: "UpdateGenderAndCityByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.SimpleType("string")},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateFields{
					{Name: "Gender", ParamIndex: 1},
					{Name: "City", ParamIndex: 2},
				},
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "ID", Comparator: spec.ComparatorEqual, ParamIndex: 3},
				}},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(structModel, testCase.Method)

			if err != nil {
				t.Errorf("Error = %s", err)
			}
			expectedOutput := spec.MethodSpec{
				Name:      testCase.Method.Name,
				Params:    testCase.Method.Params,
				Returns:   testCase.Method.Returns,
				Operation: testCase.ExpectedOperation,
			}
			if !reflect.DeepEqual(actualSpec, expectedOutput) {
				t.Errorf("Expected = %v\nReceived = %v", expectedOutput, actualSpec)
			}
		})
	}
}

func TestParseInterfaceMethod_Delete(t *testing.T) {
	testTable := []ParseInterfaceMethodTestCase{
		{
			Name: "DeleteOneByArg method",
			Method: code.Method{
				Name: "DeleteOneByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "ID", Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteOneByMultiWordArg method",
			Method: code.Method{
				Name: "DeleteOneByPhoneNumber",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "PhoneNumber", Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteByArg method",
			Method: code.Method{
				Name: "DeleteByCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteAll method",
			Method: code.Method{
				Name: "DeleteAll",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
			},
		},
		{
			Name: "DeleteByArgAndArg method",
			Method: code.Method{
				Name: "DeleteByCityAndGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorAnd,
					Predicates: []spec.Predicate{
						{Field: "City", Comparator: spec.ComparatorEqual, ParamIndex: 1},
						{Field: "Gender", Comparator: spec.ComparatorEqual, ParamIndex: 2},
					},
				},
			},
		},
		{
			Name: "DeleteByArgOrArg method",
			Method: code.Method{
				Name: "DeleteByCityOrGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorOr,
					Predicates: []spec.Predicate{
						{Field: "City", Comparator: spec.ComparatorEqual, ParamIndex: 1},
						{Field: "Gender", Comparator: spec.ComparatorEqual, ParamIndex: 2},
					},
				},
			},
		},
		{
			Name: "DeleteByArgNot method",
			Method: code.Method{
				Name: "DeleteByCityNot",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorNot, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteByArgLessThan method",
			Method: code.Method{
				Name: "DeleteByAgeLessThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Age", Comparator: spec.ComparatorLessThan, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteByArgLessThanEqual method",
			Method: code.Method{
				Name: "DeleteByAgeLessThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Age", Comparator: spec.ComparatorLessThanEqual, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteByArgGreaterThan method",
			Method: code.Method{
				Name: "DeleteByAgeGreaterThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Age", Comparator: spec.ComparatorGreaterThan, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteByArgGreaterThanEqual method",
			Method: code.Method{
				Name: "DeleteByAgeGreaterThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Age", Comparator: spec.ComparatorGreaterThanEqual, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteByArgBetween method",
			Method: code.Method{
				Name: "DeleteByAgeBetween",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "Age", Comparator: spec.ComparatorBetween, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteByArgIn method",
			Method: code.Method{
				Name: "DeleteByCityIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.SimpleType("string")}},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{Field: "City", Comparator: spec.ComparatorIn, ParamIndex: 1},
				}},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(structModel, testCase.Method)

			if err != nil {
				t.Errorf("Error = %s", err)
			}
			expectedOutput := spec.MethodSpec{
				Name:      testCase.Method.Name,
				Params:    testCase.Method.Params,
				Returns:   testCase.Method.Returns,
				Operation: testCase.ExpectedOperation,
			}
			if !reflect.DeepEqual(actualSpec, expectedOutput) {
				t.Errorf("Expected = %v\nReceived = %v", expectedOutput, actualSpec)
			}
		})
	}
}

func TestParseInterfaceMethod_Count(t *testing.T) {
	testTable := []ParseInterfaceMethodTestCase{
		{
			Name: "CountAll method",
			Method: code.Method{
				Name: "CountAll",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.CountOperation{
				Query: spec.QuerySpec{},
			},
		},
		{
			Name: "CountByArg method",
			Method: code.Method{
				Name: "CountByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedOperation: spec.CountOperation{
				Query: spec.QuerySpec{
					Predicates: []spec.Predicate{
						{Field: "Gender", Comparator: spec.ComparatorEqual, ParamIndex: 1},
					},
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
			expectedOutput := spec.MethodSpec{
				Name:      testCase.Method.Name,
				Params:    testCase.Method.Params,
				Returns:   testCase.Method.Returns,
				Operation: testCase.ExpectedOperation,
			}
			if !reflect.DeepEqual(actualSpec, expectedOutput) {
				t.Errorf("Expected = %v\nReceived = %v", expectedOutput, actualSpec)
			}
		})
	}
}

type ParseInterfaceMethodInvalidTestCase struct {
	Name          string
	Method        code.Method
	ExpectedError error
}

func TestParseInterfaceMethod_Invalid(t *testing.T) {
	_, err := spec.ParseInterfaceMethod(structModel, code.Method{
		Name: "SearchByID",
	})

	expectedError := spec.NewUnknownOperationError("Search")
	if err != expectedError {
		t.Errorf("\nExpected = %v\nReceived = %v", expectedError, err)
	}
}

func TestParseInterfaceMethod_Insert_Invalid(t *testing.T) {
	testTable := []ParseInterfaceMethodInvalidTestCase{
		{
			Name: "invalid number of returns",
			Method: code.Method{
				Name: "Insert",
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.InterfaceType{},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "unsupported return types from insert method",
			Method: code.Method{
				Name: "Insert",
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "unempty interface return from insert method",
			Method: code.Method{
				Name: "Insert",
				Returns: []code.Type{
					code.InterfaceType{
						Methods: []code.Method{
							{Name: "DoSomething"},
						},
					},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "error return not provided",
			Method: code.Method{
				Name: "Insert",
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.InterfaceType{},
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "no context parameter",
			Method: code.Method{
				Name: "Insert",
				Params: []code.Param{
					{Name: "userModel", Type: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				},
				Returns: []code.Type{
					code.InterfaceType{},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.ContextParamRequiredError,
		},
		{
			Name: "mismatched model parameter for one mode",
			Method: code.Method{
				Name: "Insert",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "userModel", Type: code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}}},
				},
				Returns: []code.Type{
					code.InterfaceType{},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidParamError,
		},
		{
			Name: "mismatched model parameter for many mode",
			Method: code.Method{
				Name: "Insert",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "userModel", Type: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.InterfaceType{}},
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

func TestParseInterfaceMethod_Find_Invalid(t *testing.T) {
	testTable := []ParseInterfaceMethodInvalidTestCase{
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
			Name: "unsupported return types from find method",
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
			Name: "find method without query",
			Method: code.Method{
				Name: "Find",
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.QueryRequiredError,
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
			ExpectedError: spec.NewInvalidQueryError([]string{"By", "And", "Gender"}),
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
			ExpectedError: spec.NewInvalidQueryError([]string{"By", "Gender", "And"}),
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
			ExpectedError: spec.NewInvalidQueryError([]string{"By", "Gender", "And", "And", "City"}),
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
			ExpectedError: spec.NewInvalidQueryError([]string{"By", "Gender", "And", "City", "Or", "Age"}),
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
			ExpectedError: spec.NewStructFieldNotFoundError("Country"),
		},
		{
			Name: "incompatible struct field for True comparator",
			Method: code.Method{
				Name: "FindByGenderTrue",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewIncompatibleComparatorError(spec.ComparatorTrue, code.StructField{
				Name: "Gender",
				Type: code.SimpleType("Gender"),
			}),
		},
		{
			Name: "incompatible struct field for False comparator",
			Method: code.Method{
				Name: "FindByGenderFalse",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewIncompatibleComparatorError(spec.ComparatorFalse, code.StructField{
				Name: "Gender",
				Type: code.SimpleType("Gender"),
			}),
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
		{
			Name: "misplaced operator token (leftmost)",
			Method: code.Method{
				Name: "FindAllOrderByAndAge",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewInvalidSortError([]string{"Order", "By", "And", "Age"}),
		},
		{
			Name: "misplaced operator token (rightmost)",
			Method: code.Method{
				Name: "FindAllOrderByAgeAnd",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewInvalidSortError([]string{"Order", "By", "Age", "And"}),
		},
		{
			Name: "misplaced operator token (double operator)",
			Method: code.Method{
				Name: "FindAllOrderByAgeAndAndGender",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewInvalidSortError([]string{"Order", "By", "Age", "And", "And", "Gender"}),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(structModel, testCase.Method)

			if err.Error() != testCase.ExpectedError.Error() {
				t.Errorf("\nExpected = %v\nReceived = %v", testCase.ExpectedError.Error(), err.Error())
			}
		})
	}
}

func TestParseInterfaceMethod_Update_Invalid(t *testing.T) {
	testTable := []ParseInterfaceMethodInvalidTestCase{
		{
			Name: "invalid number of returns",
			Method: code.Method{
				Name: "UpdateAgeByID",
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "unsupported return types from update method",
			Method: code.Method{
				Name: "UpdateAgeByID",
				Returns: []code.Type{
					code.SimpleType("float64"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "error return not provided",
			Method: code.Method{
				Name: "UpdateAgeByID",
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("bool"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "update with no field provided",
			Method: code.Method{
				Name: "UpdateByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidUpdateFieldsError,
		},
		{
			Name: "misplaced And token in update fields",
			Method: code.Method{
				Name: "UpdateAgeAndAndGenderByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidUpdateFieldsError,
		},
		{
			Name: "update method without query",
			Method: code.Method{
				Name: "UpdateCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.QueryRequiredError,
		},
		{
			Name: "ambiguous query",
			Method: code.Method{
				Name: "UpdateAgeByIDAndUsernameOrGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"By", "ID", "And", "Username", "Or", "Gender"}),
		},
		{
			Name: "update model with invalid parameter",
			Method: code.Method{
				Name: "UpdateByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidUpdateFieldsError,
		},
		{
			Name: "no context parameter",
			Method: code.Method{
				Name: "UpdateAgeByGender",
				Params: []code.Param{
					{Type: code.SimpleType("int")},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.ContextParamRequiredError,
		},
		{
			Name: "struct field not found in update fields",
			Method: code.Method{
				Name: "UpdateCountryByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError("Country"),
		},
		{
			Name: "struct field does not match parameter in update fields",
			Method: code.Method{
				Name: "UpdateAgeByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("float64")},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidUpdateFieldsError,
		},
		{
			Name: "struct field does not match parameter in query",
			Method: code.Method{
				Name: "UpdateAgeByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
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

func TestParseInterfaceMethod_Delete_Invalid(t *testing.T) {
	testTable := []ParseInterfaceMethodInvalidTestCase{
		{
			Name: "invalid number of returns",
			Method: code.Method{
				Name: "DeleteOneByID",
				Returns: []code.Type{
					code.SimpleType("UserModel"),
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "unsupported return types from delete method",
			Method: code.Method{
				Name: "DeleteOneByID",
				Returns: []code.Type{
					code.SimpleType("float64"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "error return not provided",
			Method: code.Method{
				Name: "DeleteOneByID",
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("bool"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "delete method without query",
			Method: code.Method{
				Name: "Delete",
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.QueryRequiredError,
		},
		{
			Name: "misplaced operator token (leftmost)",
			Method: code.Method{
				Name: "DeleteByAndGender",
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"By", "And", "Gender"}),
		},
		{
			Name: "misplaced operator token (rightmost)",
			Method: code.Method{
				Name: "DeleteByGenderAnd",
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"By", "Gender", "And"}),
		},
		{
			Name: "misplaced operator token (double operator)",
			Method: code.Method{
				Name: "DeleteByGenderAndAndCity",
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"By", "Gender", "And", "And", "City"}),
		},
		{
			Name: "ambiguous query",
			Method: code.Method{
				Name: "DeleteByGenderAndCityOrAge",
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"By", "Gender", "And", "City", "Or", "Age"}),
		},
		{
			Name: "no context parameter",
			Method: code.Method{
				Name: "DeleteByGender",
				Params: []code.Param{
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.ContextParamRequiredError,
		},
		{
			Name: "mismatched number of parameters",
			Method: code.Method{
				Name: "DeleteByCountry",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidParamError,
		},
		{
			Name: "struct field not found",
			Method: code.Method{
				Name: "DeleteByCountry",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError("Country"),
		},
		{
			Name: "mismatched method parameter type",
			Method: code.Method{
				Name: "DeleteByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidParamError,
		},
		{
			Name: "mismatched method parameter type for special case",
			Method: code.Method{
				Name: "DeleteByCityIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
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

func TestParseInterfaceMethod_Count_Invalid(t *testing.T) {
	testTable := []ParseInterfaceMethodInvalidTestCase{
		{
			Name: "invalid number of returns",
			Method: code.Method{
				Name: "CountAll",
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
					code.SimpleType("bool"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "invalid integer return",
			Method: code.Method{
				Name: "CountAll",
				Returns: []code.Type{
					code.SimpleType("int64"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "error return not provided",
			Method: code.Method{
				Name: "CountAll",
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("bool"),
				},
			},
			ExpectedError: spec.UnsupportedReturnError,
		},
		{
			Name: "count method without query",
			Method: code.Method{
				Name: "Count",
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.QueryRequiredError,
		},
		{
			Name: "invalid query",
			Method: code.Method{
				Name: "CountBy",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"By"}),
		},
		{
			Name: "context parameter not provided",
			Method: code.Method{
				Name: "CountAll",
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.ContextParamRequiredError,
		},
		{
			Name: "mismatched number of parameter",
			Method: code.Method{
				Name: "CountByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidParamError,
		},
		{
			Name: "mismatched method parameter type",
			Method: code.Method{
				Name: "CountByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.InvalidParamError,
		},
		{
			Name: "struct field not found",
			Method: code.Method{
				Name: "CountByCountry",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError("Country"),
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
