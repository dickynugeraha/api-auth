package app

import (
	"01_REST_Auth/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context){
		c.Set("db", db)
	})

	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	r.POST("/change-password", controllers.ChangePassword)
	r.GET("/users", controllers.AllUsers)
	r.GET("/user/:user_id", controllers.SingleUser)
	r.DELETE("/user", controllers.DeleteUser)
	
	return r
}