package mock

import (
	"api-auth/domains"
	"api-auth/services/repository"

	"github.com/stretchr/testify/mock"
)

type UserUsecaseMock struct {
	Mock mock.Mock
}

func (usecase *UserUsecaseMock) GetUsers() (users []repository.User, err error) {
	args := usecase.Mock.Called()

	if rf, ok := args.Get(0).(func() []repository.User); ok {
		users = rf()
	} else {
		if args.Get(0) != nil {
			users = args.Get(0).([]repository.User)
		} else {
			users = nil
		}
	}

	if rf, ok := args.Get(1).(func() error); ok {
		err = rf()
	} else {
		err = args.Error(1)
	}

	return
}

func (usecase *UserUsecaseMock) RegisterHandler(input *domains.Register) (err error) {
	args := usecase.Mock.Called(input)

	if args.Get(0) == nil {
		err = nil
	} else {
		err = args.Error(0)
	}

	return
}

func (usecase *UserUsecaseMock) LoginHandler(input *domains.Login) (user *repository.User, token string, err error) {
	args := usecase.Mock.Called(input)

	if rf, ok := args.Get(0).(func(*domains.Login) *repository.User); ok {
		user = rf(input)
	} else {
		if args.Get(0) != nil {
			user = args.Get(0).(*repository.User)
		}
	}

	if rf, ok := args.Get(1).(func(*domains.Login) string); ok {
		token = rf(input)
	} else {
		token = args.Get(1).(string)
	}

	if rf, ok := args.Get(2).(func(*domains.Login) error); ok {
		err = rf(input)
	} else {
		err = args.Error(2)
	}

	return
}

func (usecase *UserUsecaseMock) ChangePasswordHandler(input *domains.ChangePassword) (err error) {
	args := usecase.Mock.Called(input)

	if args.Get(0) != nil {
		err = args.Error(0)
	} else {
		err = nil
	}

	return
}

func (usecase *UserUsecaseMock) GetSingleUserHandler(userId string) (user *repository.User, err error) {
	args := usecase.Mock.Called(userId)

	if rf, ok := args.Get(0).(func(string) *repository.User); ok {
		user = rf(userId)
	} else {
		if args.Get(0) != nil {
			user = args.Get(0).(*repository.User)
		}
	}

	if rf, ok := args.Get(1).(func(string) error); ok {
		err = rf(userId)
	} else {
		err = args.Error(1)
	}

	return
}

func (usecase *UserUsecaseMock) DeleteUserHandler(userId string) (err error) {
	args := usecase.Mock.Called(userId)

	if args.Get(0) != nil {
		err = args.Error(0)
	} else {
		err = nil
	}
	return
}
