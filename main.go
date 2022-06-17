package main

import (
	"bwastartup/user"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/bwastart?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var users []user.User
	if err != nil {
		log.Fatal(err.Error())
	}
	db.First(&users)

	for _, user := range users {
		fmt.Println(user.Name)
		fmt.Println(user.Occupation)
	}
}
