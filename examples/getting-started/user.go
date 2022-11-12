package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserModel is a model of user that is stored in the database
type UserModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username"`
	DisplayName string             `bson:"display_name" json:"displayName"`
	City        string             `bson:"city" json:"city"`
}

//go:generate repogen -src=user.go -dest=user_repo.go -model=UserModel -repo=UserRepository

// UserRepository is an interface that describes the specification of querying
// user data in the database.
type UserRepository interface {
	// InsertOne stores userModel into the database and returns inserted ID
	// if insertion succeeds and returns error if insertion fails.
	InsertOne(ctx context.Context, userModel *UserModel) (interface{}, error)

	// FindByUsername queries user by username. If a user with specified
	// username exists, the user will be returned. Otherwise, error will be
	// returned.
	FindByUsername(ctx context.Context, username string) (*UserModel, error)

	// UpdateDisplayNameByID updates a user with the specified ID with a new
	// display name. If there is a user matches the query, it will return
	// true. Error will be returned only when error occurs while accessing
	// the database.
	UpdateDisplayNameByID(ctx context.Context, displayName string, id primitive.ObjectID) (bool, error)

	// DeleteByCity deletes users that have `city` value match the parameter
	// and returns the match count. The error will be returned only when
	// error occurs while accessing the database. This is a MANY mode
	// because the first return type is an integer.
	DeleteByCity(ctx context.Context, city string) (int, error)

	// CountByCity returns the number of rows that match the given city
	// parameter. If an error occurs while accessing the database, error
	// value will be returned.
	CountByCity(ctx context.Context, city string) (int, error)
}
