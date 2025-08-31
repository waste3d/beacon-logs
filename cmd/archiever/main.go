package main

import (
	"log"
	"time"

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

	_, err = nc.Subscribe(natsSubject, func(msg *nats.Msg) {
		log.Printf("Получено сообщение из темы '%s': %s\n", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to NATS: %v", err)
	}

	log.Printf("Подписались на тему '%s'. Ожидаем сообщения...", natsSubject)

	select {}
}
