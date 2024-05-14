package spec

import (
	"go/types"
	"reflect"
	"strings"

	"github.com/sunboyy/repogen/internal/code"
)

// FieldReference is a reference path to access to the field
type FieldReference []code.StructField

// ReferencedField returns the last struct field
func (r FieldReference) ReferencedField() code.StructField {
	return r[len(r)-1]
}

// ReferencingCode returns a string containing name of the referenced fields
// concatenating with period (.).
func (r FieldReference) ReferencingCode() string {
	var fieldNames []string
	for _, field := range r {
		fieldNames = append(fieldNames, field.Var.Name())
	}
	return strings.Join(fieldNames, ".")
}

func resolveStructField(structModel *types.Struct, tokens []string) (FieldReference, bool) {
	fieldName := strings.Join(tokens, "")
	for i := 0; i < structModel.NumFields(); i++ {
		field := structModel.Field(i)
		if field.Name() == fieldName {
			return FieldReference{
				code.StructField{
					Var: field,
					Tag: reflect.StructTag(structModel.Tag(i)),
				},
			}, true
		}
	}

	for i := len(tokens) - 1; i > 0; i-- {
		fieldName := strings.Join(tokens[:i], "")
		var foundField *types.Var
		var foundFieldIndex int
		for j := 0; j < structModel.NumFields(); j++ {
			field := structModel.Field(j)
			if field.Name() == fieldName {
				foundField = field
				foundFieldIndex = j
				break
			}
		}
		if foundField == nil {
			continue
		}

		underlyingStructType, ok := getUnderlyingStructType(foundField.Type())
		if !ok {
			continue
		}

		fields, ok := resolveStructField(underlyingStructType, tokens[i:])
		if !ok {
			continue
		}

		return append(FieldReference{
			code.StructField{
				Var: foundField,
				Tag: reflect.StructTag(structModel.Tag(foundFieldIndex)),
			},
		}, fields...), true
	}

	return nil, false
}

func getUnderlyingStructType(t types.Type) (*types.Struct, bool) {
	switch t := t.(type) {
	case *types.Named:
		return getUnderlyingStructType(t.Underlying())

	case *types.Struct:
		return t, true

	case *types.Pointer:
		return getUnderlyingStructType(t.Elem())

	default:
		return nil, false
	}
}
