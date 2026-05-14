package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
)

func TestHelloWorld(t *testing.T) {
	got := "hello"
	want := "hello"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestNumberOfWays(t *testing.T) {
	got := "hello"
	want := "hello"
	if got != want {
		t.Fatalf("Fatal Error")
	}

	if got != want {
		t.Errorf("Fatal Error")
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"Negative Number", 1, 2, 3},
	}
	for _, val := range tests {
		got := val.a + val.b
		if got != val.want {
			t.Run(val.name, func(t *testing.T) {
				t.Errorf("Got : %d is not matching with want : %d", got, val.want)
			})
		}
	}
	t.Log("Successfully completed the Add testing")
}

func TestInMemoryStore_SaveAndGet(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	incident := &models.IncidentData{
		Id:       "1",
		Message:  "db down",
		Severity: "SEV1",
	}
	// Save it
	store.Save(ctx, incident)
	// Get it back
	got, err := store.Get(ctx, "1")
	if err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}
	if got.Id != "1" {
		t.Errorf("got id %q, want %q", got.Id, "1")
	}
	if got.Message != "db down" {
		t.Errorf("got message %q, want %q", got.Message, "db down")
	}
}
func TestInMemoryStore_GetNotFound(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	_, err := store.Get(ctx, "nonexistent")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("got error %v, want ErrNotFound", err)
	}
}
func TestInMemoryStore_OverwriteExisting(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	store.Save(ctx, &models.IncidentData{Id: "1", Message: "first"})
	store.Save(ctx, &models.IncidentData{Id: "1", Message: "second"})
	got, _ := store.Get(ctx, "1")
	if got.Message != "second" {
		t.Errorf("expected overwrite, got %q", got.Message)
	}
}

func TestInMemoryStore_GetAll_Empty(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	data, _ := store.GetAll(ctx)
	if len(data) != 0 {
		t.Errorf("Expected 0 records but got few")
	}
}

func TestInMemoryStore_GetAll_Non_Empty(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	incident := &models.IncidentData{
		Id:       "1",
		Message:  "db down",
		Severity: "SEV1",
	}
	store.Save(ctx, incident)
	data, _ := store.GetAll(ctx)
	if len(data) != 1 {
		t.Errorf("Expected 1 records but got 0")
	}
}
