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
		return "update fields are invalid"
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

// NewInvalidSortError creates invalidSortError
func NewInvalidSortError(sortTokens []string) error {
	return invalidSortError{SortString: strings.Join(sortTokens, "")}
}

type invalidSortError struct {
	SortString string
}

func (err invalidSortError) Error() string {
	return fmt.Sprintf("invalid sort '%s'", err.SortString)
}

// NewArgumentTypeNotMatchedError creates argumentTypeNotMatchedError
func NewArgumentTypeNotMatchedError(fieldName string, requiredType code.Type, givenType code.Type) error {
	return argumentTypeNotMatchedError{
		FieldName:    fieldName,
		RequiredType: requiredType,
		GivenType:    givenType,
	}
}

type argumentTypeNotMatchedError struct {
	FieldName    string
	RequiredType code.Type
	GivenType    code.Type
}

func (err argumentTypeNotMatchedError) Error() string {
	return fmt.Sprintf("field '%s' requires an argument of type '%s' (got '%s')",
		err.FieldName, err.RequiredType.Code(), err.GivenType.Code())
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
func NewStructFieldNotFoundError(tokens []string) error {
	return structFieldNotFoundError{FieldName: strings.Join(tokens, "")}
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

// NewIncompatibleUpdateOperatorError creates incompatibleUpdateOperatorError
func NewIncompatibleUpdateOperatorError(updateOperator UpdateOperator, fieldReference FieldReference) error {
	return incompatibleUpdateOperatorError{
		UpdateOperator:  updateOperator,
		ReferencingCode: fieldReference.ReferencingCode(),
		ReferencedType:  fieldReference.ReferencedField().Type,
	}
}

type incompatibleUpdateOperatorError struct {
	UpdateOperator  UpdateOperator
	ReferencingCode string
	ReferencedType  code.Type
}

func (err incompatibleUpdateOperatorError) Error() string {
	return fmt.Sprintf("cannot use update operator %s with struct field '%s' of type '%s'",
		err.UpdateOperator, err.ReferencingCode, err.ReferencedType.Code())
}
