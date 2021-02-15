package mongo

import (
	"fmt"
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
