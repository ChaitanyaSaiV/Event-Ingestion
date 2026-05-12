package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(models.HealthCheck{
		Health:    "OK",
		TimeStamp: time.Now(),
	})

	w.Write(data)
}
