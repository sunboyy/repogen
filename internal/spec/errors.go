package spec

import (
	"fmt"
	"strings"

	"github.com/sunboyy/repogen/internal/code"
)

// ParsingError is an error from parsing interface methods
type ParsingError string

func (err ParsingError) Error() string {
	switch err {
	case UnsupportedReturnError:
		return "this type of return is not supported"
	case QueryRequiredError:
		return "query is required"
	case InvalidParamError:
		return "parameters do not match the query"
	case InvalidUpdateFieldsError:
		return "update fields is invalid"
	case ContextParamRequiredError:
		return "context parameter is required"
	}
	return string(err)
}

// parsing error constants
const (
	UnsupportedReturnError    ParsingError = "ERROR_UNSUPPORTED_RETURN"
	QueryRequiredError        ParsingError = "ERROR_QUERY_REQUIRED"
	InvalidParamError         ParsingError = "ERROR_INVALID_PARAM"
	InvalidUpdateFieldsError  ParsingError = "ERROR_INVALID_UPDATE_FIELDS"
	ContextParamRequiredError ParsingError = "ERROR_CONTEXT_PARAM_REQUIRED"
)

// NewInvalidQueryError creates invalidQueryError
func NewInvalidQueryError(queryTokens []string) error {
	return invalidQueryError{QueryString: strings.Join(queryTokens, "")}
}

type invalidQueryError struct {
	QueryString string
}

func (err invalidQueryError) Error() string {
	return fmt.Sprintf("invalid query '%s'", err.QueryString)
}

// NewUnknownOperationError creates unknownOperationError
func NewUnknownOperationError(operationName string) error {
	return unknownOperationError{OperationName: operationName}
}

type unknownOperationError struct {
	OperationName string
}

func (err unknownOperationError) Error() string {
	return fmt.Sprintf("unknown operation '%s'", err.OperationName)
}

// NewStructFieldNotFoundError creates structFieldNotFoundError
func NewStructFieldNotFoundError(fieldName string) error {
	return structFieldNotFoundError{FieldName: fieldName}
}

type structFieldNotFoundError struct {
	FieldName string
}

func (err structFieldNotFoundError) Error() string {
	return fmt.Sprintf("struct field '%s' not found", err.FieldName)
}

// NewIncompatibleComparatorError creates incompatibleComparatorError
func NewIncompatibleComparatorError(comparator Comparator, field code.StructField) error {
	return incompatibleComparatorError{
		Comparator: comparator,
		Field:      field,
	}
}

type incompatibleComparatorError struct {
	Comparator Comparator
	Field      code.StructField
}

func (err incompatibleComparatorError) Error() string {
	return fmt.Sprintf("cannot use comparator %s with struct field '%s' of type '%s'",
		err.Comparator, err.Field.Name, err.Field.Type.Code())
}
