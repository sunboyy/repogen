package mongo_test

import (
	"fmt"
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

func createSignature(params []*types.Var, results []*types.Var) *types.Signature {
	return types.NewSignatureType(nil, nil, nil, types.NewTuple(params...), types.NewTuple(results...), false)
}

func createTypeVar(t types.Type) *types.Var {
	return types.NewVar(token.NoPos, nil, "", t)
}

func TestGenerateMethod_Insert(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "insert one method",
			MethodSpec: spec.MethodSpec{
				Name: "InsertOne",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(types.NewPointer(testutils.TypeUserNamed)),
					},
					[]*types.Var{
						createTypeVar(types.NewInterfaceType(nil, nil)),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.InsertOperation{
					Mode: spec.QueryModeOne,
				},
			},
			ExpectedBody: `	result, err := r.collection.InsertOne(arg0, arg1)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil`,
		},
		{
			Name: "insert many method",
			MethodSpec: spec.MethodSpec{
				Name: "Insert",
				Signature: createSignature(
					[]*types.Var{
						createTypeVar(testutils.TypeContextNamed),
						createTypeVar(types.NewSlice(types.NewPointer(testutils.TypeUserNamed))),
					},
					[]*types.Var{
						createTypeVar(types.NewSlice(types.NewInterfaceType(nil, nil))),
						createTypeVar(code.TypeError),
					},
				),
				Operation: spec.InsertOperation{
					Mode: spec.QueryModeMany,
				},
			},
			ExpectedBody: `	var entities []interface{}
	for _, model := range arg1 {
		entities = append(entities, model)
	}
	result, err := r.collection.InsertMany(arg0, entities)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(testutils.Pkg, "User", "UserRepository")
			expectedReceiver := codegen.MethodReceiver{
				Name:    "r",
				Type:    "UserRepositoryMongo",
				Pointer: true,
			}

			params := testCase.MethodSpec.Signature.Params()
			var expectedParamVars []*types.Var
			for i := 0; i < params.Len(); i++ {
				expectedParamVars = append(expectedParamVars, types.NewVar(token.NoPos, nil, fmt.Sprintf("arg%d", i),
					params.At(i).Type()))
			}
			expectedParams := types.NewTuple(expectedParamVars...)
			returns := testCase.MethodSpec.Signature.Results()
			var expectedReturns []types.Type
			for i := 0; i < returns.Len(); i++ {
				expectedReturns = append(expectedReturns, returns.At(i).Type())
			}

			actual, err := generator.GenerateMethod(testCase.MethodSpec)

			if err != nil {
				t.Fatal(err)
			}
			if expectedReceiver != actual.Receiver {
				t.Errorf(
					"incorrect method receiver: expected %+v, got %+v",
					expectedReceiver,
					actual.Receiver,
				)
			}
			if testCase.MethodSpec.Name != actual.Name {
				t.Errorf(
					"incorrect method name: expected %s, got %s",
					testCase.MethodSpec.Name,
					actual.Name,
				)
			}
			if !reflect.DeepEqual(expectedParams, actual.Params) {
				t.Errorf(
					"incorrect struct params: expected %+v, got %+v",
					expectedParams,
					actual.Params,
				)
			}
			if !reflect.DeepEqual(expectedReturns, actual.Returns) {
				t.Errorf(
					"incorrect struct returns: expected %+v, got %+v",
					expectedReturns,
					actual.Returns,
				)
			}
			if err := testutils.ExpectMultiLineString(testCase.ExpectedBody, actual.Body.Code()); err != nil {
				t.Error(err)
			}
		})
	}
}
