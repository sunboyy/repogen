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
	Operator   Operator
	Predicates []Predicate
}

// NumberOfArguments returns number of arguments required to perform the query
func (q QuerySpec) NumberOfArguments() int {
	return len(q.Predicates)
}

// Operator is a boolean operator for merging conditions
type Operator string

// boolean operator types
const (
	OperatorAnd Operator = "AND"
	OperatorOr  Operator = "OR"
)

// Comparator is a comparison operator of the condition to query the data
type Comparator string

// comparator types
const (
	ComparatorNot              Comparator = "NOT"
	ComparatorEqual            Comparator = "EQUAL"
	ComparatorLessThan         Comparator = "LESS_THAN"
	ComparatorLessThanEqual    Comparator = "LESS_THAN_EQUAL"
	ComparatorGreaterThan      Comparator = "GREATER_THAN"
	ComparatorGreaterThanEqual Comparator = "GREATER_THAN_EQUAL"
	ComparatorIn               Comparator = "IN"
)

// ArgumentTypeFromFieldType returns a type of required argument from the given struct field type
func (c Comparator) ArgumentTypeFromFieldType(t code.Type) code.Type {
	if c == ComparatorIn {
		return code.ArrayType{ContainedType: t}
	}
	return t
}

// Predicate is a criteria for querying a field
type Predicate struct {
	Field      string
	Comparator Comparator
}

type predicateToken []string

func (t predicateToken) ToPredicate() Predicate {
	if len(t) > 1 && t[len(t)-1] == "Not" {
		return Predicate{Field: strings.Join(t[:len(t)-1], ""), Comparator: ComparatorNot}
	}
	if len(t) > 2 && t[len(t)-2] == "Less" && t[len(t)-1] == "Than" {
		return Predicate{Field: strings.Join(t[:len(t)-2], ""), Comparator: ComparatorLessThan}
	}
	if len(t) > 3 && t[len(t)-3] == "Less" && t[len(t)-2] == "Than" && t[len(t)-1] == "Equal" {
		return Predicate{Field: strings.Join(t[:len(t)-3], ""), Comparator: ComparatorLessThanEqual}
	}
	if len(t) > 2 && t[len(t)-2] == "Greater" && t[len(t)-1] == "Than" {
		return Predicate{Field: strings.Join(t[:len(t)-2], ""), Comparator: ComparatorGreaterThan}
	}
	if len(t) > 3 && t[len(t)-3] == "Greater" && t[len(t)-2] == "Than" && t[len(t)-1] == "Equal" {
		return Predicate{Field: strings.Join(t[:len(t)-3], ""), Comparator: ComparatorGreaterThanEqual}
	}
	if len(t) > 1 && t[len(t)-1] == "In" {
		return Predicate{Field: strings.Join(t[:len(t)-1], ""), Comparator: ComparatorIn}
	}
	return Predicate{Field: strings.Join(t, ""), Comparator: ComparatorEqual}
}
