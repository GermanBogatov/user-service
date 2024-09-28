package entity

import "github.com/golang-jwt/jwt/v4"

type UserClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}
