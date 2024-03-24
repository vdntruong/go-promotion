package ejwt

import "github.com/golang-jwt/jwt/v5"

type EClaims struct {
	jwt.RegisteredClaims
	Username     string `json:"username"`
	Email        string `json:"email"`
}