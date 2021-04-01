package code_test

import (
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
)

func TestStructsByName(t *testing.T) {
	userStruct := code.Struct{
		Name: "UserModel",
		Fields: code.StructFields{
			{Name: "ID", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
			{Name: "Username", Type: code.SimpleType("string")},
		},
	}
	structs := code.Structs{userStruct}

	t.Run("struct found", func(t *testing.T) {
		structModel, ok := structs.ByName("UserModel")

		if !ok {
			t.Fail()
		}
		if !reflect.DeepEqual(structModel, userStruct) {
			t.Errorf("Expected = %+v\nReceived = %+v", userStruct, structModel)
		}
	})

	t.Run("struct not found", func(t *testing.T) {
		_, ok := structs.ByName("ProductModel")

		if ok {
			t.Fail()
		}
	})
}

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

func TestInterfacesByName(t *testing.T) {
	userRepoIntf := code.InterfaceType{Name: "UserRepository"}
	interfaces := code.Interfaces{userRepoIntf}

	t.Run("struct field found", func(t *testing.T) {
		intf, ok := interfaces.ByName("UserRepository")

		if !ok {
			t.Fail()
		}
		if !reflect.DeepEqual(intf, userRepoIntf) {
			t.Errorf("Expected = %+v\nReceived = %+v", userRepoIntf, intf)
		}
	})

	t.Run("struct field not found", func(t *testing.T) {
		_, ok := interfaces.ByName("Password")

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

func TestArrayTypeCode(t *testing.T) {
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
