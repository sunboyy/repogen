package code_test

import (
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
)

func TestStructFieldsByName(t *testing.T) {
	idField := code.StructField{Name: "ID", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}}
	usernameField := code.StructField{Name: "Username", Type: code.SimpleType("string")}
	fields := code.StructFields{idField, usernameField}

	t.Run("struct field found", func(t *testing.T) {
		field, ok := fields.ByName("Username")

		if !ok {
			t.Fail()
		}
		if !reflect.DeepEqual(field, usernameField) {
			t.Errorf("Expected = %+v\nReceived = %+v", usernameField, field)
		}
	})

	t.Run("struct field not found", func(t *testing.T) {
		_, ok := fields.ByName("Password")

		if ok {
			t.Fail()
		}
	})
}

type TypeCodeTestCase struct {
	Name         string
	Type         code.Type
	ExpectedCode string
}

func TestTypeCode(t *testing.T) {
	testTable := []TypeCodeTestCase{
		{
			Name:         "simple type",
			Type:         code.SimpleType("UserModel"),
			ExpectedCode: "UserModel",
		},
		{
			Name:         "external type",
			Type:         code.ExternalType{PackageAlias: "context", Name: "Context"},
			ExpectedCode: "context.Context",
		},
		{
			Name:         "pointer type",
			Type:         code.PointerType{ContainedType: code.SimpleType("UserModel")},
			ExpectedCode: "*UserModel",
		},
		{
			Name:         "array type",
			Type:         code.ArrayType{ContainedType: code.SimpleType("UserModel")},
			ExpectedCode: "[]UserModel",
		},
		{
			Name: "map type",
			Type: code.MapType{
				KeyType:   code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
				ValueType: code.PointerType{ContainedType: code.SimpleType("UserModel")},
			},
			ExpectedCode: "map[primitive.ObjectID]*UserModel",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			code := testCase.Type.Code()

			if code != testCase.ExpectedCode {
				t.Errorf("Expected = %+v\nReceived = %+v", testCase.ExpectedCode, code)
			}
		})
	}
}

type TypeIsNumberTestCase struct {
	Name     string
	Type     code.Type
	IsNumber bool
}

func TestTypeIsNumber(t *testing.T) {
	testTable := []TypeIsNumberTestCase{
		{
			Name:     "simple type: int",
			Type:     code.SimpleType("int"),
			IsNumber: true,
		},
		{
			Name:     "simple type: other integer variants",
			Type:     code.SimpleType("int64"),
			IsNumber: true,
		},
		{
			Name:     "simple type: uint",
			Type:     code.SimpleType("uint"),
			IsNumber: true,
		},
		{
			Name:     "simple type: other unsigned integer variants",
			Type:     code.SimpleType("uint64"),
			IsNumber: true,
		},
		{
			Name:     "simple type: float32",
			Type:     code.SimpleType("float32"),
			IsNumber: true,
		},
		{
			Name:     "simple type: other float variants",
			Type:     code.SimpleType("float64"),
			IsNumber: true,
		},
		{
			Name:     "simple type: non-number primitive type",
			Type:     code.SimpleType("string"),
			IsNumber: false,
		},
		{
			Name:     "simple type: non-number custom type",
			Type:     code.SimpleType("UserModel"),
			IsNumber: false,
		},
		{
			Name:     "external type",
			Type:     code.ExternalType{PackageAlias: "context", Name: "Context"},
			IsNumber: false,
		},
		{
			Name:     "pointer type: number",
			Type:     code.PointerType{ContainedType: code.SimpleType("int")},
			IsNumber: true,
		},
		{
			Name:     "pointer type: non-number",
			Type:     code.PointerType{ContainedType: code.SimpleType("string")},
			IsNumber: false,
		},
		{
			Name:     "array type",
			Type:     code.ArrayType{ContainedType: code.SimpleType("int")},
			IsNumber: false,
		},
		{
			Name:     "map type",
			Type:     code.MapType{KeyType: code.SimpleType("int"), ValueType: code.SimpleType("float64")},
			IsNumber: false,
		},
		{
			Name:     "interface type",
			Type:     code.InterfaceType{},
			IsNumber: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			isNumber := testCase.Type.IsNumber()

			if isNumber != testCase.IsNumber {
				t.Errorf("Expected = %+v\nReceived = %+v", testCase.IsNumber, isNumber)
			}
		})
	}
}
