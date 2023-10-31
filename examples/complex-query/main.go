package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace these values with your own connection option. This connection option is hard-coded for easy
// demonstration. Make sure not to hard-code the credentials in the production code.
const (
	connectionString = "mongodb://admin:password@localhost:27017"
	databaseName     = "repogen_examples"
	collectionName   = "complexquery_user"
)

var (
	userComparatorRepository UserComparatorRepository
	userOtherRepository      UserOtherRepository
)

func init() {
	// create a connection to the database
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		panic(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		panic(err)
	}

	// instantiate a user repository from the user collection
	collection := client.Database(databaseName).Collection(collectionName)
	userComparatorRepository = NewUserComparatorRepository(collection)
	userOtherRepository = NewUserOtherRepository(collection)
}

func main() {
	demonstrateFindByExists()
	demonstrateFindByNotExists()
}

// demonstrateFindByExists shows how find method in repogen works. It receives
// query parameters through method arguments and returns matched result
func demonstrateFindByExists() {
	users, err := userComparatorRepository.FindByContactExists(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("FindByExists: found users = %+v\n", users)
}

// demonstrateFindByNotExists shows how find method in repogen works. It
// receives query parameters through method arguments and returns matched
// result.
func demonstrateFindByNotExists() {
	users, err := userComparatorRepository.FindByContactNotExists(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("FindByNotExists: found users = %+v\n", users)
}
