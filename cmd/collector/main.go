package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

const (
	natsURL     = "nats://localhost:4222"
	natsSubject = "events.raw"
)

func main() {
	nc, err := nats.Connect(natsURL, nats.PingInterval(20*time.Second), nats.MaxPingsOutstanding(5))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("connected to nats")

	router := gin.Default()

	router.POST("/track", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}

		if err := nc.Publish(natsSubject, body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event"})
			return
		}

		log.Printf("Сообщение опубликовано в тему '%s': %s", natsSubject, string(body))
		c.JSON(200, gin.H{"status": "ok"})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start collector: %v", err)
	}
}
