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
		return codegen.MapPair{Key: p.Field, Value: argStmt}
	case spec.ComparatorNot:
		return codegen.MapPair{
			Key: p.Field,
			Value: codegen.MapStatement{
				Type:  "bson.M",
				Pairs: []codegen.MapPair{{Key: "$ne", Value: argStmt}},
			},
		}
	case spec.ComparatorLessThan:
		return codegen.MapPair{
			Key: p.Field,
			Value: codegen.MapStatement{
				Type:  "bson.M",
				Pairs: []codegen.MapPair{{Key: "$lt", Value: argStmt}},
			},
		}
	case spec.ComparatorLessThanEqual:
		return codegen.MapPair{
			Key: p.Field,
			Value: codegen.MapStatement{
				Type:  "bson.M",
				Pairs: []codegen.MapPair{{Key: "$lte", Value: argStmt}},
			},
		}
	case spec.ComparatorGreaterThan:
		return codegen.MapPair{
			Key: p.Field,
			Value: codegen.MapStatement{
				Type:  "bson.M",
				Pairs: []codegen.MapPair{{Key: "$gt", Value: argStmt}},
			},
		}
	case spec.ComparatorGreaterThanEqual:
		return codegen.MapPair{
			Key: p.Field,
			Value: codegen.MapStatement{
				Type:  "bson.M",
				Pairs: []codegen.MapPair{{Key: "$gte", Value: argStmt}},
			},
		}
	case spec.ComparatorBetween:
		argStmt2 := codegen.Identifier(fmt.Sprintf("arg%d", p.ParamIndex+1))
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
	case spec.ComparatorIn:
		return codegen.MapPair{
			Key: p.Field,
			Value: codegen.MapStatement{
				Type:  "bson.M",
				Pairs: []codegen.MapPair{{Key: "$in", Value: argStmt}},
			},
		}
	case spec.ComparatorNotIn:
		return codegen.MapPair{
			Key: p.Field,
			Value: codegen.MapStatement{
				Type:  "bson.M",
				Pairs: []codegen.MapPair{{Key: "$nin", Value: argStmt}},
			},
		}
	case spec.ComparatorTrue:
		return codegen.MapPair{
			Key:   p.Field,
			Value: codegen.Identifier("true"),
		}
	case spec.ComparatorFalse:
		return codegen.MapPair{
			Key:   p.Field,
			Value: codegen.Identifier("false"),
		}
	case spec.ComparatorExists:
		return codegen.MapPair{
			Key: p.Field,
			Value: codegen.MapStatement{
				Type: "bson.M",
				Pairs: []codegen.MapPair{{
					Key:   "$exists",
					Value: codegen.Identifier("1"),
				}},
			},
		}
	case spec.ComparatorNotExists:
		return codegen.MapPair{
			Key: p.Field,
			Value: codegen.MapStatement{
				Type: "bson.M",
				Pairs: []codegen.MapPair{{
					Key:   "$exists",
					Value: codegen.Identifier("0"),
				}},
			},
		}
	}
	return codegen.MapPair{}
}
