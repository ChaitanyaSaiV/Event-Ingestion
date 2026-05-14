package models

import (
	"errors"
	"strings"
	"time"
)

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

var (
	ErrEmptyID         = errors.New("id is required")
	ErrEmptyMessage    = errors.New("message is required")
	ErrInvalidSeverity = errors.New("severity must be SEV1, SEV2, or SEV3")
)

func (r *CreateIncidentRequest) Validate() error {
	if strings.TrimSpace(r.Id) == "" {
		return ErrEmptyID
	}
	if strings.TrimSpace(r.Message) == "" {
		return ErrEmptyMessage
	}
	if r.Severity != "SEV1" && r.Severity != "SEV2" && r.Severity != "SEV3" {
		return ErrInvalidSeverity
	}
	return nil
}
