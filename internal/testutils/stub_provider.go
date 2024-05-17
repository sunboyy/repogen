package testutils

import (
	"go/types"
	"reflect"

	"github.com/sunboyy/repogen/internal/code"
	"golang.org/x/tools/go/packages"
)

var (
	TypeContextNamed    *types.Named
	TypeObjectIDNamed   *types.Named
	TypeCollectionNamed *types.Named

	Pkg                     *types.Package
	TypeUserNamed           *types.Named
	TypeUserStruct          *types.Struct
	TypeGenderNamed         *types.Named
	TypeNameStruct          *types.Struct
	TypeConsentHistoryNamed *types.Named
)

func init() {
	cfg := &packages.Config{Mode: packages.NeedTypes}

	contextPkgs, err := packages.Load(cfg, "context")
	if err != nil {
		panic(err)
	}
	TypeContextNamed = contextPkgs[0].Types.Scope().Lookup("Context").Type().(*types.Named)

	primitivePkgs, err := packages.Load(cfg, "go.mongodb.org/mongo-driver/bson/primitive")
	if err != nil {
		panic(err)
	}
	TypeObjectIDNamed = primitivePkgs[0].Types.Scope().Lookup("ObjectID").Type().(*types.Named)

	mongoPkgs, err := packages.Load(cfg, "go.mongodb.org/mongo-driver/mongo")
	if err != nil {
		panic(err)
	}
	TypeCollectionNamed = mongoPkgs[0].Types.Scope().Lookup("Collection").Type().(*types.Named)

	stubPkgs, err := packages.Load(cfg, "../teststub")
	if err != nil {
		panic(err)
	}
	Pkg = stubPkgs[0].Types
	TypeUserNamed = Pkg.Scope().Lookup("User").Type().(*types.Named)
	TypeUserStruct = TypeUserNamed.Underlying().(*types.Struct)
	TypeGenderNamed = Pkg.Scope().Lookup("Gender").Type().(*types.Named)
	TypeNameStruct = Pkg.Scope().Lookup("Name").Type().Underlying().(*types.Struct)
	TypeConsentHistoryNamed = Pkg.Scope().Lookup("ConsentHistory").Type().(*types.Named)
}

func FindStructFieldByName(s *types.Struct, name string) code.StructField {
	fieldIndex := -1
	for i := 0; i < s.NumFields(); i++ {
		if s.Field(i).Name() == name {
			fieldIndex = i
			break
		}
	}

	return code.StructField{
		Var: s.Field(fieldIndex),
		Tag: reflect.StructTag(s.Tag(fieldIndex)),
	}
}
