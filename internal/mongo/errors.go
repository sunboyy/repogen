package mongo

import (
	"fmt"

	"github.com/sunboyy/repogen/internal/spec"
)

// NewOperationNotSupportedError creates operationNotSupportedError
func NewOperationNotSupportedError(operationName string) error {
	return operationNotSupportedError{OperationName: operationName}
}

type operationNotSupportedError struct {
	OperationName string
}

func (err operationNotSupportedError) Error() string {
	return fmt.Sprintf("operation '%s' not supported", err.OperationName)
}

// NewBsonTagNotFoundError creates bsonTagNotFoundError
func NewBsonTagNotFoundError(fieldName string) error {
	return bsonTagNotFoundError{FieldName: fieldName}
}

type bsonTagNotFoundError struct {
	FieldName string
}

func (err bsonTagNotFoundError) Error() string {
	return fmt.Sprintf("bson tag of field '%s' not found", err.FieldName)
}

// NewUpdateTypeNotSupportedError creates updateTypeNotSupportedError
func NewUpdateTypeNotSupportedError(update spec.Update) error {
	return updateTypeNotSupportedError{Update: update}
}

type updateTypeNotSupportedError struct {
	Update spec.Update
}

func (err updateTypeNotSupportedError) Error() string {
	return fmt.Sprintf("update type %s not supported", err.Update.Name())
}

// NewUpdateOperatorNotSupportedError creates updateOperatorNotSupportedError
func NewUpdateOperatorNotSupportedError(operator spec.UpdateOperator) error {
	return updateOperatorNotSupportedError{Operator: operator}
}

type updateOperatorNotSupportedError struct {
	Operator spec.UpdateOperator
}

func (err updateOperatorNotSupportedError) Error() string {
	return fmt.Sprintf("update operator %s not supported", err.Operator)
}
