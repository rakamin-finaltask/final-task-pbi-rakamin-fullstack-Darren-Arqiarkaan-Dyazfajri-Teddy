package dto

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}
