package spec_test

import (
	"testing"

	"github.com/sunboyy/repogen/internal/code"
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
			Error:          spec.NewStructFieldNotFoundError([]string{"Phone", "Number"}),
			ExpectedString: "struct field 'PhoneNumber' not found",
		},
		{
			Name:           "InvalidQueryError",
			Error:          spec.NewInvalidQueryError([]string{"And"}),
			ExpectedString: "invalid query 'And'",
		},
		{
			Name: "IncompatibleComparatorError",
			Error: spec.NewIncompatibleComparatorError(spec.ComparatorTrue, code.StructField{
				Name: "Age",
				Type: code.SimpleType("int"),
			}),
			ExpectedString: "cannot use comparator EQUAL_TRUE with struct field 'Age' of type 'int'",
		},
		{
			Name:           "InvalidSortError",
			Error:          spec.NewInvalidSortError([]string{"Order", "By"}),
			ExpectedString: "invalid sort 'OrderBy'",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.Error.Error() != testCase.ExpectedString {
				t.Errorf("Expected = %+v\nReceived = %+v", testCase.ExpectedString, testCase.Error.Error())
			}
		})
	}
}
