package generator

import (
	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
)

// GenerateRepository generates repository implementation code from repository
// interface specification.
func GenerateRepository(modelPackagePath, packageName string, structModel code.Struct,
	interfaceName string, methodSpecs []spec.MethodSpec) (string, error) {

	generator := mongo.NewGenerator(structModel, interfaceName)

	codeBuilder := codegen.NewBuilder(
		"repogen",
		packageName,
		append(generator.Imports(), []code.Import{{Path: modelPackagePath}}),
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
