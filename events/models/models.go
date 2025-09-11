package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID               primitive.ObjectID     `bson:"_id,omitempty" json:"_id"`
	Title            string                 `bson:"title" json:"title"`
	Description      string                 `bson:"description,omitempty" json:"description"`
	Venue            string                 `bson:"venue" json:"venue"`
	Date             time.Time             `bson:"date" json:"date"`
	Price            float64                `bson:"price" json:"price"`
	AvailableTickets int64                   `bson:"available_tickets" json:"available_tickets"`
	TotalTickets     int64                   `bson:"total_tickets" json:"total_tickets"`
	Metadata         map[string]interface{} `bson:",inline"`
	CreatedAt        time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time              `bson:"updated_at" json:"updated_at"`
}


type UpcomingEvent struct {
	ID               primitive.ObjectID     `bson:"_id,omitempty" json:"_id"`
	Title            string                 `bson:"title" json:"title"`
	Venue            string                 `bson:"venue" json:"venue"`
	Date             time.Time             `bson:"date" json:"date"`
	AvailableTickets int64                   `bson:"available_tickets" json:"available_tickets"`
	TotalTickets     int64                   `bson:"total_tickets" json:"total_tickets"`
}

