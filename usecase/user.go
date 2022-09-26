package usecase

import (
	"api-auth/app/helper"
	"api-auth/domains"
	"api-auth/services/repository"
	"errors"
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

	err := helper.EmailRequired(input.Email)
	if err != nil {
		return err
	}
	err = helper.PasswordRequired(input.Password, input.PasswordConfirm)
	if err != nil {
		return err
	}
	user := uu.Repository.FindByEmail(input.Email)
	if user != nil {
		return errors.New("Email has been used!")
	}
	input.Password = helper.PasswordHashing(input.Password)
	err = uu.Repository.CreateUser(input)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UserUsecase) LoginHandler(input *domains.Login) (*repository.User, string, error) {
	err := helper.EmailRequired(input.Email)
	if err != nil {
		return nil, "", err
	}
	user := uu.Repository.FindByEmail(input.Email)
	if user == nil {
		return user, "", errors.New("Email not register!")
	}
	err = helper.CheckPasswordHash(input.Password, user.Password)
	if err != nil {
		return user, "", err
	}
	token, _ := helper.GenerateJWT(user.ID, user.Email)

	return user, token, nil
}

func (uu *UserUsecase) ChangePasswordHandler(input *domains.ChangePassword) error {
	err := helper.EmailRequired(input.Email)
	if err != nil {
		return err
	}
	err = helper.PasswordRequired(input.NewPassword, input.PasswordConfirm)
	if err != nil {
		return err
	}
	newPasswordHash := helper.PasswordHashing(input.NewPassword)
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
