package consumer

import (
	"bookings_consumer/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"bookings_consumer/kafka"
	"gorm.io/gorm/clause"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var decrSeatsScript = redis.NewScript(`
local available = redis.call("GET", KEYS[1])
if not available then
    return -1
end
available = tonumber(available)
local required = tonumber(ARGV[1])
if available >= required then
    redis.call("DECRBY", KEYS[1], required)
    return 1
else
    return 0
end
`)

func processBookingMessage(ctx context.Context, value []byte, redisReq *redis.Client, redisSeats *redis.Client, db *gorm.DB,  producer *kafka.Producer ) {
	var req models.KafkaEvent
	if err := json.Unmarshal(value, &req); err != nil {
		log.Printf("Invalid booking message: %v", err)
		return
	}

	reqKey := "req:" + req.RequestID

	existing, _ := redisReq.Get(ctx, reqKey).Result()
	if existing != "" {
		req.State = existing 
	} else {
		req.State = "state1"
		saveState(ctx, redisReq, reqKey, req.State)
	}

	switch req.State {

	case "state1":
		seatsKey := "event:" + req.EventID

		result, err := decrSeatsScript.Run(ctx, redisSeats, []string{seatsKey}, req.Seats).Int()
		if err != nil {
			log.Printf("Redis error: %v", err)
			return
		}

		if result == -1 {
			log.Printf("Request %s failed: event not found", req.RequestID)
			req.State = "failed"
			saveState(ctx, redisReq, reqKey, req.State)
			return
		}

		if result == 0 {
			log.Printf("Request %s failed: not enough seats", req.RequestID)
			req.State = "failed"
			saveState(ctx, redisReq, reqKey, req.State)
			return
		}

		req.State = "state2"
		saveState(ctx, redisReq, reqKey, req.State)
		fallthrough

	case "state2":

		err := db.Transaction(func(tx *gorm.DB) error {
			booking := models.Booking{
				RequestID: req.RequestID,
				EventID:   req.EventID,
				UserID:    req.UserID,
				Price:     float64(100 * req.Seats),
				Seats:     req.Seats,
			}

			return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&booking).Error
		})
		if err != nil {
			log.Printf("DB error: %v", err)
			return
		}

		req.State = "state3"
		saveState(ctx, redisReq, reqKey, req.State)
		fallthrough

	case "state3":
		err := publishSeatsUpdate(producer, req)
		if err != nil {
			log.Printf("Kafka error: %v", err)
			return
		}

		req.State = "success"
		saveState(ctx, redisReq, reqKey, req.State)
		log.Printf("Request %s processed successfully", req.RequestID)

	case "failed":
		log.Printf("Request %s already failed", req.RequestID)

	case "success":
		log.Printf("Request %s already succeeded", req.RequestID)
	}
}

func saveState(ctx context.Context, rdb *redis.Client, key string, state string) {
	rdb.Set(ctx, key, state, 5*time.Minute)
}

func publishSeatsUpdate(producer *kafka.Producer, req models.KafkaEvent) error {
	event := models.KafkaEvent{
		EventID:   req.EventID,
		Seats:     req.Seats,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal seats update event: %v", err)
		return err
	}

	err = producer.Publish(
		"update_seats",
		[]byte(req.RequestID), 
		payload,         
	)

	if err != nil {
		log.Printf("Failed to publish seats update event: %v", err)
	} else {
		log.Printf("Published seats update for event %s: %d seats left",  req.EventID, req.Seats)
	}

	return err
}
