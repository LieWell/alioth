package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password, salt string) (string, error) {
	psswd, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(psswd), nil
}

func checkPassword(password, salt, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
}
