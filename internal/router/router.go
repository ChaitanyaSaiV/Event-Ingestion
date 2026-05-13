package router

import (
	"net/http"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/handlers"
)

func Routes(path string, handler *handlers.IncidentHandler) {
	r := http.NewServeMux()
	r.HandleFunc("GET /healthCheck", handlers.HealthCheck)
	r.HandleFunc("POST /incidents", handler.SaveIncident)
	r.HandleFunc("GET /incidents/{id}", handler.GetIncident)
	r.HandleFunc("GET /incidents", handler.GetAllIncidents)

	http.ListenAndServe(path, r)
}
