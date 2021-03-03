package spec

import (
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
	FieldReference FieldReference
	Comparator     Comparator
	ParamIndex     int
}

type queryParser struct {
	fieldResolver fieldResolver
	StructModel   code.Struct
}

func (p queryParser) parseQuery(rawTokens []string, paramIndex int) (QuerySpec, error) {
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

	if len(tokens) == 0 {
		return QuerySpec{}, NewInvalidQueryError(rawTokens)
	}

	operator, predicateTokens, err := p.splitPredicateTokens(tokens)
	if err != nil {
		return QuerySpec{}, err
	}

	querySpec := QuerySpec{
		Operator: operator,
	}

	for _, predicateToken := range predicateTokens {
		predicate, err := p.parsePredicate(predicateToken, paramIndex)
		if err != nil {
			return QuerySpec{}, err
		}
		querySpec.Predicates = append(querySpec.Predicates, predicate)
		paramIndex += predicate.Comparator.NumberOfArguments()
	}

	return querySpec, nil
}

func (p queryParser) splitPredicateTokens(tokens []string) (Operator, [][]string, error) {
	var operator Operator
	var predicateTokens [][]string
	var aggregatedToken []string

	for _, token := range tokens {
		if token != "And" && token != "Or" {
			aggregatedToken = append(aggregatedToken, token)
		} else if len(aggregatedToken) == 0 {
			return "", nil, NewInvalidQueryError(tokens)
		} else if token == "And" && operator != OperatorOr {
			operator = OperatorAnd
			predicateTokens = append(predicateTokens, aggregatedToken)
			aggregatedToken = nil
		} else if token == "Or" && operator != OperatorAnd {
			operator = OperatorOr
			predicateTokens = append(predicateTokens, aggregatedToken)
			aggregatedToken = nil
		} else {
			return "", nil, NewInvalidQueryError(tokens)
		}
	}
	if len(aggregatedToken) == 0 {
		return "", nil, NewInvalidQueryError(tokens)
	}
	predicateTokens = append(predicateTokens, aggregatedToken)

	return operator, predicateTokens, nil
}

func (p queryParser) parsePredicate(t []string, paramIndex int) (Predicate, error) {
	if len(t) > 1 && t[len(t)-1] == "Not" {
		return p.createPredicate(t[:len(t)-1], ComparatorNot, paramIndex)
	}
	if len(t) > 2 && t[len(t)-2] == "Less" && t[len(t)-1] == "Than" {
		return p.createPredicate(t[:len(t)-2], ComparatorLessThan, paramIndex)
	}
	if len(t) > 3 && t[len(t)-3] == "Less" && t[len(t)-2] == "Than" && t[len(t)-1] == "Equal" {
		return p.createPredicate(t[:len(t)-3], ComparatorLessThanEqual, paramIndex)
	}
	if len(t) > 2 && t[len(t)-2] == "Greater" && t[len(t)-1] == "Than" {
		return p.createPredicate(t[:len(t)-2], ComparatorGreaterThan, paramIndex)
	}
	if len(t) > 3 && t[len(t)-3] == "Greater" && t[len(t)-2] == "Than" && t[len(t)-1] == "Equal" {
		return p.createPredicate(t[:len(t)-3], ComparatorGreaterThanEqual, paramIndex)
	}
	if len(t) > 2 && t[len(t)-2] == "Not" && t[len(t)-1] == "In" {
		return p.createPredicate(t[:len(t)-2], ComparatorNotIn, paramIndex)
	}
	if len(t) > 1 && t[len(t)-1] == "In" {
		return p.createPredicate(t[:len(t)-1], ComparatorIn, paramIndex)
	}
	if len(t) > 1 && t[len(t)-1] == "Between" {
		return p.createPredicate(t[:len(t)-1], ComparatorBetween, paramIndex)
	}
	if len(t) > 1 && t[len(t)-1] == "True" {
		return p.createPredicate(t[:len(t)-1], ComparatorTrue, paramIndex)
	}
	if len(t) > 1 && t[len(t)-1] == "False" {
		return p.createPredicate(t[:len(t)-1], ComparatorFalse, paramIndex)
	}
	return p.createPredicate(t, ComparatorEqual, paramIndex)
}

func (p queryParser) createPredicate(t []string, comparator Comparator, paramIndex int) (Predicate, error) {
	fields, ok := p.fieldResolver.ResolveStructField(p.StructModel, t)
	if !ok {
		return Predicate{}, NewStructFieldNotFoundError(t)
	}

	return Predicate{
		FieldReference: fields,
		Comparator:     comparator,
		ParamIndex:     paramIndex,
	}, nil
}
