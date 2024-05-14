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
		Fields: []code.LegacyStructField{
			{
				Name: "ID",
				Type: code.ExternalType{
					PackageAlias: "primitive",
					Name:         "ObjectID",
				},
				Tag: `bson:"id,omitempty" json:"id,omitempty"`,
			},
			{
				Name: "Username",
				Type: code.SimpleType("string"),
				Tag:  `bson:"username" json:"username"`,
			},
			{
				Name: "Age",
				Type: code.SimpleType("int"),
				Tag:  `bson:"age"`,
			},
			{
				Name: "orderCount",
				Type: code.PointerType{
					ContainedType: code.SimpleType("int"),
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
