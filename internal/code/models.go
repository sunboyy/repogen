package code

import (
	"fmt"
	"reflect"
)

// File is a container of all required components for code generation in the file
type File struct {
	PackageName string
	Imports     []Import
	Structs     []Struct
	Interfaces  []InterfaceType
}

// Import is a model for package imports
type Import struct {
	Name string
	Path string
}

// Struct is a definition of the struct
type Struct struct {
	Name   string
	Fields StructFields
}

// ReferencedType returns a type variable of this struct
func (str Struct) ReferencedType() Type {
	return SimpleType(str.Name)
}

// StructFields is a group of the StructField model
type StructFields []StructField

// ByName return struct field with matching name
func (fields StructFields) ByName(name string) (StructField, bool) {
	for _, field := range fields {
		if field.Name == name {
			return field, true
		}
	}
	return StructField{}, false
}

// StructField is a definition of the struct field
type StructField struct {
	Name string
	Type Type
	Tag  reflect.StructTag
}

// InterfaceType is a definition of the interface
type InterfaceType struct {
	Name    string
	Methods []Method
}

// Code returns token string in code format
func (intf InterfaceType) Code() string {
	return `interface{}`
}

// IsNumber returns false
func (intf InterfaceType) IsNumber() bool {
	return false
}

// Method is a definition of the method inside the interface
type Method struct {
	Name     string
	Comments []string
	Params   []Param
	Returns  []Type
}

// Param is a model of method parameter
type Param struct {
	Name string
	Type Type
}

// Type is an interface for value types
type Type interface {
	Code() string
	IsNumber() bool
}

// SimpleType is a type that can be called directly
type SimpleType string

// Code returns token string in code format
func (t SimpleType) Code() string {
	return string(t)
}

// IsNumber returns true id a SimpleType is integer or float variants.
func (t SimpleType) IsNumber() bool {
	return t == "uint" || t == "uint8" || t == "uint16" || t == "uint32" || t == "uint64" ||
		t == "int" || t == "int8" || t == "int16" || t == "int32" || t == "int64" ||
		t == "float32" || t == "float64"
}

// commonly-used types
const (
	TypeBool    = SimpleType("bool")
	TypeInt     = SimpleType("int")
	TypeFloat64 = SimpleType("float64")
	TypeString  = SimpleType("string")
	TypeError   = SimpleType("error")
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

// IsNumber returns false
func (t ExternalType) IsNumber() bool {
	return false
}

// PointerType is a model of pointer
type PointerType struct {
	ContainedType Type
}

// Code returns token string in code format
func (t PointerType) Code() string {
	return fmt.Sprintf("*%s", t.ContainedType.Code())
}

// IsNumber returns IsNumber of its contained type
func (t PointerType) IsNumber() bool {
	return t.ContainedType.IsNumber()
}

// ArrayType is a model of array
type ArrayType struct {
	ContainedType Type
}

// Code returns token string in code format
func (t ArrayType) Code() string {
	return fmt.Sprintf("[]%s", t.ContainedType.Code())
}

// IsNumber returns false
func (t ArrayType) IsNumber() bool {
	return false
}

// MapType is a model of map
type MapType struct {
	KeyType   Type
	ValueType Type
}

// Code returns token string in code format
func (t MapType) Code() string {
	return fmt.Sprintf("map[%s]%s", t.KeyType.Code(), t.ValueType.Code())
}

// IsNumber returns false
func (t MapType) IsNumber() bool {
	return false
}
