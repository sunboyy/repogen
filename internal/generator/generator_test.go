package generator_test

import (
	"os"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/generator"
	"github.com/sunboyy/repogen/internal/spec"
	"github.com/sunboyy/repogen/internal/testutils"
)

var (
	idField = code.StructField{
		Name: "ID",
		Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
		Tag:  `bson:"_id,omitempty"`,
	}
	genderField = code.StructField{
		Name: "Gender",
		Type: code.SimpleType("Gender"),
		Tag:  `bson:"gender"`,
	}
	ageField = code.StructField{
		Name: "Age",
		Type: code.TypeInt,
		Tag:  `bson:"age"`,
	}
)

func TestGenerateMongoRepository(t *testing.T) {
	userModel := code.Struct{
		Name: "UserModel",
		Fields: code.StructFields{
			idField,
			code.StructField{
				Name: "Username",
				Type: code.TypeString,
				Tag:  `bson:"username"`,
			},
			genderField,
			ageField,
		},
	}
	methods := []spec.MethodSpec{
		// test find: One mode
		{
			Name: "FindByID",
			Params: []code.Param{
				{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
			},
			Returns: []code.Type{code.PointerType{ContainedType: code.SimpleType("UserModel")}, code.TypeError},
			Operation: spec.FindOperation{
				Mode: spec.QueryModeOne,
				Query: spec.QuerySpec{
					Predicates: []spec.Predicate{
						{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
					},
				},
			},
		},
		// test find: Many mode, And operator, NOT and LessThan comparator
		{
			Name: "FindByGenderNotAndAgeLessThan",
			Params: []code.Param{
				{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				{Name: "gender", Type: code.SimpleType("Gender")},
				{Name: "age", Type: code.TypeInt},
			},
			Returns: []code.Type{
				code.PointerType{ContainedType: code.SimpleType("UserModel")},
				code.TypeError,
			},
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorAnd,
					Predicates: []spec.Predicate{
						{
							FieldReference: spec.FieldReference{genderField},
							Comparator:     spec.ComparatorNot,
							ParamIndex:     1,
						},
						{
							FieldReference: spec.FieldReference{ageField},
							Comparator:     spec.ComparatorLessThan,
							ParamIndex:     2,
						},
					},
				},
			},
		},
		{
			Name: "FindByAgeLessThanEqualOrderByAge",
			Params: []code.Param{
				{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				{Name: "age", Type: code.TypeInt},
			},
			Returns: []code.Type{
				code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				code.TypeError,
			},
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Predicates: []spec.Predicate{
						{
							FieldReference: spec.FieldReference{ageField},
							Comparator:     spec.ComparatorLessThanEqual,
							ParamIndex:     1,
						},
					},
				},
				Sorts: []spec.Sort{
					{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingAscending},
				},
			},
		},
		{
			Name: "FindByAgeGreaterThanOrderByAgeAsc",
			Params: []code.Param{
				{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				{Name: "age", Type: code.TypeInt},
			},
			Returns: []code.Type{
				code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				code.TypeError,
			},
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Predicates: []spec.Predicate{
						{
							FieldReference: spec.FieldReference{ageField},
							Comparator:     spec.ComparatorGreaterThan,
							ParamIndex:     1,
						},
					},
				},
				Sorts: []spec.Sort{
					{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingAscending},
				},
			},
		},
		{
			Name: "FindByAgeGreaterThanEqualOrderByAgeDesc",
			Params: []code.Param{
				{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				{Name: "age", Type: code.TypeInt},
			},
			Returns: []code.Type{
				code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				code.TypeError,
			},
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Predicates: []spec.Predicate{
						{
							FieldReference: spec.FieldReference{ageField},
							Comparator:     spec.ComparatorGreaterThanEqual,
							ParamIndex:     1,
						},
					},
				},
				Sorts: []spec.Sort{
					{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingDescending},
				},
			},
		},
		{
			Name: "FindByAgeBetween",
			Params: []code.Param{
				{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				{Name: "fromAge", Type: code.TypeInt},
				{Name: "toAge", Type: code.TypeInt},
			},
			Returns: []code.Type{
				code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				code.TypeError,
			},
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Predicates: []spec.Predicate{
						{
							FieldReference: spec.FieldReference{ageField},
							Comparator:     spec.ComparatorBetween,
							ParamIndex:     1,
						},
					},
				},
			},
		},
		{
			Name: "FindByGenderOrAge",
			Params: []code.Param{
				{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				{Name: "gender", Type: code.SimpleType("Gender")},
				{Name: "age", Type: code.TypeInt},
			},
			Returns: []code.Type{
				code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				code.TypeError,
			},
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorOr,
					Predicates: []spec.Predicate{
						{
							FieldReference: spec.FieldReference{genderField},
							Comparator:     spec.ComparatorEqual,
							ParamIndex:     1,
						},
						{
							FieldReference: spec.FieldReference{ageField},
							Comparator:     spec.ComparatorEqual,
							ParamIndex:     2,
						},
					},
				},
			},
		},
	}
	expectedBytes, err := os.ReadFile("../../test/generator_test_expected.txt")
	if err != nil {
		t.Fatal(err)
	}
	expectedCode := string(expectedBytes)

	code, err := generator.GenerateRepository("user", userModel, "UserRepository", methods)

	if err != nil {
		t.Fatal(err)
	}
	if err := testutils.ExpectMultiLineString(expectedCode, code); err != nil {
		t.Error(err)
	}
}
