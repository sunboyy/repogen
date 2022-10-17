package codegen_test

import (
	"bytes"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/testutils"
)

func TestFunctionBuilderBuild_NoReturn(t *testing.T) {
	fb := codegen.FunctionBuilder{
		Name:    "init",
		Params:  nil,
		Returns: nil,
		Body:    `	logrus.SetLevel(logrus.DebugLevel)`,
	}
	expectedCode := `
func init() {
	logrus.SetLevel(logrus.DebugLevel)
}
`
	buffer := new(bytes.Buffer)

	err := fb.Impl(buffer)

	if err != nil {
		t.Fatal(err)
	}
	actual := buffer.String()
	if err := testutils.ExpectMultiLineString(
		expectedCode,
		actual,
	); err != nil {
		t.Error(err)
	}
}

func TestFunctionBuilderBuild_OneReturn(t *testing.T) {
	fb := codegen.FunctionBuilder{
		Name: "NewUser",
		Params: []code.Param{
			{
				Name: "username",
				Type: code.TypeString,
			},
			{
				Name: "age",
				Type: code.TypeInt,
			},
			{
				Name: "parent",
				Type: code.PointerType{ContainedType: code.SimpleType("User")},
			},
		},
		Returns: []code.Type{
			code.SimpleType("User"),
		},
		Body: `	return User{
		Username: username,
		Age: age,
		Parent: parent
	}`,
	}
	expectedCode := `
func NewUser(username string, age int, parent *User) User {
	return User{
		Username: username,
		Age: age,
		Parent: parent
	}
}
`
	buffer := new(bytes.Buffer)

	err := fb.Impl(buffer)

	if err != nil {
		t.Fatal(err)
	}
	actual := buffer.String()
	if err := testutils.ExpectMultiLineString(
		expectedCode,
		actual,
	); err != nil {
		t.Error(err)
	}
}

func TestFunctionBuilderBuild_MultiReturn(t *testing.T) {
	fb := codegen.FunctionBuilder{
		Name: "Save",
		Params: []code.Param{
			{
				Name: "user",
				Type: code.SimpleType("User"),
			},
		},
		Returns: []code.Type{
			code.SimpleType("User"),
			code.TypeError,
		},
		Body: `	return collection.Save(user)`,
	}
	expectedCode := `
func Save(user User) (User, error) {
	return collection.Save(user)
}
`
	buffer := new(bytes.Buffer)

	err := fb.Impl(buffer)

	if err != nil {
		t.Fatal(err)
	}
	actual := buffer.String()
	if err := testutils.ExpectMultiLineString(
		expectedCode,
		actual,
	); err != nil {
		t.Error(err)
	}
}
