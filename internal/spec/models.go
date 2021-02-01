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

// InsertOperation is a method specification for insert operations
type InsertOperation struct {
	Mode QueryMode
}

// FindOperation is a method specification for find operations
type FindOperation struct {
	Mode  QueryMode
	Query QuerySpec
}

// UpdateOperation is a method specification for update operations
type UpdateOperation struct {
	Fields []UpdateField
	Mode   QueryMode
	Query  QuerySpec
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

// QuerySpec is a set of conditions of querying the database
type QuerySpec struct {
	Operator   Operator
	Predicates []Predicate
}

// NumberOfArguments returns number of arguments required to perform the query
func (q QuerySpec) NumberOfArguments() int {
	var totalArgs int
	for _, predicate := range q.Predicates {
		totalArgs += predicate.Comparator.NumberOfArguments()
	}
	return totalArgs
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
	ComparatorBetween          Comparator = "BETWEEN"
	ComparatorIn               Comparator = "IN"
)

// ArgumentTypeFromFieldType returns a type of required argument from the given struct field type
func (c Comparator) ArgumentTypeFromFieldType(t code.Type) code.Type {
	if c == ComparatorIn {
		return code.ArrayType{ContainedType: t}
	}
	return t
}

// NumberOfArguments returns the number of arguments required to perform the comparison
func (c Comparator) NumberOfArguments() int {
	if c == ComparatorBetween {
		return 2
	}
	return 1
}

// Predicate is a criteria for querying a field
type Predicate struct {
	Field      string
	Comparator Comparator
	ParamIndex int
}

type predicateToken []string

func (t predicateToken) ToPredicate(paramIndex int) Predicate {
	if len(t) > 1 && t[len(t)-1] == "Not" {
		return Predicate{Field: strings.Join(t[:len(t)-1], ""), Comparator: ComparatorNot, ParamIndex: paramIndex}
	}
	if len(t) > 2 && t[len(t)-2] == "Less" && t[len(t)-1] == "Than" {
		return Predicate{Field: strings.Join(t[:len(t)-2], ""), Comparator: ComparatorLessThan, ParamIndex: paramIndex}
	}
	if len(t) > 3 && t[len(t)-3] == "Less" && t[len(t)-2] == "Than" && t[len(t)-1] == "Equal" {
		return Predicate{Field: strings.Join(t[:len(t)-3], ""), Comparator: ComparatorLessThanEqual, ParamIndex: paramIndex}
	}
	if len(t) > 2 && t[len(t)-2] == "Greater" && t[len(t)-1] == "Than" {
		return Predicate{Field: strings.Join(t[:len(t)-2], ""), Comparator: ComparatorGreaterThan, ParamIndex: paramIndex}
	}
	if len(t) > 3 && t[len(t)-3] == "Greater" && t[len(t)-2] == "Than" && t[len(t)-1] == "Equal" {
		return Predicate{Field: strings.Join(t[:len(t)-3], ""), Comparator: ComparatorGreaterThanEqual, ParamIndex: paramIndex}
	}
	if len(t) > 1 && t[len(t)-1] == "In" {
		return Predicate{Field: strings.Join(t[:len(t)-1], ""), Comparator: ComparatorIn, ParamIndex: paramIndex}
	}
	if len(t) > 1 && t[len(t)-1] == "Between" {
		return Predicate{Field: strings.Join(t[:len(t)-1], ""), Comparator: ComparatorBetween, ParamIndex: paramIndex}
	}
	return Predicate{Field: strings.Join(t, ""), Comparator: ComparatorEqual, ParamIndex: paramIndex}
}
