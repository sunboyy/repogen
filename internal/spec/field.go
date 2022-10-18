package spec

import (
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
		fieldNames = append(fieldNames, field.Name)
	}
	return strings.Join(fieldNames, ".")
}

type fieldResolver struct {
	Structs map[string]code.Struct
}

func (r fieldResolver) ResolveStructField(structModel code.Struct, tokens []string) (FieldReference, bool) {
	fieldName := strings.Join(tokens, "")
	field, ok := structModel.Fields.ByName(fieldName)
	if ok {
		return FieldReference{field}, true
	}

	for i := len(tokens) - 1; i > 0; i-- {
		fieldName := strings.Join(tokens[:i], "")
		field, ok := structModel.Fields.ByName(fieldName)
		if !ok {
			continue
		}

		fieldSimpleType, ok := getSimpleType(field.Type)
		if !ok {
			continue
		}

		childStruct, ok := r.Structs[fieldSimpleType.Code()]
		if !ok {
			continue
		}

		fields, ok := r.ResolveStructField(childStruct, tokens[i:])
		if !ok {
			continue
		}
		return append(FieldReference{field}, fields...), true
	}

	return nil, false
}

func getSimpleType(t code.Type) (code.SimpleType, bool) {
	switch t := t.(type) {
	case code.SimpleType:
		return t, true
	case code.PointerType:
		return getSimpleType(t.ContainedType)
	default:
		return "", false
	}
}
