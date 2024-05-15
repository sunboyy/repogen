package codegen_test

import (
	"bytes"
	"go/token"
	"go/types"
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
		Pkg:  testutils.Pkg,
		Name: "User",
		Fields: []code.StructField{
			{
				Var: types.NewVar(token.NoPos, nil, "ID", testutils.TypeObjectIDNamed),
				Tag: `bson:"id,omitempty" json:"id,omitempty"`,
			},
			{
				Var: types.NewVar(token.NoPos, nil, "Username", code.TypeString),
				Tag: `bson:"username" json:"username"`,
			},
			{
				Var: types.NewVar(token.NoPos, nil, "Age", code.TypeInt),
				Tag: `bson:"age"`,
			},
			{
				Var: types.NewVar(token.NoPos, nil, "orderCount", types.NewPointer(code.TypeInt)),
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
