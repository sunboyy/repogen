package spec

import (
	"github.com/sunboyy/repogen/internal/code"
)

// QueryMode one or many
type QueryMode string

// query mode constants
const (
	QueryModeOne  QueryMode = "ONE"
	QueryModeMany QueryMode = "MANY"
)

// MethodSpec is a method specification inside repository specification
type MethodSpec struct {
	Name      string
	Params    []code.Param
	Returns   []code.Type
	Operation Operation
}

// Operation is an interface for any kind of operation
type Operation interface {
	Name() string
}

// InsertOperation is a method specification for insert operations
type InsertOperation struct {
	Mode QueryMode
}

// Name returns "Insert" operation name
func (o InsertOperation) Name() string {
	return "Insert"
}

// FindOperation is a method specification for find operations
type FindOperation struct {
	Mode  QueryMode
	Query QuerySpec
}

// Name returns "Find" operation name
func (o FindOperation) Name() string {
	return "Find"
}

// UpdateOperation is a method specification for update operations
type UpdateOperation struct {
	Fields []UpdateField
	Mode   QueryMode
	Query  QuerySpec
}

// Name returns "Update" operation name
func (o UpdateOperation) Name() string {
	return "Update"
}

// UpdateField stores mapping between field name in the model and the parameter index
type UpdateField struct {
	Name       string
	ParamIndex int
}

// DeleteOperation is a method specification for delete operations
type DeleteOperation struct {
	Mode  QueryMode
	Query QuerySpec
}

// Name returns "Delete" operation name
func (o DeleteOperation) Name() string {
	return "Delete"
}

// CountOperation is a method specification for count operations
type CountOperation struct {
	Query QuerySpec
}

// Name returns "Count" operation name
func (o CountOperation) Name() string {
	return "Count"
}
