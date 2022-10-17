package code_test

import (
	"testing"

	"github.com/sunboyy/repogen/internal/code"
)

type ErrorTestCase struct {
	Name           string
	Error          error
	ExpectedString string
}

func TestError(t *testing.T) {
	testTable := []ErrorTestCase{
		{
			Name:           "DuplicateStructError",
			Error:          code.DuplicateStructError("User"),
			ExpectedString: "code: duplicate implementation of struct 'User'",
		},
		{
			Name:           "DuplicateInterfaceError",
			Error:          code.DuplicateInterfaceError("UserRepository"),
			ExpectedString: "code: duplicate implementation of interface 'UserRepository'",
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
