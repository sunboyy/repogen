package generator

import (
	"bytes"
	"html/template"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
	"golang.org/x/tools/imports"
)

// GenerateRepository generates repository implementation from repository interface specification
func GenerateRepository(packageName string, structModel code.Struct, interfaceName string,
	methodSpecs []spec.MethodSpec) (string, error) {

	repositoryGenerator := repositoryGenerator{
		PackageName:   packageName,
		StructModel:   structModel,
		InterfaceName: interfaceName,
		MethodSpecs:   methodSpecs,
		Generator:     mongo.NewGenerator(structModel, interfaceName),
	}

	return repositoryGenerator.Generate()
}

type repositoryGenerator struct {
	PackageName   string
	StructModel   code.Struct
	InterfaceName string
	MethodSpecs   []spec.MethodSpec
	Generator     mongo.RepositoryGenerator
}

func (g repositoryGenerator) Generate() (string, error) {
	buffer := new(bytes.Buffer)
	if err := g.generateBase(buffer); err != nil {
		return "", err
	}

	if err := g.Generator.GenerateConstructor(buffer); err != nil {
		return "", err
	}

	for _, method := range g.MethodSpecs {
		if err := g.Generator.GenerateMethod(method, buffer); err != nil {
			return "", err
		}
	}

	formattedCode, err := imports.Process("", buffer.Bytes(), nil)
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}

func (g repositoryGenerator) generateBase(buffer *bytes.Buffer) error {
	tmpl, err := template.New("file_base").Parse(baseTemplate)
	if err != nil {
		return err
	}

	tmplData := baseTemplateData{
		PackageName: g.PackageName,
	}

	if err := tmpl.Execute(buffer, tmplData); err != nil {
		return err
	}

	return nil
}
