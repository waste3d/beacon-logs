package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type Event struct {
	Event string `json:"event"`
	Url   string `json:"url"`
}

type Service struct {
	natsConn *nats.Conn
	dbPool   *pgxpool.Pool
}

func New(natsConn *nats.Conn, dbPool *pgxpool.Pool) *Service {
	return &Service{
		natsConn: natsConn,
		dbPool:   dbPool,
	}
}

func (s *Service) Run(subj string) error {
	_, err := s.natsConn.Subscribe(subj, s.handleEvent)
	if err != nil {
		return err
	}

	log.Printf("Подписались на тему '%s'. Ожидаем сообщения...", subj)

	select {}

}

func (s *Service) handleEvent(msg *nats.Msg) {
	var event Event
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Failed to unmarshal event data: %v", err)
		return
	}

	tx, err := s.dbPool.Begin(context.Background())
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback(context.Background())

	query := `
	        INSERT INTO page_views (url, view_count, last_seen) 
        VALUES ($1, 1, $2)
        ON CONFLICT (url) 
        DO UPDATE SET 
            view_count = page_views.view_count + 1, 
            last_seen = $2;`

	_, err = tx.Exec(context.Background(), query, event.Url, time.Now())
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
	}

	log.Printf("Event processed: %+v", event)

}
