package controllers

import (
	"01_REST_Auth/domains"
	"01_REST_Auth/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	caseUser usecase.UserUsecaseInterface
}

func NewInitController(caseUser usecase.UserUsecaseInterface) *AuthController {
	return &AuthController{
		caseUser: caseUser,
	}
}

func (ac *AuthController) Register(c *gin.Context){
	var inputRegis domains.Register

	if err := c.ShouldBindJSON(&inputRegis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}

	err := ac.caseUser.RegisterHandler(&inputRegis)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message" : "Successfully created user!",
		"name" : inputRegis.Name,
		"email" : inputRegis.Email,
	})
}

func (ac *AuthController) Login(c *gin.Context){
	var inputLogin domains.Login

	if err := c.ShouldBindJSON(&inputLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}

	user, token, err := ac.caseUser.LoginHandler(&inputLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successfully!",
		"user" : user,
		"token": token,
	})
}

func (ac *AuthController) ChangePassword(c *gin.Context){
	var inputChangePass domains.ChangePassword

	if err := c.ShouldBindJSON(&inputChangePass); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
	}
	
	err := ac.caseUser.ChangePasswordHandler(&inputChangePass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "Password changed!",
	})
}

func (ac *AuthController) AllUsers(c *gin.Context){
	fmt.Println("Inner controller")
	users, err := ac.caseUser.GetUsers()
	
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error" : err.Error(),
		})
		return 
	}
	c.JSON(http.StatusOK, gin.H{
		"message" : "Successflly fetch all users",
		"users" : users,
	})
} 

func (ac *AuthController) SingleUser(c *gin.Context){
	userId := c.Param("user_id")
	user, err := ac.caseUser.GetSingleUserHandler(userId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error" : err.Error(),
		})
		return 
	}
	c.JSON(http.StatusOK, gin.H{
		"message" : "Successflly fetch single user",
		"user" : user,
	})
}

func (ac *AuthController) DeleteUser(c *gin.Context){
	var inputId domains.UserId

	err := c.ShouldBindJSON(&inputId);
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
	}

	err = ac.caseUser.DeleteUserHadler(inputId.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})
		return 
	}
	c.JSON(http.StatusOK, gin.H{
		"message" : "Successflly delete user",
	})
}