package claims

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type RegisterClaims struct {
	jwt.StandardClaims
	Email    string
	Password string
}

func NewRegisterClaimsIngot() *RegisterClaims {
	return &RegisterClaims{}
}

func NewRegisterClaims(email string, password string, duration time.Duration) *RegisterClaims {
	return &RegisterClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Email:    email,
		Password: password,
	}
}
