package main

import (
	"gateway/kafka"
	"gateway/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Gateway service...")

	r := gin.Default()

	kafkaBrokers := mustGetEnv("KAFKA_BROKER")
	producer := kafka.NewProducer(kafkaBrokers)

	if producer == nil {
		log.Fatal("Failed to create Kafka producer")
	} else {
		log.Println("Kafka producer initialized for broker kafka:9092")
	}

	routes.RegisterRoutes(r, producer)

	port := mustGetEnv("PORT_GATEWAY")
	log.Println("Gateway service running on port " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start Gateway service:", err)
	}

}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
