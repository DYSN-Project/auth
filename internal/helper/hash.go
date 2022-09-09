package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func GetHash(pwd string, salt string) (string, error) {
	password := []byte(pwd + salt)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHash(inputPwd string, userPwd string, salt string) error {
	return bcrypt.CompareHashAndPassword([]byte(userPwd), []byte(inputPwd+salt))
}
