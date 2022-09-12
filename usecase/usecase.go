package usecase

import (
	"01_REST_Auth/services/repository"
	"01_REST_Auth/domains"
)

type UserUsecaseInterface interface {
	AllUsers() ([]*repository.User, error)
	RegisterHandler(*domains.Register) (int, error)
	LoginHandler(*domains.Login) (repository.User)
	GetSingleUserHandler()
	DeleteUserHadler() error
}