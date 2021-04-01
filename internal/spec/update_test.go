package spec_test

import (
	"testing"

	"github.com/sunboyy/repogen/internal/spec"
)

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
				t.Errorf("Expected = %+v\nReceived = %+v", testCase.ExpectedName, testCase.Update.Name())
			}
		})
	}
}
