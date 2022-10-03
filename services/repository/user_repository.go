package repository

import (
	"api-auth/domains"
)

type UserRepositoryInterface interface {
	FindByEmail(email string) *User
	FindById(userId string) *User
	CreateUser(input *domains.Register) error
	UpdatePassword(input *domains.ChangePassword, userId string) error
	Users() ([]User, error)
	DeleteUserById(userId string) error
}
