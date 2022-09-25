package usecase

import (
	"api-auth/domains"
	"api-auth/services/repository"
)

type UserUsecaseInterface interface {
	GetUsers() ([]repository.User, error)
	RegisterHandler(input *domains.Register) error
	LoginHandler(input *domains.Login) (*repository.User, string, error)
	ChangePasswordHandler(input *domains.ChangePassword) error
	GetSingleUserHandler(userId string) (*repository.User, error)
	DeleteUserHandler(userId string) error
}
