package mongo

import (
	"fmt"
	"strings"

	"github.com/sunboyy/repogen/internal/spec"
)

type querySpec struct {
	Operator   spec.Operator
	Predicates []predicate
}

func (q querySpec) Code() string {
	var predicateCodes []string
	for i, predicate := range q.Predicates {
		predicateCodes = append(predicateCodes, predicate.Code(i))
	}

	var lines []string
	switch q.Operator {
	case spec.OperatorOr:
		lines = append(lines, `		"$or": []bson.M{`)
		for _, predicateCode := range predicateCodes {
			lines = append(lines, fmt.Sprintf(`			{%s},`, predicateCode))
		}
		lines = append(lines, `		},`)
	default:
		for _, predicateCode := range predicateCodes {
			lines = append(lines, fmt.Sprintf(`		%s,`, predicateCode))
		}
	}
	return strings.Join(lines, "\n")
}

type predicate struct {
	Field      string
	Comparator spec.Comparator
}

func (p predicate) Code(argIndex int) string {
	switch p.Comparator {
	case spec.ComparatorEqual:
		return fmt.Sprintf(`"%s": arg%d`, p.Field, argIndex)
	case spec.ComparatorNot:
		return fmt.Sprintf(`"%s": bson.M{"$ne": arg%d}`, p.Field, argIndex)
	case spec.ComparatorLessThan:
		return fmt.Sprintf(`"%s": bson.M{"$lt": arg%d}`, p.Field, argIndex)
	case spec.ComparatorLessThanEqual:
		return fmt.Sprintf(`"%s": bson.M{"$lte": arg%d}`, p.Field, argIndex)
	case spec.ComparatorGreaterThan:
		return fmt.Sprintf(`"%s": bson.M{"$gt": arg%d}`, p.Field, argIndex)
	case spec.ComparatorGreaterThanEqual:
		return fmt.Sprintf(`"%s": bson.M{"$gte": arg%d}`, p.Field, argIndex)
	}
	return ""
}
