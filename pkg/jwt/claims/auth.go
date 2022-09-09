package claims

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type AuthClaims struct {
	jwt.StandardClaims
	Uid uuid.UUID
}

func NewAuthClaimsIngot() *AuthClaims {
	return &AuthClaims{}
}

func NewAuthClaims(userId uuid.UUID, duration time.Duration) *AuthClaims {
	return &AuthClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Uid: userId,
	}
}
