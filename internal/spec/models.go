package spec

import (
	"strings"

	"github.com/sunboyy/repogen/internal/code"
)

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

// QuerySpec is a set of conditions of querying the database
type QuerySpec struct {
	Predicates []Predicate
}

// NumberOfArguments returns number of arguments required to perform the query
func (q QuerySpec) NumberOfArguments() int {
	return len(q.Predicates)
}

// Operator is an operator of the condition to query the data
type Operator string

// operator constants
const (
	OperatorEqual            Operator = "EQUAL"
	OperatorNot              Operator = "NOT"
	OperatorLessThan         Operator = "LESS_THAN"
	OperatorLessThanEqual    Operator = "LESS_THAN_EQUAL"
	OperatorGreaterThan      Operator = "GREATER_THAN"
	OperatorGreaterThanEqual Operator = "GREATER_THAN_EQUAL"
)

// Predicate is a criteria for querying a field
type Predicate struct {
	Field    string
	Operator Operator
}

type predicateToken []string

func (t predicateToken) ToPredicate() Predicate {
	if len(t) > 1 && t[len(t)-1] == "Not" {
		return Predicate{Field: strings.Join(t[:len(t)-1], ""), Operator: OperatorNot}
	}
	if len(t) > 2 && t[len(t)-2] == "Less" && t[len(t)-1] == "Than" {
		return Predicate{Field: strings.Join(t[:len(t)-2], ""), Operator: OperatorLessThan}
	}
	if len(t) > 3 && t[len(t)-3] == "Less" && t[len(t)-2] == "Than" && t[len(t)-1] == "Equal" {
		return Predicate{Field: strings.Join(t[:len(t)-3], ""), Operator: OperatorLessThanEqual}
	}
	if len(t) > 2 && t[len(t)-2] == "Greater" && t[len(t)-1] == "Than" {
		return Predicate{Field: strings.Join(t[:len(t)-2], ""), Operator: OperatorGreaterThan}
	}
	if len(t) > 3 && t[len(t)-3] == "Greater" && t[len(t)-2] == "Than" && t[len(t)-1] == "Equal" {
		return Predicate{Field: strings.Join(t[:len(t)-3], ""), Operator: OperatorGreaterThanEqual}
	}
	return Predicate{Field: strings.Join(t, ""), Operator: OperatorEqual}
}
