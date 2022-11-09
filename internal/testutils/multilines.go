package testutils

import (
	"fmt"
	"strings"
)

// lineMismatchedError is an error indicating unmatched line when comparing
// string using ExpectMultiLineString.
type lineMismatchedError struct {
	LineNumber int
	Expected   string
	Received   string
}

// NewLineMismatchedError is a constructor for lineMismatchedError.
func NewLineMismatchedError(lineNumber int, expected, received string) error {
	return lineMismatchedError{
		LineNumber: lineNumber,
		Expected:   expected,
		Received:   received,
	}
}

func (err lineMismatchedError) Error() string {
	return fmt.Sprintf("at line %d\nexpected: %v\nreceived: %v", err.LineNumber,
		err.Expected, err.Received)
}

type missingLinesError struct {
	MissingLines []string
}

// NewMissingLinesError is a constructor for missingLinesError.
func NewMissingLinesError(missingLines []string) error {
	return missingLinesError{
		MissingLines: missingLines,
	}
}

func (err missingLinesError) Error() string {
	return fmt.Sprintf("missing lines:\n%s",
		strings.Join(err.MissingLines, "\n"))
}

type unexpectedLinesError struct {
	UnexpectedLines []string
}

// NewUnexpectedLinesError is a constructor for unexpectedLinesError.
func NewUnexpectedLinesError(unexpectedLines []string) error {
	return unexpectedLinesError{
		UnexpectedLines: unexpectedLines,
	}
}

func (err unexpectedLinesError) Error() string {
	return fmt.Sprintf("unexpected lines:\n%s",
		strings.Join(err.UnexpectedLines, "\n"))
}

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
			return NewLineMismatchedError(i+1, expectedLines[i], actualLines[i])
		}
	}

	if len(expectedLines) < len(actualLines) {
		return NewUnexpectedLinesError(actualLines[len(expectedLines):])
	} else if len(expectedLines) > len(actualLines) {
		return NewMissingLinesError(expectedLines[len(actualLines):])
	}

	return nil
}
