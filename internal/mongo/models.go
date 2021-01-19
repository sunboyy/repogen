package mongo

import (
	"fmt"

	"github.com/sunboyy/repogen/internal/spec"
)

type predicate struct {
	Field    string
	Operator spec.Operator
}

func (p predicate) Code(argIndex int) string {
	switch p.Operator {
	case spec.OperatorEqual:
		return fmt.Sprintf(`"%s": arg%d`, p.Field, argIndex)
	case spec.OperatorNot:
		return fmt.Sprintf(`"%s": bson.M{"$ne": arg%d}`, p.Field, argIndex)
	case spec.OperatorLessThan:
		return fmt.Sprintf(`"%s": bson.M{"$lt": arg%d}`, p.Field, argIndex)
	case spec.OperatorLessThanEqual:
		return fmt.Sprintf(`"%s": bson.M{"$lte": arg%d}`, p.Field, argIndex)
	case spec.OperatorGreaterThan:
		return fmt.Sprintf(`"%s": bson.M{"$gt": arg%d}`, p.Field, argIndex)
	case spec.OperatorGreaterThanEqual:
		return fmt.Sprintf(`"%s": bson.M{"$gte": arg%d}`, p.Field, argIndex)
	}
	return ""
}
