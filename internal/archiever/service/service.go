package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	nc        *nats.Conn
	eventsCol *mongo.Collection
}

func NewService(nc *nats.Conn, eventsCol *mongo.Collection) *Service {
	return &Service{
		nc:        nc,
		eventsCol: eventsCol,
	}
}

func (s *Service) Start(subj string) error {
	_, err := s.nc.Subscribe(subj, s.handleEvent)

	if err != nil {
		return err
	}

	log.Printf("Подписались на тему '%s'. Ожидаем сообщения...", subj)

	select {}
}

func (s *Service) handleEvent(msg *nats.Msg) {
	log.Printf("Получено сообщение из темы '%s': %s\n", msg.Subject, string(msg.Data))

	var eventData bson.M
	if err := json.Unmarshal(msg.Data, &eventData); err != nil {
		log.Printf("Failed to unmarshal event data: %v", err)
		return
	}

	eventData["processed_at"] = time.Now()

	insertCtx, insertCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer insertCancel()

	insertResult, err := s.eventsCol.InsertOne(insertCtx, eventData)
	if err != nil {
		log.Printf("Failed to insert event data: %v", err)
		return
	}
	log.Printf("Документ успешно вставлен в MongoDB с ID: %v", insertResult.InsertedID)
}
