package router

import (
	"net/http"
	"time"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/handler"
)

func NewRouter(path string, app *handler.App) {

	router := http.NewServeMux()
	router.HandleFunc("GET /healthCheck", handler.HealthCheck)
	router.HandleFunc("GET /incident/{id}", app.GetIncident)
	router.HandleFunc("POST /incident", app.PostIncident)

	s := &http.Server{
		Addr:           path,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
