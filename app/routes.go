package app

import (
	"api-auth/controllers"
	"api-auth/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Routes(db *gorm.DB, uc usecase.UserUsecaseInterface) *gin.Engine {

	c := controllers.NewInitController(uc)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	r.POST("/login", c.Login)
	r.POST("/register", c.Register)
	r.POST("/change-password", c.ChangePassword)
	r.GET("/users", c.AllUsers)
	r.GET("/user/:userId", c.SingleUser)
	r.DELETE("/user", c.DeleteUser)

	return r
}
