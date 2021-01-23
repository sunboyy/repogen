package testutils_test

import (
	"testing"

	"github.com/sunboyy/repogen/internal/testutils"
)

func TestExpectMultiLineString(t *testing.T) {
	t.Run("same string should return nil", func(t *testing.T) {
		text := `  Hello world
			this is a test text  `

		err := testutils.ExpectMultiLineString(text, text)

		if err != nil {
			t.Errorf("Expected = <nil>\nReceived = %s", err.Error())
		}
	})

	t.Run("different string with same number of lines", func(t *testing.T) {
		expectedText := `  Hello world
this is an expected text
how are you?`
		actualText := `  Hello world
this is a real text
How are you?`

		err := testutils.ExpectMultiLineString(expectedText, actualText)

		expectedError := "On line 2\nExpected: this is an expected text\nReceived: this is a real text"
		if err == nil || err.Error() != expectedError {
			t.Errorf("Expected = %s\nReceived = %s", expectedError, err.Error())
		}
	})

	t.Run("expected text longer than actual text", func(t *testing.T) {
		expectedText := `  Hello world
this is an expected text
how are you?
I'm fine...
Thank you...`
		actualText := `  Hello world
this is an expected text
how are you?`

		err := testutils.ExpectMultiLineString(expectedText, actualText)

		expectedError := "Missing lines:\nI'm fine...\nThank you..."
		if err == nil || err.Error() != expectedError {
			t.Errorf("Expected = %s\nReceived = %s", expectedError, err.Error())
		}
	})

	t.Run("actual text longer than expected text", func(t *testing.T) {
		expectedText := `  Hello world
this is an expected text
how are you?`
		actualText := `  Hello world
this is an expected text
how are you?
I'm fine...
Thank you...`

		err := testutils.ExpectMultiLineString(expectedText, actualText)

		expectedError := "Unexpected lines:\nI'm fine...\nThank you..."
		if err == nil || err.Error() != expectedError {
			t.Errorf("Expected = %s\nReceived = %s", expectedError, err.Error())
		}
	})
}
