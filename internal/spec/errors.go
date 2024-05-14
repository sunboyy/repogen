package spec

import (
	"errors"
	"fmt"
	"go/types"
	"strings"

	"github.com/sunboyy/repogen/internal/code"
)

// parsing error constants
var (
	ErrQueryRequired        = errors.New("spec: query is required")
	ErrInvalidParam         = errors.New("spec: parameters do not match the query")
	ErrInvalidUpdateFields  = errors.New("spec: update fields are invalid")
	ErrContextParamRequired = errors.New("spec: context parameter is required")
	ErrLimitAmountRequired  = errors.New("spec: limit amount is required")
	ErrLimitNonPositive     = errors.New("spec: limit value must be positive")
	ErrLimitOnFindOne       = errors.New("spec: cannot specify limit on find one")
)

// NewUnsupportedReturnError creates unsupportedReturnError
func NewUnsupportedReturnError(givenType types.Type, index int) error {
	return unsupportedReturnError{
		GivenType: givenType,
		Index:     index,
	}
}

type unsupportedReturnError struct {
	GivenType types.Type
	Index     int
}

func (err unsupportedReturnError) Error() string {
	return fmt.Sprintf("return type '%s' at index %d is not supported", err.GivenType.String(), err.Index)
}

// NewOperationReturnCountUnmatchedError creates
// operationReturnCountUnmatchedError.
func NewOperationReturnCountUnmatchedError(returnCount int) error {
	return operationReturnCountUnmatchedError{
		ReturnCount: returnCount,
	}
}

type operationReturnCountUnmatchedError struct {
	ReturnCount int
}

func (err operationReturnCountUnmatchedError) Error() string {
	return fmt.Sprintf("operation requires return count of %d", err.ReturnCount)
}

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
func NewArgumentTypeNotMatchedError(fieldName string, requiredType types.Type, givenType types.Type) error {
	return argumentTypeNotMatchedError{
		FieldName:    fieldName,
		RequiredType: requiredType,
		GivenType:    givenType,
	}
}

type argumentTypeNotMatchedError struct {
	FieldName    string
	RequiredType types.Type
	GivenType    types.Type
}

func (err argumentTypeNotMatchedError) Error() string {
	return fmt.Sprintf("field '%s' requires an argument of type '%s' (got '%s')",
		err.FieldName, err.RequiredType.String(), err.GivenType.String())
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
		err.Comparator, err.Field.Var.Name(), err.Field.Var.Type())
}

// NewIncompatibleUpdateOperatorError creates incompatibleUpdateOperatorError
func NewIncompatibleUpdateOperatorError(updateOperator UpdateOperator, fieldReference FieldReference) error {
	return incompatibleUpdateOperatorError{
		UpdateOperator:  updateOperator,
		ReferencingCode: fieldReference.ReferencingCode(),
		ReferencedType:  fieldReference.ReferencedField().Var.Type(),
	}
}

type incompatibleUpdateOperatorError struct {
	UpdateOperator  UpdateOperator
	ReferencingCode string
	ReferencedType  types.Type
}

func (err incompatibleUpdateOperatorError) Error() string {
	return fmt.Sprintf("cannot use update operator %s with struct field '%s' of type '%s'",
		err.UpdateOperator, err.ReferencingCode, err.ReferencedType.String())
}
