package storage

import (
	"context"
	"log"
	"time"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/store"
)

type LoggingStore struct {
	next store.IncidentStore
}

func NewLoggingStore(next store.IncidentStore) *LoggingStore {
	return &LoggingStore{next: next}
}

func (l *LoggingStore) Save(ctx context.Context, incident *models.IncidentData) {
	start := time.Now()
	l.next.Save(ctx, incident)
	log.Printf("Save id=%s took=%v", incident.Id, time.Since(start))
}
func (l *LoggingStore) Get(ctx context.Context, id string) (models.IncidentData, error) {
	start := time.Now()
	incident, err := l.next.Get(ctx, id)
	log.Printf("Get id=%s err=%v took=%v", id, err, time.Since(start))
	return incident, err
}

func (l *LoggingStore) GetAll(ctx context.Context) ([]models.IncidentData, error) {
	start := time.Now()
	incidents, err := l.next.GetAll(ctx)
	log.Printf("GetAll count=%d err=%v took=%v", len(incidents), err, time.Since(start))
	return incidents, err
}

func (l *LoggingStore) Delete(ctx context.Context, id string) error {
	start := time.Now()
	err := l.next.Delete(ctx, id)
	log.Printf("Delete id=%s err=%v took=%v", id, err, time.Since(start))
	return err
}
