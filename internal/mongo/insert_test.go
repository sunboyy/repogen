package mongo_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
	"github.com/sunboyy/repogen/internal/testutils"
)

func TestGenerateMethod_Insert(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "insert one method",
			MethodSpec: spec.MethodSpec{
				Name: "InsertOne",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "userModel", Type: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				},
				Returns: []code.Type{
					code.InterfaceType{},
					code.TypeError,
				},
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
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "userModel", Type: code.ArrayType{
						ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")},
					}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.InterfaceType{}},
					code.TypeError,
				},
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
			generator := mongo.NewGenerator(userModel, "UserRepository")
			expectedReceiver := codegen.MethodReceiver{
				Name:    "r",
				Type:    "UserRepositoryMongo",
				Pointer: true,
			}
			var expectedParams []code.Param
			for i, param := range testCase.MethodSpec.Params {
				expectedParams = append(expectedParams, code.Param{
					Name: fmt.Sprintf("arg%d", i),
					Type: param.Type,
				})
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
			if !reflect.DeepEqual(testCase.MethodSpec.Returns, actual.Returns) {
				t.Errorf(
					"incorrect struct returns: expected %+v, got %+v",
					testCase.MethodSpec.Returns,
					actual.Returns,
				)
			}
			if err := testutils.ExpectMultiLineString(testCase.ExpectedBody, actual.Body.Code()); err != nil {
				t.Error(err)
			}
		})
	}
}
