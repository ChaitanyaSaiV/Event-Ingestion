package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
)

type InMemoryStore struct {
	Mu        sync.RWMutex
	Incidents map[string]models.IncidentData
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Incidents: make(map[string]models.IncidentData),
	}
}

func (i *InMemoryStore) Save(ctx context.Context, incident *models.IncidentData) {
	i.Mu.Lock()
	defer i.Mu.Unlock()
	i.Incidents[incident.Id] = *incident
}

func (i *InMemoryStore) Get(ctx context.Context, id string) (models.IncidentData, error) {
	i.Mu.RLock()
	defer i.Mu.RUnlock()
	data, ok := i.Incidents[id]
	if !ok {
		return models.IncidentData{}, fmt.Errorf("Record not found")
	}
	return data, nil
}

func (i *InMemoryStore) GetAll(ctx context.Context) ([]models.IncidentData, error) {
	i.Mu.RLock()
	defer i.Mu.RUnlock()
	incidents := make([]models.IncidentData, 0, len(i.Incidents))
	for _, incident := range i.Incidents {
		incidents = append(incidents, incident)
	}
	return incidents, nil
}
