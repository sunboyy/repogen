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

type StubUpdate struct {
}

func (update StubUpdate) Name() string {
	return "Stub"
}

func (update StubUpdate) NumberOfArguments() int {
	return 1
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
		{
			Name:           "UpdateTypeNotSupportedError",
			Error:          mongo.NewUpdateTypeNotSupportedError(StubUpdate{}),
			ExpectedString: "update type Stub not supported",
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
