package mongo_test

import (
	"bytes"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
	"github.com/sunboyy/repogen/internal/testutils"
)

var userModel = code.Struct{
	Name: "UserModel",
	Fields: code.StructFields{
		{
			Name: "ID",
			Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"},
			Tags: map[string][]string{"bson": {"_id", "omitempty"}},
		},
		{
			Name: "Username",
			Type: code.SimpleType("string"),
			Tags: map[string][]string{"bson": {"username"}},
		},
		{
			Name: "Gender",
			Type: code.SimpleType("Gender"),
			Tags: map[string][]string{"bson": {"gender"}},
		},
		{
			Name: "Age",
			Type: code.SimpleType("int"),
			Tags: map[string][]string{"bson": {"age"}},
		},
		{
			Name: "AccessToken",
			Type: code.SimpleType("string"),
		},
	},
}

const expectedConstructorResult = `
import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
				Returns: []code.Type{code.PointerType{ContainedType: code.SimpleType("UserModel")}, code.SimpleType("error")},
				Operation: spec.FindOperation{
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Comparator: spec.ComparatorEqual, Field: "ID"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByID(ctx context.Context, arg0 primitive.ObjectID) (*UserModel, error) {
	var entity UserModel
	if err := r.collection.FindOne(ctx, bson.M{
		"_id": arg0,
	}).Decode(&entity); err != nil {
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
							{Comparator: spec.ComparatorEqual, Field: "Gender"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGender(ctx context.Context, arg0 Gender) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"gender": arg0,
	})
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, err
	}
	return entities, nil
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
							{Comparator: spec.ComparatorEqual, Field: "Gender"},
							{Comparator: spec.ComparatorEqual, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderAndAge(ctx context.Context, arg0 Gender, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"gender": arg0,
		"age": arg1,
	})
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(ctx, &entities); err != nil {
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
							{Comparator: spec.ComparatorEqual, Field: "Gender"},
							{Comparator: spec.ComparatorEqual, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderOrAge(ctx context.Context, arg0 Gender, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"gender": arg0},
			{"age": arg1},
		},
	})
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(ctx, &entities); err != nil {
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
							{Comparator: spec.ComparatorNot, Field: "Gender"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderNot(ctx context.Context, arg0 Gender) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"gender": bson.M{"$ne": arg0},
	})
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(ctx, &entities); err != nil {
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
							{Comparator: spec.ComparatorLessThan, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeLessThan(ctx context.Context, arg0 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"age": bson.M{"$lt": arg0},
	})
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(ctx, &entities); err != nil {
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
							{Comparator: spec.ComparatorLessThanEqual, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeLessThanEqual(ctx context.Context, arg0 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"age": bson.M{"$lte": arg0},
	})
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(ctx, &entities); err != nil {
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
							{Comparator: spec.ComparatorGreaterThan, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeGreaterThan(ctx context.Context, arg0 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"age": bson.M{"$gt": arg0},
	})
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(ctx, &entities); err != nil {
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
							{Comparator: spec.ComparatorGreaterThanEqual, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeGreaterThanEqual(ctx context.Context, arg0 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"age": bson.M{"$gte": arg0},
	})
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(ctx, &entities); err != nil {
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
							{Comparator: spec.ComparatorBetween, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeBetween(ctx context.Context, arg0 int, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"age": bson.M{"$gte": arg0, "$lte": arg1},
	})
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(ctx, &entities); err != nil {
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
							{Comparator: spec.ComparatorIn, Field: "Gender"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderIn(ctx context.Context, arg0 []Gender) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"gender": bson.M{"$in": arg0},
	})
	if err != nil {
		return nil, err
	}
	var entities []*UserModel
	if err := cursor.All(ctx, &entities); err != nil {
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
							{Comparator: spec.ComparatorEqual, Field: "ID"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByID(ctx context.Context, arg0 primitive.ObjectID) (bool, error) {
	result, err := r.collection.DeleteOne(ctx, bson.M{
		"_id": arg0,
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
							{Comparator: spec.ComparatorEqual, Field: "Gender"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGender(ctx context.Context, arg0 Gender) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"gender": arg0,
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
							{Comparator: spec.ComparatorEqual, Field: "Gender"},
							{Comparator: spec.ComparatorEqual, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGenderAndAge(ctx context.Context, arg0 Gender, arg1 int) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"gender": arg0,
		"age": arg1,
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
							{Comparator: spec.ComparatorEqual, Field: "Gender"},
							{Comparator: spec.ComparatorEqual, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGenderOrAge(ctx context.Context, arg0 Gender, arg1 int) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"$or": []bson.M{
			{"gender": arg0},
			{"age": arg1},
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
							{Comparator: spec.ComparatorNot, Field: "Gender"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGenderNot(ctx context.Context, arg0 Gender) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"gender": bson.M{"$ne": arg0},
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
							{Comparator: spec.ComparatorLessThan, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByAgeLessThan(ctx context.Context, arg0 int) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"age": bson.M{"$lt": arg0},
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
							{Comparator: spec.ComparatorLessThanEqual, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByAgeLessThanEqual(ctx context.Context, arg0 int) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"age": bson.M{"$lte": arg0},
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
							{Comparator: spec.ComparatorGreaterThan, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByAgeGreaterThan(ctx context.Context, arg0 int) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"age": bson.M{"$gt": arg0},
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
							{Comparator: spec.ComparatorGreaterThanEqual, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByAgeGreaterThanEqual(ctx context.Context, arg0 int) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"age": bson.M{"$gte": arg0},
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
							{Comparator: spec.ComparatorBetween, Field: "Age"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByAgeBetween(ctx context.Context, arg0 int, arg1 int) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"age": bson.M{"$gte": arg0, "$lte": arg1},
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
							{Comparator: spec.ComparatorIn, Field: "Gender"},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGenderIn(ctx context.Context, arg0 []Gender) (int, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"gender": bson.M{"$in": arg0},
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

type GenerateMethodInvalidTestCase struct {
	Name          string
	Method        spec.MethodSpec
	ExpectedError error
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
				Operation: "search",
			},
			ExpectedError: mongo.OperationNotSupportedError,
		},
		{
			Name: "bson tag not found",
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
							{Field: "AccessToken", Comparator: spec.ComparatorEqual},
						},
					},
				},
			},
			ExpectedError: mongo.BsonTagNotFoundError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			generator := mongo.NewGenerator(userModel, "UserRepository")
			buffer := new(bytes.Buffer)

			err := generator.GenerateMethod(testCase.Method, buffer)

			if err != testCase.ExpectedError {
				t.Errorf("\nExpected = %v\nReceived = %v", testCase.ExpectedError, err)
			}
		})
	}
}
