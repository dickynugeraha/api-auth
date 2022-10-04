package repository

import (
	"api-auth/domains"
	"errors"

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
	user := User{}

	result := ur.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil
	}

	return &user
}

func (ur *UserRepository) FindById(userId string) *User {
	user := User{}

	result := ur.db.First(&user, "id = ?", userId)
	if result.Error != nil {
		return nil
	}
	return &user
}

func (ur *UserRepository) CreateUser(uuid string, input *domains.Register) error {
	user := User{}

	newUser := User{
		ID:       uuid,
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	result := ur.db.Model(&user).Create(&newUser)

	if result.Error != nil {
		return errors.New("Cannot create user!")
	}

	return nil
}

func (ur *UserRepository) UpdatePassword(input *domains.ChangePassword, userId string) (err error) {

	result := ur.db.Model(&User{}).Where("id = ?", userId).Update("password", input.NewPassword)
	// db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
	// user.Password = input.NewPassword
	// ur.db.Save(&user)

	err = result.Error

	if err != nil {
		return 
	}
	return 
}

func (ur *UserRepository) Users() ([]User, error) {
	var users []User

	result := ur.db.Find(&users)

	if result.Error != nil {
		return nil, errors.New("User not found!")
	}

	return users, nil
}

func (ur *UserRepository) DeleteUserById(userId string) error {
	user := User{}

	result := ur.db.Where("id = ?", userId).Delete(&user)
	if result.Error != nil {
		return errors.New("Something wrong, cannot delete single user!")
	}
	return nil
}
