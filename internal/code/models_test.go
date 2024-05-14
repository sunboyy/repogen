package code_test

import (
	"testing"

	"github.com/sunboyy/repogen/internal/code"
)

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
