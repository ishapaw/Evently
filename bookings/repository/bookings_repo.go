package repository

import (
	"bookings/model"
)

type BookingRepository struct {
	bookings map[string]*model.Booking
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{
		bookings: make(map[string]*model.Booking),
	}
}
