package code_test

import (
	"go/parser"
	"go/token"
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
)

type TestCase struct {
	Name           string
	Source         string
	ExpectedOutput code.File
}

func TestExtractComponents(t *testing.T) {
	testTable := []TestCase{
		{
			Name:   "package name",
			Source: `package user`,
			ExpectedOutput: code.File{
				PackageName: "user",
			},
		},
		{
			Name: "single line imports",
			Source: `package user

import ctx "context"
import "go.mongodb.org/mongo-driver/bson/primitive"`,
			ExpectedOutput: code.File{
				PackageName: "user",
				Imports: []code.Import{
					{Name: "ctx", Path: "context"},
					{Path: "go.mongodb.org/mongo-driver/bson/primitive"},
				},
			},
		},
		{
			Name: "multiple line imports",
			Source: `package user

import (
	ctx "context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)`,
			ExpectedOutput: code.File{
				PackageName: "user",
				Imports: []code.Import{
					{Name: "ctx", Path: "context"},
					{Path: "go.mongodb.org/mongo-driver/bson/primitive"},
				},
			},
		},
		{
			Name: "struct declaration",
			Source: `package user

type UserModel struct {
	ID       primitive.ObjectID ` + "`bson:\"_id,omitempty\" json:\"id\"`" + `
	Username string             ` + "`bson:\"username\" json:\"username\"`" + `
}`,
			ExpectedOutput: code.File{
				PackageName: "user",
				Structs: code.Structs{
					{
						Name: "UserModel",
						Fields: code.StructFields{
							{
								Name: "ID",
								Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
								Tags: map[string][]string{
									"bson": {"_id", "omitempty"},
									"json": {"id"},
								},
							},
							{
								Name: "Username",
								Type: code.SimpleType("string"),
								Tags: map[string][]string{
									"bson": {"username"},
									"json": {"username"},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "interface declaration",
			Source: `package user

type UserRepository interface {
	FindOneByID(ctx context.Context, id primitive.ObjectID) (*UserModel, error)
	FindAll(context.Context) ([]*UserModel, error)
}`,
			ExpectedOutput: code.File{
				PackageName: "user",
				Interfaces: code.Interfaces{
					{
						Name: "UserRepository",
						Methods: []code.Method{
							{
								Name: "FindOneByID",
								Params: []code.Param{
									{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
									{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
								},
								Returns: []code.Type{
									code.PointerType{ContainedType: code.SimpleType("UserModel")},
									code.SimpleType("error"),
								},
							},
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
				},
			},
		},
		{
			Name: "integration",
			Source: `package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	ID       primitive.ObjectID ` + "`bson:\"_id,omitempty\" json:\"id\"`" + `
	Username string             ` + "`bson:\"username\" json:\"username\"`" + `
}

type UserRepository interface {
	FindOneByID(ctx context.Context, id primitive.ObjectID) (*UserModel, error)
	FindAll(ctx context.Context) ([]*UserModel, error)
}
`,
			ExpectedOutput: code.File{
				PackageName: "user",
				Imports: []code.Import{
					{Path: "context"},
					{Path: "go.mongodb.org/mongo-driver/bson/primitive"},
				},
				Structs: code.Structs{
					{
						Name: "UserModel",
						Fields: code.StructFields{
							{
								Name: "ID",
								Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
								Tags: map[string][]string{
									"bson": {"_id", "omitempty"},
									"json": {"id"},
								},
							},
							{
								Name: "Username",
								Type: code.SimpleType("string"),
								Tags: map[string][]string{
									"bson": {"username"},
									"json": {"username"},
								},
							},
						},
					},
				},
				Interfaces: code.Interfaces{
					{
						Name: "UserRepository",
						Methods: []code.Method{
							{
								Name: "FindOneByID",
								Params: []code.Param{
									{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
									{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
								},
								Returns: []code.Type{
									code.PointerType{ContainedType: code.SimpleType("UserModel")},
									code.SimpleType("error"),
								},
							},
							{
								Name: "FindAll",
								Params: []code.Param{
									{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
								},
								Returns: []code.Type{
									code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
									code.SimpleType("error"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, _ := parser.ParseFile(fset, "", testCase.Source, parser.ParseComments)

			file := code.ExtractComponents(f)

			if !reflect.DeepEqual(file, testCase.ExpectedOutput) {
				t.Errorf("Expected = %v\nReceived = %v", testCase.ExpectedOutput, file)
			}
		})
	}
}
