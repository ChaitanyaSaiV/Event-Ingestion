package main

import (
	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/db"
	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/handler"
	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/router"
)

func main() {

	store := db.NewInMemoryStore()

	app := &handler.App{
		Incidents: store,
	}
	router.NewRouter(":8080", app)
}
