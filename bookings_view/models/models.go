package models

import (
	"time"
)

type Booking struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RequestID string    `gorm:"type:varchar(255);not null" json:"requestId"`
	UserID    string    `gorm:"type:varchar(255);not null" json:"userId"`
	EventID   string    `gorm:"type:varchar(255);not null" json:"eventId"`
	Price     float64   `gorm:"type:numeric;not null" json:"price"`
	Tickets   int64     `gorm:"not null" json:"tickets"`
	Status    string    `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
