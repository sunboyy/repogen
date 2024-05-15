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
		Pkg:  testutils.Pkg,
		Name: "NewUser",
		Params: types.NewTuple(
			types.NewVar(token.NoPos, nil, "username", code.TypeString),
			types.NewVar(token.NoPos, nil, "age", code.TypeInt),
			types.NewVar(token.NoPos, nil, "parent", types.NewPointer(testutils.TypeUserNamed)),
		),
		Returns: []types.Type{
			testutils.TypeUserNamed,
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
		Pkg:  testutils.Pkg,
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

func TestTypeToString(t *testing.T) {
	internalPkg := types.NewPackage("github.com/sunboyy/repogen/internal/foo", "foo")
	externalPkg := types.NewPackage("github.com/sunboyy/repogen/internal/bar", "bar")

	tests := []struct {
		name string
		typ  types.Type
		want string
	}{
		{
			name: "basic type",
			typ:  code.TypeString,
			want: "string",
		},
		{
			name: "pointer type",
			typ:  types.NewPointer(code.TypeString),
			want: "*string",
		},
		{
			name: "slice type",
			typ:  types.NewSlice(code.TypeString),
			want: "[]string",
		},
		{
			name: "named type internal",
			typ:  types.NewNamed(types.NewTypeName(token.NoPos, internalPkg, "User", nil), nil, nil),
			want: "User",
		},
		{
			name: "named type external",
			typ:  types.NewNamed(types.NewTypeName(token.NoPos, externalPkg, "User", nil), nil, nil),
			want: "bar.User",
		},
		{
			name: "integration internal",
			typ: types.NewSlice(
				types.NewPointer(
					types.NewNamed(types.NewTypeName(token.NoPos, internalPkg, "User", nil), nil, nil),
				),
			),
			want: "[]*User",
		},
		{
			name: "integration external",
			typ: types.NewSlice(
				types.NewPointer(
					types.NewNamed(types.NewTypeName(token.NoPos, externalPkg, "User", nil), nil, nil),
				),
			),
			want: "[]*bar.User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := codegen.TypeToString(internalPkg, tt.typ); got != tt.want {
				t.Errorf("TypeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeToString_PkgNil(t *testing.T) {
	externalPkg := types.NewPackage("github.com/sunboyy/repogen/internal/bar", "bar")

	tests := []struct {
		name string
		typ  types.Type
		want string
	}{
		{
			name: "basic type",
			typ:  code.TypeString,
			want: "string",
		},
		{
			name: "pointer type",
			typ:  types.NewPointer(code.TypeString),
			want: "*string",
		},
		{
			name: "slice type",
			typ:  types.NewSlice(code.TypeString),
			want: "[]string",
		},
		{
			name: "named type external",
			typ:  types.NewNamed(types.NewTypeName(token.NoPos, externalPkg, "User", nil), nil, nil),
			want: "bar.User",
		},
		{
			name: "integration external",
			typ: types.NewSlice(
				types.NewPointer(
					types.NewNamed(types.NewTypeName(token.NoPos, externalPkg, "User", nil), nil, nil),
				),
			),
			want: "[]*bar.User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := codegen.TypeToString(nil, tt.typ); got != tt.want {
				t.Errorf("TypeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
