package main

import (
	"bookings_view/controllers"
	"bookings_view/repository"
	"bookings_view/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=postgres-bookings user=admin password=secret dbname=bookingsdb port=5432 sslmode=disable"

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

	bookingRepo := repository.NewBookingsViewRepository(db)
	bookingService := service.NewBookingsViewService(bookingRepo)
	bookingController := controllers.NewBookingsViewController(bookingService)

	r := gin.Default()

	r.GET("/api/v1/bookings/:id", bookingController.GetBookingByID)

	r.GET("/api/v1/bookings/user/:user_id", bookingController.GetBookingsByUserID)

	r.GET("/api/v1/bookings/event/:event_id", bookingController.GetBookingsByEventID)
	r.Run(":8084")
}
	
