package usecase

import (
	"api-auth/domains"
	"api-auth/services/repository"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Repository repository.UserRepositoryInterface
}

func NewUserUsecase(Repository repository.UserRepositoryInterface) UserUsecaseInterface {
	return &UserUsecase{
		Repository: Repository,
	}
}

func (uu *UserUsecase) GetUsers() ([]repository.User, error) {
	users, err := uu.Repository.Users()

	if err != nil || users == nil {
		return nil, errors.New("Users not found!")
	}

	return users, nil
}

func (uu *UserUsecase) RegisterHandler(input *domains.Register) error {
	err := EmailRequired(input.Email)
	if err != nil {
		return err
	}
	err = PasswordRequired(input.Password, input.PasswordConfirm)
	if err != nil {
		return err
	}
	user := uu.Repository.FindByEmail(input.Email)
	if user != nil {
		return errors.New("Email has been used!")
	}
	hashedPassword := PasswordHashing(input.Password)
	input.Password = hashedPassword
	err = uu.Repository.CreateUser(input)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UserUsecase) LoginHandler(input *domains.Login) (*repository.User, string, error) {
	err := EmailRequired(input.Email)
	if err != nil {
		return nil, "", err
	}
	user := uu.Repository.FindByEmail(input.Email)
	if user == nil {
		return user, "", errors.New("Email not register!")
	}
	err = CheckPasswordHash(input.Password, user.Password)
	if err != nil {
		return user, "", err
	}
	token, err := generateJWT(user.ID, user.Email)
	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

func (uu *UserUsecase) ChangePasswordHandler(input *domains.ChangePassword) error {
	err := EmailRequired(input.Email)
	if err != nil {
		return err
	}
	err = PasswordRequired(input.NewPassword, input.PasswordConfirm)
	if err != nil {
		return err
	}
	newPasswordHash := PasswordHashing(input.NewPassword)
	input.NewPassword = newPasswordHash
	err = uu.Repository.UpdatePassword(input)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UserUsecase) GetSingleUserHandler(userId string) (*repository.User, error) {
	user := uu.Repository.FindById(userId)
	if user == nil {
		return nil, errors.New("User not found!")
	}
	return user, nil
}

func (uu *UserUsecase) DeleteUserHandler(userId string) error {
	if err := uu.Repository.DeleteUserById(userId); err != nil {
		return err
	}
	return nil
}

func generateJWT(id string, email string) (string, error) {
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
