package router

import (
	"net/http"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/handlers"
)

func Routes(path string, handler *handlers.IncidentHandler) {
	r := http.NewServeMux()
	r.HandleFunc("GET /healthCheck", handlers.HealthCheck)
	r.HandleFunc("POST /incident", handler.SaveIncident)
	r.HandleFunc("GET /incident/{id}", handler.GetIncident)

	http.ListenAndServe(path, r)
}
