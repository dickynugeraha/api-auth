package controllers

import (
	"api-auth/domains"
	mokz "api-auth/mock"
	"api-auth/services/repository"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var userUsecase = &mokz.UserUsecaseMock{Mock: mock.Mock{}}
var userController = AuthController{caseUser: userUsecase}

func SetRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func TestSuccessRegister(t *testing.T) {
	input := domains.Register{
		Name:            "kale",
		Email:           "kale@gmail.com",
		Password:        "password",
		PasswordConfirm: "password",
	}

	mockResponse := `{"email":"%s","message":"Successfully created user!","name":"%s"}`
	mockResponse = fmt.Sprintf(mockResponse, input.Email, input.Name)

	r := SetRouter()
	r.POST("/register", userController.Register)

	userUsecase.Mock.On("RegisterHandler", &input).Return(nil)

	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, mockResponse, string(res))
}

func TestFailRegister(t *testing.T) {
	r := SetRouter()
	r.POST("/register", userController.Register)

	t.Run("fail_case1", func(t *testing.T) {
		input := domains.Register{
			Name:     "kale",
			Email:    "kale@gmail.com",
			Password: "password",
		}
		jsonValue, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("fail_case2", func(t *testing.T) {
		input := domains.Register{
			Name:            "kale",
			Email:           "kale@gmail.com",
			Password:        "password",
			PasswordConfirm: "passwords",
		}

		userUsecase.Mock.On("RegisterHandler", &input).Return(errors.New(""))

		jsonValue, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSuccessLogin(t *testing.T) {
	user := &repository.User{
		ID:       "uuid",
		Name:     "kale",
		Email:    "kale@gmail.com",
		Password: "password",
	}

	input := domains.Login{
		Name:     "kale",
		Email:    "kale@gmail.com",
		Password: "passwords",
	}

	mockResponse := `{"message":"Login successfully!","token":"valid_token","user":{"id":"%s","name":"%s","email":"%s","password":"%s"}}`
	mockResponse = fmt.Sprintf(mockResponse, user.ID, user.Name, user.Email, user.Password)

	r := SetRouter()
	r.POST("/login", userController.Login)

	userUsecase.Mock.On("LoginHandler", &input).Return(user, "valid_token", nil)

	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockResponse, string(res))
}

func TestFailLogin(t *testing.T) {

	r := SetRouter()
	r.POST("/login", userController.Login)

	t.Run("fail_case1", func(t *testing.T) {
		input := domains.Login{
			Name:  "kale",
			Email: "kale@gmail.com",
		}

		userUsecase.Mock.On("LoginHandler", &input).Return(nil, "", errors.New(""))

		jsonValue, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("fail_case2", func(t *testing.T) {
		input := domains.Login{
			Name:     "user1",
			Email:    "kalegmailcom",
			Password: "password",
		}

		userUsecase.Mock.On("LoginHandler", &input).Return(nil, "", errors.New(""))

		jsonValue, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSuccessChangePassword(t *testing.T) {
	mockResponse := `{"message":"Password changed!"}`

	input := domains.ChangePassword{
		Email:           "kale@gmail.com",
		NewPassword:     "passwords",
		PasswordConfirm: "passwords",
	}

	r := SetRouter()
	r.POST("/change-password", userController.ChangePassword)

	userUsecase.Mock.On("ChangePasswordHandler", &input).Return(nil)

	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/change-password", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, string(res), mockResponse)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFailChangePassword(t *testing.T) {

	input := domains.ChangePassword{
		Email:           "kale@gmail.com",
		NewPassword:     "passwor",
		PasswordConfirm: "passwords",
	}

	r := SetRouter()
	r.POST("/change-password", userController.ChangePassword)

	userUsecase.Mock.On("ChangePasswordHandler", &input).Return(errors.New(""))

	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/change-password", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestSuccessDeleteUser(t *testing.T) {
	mockResponse := `{"message":"Successflly delete user"}`
	input := domains.UserId{
		ID: "valid_uuid",
	}

	r := SetRouter()
	r.GET("/user")
	r.DELETE("/user", userController.DeleteUser)

	userUsecase.Mock.On("DeleteUserHandler", input.ID).Return(nil)

	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, string(res), mockResponse)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFailDeleteUser(t *testing.T) {
	input := domains.UserId{
		ID: "invalid_uuid",
	}

	r := SetRouter()
	r.DELETE("/user", userController.DeleteUser)

	userUsecase.Mock.On("DeleteUserHandler", input.ID).Return(errors.New(""))

	jsonValue, _ := json.Marshal(input)
	req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSuccessGetAllUsers(t *testing.T) {

	users := []repository.User{
		{
			ID:       "uuid1",
			Name:     "user1",
			Email:    "email1",
			Password: "password1",
		},
		{
			ID:       "uuid2",
			Name:     "user2",
			Email:    "email2",
			Password: "password2",
		},
	}

	var usersString []string
	for _, v := range users {
		usersString = append(usersString, fmt.Sprintf(`{"id":"%s",`, v.ID), fmt.Sprintf(`"name":"%s",`, v.Name), fmt.Sprintf(`"email":"%s",`, v.Email), fmt.Sprintf(`"password":"%s"},`, v.Password))
	}

	mockResponse := `{"message":"Successflly fetch all users","users":%s}`
	mockResponse = fmt.Sprintf(mockResponse, (usersString))
	mockResponse = strings.ReplaceAll(mockResponse, " ", "")
	mockResponse = strings.ReplaceAll(mockResponse, ",", "")

	r := SetRouter()
	r.GET("/users", userController.AllUsers)

	userUsecase.Mock.On("GetUsers").Return(users, nil)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res, _ := ioutil.ReadAll(w.Body)
	resManip := strings.ReplaceAll(string(res), " ", "")
	resManip = strings.ReplaceAll(resManip, ",", "")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockResponse, resManip)
}

func TestFailGetAllUsers(t *testing.T) {
	r := SetRouter()
	r.GET("/users", userController.AllUsers)
	userUsecase.Mock.On("GetUsers").Return(nil, errors.New("Error"))

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSuccessGetSingleUser(t *testing.T) {
	userId := domains.UserId{
		ID: "valid_uuid",
	}

	r := SetRouter()
	r.GET("/user/:userId", userController.SingleUser)

	user := &repository.User{
		ID:       "valid_uuid",
		Name:     "kale",
		Email:    "kale@gmail.com",
		Password: "password",
	}

	userUsecase.Mock.On("GetSingleUserHandler", userId.ID).Return(user, nil)

	mockResponse := `{"message":"Successflly fetch single user","user":{"id":"%s","name":"%s","email":"%s","password":"%s"}}`
	mockResponse = fmt.Sprintf(mockResponse, user.ID, user.Name, user.Email, user.Password)

	jsonValue, _ := json.Marshal(userId)
	req, _ := http.NewRequest("GET", "/user/"+userId.ID, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	res, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(res), mockResponse)

}

func TestFailGetSingleUser(t *testing.T) {
	userId := domains.UserId{
		ID: "invalid_uuid",
	}

	r := SetRouter()
	r.GET("/user/:userId", userController.SingleUser)

	userUsecase.Mock.On("GetSingleUserHandler", userId.ID).Return(nil, errors.New(""))

	jsonValue, _ := json.Marshal(userId)
	req, _ := http.NewRequest("GET", "/user/"+userId.ID, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

}
