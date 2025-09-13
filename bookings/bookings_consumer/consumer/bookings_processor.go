package consumer

import (
	"bookings_consumer/models"
	"context"
	"encoding/json"
	"log"
	"strconv"
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


func processBookingMessage(ctx context.Context, value []byte, deps *models.ProcessorDeps){
	var req models.KafkaEvent
	if err := json.Unmarshal(value, &req); err != nil {
		log.Printf("Invalid booking message: %v", err)
		return
	}

	reqKey := "req:" + req.RequestID

	existing, _ := deps.RedisReq.Get(ctx, reqKey).Result()

	if existing != "" {
		req.State = existing
	} else {
		req.State = "state1"
		saveState(ctx, deps.RedisReq, reqKey, req.State)
	}

	switch req.State {

	case "state1":
		stateHandlerFunc1(ctx, req, deps)

	case "state2":
		stateHandlerFunc2(ctx, req, deps)

	case "state3":
		stateHandlerFunc3(ctx,req, deps)

	case "failed":
		log.Printf("Request %s already failed", req.RequestID)

	case "success":
		log.Printf("Request %s already succeeded", req.RequestID)

	case "cancelled":
		log.Printf("Request %s is already cancelled", req.RequestID)
	}
}

func stateHandlerFunc1(ctx context.Context, req models.KafkaEvent, deps *models.ProcessorDeps) {
	reqKey := "req:" + req.RequestID
	seatsKey := "event:" + req.EventID

	if isCancelled(ctx, deps.RedisReq, reqKey) {
		log.Printf("Request %s was cancelled before seat allocation", req.RequestID)
		return
	}


	result, err := decrSeatsScript.Run(ctx, deps.RedisSeats, []string{seatsKey}, req.Seats).Int()
	if err != nil {
		log.Printf("Redis error: %v", err)
		return
	}

	if result == -1 {
		log.Printf("Request %s failed: event not found", req.RequestID)
		req.State = "failed"
		saveState(ctx, deps.RedisReq, reqKey, req.State)
		return
	}

	if result == 0 {
		log.Printf("Request %s failed: not enough seats", req.RequestID)
		req.State = "failed"
		saveState(ctx, deps.RedisReq, reqKey, req.State)
		return
	}

	req.State = "state2"
	saveState(ctx, deps.RedisReq, reqKey, req.State)
	stateHandlerFunc2(ctx, req, deps)
}


func stateHandlerFunc2(ctx context.Context, req models.KafkaEvent, deps *models.ProcessorDeps) {
	reqKey := "req:" + req.RequestID
	seatsKey := "event:" + req.EventID

	if isCancelled(ctx, deps.RedisReq, reqKey) {
		// revert seats
		deps.RedisSeats.IncrBy(ctx, seatsKey, int64(req.Seats))
		log.Printf("Request %s cancelled during processing, seats reverted", req.RequestID)

		saveState(ctx, deps.RedisReq, reqKey, "cancelled")
		return
	}

	priceKey := "event:" + req.EventID
	priceStr, _ := deps.RedisPrice.Get(ctx, priceKey).Result()

	price, _ := strconv.ParseFloat(priceStr, 64)

	err := deps.DB.Transaction(func(tx *gorm.DB) error {
		booking := models.Booking{
			RequestID: req.RequestID,
			EventID:   req.EventID,
			UserID:    req.UserID,
			Price:     price * float64(req.Seats),
			Seats:     req.Seats,
		}

		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&booking).Error
	})
	if err != nil {
		log.Printf("DB error: %v", err)
		return
	}

	req.State = "state3"
	saveState(ctx, deps.RedisReq, reqKey, req.State)
	stateHandlerFunc3(ctx, req, deps)
}

func stateHandlerFunc3(ctx context.Context, req models.KafkaEvent, deps *models.ProcessorDeps) {
	reqKey := "req:" + req.RequestID
	seatsKey := "event:" + req.EventID

	if isCancelled(ctx, deps.RedisReq, reqKey) {
		// update booking status + restore seats
		if err := deps.DB.Model(&models.Booking{}).Where("request_id = ?", req.RequestID).Update("status","cancelled").Error; err != nil {
			log.Printf("Error deleting cancelled booking %s: %v", req.RequestID, err)
			return
		}

		deps.RedisSeats.IncrBy(ctx, seatsKey, int64(req.Seats))
		log.Printf("Request %s cancelled after DB insert, booking deleted + seats reverted", req.RequestID)
		saveState(ctx, deps.RedisReq, reqKey, "cancelled")
		return
	}

	err := publishSeatsUpdate(deps.Producer, req)
	if err != nil {
		log.Printf("Kafka error: %v", err)
		return
	}

	req.State = "success"
	saveState(ctx, deps.RedisReq, reqKey, req.State)
	log.Printf("Request %s processed successfully", req.RequestID)
}

func saveState(ctx context.Context, rdb *redis.Client, key string, state string) {
	rdb.Set(ctx, key, state, 5*time.Minute)
}

func publishSeatsUpdate(producer *kafka.Producer, req models.KafkaEvent) error {
	event := models.KafkaUpdateEvent{
		EventId: req.EventID,
		Seats:   req.Seats,
		Operation: "subtract",
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal seats update event: %v", err)
		return err
	}

	err = producer.Publish(
		"update.seats",
		[]byte(req.RequestID),
		payload,
	)

	if err != nil {
		log.Printf("Failed to publish seats update event: %v", err)
	} else {
		log.Printf("Published seats update for event %s: %d seats left", req.EventID, req.Seats)
	}

	return err
}


func isCancelled(ctx context.Context, rdb *redis.Client, key string) bool {
	state, _ := rdb.Get(ctx, key).Result()
	return state == "cancelled"
}


