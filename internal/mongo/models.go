package mongo

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sunboyy/repogen/internal/spec"
)

type updateField struct {
	BsonTag    string
	ParamIndex int
}

type update interface {
	Code() string
}

type updateModel struct {
}

func (u updateModel) Code() string {
	return `		"$set": arg1,`
}

type updateFields map[string][]updateField

func (u updateFields) Code() string {
	var keys []string
	for k := range u {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var lines []string
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf(`		"%s": bson.M{`, key))

		for _, field := range u[key] {
			lines = append(lines, fmt.Sprintf(`			"%s": arg%d,`, field.BsonTag, field.ParamIndex))
		}

		lines = append(lines, `		},`)
	}
	return strings.Join(lines, "\n")
}

type querySpec struct {
	Operator   spec.Operator
	Predicates []predicate
}

func (q querySpec) Code() string {
	var predicateCodes []string
	for _, predicate := range q.Predicates {
		predicateCodes = append(predicateCodes, predicate.Code())
	}

	var lines []string
	switch q.Operator {
	case spec.OperatorOr:
		lines = append(lines, `		"$or": []bson.M{`)
		for _, predicateCode := range predicateCodes {
			lines = append(lines, fmt.Sprintf(`			{%s},`, predicateCode))
		}
		lines = append(lines, `		},`)
	case spec.OperatorAnd:
		lines = append(lines, `		"$and": []bson.M{`)
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
	ParamIndex int
}

func (p predicate) Code() string {
	switch p.Comparator {
	case spec.ComparatorEqual:
		return fmt.Sprintf(`"%s": arg%d`, p.Field, p.ParamIndex)
	case spec.ComparatorNot:
		return fmt.Sprintf(`"%s": bson.M{"$ne": arg%d}`, p.Field, p.ParamIndex)
	case spec.ComparatorLessThan:
		return fmt.Sprintf(`"%s": bson.M{"$lt": arg%d}`, p.Field, p.ParamIndex)
	case spec.ComparatorLessThanEqual:
		return fmt.Sprintf(`"%s": bson.M{"$lte": arg%d}`, p.Field, p.ParamIndex)
	case spec.ComparatorGreaterThan:
		return fmt.Sprintf(`"%s": bson.M{"$gt": arg%d}`, p.Field, p.ParamIndex)
	case spec.ComparatorGreaterThanEqual:
		return fmt.Sprintf(`"%s": bson.M{"$gte": arg%d}`, p.Field, p.ParamIndex)
	case spec.ComparatorBetween:
		return fmt.Sprintf(`"%s": bson.M{"$gte": arg%d, "$lte": arg%d}`, p.Field, p.ParamIndex, p.ParamIndex+1)
	case spec.ComparatorIn:
		return fmt.Sprintf(`"%s": bson.M{"$in": arg%d}`, p.Field, p.ParamIndex)
	case spec.ComparatorNotIn:
		return fmt.Sprintf(`"%s": bson.M{"$nin": arg%d}`, p.Field, p.ParamIndex)
	case spec.ComparatorTrue:
		return fmt.Sprintf(`"%s": true`, p.Field)
	case spec.ComparatorFalse:
		return fmt.Sprintf(`"%s": false`, p.Field)
	}
	return ""
}
