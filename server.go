package main

import (
	"api-auth/app"
	"api-auth/app/config"
	"api-auth/services/repository"
	"api-auth/usecase"
)


func main() {
	db := config.SetupMysql()
	db.AutoMigrate(&repository.User{})
	
	userRepo := repository.NewUserRepository(db)
	userCase := usecase.NewUserUsecase(userRepo)

	r := app.Routes(db, userCase)
	r.Run(":3000")
}