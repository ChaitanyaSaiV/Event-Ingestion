package store

import (
	"context"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
)

// IncidentStore is the shared interface — single source of truth.
type IncidentStore interface {
	Reader
	Writer
}

type Reader interface {
	Get(ctx context.Context, id string) (models.IncidentData, error)
	GetAll(ctx context.Context) ([]models.IncidentData, error)
}

type Writer interface {
	Save(ctx context.Context, incident *models.IncidentData)
	Delete(ctx context.Context, id string) error
}
