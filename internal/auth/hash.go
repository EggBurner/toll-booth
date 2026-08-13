package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (result string, err error) {

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func ComparePassword(hashedPassword string, toCheckPassword string) error {

	hashedPasswordByte := []byte(hashedPassword)
	toCheckPasswordByte := []byte(toCheckPassword)
	err := bcrypt.CompareHashAndPassword(hashedPasswordByte, toCheckPasswordByte)

	return err
}
