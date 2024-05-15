package generator

import (
	"go/types"

	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
)

// GenerateRepository generates repository implementation code from repository
// interface specification.
func GenerateRepository(pkg *types.Package, namedStruct *types.Named,
	interfaceName string, methodSpecs []spec.MethodSpec) (string, error) {

	generator := mongo.NewGenerator(pkg, namedStruct, interfaceName)

	codeBuilder := codegen.NewBuilder(
		"repogen",
		pkg.Name(),
		generator.Imports(),
	)

	constructorBuilder, err := generator.GenerateConstructor()
	if err != nil {
		return "", err
	}

	codeBuilder.AddImplementer(constructorBuilder)
	codeBuilder.AddImplementer(generator.GenerateStruct())

	for _, method := range methodSpecs {
		methodBuilder, err := generator.GenerateMethod(method)
		if err != nil {
			return "", err
		}
		codeBuilder.AddImplementer(methodBuilder)
	}

	return codeBuilder.Build()
}
