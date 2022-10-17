package codegen

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/sunboyy/repogen/internal/code"
)

const structTemplate = `
type {{.Name}} struct {
{{.GenFields}}
}
`

// StructBuilder is an implementer of a struct.
type StructBuilder struct {
	Name   string
	Fields code.StructFields
}

// Impl writes struct declatation code to the buffer.
func (sb StructBuilder) Impl(buffer *bytes.Buffer) error {
	tmpl, err := template.New("struct").Parse(structTemplate)
	if err != nil {
		return err
	}

	// writing to a buffer should not cause errors.
	_ = tmpl.Execute(buffer, sb)

	return nil
}

func (sb StructBuilder) GenFields() string {
	var fieldLines []string
	for _, field := range sb.Fields {
		fieldLine := fmt.Sprintf("\t%s %s", field.Name, field.Type.Code())
		if len(field.Tags) > 0 {
			fieldLine += fmt.Sprintf(" `%s`", sb.generateStructTag(field.Tags))
		}
		fieldLines = append(fieldLines, fieldLine)
	}
	return strings.Join(fieldLines, "\n")
}

func (sb StructBuilder) generateStructTag(tags map[string][]string) string {
	var tagKeys []string
	for key := range tags {
		tagKeys = append(tagKeys, key)
	}
	sort.Strings(tagKeys)

	var tagGroups []string
	for _, key := range tagKeys {
		tagValue := strings.Join(tags[key], ",")
		tagGroups = append(tagGroups, fmt.Sprintf("%s:\"%s\"", key, tagValue))
	}
	return strings.Join(tagGroups, " ")
}
