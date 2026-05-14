package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/models"
)

type FileStorage struct {
	mu       sync.RWMutex
	fileName string
}

func NewFileStorage(filename string) (*FileStorage, error) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(filename, []byte("{}"), 0644); err != nil {
			return nil, fmt.Errorf("initialize file store: %w", err)
		}
	}
	return &FileStorage{fileName: filename}, nil
}

func (fs *FileStorage) Get(ctx context.Context, id string) (models.IncidentData, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	data, err := os.ReadFile(fs.fileName)
	if err != nil {
		return models.IncidentData{}, fmt.Errorf("Having issue reading the file : %w", err.Error())
	}
	var incidents map[string]models.IncidentData

	json.Unmarshal(data, &incidents)

	incident, ok := incidents[id]

	if !ok {
		return models.IncidentData{}, ErrNotFound
	}

	return incident, nil
}

func (fs *FileStorage) Save(ctx context.Context, incident *models.IncidentData) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	data, _ := os.ReadFile(fs.fileName)

	var incidents map[string]models.IncidentData
	json.Unmarshal(data, &incidents)
	incidents[incident.Id] = *incident

	newData, _ := json.MarshalIndent(incidents, "", "	")

	os.WriteFile(fs.fileName, newData, 0644)
}

func (fs *FileStorage) Delete(ctx context.Context, id string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := os.ReadFile(fs.fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return fmt.Errorf("Error reading the file")
	}
	var incidents map[string]models.IncidentData
	err = json.Unmarshal(data, &incidents)
	if err != nil {
		return err
	}

	delete(incidents, id)

	newData, err := json.MarshalIndent(incidents, "", "	")

	if err != nil {
		return fmt.Errorf("Error while converting the data to JSON")
	}

	err = os.WriteFile(fs.fileName, newData, 0644)

	if err != nil {
		return fmt.Errorf("Error writing the data to a file")
	}

	return nil
}

func (fs *FileStorage) GetAll(ctx context.Context) ([]models.IncidentData, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	data, err := os.ReadFile(fs.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.IncidentData{}, fmt.Errorf("File or data source does not exist")
		}
		return []models.IncidentData{}, fmt.Errorf("Error reading the file")
	}
	var incidents map[string]models.IncidentData
	err = json.Unmarshal(data, &incidents)
	if err != nil {
		return []models.IncidentData{}, fmt.Errorf("Error reading the file")
	}
	returnData := make([]models.IncidentData, 0, len(incidents))
	for _, val := range incidents {
		returnData = append(returnData, val)
	}
	return returnData, nil
}
