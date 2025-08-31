package main

import (
	"beacon-logs/internal/aggregator/service"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

const (
	natsURL     = "nats://localhost:4222"
	natsSubject = "events.raw"
	postgresURL = "postgres://waste3d:waste3d@localhost:5432/beacon"
)

func main() {
	nc, err := nats.Connect(natsURL, nats.PingInterval(20*time.Second), nats.MaxPingsOutstanding(5))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("connected to nats")

	dbPool, err := pgxpool.New(context.Background(), postgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	log.Println("connected to postgresql")

	aggregatorService := service.New(nc, dbPool)

	if err := aggregatorService.Run(natsSubject); err != nil {
		log.Fatalf("Failed to start aggregator service: %v", err)
	}
}
