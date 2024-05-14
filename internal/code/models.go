package code

import (
	"fmt"
	"go/types"
	"reflect"
)

// Import is a model for package imports
type Import struct {
	Name string
	Path string
}

// LegacyStructField is a definition of the struct field
type LegacyStructField struct {
	Name string
	Type Type
	Tag  reflect.StructTag
}

// StructField is a definition of the struct field
type StructField struct {
	Var *types.Var
	Tag reflect.StructTag
}

// InterfaceType is a definition of the interface
type InterfaceType struct {
}

// Code returns token string in code format
func (intf InterfaceType) Code() string {
	return `interface{}`
}

// Type is an interface for value types
type Type interface {
	Code() string
}

// SimpleType is a type that can be called directly
type SimpleType string

// Code returns token string in code format
func (t SimpleType) Code() string {
	return string(t)
}

var (
	TypeBool    = types.Typ[types.Bool]
	TypeInt     = types.Typ[types.Int]
	TypeInt64   = types.Typ[types.Int64]
	TypeFloat64 = types.Typ[types.Float64]
	TypeString  = types.Typ[types.String]
	TypeError   = types.Universe.Lookup("error").Type()
)

// ExternalType is a type that is called to another package
type ExternalType struct {
	PackageAlias string
	Name         string
}

// Code returns token string in code format
func (t ExternalType) Code() string {
	return fmt.Sprintf("%s.%s", t.PackageAlias, t.Name)
}

// PointerType is a model of pointer
type PointerType struct {
	ContainedType Type
}

// Code returns token string in code format
func (t PointerType) Code() string {
	return fmt.Sprintf("*%s", t.ContainedType.Code())
}

// ArrayType is a model of array
type ArrayType struct {
	ContainedType Type
}

// Code returns token string in code format
func (t ArrayType) Code() string {
	return fmt.Sprintf("[]%s", t.ContainedType.Code())
}
