package mongo_test

import (
	"strings"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/mongo"
)

func TestGenerateMongoRepository(t *testing.T) {
	userModel := code.Struct{
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
		},
	}
	intf := code.Interface{
		Name: "UserRepository",
		Methods: []code.Method{
			{
				Name: "FindByID",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
				},
				Returns: []code.Type{code.PointerType{ContainedType: code.SimpleType("UserModel")}, code.SimpleType("error")},
			},
			{
				Name: "FindOneByUsername",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "username", Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.SimpleType("error"),
				},
			},
			{
				Name: "FindByUsername",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "username", Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			{
				Name: "FindByIDAndUsername",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "id", Type: code.ExternalType{PackageAlias: "primitive", Name: "ObjectID"}},
					{Name: "username", Type: code.SimpleType("string")},
				},
				Returns: []code.Type{
					code.PointerType{ContainedType: code.SimpleType("UserModel")},
					code.SimpleType("error"),
				},
			},
			{
				Name: "FindByGenderNot",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			{
				Name: "FindByAgeLessThan",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			{
				Name: "FindByAgeLessThanEqual",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			{
				Name: "FindByAgeGreaterThan",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			{
				Name: "FindByAgeGreaterThanEqual",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			{
				Name: "FindByGenderOrAgeLessThan",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.SimpleType("Gender")},
					{Name: "age", Type: code.SimpleType("int")},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
			{
				Name: "FindByGenderIn",
				Params: []code.Param{
					{Name: "ctx", Type: code.ExternalType{PackageAlias: "context", Name: "Context"}},
					{Name: "gender", Type: code.ArrayType{ContainedType: code.SimpleType("Gender")}},
				},
				Returns: []code.Type{
					code.ArrayType{ContainedType: code.PointerType{ContainedType: code.SimpleType("UserModel")}},
					code.SimpleType("error"),
				},
			},
		},
	}

	code, err := mongo.GenerateMongoRepository("user", userModel, intf)

	if err != nil {
		t.Error(err)
	}
	expectedCode := `// Code generated by repogen. DO NOT EDIT.
package user

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

func (r *UserRepositoryMongo) FindByID(ctx context.Context, arg0 primitive.ObjectID) (*UserModel, error) {
	var entity UserModel
	if err := r.collection.FindOne(ctx, bson.M{
		"_id": arg0,
	}).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *UserRepositoryMongo) FindOneByUsername(ctx context.Context, arg0 string) (*UserModel, error) {
	var entity UserModel
	if err := r.collection.FindOne(ctx, bson.M{
		"username": arg0,
	}).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *UserRepositoryMongo) FindByUsername(ctx context.Context, arg0 string) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"username": arg0,
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

func (r *UserRepositoryMongo) FindByIDAndUsername(ctx context.Context, arg0 primitive.ObjectID, arg1 string) (*UserModel, error) {
	var entity UserModel
	if err := r.collection.FindOne(ctx, bson.M{
		"_id":      arg0,
		"username": arg1,
	}).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

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

func (r *UserRepositoryMongo) FindByGenderOrAgeLessThan(ctx context.Context, arg0 Gender, arg1 int) ([]*UserModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"gender": arg0},
			{"age": bson.M{"$lt": arg1}},
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
`
	expectedCodeLines := strings.Split(expectedCode, "\n")
	actualCodeLines := strings.Split(code, "\n")

	for i, line := range expectedCodeLines {
		if line != actualCodeLines[i] {
			t.Errorf("On line %d\nExpected = %v\nActual = %v", i, line, actualCodeLines[i])
		}
	}
}
