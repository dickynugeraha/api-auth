package main

import (
	"api-auth/app"
	"api-auth/app/config"
	"api-auth/services/logic"
	"api-auth/services/repository"
)


func main() {
	db := config.SetupMysql()
	db.AutoMigrate(&repository.User{})
	
	userRepo := repository.NewUserRepository(db)
	userCase := logic.NewUserUsecase(userRepo)

	r := app.Routes(db, userCase)
	r.Run(":3000")
}