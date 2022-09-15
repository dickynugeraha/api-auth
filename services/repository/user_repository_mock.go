package repository

import (
	"01_REST_Auth/domains"
	"errors"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) FindByEmail(email string) (*User) {
	args := repository.Mock.Called(email)
	if args.Get(0) == nil {
		return nil
	}
	user := args.Get(0).(User)
	return &user
}

func (repository *UserRepositoryMock) FindById(userId string) (*User) {
	args := repository.Mock.Called(userId)
	if args.Get(0) == nil {
		return nil
	}
	user := args.Get(0).(User)
	return &user
}

func (repository *UserRepositoryMock) CreateUser(input *domains.Register) error {
	args := repository.Mock.Called(input)
	if args.Get(0) != nil {
		return errors.New("Cannot create user!")
	}
	return nil
}

func (repository *UserRepositoryMock) DeleteUserById(userId string) error {
	args := repository.Mock.Called(userId)
	if args.Get(0) != nil {
		return errors.New("Cannot delete user!")
	}
	return nil
}

func (repository *UserRepositoryMock) UpdatePassword(input *domains.ChangePassword) error {
	args := repository.Mock.Called(input)
	if args.Get(0) != nil {
		return errors.New("Cannot update password!")
	}
	return nil
}

func (repository *UserRepositoryMock) Users() *[]User {
	var users *[]User
	args := repository.Mock.Called()
	if args.Get(0) == nil {
		return nil
	}

	users = args.Get(0).(*[]User)
	return users
}