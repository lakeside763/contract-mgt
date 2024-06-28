package models

import "time"

type Role BaseType

const (
	ADMIN				Role = "ADMIN"
	CLIENT			Role = "CLIENT"
	CONTRACTOR	Role = "CONTRACTOR"
)
type User struct {
	ID 				string 			`json:"id" gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	Username 	string 			`json:"username" gorm:"type:varchar(225);column:username"`
	Password 	string 			`json:"password" gorm:"type:varchar(225);column:password"`
	Role 			Role				`json:"role" gorm:"type:varchar(50)"`
	CreatedAt  time.Time  `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}

var Users []User

type UserResponse struct {
	ID        string    	`json:"id"`
	Username  string    	`json:"username"`
	Role      string    	`json:"role"`
	CreatedAt  time.Time 	`json:"createdAt" gorm:"type:timestamp;column:createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt" gorm:"type:timestamp;column:updatedAt"`
}

type UserResponseWithToken struct {
	UserResponse
	Token string `json:"token"`
}