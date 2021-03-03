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
	Sorts []Sort
}

// Name returns "Find" operation name
func (o FindOperation) Name() string {
	return "Find"
}

// Sort is a detail of sorting find result
type Sort struct {
	FieldReference FieldReference
	Ordering       Ordering
}

// Ordering is a sort order
type Ordering string

// Ordering constants
const (
	OrderingAscending  = "ASC"
	OrderingDescending = "DESC"
)

// UpdateOperation is a method specification for update operations
type UpdateOperation struct {
	Update Update
	Mode   QueryMode
	Query  QuerySpec
}

// Update is an interface of update operation type
type Update interface {
	Name() string
	NumberOfArguments() int
}

// UpdateFields is a type of update operation that update specific fields
type UpdateFields []UpdateField

// Name returns UpdateFields name 'Fields'
func (u UpdateFields) Name() string {
	return "Fields"
}

// NumberOfArguments returns number of update fields
func (u UpdateFields) NumberOfArguments() int {
	return len(u)
}

// UpdateModel is a type of update operation that update the whole model
type UpdateModel struct {
}

// Name returns UpdateModel name 'Model'
func (u UpdateModel) Name() string {
	return "Model"
}

// NumberOfArguments returns 1
func (u UpdateModel) NumberOfArguments() int {
	return 1
}

// Name returns "Update" operation name
func (o UpdateOperation) Name() string {
	return "Update"
}

// UpdateField stores mapping between field name in the model and the parameter index
type UpdateField struct {
	FieldReference FieldReference
	ParamIndex     int
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
