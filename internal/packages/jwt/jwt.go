package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

var (
	errInvalidSign  = errors.New("invalid method sign")
	errTokenInvalid = errors.New("invalid token")
	errClaims       = errors.New("get token claims error")
)

type JwtInterface interface {
	GenerateToken(data map[string]interface{},
		secretKey string,
		duration time.Duration) (string, error)
	ParseToken(strToken string, key string) (map[string]interface{}, error)
	Verify(strToken string, secretKey string) (bool, error)
}

type JwtService struct{}

func NewJwtService() *JwtService {
	return &JwtService{}
}

func (j *JwtService) GenerateToken(data map[string]interface{},
	secretKey string,
	duration time.Duration) (string, error) {
	cls := j.getClaims(data, duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cls)

	return token.SignedString([]byte(secretKey))
}

func (j *JwtService) ParseToken(strToken string, key string) (map[string]interface{}, error) {
	token, err := j.parse(strToken, jwt.MapClaims{}, key)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errTokenInvalid
	}

	clms, ok := token.Claims.(jwt.MapClaims)
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

func (j *JwtService) getClaims(data map[string]interface{}, dur time.Duration) jwt.MapClaims {
	cls := jwt.MapClaims{}
	for k, v := range data {
		cls[k] = v
	}
	cls["exp"] = time.Now().Add(dur).Unix()

	return cls
}
