package spec_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/spec"
)

var (
	idField = code.StructField{
		Name: "ID",
		Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
	}
	phoneNumberField = code.StructField{
		Name: "PhoneNumber",
		Type: code.TypeString,
	}
	genderField = code.StructField{
		Name: "Gender",
		Type: code.SimpleType("Gender"),
	}
	cityField = code.StructField{
		Name: "City",
		Type: code.TypeString,
	}
	ageField = code.StructField{
		Name: "Age",
		Type: code.TypeInt,
	}
	nameField = code.StructField{
		Name: "Name",
		Type: code.SimpleType("NameModel"),
	}
	contactField = code.StructField{
		Name: "Contact",
		Type: code.SimpleType("ContactModel"),
	}
	referrerField = code.StructField{
		Name: "Referrer",
		Type: code.PointerType{ContainedType: code.SimpleType("UserModel")},
	}
	defaultPaymentField = code.StructField{
		Name: "DefaultPayment",
		Type: code.ExternalType{PackageAlias: "payment", Name: "Payment"},
	}
	enabledField = code.StructField{
		Name: "Enabled",
		Type: code.TypeBool,
	}
	consentHistoryField = code.StructField{
		Name: "ConsentHistory",
		Type: code.ArrayType{ContainedType: code.SimpleType("ConsentHistoryItem")},
	}

	firstNameField = code.StructField{
		Name: "First",
		Type: code.TypeString,
	}
	lastNameField = code.StructField{
		Name: "Last",
		Type: code.TypeString,
	}
)

var (
	nameStruct = code.Struct{
		Name: "NameModel",
		Fields: code.StructFields{
			firstNameField,
			lastNameField,
		},
	}

	structModel = code.Struct{
		Name: "UserModel",
		Fields: code.StructFields{
			idField,
			phoneNumberField,
			genderField,
			cityField,
			ageField,
			nameField,
			contactField,
			referrerField,
			defaultPaymentField,
			consentHistoryField,
			enabledField,
		},
	}
)

var structs = map[string]code.Struct{
	nameStruct.Name:  nameStruct,
	structModel.Name: structModel,
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
					code.TypeError,
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
					{
						Type: code.ExternalType{
							PackageAlias: "context",
							Name:         "Context",
						},
					},
					{
						Type: code.ArrayType{
							ContainedType: code.PointerType{
								ContainedType: code.SimpleType("UserModel"),
							},
						},
					},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.InterfaceType{}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.InsertOperation{
				Mode: spec.QueryModeMany,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(structs, structModel, testCase.Method)

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
				t.Errorf("Expected = %+v\nReceived = %+v", expectedOutput, actualSpec)
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
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindOneByMultiWordArg method",
			Method: code.Method{
				Name: "FindOneByPhoneNumber",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{phoneNumberField},
						Comparator:     spec.ComparatorEqual,
						ParamIndex:     1,
					},
				}},
			},
		},
		{
			Name: "FindByArg method",
			Method: code.Method{
				Name: "FindByCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindByDeepArg method",
			Method: code.Method{
				Name: "FindByNameFirst",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{nameField, firstNameField},
						Comparator:     spec.ComparatorEqual,
						ParamIndex:     1,
					},
				}},
			},
		},
		{
			Name: "FindByDeepPointerArg method",
			Method: code.Method{
				Name: "FindByReferrerID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{referrerField, idField},
						Comparator:     spec.ComparatorEqual,
						ParamIndex:     1,
					},
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
					code.TypeError,
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
					{Type: code.TypeString},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorAnd,
					Predicates: []spec.Predicate{
						{
							FieldReference: spec.FieldReference{cityField},
							Comparator:     spec.ComparatorEqual,
							ParamIndex:     1,
						},
						{
							FieldReference: spec.FieldReference{genderField},
							Comparator:     spec.ComparatorEqual,
							ParamIndex:     2,
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
					{Type: code.TypeString},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorOr,
					Predicates: []spec.Predicate{
						{
							FieldReference: spec.FieldReference{cityField},
							Comparator:     spec.ComparatorEqual,
							ParamIndex:     1,
						},
						{
							FieldReference: spec.FieldReference{genderField},
							Comparator:     spec.ComparatorEqual,
							ParamIndex:     2,
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
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorNot, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindByArgLessThan method",
			Method: code.Method{
				Name: "FindByAgeLessThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorLessThan, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindByArgLessThanEqual method",
			Method: code.Method{
				Name: "FindByAgeLessThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{ageField},
						Comparator:     spec.ComparatorLessThanEqual,
						ParamIndex:     1,
					},
				}},
			},
		},
		{
			Name: "FindByArgGreaterThan method",
			Method: code.Method{
				Name: "FindByAgeGreaterThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{ageField},
						Comparator:     spec.ComparatorGreaterThan,
						ParamIndex:     1,
					},
				}},
			},
		},
		{
			Name: "FindByArgGreaterThanEqual method",
			Method: code.Method{
				Name: "FindByAgeGreaterThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{ageField},
						Comparator:     spec.ComparatorGreaterThanEqual,
						ParamIndex:     1,
					},
				}},
			},
		},
		{
			Name: "FindByArgBetween method",
			Method: code.Method{
				Name: "FindByAgeBetween",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorBetween, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindByArgIn method",
			Method: code.Method{
				Name: "FindByCityIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.TypeString}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorIn, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "FindByArgNotIn method",
			Method: code.Method{
				Name: "FindByCityNotIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.TypeString}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorNotIn, ParamIndex: 1},
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
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{enabledField}, Comparator: spec.ComparatorTrue, ParamIndex: 1},
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
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{enabledField},
						Comparator:     spec.ComparatorFalse,
						ParamIndex:     1,
					},
				}},
			},
		},
		{
			Name: "FindByArgOrderByArg method",
			Method: code.Method{
				Name: "FindByCityOrderByAge",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
				Sorts: []spec.Sort{
					{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingAscending},
				},
			},
		},
		{
			Name: "FindByArgOrderByArgAsc method",
			Method: code.Method{
				Name: "FindByCityOrderByAgeAsc",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
				Sorts: []spec.Sort{
					{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingAscending},
				},
			},
		},
		{
			Name: "FindByArgOrderByArgDesc method",
			Method: code.Method{
				Name: "FindByCityOrderByAgeDesc",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
				Sorts: []spec.Sort{
					{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingDescending},
				},
			},
		},
		{
			Name: "FindByArgOrderByDeepArg method",
			Method: code.Method{
				Name: "FindByCityOrderByNameFirst",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
				Sorts: []spec.Sort{
					{FieldReference: spec.FieldReference{nameField, firstNameField}, Ordering: spec.OrderingAscending},
				},
			},
		},
		{
			Name: "FindByArgOrderByArgAndArg method",
			Method: code.Method{
				Name: "FindByCityOrderByCityAndAgeDesc",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedOperation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
				Sorts: []spec.Sort{
					{FieldReference: spec.FieldReference{cityField}, Ordering: spec.OrderingAscending},
					{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingDescending},
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(structs, structModel, testCase.Method)

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
				t.Errorf("Expected = %+v\nReceived = %+v", expectedOutput, actualSpec)
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
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateModel{},
				Mode:   spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
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
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateFields{
					spec.UpdateField{
						FieldReference: spec.FieldReference{genderField},
						ParamIndex:     1,
						Operator:       spec.UpdateOperatorSet,
					},
				},
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{idField},
						Comparator:     spec.ComparatorEqual,
						ParamIndex:     2,
					},
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
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateFields{
					spec.UpdateField{
						FieldReference: spec.FieldReference{genderField},
						ParamIndex:     1,
						Operator:       spec.UpdateOperatorSet,
					},
				},
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{idField},
						Comparator:     spec.ComparatorEqual,
						ParamIndex:     2,
					},
				}},
			},
		},
		{
			Name: "UpdateArgByArg one with deeply referenced update field method",
			Method: code.Method{
				Name: "UpdateNameFirstByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateFields{
					spec.UpdateField{
						FieldReference: spec.FieldReference{nameField, firstNameField},
						ParamIndex:     1,
						Operator:       spec.UpdateOperatorSet,
					},
				},
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{idField},
						Comparator:     spec.ComparatorEqual,
						ParamIndex:     2,
					},
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
					{Type: code.TypeString},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateFields{
					spec.UpdateField{
						FieldReference: spec.FieldReference{genderField},
						ParamIndex:     1,
						Operator:       spec.UpdateOperatorSet,
					},
					spec.UpdateField{
						FieldReference: spec.FieldReference{cityField},
						ParamIndex:     2,
						Operator:       spec.UpdateOperatorSet,
					},
				},
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 3},
				}},
			},
		},
		{
			Name: "UpdateArgPushByArg method",
			Method: code.Method{
				Name: "UpdateConsentHistoryPushByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("ConsentHistoryItem")},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateFields{
					spec.UpdateField{
						FieldReference: spec.FieldReference{consentHistoryField},
						ParamIndex:     1,
						Operator:       spec.UpdateOperatorPush,
					},
				},
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{idField},
						Comparator:     spec.ComparatorEqual,
						ParamIndex:     2,
					},
				}},
			},
		},
		{
			Name: "UpdateArgPushByArg method",
			Method: code.Method{
				Name: "UpdateAgeIncByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.UpdateOperation{
				Update: spec.UpdateFields{
					spec.UpdateField{
						FieldReference: spec.FieldReference{ageField},
						ParamIndex:     1,
						Operator:       spec.UpdateOperatorInc,
					},
				},
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{idField},
						Comparator:     spec.ComparatorEqual,
						ParamIndex:     2,
					},
				}},
			},
		},
		{
			Name: "UpdateArgAndArgPushByArg method",
			Method: code.Method{
				Name: "UpdateEnabledAndConsentHistoryPushByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeBool},
					{Type: code.SimpleType("ConsentHistoryItem")},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.UpdateOperation{
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
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 3},
				}},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(structs, structModel, testCase.Method)

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
				t.Errorf("Expected = %+v\nReceived = %+v", expectedOutput, actualSpec)
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
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteOneByMultiWordArg method",
			Method: code.Method{
				Name: "DeleteOneByPhoneNumber",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{phoneNumberField},
						Comparator:     spec.ComparatorEqual,
						ParamIndex:     1,
					},
				}},
			},
		},
		{
			Name: "DeleteByArg method",
			Method: code.Method{
				Name: "DeleteByCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
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
					code.TypeInt,
					code.TypeError,
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
					{Type: code.TypeString},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorAnd,
					Predicates: []spec.Predicate{
						{
							FieldReference: spec.FieldReference{cityField},
							Comparator:     spec.ComparatorEqual,
							ParamIndex:     1,
						},
						{
							FieldReference: spec.FieldReference{genderField},
							Comparator:     spec.ComparatorEqual,
							ParamIndex:     2,
						},
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
					{Type: code.TypeString},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorOr,
					Predicates: []spec.Predicate{
						{
							FieldReference: spec.FieldReference{cityField},
							Comparator:     spec.ComparatorEqual,
							ParamIndex:     1,
						},
						{
							FieldReference: spec.FieldReference{genderField},
							Comparator:     spec.ComparatorEqual,
							ParamIndex:     2,
						},
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
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorNot, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteByArgLessThan method",
			Method: code.Method{
				Name: "DeleteByAgeLessThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorLessThan, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteByArgLessThanEqual method",
			Method: code.Method{
				Name: "DeleteByAgeLessThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{ageField},
						Comparator:     spec.ComparatorLessThanEqual,
						ParamIndex:     1,
					},
				}},
			},
		},
		{
			Name: "DeleteByArgGreaterThan method",
			Method: code.Method{
				Name: "DeleteByAgeGreaterThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{ageField},
						Comparator:     spec.ComparatorGreaterThan,
						ParamIndex:     1,
					},
				}},
			},
		},
		{
			Name: "DeleteByArgGreaterThanEqual method",
			Method: code.Method{
				Name: "DeleteByAgeGreaterThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{
						FieldReference: spec.FieldReference{ageField},
						Comparator:     spec.ComparatorGreaterThanEqual,
						ParamIndex:     1,
					},
				}},
			},
		},
		{
			Name: "DeleteByArgBetween method",
			Method: code.Method{
				Name: "DeleteByAgeBetween",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorBetween, ParamIndex: 1},
				}},
			},
		},
		{
			Name: "DeleteByArgIn method",
			Method: code.Method{
				Name: "DeleteByCityIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.TypeString}},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.DeleteOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{Predicates: []spec.Predicate{
					{FieldReference: spec.FieldReference{cityField}, Comparator: spec.ComparatorIn, ParamIndex: 1},
				}},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(structs, structModel, testCase.Method)

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
				t.Errorf("Expected = %+v\nReceived = %+v", expectedOutput, actualSpec)
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
					code.TypeInt,
					code.TypeError,
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
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedOperation: spec.CountOperation{
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
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSpec, err := spec.ParseInterfaceMethod(structs, structModel, testCase.Method)

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
				t.Errorf("Expected = %+v\nReceived = %+v", expectedOutput, actualSpec)
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
	_, err := spec.ParseInterfaceMethod(structs, structModel, code.Method{
		Name: "SearchByID",
	})

	expectedError := spec.NewUnknownOperationError("Search")
	if !errors.Is(err, expectedError) {
		t.Errorf("\nExpected = %+v\nReceived = %+v", expectedError, err)
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
					code.TypeError,
				},
			},
			ExpectedError: spec.NewOperationReturnCountUnmatchedError(2),
		},
		{
			Name: "unsupported return types from insert method",
			Method: code.Method{
				Name: "Insert",
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewUnsupportedReturnError(
				code.PointerType{ContainedType: code.SimpleType("UserModel")},
				0,
			),
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
					code.TypeError,
				},
			},
			ExpectedError: spec.NewUnsupportedReturnError(code.InterfaceType{}, 0),
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
			ExpectedError: spec.NewUnsupportedReturnError(code.InterfaceType{}, 1),
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
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrContextParamRequired,
		},
		{
			Name: "mismatched model parameter for one mode",
			Method: code.Method{
				Name: "Insert",
				Params: []code.Param{
					{
						Name: "ctx",
						Type: code.ExternalType{PackageAlias: "context", Name: "Context"},
					},
					{
						Name: "userModel",
						Type: code.ArrayType{
							ContainedType: code.PointerType{
								ContainedType: code.SimpleType("UserModel"),
							},
						},
					},
				},
				Returns: []code.Type{
					code.InterfaceType{},
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrInvalidParam,
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
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrInvalidParam,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(structs, structModel, testCase.Method)

			if err.Error() != testCase.ExpectedError.Error() {
				t.Errorf("\nExpected = %+v\nReceived = %+v", testCase.ExpectedError, err)
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
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewOperationReturnCountUnmatchedError(2),
		},
		{
			Name: "unsupported return types from find method",
			Method: code.Method{
				Name: "FindOneByID",
				Returns: []code.Type{
					code.SimpleType("UserModel"),
					code.TypeError,
				},
			},
			ExpectedError: spec.NewUnsupportedReturnError(code.SimpleType("UserModel"), 0),
		},
		{
			Name: "error return not provided",
			Method: code.Method{
				Name: "FindOneByID",
				Returns: []code.Type{
					code.SimpleType("UserModel"),
					code.TypeInt,
				},
			},
			ExpectedError: spec.NewUnsupportedReturnError(code.TypeInt, 1),
		},
		{
			Name: "find method without query",
			Method: code.Method{
				Name: "Find",
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrQueryRequired,
		},
		{
			Name: "misplaced operator token (leftmost)",
			Method: code.Method{
				Name: "FindByAndGender",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"And", "Gender"}),
		},
		{
			Name: "misplaced operator token (rightmost)",
			Method: code.Method{
				Name: "FindByGenderAnd",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"Gender", "And"}),
		},
		{
			Name: "misplaced operator token (double operator)",
			Method: code.Method{
				Name: "FindByGenderAndAndCity",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"Gender", "And", "And", "City"}),
		},
		{
			Name: "ambiguous query",
			Method: code.Method{
				Name: "FindByGenderAndCityOrAge",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"Gender", "And", "City", "Or", "Age"}),
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
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrContextParamRequired,
		},
		{
			Name: "mismatched number of parameters",
			Method: code.Method{
				Name: "FindByCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrInvalidParam,
		},
		{
			Name: "struct field not found",
			Method: code.Method{
				Name: "FindByCountry",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError([]string{"Country"}),
		},
		{
			Name: "deeply referenced struct field not found",
			Method: code.Method{
				Name: "FindByNameMiddle",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError([]string{"Name", "Middle"}),
		},
		{
			Name: "deeply referenced struct not found",
			Method: code.Method{
				Name: "FindByContactPhone",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError([]string{"Contact", "Phone"}),
		},
		{
			Name: "deeply referenced external struct field",
			Method: code.Method{
				Name: "FindByDefaultPaymentMethod",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError([]string{"Default", "Payment", "Method"}),
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
					code.TypeError,
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
					code.TypeError,
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
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewArgumentTypeNotMatchedError(genderField.Name, genderField.Type, code.TypeString),
		},
		{
			Name: "mismatched method parameter type for special case",
			Method: code.Method{
				Name: "FindByCityIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewArgumentTypeNotMatchedError(cityField.Name,
				code.ArrayType{ContainedType: code.TypeString}, code.TypeString),
		},
		{
			Name: "misplaced operator token (leftmost)",
			Method: code.Method{
				Name: "FindAllOrderByAndAge",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
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
					code.TypeError,
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
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidSortError([]string{"Order", "By", "Age", "And", "And", "Gender"}),
		},
		{
			Name: "sort field not found",
			Method: code.Method{
				Name: "FindAllOrderByCountry",
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.TypeError,
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError([]string{"Country"}),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(structs, structModel, testCase.Method)

			if err.Error() != testCase.ExpectedError.Error() {
				t.Errorf("\nExpected = %+v\nReceived = %+v", testCase.ExpectedError.Error(), err.Error())
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
					code.TypeBool,
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewOperationReturnCountUnmatchedError(2),
		},
		{
			Name: "unsupported return types from update method",
			Method: code.Method{
				Name: "UpdateAgeByID",
				Returns: []code.Type{
					code.TypeFloat64,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewUnsupportedReturnError(code.TypeFloat64, 0),
		},
		{
			Name: "error return not provided",
			Method: code.Method{
				Name: "UpdateAgeByID",
				Returns: []code.Type{
					code.TypeBool,
					code.TypeBool,
				},
			},
			ExpectedError: spec.NewUnsupportedReturnError(code.TypeBool, 1),
		},
		{
			Name: "update with no field provided",
			Method: code.Method{
				Name: "UpdateByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrInvalidUpdateFields,
		},
		{
			Name: "misplaced And token in update fields",
			Method: code.Method{
				Name: "UpdateAgeAndAndGenderByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrInvalidUpdateFields,
		},
		{
			Name: "push operator in non-array field",
			Method: code.Method{
				Name: "UpdateGenderPushByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewIncompatibleUpdateOperatorError(spec.UpdateOperatorPush, spec.FieldReference{
				code.StructField{
					Name: "Gender",
					Type: code.SimpleType("Gender"),
				},
			}),
		},
		{
			Name: "inc operator in non-number field",
			Method: code.Method{
				Name: "UpdateCityIncByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewIncompatibleUpdateOperatorError(spec.UpdateOperatorInc, spec.FieldReference{
				code.StructField{
					Name: "City",
					Type: code.TypeString,
				},
			}),
		},
		{
			Name: "update method without query",
			Method: code.Method{
				Name: "UpdateCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrQueryRequired,
		},
		{
			Name: "ambiguous query",
			Method: code.Method{
				Name: "UpdateAgeByIDAndUsernameOrGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"ID", "And", "Username", "Or", "Gender"}),
		},
		{
			Name: "parameters for push operator is not array's contained type",
			Method: code.Method{
				Name: "UpdateConsentHistoryPushByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.SimpleType("ConsentHistoryItem")}},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewArgumentTypeNotMatchedError(
				consentHistoryField.Name,
				code.SimpleType("ConsentHistoryItem"),
				code.ArrayType{
					ContainedType: code.SimpleType("ConsentHistoryItem"),
				},
			),
		},
		{
			Name: "insufficient function parameters",
			Method: code.Method{
				Name: "UpdateEnabledAll",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					// {Type: code.SimpleType("Enabled")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrInvalidUpdateFields,
		},
		{
			Name: "update model with invalid parameter",
			Method: code.Method{
				Name: "UpdateByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeBool,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrInvalidUpdateFields,
		},
		{
			Name: "no context parameter",
			Method: code.Method{
				Name: "UpdateAgeByGender",
				Params: []code.Param{
					{Type: code.TypeInt},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrContextParamRequired,
		},
		{
			Name: "struct field not found in update fields",
			Method: code.Method{
				Name: "UpdateCountryByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError([]string{"Country"}),
		},
		{
			Name: "struct field does not match parameter in update fields",
			Method: code.Method{
				Name: "UpdateAgeByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeFloat64},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewArgumentTypeNotMatchedError(ageField.Name, ageField.Type, code.TypeFloat64),
		},
		{
			Name: "struct field does not match parameter in query",
			Method: code.Method{
				Name: "UpdateAgeByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeInt},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewArgumentTypeNotMatchedError(genderField.Name, genderField.Type, code.TypeString),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(structs, structModel, testCase.Method)

			if err.Error() != testCase.ExpectedError.Error() {
				t.Errorf("\nExpected = %+v\nReceived = %+v", testCase.ExpectedError, err)
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
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewOperationReturnCountUnmatchedError(2),
		},
		{
			Name: "unsupported return types from delete method",
			Method: code.Method{
				Name: "DeleteOneByID",
				Returns: []code.Type{
					code.TypeFloat64,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewUnsupportedReturnError(code.TypeFloat64, 0),
		},
		{
			Name: "error return not provided",
			Method: code.Method{
				Name: "DeleteOneByID",
				Returns: []code.Type{
					code.TypeInt,
					code.TypeBool,
				},
			},
			ExpectedError: spec.NewUnsupportedReturnError(code.TypeBool, 1),
		},
		{
			Name: "delete method without query",
			Method: code.Method{
				Name: "Delete",
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrQueryRequired,
		},
		{
			Name: "misplaced operator token (leftmost)",
			Method: code.Method{
				Name: "DeleteByAndGender",
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"And", "Gender"}),
		},
		{
			Name: "misplaced operator token (rightmost)",
			Method: code.Method{
				Name: "DeleteByGenderAnd",
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"Gender", "And"}),
		},
		{
			Name: "misplaced operator token (double operator)",
			Method: code.Method{
				Name: "DeleteByGenderAndAndCity",
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"Gender", "And", "And", "City"}),
		},
		{
			Name: "ambiguous query",
			Method: code.Method{
				Name: "DeleteByGenderAndCityOrAge",
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"Gender", "And", "City", "Or", "Age"}),
		},
		{
			Name: "no context parameter",
			Method: code.Method{
				Name: "DeleteByGender",
				Params: []code.Param{
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrContextParamRequired,
		},
		{
			Name: "mismatched number of parameters",
			Method: code.Method{
				Name: "DeleteByCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrInvalidParam,
		},
		{
			Name: "struct field not found",
			Method: code.Method{
				Name: "DeleteByCountry",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError([]string{"Country"}),
		},
		{
			Name: "mismatched method parameter type",
			Method: code.Method{
				Name: "DeleteByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewArgumentTypeNotMatchedError("Gender", code.SimpleType("Gender"), code.TypeString),
		},
		{
			Name: "mismatched method parameter type for special case",
			Method: code.Method{
				Name: "DeleteByCityIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewArgumentTypeNotMatchedError("City",
				code.ArrayType{ContainedType: code.TypeString}, code.TypeString),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(structs, structModel, testCase.Method)

			if err.Error() != testCase.ExpectedError.Error() {
				t.Errorf("\nExpected = %+v\nReceived = %+v", testCase.ExpectedError, err)
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
					code.TypeInt,
					code.TypeError,
					code.TypeBool,
				},
			},
			ExpectedError: spec.NewOperationReturnCountUnmatchedError(2),
		},
		{
			Name: "invalid integer return",
			Method: code.Method{
				Name: "CountAll",
				Returns: []code.Type{
					code.SimpleType("int64"),
					code.TypeError,
				},
			},
			ExpectedError: spec.NewUnsupportedReturnError(code.SimpleType("int64"), 0),
		},
		{
			Name: "error return not provided",
			Method: code.Method{
				Name: "CountAll",
				Returns: []code.Type{
					code.TypeInt,
					code.TypeBool,
				},
			},
			ExpectedError: spec.NewUnsupportedReturnError(code.TypeBool, 1),
		},
		{
			Name: "count method without query",
			Method: code.Method{
				Name: "Count",
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrQueryRequired,
		},
		{
			Name: "invalid query",
			Method: code.Method{
				Name: "CountBy",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewInvalidQueryError([]string{"By"}),
		},
		{
			Name: "context parameter not provided",
			Method: code.Method{
				Name: "CountAll",
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrContextParamRequired,
		},
		{
			Name: "mismatched number of parameter",
			Method: code.Method{
				Name: "CountByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.TypeInt},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.ErrInvalidParam,
		},
		{
			Name: "mismatched method parameter type",
			Method: code.Method{
				Name: "CountByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewArgumentTypeNotMatchedError("Gender", code.SimpleType("Gender"), code.TypeString),
		},
		{
			Name: "struct field not found",
			Method: code.Method{
				Name: "CountByCountry",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.TypeString},
				},
				Returns: []code.Type{
					code.TypeInt,
					code.TypeError,
				},
			},
			ExpectedError: spec.NewStructFieldNotFoundError([]string{"Country"}),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := spec.ParseInterfaceMethod(structs, structModel, testCase.Method)

			if err.Error() != testCase.ExpectedError.Error() {
				t.Errorf("\nExpected = %+v\nReceived = %+v", testCase.ExpectedError, err)
			}
		})
	}
}
