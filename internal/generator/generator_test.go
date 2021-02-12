package generator_test

import (
	"strings"
	"testing"

	"github.com/sunboyy/repogen/internal/code"
	"github.com/sunboyy/repogen/internal/generator"
	"github.com/sunboyy/repogen/internal/spec"
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
	methods := []spec.MethodSpec{
		// test find: One mode
		{
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
						{Field: "ID", Comparator: spec.ComparatorEqual, ParamIndex: 1},
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
				{Name: "age", Type: code.SimpleType("int")},
			},
			Returns: []code.Type{
				code.PointerType{ContainedType: code.SimpleType("UserModel")},
				code.SimpleType("error"),
			},
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Operator: spec.OperatorAnd,
					Predicates: []spec.Predicate{
						{Field: "Gender", Comparator: spec.ComparatorNot, ParamIndex: 1},
						{Field: "Age", Comparator: spec.ComparatorLessThan, ParamIndex: 2},
					},
				},
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
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Predicates: []spec.Predicate{
						{Field: "Age", Comparator: spec.ComparatorLessThanEqual, ParamIndex: 1},
					},
				},
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
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Predicates: []spec.Predicate{
						{Field: "Age", Comparator: spec.ComparatorGreaterThan, ParamIndex: 1},
					},
				},
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
			Operation: spec.FindOperation{
				Mode: spec.QueryModeMany,
				Query: spec.QuerySpec{
					Predicates: []spec.Predicate{
						{Field: "Age", Comparator: spec.ComparatorGreaterThanEqual, ParamIndex: 1},
					},
				},
			},
		},
		{
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
						{Field: "Age", Comparator: spec.ComparatorBetween, ParamIndex: 1},
					},
				},
			},
		},
		{
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
						{Field: "Gender", Comparator: spec.ComparatorEqual, ParamIndex: 1},
						{Field: "Age", Comparator: spec.ComparatorEqual, ParamIndex: 2},
					},
				},
			},
		},
	}

	code, err := generator.GenerateRepository("user", userModel, "UserRepository", methods)

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

func (r *UserRepositoryMongo) FindByID(arg0 context.Context, arg1 primitive.ObjectID) (*UserModel, error) {
	var entity UserModel
	if err := r.collection.FindOne(arg0, bson.M{
		"_id": arg1,
	}).Decode(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *UserRepositoryMongo) FindByGenderNotAndAgeLessThan(arg0 context.Context, arg1 Gender, arg2 int) (*UserModel, error) {
	cursor, err := r.collection.Find(arg0, bson.M{
		"$and": []bson.M{
			{"gender": bson.M{"$ne": arg1}},
			{"age": bson.M{"$lt": arg2}},
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
`
	expectedCodeLines := strings.Split(expectedCode, "\n")
	actualCodeLines := strings.Split(code, "\n")

	for i, line := range expectedCodeLines {
		if line != actualCodeLines[i] {
			t.Errorf("On line %d\nExpected = %v\nActual = %v", i, line, actualCodeLines[i])
		}
	}
}
