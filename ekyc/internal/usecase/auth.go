package usecase

import (
	"time"

	"ekyc/internal/model"
	"ekyc/internal/pkg/ejwt"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	generator *ejwt.Generator
}

func NewAuth(g *ejwt.Generator) *Auth {
	return &Auth{generator: g}
}

func (a *Auth) GetToken(u model.User) (string, error) {
	var now = time.Now()
	var claims = ejwt.EClaims{
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    a.generator.Issuer,
			Subject:   u.Email,
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        u.Email,
		},
	}
	if u.Username != nil {
		claims.Username = *u.Username
	}

	t, err := a.generator.Sign(claims)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (a *Auth) VerifyToken(token string) error {
	_, err := a.generator.Verify(token)
	if err != nil {
		return err
	}
	return nil
}
