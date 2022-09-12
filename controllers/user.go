package controllers

import (
	"01_REST_Auth/domains"
	repo "01_REST_Auth/services/repository"
	caseUser "01_REST_Auth/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context){
	var inputRegis domains.Register

	if err := c.ShouldBindJSON(&inputRegis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}

	statusHttp, err := caseUser.RegisterHandler(&inputRegis)
	if err != nil {
		c.JSON(statusHttp, gin.H{
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

func Login(c *gin.Context){
	var inputLogin domains.Login

	if err := c.ShouldBindJSON(&inputLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}

	user, token, err := caseUser.LoginHandler(&inputLogin)
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

func ChangePassword(c *gin.Context){
	var inputChangePass domains.ChangePassword

	if err := c.ShouldBindJSON(&inputChangePass); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
	}
	_, err := repo.FindbyEmail(inputChangePass.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Email not register!",
		})
		return
	}
	err = caseUser.PasswordRequired(inputChangePass.NewPassword, inputChangePass.PasswordConfirm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}
	hashedPassword := caseUser.PasswordHashing(inputChangePass.NewPassword)
	_ = repo.UpdatePassword(inputChangePass.Email, hashedPassword)

	c.JSON(http.StatusOK, gin.H{
		"message" : "Password changed!",
	})
}


func AllUsers(c *gin.Context){
	users, err := caseUser.AllUsers()
	
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

func SingleUser(c *gin.Context){
	userId := c.Param("user_id")
	user, err := repo.GetUserById(userId)
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

func DeleteUser(c *gin.Context){
	var inputId domains.UserId

	err := c.ShouldBindJSON(&inputId);
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
	}

	err = repo.DeleteUserById(inputId.ID)
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