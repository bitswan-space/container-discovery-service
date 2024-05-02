package repository

import (
	"context"
	"encoding/json"
	"os"

	"bitswan.space/container-discovery-service/internal/config"
	"bitswan.space/container-discovery-service/internal/models"
	"bitswan.space/container-discovery-service/pkg"
	"github.com/google/uuid"
)

const (
	fetchEntriesErrorMsg = "could not fetch entries: %v"
	writeEntriesErrorMsg = "could not write entries: %v"
)

type CDSRepository interface {
	FetchDashboardEntries(ctx context.Context) (*models.DashboardEntryListResponse, error)
	CreateDashboardEntry(ctx context.Context, entry *models.DashboardEntry) (*models.DashboardEntry, error)
	DeleteDashboardEntry(ctx context.Context, id string) error
	UpdateDashboardEntry(ctx context.Context, id string, entry *models.DashboardEntry) (*models.DashboardEntry, error)
}

type cdsRepository struct {
	appConfig *config.Configuration
}

func NewCDSRepository(appConfig *config.Configuration) CDSRepository {

	return &cdsRepository{
		appConfig: appConfig,
	}
}

func (r *cdsRepository) FetchDashboardEntries(ctx context.Context) (*models.DashboardEntryListResponse, error) {

	entries, err := r.fetchDashboardEntries()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, fetchEntriesErrorMsg, err)
	}

	return &models.DashboardEntryListResponse{
		Results: entries,
	}, nil
}

func (r *cdsRepository) fetchDashboardEntries() ([]*models.DashboardEntry, error) {
	var entries []*models.DashboardEntry

	// read from json file
	jsonData, err := os.ReadFile(r.appConfig.DashboardEntriesFile)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "could not read dashboard entries from file: %v", err)
	}

	err = json.Unmarshal(jsonData, &entries)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "could not unmarshal dashboard entries: %v", err)
	}

	return entries, nil
}

func (r *cdsRepository) CreateDashboardEntry(
	ctx context.Context, entry *models.DashboardEntry) (*models.DashboardEntry, error) {
	// read from json file
	entries, err := r.fetchDashboardEntries()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, fetchEntriesErrorMsg, err)
	}

	entries = append(entries, &models.DashboardEntry{
		Id:          uuid.New().String(),
		Name:        entry.Name,
		Url:         entry.Url,
		Description: entry.Description,
	})

	// write to json file
	err = r.writeDashboardEntries(entries)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, writeEntriesErrorMsg, err)
	}

	return entry, nil
}

func (r *cdsRepository) writeDashboardEntries(entries []*models.DashboardEntry) error {
	// write to json file
	jsonData, err := json.Marshal(entries)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "could not marshal dashboard entries: %v", err)
	}

	err = os.WriteFile(r.appConfig.DashboardEntriesFile, jsonData, 0644)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "could not write dashboard entries to file: %v", err)
	}

	return nil
}

func (r *cdsRepository) DeleteDashboardEntry(ctx context.Context, id string) error {
	// read from json file
	entries, err := r.fetchDashboardEntries()
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "could not fetch entries: %v", err)
	}

	for i, entry := range entries {
		if entry.Id == id {
			entries = append(entries[:i], entries[i+1:]...)
			break
		}
	}

	// write to json file
	err = r.writeDashboardEntries(entries)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, writeEntriesErrorMsg, err)
	}

	return nil
}

func (r *cdsRepository) UpdateDashboardEntry(
	ctx context.Context, id string, entry *models.DashboardEntry) (*models.DashboardEntry, error) {
	// read from json file
	entries, err := r.fetchDashboardEntries()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, fetchEntriesErrorMsg, err)
	}

	for i, e := range entries {
		if e.Id == id {
			entries[i] = &models.DashboardEntry{
				Id:          id,
				Name:        entry.Name,
				Url:         entry.Url,
				Description: entry.Description,
			}
			break
		}
	}

	// write to json file
	err = r.writeDashboardEntries(entries)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, writeEntriesErrorMsg, err)
	}

	return entry, nil
}
