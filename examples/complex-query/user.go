package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

type UserModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Gender   Gender             `bson:"gender" json:"gender"`
	Age      int                `bson:"age" json:"age"`
	City     string             `bson:"city" json:"city"`
	Contact  *UserContactModel  `bson:"contact,omitempty" json:"contact"`
	Banned   bool               `bson:"banned" json:"banned"`
}

type UserContactModel struct {
	Phone string `bson:"phone" json:"phone"`
	Email string `bson:"email" json:"email"`
}

//go:generate repogen -dest=user_comparator_repo.go -model=UserModel -repo=UserComparatorRepository

// UserComparatorRepository is an interface that describes the specification of
// querying user data in the database.
type UserComparatorRepository interface {
	FindByUsername(ctx context.Context, username string) (*UserModel, error)
	FindByAgeGreaterThan(ctx context.Context, age int) ([]*UserModel, error)
	FindByAgeGreaterThanEqual(ctx context.Context, age int) ([]*UserModel, error)
	FindByAgeLessThan(ctx context.Context, age int) ([]*UserModel, error)
	FindByAgeLessThanEqual(ctx context.Context, age int) ([]*UserModel, error)
	FindByAgeBetween(ctx context.Context, fromAge int, toAge int) ([]*UserModel, error)
	FindByCityNot(ctx context.Context, city string) ([]*UserModel, error)
	FindByCityIn(ctx context.Context, cities []string) ([]*UserModel, error)
	FindByCityNotIn(ctx context.Context, cities []string) ([]*UserModel, error)
	FindByBannedTrue(ctx context.Context) ([]*UserModel, error)
	FindByBannedFalse(ctx context.Context) ([]*UserModel, error)
	FindByContactExists(ctx context.Context) (*UserModel, error)
	FindByContactNotExists(ctx context.Context) (*UserModel, error)
}

//go:generate repogen -dest=user_other_repo.go -model=UserModel -repo=UserOtherRepository

type UserOtherRepository interface {
	// FindByContactEmail demonstrates deeply-reference field (contect.email).
	FindByContactEmail(ctx context.Context, email string) (*UserModel, error)

	// FindByAgeAndCity demonstrates $and operation between two comparisons.
	FindByAgeAndCity(ctx context.Context, age int, city string) ([]*UserModel, error)

	// FindByGenderOrAgeGreaterThan demonstrates $or operation between two
	// comparisons.
	FindByGenderOrAgeGreaterThan(ctx context.Context, gender Gender, age int) ([]*UserModel, error)

	// FindTop5AllOrderByAge demonstrates limiting find many results
	FindTop5AllOrderByAge(ctx context.Context) ([]*UserModel, error)
}
