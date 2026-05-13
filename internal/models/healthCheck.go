package models

import "time"

type HealthCheck struct {
	Status    string    `json:"status"`
	TimeStamp time.Time `json:"timeStamp"`
}
