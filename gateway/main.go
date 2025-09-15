package main

import (
	"gateway/kafka"
	"gateway/routes"
	"log"
	"os"
	"time"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("Starting Gateway service...")

	r := gin.Default()

	kafkaBrokers := mustGetEnv("KAFKA_BROKERS")
	producer := kafka.NewProducer(kafkaBrokers)

	if producer == nil {
		log.Fatal("Failed to create Kafka producer")
	} else {
		log.Println("Kafka producer initialized for broker kafka:9092")
	}

	redis := newRedisClient(
		mustGetEnv("REDIS_RATE_LIMITER_HOST"),
		mustGetEnv("REDIS_RATE_LIMITER_PORT"),
		mustGetEnv("REDIS_RATE_LIMITER_PASSWORD"),
	)

	dbHost := mustGetEnv("DB_HOST")
	dbPort := mustGetEnv("DB_PORT")
	dbUser := mustGetEnv("DB_USER")
	dbPass := mustGetEnv("DB_PASSWORD")
	dbName := mustGetEnv("DB_NAME")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Println("Waiting for Postgres to be ready...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	routes.RegisterRoutes(r, producer, redis, db)

	port := mustGetEnv("PORT")
	log.Println("Gateway service running on port " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start Gateway service:", err)
	}

}

func newRedisClient(host, port, pass string) *redis.Client {
	addr := host + ":" + port
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis at", addr)
	return rdb
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
