package code_test

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
)

const goImplFile1Data = `
package codepkgsuccess

import (
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

type User struct {
	ID       primitive.ObjectID ` + "`json:\"id\"`" + `
	Name     string             ` + "`json:\"name\"`" + `
	Gender   Gender             ` + "`json:\"gender\"`" + `
	Birthday time.Time          ` + "`json:\"birthday\"`" + `
}

func (u User) Age() int {
	return int(math.Floor(time.Since(u.Birthday).Hours() / 24 / 365))
}

type (
	Product struct {
		ID    primitive.ObjectID ` + "`json:\"id\"`" + `
		Name  string             ` + "`json:\"name\"`" + `
		Price float64            ` + "`json:\"price\"`" + `
	}

	Order struct {
		ID         primitive.ObjectID         ` + "`json:\"id\"`" + `
		ItemIDs    map[primitive.ObjectID]int ` + "`json:\"itemIds\"`" + `
		TotalPrice float64                    ` + "`json:\"totalPrice\"`" + `
		UserID     primitive.ObjectID         ` + "`json:\"userId\"`" + `
		CreatedAt  time.Time                  ` + "`json:\"createdAt\"`" + `
	}
)
`

const goImplFile2Data = `
package codepkgsuccess

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService interface {
	CreateOrder(u User, products map[Product]int) Order
}

type OrderServiceImpl struct{}

func (s *OrderServiceImpl) CreateOrder(u User, products map[Product]int) Order {
	itemIDs := map[primitive.ObjectID]int{}
	var totalPrice float64
	for product, amount := range products {
		itemIDs[product.ID] = amount
		totalPrice += product.Price * float64(amount)
	}

	return Order{
		ID:         primitive.NewObjectID(),
		ItemIDs:    map[primitive.ObjectID]int{},
		TotalPrice: totalPrice,
		UserID:     u.ID,
		CreatedAt:  time.Now(),
	}
}
`

const goImplFile3Data = `
package success
`

const goImplFile4Data = `
package codepkgsuccess

type User struct {
	Name     string
}
`

const goImplFile5Data = `
package codepkgsuccess

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderService interface {
	CancelOrder(orderID primitive.ObjectID) error
}
`

const goTestFileData = `
package codepkgsuccess

type TestCase struct {
	Name     string
	Params   []interface{}
	Expected string
	Actual   string
}
`

var (
	goImplFile1 *ast.File
	goImplFile2 *ast.File
	goImplFile3 *ast.File
	goImplFile4 *ast.File
	goImplFile5 *ast.File
	goTestFile  *ast.File
)

func init() {
	fset := token.NewFileSet()
	goImplFile1, _ = parser.ParseFile(fset, "", goImplFile1Data, parser.ParseComments)
	goImplFile2, _ = parser.ParseFile(fset, "", goImplFile2Data, parser.ParseComments)
	goImplFile3, _ = parser.ParseFile(fset, "", goImplFile3Data, parser.ParseComments)
	goImplFile4, _ = parser.ParseFile(fset, "", goImplFile4Data, parser.ParseComments)
	goImplFile5, _ = parser.ParseFile(fset, "", goImplFile5Data, parser.ParseComments)
	goTestFile, _ = parser.ParseFile(fset, "", goTestFileData, parser.ParseComments)
}

const testpkgpath = "github/usr/codepkgsuccess"

var (
	mockPackPathParser = func(pkgName string) (string, error) {
		return testpkgpath, nil
	}
)

func TestParsePackage_Success(t *testing.T) {
	mockDirParser := func(dir string) (pkgs map[string]*ast.Package, err error) {
		return map[string]*ast.Package{
			"codepkgsuccess": {
				Files: map[string]*ast.File{
					"file1.go":      goImplFile1,
					"file2.go":      goImplFile2,
					"file1_test.go": goTestFile,
				},
			},
		}, nil
	}

	parser := code.NewPackageParser()
	parser.DirParser = mockDirParser
	parser.PkgPathParser = mockPackPathParser

	pkg, err := parser.ParsePackage(testpkgpath)
	if err != nil {
		t.Fatal(err)
	}

	if pkg.Name != "codepkgsuccess" {
		t.Errorf("expected package name 'codepkgsuccess', got '%s'", pkg.Name)
	}
	if _, ok := pkg.Structs["User"]; !ok {
		t.Error("struct 'User' not found")
	}
	if _, ok := pkg.Structs["Product"]; !ok {
		t.Error("struct 'Product' not found")
	}
	if _, ok := pkg.Structs["Order"]; !ok {
		t.Error("struct 'Order' not found")
	}
	if _, ok := pkg.Structs["OrderServiceImpl"]; !ok {
		t.Error("struct 'OrderServiceImpl' not found")
	}
	if _, ok := pkg.Interfaces["OrderService"]; !ok {
		t.Error("interface 'OrderService' not found")
	}
	if _, ok := pkg.Structs["TestCase"]; ok {
		t.Error("unexpected struct 'TestCase' in test file")
	}

	if pkg.Path != testpkgpath {
		t.Errorf("expected package path '%s', got '%s'", testpkgpath, pkg.Path)
	}
}

func TestParsePackage_AmbiguousPackageName(t *testing.T) {
	mockDirParser := func(dir string) (pkgs map[string]*ast.Package, err error) {
		return map[string]*ast.Package{
			"codepkgsuccess": {
				Files: map[string]*ast.File{
					"file1.go": goImplFile1,
					"file2.go": goImplFile2,
					"file3.go": goImplFile3,
				},
			},
		}, nil
	}
	parser := code.NewPackageParser()
	parser.DirParser = mockDirParser
	parser.PkgPathParser = mockPackPathParser

	_, err := parser.ParsePackage(testpkgpath)

	if !errors.Is(err, code.ErrAmbiguousPackageName) {
		t.Errorf(
			"expected error '%s', got '%s'",
			code.ErrAmbiguousPackageName.Error(),
			err.Error(),
		)
	}
}

func TestParsePackage_DuplicateStructs(t *testing.T) {
	mockDirParser := func(dir string) (pkgs map[string]*ast.Package, err error) {
		return map[string]*ast.Package{
			"codepkgsuccess": {
				Files: map[string]*ast.File{
					"file1.go": goImplFile1,
					"file2.go": goImplFile2,
					"file4.go": goImplFile4,
				},
			},
		}, nil
	}
	parser := code.NewPackageParser()
	parser.DirParser = mockDirParser
	parser.PkgPathParser = mockPackPathParser

	_, err := parser.ParsePackage(testpkgpath)

	if !errors.Is(err, code.DuplicateStructError("User")) {
		t.Errorf(
			"expected error '%s', got '%s'",
			code.ErrAmbiguousPackageName.Error(),
			err.Error(),
		)
	}
}

func TestParsePackage_DuplicateInterfaces(t *testing.T) {
	mockDirParser := func(dir string) (pkgs map[string]*ast.Package, err error) {
		return map[string]*ast.Package{
			"codepkgsuccess": {
				Files: map[string]*ast.File{
					"file1.go": goImplFile1,
					"file2.go": goImplFile2,
					"file5.go": goImplFile5,
				},
			},
		}, nil
	}
	parser := code.NewPackageParser()
	parser.DirParser = mockDirParser
	parser.PkgPathParser = mockPackPathParser
	_, err := parser.ParsePackage(testpkgpath)

	if !errors.Is(err, code.DuplicateInterfaceError("OrderService")) {
		t.Errorf(
			"expected error '%s', got '%s'",
			code.ErrAmbiguousPackageName.Error(),
			err.Error(),
		)
	}
}
