package models

import (
	"errors"
	"testing"
)

func TestCreateIncidentRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateIncidentRequest
		wantErr error // exact sentinel error we expect (nil for no error)
	}{
		{
			name: "valid request",
			req: CreateIncidentRequest{
				Id:       "1",
				Message:  "db is down",
				Severity: "SEV1",
			},
			wantErr: nil,
		},
		{
			name: "empty id",
			req: CreateIncidentRequest{
				Id:       "",
				Message:  "db is down",
				Severity: "SEV1",
			},
			wantErr: ErrEmptyID,
		},
		{
			name: "whitespace-only id",
			req: CreateIncidentRequest{
				Id:       "   ",
				Message:  "db is down",
				Severity: "SEV1",
			},
			wantErr: ErrEmptyID,
		},
		{
			name: "empty message",
			req: CreateIncidentRequest{
				Id:       "1",
				Message:  "",
				Severity: "SEV1",
			},
			wantErr: ErrEmptyMessage,
		},
		{
			name: "invalid severity",
			req: CreateIncidentRequest{
				Id:       "1",
				Message:  "db is down",
				Severity: "SEVS1",
			},
			wantErr: ErrInvalidSeverity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}
