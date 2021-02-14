package spec_test

import (
	"testing"

	"github.com/sunboyy/repogen/internal/spec"
)

type ErrorTestCase struct {
	Name           string
	Error          error
	ExpectedString string
}

func TestError(t *testing.T) {
	testTable := []ErrorTestCase{
		{
			Name:           "UnknownOperationError",
			Error:          spec.NewUnknownOperationError("Search"),
			ExpectedString: "unknown operation 'Search'",
		},
		{
			Name:           "StructFieldNotFoundError",
			Error:          spec.NewStructFieldNotFoundError("Country"),
			ExpectedString: "struct field 'Country' not found",
		},
		{
			Name:           "InvalidQueryError",
			Error:          spec.NewInvalidQueryError([]string{"By", "And"}),
			ExpectedString: "invalid query 'ByAnd'",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.Error.Error() != testCase.ExpectedString {
				t.Errorf("Expected = %v\nReceived = %v", testCase.ExpectedString, testCase.Error.Error())
			}
		})
	}
}
