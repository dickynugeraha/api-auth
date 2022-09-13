package usecase

import (
	"01_REST_Auth/services/repository"
	"01_REST_Auth/domains"
)

type UserUsecaseInterface interface {
	AllUsers() ([]repository.User, error)
	RegisterHandler(input *domains.Register) (int, error)
	LoginHandler(input *domains.Login) (repository.User, string, error)
	ChangePasswordHandler(input *domains.ChangePassword) (error)
	GetSingleUserHandler(userId string) (repository.User, error)
	DeleteUserHadler(userId string) error
}