package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// 1. Define the interface your handler needs
type IncidentStore interface {
	Save(ctx context.Context, incident *models.IncidentData)
	Get(ctx context.Context, id string) (models.IncidentData, error)
	GetAll(ctx context.Context) ([]models.IncidentData, error)
}

// 2. Create the handler struct that holds the database
type IncidentHandler struct {
	store IncidentStore
}

// 3. Create a constructor to initialize the handler
func NewIncidentHandler(s IncidentStore) *IncidentHandler {
	return &IncidentHandler{
		store: s,
	}
}

func (h *IncidentHandler) SaveIncident(w http.ResponseWriter, r *http.Request) {
	var req models.CreateIncidentRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Passed unsupported JSON data"))
		return
	}

	err = validate.Struct(req)

	incident := models.IncidentData{
		Id:        req.Id,
		Severity:  req.Severity,
		Message:   req.Message,
		TimeStamp: time.Now().UTC(),
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing Required Fields"))
	}

	h.store.Save(r.Context(), &incident)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Saved the Record"))

}

func (h *IncidentHandler) GetIncident(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := h.store.Get(r.Context(), id)
	if err != nil {

		http.Error(w, "No record found with supplied ID", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(&data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *IncidentHandler) GetAllIncidents(w http.ResponseWriter, r *http.Request) {
	data, err := h.store.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
