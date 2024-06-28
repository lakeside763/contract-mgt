package models

import (
	"time"
)


type ProfileType BaseType

const (
	PROFILE_CLIENT			ProfileType = "CLIENT"
	PROFILE_CONTRACTOR	ProfileType = "CONTRACTOR"
)


type Profile struct {
	ID         string       `json:"id" gorm:"type:uuid;default:uuid_generate_v4();column:id"`
	FirstName  string       `json:"first_name" gorm:"type:varchar(255);column:first_name"`
	LastName   string       `json:"last_name" gorm:"type:varchar(255);column:last_name"`
	Profession string       `json:"profession"`
	Balance    float64      `json:"balance"`
	Type       ProfileType  `json:"type"`
	CreatedAt  time.Time    `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt  time.Time    `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}

type ProfileAndUserCredentials struct {
	Profile
	Credentials
}