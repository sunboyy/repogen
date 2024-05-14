package mongo

import (
	"fmt"
	"go/token"
	"go/types"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/spec"
	"golang.org/x/tools/go/packages"
)

// NewGenerator creates a new instance of MongoDB repository generator
func NewGenerator(pkg *types.Package, structModelName string, interfaceName string) RepositoryGenerator {
	return RepositoryGenerator{
		baseMethodGenerator: baseMethodGenerator{
			pkg:             pkg,
			structModelName: structModelName,
		},
		InterfaceName: interfaceName,
	}
}

// RepositoryGenerator is a MongoDB repository generator that provides
// necessary information required to construct an implementation.
type RepositoryGenerator struct {
	baseMethodGenerator
	InterfaceName string
}

// Imports returns necessary imports for the mongo repository implementation.
func (g RepositoryGenerator) Imports() [][]code.Import {
	return [][]code.Import{
		{
			{Path: "context"},
		},
		{
			{Path: "go.mongodb.org/mongo-driver/bson"},
			{Path: "go.mongodb.org/mongo-driver/bson/primitive"},
			{Path: "go.mongodb.org/mongo-driver/mongo"},
			{Path: "go.mongodb.org/mongo-driver/mongo/options"},
		},
	}
}

// GenerateStruct creates codegen.StructBuilder of mongo repository
// implementation struct.
func (g RepositoryGenerator) GenerateStruct() codegen.StructBuilder {
	return codegen.StructBuilder{
		Name: g.repoImplStructName(),
		Fields: []code.LegacyStructField{
			{
				Name: "collection",
				Type: code.PointerType{
					ContainedType: code.ExternalType{
						PackageAlias: "mongo",
						Name:         "Collection",
					},
				},
			},
		},
	}
}

// GenerateConstructor creates codegen.FunctionBuilder of a constructor for
// mongo repository implementation struct.
func (g RepositoryGenerator) GenerateConstructor() (codegen.FunctionBuilder, error) {
	mongoPkgs, err := packages.Load(&packages.Config{Mode: packages.NeedTypes}, "go.mongodb.org/mongo-driver/mongo")
	if err != nil {
		return codegen.FunctionBuilder{}, err
	}
	mongoPkg := mongoPkgs[0]
	collectionObj := mongoPkg.Types.Scope().Lookup("Collection")
	collectionType := collectionObj.Type()

	return codegen.FunctionBuilder{
		Pkg:    g.pkg,
		Name:   "New" + g.InterfaceName,
		Params: types.NewTuple(types.NewVar(token.NoPos, nil, "collection", types.NewPointer(collectionType))),
		Returns: []types.Type{
			types.NewPointer(types.NewNamed(
				types.NewTypeName(token.NoPos, nil, g.repoImplStructName(), nil), nil, nil)),
		},
		Body: codegen.FunctionBody{
			codegen.ReturnStatement{
				codegen.StructStatement{
					Type: fmt.Sprintf("&%s", g.repoImplStructName()),
					Pairs: []codegen.StructFieldPair{{
						Key:   "collection",
						Value: codegen.Identifier("collection"),
					}},
				},
			},
		},
	}, nil
}

// GenerateMethod creates codegen.MethodBuilder of repository method from the
// provided method specification.
func (g RepositoryGenerator) GenerateMethod(methodSpec spec.MethodSpec) (codegen.MethodBuilder, error) {
	var paramVars []*types.Var
	for i := 0; i < methodSpec.Signature.Params().Len(); i++ {
		param := types.NewVar(token.NoPos, nil, fmt.Sprintf("arg%d", i),
			methodSpec.Signature.Params().At(i).Type())
		paramVars = append(paramVars, param)
	}

	var returns []types.Type
	for i := 0; i < methodSpec.Signature.Results().Len(); i++ {
		returns = append(returns, methodSpec.Signature.Results().At(i).Type())
	}

	implementation, err := g.generateMethodImplementation(methodSpec)
	if err != nil {
		return codegen.MethodBuilder{}, err
	}

	return codegen.MethodBuilder{
		Pkg: g.pkg,
		Receiver: codegen.MethodReceiver{
			Name:    "r",
			Type:    code.SimpleType(g.repoImplStructName()),
			Pointer: true,
		},
		Name:    methodSpec.Name,
		Params:  types.NewTuple(paramVars...),
		Returns: returns,
		Body:    implementation,
	}, nil
}

func (g RepositoryGenerator) generateMethodImplementation(
	methodSpec spec.MethodSpec) (codegen.FunctionBody, error) {

	switch operation := methodSpec.Operation.(type) {
	case spec.InsertOperation:
		return g.generateInsertBody(operation), nil
	case spec.FindOperation:
		return g.generateFindBody(operation)
	case spec.UpdateOperation:
		return g.generateUpdateBody(operation)
	case spec.DeleteOperation:
		return g.generateDeleteBody(operation)
	case spec.CountOperation:
		return g.generateCountBody(operation)
	default:
		return nil, NewOperationNotSupportedError(operation.Name())
	}
}

func (g RepositoryGenerator) repoImplStructName() string {
	return g.InterfaceName + "Mongo"
}
