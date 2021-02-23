package spec

import (
	"strings"

	"github.com/sunboyy/repogen/internal/code"
)

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
	ComparatorNotIn            Comparator = "NOT_IN"
	ComparatorTrue             Comparator = "EQUAL_TRUE"
	ComparatorFalse            Comparator = "EQUAL_FALSE"
)

// ArgumentTypeFromFieldType returns a type of required argument from the given struct field type
func (c Comparator) ArgumentTypeFromFieldType(t code.Type) code.Type {
	switch c {
	case ComparatorIn, ComparatorNotIn:
		return code.ArrayType{ContainedType: t}
	default:
		return t
	}
}

// NumberOfArguments returns the number of arguments required to perform the comparison
func (c Comparator) NumberOfArguments() int {
	switch c {
	case ComparatorBetween:
		return 2
	case ComparatorTrue, ComparatorFalse:
		return 0
	default:
		return 1
	}
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
	if len(t) > 2 && t[len(t)-2] == "Not" && t[len(t)-1] == "In" {
		return Predicate{Field: strings.Join(t[:len(t)-2], ""), Comparator: ComparatorNotIn, ParamIndex: paramIndex}
	}
	if len(t) > 1 && t[len(t)-1] == "In" {
		return Predicate{Field: strings.Join(t[:len(t)-1], ""), Comparator: ComparatorIn, ParamIndex: paramIndex}
	}
	if len(t) > 1 && t[len(t)-1] == "Between" {
		return Predicate{Field: strings.Join(t[:len(t)-1], ""), Comparator: ComparatorBetween, ParamIndex: paramIndex}
	}
	if len(t) > 1 && t[len(t)-1] == "True" {
		return Predicate{Field: strings.Join(t[:len(t)-1], ""), Comparator: ComparatorTrue, ParamIndex: paramIndex}
	}
	if len(t) > 1 && t[len(t)-1] == "False" {
		return Predicate{Field: strings.Join(t[:len(t)-1], ""), Comparator: ComparatorFalse, ParamIndex: paramIndex}
	}
	return Predicate{Field: strings.Join(t, ""), Comparator: ComparatorEqual, ParamIndex: paramIndex}
}

func parseQuery(rawTokens []string, paramIndex int) (QuerySpec, error) {
	if len(rawTokens) == 0 {
		return QuerySpec{}, QueryRequiredError
	}

	tokens := rawTokens
	if len(tokens) == 1 && tokens[0] == "All" {
		return QuerySpec{}, nil
	}

	if tokens[0] == "One" {
		tokens = tokens[1:]
	}
	if tokens[0] == "By" {
		tokens = tokens[1:]
	}

	if len(tokens) == 0 || tokens[0] == "And" || tokens[0] == "Or" {
		return QuerySpec{}, NewInvalidQueryError(rawTokens)
	}

	var operator Operator
	var predicates []Predicate
	var aggregatedToken predicateToken
	for _, token := range tokens {
		if token != "And" && token != "Or" {
			aggregatedToken = append(aggregatedToken, token)
		} else if len(aggregatedToken) == 0 {
			return QuerySpec{}, NewInvalidQueryError(rawTokens)
		} else if token == "And" && operator != OperatorOr {
			operator = OperatorAnd
			predicate := aggregatedToken.ToPredicate(paramIndex)
			predicates = append(predicates, predicate)
			paramIndex += predicate.Comparator.NumberOfArguments()
			aggregatedToken = predicateToken{}
		} else if token == "Or" && operator != OperatorAnd {
			operator = OperatorOr
			predicate := aggregatedToken.ToPredicate(paramIndex)
			predicates = append(predicates, predicate)
			paramIndex += predicate.Comparator.NumberOfArguments()
			aggregatedToken = predicateToken{}
		} else {
			return QuerySpec{}, NewInvalidQueryError(rawTokens)
		}
	}
	if len(aggregatedToken) == 0 {
		return QuerySpec{}, NewInvalidQueryError(rawTokens)
	}
	predicates = append(predicates, aggregatedToken.ToPredicate(paramIndex))

	return QuerySpec{Operator: operator, Predicates: predicates}, nil
}
