package helper

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(id string, email string) (string, error) {
	var mySigningKey = []byte("secretkey")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = id
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func EmailRequired(email string) error {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("Email does'n contains @ or .")
	}
	return nil
}

func PasswordRequired(pass, confPass string) error {
	if len(pass) < 8 {
		return errors.New("Password must be greater than 8 characters!")
	}
	if strings.Compare(pass, confPass) != 0 {
		return errors.New("Password not match!")
	}
	return nil
}

func CheckPasswordHash(passEntered, passHashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passHashed), []byte(passEntered))

	if err != nil {
		return errors.New("Invalid password!")
	}

	return nil
}

func PasswordHashing(pw string) string {
	pen, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)

	if err != nil {
		return "Cannot hash password"
	}
	return string(pen)
}
