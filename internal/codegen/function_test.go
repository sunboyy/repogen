package codegen_test

import (
	"bytes"
	"go/token"
	"go/types"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/testutils"
)

func TestFunctionBuilderBuild_NoReturn(t *testing.T) {
	fb := codegen.FunctionBuilder{
		Name:    "init",
		Params:  types.NewTuple(),
		Returns: nil,
		Body: codegen.FunctionBody{
			codegen.ChainStatement{
				codegen.Identifier("logrus"),
				codegen.CallStatement{
					FuncName: "SetLevel",
					Params: codegen.StatementList{
						codegen.ChainStatement{
							codegen.Identifier("logrus"),
							codegen.Identifier("DebugLevel"),
						},
					},
				},
			},
		},
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
		Params: types.NewTuple(
			types.NewVar(token.NoPos, nil, "username", code.TypeString),
			types.NewVar(token.NoPos, nil, "age", code.TypeInt),
			types.NewVar(token.NoPos, nil, "parent",
				types.NewPointer(types.NewNamed(types.NewTypeName(token.NoPos, nil, "User", nil), nil, nil))),
		),
		Returns: []types.Type{
			types.NewNamed(types.NewTypeName(token.NoPos, nil, "User", nil), nil, nil),
		},
		Body: codegen.FunctionBody{
			codegen.ReturnStatement{
				codegen.StructStatement{
					Type: "User",
					Pairs: []codegen.StructFieldPair{
						{Key: "Username", Value: codegen.Identifier("username")},
						{Key: "Age", Value: codegen.Identifier("age")},
						{Key: "Parent", Value: codegen.Identifier("parent")},
					},
				},
			},
		},
	}
	expectedCode := `
func NewUser(username string, age int, parent *User) User {
	return User{
		Username: username,
		Age: age,
		Parent: parent,
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
		Params: types.NewTuple(
			types.NewVar(token.NoPos, nil, "user",
				types.NewNamed(types.NewTypeName(token.NoPos, nil, "User", nil), nil, nil)),
		),
		Returns: []types.Type{
			types.NewNamed(types.NewTypeName(token.NoPos, nil, "User", nil), nil, nil),
			code.TypeError,
		},
		Body: codegen.FunctionBody{
			codegen.ReturnStatement{
				codegen.ChainStatement{
					codegen.Identifier("collection"),
					codegen.CallStatement{
						FuncName: "Save",
						Params: codegen.StatementList{
							codegen.Identifier("user"),
						},
					},
				},
			},
		},
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
