package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("d77bb076-ac99-46ee-8a50-79ebc28ce154")

type CustomClaims struct {
	Name     string `json:"user_name"`
	Password string `json:"user_password"`
	jwt.RegisteredClaims
}

type TokenService interface {
	GenerateToken(name, password string) (string, error)
}

type JWTService struct{}

func (j *JWTService) GenerateToken(name, password string) (string, error) {
	claims := CustomClaims{
		Name:     name,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "move-ass",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
