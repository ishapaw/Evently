package main

import (
	"context"
	"log"
	"os"
	"time"

	"events/auth"
	"events/controllers"
	"events/repository"
	"events/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	dbHost := getEnv("DB_HOST", "mongodb")
	dbPort := getEnv("DB_PORT", "27017")
	dbUser := getEnv("DB_USER", "admin")
	dbPass := getEnv("DB_PASSWORD", "secret")
	dbName := getEnv("DB_NAME", "eventsdb")

	uri := "mongodb://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort

	var client *mongo.Client
	var err error

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		cancel()

		if err == nil {
			ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
			err = client.Ping(ctx, nil)
			cancel()
			if err == nil {
				log.Println("Connected to MongoDB")
				break
			}
		}

		log.Println("Waiting for MongoDB to be ready...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	db := client.Database(dbName)
	repo := repository.NewEventRepository(db)

	redisClient := newRedisClient("localhost","6379")
	redisSeats := newRedisClient("localhost","6380")

	eventService := service.NewEventService(repo, redisClient, redisSeats)
	eventController := controllers.NewEventController(eventService)

	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.GET("/events/all", eventController.GetAllEvents)
		api.GET("/events/upcoming", eventController.GetAllUpcomingEvents)
		api.GET("/events/:id", eventController.GetEventByID)

		admin := api.Group("/events")
		admin.Use(auth.AdminOnly())
		{
			admin.POST("/create", eventController.CreateEvent)
			admin.PUT("/:id", eventController.UpdateEvent)
		}
	}

	port := getEnv("PORT", "8082")
	log.Println("Events service running on port " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func newRedisClient(host string, port string) *redis.Client {
	addr := host + ":" + port
	pass := getEnv("REDIS_PASSWORD", "")
	db := 0

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Connected to Redis at", addr)
	return rdb
}
