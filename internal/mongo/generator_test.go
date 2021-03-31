package mongo_test

import (
	"bytes"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
	"github.com/sunboyy/repogen/internal/testutils"
)

var (
	idField = code.StructField{
		Name: "ID",
		Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
		Tags: map[string][]string{"bson": {"_id", "omitempty"}},
	}
	genderField = code.StructField{
		Name: "Gender",
		Type: code.SimpleType("Gender"),
		Tags: map[string][]string{"bson": {"gender"}},
	}
	ageField = code.StructField{
		Name: "Age",
		Type: code.SimpleType("int"),
		Tags: map[string][]string{"bson": {"age"}},
	}
	nameField = code.StructField{
		Name: "Name",
		Type: code.SimpleType("NameModel"),
		Tags: map[string][]string{"bson": {"name"}},
	}
	consentHistoryField = code.StructField{
		Name: "ConsentHistory",
		Type: code.ArrayType{ContainedType: code.SimpleType("ConsentHistory")},
		Tags: map[string][]string{"bson": {"consent_history"}},
	}
	enabledField = code.StructField{
		Name: "Enabled",
		Type: code.SimpleType("bool"),
		Tags: map[string][]string{"bson": {"enabled"}},
	}
	accessTokenField = code.StructField{
		Name: "AccessToken",
		Type: code.SimpleType("string"),
	}

	firstNameField = code.StructField{
		Name: "First",
		Type: code.SimpleType("string"),
		Tags: map[string][]string{"bson": {"first"}},
	}
)

var userModel = code.Struct{
	Name: "UserModel",
	Fields: code.StructFields{
		idField,
		{
			Name: "Username",
			Type: code.SimpleType("string"),
			Tags: map[string][]string{"bson": {"username"}},
		},
		genderField,
		ageField,
		nameField,
		consentHistoryField,
		enabledField,
		accessTokenField,
	},
}

const expectedConstructorResult = `
import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &UserRepositoryMongo{
		collection: collection,
	}
}

type UserRepositoryMongo struct {
	collection *mongo.Collection
}
`

func TestGenerateConstructor(t *testing.T) {
	generator := mongo.NewGenerator(userModel, "UserRepository")
	buffer := new(bytes.Buffer)

	err := generator.GenerateConstructor(buffer)

	if err != nil {
		t.Error(err)
	}
	if err := testutils.ExpectMultiLineString(expectedConstructorResult, buffer.String()); err != nil {
		t.Error(err)
	}
}

type GenerateMethodTestCase struct {
	Name         string
	MethodSpec   spec.MethodSpec
	ExpectedCode string
}

func TestGenerateMethod_Insert(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "insert one method",
			MethodSpec: spec.MethodSpec{
				Name: "InsertOne",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "userModel", Type: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
				},
				Returns: []code.Type{
					code.InterfaceType{},
					code.SimpleType("error"),
				},
				Operation: spec.InsertOperation{
					Mode: spec.QueryModeOne,
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) InsertOne(arg0 context.Context, arg1 *UserModel) (interface{}, error) {
	result, err := r.collection.InsertOne(arg0, arg1)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}
`,
		},
		{
			Name: "insert many method",
			MethodSpec: spec.MethodSpec{
				Name: "Insert",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "userModel", Type: code.ArrayType{
						ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")},
					}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.InterfaceType{}},
					code.SimpleType("error"),
				},
				Operation: spec.InsertOperation{
					Mode: spec.QueryModeMany,
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) Insert(arg0 context.Context, arg1 []*UserModel) ([]interface{}, error) {
	var entities []interface{}
	for _, model := range arg1 {
		entities = append(entities, model)
	}
	result, err := r.collection.InsertMany(arg0, entities)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil
}
`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			buffer := new(bytes.Buffer)

			err := generator.GenerateMethod(testCase.MethodSpec, buffer)

			if err != nil {
				t.Error(err)
			}
			if err := testutils.ExpectMultiLineString(testCase.ExpectedCode, buffer.String()); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGenerateMethod_Find(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "simple find one method",
			MethodSpec: spec.MethodSpec{
				Name: "FindByID",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{idField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByID(arg0 context.Context, arg1 primitive.ObjectID) (*UserModel, error) {
	var entity UserModel
	if err := r.collection.FindOne(arg0, bson.M{
		"_id": arg1,
	}, options.FindOne().SetSort(bson.M{
	})).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
`,
		},
		{
			Name: "simple find many method",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGender",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGender(arg0 context.Context, arg1 Gender) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": arg1,
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with deep field reference",
			MethodSpec: spec.MethodSpec{
				Name: "FindByNameFirst",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "firstName", Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{
								Comparator:     spec.ComparatorEqual,
								FieldReference: spec.FieldReference{nameField, firstNameField},
								ParamIndex:     1,
							},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByNameFirst(arg0 context.Context, arg1 string) (*UserModel, error) {
	var entity UserModel
	if err := r.collection.FindOne(arg0, bson.M{
		"name.first": arg1,
	}, options.FindOne().SetSort(bson.M{
	})).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
`,
		},
		{
			Name: "find with And operator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderAndAge",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Operator: spec.OperatorAnd,
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{ageField}, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderAndAge(arg0 context.Context, arg1 Gender, arg2 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"$and": []bson.M{
			{"gender": arg1},
			{"age": arg2},
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with Or operator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderOrAge",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Operator: spec.OperatorOr,
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{ageField}, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderOrAge(arg0 context.Context, arg1 Gender, arg2 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"$or": []bson.M{
			{"gender": arg1},
			{"age": arg2},
		},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with Not comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderNot",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorNot, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderNot(arg0 context.Context, arg1 Gender) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": bson.M{"$ne": arg1},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with LessThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeLessThan",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorLessThan, FieldReference: spec.FieldReference{ageField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeLessThan(arg0 context.Context, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{"$lt": arg1},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with LessThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeLessThanEqual",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorLessThanEqual, FieldReference: spec.FieldReference{ageField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeLessThanEqual(arg0 context.Context, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{"$lte": arg1},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with GreaterThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeGreaterThan",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorGreaterThan, FieldReference: spec.FieldReference{ageField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeGreaterThan(arg0 context.Context, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{"$gt": arg1},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with GreaterThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeGreaterThanEqual",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorGreaterThanEqual, FieldReference: spec.FieldReference{ageField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeGreaterThanEqual(arg0 context.Context, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{"$gte": arg1},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with Between comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByAgeBetween",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "fromAge", Type: code.SimpleType("int")},
					{Name: "toAge", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorBetween, FieldReference: spec.FieldReference{ageField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeBetween(arg0 context.Context, arg1 int, arg2 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{"$gte": arg1, "$lte": arg2},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with In comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderIn",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.ArrayType{ContainedType: code.SimpleType("Gender")}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorIn, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderIn(arg0 context.Context, arg1 []Gender) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": bson.M{"$in": arg1},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with NotIn comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByGenderNotIn",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.ArrayType{ContainedType: code.SimpleType("Gender")}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorNotIn, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderNotIn(arg0 context.Context, arg1 []Gender) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": bson.M{"$nin": arg1},
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with True comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByEnabledTrue",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorTrue, FieldReference: spec.FieldReference{enabledField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByEnabledTrue(arg0 context.Context) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"enabled": true,
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with False comparator",
			MethodSpec: spec.MethodSpec{
				Name: "FindByEnabledFalse",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorFalse, FieldReference: spec.FieldReference{enabledField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByEnabledFalse(arg0 context.Context) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"enabled": false,
	}, options.Find().SetSort(bson.M{
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with sort ascending",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByAge",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingAscending},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindAllOrderByAge(arg0 context.Context) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{

	}, options.Find().SetSort(bson.M{
		"age": 1,
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with sort descending",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByAgeDesc",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingDescending},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindAllOrderByAgeDesc(arg0 context.Context) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{

	}, options.Find().SetSort(bson.M{
		"age": -1,
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with deep sort ascending",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByNameFirst",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{FieldReference: spec.FieldReference{nameField, firstNameField}, Ordering: spec.OrderingAscending},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindAllOrderByNameFirst(arg0 context.Context) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{

	}, options.Find().SetSort(bson.M{
		"name.first": 1,
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
		{
			Name: "find with multiple sorts",
			MethodSpec: spec.MethodSpec{
				Name: "FindAllOrderByGenderAndAgeDesc",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeMany,
					Sorts: []spec.Sort{
						{FieldReference: spec.FieldReference{genderField}, Ordering: spec.OrderingAscending},
						{FieldReference: spec.FieldReference{ageField}, Ordering: spec.OrderingDescending},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindAllOrderByGenderAndAgeDesc(arg0 context.Context) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{

	}, options.Find().SetSort(bson.M{
		"gender": 1,
		"age": -1,
	}))
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(arg0, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			buffer := new(bytes.Buffer)

			err := generator.GenerateMethod(testCase.MethodSpec, buffer)

			if err != nil {
				t.Error(err)
			}
			if err := testutils.ExpectMultiLineString(testCase.ExpectedCode, buffer.String()); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGenerateMethod_Update(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "update model method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateByID",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "model", Type: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateModel{},
					Mode:   spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) UpdateByID(arg0 context.Context, arg1 *UserModel, arg2 primitive.ObjectID) (bool, error) {
	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg2,
	}, bson.M{
		"$set": arg1,
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, err
}
`,
		},
		{
			Name: "simple update one method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateAgeByID",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						{FieldReference: spec.FieldReference{ageField}, ParamIndex: 1, Operator: spec.UpdateOperatorSet},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) UpdateAgeByID(arg0 context.Context, arg1 int, arg2 primitive.ObjectID) (bool, error) {
	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg2,
	}, bson.M{
		"$set": bson.M{
			"age": arg1,
		},
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, err
}
`,
		},
		{
			Name: "simple update many method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateAgeByGender",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						{FieldReference: spec.FieldReference{ageField}, ParamIndex: 1, Operator: spec.UpdateOperatorSet},
					},
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{genderField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) UpdateAgeByGender(arg0 context.Context, arg1 int, arg2 Gender) (int, error) {
	result, err := r.collection.UpdateMany(arg0, bson.M{
		"gender": arg2,
	}, bson.M{
		"$set": bson.M{
			"age": arg1,
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.MatchedCount), err
}
`,
		},
		{
			Name: "simple update push method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateConsentHistoryPushByID",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "consentHistory", Type: code.SimpleType("ConsentHistory")},
					{Name: "gender", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						{FieldReference: spec.FieldReference{consentHistoryField}, ParamIndex: 1, Operator: spec.UpdateOperatorPush},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) UpdateConsentHistoryPushByID(arg0 context.Context, arg1 ConsentHistory, arg2 primitive.ObjectID) (bool, error) {
	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg2,
	}, bson.M{
		"$push": bson.M{
			"consent_history": arg1,
		},
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, err
}
`,
		},
		{
			Name: "simple update set and push method",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateEnabledAndConsentHistoryPushByID",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "enabled", Type: code.SimpleType("bool")},
					{Name: "consentHistory", Type: code.SimpleType("ConsentHistory")},
					{Name: "gender", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						{FieldReference: spec.FieldReference{enabledField}, ParamIndex: 1, Operator: spec.UpdateOperatorSet},
						{FieldReference: spec.FieldReference{consentHistoryField}, ParamIndex: 2, Operator: spec.UpdateOperatorPush},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 3},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) UpdateEnabledAndConsentHistoryPushByID(arg0 context.Context, arg1 bool, arg2 ConsentHistory, arg3 primitive.ObjectID) (bool, error) {
	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg3,
	}, bson.M{
		"$set": bson.M{
			"enabled": arg1,
		},
		"$push": bson.M{
			"consent_history": arg2,
		},
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, err
}
`,
		},
		{
			Name: "update with deeply referenced field",
			MethodSpec: spec.MethodSpec{
				Name: "UpdateNameFirstByID",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "firstName", Type: code.SimpleType("string")},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						{FieldReference: spec.FieldReference{nameField, firstNameField}, ParamIndex: 1, Operator: spec.UpdateOperatorSet},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) UpdateNameFirstByID(arg0 context.Context, arg1 string, arg2 primitive.ObjectID) (bool, error) {
	result, err := r.collection.UpdateOne(arg0, bson.M{
		"_id": arg2,
	}, bson.M{
		"$set": bson.M{
			"name.first": arg1,
		},
	})
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, err
}
`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			buffer := new(bytes.Buffer)

			err := generator.GenerateMethod(testCase.MethodSpec, buffer)

			if err != nil {
				t.Error(err)
			}
			if err := testutils.ExpectMultiLineString(testCase.ExpectedCode, buffer.String()); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGenerateMethod_Delete(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "simple delete one method",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByID",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{code.SimpleType("bool"), code.SimpleType("error")},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{idField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByID(arg0 context.Context, arg1 primitive.ObjectID) (bool, error) {
	result, err := r.collection.DeleteOne(arg0, bson.M{
		"_id": arg1,
	})
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}
`,
		},
		{
			Name: "simple delete many method",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByGender",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGender(arg0 context.Context, arg1 Gender) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"gender": arg1,
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
`,
		},
		{
			Name: "delete with And operator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByGenderAndAge",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Operator: spec.OperatorAnd,
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{ageField}, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGenderAndAge(arg0 context.Context, arg1 Gender, arg2 int) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"$and": []bson.M{
			{"gender": arg1},
			{"age": arg2},
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
`,
		},
		{
			Name: "delete with Or operator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByGenderOrAge",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Operator: spec.OperatorOr,
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
							{Comparator: spec.ComparatorEqual, FieldReference: spec.FieldReference{ageField}, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGenderOrAge(arg0 context.Context, arg1 Gender, arg2 int) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"$or": []bson.M{
			{"gender": arg1},
			{"age": arg2},
		},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
`,
		},
		{
			Name: "delete with Not comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByGenderNot",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorNot, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGenderNot(arg0 context.Context, arg1 Gender) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"gender": bson.M{"$ne": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
`,
		},
		{
			Name: "delete with LessThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByAgeLessThan",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorLessThan, FieldReference: spec.FieldReference{ageField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByAgeLessThan(arg0 context.Context, arg1 int) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"age": bson.M{"$lt": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
`,
		},
		{
			Name: "delete with LessThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByAgeLessThanEqual",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorLessThanEqual, FieldReference: spec.FieldReference{ageField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByAgeLessThanEqual(arg0 context.Context, arg1 int) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"age": bson.M{"$lte": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
`,
		},
		{
			Name: "delete with GreaterThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByAgeGreaterThan",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorGreaterThan, FieldReference: spec.FieldReference{ageField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByAgeGreaterThan(arg0 context.Context, arg1 int) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"age": bson.M{"$gt": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
`,
		},
		{
			Name: "delete with GreaterThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByAgeGreaterThanEqual",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorGreaterThanEqual, FieldReference: spec.FieldReference{ageField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByAgeGreaterThanEqual(arg0 context.Context, arg1 int) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"age": bson.M{"$gte": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
`,
		},
		{
			Name: "delete with Between comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByAgeBetween",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "fromAge", Type: code.SimpleType("int")},
					{Name: "toAge", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorBetween, FieldReference: spec.FieldReference{ageField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByAgeBetween(arg0 context.Context, arg1 int, arg2 int) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"age": bson.M{"$gte": arg1, "$lte": arg2},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
`,
		},
		{
			Name: "delete with In comparator",
			MethodSpec: spec.MethodSpec{
				Name: "DeleteByGenderIn",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.ArrayType{ContainedType: code.SimpleType("Gender")}},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.DeleteOperation{
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorIn, FieldReference: spec.FieldReference{genderField}, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGenderIn(arg0 context.Context, arg1 []Gender) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"gender": bson.M{"$in": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}
`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			buffer := new(bytes.Buffer)

			err := generator.GenerateMethod(testCase.MethodSpec, buffer)

			if err != nil {
				t.Error(err)
			}
			if err := testutils.ExpectMultiLineString(testCase.ExpectedCode, buffer.String()); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGenerateMethod_Count(t *testing.T) {
	testTable := []GenerateMethodTestCase{
		{
			Name: "simple count method",
			MethodSpec: spec.MethodSpec{
				Name: "CountByGender",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{genderField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) CountByGender(arg0 context.Context, arg1 Gender) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"gender": arg1,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
`,
		},
		{
			Name: "count with And operator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByGenderAndCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Operator: spec.OperatorAnd,
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{genderField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
							{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) CountByGenderAndCity(arg0 context.Context, arg1 Gender, arg2 int) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"$and": []bson.M{
			{"gender": arg1},
			{"age": arg2},
		},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
`,
		},
		{
			Name: "count with Or operator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByGenderOrCity",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Operator: spec.OperatorOr,
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{genderField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
							{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) CountByGenderOrCity(arg0 context.Context, arg1 Gender, arg2 int) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"$or": []bson.M{
			{"gender": arg1},
			{"age": arg2},
		},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
`,
		},
		{
			Name: "count with Not comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByGenderNot",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{genderField}, Comparator: spec.ComparatorNot, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) CountByGenderNot(arg0 context.Context, arg1 Gender) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"gender": bson.M{"$ne": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
`,
		},
		{
			Name: "count with LessThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeLessThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorLessThan, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) CountByAgeLessThan(arg0 context.Context, arg1 int) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{"$lt": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
`,
		},
		{
			Name: "count with LessThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeLessThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorLessThanEqual, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) CountByAgeLessThanEqual(arg0 context.Context, arg1 int) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{"$lte": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
`,
		},
		{
			Name: "count with GreaterThan comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeGreaterThan",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorGreaterThan, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) CountByAgeGreaterThan(arg0 context.Context, arg1 int) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{"$gt": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
`,
		},
		{
			Name: "count with GreaterThanEqual comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeGreaterThanEqual",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorGreaterThanEqual, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) CountByAgeGreaterThanEqual(arg0 context.Context, arg1 int) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{"$gte": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
`,
		},
		{
			Name: "count with Between comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeBetween",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
					{Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorBetween, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) CountByAgeBetween(arg0 context.Context, arg1 int, arg2 int) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{"$gte": arg1, "$lte": arg2},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
`,
		},
		{
			Name: "count with In comparator",
			MethodSpec: spec.MethodSpec{
				Name: "CountByAgeIn",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ArrayType{ContainedType: code.SimpleType("int")}},
				},
				Returns: []code.Type{
					code.SimpleType("int"),
					code.SimpleType("error"),
				},
				Operation: spec.CountOperation{
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{ageField}, Comparator: spec.ComparatorIn, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) CountByAgeIn(arg0 context.Context, arg1 []int) (int, error) {
	count, err := r.collection.CountDocuments(arg0, bson.M{
		"age": bson.M{"$in": arg1},
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			buffer := new(bytes.Buffer)

			err := generator.GenerateMethod(testCase.MethodSpec, buffer)

			if err != nil {
				t.Error(err)
			}
			if err := testutils.ExpectMultiLineString(testCase.ExpectedCode, buffer.String()); err != nil {
				t.Error(err)
			}
		})
	}
}

type GenerateMethodInvalidTestCase struct {
	Name          string
	Method        spec.MethodSpec
	ExpectedError error
}

type StubOperation struct {
}

func (o StubOperation) Name() string {
	return "Stub"
}

func TestGenerateMethod_Invalid(t *testing.T) {
	testTable := []GenerateMethodInvalidTestCase{
		{
			Name: "operation not supported",
			Method: spec.MethodSpec{
				Name: "SearchByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: StubOperation{},
			},
			ExpectedError: mongo.NewOperationNotSupportedError("Stub"),
		},
		{
			Name: "bson tag not found in query",
			Method: spec.MethodSpec{
				Name: "FindByAccessToken",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{accessTokenField}, Comparator: spec.ComparatorEqual, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedError: mongo.NewBsonTagNotFoundError("AccessToken"),
		},
		{
			Name: "bson tag not found in sort",
			Method: spec.MethodSpec{
				Name: "FindAllOrderByAccessToken",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Sorts: []spec.Sort{
						{FieldReference: spec.FieldReference{accessTokenField}, Ordering: spec.OrderingAscending},
					},
				},
			},
			ExpectedError: mongo.NewBsonTagNotFoundError("AccessToken"),
		},
		{
			Name: "bson tag not found in update field",
			Method: spec.MethodSpec{
				Name: "UpdateAccessTokenByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("string")},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						{FieldReference: spec.FieldReference{accessTokenField}, ParamIndex: 1, Operator: spec.UpdateOperatorSet},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedError: mongo.NewBsonTagNotFoundError("AccessToken"),
		},
		{
			Name: "update type not supported",
			Method: spec.MethodSpec{
				Name: "UpdateAgeByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
				Operation: spec.UpdateOperation{
					Update: StubUpdate{},
					Mode:   spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedError: mongo.NewUpdateTypeNotSupportedError(StubUpdate{}),
		},
		{
			Name: "update operator not supported",
			Method: spec.MethodSpec{
				Name: "UpdateConsentHistoryAppendByID",
				Params: []code.Param{
					{Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Type: code.SimpleType("int")},
					{Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{
					code.SimpleType("bool"),
					code.SimpleType("error"),
				},
				Operation: spec.UpdateOperation{
					Update: spec.UpdateFields{
						{FieldReference: spec.FieldReference{consentHistoryField}, ParamIndex: 1, Operator: "APPEND"},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{FieldReference: spec.FieldReference{idField}, Comparator: spec.ComparatorEqual, ParamIndex: 2},
						},
					},
				},
			},
			ExpectedError: mongo.NewUpdateOperatorNotSupportedError("APPEND"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			buffer := new(bytes.Buffer)

			err := generator.GenerateMethod(testCase.Method, buffer)

			if err != testCase.ExpectedError {
				t.Errorf("\nExpected = %+v\nReceived = %+v", testCase.ExpectedError, err)
			}
		})
	}
}
