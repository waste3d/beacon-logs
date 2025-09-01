package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PageViewStat struct {
	URL       string    `json:"url"`
	ViewCount int       `json:"view_count"`
	LastSeen  time.Time `json:"last_seen"`
}

type Service struct {
	dbPool *pgxpool.Pool
}

func NewService(dbPool *pgxpool.Pool) *Service {
	return &Service{
		dbPool: dbPool,
	}
}

func (s *Service) Run(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Наш единственный эндпоинт
	router.GET("/stats/pageviews", s.handleGetPageViews)

	log.Printf("API сервер запущен на %s", addr)
	return router.Run(addr)
}

func (s *Service) handleGetPageViews(c *gin.Context) {
	stats, err := s.handleGetPageViewsFromDB(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get page views"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (s *Service) handleGetPageViewsFromDB(c *gin.Context) ([]PageViewStat, error) {
	query := `        
	SELECT url, view_count, last_seen 
    FROM page_views 
    ORDER BY view_count DESC;`

	rows, err := s.dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []PageViewStat

	for rows.Next() {
		var stat PageViewStat
		if err := rows.Scan(&stat.URL, &stat.ViewCount, &stat.LastSeen); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	if stats == nil {
		stats = make([]PageViewStat, 0)
	}

	return stats, nil
}
