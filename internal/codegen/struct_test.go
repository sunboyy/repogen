package codegen_test

import (
	"bytes"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/testutils"
)

const expectedStructBuilderCode = `
type User struct {
	ID primitive.ObjectID ` + "`bson:\"id,omitempty\" json:\"id,omitempty\"`" + `
	Username string ` + "`bson:\"username\" json:\"username\"`" + `
	Age int ` + "`bson:\"age\"`" + `
	orderCount *int
}
`

func TestStructBuilderBuild(t *testing.T) {
	sb := codegen.StructBuilder{
		Name: "User",
		Fields: []code.StructField{
			{
				Name: "ID",
				Type: code.ExternalType{
					PackageAlias: "primitive",
					Name:         "ObjectID",
				},
				Tags: map[string][]string{
					"json": {"id", "omitempty"},
					"bson": {"id", "omitempty"},
				},
			},
			{
				Name: "Username",
				Type: code.TypeString,
				Tags: map[string][]string{
					"json": {"username"},
					"bson": {"username"},
				},
			},
			{
				Name: "Age",
				Type: code.TypeInt,
				Tags: map[string][]string{
					"bson": {"age"},
				},
			},
			{
				Name: "orderCount",
				Type: code.PointerType{
					ContainedType: code.TypeInt,
				},
			},
		},
	}
	buffer := new(bytes.Buffer)

	err := sb.Impl(buffer)

	if err != nil {
		t.Fatal(err)
	}
	actual := buffer.String()
	if err := testutils.ExpectMultiLineString(
		expectedStructBuilderCode,
		actual,
	); err != nil {
		t.Error(err)
	}
}
