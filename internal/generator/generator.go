package generator

import (
	"go/types"
	"log"

	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
)

func GenerateRepositoryImpl(pkg, modelPkg *types.Package, structModelName,
	repoInterfaceName, outputPkgName string) (string, error) {

	namedStruct, intf, err := deriveSourceTypes(pkg, modelPkg, structModelName,
		repoInterfaceName)
	if err != nil {
		return "", err
	}

	methodSpecs, err := constructRepositorySpec(pkg, namedStruct, intf)
	if err != nil {
		return "", err
	}

	codeBuilder, err := constructCodeBuilder(pkg, namedStruct,
		repoInterfaceName, methodSpecs, outputPkgName)
	if err != nil {
		return "", err
	}

	return codeBuilder.Build()
}

func deriveSourceTypes(pkg, modelPkg *types.Package, structModelName string,
	repositoryInterfaceName string) (*types.Named, *types.Interface, error) {

	structModelObj := modelPkg.Scope().Lookup(structModelName)
	if structModelObj == nil {
		return nil, nil, ErrStructNotFound
	}
	namedStruct := structModelObj.Type().(*types.Named)
	if _, ok := namedStruct.Underlying().(*types.Struct); !ok {
		return nil, nil, ErrNotNamedStruct
	}

	intfObj := pkg.Scope().Lookup(repositoryInterfaceName)
	if intfObj == nil {
		return nil, nil, ErrInterfaceNotFound
	}
	intf, ok := intfObj.Type().Underlying().(*types.Interface)
	if !ok {
		return nil, nil, ErrNotInterface
	}

	return namedStruct, intf, nil
}

func constructRepositorySpec(pkg *types.Package, namedStruct *types.Named,
	intf *types.Interface) ([]spec.MethodSpec, error) {

	var methodSpecs []spec.MethodSpec
	for i := 0; i < intf.NumMethods(); i++ {
		method := intf.Method(i)
		log.Println("Generating method:", method.Name())

		methodSpec, err := spec.ParseInterfaceMethod(pkg, namedStruct, method)
		if err != nil {
			return nil, err
		}
		methodSpecs = append(methodSpecs, methodSpec)
	}

	return methodSpecs, nil
}

func constructCodeBuilder(pkg *types.Package, namedStruct *types.Named,
	interfaceName string, methodSpecs []spec.MethodSpec, outputPkgName string) (*codegen.Builder, error) {

	generator := mongo.NewGenerator(pkg, namedStruct, interfaceName)
	if outputPkgName == "" {
		outputPkgName = pkg.Name()
	}
	codeBuilder := codegen.NewBuilder(
		"repogen",
		outputPkgName,
		generator.Imports(),
	)

	constructorBuilder, err := generator.GenerateConstructor()
	if err != nil {
		return nil, err
	}

	codeBuilder.AddImplementer(constructorBuilder)
	codeBuilder.AddImplementer(generator.GenerateStruct())

	for _, method := range methodSpecs {
		methodBuilder, err := generator.GenerateMethod(method)
		if err != nil {
			return nil, err
		}
		codeBuilder.AddImplementer(methodBuilder)
	}

	return codeBuilder, nil
}
