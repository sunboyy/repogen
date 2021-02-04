package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace these values with your own connection option. This connection option is hard-coded for easy
// demonstration. Make sure not to hard-code the credentials in the production code.
const (
	connectionString = "mongodb://lineman:lineman@localhost:27017"
	databaseName     = "repogen_examples"
	collectionName   = "gettingstarted_user"
)

var (
	userRepository UserRepository
	userID         primitive.ObjectID
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
	userRepository = NewUserRepository(client.Database(databaseName).Collection(collectionName))
}

func main() {
	demonstrateInsertion()
	demonstrateFind()
	demonstrateUpdate()
	demonstrateDelete()
}

// demonstrateInsertion shows how insert method of repogen works. It receives a model struct and returns
// an insertedID.
func demonstrateInsertion() {
	userModel := &UserModel{
		Username:    "sunboyy",
		DisplayName: "SB",
		City:        "Bangkok, Thailand",
	}

	insertedID, err := userRepository.InsertOne(context.Background(), userModel)
	if err != nil {
		panic(err)
	}
	userID = insertedID.(primitive.ObjectID)
	fmt.Printf("Insert (one): inserted id = %v\n", insertedID)
}

// demonstrateFind shows how find method in repogen works. It receives query parameters through method
// arguments and returns matched result
func demonstrateFind() {
	userModel, err := userRepository.FindByUsername(context.Background(), "sunboyy")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Find (one): found user = %+v\n", userModel)
}

// demonstrateUpdate shows how update method in repogen works. It receives updates and query parameters
// through method arguments and returns true/false whether there is a matched query.
func demonstrateUpdate() {
	matched, err := userRepository.UpdateDisplayNameByID(context.Background(), "Sunboy", userID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Update (one): matched = %v\n", matched)
}

// demonstrateDelete shows how delete method in repogen works. It receives query parameters through
// method arguments and returns matched count.
func demonstrateDelete() {
	matchedCount, err := userRepository.DeleteByCity(context.Background(), "Bangkok, Thailand")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Delete (many): matched count = %d\n", matchedCount)
}
