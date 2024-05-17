package teststub

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepositoryIntegration interface {
	FindAll(ctx context.Context) ([]*User, error)
	FindByAgeLessThanEqualOrderByAge(ctx context.Context, age int) ([]*User, error)
	FindByAgeGreaterThanEqualOrderByAgeDesc(ctx context.Context, age int) ([]*User, error)
	FindByAgeGreaterThanOrderByAgeAsc(ctx context.Context, age int) ([]*User, error)
	FindByAgeBetween(ctx context.Context, ageFrom int, ageTo int) ([]*User, error)
	FindByGenderNotAndAgeLessThan(ctx context.Context, gender Gender, age int) ([]*User, error)
	FindByGenderOrAge(ctx context.Context, gender Gender, age int) ([]*User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	InsertMany(ctx context.Context, users []*User) ([]interface{}, error)
	InsertOne(ctx context.Context, user *User) (interface{}, error)
}
