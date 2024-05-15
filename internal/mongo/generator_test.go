package mongo_test

import (
	"errors"
	"go/token"
	"go/types"
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
	"github.com/sunboyy/repogen/internal/testutils"
)

func TestImports(t *testing.T) {
	generator := mongo.NewGenerator(testutils.Pkg, testutils.TypeUserNamed, "UserRepository")
	expected := [][]codegen.Import{
		{
			{Path: "context"},
		},
		{
			{Path: "go.mongodb.org/mongo-driver/bson"},
			{Path: "go.mongodb.org/mongo-driver/bson/primitive"},
			{Path: "go.mongodb.org/mongo-driver/mongo"},
			{Path: "go.mongodb.org/mongo-driver/mongo/options"},
		},
	}

	actual := generator.Imports()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("incorrect imports: expected %+v, got %+v", expected, actual)
	}
}

func TestGenerateStruct(t *testing.T) {
	bareMongoPkg := types.NewPackage("go.mongodb.org/mongo-driver/mongo", "mongo")
	bareCollectionType := types.NewNamed(types.NewTypeName(token.NoPos, bareMongoPkg, "Collection", nil), nil, nil)
	generator := mongo.NewGenerator(testutils.Pkg, testutils.TypeUserNamed, "UserRepository")
	expected := codegen.StructBuilder{
		Name: "UserRepositoryMongo",
		Fields: []code.StructField{
			{
				Var: types.NewVar(token.NoPos, nil, "collection", types.NewPointer(bareCollectionType)),
			},
		},
	}

	actual := generator.GenerateStruct()

	if expected.Name != actual.Name {
		t.Errorf(
			"incorrect struct name: expected %s, got %s",
			expected.Name,
			actual.Name,
		)
	}
	if !reflect.DeepEqual(expected.Fields, actual.Fields) {
		t.Errorf(
			"incorrect struct fields: expected %+v, got %+v",
			expected.Fields,
			actual.Fields,
		)
	}
}

func TestGenerateConstructor(t *testing.T) {
	generator := mongo.NewGenerator(testutils.Pkg, testutils.TypeUserNamed, "UserRepository")
	expected := codegen.FunctionBuilder{
		Name: "NewUserRepository",
		Params: types.NewTuple(
			types.NewVar(token.NoPos, nil, "collection", types.NewPointer(testutils.TypeCollectionNamed)),
		),
		Returns: []types.Type{
			types.NewNamed(types.NewTypeName(token.NoPos, nil, "UserRepository", nil), nil, nil),
		},
		Body: codegen.FunctionBody{
			codegen.ReturnStatement{
				codegen.StructStatement{
					Type: "&UserRepositoryMongo",
					Pairs: []codegen.StructFieldPair{{
						Key:   "collection",
						Value: codegen.Identifier("collection"),
					}},
				},
			},
		},
	}

	actual, err := generator.GenerateConstructor()

	if err != nil {
		t.Fatal(err)
	}
	if expected.Name != actual.Name {
		t.Errorf(
			"incorrect function name: expected %s, got %s",
			expected.Name,
			actual.Name,
		)
	}
	if expected.Params.Len() != actual.Params.Len() {
		t.Errorf(
			"incorrect function params length: expected %d, got %d",
			expected.Params.Len(),
			actual.Params.Len(),
		)
	}
	for i := 0; i < expected.Params.Len(); i++ {
		if expected.Params.At(i).Name() != actual.Params.At(i).Name() {
			t.Errorf(
				"incorrect function param name: expected %s, got %s",
				expected.Params.At(i).Name(),
				actual.Params.At(i).Name(),
			)
		}
		if expected.Params.At(i).Type().String() != actual.Params.At(i).Type().String() {
			t.Errorf(
				"incorrect function param type at %d: expected %s, got %s",
				i,
				expected.Params.At(i).Type(),
				actual.Params.At(i).Type(),
			)
		}
	}
	if !reflect.DeepEqual(expected.Body, actual.Body) {
		t.Errorf("incorrect function body: expected %+v got %+v",
			expected.Body,
			actual.Body,
		)
	}
}

type GenerateMethodTestCase struct {
	Name         string
	MethodSpec   spec.MethodSpec
	ExpectedBody string
}

type GenerateMethodInvalidTestCase struct {
	Name          string
	Method        spec.MethodSpec
	ExpectedError error
}

type StubOperation struct {
}

func (o StubOperation) Name() string {
	return "Stub"
}

func TestGenerateMethod_Invalid(t *testing.T) {
	testTable := []GenerateMethodInvalidTestCase{
		{
			Name: "operation not supported",
			Method: spec.MethodSpec{
				Name: "SearchByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewPointer(testutils.TypeUserNamed)),
						createTypeVar(code.TypeError),
					},
				),
				Operation: StubOperation{},
			},
			ExpectedError: mongo.NewOperationNotSupportedError("Stub"),
		},
		{
			Name: "bson tag not found in query",
			Method: spec.MethodSpec{
				Name: "FindByAccessToken",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeString),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserNamed))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "AccessToken"),
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 1,
							},
						},
					},
				},
			},
			ExpectedError: mongo.NewBsonTagNotFoundError("AccessToken"),
		},
		{
			Name: "bson tag not found in sort",
			Method: spec.MethodSpec{
				Name: "FindAllOrderByAccessToken",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserNamed))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Sorts: []spec.Sort{
						{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "AccessToken"),
								},
							},
							Ordering: spec.OrderingAscending,
						},
					},
				},
			},
			ExpectedError: mongo.NewBsonTagNotFoundError("AccessToken"),
		},
		{
			Name: "bson tag not found in update field",
			Method: spec.MethodSpec{
				Name: "UpdateAccessTokenByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeString),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeBool),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "AccessToken"),
								},
							},
							ParamIndex: 1,
							Operator:   spec.UpdateOperatorSet,
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
										Tag: `bson:"_id,omitempty"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedError: mongo.NewBsonTagNotFoundError("AccessToken"),
		},
		{
			Name: "update type not supported",
			Method: spec.MethodSpec{
				Name: "UpdateAgeByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeBool),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.UpdateOperation{
					Update: StubUpdate{},
					Mode:   spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
										Tag: `bson:"_id,omitempty"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedError: mongo.NewUpdateTypeNotSupportedError(StubUpdate{}),
		},
		{
			Name: "update operator not supported",
			Method: spec.MethodSpec{
				Name: "UpdateConsentHistoryAppendByID",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(code.TypeInt),
						createTypeVar(testutils.TypeObjectIDNamed),
					},
					[]*types.Var{
						createTypeVar(code.TypeBool),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						spec.UpdateField{
							FieldReference: spec.FieldReference{
								{
									Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ConsentHistory"),
									Tag: `bson:"consent_history"`,
								},
							},
							ParamIndex: 1,
							Operator:   "APPEND",
						},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								FieldReference: spec.FieldReference{
									{
										Var: testutils.FindStructFieldByName(testutils.TypeUserStruct, "ID"),
										Tag: `bson:"_id,omitempty"`,
									},
								},
								Comparator: spec.ComparatorEqual,
								ParamIndex: 2,
							},
						},
					},
				},
			},
			ExpectedError: mongo.NewUpdateOperatorNotSupportedError("APPEND"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(testutils.Pkg, testutils.TypeUserNamed, "UserRepository")

			_, err := generator.GenerateMethod(testCase.Method)

			if !errors.Is(err, testCase.ExpectedError) {
				t.Errorf("\nExpected = %+v\nReceived = %+v", testCase.ExpectedError, err)
			}
		})
	}
}
