package repository

import (
	"api-auth/domains"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepository struct{
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{
		db: db, 
	}
}

func (ur *UserRepository) FindByEmail(email string) (*User) {
	var user User
	
	result := ur.db.First(&user, "email = ?", email);
	if result.Error != nil {
		return nil
	}

	return &user
}

func (ur *UserRepository) FindById(userId string) (*User) {
	var user User

	result := ur.db.First(&user, "id = ?", userId)
	if result.Error != nil {
		return nil
	}
	return &user
}


func (ur *UserRepository) CreateUser(input *domains.Register) error {
	uuid := uuid.New()

	user := User{
		ID: uuid.String(),
		Name: input.Name,
		Email: input.Email,
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
	user.Password = input.NewPassword
	ur.db.Save(&user)

	return nil
}

func (ur *UserRepository) Users() ([]User) {
	var users []User
	
	result := ur.db.Find(&users)
	if result.Error != nil {
		return nil
	}
	return users
}

func (ur *UserRepository) DeleteUserById(userId string) error {
	result := ur.db.Delete(&User{}, "id = ?", userId)
	if result.Error != nil {
		fmt.Println(result.Error)
		return errors.New("Something wrong, cannot delete single user!")
	}
	return nil
}