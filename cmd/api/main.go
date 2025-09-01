package main

import (
	"beacon-logs/internal/api"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	postgresURL = "postgres://waste3d:waste3d@localhost:5432/beacon"
	apiPort     = ":8081"
)

func main() {
	dbPool, err := pgxpool.New(context.Background(), postgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	apiService := api.NewService(dbPool)

	if err := apiService.Run(apiPort); err != nil {
		log.Fatalf("Failed to start API service: %v", err)
	}
}
