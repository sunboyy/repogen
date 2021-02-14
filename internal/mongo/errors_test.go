package mongo_test

import (
	"testing"

	"github.com/sunboyy/repogen/internal/mongo"
)

type ErrorTestCase struct {
	Name           string
	Error          error
	ExpectedString string
}

func TestError(t *testing.T) {
	testTable := []ErrorTestCase{
		{
			Name:           "OperationNotSupportedError",
			Error:          mongo.NewOperationNotSupportedError("Stub"),
			ExpectedString: "operation 'Stub' not supported",
		},
		{
			Name:           "BsonTagNotFoundError",
			Error:          mongo.NewBsonTagNotFoundError("AccessToken"),
			ExpectedString: "bson tag of field 'AccessToken' not found",
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
