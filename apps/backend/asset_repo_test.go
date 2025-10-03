package main

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB - мок для базы данных
type MockDB struct {
	mock.Mock
}

// Убеждаемся, что MockDB реализует интерфейс
var _ repo.DBInterface = (*MockDB)(nil)

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	arguments := m.Called(query, args)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*sql.Rows), arguments.Error(1)
}

func (m *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	arguments := m.Called(ctx, query, args)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*sql.Rows), arguments.Error(1)
}

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	arguments := m.Called(query, args)
	if arguments.Get(0) == nil {
		return nil
	}
	return arguments.Get(0).(*sql.Row)
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	arguments := m.Called(query, args)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(sql.Result), arguments.Error(1)
}

func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	arguments := m.Called(ctx, query, args)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(sql.Result), arguments.Error(1)
}

func TestAssetRepo_Create(t *testing.T) {
	mockDB := new(MockDB)
	assetRepo := repo.NewAssetRepo(mockDB)

	ctx := context.Background()
	asset := repo.Asset{
		ID:              uuid.New().String(),
		TenantID:        "tenant-123",
		InventoryNumber: "AST-20241201-ABC12345",
		Name:            "Test Server",
		Type:            "server",
		Class:           "hardware",
		OwnerID:         stringPtr("owner-123"),
		Location:        stringPtr("Data Center A"),
		Criticality:     "high",
		Confidentiality: "high",
		Integrity:       "high",
		Availability:    "high",
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Настройка мока
	mockDB.On("ExecContext", ctx, mock.AnythingOfType("string"), mock.Anything).Return(nil, nil)

	// Выполнение теста
	err := assetRepo.Create(ctx, asset)

	// Проверки
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestAssetRepo_Update(t *testing.T) {
	mockDB := new(MockDB)
	assetRepo := repo.NewAssetRepo(mockDB)

	ctx := context.Background()
	asset := repo.Asset{
		ID:              uuid.New().String(),
		TenantID:        "tenant-123",
		InventoryNumber: "AST-20241201-ABC12345",
		Name:            "Updated Server",
		Type:            "server",
		Class:           "hardware",
		Criticality:     "high",
		Confidentiality: "high",
		Integrity:       "high",
		Availability:    "high",
		Status:          "active",
		UpdatedAt:       time.Now(),
	}

	// Настройка мока
	mockDB.On("Exec", mock.AnythingOfType("string"), mock.Anything).Return(nil, nil)

	// Выполнение теста
	err := assetRepo.Update(ctx, asset)

	// Проверки
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestAssetRepo_SoftDelete(t *testing.T) {
	mockDB := new(MockDB)
	assetRepo := repo.NewAssetRepo(mockDB)

	ctx := context.Background()
	assetID := uuid.New().String()

	// Настройка мока
	mockDB.On("Exec", mock.AnythingOfType("string"), mock.Anything).Return(nil, nil)

	// Выполнение теста
	err := assetRepo.SoftDelete(ctx, assetID)

	// Проверки
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}
