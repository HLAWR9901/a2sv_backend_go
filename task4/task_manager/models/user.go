package models

type Class string

const (
	Admin   Class = "admin"
	Regular Class = "regular"
)

type User struct {
	ID       string `json:"id" bson:"id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Role     Class  `json:"role" bson:"role"`
}

type AuthUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Class  `json:"role"`
}