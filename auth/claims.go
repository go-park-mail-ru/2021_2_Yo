package auth

import (
	"github.com/dgrijalva/jwt-go/v4"
)

type Claims struct {
	jwt.StandardClaims
	Name string `json:"name"`
	Surname string `json:"Surname"`
	Mail string `json:"mail"`
}
