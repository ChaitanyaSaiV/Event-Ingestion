package handlers

import (
	"context"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
)

type mockStore2 struct {
	getErr           error
	saveErr          error
	gettAllIncidents []models.IncidentData
	getIncident      models.IncidentData
	saved            models.IncidentData
	allErr           error
}

func (m *mockStore2) Get(ctx context.Context, id string) (models.IncidentData, error) {
	return m.getIncident, m.getErr
}

func (m *mockStore2) Save(ctx context.Context, incident *models.IncidentData) {
	m.saved = *incident
}

func (m *mockStore2) GetAll(ctx context.Context) ([]models.IncidentData, error) {
	return m.gettAllIncidents, m.allErr
}

var _ IncidentStore = (*mockStore2)(nil)
