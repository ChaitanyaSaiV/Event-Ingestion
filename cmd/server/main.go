package main

import (
	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/handlers"
	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/router"
	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/storage"
)

func main() {

	memoryStore := storage.NewInMemoryStore()

	handler := handlers.NewIncidentHandler(memoryStore)
	router.Routes(":8080", handler)
}
