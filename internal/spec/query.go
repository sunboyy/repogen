package spec

import "go/types"

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
	ComparatorExists           Comparator = "EXISTS"
	ComparatorNotExists        Comparator = "NOT_EXISTS"
)

// ArgumentTypeFromFieldType returns a type of required argument from the given
// struct field type.
func (c Comparator) ArgumentTypeFromFieldType(t types.Type) types.Type {
	switch c {
	case ComparatorIn, ComparatorNotIn:
		return types.NewSlice(t)
	default:
		return t
	}
}

// NumberOfArguments returns the number of arguments required to perform the
// comparison.
func (c Comparator) NumberOfArguments() int {
	switch c {
	case ComparatorBetween:
		return 2
	case ComparatorTrue, ComparatorFalse, ComparatorExists, ComparatorNotExists:
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
	UnderlyingStruct *types.Struct
}

func (p queryParser) parseQuery(rawTokens []string, paramIndex int) (QuerySpec,
	error) {

	if len(rawTokens) == 0 {
		return QuerySpec{}, ErrQueryRequired
	}

	switch rawTokens[0] {
	case "All":
		if len(rawTokens) == 1 {
			return QuerySpec{}, nil
		}
	case "By":
		return p.parseQueryBy(rawTokens, paramIndex)
	}

	return QuerySpec{}, NewInvalidQueryError(rawTokens)
}

func (p queryParser) parseQueryBy(rawTokens []string, paramIndex int) (QuerySpec, error) {
	tokens := rawTokens[1:]
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

func (p queryParser) parsePredicate(t []string, paramIndex int) (Predicate,
	error) {

	switch {
	case endsWith(t, "Not"):
		return p.createPredicate(t[:len(t)-1], ComparatorNot, paramIndex)

	case endsWith(t, "Less", "Than"):
		return p.createPredicate(t[:len(t)-2], ComparatorLessThan, paramIndex)

	case endsWith(t, "Less", "Than", "Equal"):
		return p.createPredicate(t[:len(t)-3], ComparatorLessThanEqual, paramIndex)

	case endsWith(t, "Greater", "Than"):
		return p.createPredicate(t[:len(t)-2], ComparatorGreaterThan, paramIndex)

	case endsWith(t, "Greater", "Than", "Equal"):
		return p.createPredicate(t[:len(t)-3], ComparatorGreaterThanEqual, paramIndex)

	case endsWith(t, "Not", "In"):
		return p.createPredicate(t[:len(t)-2], ComparatorNotIn, paramIndex)

	case endsWith(t, "Not", "Exists"):
		return p.createPredicate(t[:len(t)-2], ComparatorNotExists, paramIndex)

	case endsWith(t, "In"):
		return p.createPredicate(t[:len(t)-1], ComparatorIn, paramIndex)

	case endsWith(t, "Between"):
		return p.createPredicate(t[:len(t)-1], ComparatorBetween, paramIndex)

	case endsWith(t, "True"):
		return p.createPredicate(t[:len(t)-1], ComparatorTrue, paramIndex)

	case endsWith(t, "False"):
		return p.createPredicate(t[:len(t)-1], ComparatorFalse, paramIndex)

	case endsWith(t, "Exists"):
		return p.createPredicate(t[:len(t)-1], ComparatorExists, paramIndex)
	}

	return p.createPredicate(t, ComparatorEqual, paramIndex)
}

func endsWith(t []string, suffix ...string) bool {
	if len(t) < len(suffix) {
		return false
	}

	for i, suffixToken := range suffix {
		if t[len(t)-len(suffix)+i] != suffixToken {
			return false
		}
	}
	return true
}

func (p queryParser) createPredicate(t []string, comparator Comparator,
	paramIndex int) (Predicate, error) {

	fields, ok := resolveStructField(p.UnderlyingStruct, t)
	if !ok {
		return Predicate{}, NewStructFieldNotFoundError(t)
	}

	return Predicate{
		FieldReference: fields,
		Comparator:     comparator,
		ParamIndex:     paramIndex,
	}, nil
}
