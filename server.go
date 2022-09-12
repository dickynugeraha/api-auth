package main

import (
	"01_REST_Auth/app"
	"01_REST_Auth/app/config"
	"01_REST_Auth/services/repository"
)

func main () {
	db := config.SetupMysql()
	db.AutoMigrate(&repository.User{})

	r := app.Routes(db)
	r.Run(":3000")
}


