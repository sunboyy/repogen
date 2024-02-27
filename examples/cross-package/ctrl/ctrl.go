package ctrl

import (
	"context"

	"github.com/sunboyy/repogen/examples/cross_package"
)

//go:generate repogen -dest=../repo/user_repo.go -dest-pkg=repo -model-dir=../ -model=UserModel -repo=UserRepository

// UserRepository is an interface that describes the specification of querying
// user data in the database.
type UserRepository interface {
	// InsertOne stores userModel into the database and returns inserted ID
	// if insertion succeeds and returns error if insertion fails.
	InsertOne(ctx context.Context, userModel *cross_package.UserModel) (interface{}, error)
	// FindByUsername queries user by username. If a user with specified
	// username exists, the user will be returned. Otherwise, error will be
	// returned.
	FindByUsername(ctx context.Context, username string) (*cross_package.UserModel, error)

	// UpdateByID updates a single document by ID
	UpdateByID(ctx context.Context, user *cross_package.UserModel, id string) (bool, error)

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
