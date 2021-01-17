package spec

import "github.com/sunboyy/repogen/internal/code"

// QueryMode one or many
type QueryMode string

// query mode constants
const (
	QueryModeOne  QueryMode = "ONE"
	QueryModeMany QueryMode = "MANY"
)

// RepositorySpec is a specification generated from the repository interface
type RepositorySpec struct {
	InterfaceName string
	Methods       []MethodSpec
}

// MethodSpec is a method specification inside repository specification
type MethodSpec struct {
	Name      string
	Params    []code.Param
	Returns   []code.Type
	Operation Operation
}

// Operation is an interface for any kind of operation
type Operation interface {
}

// FindOperation is a method specification for find operations
type FindOperation struct {
	Mode  QueryMode
	Query QuerySpec
}

// QuerySpec is a condition of querying the database
type QuerySpec struct {
	Fields []string
}

// NumberOfArguments returns number of arguments required to perform the query
func (q QuerySpec) NumberOfArguments() int {
	return len(q.Fields)
}
