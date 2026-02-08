package crosspackage

import "go.mongodb.org/mongo-driver/v2/bson"

// UserModel is a model of user that is stored in the database
type UserModel struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string        `bson:"username" json:"username"`
	DisplayName string        `bson:"display_name" json:"displayName"`
	City        string        `bson:"city" json:"city"`
}
