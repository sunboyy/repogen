package code

import (
	"fmt"
)

// File is a container of all required components for code generation in the file
type File struct {
	PackageName string
	Imports     []Import
	Structs     Structs
	Interfaces  Interfaces
}

// Import is a model for package imports
type Import struct {
	Name string
	Path string
}

// Structs is a group of Struct model
type Structs []Struct

// ByName return struct with matching name. Another return value shows whether there is a struct
// with that name exists.
func (strs Structs) ByName(name string) (Struct, bool) {
	for _, str := range strs {
		if str.Name == name {
			return str, true
		}
	}
	return Struct{}, false
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
	Tags map[string][]string
}

// Interfaces is a group of Interface model
type Interfaces []InterfaceType

// ByName return interface by name Another return value shows whether there is an interface
// with that name exists.
func (intfs Interfaces) ByName(name string) (InterfaceType, bool) {
	for _, intf := range intfs {
		if intf.Name == name {
			return intf, true
		}
	}
	return InterfaceType{}, false
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

// Method is a definition of the method inside the interface
type Method struct {
	Name    string
	Params  []Param
	Returns []Type
}

// Param is a model of method parameter
type Param struct {
	Name string
	Type Type
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
