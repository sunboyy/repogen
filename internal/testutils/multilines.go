package testutils

import (
	"fmt"
	"strings"
)

// ExpectMultiLineString compares two multi-line strings and report the
// difference.
func ExpectMultiLineString(expected, actual string) error {
	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(actual, "\n")

	numberOfComparableLines := len(expectedLines)
	if len(actualLines) < numberOfComparableLines {
		numberOfComparableLines = len(actualLines)
	}

	for i := 0; i < numberOfComparableLines; i++ {
		if expectedLines[i] != actualLines[i] {
			return fmt.Errorf("at line %d\nexpected: %v\nreceived: %v", i+1, expectedLines[i], actualLines[i])
		}
	}

	if len(expectedLines) < len(actualLines) {
		return fmt.Errorf("unexpected lines:\n%s", strings.Join(actualLines[len(expectedLines):], "\n"))
	} else if len(expectedLines) > len(actualLines) {
		return fmt.Errorf("missing lines:\n%s", strings.Join(expectedLines[len(actualLines):], "\n"))
	}

	return nil
}
