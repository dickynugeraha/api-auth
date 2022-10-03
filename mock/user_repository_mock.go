package mock

import (
	"api-auth/domains"
	repo "api-auth/services/repository"
	"errors"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) FindByEmail(email string) *repo.User {

	args := repository.Mock.Called(email)
	if args.Get(0) == nil {
		return nil
	}
	user := args.Get(0).(repo.User)
	return &user
}

func (repository *UserRepositoryMock) FindById(userId string) *repo.User {
	args := repository.Mock.Called(userId)
	if args.Get(0) == nil {
		return nil
	}
	user := args.Get(0).(repo.User)
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

func (repository *UserRepositoryMock) UpdatePassword(input *domains.ChangePassword, userId string) error {
	args := repository.Mock.Called(input,userId)
	if args.Get(0) != nil {
		return errors.New("Cannot update password!")
	}
	return nil
}

// func (repository *UserRepositoryMock) Users() (users []repo.User, err error) {
// 	args := repository.Mock.Called()

// 	if rf, ok := args.Get(0).(func() []repo.User); ok {
// 		users = rf()
// 	} else {
// 		if args.Get(0) != nil {
// 			users = args.Get(0).([]repo.User)
// 		}
// 	}

// 	if rf, ok := args.Get(1).(func() error); ok {
// 		err = rf()
// 	} else {
// 		err = args.Error(1)
// 	}

// 	return
// }

func (repository *UserRepositoryMock) Users() ([]repo.User, error) {
	args := repository.Mock.Called()

	if args.Get(1) != nil {
		return nil, errors.New("users not found")
	} else {
		users := args.Get(0).([]repo.User)
		return users, nil
	}
}
