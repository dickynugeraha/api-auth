package logic

import (
	"api-auth/app/helper"
	"api-auth/domains"
	mokz "api-auth/mock"
	"api-auth/services/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var userRepository = &mokz.UserRepositoryMock{Mock: mock.Mock{}}
var userUsecase = UserUsecase{Repository: userRepository}

func TestUserUsecase_SuccessRegisterHandler(t *testing.T) {

	input := &domains.Register{
		Name:            "Ardhito",
		Email:           "ardhito@gmail.com",
		Password:        "password",
		PasswordConfirm: "password",
	}

	userRepository.Mock.On("FindByEmail", input.Email).Return(nil).Once()

	userRepository.Mock.On("CreateUser", input).Return(nil).Once()

	err := userUsecase.RegisterHandler(input)
	assert.Nil(t, err)
}

func TestUserUsecase_FailedRegisterHandler(t *testing.T) {
	type Tested struct {
		name     string
		request  *domains.Register
		expected error
	}

	tests := []Tested{
		{
			name: "user_fail_1",
			request: &domains.Register{
				Name:            "kale",
				Email:           "kale",
				Password:        "123456789",
				PasswordConfirm: "123456789",
			},
			expected: errors.New("Email does'n contains @ or ."),
		},
		{
			name: "user_fail_2",
			request: &domains.Register{
				Name:            "kale",
				Email:           "kale@gmail.com",
				Password:        "12345678",
				PasswordConfirm: "12345678",
			},
			expected: errors.New("Email has been used!"),
		},
		{
			name: "user_fail_3",
			request: &domains.Register{
				Name:            "leo",
				Email:           "leo@gmail.com",
				Password:        "1234",
				PasswordConfirm: "1234",
			},
			expected: errors.New("Password must be greater than 8 characters!"),
		},
		{
			name: "user_fail_4",
			request: &domains.Register{
				Name:            "leo",
				Email:           "leo@gmail.com",
				Password:        "123456789",
				PasswordConfirm: "1234987655",
			},
			expected: errors.New("Password not match!"),
		},
	}

	for _, test := range tests {
		user := repository.User{
			ID:       "uuid",
			Name:     test.request.Name,
			Email:    test.request.Email,
			Password: test.request.Password,
		}

		t.Run(test.name, func(t *testing.T) {
			userRepository.Mock.On("FindByEmail", test.request.Email).Return(user).Once()
			userRepository.Mock.On("CreateUser", user.ID, test.request).Return(errors.New("Cannot create user!")).Once()

			err := userUsecase.RegisterHandler(test.request)

			assert.Equal(t, test.expected, err)
			assert.NotNil(t, err)
		})
	}

	test2 := Tested{
		name: "user_fail_5",
		request: &domains.Register{
			Name:            "leod",
			Email:           "leod@gmail.com",
			Password:        "123456789",
			PasswordConfirm: "123456789",
		},
		expected: errors.New("Cannot create user!"),
	}
	t.Run(test2.name, func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", test2.request.Email).Return(nil).Once()
		userRepository.Mock.On("CreateUser", test2.request).Return(errors.New("Cannot create user!")).Once()

		err := userUsecase.RegisterHandler(test2.request)

		assert.Equal(t, test2.expected, err)
		assert.NotNil(t, err)
	})
}

func TestUserUsecase_SuccessLoginHandler(t *testing.T) {
	tests := []struct {
		name    string
		request *domains.Login
	}{
		{
			name: "user_success_1",
			request: &domains.Login{
				Name:     "joko",
				Email:    "joko@gmail.com",
				Password: "passwords",
			},
		},
		{
			name: "user_success_2",
			request: &domains.Login{
				Name:     "daniel",
				Email:    "daniel@gmail.com",
				Password: "12345678",
			},
		},
	}

	for _, test := range tests {
		user1 := repository.User{
			ID:       "uuid",
			Name:     test.request.Name,
			Email:    test.request.Email,
			Password: helper.PasswordHashing(test.request.Password),
		}

		t.Run(test.name, func(t *testing.T) {
			userRepository.Mock.On("FindByEmail", test.request.Email).Return(user1).Once()
			user, token, err := userUsecase.LoginHandler(test.request)

			assert.NotNil(t, user)
			assert.Nil(t, err)
			assert.NotEmpty(t, token)
			assert.Equal(t, user1.Email, user.Email)
			assert.Equal(t, user1.Name, user.Name)
		})
	}
}

func TestUserUsecase_FailedLoginHandler(t *testing.T) {
	tests := []struct {
		name     string
		request  *domains.Login
		expected error
	}{
		{
			name: "user_fail_1",
			request: &domains.Login{
				Name:     "kale",
				Email:    "kale@gmailcom",
				Password: "password",
			},
			expected: errors.New("Email does'n contains @ or ."),
		},
		{
			name: "user_fail_2",
			request: &domains.Login{
				Name:     "rahmad",
				Email:    "rahmad@gmail.com",
				Password: "password",
			},
			expected: errors.New("Email not register!"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userRepository.Mock.On("FindByEmail", test.request.Email).Return(nil).Once()

			user, token, _ := userUsecase.LoginHandler(test.request)

			assert.Nil(t, user)
			assert.Empty(t, token)
			// assert.NotNil(t, err)
			// assert.Equal(t, err, test.expected)
		})
	}

	inputUser := &domains.Login{
		Name:     "kale",
		Email:    "kale@gmail.com",
		Password: "12345678",
	}
	userQuery := repository.User{
		ID:       "id",
		Name:     "kale",
		Email:    "kale@gmail.com",
		Password: "123456789",
	}
	t.Run("user_fail_3", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", inputUser.Email).Return(userQuery)

		user, token, err := userUsecase.LoginHandler(inputUser)

		assert.NotNil(t, user)
		assert.Empty(t, token)
		assert.NotNil(t, err)
		assert.Equal(t, err, errors.New("Invalid password!"))
	})
}

func TestUserUsecase_GetSingleUserHandler(t *testing.T) {
	userId := "random_id"

	t.Run("user_not_found", func(t *testing.T) {
		userRepository.Mock.On("FindById", userId).Return(nil)
		user, err := userUsecase.GetSingleUserHandler(userId)
		customError := errors.New("User not found!")

		assert.Nil(t, user)
		assert.Equal(t, err, customError)
		assert.NotNil(t, err)
	})

	userId = "uuid_hash"

	user1 := repository.User{
		ID:       userId,
		Name:     "Bjorka",
		Email:    "bjorka@gmail.com",
		Password: "hash_password",
	}

	t.Run("success_get_user", func(t *testing.T) {
		userRepository.Mock.On("FindById", userId).Return(user1)

		user, err := userUsecase.GetSingleUserHandler(userId)

		assert.NotNil(t, user)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, user1.ID)

	})
}

func TestUserUsecase_FailedChangePasswordHandler(t *testing.T) {
	tests := []struct {
		name     string
		request  *domains.ChangePassword
		expected error
	}{
		{
			name: "user_fail_1",
			request: &domains.ChangePassword{
				Email:           "kalegmail.com",
				NewPassword:     "password123",
				PasswordConfirm: "password123",
			},
			expected: errors.New("Email does'n contains @ or ."),
		},
		{
			name: "user_fail_2",
			request: &domains.ChangePassword{
				Email:           "kale@gmail.com",
				NewPassword:     "password12",
				PasswordConfirm: "password123",
			},
			expected: errors.New("Password not match!"),
		},
		{
			name: "user_fail_3",
			request: &domains.ChangePassword{
				Email:           "kale@gmail.com",
				NewPassword:     "pass",
				PasswordConfirm: "pass",
			},
			expected: errors.New("Password must be greater than 8 characters!"),
		},
		{
			name: "user_fail_4",
			request: &domains.ChangePassword{
				Email:           "email_not_found@gmail.com",
				NewPassword:     "password",
				PasswordConfirm: "password",
			},
			expected: errors.New("Cannot update password!"),
		},
	}

	for _, test := range tests {
		userRepository.Mock.On("FindByEmail", test.request.Email).Return(nil)
		userRepository.Mock.On("UpdatePassword", test.request).Return(errors.New("Error"))
		err := userUsecase.ChangePasswordHandler(test.request)
		assert.Error(t, err)
		// assert.Equal(t, err, test.expected)
	}
}

func TestUserUsecase_SuccessChangePasswordHandler(t *testing.T) {
	userInput := &domains.ChangePassword{
		Email:           "kale@gmail.com",
		NewPassword:     "password",
		PasswordConfirm: "password",
	}

	userRepository.Mock.On("FindByEmail", userInput.Email).Return(repository.User{ID: "id", Name: "kale", Email: userInput.Email, Password: userInput.PasswordConfirm})
	userRepository.Mock.On("UpdatePassword", userInput, "id").Return(nil)
	err := userUsecase.ChangePasswordHandler(userInput)

	assert.Nil(t, err)
}

func TestUserUsecase_DeleteUserHandler(t *testing.T) {
	userId := "random_uuid"

	t.Run("failed_delete_user", func(t *testing.T) {
		userRepository.Mock.On("DeleteUserById", userId).Return(errors.New("")).Once()

		err := userUsecase.DeleteUserHandler(userId)
		assert.NotNil(t, err)

	})

	userId = "real_uuid"
	t.Run("success_delete_user", func(t *testing.T) {
		userRepository.Mock.On("DeleteUserById", userId).Return(nil)

		err := userUsecase.DeleteUserHandler(userId)
		assert.Nil(t, err)
	})
}

func TestUserUsecase_SuccessGetUsers(t *testing.T) {

	users_test := []repository.User{
		{
			ID:       "uuid1",
			Name:     "kale",
			Email:    "kale@gmail.com",
			Password: "password",
		},
		{
			ID:       "uuid2",
			Name:     "kal2",
			Email:    "kale2@gmail.com",
			Password: "password2",
		},
	}

	userRepository.Mock.On("Users").Return(users_test, nil).Once()

	users_got, err := userUsecase.GetUsers()

	// assert.Nil(t, users_got)
	// assert.NotNil(t, err)

	assert.NotEmpty(t, users_got)
	assert.Nil(t, err)
}

func TestUserUsecase_FailGetUsers(t *testing.T) {

	userRepository.Mock.On("Users").Return(nil, errors.New(""))

	users, err := userUsecase.GetUsers()
	errExpected := errors.New("Users not found!")

	assert.NotNil(t, err)
	assert.Nil(t, users)
	assert.Equal(t, errExpected, err)

}
