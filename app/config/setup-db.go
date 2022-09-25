package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func SetupMysql() *gorm.DB {
	USER := "root"
	PASS := ""
	HOST := "localhost"
	// PORT := "3306"
	DBNAME := "intern_sekolahmu"

	URL := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, DBNAME)

	db, err := gorm.Open("mysql", URL)

	if err != nil {
		panic(err.Error())
	}
	return db
}
