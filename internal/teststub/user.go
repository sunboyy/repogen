package teststub

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gender string

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	PhoneNumber    string             `bson:"phone_number"`
	Gender         Gender             `bson:"gender"`
	City           string             `bson:"city"`
	Age            int                `bson:"age"`
	Name           Name               `bson:"name"`
	Contact        Contact            `bson:"contact"`
	Referrer       *User              `bson:"referrer"`
	Enabled        bool               `bson:"enabled"`
	ConsentHistory []ConsentHistory   `bson:"consent_history"`
	AccessToken    string
}

type Name struct {
	First string `bson:"first"`
	Last  string `bson:"last"`
}

type Contact struct {
	Phone string
}

type ConsentHistory struct {
	ID    primitive.ObjectID
	Value bool
}

type UserRepositoryInsert interface {
	InsertMany(ctx context.Context, users []*User) ([]interface{}, error)
	InsertOne(ctx context.Context, user *User) (interface{}, error)
}

type UserRepositoryFind interface {
	// Test find all
	FindAll(ctx context.Context) ([]*User, error)
	// Test find with Between operator
	FindByAgeBetween(ctx context.Context, fromAge int, toAge int) ([]*User, error)
	// Test find with GreaterThan operator
	FindByAgeGreaterThan(ctx context.Context, age int) ([]*User, error)
	// Test find with GreaterThanEqual operator
	FindByAgeGreaterThanEqual(ctx context.Context, age int) ([]*User, error)
	// Test find with LessThan operator
	FindByAgeLessThan(ctx context.Context, age int) ([]*User, error)
	// Test find with LessThanEqual operator
	FindByAgeLessThanEqual(ctx context.Context, age int) ([]*User, error)
	// Test find MANY mode
	FindByCity(ctx context.Context, city string) ([]*User, error)
	// Test find with And operator
	FindByCityAndGender(ctx context.Context, city string, gender Gender) ([]*User, error)
	// Test find with In operator
	FindByCityIn(ctx context.Context, cities []string) ([]*User, error)
	// Test find with Not operator
	FindByCityNot(ctx context.Context, city string) ([]*User, error)
	// Test find with NotIn operator
	FindByCityNotIn(ctx context.Context, cities []string) ([]*User, error)
	// Test find with Or operator
	FindByCityOrGender(ctx context.Context, city string, gender Gender) ([]*User, error)
	// Test find ordering without explicit direction
	FindByCityOrderByAge(ctx context.Context, city string) ([]*User, error)
	// Test find ordering with explicit ascending direction
	FindByCityOrderByAgeAsc(ctx context.Context, city string) ([]*User, error)
	// Test find ordering with explicit descending direction
	FindByCityOrderByAgeDesc(ctx context.Context, city string) ([]*User, error)
	// Test find with multiple ordering
	FindByCityOrderByCityAndAgeDesc(ctx context.Context, city string) ([]*User, error)
	// Test find with deep reference ordering
	FindByCityOrderByNameFirst(ctx context.Context, city string) ([]*User, error)
	// Test find with False operator
	FindByEnabledFalse(ctx context.Context) ([]*User, error)
	// Test find with True operator
	FindByEnabledTrue(ctx context.Context) ([]*User, error)
	// Test find ONE mode
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	// Test find with deep referencing
	FindByNameFirst(ctx context.Context, firstName string) ([]*User, error)
	// Test find with multi-word arg
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	// Test find with Exists operator
	FindByReferrerExists(ctx context.Context) ([]*User, error)
	// Test find with deep pointer referencing
	FindByReferrerID(ctx context.Context, id primitive.ObjectID) ([]*User, error)
	// Test find with NotExists operator
	FindByReferrerNotExists(ctx context.Context) ([]*User, error)
	// Test find Top N
	FindTop5ByGenderOrderByAgeDesc(ctx context.Context, gender Gender) ([]*User, error)
}

type UserRepositoryUpdate interface {
	// Test update inc operator
	UpdateAgeIncByID(ctx context.Context, age int, id primitive.ObjectID) (bool, error)
	// Test update model ONE mode
	UpdateByID(ctx context.Context, user *User, id primitive.ObjectID) (bool, error)
	// Test update push operator
	UpdateConsentHistoryPushByID(ctx context.Context, consentHistoryItem ConsentHistory,
		id primitive.ObjectID) (int, error)
	// Test update multiple fields with push operator
	UpdateEnabledAndConsentHistoryPushByID(ctx context.Context, enabled bool,
		consentHistoryItem ConsentHistory, id primitive.ObjectID) (int, error)
	// Test update multiple fields
	UpdateGenderAndCityByID(ctx context.Context, gender Gender, city string, id primitive.ObjectID) (int, error)
	// Test update field MANY mode
	UpdateGenderByAge(ctx context.Context, gender Gender, age int) (int, error)
	// Test update field ONE mode
	UpdateGenderByID(ctx context.Context, gender Gender, id primitive.ObjectID) (bool, error)
	// Test update deep reference field
	UpdateNameFirstByID(ctx context.Context, firstName string, id primitive.ObjectID) (bool, error)
}

type UserRepositoryDelete interface {
	// Test delete all
	DeleteAll(ctx context.Context) (int, error)
	// Test delete with Between operator
	DeleteByAgeBetween(ctx context.Context, fromAge int, toAge int) (int, error)
	// Test delete with GreaterThan operator
	DeleteByAgeGreaterThan(ctx context.Context, age int) (int, error)
	// Test delete with GreaterThanEqual operator
	DeleteByAgeGreaterThanEqual(ctx context.Context, age int) (int, error)
	// Test delete with LessThan operator
	DeleteByAgeLessThan(ctx context.Context, age int) (int, error)
	// Test delete with LessThanEqual operator
	DeleteByAgeLessThanEqual(ctx context.Context, age int) (int, error)
	// Test delete MANY mode
	DeleteByCity(ctx context.Context, city string) (int, error)
	// Test delete with And operator
	DeleteByCityAndGender(ctx context.Context, city string, gender Gender) (int, error)
	// Test delete with In operator
	DeleteByCityIn(ctx context.Context, cities []string) (int, error)
	// Test delete with Not operator
	DeleteByCityNot(ctx context.Context, city string) (int, error)
	// Test delete with Or operator
	DeleteByCityOrGender(ctx context.Context, city string, gender Gender) (int, error)
	// Test delete ONE mode
	DeleteByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	// Test delete with deep reference
	DeleteByNameFirst(ctx context.Context, firstName string) (int, error)
	// Test delete multi-word arg
	DeleteByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
}

type UserRepositoryCount interface {
	// Test count all
	CountAll(ctx context.Context) (int, error)
	// Test count with query
	CountByGender(ctx context.Context, gender Gender) (int, error)
	// Test count with deep reference
	CountByNameFirst(ctx context.Context, firstName string) (int, error)
}

type UserRepositoryInvalidOperation interface {
	SearchByID(ctx context.Context, id primitive.ObjectID) (*User, error)
}

type UserRepositoryInvalidInsert interface {
	// Test insert with invalid number of returns
	Insert1(ctx context.Context, user *User) (*User, interface{}, error)
	// Test insert with invalid return type
	Insert2(ctx context.Context, user *User) (*User, error)
	// Test insert with unempty interface return
	Insert3(ctx context.Context, user *User) (interface{ Foo() }, error)
	// Test insert with no error return
	Insert4(ctx context.Context, user *User) (interface{}, bool)
	// Test insert with no context parameter
	Insert5(user *User) (interface{}, error)
	// Test insert with mismatched model parameter for ONE mode
	Insert6(ctx context.Context, user []*User) (interface{}, error)
	// Test insert with mismatched model parameter for MANY mode
	Insert7(ctx context.Context, user []*User) (interface{}, error)
}

type UserRepositoryInvalidFind interface {
	// Test find without query
	Find(ctx context.Context) ([]*User, error)
	// Test find with invalid number of returns
	FindAll(ctx context.Context) ([]*User, int, error)
	// Test find with misplaced sort operator token (rightmost)
	FindAllOrderByAgeAnd(ctx context.Context) ([]*User, error)
	// Test find with misplaced sort operator token (double operator)
	FindAllOrderByAgeAndAndGender(ctx context.Context) ([]*User, error)
	// Test find with misplaced sort operator token (leftmost)
	FindAllOrderByAndAge(ctx context.Context) ([]*User, error)
	// Test find with sort struct field not found
	FindAllOrderByCountry(ctx context.Context) ([]*User, error)
	// Test find with no context parameter
	FindByAge(age int) ([]*User, error)
	// Test find with misplaced query operator token (leftmost)
	FindByAndGender(ctx context.Context, gender Gender) ([]*User, error)
	// Test find with mismatched number of parameters
	FindByCity(ctx context.Context, city string, gender Gender) ([]*User, error)
	// Test find with mismatched parameter with In query
	FindByCityIn(ctx context.Context, city string) ([]*User, error)
	// Test find with query struct field not found
	FindByCountry(ctx context.Context, country string) ([]*User, error)
	// test find with mismatched parameter type
	FindByGender(ctx context.Context, gender string) ([]*User, error)
	// Test find with misplaced query operator token (rightmost)
	FindByGenderAnd(ctx context.Context, gender Gender) ([]*User, error)
	// Test find with misplaced query operator token (double operator)
	FindByGenderAndAndCity(ctx context.Context, gender Gender, city string) ([]*User, error)
	// Test find with ambiguous operator
	FindByGenderAndCityOrAge(ctx context.Context, gender Gender, city string, age int) ([]*User, error)
	// Test find with incompatible struct field for False comparator
	FindByGenderFalse(ctx context.Context) ([]*User, error)
	// Test find with incompatible struct field for True comparator
	FindByGenderTrue(ctx context.Context) ([]*User, error)
	// Test find with invalid return type
	FindByID(ctx context.Context, id primitive.ObjectID) (User, error)
	// Test find with deep reference field not found
	FindByNameMiddle(ctx context.Context, middleName string) ([]*User, error)
	// Test find top with no number and query
	FindTop(ctx context.Context) ([]*User, error)
	// Test find top 0
	FindTop0All(ctx context.Context) ([]*User, error)
	// Test find top in ONE mode
	FindTop5All(ctx context.Context) (*User, error)
	// Test find top with no number
	FindTopAll(ctx context.Context) ([]*User, error)
}

type UserRepositoryInvalidUpdate interface {
	// Test update with mismatched And token in update fields
	UpdateAgeAndAndGenderByID(ctx context.Context, age int, gender Gender,
		id primitive.ObjectID) (bool, error)
	// Test update without context parameter
	UpdateAgeByGender(age int, gender Gender) (int, error)
	// Test update with invalid number of returns
	UpdateAgeByID(ctx context.Context, age int, id primitive.ObjectID) (bool, int, error)
	// Test update with ambiguous query
	UpdateAgeByIDAndUsernameOrGender(ctx context.Context, age int, id primitive.ObjectID,
		username string, gender Gender) (bool, error)
	// Test update model with invalid parameter type
	UpdateByGender(ctx context.Context, gender Gender) (bool, error)
	// Test update with no update parameter provided
	UpdateByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	// Test update without query
	UpdateCity(ctx context.Context, city string) (bool, error)
	// Test update with invalid return type
	UpdateCityByID(ctx context.Context, city string, id primitive.ObjectID) (float64, error)
	// Test update with inc operator in non-number field
	UpdateCityIncByID(ctx context.Context, city string, id primitive.ObjectID) (bool, error)
	// Test update with push operator with incorrect parameter type
	UpdateConsentHistoryPushByID(ctx context.Context, consentHistoryItem []ConsentHistory,
		id primitive.ObjectID) (int, error)
	// Test update field not found in struct
	UpdateCountryByGender(ctx context.Context, country string, gender Gender) (int, error)
	// Test update with insufficient function parameters
	UpdateEnabledAll(ctx context.Context) (int, error)
	// Test update with incorrect parameter type for query
	UpdateEnabledByCity(ctx context.Context, enabled bool, city int) (bool, error)
	// Test update with incorrect parameter type for update field
	UpdateEnabledByGender(ctx context.Context, enabled int, gender Gender) (bool, error)
	// Test update with no error return
	UpdateEnabledByID(ctx context.Context, enabled bool, id primitive.ObjectID) (bool, bool)
	// Test update with push operator in non-array field
	UpdateGenderPushByID(ctx context.Context, gender Gender, id primitive.ObjectID) (bool, error)
}

type UserRepositoryInvalidDelete interface {
	// Test delete without query
	Delete(ctx context.Context) (int, error)
	// Test delete with invalid number of returns
	DeleteAll(ctx context.Context) (*User, int, error)
	// Test delete with unsupported return type
	DeleteByAge(ctx context.Context, age int) (float64, error)
	// Test delete with misplaced operator token (leftmost)
	DeleteByAndGender(ctx context.Context, gender Gender) (bool, error)
	// Test delete with no error return
	DeleteByCity(ctx context.Context, city string) (int, bool)
	// Test delete with mismatched parameter type for In operator
	DeleteByCityIn(ctx context.Context, city string) (int, error)
	// Test delete with query struct field not found
	DeleteByCountry(ctx context.Context, country string) (int, error)
	// Test delete with mismatched number of parameters
	DeleteByEnabled(ctx context.Context, enabled bool, enabled2 bool) (int, error)
	// Test delete without context parameter
	DeleteByGender(gender Gender) (int, error)
	// Test delete with misplaced operator token (rightmost)
	DeleteByGenderAnd(ctx context.Context, gender Gender) (bool, error)
	// Test delete with misplaced operator token (double operator)
	DeleteByGenderAndAndCity(ctx context.Context, gender Gender, city string) (bool, error)
	// Test delete with ambiguous query
	DeleteByGenderAndCityOrAge(ctx context.Context, gender Gender, city string, age int) (bool, error)
	// Test delete with mismatched parameter type
	DeleteByPhoneNumber(ctx context.Context, phoneNumber int) (bool, error)
}

type UserRepositoryInvalidCount interface {
	// Test count with query
	Count(ctx context.Context) (int, error)
	// Test count with invalid number of returns
	CountAll(ctx context.Context) (int, error, bool)
	// Test count with invalid query
	CountBy(ctx context.Context) (int, error)
	// Test count with invalid integer return
	CountByAge(ctx context.Context, age int) (int64, error)
	// Test count with no error return
	CountByCity(ctx context.Context, city string) (int, bool)
	// Test count with struct field not found
	CountByCountry(ctx context.Context, country string) (int, error)
	// Test count with mismatched number of parameters
	CountByEnabled(ctx context.Context, enabled bool, enabled2 bool) (int, error)
	// Test count without context parameter
	CountByGender(gender Gender) (int, error)
	// Test count with mismatched parameter type
	CountByPhoneNumber(ctx context.Context, phoneNumber int) (int, error)
}
