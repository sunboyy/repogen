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
				Structs: []code.Struct{
					{
						Name: "UserModel",
						Fields: code.StructFields{
							code.StructField{
								Name: "ID",
								Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
								Tags: map[string][]string{
									"bson": {"_id", "omitempty"},
									"json": {"id"},
								},
							},
							code.StructField{
								Name: "Username",
								Type: code.TypeString,
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
	FindByID(ctx context.Context, id primitive.ObjectID) (*UserModel, error)
	FindAll(context.Context) ([]*UserModel, error)
	FindByAgeBetween(ctx context.Context, fromAge, toAge int) ([]*UserModel, error)
	InsertOne(ctx context.Context, user *UserModel) (interface{}, error)
	UpdateAgreementByID(ctx context.Context, agreement map[string]bool, id primitive.ObjectID) (bool, error)
	// CustomMethod does custom things.
	CustomMethod(interface {
		Run(arg1 int)
	}) interface {
		Do(arg2 string)
	}
}`,
			ExpectedOutput: code.File{
				PackageName: "user",
				Interfaces: []code.InterfaceType{
					{
						Name: "UserRepository",
						Methods: []code.Method{
							{
								Name: "FindByID",
								Params: []code.Param{
									{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
									{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
								},
								Returns: []code.Type{
									code.PointerType{ContainedType: code.SimpleType("UserModel")},
									code.TypeError,
								},
							},
							{
								Name: "FindAll",
								Params: []code.Param{
									{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
								},
								Returns: []code.Type{
									code.ArrayType{
										ContainedType: code.PointerType{
											ContainedType: code.SimpleType("UserModel"),
										},
									},
									code.TypeError,
								},
							},
							{
								Name: "FindByAgeBetween",
								Params: []code.Param{
									{
										Name: "ctx",
										Type: code.ExternalType{PackageAlias: "context", Name: "Context"},
									},
									{
										Name: "fromAge",
										Type: code.TypeInt,
									},
									{
										Name: "toAge",
										Type: code.TypeInt,
									},
								},
								Returns: []code.Type{
									code.ArrayType{
										ContainedType: code.PointerType{
											ContainedType: code.SimpleType("UserModel"),
										},
									},
									code.TypeError,
								},
							},
							{
								Name: "InsertOne",
								Params: []code.Param{
									{
										Name: "ctx",
										Type: code.ExternalType{PackageAlias: "context", Name: "Context"},
									},
									{
										Name: "user",
										Type: code.PointerType{ContainedType: code.SimpleType("UserModel")},
									},
								},
								Returns: []code.Type{
									code.InterfaceType{},
									code.TypeError,
								},
							},
							{
								Name: "UpdateAgreementByID",
								Params: []code.Param{
									{
										Name: "ctx",
										Type: code.ExternalType{PackageAlias: "context", Name: "Context"},
									},
									{
										Name: "agreement",
										Type: code.MapType{KeyType: code.TypeString, ValueType: code.TypeBool},
									},
									{
										Name: "id",
										Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
									},
								},
								Returns: []code.Type{
									code.TypeBool,
									code.TypeError,
								},
							},
							{
								Name:     "CustomMethod",
								Comments: []string{"CustomMethod does custom things."},
								Params: []code.Param{
									{
										Type: code.InterfaceType{
											Methods: []code.Method{
												{
													Name: "Run",
													Params: []code.Param{
														{Name: "arg1", Type: code.TypeInt},
													},
												},
											},
										},
									},
								},
								Returns: []code.Type{
									code.InterfaceType{
										Methods: []code.Method{
											{
												Name: "Do",
												Params: []code.Param{
													{Name: "arg2", Type: code.TypeString},
												},
											},
										},
									},
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
	FindByID(ctx context.Context, id primitive.ObjectID) (*UserModel, error)
	FindAll(ctx context.Context) ([]*UserModel, error)
}
`,
			ExpectedOutput: code.File{
				PackageName: "user",
				Imports: []code.Import{
					{Path: "context"},
					{Path: "go.mongodb.org/mongo-driver/bson/primitive"},
				},
				Structs: []code.Struct{
					{
						Name: "UserModel",
						Fields: code.StructFields{
							code.StructField{
								Name: "ID",
								Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
								Tags: map[string][]string{
									"bson": {"_id", "omitempty"},
									"json": {"id"},
								},
							},
							code.StructField{
								Name: "Username",
								Type: code.TypeString,
								Tags: map[string][]string{
									"bson": {"username"},
									"json": {"username"},
								},
							},
						},
					},
				},
				Interfaces: []code.InterfaceType{
					{
						Name: "UserRepository",
						Methods: []code.Method{
							{
								Name: "FindByID",
								Params: []code.Param{
									{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
									{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
								},
								Returns: []code.Type{
									code.PointerType{ContainedType: code.SimpleType("UserModel")},
									code.TypeError,
								},
							},
							{
								Name: "FindAll",
								Params: []code.Param{
									{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
								},
								Returns: []code.Type{
									code.ArrayType{
										ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")},
									},
									code.TypeError,
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
				t.Errorf("Expected = %+v\nReceived = %+v", testCase.ExpectedOutput, file)
			}
		})
	}
}
