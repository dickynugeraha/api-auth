package repository

import "01_REST_Auth/domains"

import (
	// "context"
)

type UserRepositoryInterface interface {
	FindbyEmail(email string) (User, error)
	CreateUser(input *domains.Register, passwordHash string) (error)
	UpdatePassword(email, password string) (error)
	Users() (*[]User, error)
	GetUserById(userId string) (User, error)
	DeleteUserById(user_id string) (User, error)
}