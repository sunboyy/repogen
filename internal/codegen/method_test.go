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

func TestMethodBuilderBuild_IgnoreReceiverNoReturn(t *testing.T) {
	fb := codegen.MethodBuilder{
		Receiver: codegen.MethodReceiver{Type: "User"},
		Name:     "Init",
		Params:   nil,
		Returns:  nil,
		Body: codegen.FunctionBody{
			codegen.ChainStatement{
				codegen.Identifier("db"),
				codegen.CallStatement{
					FuncName: "Init",
					Params: codegen.StatementList{
						codegen.RawStatement("&User{}"),
					},
				},
			},
		},
	}
	expectedCode := `
func (User) Init() {
	db.Init(&User{})
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

func TestMethodBuilderBuild_IgnorePoinerReceiverOneReturn(t *testing.T) {
	fb := codegen.MethodBuilder{
		Receiver: codegen.MethodReceiver{
			Type:    "User",
			Pointer: true,
		},
		Name:    "Init",
		Params:  nil,
		Returns: []types.Type{code.TypeError},
		Body: codegen.FunctionBody{
			codegen.ReturnStatement{
				codegen.ChainStatement{
					codegen.Identifier("db"),
					codegen.CallStatement{
						FuncName: "Init",
						Params: codegen.StatementList{
							codegen.RawStatement("&User{}"),
						},
					},
				},
			},
		},
	}
	expectedCode := `
func (*User) Init() error {
	return db.Init(&User{})
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

func TestMethodBuilderBuild_UseReceiverMultiReturn(t *testing.T) {
	fb := codegen.MethodBuilder{
		Receiver: codegen.MethodReceiver{
			Name: "u",
			Type: "User",
		},
		Name: "WithAge",
		Params: types.NewTuple(
			types.NewVar(token.NoPos, nil, "age", code.TypeInt),
		),
		Returns: []types.Type{
			types.NewNamed(types.NewTypeName(token.NoPos, nil, "User", nil), nil, nil),
			code.TypeError,
		},
		Body: codegen.FunctionBody{
			codegen.AssignStatement{
				Vars: []string{"u.Age"},
				Values: codegen.StatementList{
					codegen.Identifier("age"),
				},
			},
			codegen.ReturnStatement{
				codegen.Identifier("u"),
			},
		},
	}
	expectedCode := `
func (u User) WithAge(age int) (User, error) {
	u.Age = age
	return u
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

func TestMethodBuilderBuild_UsePointerReceiverNoReturn(t *testing.T) {
	fb := codegen.MethodBuilder{
		Receiver: codegen.MethodReceiver{
			Name:    "u",
			Type:    "User",
			Pointer: true,
		},
		Name: "SetAge",
		Params: types.NewTuple(
			types.NewVar(token.NoPos, nil, "age", code.TypeInt),
		),
		Returns: nil,
		Body: codegen.FunctionBody{
			codegen.AssignStatement{
				Vars: []string{"u.Age"},
				Values: codegen.StatementList{
					codegen.Identifier("age"),
				},
			},
		},
	}
	expectedCode := `
func (u *User) SetAge(age int) {
	u.Age = age
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
