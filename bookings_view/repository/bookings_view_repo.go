package repository

import (
	"bookings_view/models"
	"errors"

	"gorm.io/gorm"
)

type BookingsViewRepository interface {
	GetByID(id string) (*models.Booking, error)
	GetByEventID(eventID string, limit, page int64) ([]models.Booking, error)
	GetByUserID(userID string, limit, page int64) ([]models.Booking, error)
}

type bookingsViewRepository struct {
	db *gorm.DB
}

func NewBookingsViewRepository(db *gorm.DB) BookingsViewRepository {
	return &bookingsViewRepository{db}
}

func (r *bookingsViewRepository) GetByID(id string) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.First(&booking, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	return &booking, nil
}

func (r *bookingsViewRepository) GetByEventID(eventID string, limit, page int64) ([]models.Booking, error) {
	var bookings []models.Booking
	offset := int((page - 1) * limit)
	if err := r.db.Where("event_id = ?", eventID).
		Limit(int(limit)).
		Offset(offset).
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingsViewRepository) GetByUserID(userID string, limit, page int64) ([]models.Booking, error) {
	var bookings []models.Booking
	offset := int((page - 1) * limit)
	if err := r.db.Where("user_id = ?", userID).
		Limit(int(limit)).
		Offset(offset).
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}
