package repository

import (
	"api-auth/domains"
)

type UserRepositoryInterface interface {
	FindByEmail(email string) *User
	FindById(userId string) *User
	CreateUser(input *domains.Register) error
	UpdatePassword(input *domains.ChangePassword) error
	Users() ([]User, error)
	DeleteUserById(userId string) error
}
