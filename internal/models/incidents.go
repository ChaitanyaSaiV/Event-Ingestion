package models

import "time"

type CreateIncidentRequest struct {
	Id       string `json:"id"        validate:"required"`
	Message  string `json:"message"   validate:"required"`
	Severity string `json:"severity"  validate:"required,oneof=SEV1 SEV2 SEV3"`
}

// What we store and return
type IncidentData struct {
	Id        string    `json:"id"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	TimeStamp time.Time `json:"timeStamp"`
}
