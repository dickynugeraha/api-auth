package repository

import (
	// "context"
	"01_REST_Auth/domains"
)

type UserRepositoryInterface interface {
	FindByEmail(email string) (*User)
	FindById(userId string) (*User)
	CreateUser(input *domains.Register) (error)
	UpdatePassword(input *domains.ChangePassword) (error)
	Users() (*[]User)
	DeleteUserById(userId string) (error)
}