package model

import "time"

type Booking struct {
	ID        string    `json:"id" bson:"_id"`
	RequestID string    `json:"requestId" bson:"request_id"`
	UserID    string    `json:"userId" bson:"user_id"`
	EventID   string    `json:"eventId" bson:"event_id"`
	Price     float64   `json:"price" bson:"total_price"`
	Tickets   int       `json:"tickets" bson:"no_of_tickets"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
