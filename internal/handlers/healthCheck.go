package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	responseData := &models.HealthCheck{
		Status:    "Okay",
		TimeStamp: time.Now(),
	}
	json.NewEncoder(w).Encode(responseData)
}
