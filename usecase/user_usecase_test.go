package usecase

import (
	"api-auth/domains"
	"api-auth/services/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var userRepository = &repository.UserRepositoryMock{Mock: mock.Mock{}}
var userUsecase = UserUsecase{Repository: userRepository}

func TestUsecase_SuccessRegisterHandler(t *testing.T){
	user := &domains.Register{
		Name: "Ardhito",
		Email: "ardhito@gmail.com",
		Password: "password",
		PasswordConfirm: "password",
	}

	t.Run("user_success", func(t *testing.T){
		userRepository.Mock.On("FindByEmail", user.Email).Return(nil)
		userRepository.Mock.On("CreateUser", user).Return(nil)

		err := userUsecase.RegisterHandler(user)

		assert.Nil(t, err)
	})
}

func TestUsecase_FailedRegisterHandler(t *testing.T){
tests := []struct{
		name string
		request *domains.Register
		expected error
	}{
		{
			name: "user_fail_1",
			request: &domains.Register{
				Name: "kale",
				Email: "kale",
				Password: "123456789",
				PasswordConfirm: "123456789",
			},
			expected: errors.New("Email does'n contains @ or ."),
		},
		{
			name: "user_fail_2",
			request: &domains.Register{
				Name: "kale",
				Email: "kale@gmail.com",
				Password: "12345678",
				PasswordConfirm: "12345678",
			},
			expected: errors.New("Email has been used!"),
		},
		{
			name: "user_fail_3",
			request: &domains.Register{
				Name: "leo",
				Email: "leo@gmail.com",
				Password: "1234",
				PasswordConfirm: "1234",
			},
			expected: errors.New("Password must be greater than 8 characters!"),
		},
		{
			name: "user_fail_4",
			request: &domains.Register{
				Name: "leo",
				Email: "leo@gmail.com",
				Password: "123456789",
				PasswordConfirm: "1234987655",
			},
			expected: errors.New("Password not match!"),
		},
	}


	for _, test := range tests{
		user := repository.User{
			ID: "kbkas",
			Name: test.request.Name,
			Email: test.request.Email,
			Password: test.request.Password,
		}

		t.Run(test.name, func(t *testing.T){
			userRepository.Mock.On("FindByEmail", test.request.Email).Return(user)
			userRepository.Mock.On("CreateUser", test.request).Return(errors.New("Cannot create user!"))

			err := userUsecase.RegisterHandler(test.request)

			assert.Equal(t, test.expected, err)
			assert.NotNil(t, err)
		})
	}
}

func TestUsecase_SuccessLoginHandler(t *testing.T){
	tests := []struct{
		name string
		request *domains.Login
	}{
		{
			name: "user_success_1",
			request: &domains.Login{
				Name: "joko",
				Email: "joko@gmail.com",
				Password: "password",
			},
		},
		{
			name: "user_success_2",
			request: &domains.Login{
				Name: "daniel",
				Email: "daniel@gmail.com",
				Password: "12345678",
			},
		},
	}

	for _, test := range tests {
		user1 := repository.User{
			ID: "kalsnl",
			Name: test.request.Name,
			Email: test.request.Email,
			Password: PasswordHashing(test.request.Password),
		}

		t.Run(test.name, func(t *testing.T){
			userRepository.Mock.On("FindByEmail", test.request.Email).Return(user1)
			user, token, err := userUsecase.LoginHandler(test.request)
			
			assert.NotNil(t, user)
			assert.Nil(t, err)
			assert.NotEmpty(t, token)
			assert.Equal(t, user1.Email, user.Email)
			assert.Equal(t, user1.Name, user.Name)
		})
	}
}

func TestUsecase_FailedLoginHandler(t *testing.T){
	tests := []struct{
		name string
		request *domains.Login
		expected error
	}{
		{
			name: "user_fail_1",
			request: &domains.Login{
				Name: "kale",
				Email: "kale@gmailcom",
				Password: "password",
			},
			expected: errors.New("Email does'n contains @ or ."),
		},
		{
			name: "user_fail_2",
			request: &domains.Login{
				Name: "rahmad",
				Email: "rahmad@gmail.com",
				Password: "password",
			},
			expected: errors.New("Email not register!"),
		},
	}

	for _, test := range tests{
		t.Run(test.name, func(t *testing.T){
			userRepository.Mock.On("FindByEmail", test.request.Email).Return(nil)
			
			user, token, err := userUsecase.LoginHandler(test.request)

			assert.Nil(t, user)
			assert.Empty(t, token)
			assert.NotNil(t, err)
			assert.Equal(t, err, test.expected)
		})
	}

	inputUser := &domains.Login{
		Name: "kale",
		Email: "kale@gmail.com",
		Password: "12345678",
	}
	userQuery := repository.User{
		ID: "id",
		Name: "kale",
		Email: "kale@gmail.com",
		Password: "123456789",
	}
	t.Run("user_fail_3", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", inputUser.Email).Return(userQuery)

		user, token, err := userUsecase.LoginHandler(inputUser)

		assert.NotNil(t, user)
		assert.Empty(t, token)
		assert.NotNil(t, err)
		assert.Equal(t, err, errors.New("Invalid password!"))
		assert.NotEqual(t, user.Password, inputUser.Password)
	})
}

func TestUsecase_GetSingleUserHandler(t *testing.T){
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
		ID: userId,
		Name: "Bjorka",
		Email: "bjorka@gmail.com",
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

func TestUsecase_FailedChangePasswordHandler(t *testing.T){
	tests := []struct{
		name string
		request *domains.ChangePassword
		expected error
	}{
		{
			name: "user_fail_1",
			request: &domains.ChangePassword{
				Email: "kalegmail.com",
				NewPassword: "password123",
				PasswordConfirm: "password123",
			},
			expected: errors.New("Email does'n contains @ or ."),
		},
		{
			name: "user_fail_2",
			request: &domains.ChangePassword{
				Email: "kale@gmail.com",
				NewPassword: "password12",
				PasswordConfirm: "password123",
			},
			expected: errors.New("Password not match!"),
		},
		{
			name: "user_fail_3",
			request: &domains.ChangePassword{
				Email: "kale@gmail.com",
				NewPassword: "pass",
				PasswordConfirm: "pass",
			},
			expected: errors.New("Password must be greater than 8 characters!"),
		},
		{
			name: "user_fail_4",
			request: &domains.ChangePassword{
				Email: "kale@gmail.com",
				NewPassword: "password",
				PasswordConfirm: "password",
			},
			expected: errors.New("Email not register!"),
		},
	}

	for _, test := range tests{
		userRepository.Mock.On("FindByEmail", test.request.Email).Return(nil)
		userRepository.Mock.On("UpdatePassword", test.request).Return(errors.New(""))
		err := userUsecase.ChangePasswordHandler(test.request)
		assert.NotNil(t, err)
		assert.Equal(t, err, test.expected)
	}
}

func TestUsecase_SuccessChangePasswordHandler(t *testing.T){
	userInput := &domains.ChangePassword{
		Email: "kale@gmail.com",
		NewPassword: "password",
		PasswordConfirm: "password",
	}
	user1 := repository.User{
		ID: "id",
		Name: "kale",
		Email: "kale@gmail.com",
		Password: "old_password",
	}

	userRepository.Mock.On("FindByEmail", userInput.Email).Return(user1)
	userRepository.Mock.On("UpdatePassword", userInput).Return(nil)
	err := userUsecase.ChangePasswordHandler(userInput)

	assert.Nil(t, err)
}

func TestUsecase_DeleteUserHadler(t *testing.T) {
	userId := "random_uuid"

	t.Run("failed_delete_user", func(t *testing.T) {
		userRepository.Mock.On("DeleteUserById", userId).Return(errors.New(""))

		err := userUsecase.DeleteUserHadler(userId)
		assert.NotNil(t, err)
		
	})

	userId = "real_uuid"
	t.Run("success_delete_user", func(t *testing.T) {
		userRepository.Mock.On("DeleteUserById", userId).Return(nil)

		err := userUsecase.DeleteUserHadler(userId)
		assert.Nil(t, err)
	})
}

func TestUsecase_FailGetUsers(t *testing.T) {
	// type users *[]repository.User
	
	t.Run("fail_get_users", func(t *testing.T) {
		userRepository.Mock.On("Users").Return(nil)
		
		users, err  := userUsecase.GetUsers()
		errExpected := errors.New("Users not found!")

		assert.NotNil(t, err)
		assert.Nil(t, users)
		assert.Equal(t, errExpected, err)
	})

}

func TestUsecase_SuccessGetUsers(t *testing.T) {

	users := []repository.User{
		{
			ID: "uuid1",
			Name: "kale",
			Email: "kale@gmail.com",
			Password: "password",
		},
		{
			ID: "uuid2",
			Name: "kal2",
			Email: "kale2@gmail.com",
			Password: "password",
		},
	}

	t.Run("success_get_users", func(t *testing.T) {
		userRepository.Mock.On("Users").Return(users)
		
		users, err  := userUsecase.GetUsers()

		assert.NotNil(t, users)
		assert.Nil(t, err)
	})
}