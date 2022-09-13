package repository

import (
	// "context"
	"01_REST_Auth/domains"
)

type UserRepositoryInterface interface {
	FindbyEmail(email string) (User, error)
	CreateUser(input *domains.Register, passwordHashing string) (error)
	UpdatePassword(email, newPassword string) (error)
	Users() ([]User, error)
	GetUserById(userId string) (User, error)
	DeleteUserById(userId string) (error)
}