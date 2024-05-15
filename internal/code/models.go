package code

import (
	"go/types"
	"reflect"
)

// StructField is a definition of the struct field
type StructField struct {
	Var *types.Var
	Tag reflect.StructTag
}

var (
	TypeBool    = types.Typ[types.Bool]
	TypeInt     = types.Typ[types.Int]
	TypeInt64   = types.Typ[types.Int64]
	TypeFloat64 = types.Typ[types.Float64]
	TypeString  = types.Typ[types.String]
	TypeError   = types.Universe.Lookup("error").Type()
)
