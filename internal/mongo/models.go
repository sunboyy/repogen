package mongo

import (
	"fmt"
	"sort"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
)

type updateField struct {
	BsonTag    string
	ParamIndex int
}

type update interface {
	Code() codegen.Statement
}

type updateModel struct {
}

func (u updateModel) Code() codegen.Statement {
	return codegen.MapStatement{
		Type: "bson.M",
		Pairs: []codegen.MapPair{
			{
				Key:   "$set",
				Value: codegen.Identifier("arg1"),
			},
		},
	}
}

type updateFields map[string][]updateField

func (u updateFields) Code() codegen.Statement {
	var keys []string
	for k := range u {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	stmt := codegen.MapStatement{
		Type: "bson.M",
	}
	for _, key := range keys {
		applicationMap := codegen.MapStatement{
			Type: "bson.M",
		}

		for _, field := range u[key] {
			applicationMap.Pairs = append(applicationMap.Pairs, codegen.MapPair{
				Key:   field.BsonTag,
				Value: codegen.Identifier(fmt.Sprintf("arg%d", field.ParamIndex)),
			})
		}

		stmt.Pairs = append(stmt.Pairs, codegen.MapPair{
			Key:   key,
			Value: applicationMap,
		})
	}
	return stmt
}

type querySpec struct {
	Operator   spec.Operator
	Predicates []predicate
}

func (q querySpec) Code() codegen.Statement {
	var predicatePairs []codegen.MapPair
	for _, predicate := range q.Predicates {
		predicatePairs = append(predicatePairs, predicate.Code())
	}
	var predicateMaps []codegen.Statement
	for _, pair := range predicatePairs {
		predicateMaps = append(predicateMaps, codegen.MapStatement{
			Pairs: []codegen.MapPair{pair},
		})
	}

	stmt := codegen.MapStatement{
		Type: "bson.M",
	}
	switch q.Operator {
	case spec.OperatorOr:
		stmt.Pairs = append(stmt.Pairs, codegen.MapPair{
			Key: "$or",
			Value: codegen.SliceStatement{
				Type: code.ArrayType{
					ContainedType: code.ExternalType{
						PackageAlias: "bson",
						Name:         "M",
					},
				},
				Values: predicateMaps,
			},
		})
	case spec.OperatorAnd:
		stmt.Pairs = append(stmt.Pairs, codegen.MapPair{
			Key: "$and",
			Value: codegen.SliceStatement{
				Type: code.ArrayType{
					ContainedType: code.ExternalType{
						PackageAlias: "bson",
						Name:         "M",
					},
				},
				Values: predicateMaps,
			},
		})
	default:
		stmt.Pairs = predicatePairs
	}
	return stmt
}

type predicate struct {
	Field      string
	Comparator spec.Comparator
	ParamIndex int
}

func (p predicate) Code() codegen.MapPair {
	argStmt := codegen.Identifier(fmt.Sprintf("arg%d", p.ParamIndex))

	switch p.Comparator {
	case spec.ComparatorEqual:
		return p.createValueMapPair(argStmt)
	case spec.ComparatorNot:
		return p.createSingleComparisonMapPair("$ne", argStmt)
	case spec.ComparatorLessThan:
		return p.createSingleComparisonMapPair("$lt", argStmt)
	case spec.ComparatorLessThanEqual:
		return p.createSingleComparisonMapPair("$lte", argStmt)
	case spec.ComparatorGreaterThan:
		return p.createSingleComparisonMapPair("$gt", argStmt)
	case spec.ComparatorGreaterThanEqual:
		return p.createSingleComparisonMapPair("$gte", argStmt)
	case spec.ComparatorBetween:
		argStmt2 := codegen.Identifier(fmt.Sprintf("arg%d", p.ParamIndex+1))
		return p.createBetweenMapPair(argStmt, argStmt2)
	case spec.ComparatorIn:
		return p.createSingleComparisonMapPair("$in", argStmt)
	case spec.ComparatorNotIn:
		return p.createSingleComparisonMapPair("$nin", argStmt)
	case spec.ComparatorTrue:
		return p.createValueMapPair(codegen.Identifier("true"))
	case spec.ComparatorFalse:
		return p.createValueMapPair(codegen.Identifier("false"))
	case spec.ComparatorExists:
		return p.createExistsMapPair("1")
	case spec.ComparatorNotExists:
		return p.createExistsMapPair("0")
	}
	return codegen.MapPair{}
}

func (p predicate) createValueMapPair(
	argStmt codegen.Statement) codegen.MapPair {

	return codegen.MapPair{
		Key:   p.Field,
		Value: argStmt,
	}
}

func (p predicate) createSingleComparisonMapPair(comparatorKey string,
	argStmt codegen.Statement) codegen.MapPair {

	return codegen.MapPair{
		Key: p.Field,
		Value: codegen.MapStatement{
			Type:  "bson.M",
			Pairs: []codegen.MapPair{{Key: comparatorKey, Value: argStmt}},
		},
	}
}

func (p predicate) createBetweenMapPair(argStmt codegen.Statement,
	argStmt2 codegen.Statement) codegen.MapPair {

	return codegen.MapPair{
		Key: p.Field,
		Value: codegen.MapStatement{
			Type: "bson.M",
			Pairs: []codegen.MapPair{
				{Key: "$gte", Value: argStmt},
				{Key: "$lte", Value: argStmt2},
			},
		},
	}
}

func (p predicate) createExistsMapPair(existsValue string) codegen.MapPair {
	return codegen.MapPair{
		Key: p.Field,
		Value: codegen.MapStatement{
			Type: "bson.M",
			Pairs: []codegen.MapPair{{
				Key:   "$exists",
				Value: codegen.Identifier(existsValue),
			}},
		},
	}
}
