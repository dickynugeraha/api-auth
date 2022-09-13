package repository

import (
	// db "01_REST_Auth/app/config"
	"01_REST_Auth/domains"
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

func (ur UserRepository) FindbyEmail(email string) (User, error) {
	// db := db.SetupMysql()
	var user User
	
	result := ur.db.First(&user, "email = ?", email);
	if result.Error != nil {
		return user, errors.New("Email not found!")
	}

	return user, nil
}

func (ur UserRepository) CreateUser(input *domains.Register, passwordHashed string) error {
	// db := db.SetupMysql()
	uuid := uuid.New()

	user := User{
		ID: uuid.String(),
		Name: input.Name,
		Email: input.Email,
		Password: passwordHashed,
	}

	result := ur.db.Create(&user)

	if result.RowsAffected < 0 {
		return result.Error
	}
	
	return nil
}

func (ur UserRepository) UpdatePassword(email, newPassword string) error {
	// db := db.SetupMysql()
	user, _ := ur.FindbyEmail(email)
	user.Password = newPassword
	ur.db.Save(&user)

	return nil
}

func (ur UserRepository) Users() ([]User, error) {
	// db := db.SetupMysql()
	var users []User
	
	result := ur.db.Find(&users)
	if result.Error != nil {
		return users, errors.New("Users not found!")
	}
	return users, nil
}

func (ur UserRepository) GetUserById(userId string) (User, error) {
	// db := db.SetupMysql()
	var user User

	result := ur.db.First(&user, "id = ?", userId)
	if result.Error != nil {
		return user, errors.New("User by id not found!")
	}
	return user, nil
}

func (ur UserRepository) DeleteUserById(userId string) error {
	// db := db.SetupMysql()

	result := ur.db.Delete(&User{}, "id = ?", userId)
	if result.Error != nil {
		fmt.Println(result.Error)
		return errors.New("Something wrong, cannot delete single user!")
	}
	return nil
}