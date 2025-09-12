package service

import (
	"bookings_view/models"
	"bookings_view/repository"
	"errors"
)

type BookingsViewService interface {
	GetBookingByID(id string) (*models.Booking, error)
	GetBookingsByEventID(eventID string, limit, page int64) ([]models.Booking, error)
	GetBookingsByUserID(userID string, limit, page int64) ([]models.Booking, error)
}

type bookingsViewService struct {
	repo repository.BookingsViewRepository
}

func NewBookingsViewService(repo repository.BookingsViewRepository) BookingsViewService {
	return &bookingsViewService{repo: repo}
}

func (s *bookingsViewService) GetBookingByID(id string) (*models.Booking, error) {
	booking, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, errors.New("booking not found")
	}

	return booking, nil
}

func (s *bookingsViewService) GetBookingsByEventID(eventID string, limit, page int64) ([]models.Booking, error) {
	return s.repo.GetByEventID(eventID, limit, page)
}

func (s *bookingsViewService) GetBookingsByUserID(userID string, limit, page int64) ([]models.Booking, error) {
	return s.repo.GetByUserID(userID, limit, page)
}
