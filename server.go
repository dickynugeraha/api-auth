package main

import (
	"01_REST_Auth/app"
	"01_REST_Auth/app/config"
	"01_REST_Auth/services/repository"
	"01_REST_Auth/usecase"
)


func main() {
	db := config.SetupMysql()
	db.AutoMigrate(&repository.User{})
	
	userRepo := repository.NewUserRepository(db)
	userCase := usecase.NewUserUsecase(userRepo)

	r := app.Routes(db, userCase)
	r.Run(":3000")
}


