package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
)

var ErrNotFound = fmt.Errorf("Incident Record Not Found")

type InMemoryStore struct {
	mu        sync.RWMutex
	incidents map[string]models.IncidentData
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		incidents: make(map[string]models.IncidentData),
	}
}

func (i *InMemoryStore) Save(ctx context.Context, incident *models.IncidentData) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.incidents[incident.Id] = *incident
}

func (i *InMemoryStore) Get(ctx context.Context, id string) (models.IncidentData, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	data, ok := i.incidents[id]
	if !ok {
		return models.IncidentData{}, ErrNotFound
	}
	return data, nil
}

func (i *InMemoryStore) GetAll(ctx context.Context) ([]models.IncidentData, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	incidents := make([]models.IncidentData, 0, len(i.incidents))
	for _, incident := range i.incidents {
		incidents = append(incidents, incident)
	}
	return incidents, nil
}
