package repository

import (
	"api-auth/domains"
	"errors"
	"fmt"

	"github.com/google/uuid"
	// "gorm.io/gorm"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID       string `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) FindByEmail(email string) *User {
	var user User

	result := ur.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil
	}

	return &user
}

func (ur *UserRepository) FindById(userId string) *User {
	var user User

	result := ur.db.First(&user, "id = ?", userId)
	if result.Error != nil {
		return nil
	}
	return &user
}

func (ur *UserRepository) CreateUser(input *domains.Register) error {
	user := User{
		ID:       uuid.New().String(),
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	result := ur.db.Create(&user)

	if result.RowsAffected < 0 {
		return errors.New("Cannot create user!")
	}

	return nil
}

func (ur *UserRepository) UpdatePassword(input *domains.ChangePassword) error {
	user := ur.FindByEmail(input.Email)
	if user == nil {
		return errors.New("Email not register!")
	}

	user.Password = input.NewPassword
	ur.db.Save(&user)

	return nil
}

func (ur *UserRepository) Users() ([]User, error) {
	var users []User

	result := ur.db.Find(&users)

	if result.Error != nil {
		return nil, errors.New("")
	}

	return users, nil
}

func (ur *UserRepository) DeleteUserById(userId string) error {
	result := ur.db.Delete(&User{}, "id = ?", userId)
	if result.Error != nil {
		fmt.Println(result.Error)
		return errors.New("Something wrong, cannot delete single user!")
	}
	return nil
}
