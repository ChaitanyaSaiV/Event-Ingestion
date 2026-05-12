package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
)

type InMemoryStore struct {
	mu            sync.RWMutex
	incident_data map[string]models.IncidentData
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{incident_data: make(map[string]models.IncidentData)}
}

func (s *InMemoryStore) Save(ctx context.Context, i *models.IncidentData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.incident_data[i.Id] = *i
}

func (s *InMemoryStore) Get(ctx context.Context, id string) (models.IncidentData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.incident_data[id]
	if !ok {
		return models.IncidentData{}, fmt.Errorf("Record not found")
	}

	return data, nil

}
