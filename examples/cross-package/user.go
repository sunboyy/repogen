package cross_package

// UserModel is a model of user that is stored in the database
type UserModel struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Username    string `bson:"username" json:"username"`
	DisplayName string `bson:"display_name" json:"displayName"`
	City        string `bson:"city" json:"city"`
}
