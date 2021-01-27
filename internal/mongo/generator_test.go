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
							{Comparator: spec.ComparatorEqual, Field: "ID", ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByID(arg0 context.Context, arg1 primitive.ObjectID) (*UserModel, error) {
	var entity UserModel
	if err := r.collection.FindOne(arg0, bson.M{
		"_id": arg1,
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
							{Comparator: spec.ComparatorEqual, Field: "Gender", ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGender(arg0 context.Context, arg1 Gender) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": arg1,
	})
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
							{Comparator: spec.ComparatorEqual, Field: "Gender", ParamIndex: 1},
							{Comparator: spec.ComparatorEqual, Field: "Age", ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderAndAge(arg0 context.Context, arg1 Gender, arg2 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": arg1,
		"age": arg2,
	})
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
							{Comparator: spec.ComparatorEqual, Field: "Gender", ParamIndex: 1},
							{Comparator: spec.ComparatorEqual, Field: "Age", ParamIndex: 2},
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
	})
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
							{Comparator: spec.ComparatorNot, Field: "Gender", ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderNot(arg0 context.Context, arg1 Gender) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": bson.M{"$ne": arg1},
	})
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
							{Comparator: spec.ComparatorLessThan, Field: "Age", ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeLessThan(arg0 context.Context, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{"$lt": arg1},
	})
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
							{Comparator: spec.ComparatorLessThanEqual, Field: "Age", ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeLessThanEqual(arg0 context.Context, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{"$lte": arg1},
	})
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
							{Comparator: spec.ComparatorGreaterThan, Field: "Age", ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeGreaterThan(arg0 context.Context, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{"$gt": arg1},
	})
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
							{Comparator: spec.ComparatorGreaterThanEqual, Field: "Age", ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeGreaterThanEqual(arg0 context.Context, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{"$gte": arg1},
	})
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
							{Comparator: spec.ComparatorBetween, Field: "Age", ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByAgeBetween(arg0 context.Context, arg1 int, arg2 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"age": bson.M{"$gte": arg1, "$lte": arg2},
	})
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
							{Comparator: spec.ComparatorIn, Field: "Gender", ParamIndex: 1},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) FindByGenderIn(arg0 context.Context, arg1 []Gender) ([]*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"gender": bson.M{"$in": arg1},
	})
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
					Fields: []spec.UpdateField{
						{Name: "Age", ParamIndex: 1},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Field: "ID", Comparator: spec.ComparatorEqual, ParamIndex: 2},
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
					Fields: []spec.UpdateField{
						{Name: "Age", ParamIndex: 1},
					},
					Mode: spec.QueryModeMany,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Field: "Gender", Comparator: spec.ComparatorEqual, ParamIndex: 2},
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
							{Comparator: spec.ComparatorEqual, Field: "ID", ParamIndex: 1},
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
							{Comparator: spec.ComparatorEqual, Field: "Gender", ParamIndex: 1},
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
							{Comparator: spec.ComparatorEqual, Field: "Gender", ParamIndex: 1},
							{Comparator: spec.ComparatorEqual, Field: "Age", ParamIndex: 2},
						},
					},
				},
			},
			ExpectedCode: `
func (r *UserRepositoryMongo) DeleteByGenderAndAge(arg0 context.Context, arg1 Gender, arg2 int) (int, error) {
	result, err := r.collection.DeleteMany(arg0, bson.M{
		"gender": arg1,
		"age": arg2,
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
							{Comparator: spec.ComparatorEqual, Field: "Gender", ParamIndex: 1},
							{Comparator: spec.ComparatorEqual, Field: "Age", ParamIndex: 2},
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
							{Comparator: spec.ComparatorNot, Field: "Gender", ParamIndex: 1},
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
							{Comparator: spec.ComparatorLessThan, Field: "Age", ParamIndex: 1},
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
							{Comparator: spec.ComparatorLessThanEqual, Field: "Age", ParamIndex: 1},
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
							{Comparator: spec.ComparatorGreaterThan, Field: "Age", ParamIndex: 1},
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
							{Comparator: spec.ComparatorGreaterThanEqual, Field: "Age", ParamIndex: 1},
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
							{Comparator: spec.ComparatorBetween, Field: "Age", ParamIndex: 1},
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
							{Comparator: spec.ComparatorIn, Field: "Gender", ParamIndex: 1},
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
							{Field: "AccessToken", Comparator: spec.ComparatorEqual, ParamIndex: 1},
						},
					},
				},
			},
			ExpectedError: mongo.BsonTagNotFoundError,
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
					Fields: []spec.UpdateField{
						{Name: "AccessToken", ParamIndex: 1},
					},
					Mode: spec.QueryModeOne,
					Query: spec.QuerySpec{
						Predicates: []spec.Predicate{
							{Field: "ID", Comparator: spec.ComparatorEqual, ParamIndex: 2},
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
