package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
)

type App struct {
	Incidents IncidentStore
}

type IncidentStore interface {
	Get(ctx context.Context, id string) (models.IncidentData, error)
	Save(ctx context.Context, i *models.IncidentData)
}

func (a *App) GetIncident(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Print(id)
	data, err := a.Incidents.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	return_data, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(return_data)
}

func (a *App) PostIncident(w http.ResponseWriter, r *http.Request) {

	var data models.IncidentData

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	a.Incidents.Save(r.Context(), &data)

	w.WriteHeader(http.StatusCreated) // 201 Created
	w.Write([]byte("Incident created successfully\n"))
}
