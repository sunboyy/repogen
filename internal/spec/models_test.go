package spec_test

import (
	"testing"

	"github.com/sunboyy/repogen/internal/spec"
)

type OperationTestCase struct {
	Operation    spec.Operation
	ExpectedName string
}

func TestOperationName(t *testing.T) {
	testTable := []OperationTestCase{
		{
			Operation:    spec.InsertOperation{},
			ExpectedName: "Insert",
		},
		{
			Operation:    spec.FindOperation{},
			ExpectedName: "Find",
		},
		{
			Operation:    spec.UpdateOperation{},
			ExpectedName: "Update",
		},
		{
			Operation:    spec.DeleteOperation{},
			ExpectedName: "Delete",
		},
		{
			Operation:    spec.CountOperation{},
			ExpectedName: "Count",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.ExpectedName, func(t *testing.T) {
			if testCase.Operation.Name() != testCase.ExpectedName {
				t.Errorf("Expected = %v\nReceived = %v", testCase.ExpectedName, testCase.Operation.Name())
			}
		})
	}
}

type UpdateTypeTestCase struct {
	Update       spec.Update
	ExpectedName string
}

func TestUpdateTypeName(t *testing.T) {
	testTable := []UpdateTypeTestCase{
		{
			Update:       spec.UpdateModel{},
			ExpectedName: "Model",
		},
		{
			Update:       spec.UpdateFields{},
			ExpectedName: "Fields",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.ExpectedName, func(t *testing.T) {
			if testCase.Update.Name() != testCase.ExpectedName {
				t.Errorf("Expected = %v\nReceived = %v", testCase.ExpectedName, testCase.Update.Name())
			}
		})
	}
}
