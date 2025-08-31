package main

import (
	"beacon-logs/internal/archiever/service"
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	natsURL     = "nats://localhost:4222"
	natsSubject = "events.raw"

	mongoURL = "mongodb://localhost:27017"
	mongoDB  = "beacon"
	mongoCol = "events"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("connected to mongodb")

	nc, err := nats.Connect(natsURL, nats.PingInterval(20*time.Second), nats.MaxPingsOutstanding(5))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("connected to nats")

	achieverService := service.NewService(nc, client.Database(mongoDB).Collection(mongoCol))

	if err := achieverService.Start(natsSubject); err != nil {
		log.Fatalf("Failed to start archiever service: %v", err)
	}

	select {}
}
