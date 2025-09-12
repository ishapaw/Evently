package service

import (
	"bookings/model"
	"bookings/repository"
)

type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) CreateBooking(booking *model.Booking) error {
	return s.repo.Save(booking)
}

func (s *BookingService) GetBooking(id string) (*model.Booking, error) {
	return s.repo.Get(id)
}
