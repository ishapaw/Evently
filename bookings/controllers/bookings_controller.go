package controller

import (
	"bookings/service"
)

type BookingController struct {
	service *service.BookingService
}

func NewBookingController(service *service.BookingService) *BookingController {
	return &BookingController{service: service}
}
