package models

import (
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("secret-1234")

type Credentials struct {
	Username 	string	`json:"username"`
	Password	string	`json:"password"`
}

type Claims struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
