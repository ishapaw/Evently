package consumer

import (
	"log"

	"github.com/redis/go-redis/v9"
	"cancel_consumer/kafka"
	"gorm.io/gorm"
)

func StartCancelConsumer(
	broker string,
	topic string,
	groupID string,
	redisReq *redis.Client,
	redisSeats *redis.Client,
	db *gorm.DB,
) {

	reader := kafka.NewReader(broker, topic, groupID)
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v", err)
		}
	}()

	log.Printf("Cancelable Kafka consumer started: topic=%s, groupID=%s", topic, groupID)

	
	log.Println("Cancelable Kafka consumer stopped")
}
