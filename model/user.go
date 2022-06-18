package model

type User struct {
	ID        string `json:"id" bson:"_id"`
	Firstname string `json:"firstname" bson:"firstname"`
	Lastname  string `json:"lastname" bson:"lastname"`
	Username  string `json:"username" bson:"username"`
	Password  string `json:"password" bson:"password"`
}
