package models

import "time"

type Response struct {
	Key string `json:"key"`
}

type HealthCheck struct {
	Health    string    `json:"health"`
	TimeStamp time.Time `json:"time"`
}

type IncidentData struct {
	Id        string    `json:"id"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"time"`
}
