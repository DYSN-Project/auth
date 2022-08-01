package jwt

import (
	"errors"
	"github.com/DYSN-Project/auth/internal/packages/jwt/claims"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

var (
	errInvalidSign  = errors.New("invalid method sign")
	errTokenInvalid = errors.New("invalid token")
	errClaims       = errors.New("get token claims error")
)

type JwtInterface interface {
	GenerateAuthToken(userId uuid.UUID, duration time.Duration) (string, error)
	GenerateRegisterToken(email, password string, duration time.Duration) (string, error)
	GenerateToken(claims jwt.Claims, secretKey string) (string, error)
	ParseRegisterToken(strToken string) (*claims.RegisterClaims, error)
	ParseAuthToken(strToken string) (*claims.AuthClaims, error)
	Verify(strToken string, secretKey string) (bool, error)
}

type JwtService struct {
	authSecretKey     string
	registerSecretKey string
}

func NewJwtService(authSecretKey, registerSecretKey string) *JwtService {
	return &JwtService{
		authSecretKey:     authSecretKey,
		registerSecretKey: registerSecretKey,
	}
}

func (j *JwtService) GenerateAuthToken(userId uuid.UUID, duration time.Duration) (string, error) {
	clms := claims.NewAuthClaims(userId, duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clms)

	return token.SignedString([]byte(j.authSecretKey))
}

func (j *JwtService) GenerateRegisterToken(email, password string, duration time.Duration) (string, error) {
	clms := claims.NewRegisterClaims(email, password, duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clms)

	return token.SignedString([]byte(j.registerSecretKey))
}

func (j *JwtService) GenerateToken(claims jwt.Claims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func (j *JwtService) ParseRegisterToken(strToken string) (*claims.RegisterClaims, error) {
	token, err := j.parse(strToken, claims.NewRegisterClaimsIngot(), j.registerSecretKey)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errTokenInvalid
	}

	clms, ok := token.Claims.(*claims.RegisterClaims)
	if !ok {
		return nil, errClaims
	}

	return clms, nil
}

func (j *JwtService) ParseAuthToken(strToken string) (*claims.AuthClaims, error) {
	token, err := j.parse(strToken, claims.NewAuthClaimsIngot(), j.authSecretKey)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errTokenInvalid
	}

	clms, ok := token.Claims.(*claims.AuthClaims)
	if !ok {
		return nil, errClaims
	}

	return clms, nil
}

func (j *JwtService) Verify(strToken string, secretKey string) (bool, error) {
	token, parseErr := j.parse(strToken, jwt.MapClaims{}, secretKey)
	if parseErr != nil {
		return false, parseErr
	}

	return token.Valid, nil
}

func (j *JwtService) parse(strToken string, model jwt.Claims, key string) (*jwt.Token, error) {
	token, parseErr := jwt.ParseWithClaims(strToken, model, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidSign
		}
		return []byte(key), nil
	})

	return token, parseErr
}
