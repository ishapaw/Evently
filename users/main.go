package main

import (
	"log"
	"time"
	"users/controllers"
	"users/repository"
	"users/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=postgres user=admin password=secret dbname=usersdb port=5432 sslmode=disable"

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Println("Waiting for Postgres to be ready...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	r := gin.Default()

	r.POST("/api/users/register", userController.Register)
	r.POST("/api/users/login", userController.Login)


	r.Run(":8081")
}
