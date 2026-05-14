package router

import (
	"net/http"
	"time"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/handlers"
)

func NewServer(path string, handler *handlers.IncidentHandler) *http.Server {
	r := http.NewServeMux()
	r.HandleFunc("GET /healthCheck", handlers.HealthCheck)
	r.HandleFunc("POST /incidents", handler.SaveIncident)
	r.HandleFunc("GET /incidents/{id}", handler.GetIncident)
	r.HandleFunc("GET /incidents", handler.GetAllIncidents)
	r.HandleFunc("DELETE /incidents/{id}", handler.DeleteIncident)

	return &http.Server{
		Addr:         path,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
