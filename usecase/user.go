package usecase

import (
	"01_REST_Auth/domains"
	"01_REST_Auth/services/repository"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo repository.UserRepositoryInterface
}

func (uu *UserUsecase) AllUsers() (*[]repository.User, error) {
	
	users, err := uu.repo.Users()


	if err != nil {
		return users, err
	}
	return users, nil
}

func (uu *UserUsecase) RegisterHandler(input *domains.Register) (int, error) {
	err := emailRequired(input.Email)
	statusHttp := 500
	if err != nil {
		statusHttp = 400
		return statusHttp, err
	}
	_, err = uu.repo.FindbyEmail(input.Email)
	if err == nil {
		statusHttp = 400
		return statusHttp, errors.New("Email has been used!")
	}
	err = passwordRequired(input.Password, input.PasswordConfirm)
	if err != nil {
		statusHttp = 400
		return statusHttp, err
	}
	hashedPassword := passwordHashing(input.Password)
	err = uu.repo.CreateUser(input, hashedPassword)
	if err != nil {
		return statusHttp, err
	}
	return 201, nil
}

func (uu *UserUsecase) LoginHandler(input *domains.Login) (*repository.User, string, error) {
	err := emailRequired(input.Email)
	if err != nil {
		return &repository.User{}, "", err
	}

	user, err := uu.repo.FindbyEmail(input.Email)
	if err != nil {
		return &user, "",err
	}

	err = checkPasswordHash(input.Password, user.Password)
	if err != nil {
		return &user, "", err
	}

	token, err := generateJWT(user.ID, user.Email)
	if err != nil {
		return &user, "", err
	}
	return &user, token, nil
}

func (uu *UserUsecase) ChangePasswordHandler(input *domains.ChangePassword) error {
	_, err := uu.repo.FindbyEmail(input.Email);
	if err != nil {
		return err
	}
	err = passwordRequired(input.NewPassword, input.PasswordConfirm);
	if err != nil {
		return err
	}
	newPasswordHash := passwordHashing(input.NewPassword)
	err =	uu.repo.UpdatePassword(input.Email, newPasswordHash);
	if err != nil {
		return err
	}
	return err
}

func generateJWT(id string, email string) (string, error){
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

func emailRequired(email string) error {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return	errors.New("Email does'n contains @ or .")
	}
		return nil
}

func passwordRequired(pass, confPass string) error {
	if len(pass) < 8 {
		return errors.New("Password must be greater than 8 characters!")
	}
	if strings.Compare(pass, confPass) != 0 {
		return errors.New("Password not match!")
	}
	return nil
}

func checkPasswordHash(passEntered, passHashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passHashed), []byte(passEntered))

	if err != nil {
		return errors.New("Invalid password!")
	}

	return nil
}

func passwordHashing(pw string) string {
	pen, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)

	if err != nil {
		return "Cannot hash password"
	}
	return string(pen)
}