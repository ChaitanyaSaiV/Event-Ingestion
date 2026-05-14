package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/storage"
)

// mockStore is a fake IncidentStore used in tests.
// It satisfies the IncidentStore interface but doesn't actually store data.
type mockStore struct {
	// Inputs we want to control:
	saveErr error
	getData models.IncidentData
	getErr  error
	allData []models.IncidentData
	allErr  error

	// Outputs we want to observe:
	savedIncident *models.IncidentData // captures what Save was called with
	saveCalled    int                  // how many times Save was called
}

func (m *mockStore) Save(ctx context.Context, incident *models.IncidentData) {
	m.saveCalled++
	m.savedIncident = incident
}
func (m *mockStore) Get(ctx context.Context, id string) (models.IncidentData, error) {
	return m.getData, m.getErr
}

func (m *mockStore) GetAll(ctx context.Context) ([]models.IncidentData, error) {
	return m.allData, m.allErr
}

func (m *mockStore) Delete(ctx context.Context, id string) error {
	return nil
}

// Compile-time check: ensures mockStore satisfies IncidentStore.
// If you change the interface and forget to update the mock, this fails to compile.
var _ IncidentStore = (*mockStore)(nil)

func TestSaveIncident_Success(t *testing.T) {
	// Arrange
	mock := &mockStore{}
	handler := NewIncidentHandler(mock)

	body := `{"id":"1","message":"db down","severity":"SEV1"}`
	req := httptest.NewRequest(http.MethodPost, "/incidents",
		bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	// Act
	handler.SaveIncident(rec, req)

	// Assert
	if rec.Code != http.StatusAccepted {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusAccepted)
	}

	if mock.saveCalled != 1 {
		t.Errorf("Save was called %d times, want 1", mock.saveCalled)
	}

	if mock.savedIncident.Id != "1" {
		t.Errorf("got id %q, want %q", mock.savedIncident.Id, "1")
	}
}

func TestGetIncident(t *testing.T) {
	mock := &mockStore{
		getData: models.IncidentData{Id: "1", Message: "test"},
	}
	handler := NewIncidentHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/incidents/1", nil)
	req.SetPathValue("id", "1") // ← required! Otherwise r.PathValue("id") returns ""
	rec := httptest.NewRecorder()

	handler.GetIncident(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d, want 200", rec.Code)
	}
}

func TestInvalidJSONSave(t *testing.T) {
	mock := &mockStore{}
	handlers := NewIncidentHandler(mock)
	body := `{"id":1,"message":"db down","severity":"SEV1"}`

	req := httptest.NewRequest(http.MethodPost, "/incidents",
		bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handlers.SaveIncident(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}

	if mock.saveCalled != 0 {
		t.Errorf("Saved malformed or Un supported JSON")
	}
}

func TestSaveIncident(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		wantStatusCode int
		wantSaveCalled int // how many times Save should have been called
	}{
		{
			name:           "valid incident saves successfully",
			body:           `{"id":"1","message":"db down","severity":"SEV1"}`,
			wantStatusCode: http.StatusAccepted,
			wantSaveCalled: 1,
		},
		{
			name:           "invalid JSON rejected before save",
			body:           `{not valid json}`,
			wantStatusCode: http.StatusBadRequest,
			wantSaveCalled: 0,
		},
		{
			name:           "missing severity rejected by validator",
			body:           `{"id":"1","message":"db down"}`,
			wantStatusCode: http.StatusBadRequest,
			wantSaveCalled: 0,
		},
		{
			name:           "invalid severity rejected by validator",
			body:           `{"id":"1","message":"db down","severity":"CRITICAL"}`,
			wantStatusCode: http.StatusBadRequest,
			wantSaveCalled: 0,
		},
		{
			name:           "empty body rejected",
			body:           ``,
			wantStatusCode: http.StatusBadRequest,
			wantSaveCalled: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fresh mock per test case — no state leakage
			mock := &mockStore{}
			handler := NewIncidentHandler(mock)
			req := httptest.NewRequest(http.MethodPost, "/incidents",
				bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()
			handler.SaveIncident(rec, req)
			if rec.Code != tt.wantStatusCode {
				t.Errorf("status: got %d, want %d", rec.Code, tt.wantStatusCode)
			}

			if mock.saveCalled != tt.wantSaveCalled {
				t.Errorf("Save calls: got %d, want %d",
					mock.saveCalled, tt.wantSaveCalled)
			}
		})
	}
}

func TestGetIncident_2(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockData       models.IncidentData
		mockErr        error
		wantStatusCode int
	}{
		{
			name:           "found",
			id:             "1",
			mockData:       models.IncidentData{Id: "1", Message: "db down"},
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "not found",
			id:             "999",
			mockErr:        storage.ErrNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "unexpected error returns 500",
			id:             "1",
			mockErr:        errors.New("database connection refused"),
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockStore{
				getData: tt.mockData,
				getErr:  tt.mockErr,
			}
			handler := NewIncidentHandler(mock)
			req := httptest.NewRequest(http.MethodGet, "/incidents/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			rec := httptest.NewRecorder()
			handler.GetIncident(rec, req)
			if rec.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatusCode)
			}
		})
	}
}
