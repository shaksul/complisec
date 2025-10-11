//go:build legacy_tests
// +build legacy_tests

package main

import (
	"context"
	"testing"
	"time"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDocumentRepo - мок для DocumentRepo
type MockDocumentRepo struct {
	mock.Mock
}

func (m *MockDocumentRepo) CreateFolder(ctx context.Context, folder repo.Folder) error {
	args := m.Called(ctx, folder)
	return args.Error(0)
}

func (m *MockDocumentRepo) GetFolderByID(ctx context.Context, id, tenantID string) (*repo.Folder, error) {
	args := m.Called(ctx, id, tenantID)
	return args.Get(0).(*repo.Folder), args.Error(1)
}

func (m *MockDocumentRepo) ListFolders(ctx context.Context, tenantID string, parentID *string) ([]repo.Folder, error) {
	args := m.Called(ctx, tenantID, parentID)
	return args.Get(0).([]repo.Folder), args.Error(1)
}

func (m *MockDocumentRepo) UpdateFolder(ctx context.Context, folder repo.Folder) error {
	args := m.Called(ctx, folder)
	return args.Error(0)
}

func (m *MockDocumentRepo) DeleteFolder(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *MockDocumentRepo) CreateDocument(ctx context.Context, document repo.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *MockDocumentRepo) GetDocumentByID(ctx context.Context, id, tenantID string) (*repo.Document, error) {
	args := m.Called(ctx, id, tenantID)
	return args.Get(0).(*repo.Document), args.Error(1)
}

func (m *MockDocumentRepo) ListDocuments(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.Document, error) {
	args := m.Called(ctx, tenantID, filters)
	return args.Get(0).([]repo.Document), args.Error(1)
}

func (m *MockDocumentRepo) UpdateDocument(ctx context.Context, document repo.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *MockDocumentRepo) DeleteDocument(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *MockDocumentRepo) AddDocumentTag(ctx context.Context, documentID, tag string) error {
	args := m.Called(ctx, documentID, tag)
	return args.Error(0)
}

func (m *MockDocumentRepo) RemoveDocumentTag(ctx context.Context, documentID, tag string) error {
	args := m.Called(ctx, documentID, tag)
	return args.Error(0)
}

func (m *MockDocumentRepo) GetDocumentTags(ctx context.Context, documentID string) ([]string, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockDocumentRepo) AddDocumentLink(ctx context.Context, link repo.DocumentLink) error {
	args := m.Called(ctx, link)
	return args.Error(0)
}

func (m *MockDocumentRepo) GetDocumentLinks(ctx context.Context, documentID string) ([]repo.DocumentLink, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).([]repo.DocumentLink), args.Error(1)
}

func (m *MockDocumentRepo) CreateOCRText(ctx context.Context, ocrText repo.OCRText) error {
	args := m.Called(ctx, ocrText)
	return args.Error(0)
}

func (m *MockDocumentRepo) GetOCRText(ctx context.Context, documentID string) (*repo.OCRText, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(*repo.OCRText), args.Error(1)
}

func (m *MockDocumentRepo) CreateDocumentPermission(ctx context.Context, permission repo.DocumentPermission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockDocumentRepo) GetDocumentPermissions(ctx context.Context, objectType, objectID, tenantID string) ([]repo.DocumentPermission, error) {
	args := m.Called(ctx, objectType, objectID, tenantID)
	return args.Get(0).([]repo.DocumentPermission), args.Error(1)
}

func (m *MockDocumentRepo) CreateDocumentVersion(ctx context.Context, version repo.DocumentVersion) error {
	args := m.Called(ctx, version)
	return args.Error(0)
}

func (m *MockDocumentRepo) GetDocumentVersions(ctx context.Context, documentID string) ([]repo.DocumentVersion, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).([]repo.DocumentVersion), args.Error(1)
}

func (m *MockDocumentRepo) CreateDocumentAuditLog(ctx context.Context, log repo.DocumentAuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockDocumentRepo) GetDocumentAuditLog(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.DocumentAuditLog, error) {
	args := m.Called(ctx, tenantID, filters)
	return args.Get(0).([]repo.DocumentAuditLog), args.Error(1)
}

func (m *MockDocumentRepo) SearchDocuments(ctx context.Context, tenantID, searchTerm string) ([]repo.Document, error) {
	args := m.Called(ctx, tenantID, searchTerm)
	return args.Get(0).([]repo.Document), args.Error(1)
}

func TestDocumentService_CreateFolder(t *testing.T) {
	mockRepo := new(MockDocumentRepo)
	service := domain.NewDocumentService(mockRepo, "./test-storage")

	ctx := context.Background()
	tenantID := "test-tenant"
	userID := "test-user"
	req := dto.CreateFolderDTO{
		Name:        "Test Folder",
		Description: stringPtr("Test Description"),
		ParentID:    nil,
		Metadata:    stringPtr("{}"),
	}

	// Настраиваем мок
	mockRepo.On("CreateFolder", ctx, mock.AnythingOfType("repo.Folder")).Return(nil)
	mockRepo.On("CreateDocumentAuditLog", ctx, mock.AnythingOfType("repo.DocumentAuditLog")).Return(nil)

	// Выполняем тест
	result, err := service.CreateFolder(ctx, tenantID, req, userID)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Folder", result.Name)
	assert.Equal(t, "Test Description", *result.Description)
	assert.Equal(t, userID, result.OwnerID)
	assert.Equal(t, userID, result.CreatedBy)

	// Проверяем, что все методы были вызваны
	mockRepo.AssertExpectations(t)
}

func TestDocumentService_GetFolder(t *testing.T) {
	mockRepo := new(MockDocumentRepo)
	service := domain.NewDocumentService(mockRepo, "./test-storage")

	ctx := context.Background()
	tenantID := "test-tenant"
	folderID := "test-folder-id"

	expectedFolder := &repo.Folder{
		ID:          folderID,
		TenantID:    tenantID,
		Name:        "Test Folder",
		Description: stringPtr("Test Description"),
		ParentID:    nil,
		OwnerID:     "test-user",
		CreatedBy:   "test-user",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
		Metadata:    stringPtr("{}"),
	}

	// Настраиваем мок
	mockRepo.On("GetFolderByID", ctx, folderID, tenantID).Return(expectedFolder, nil)

	// Выполняем тест
	result, err := service.GetFolder(ctx, folderID, tenantID)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, folderID, result.ID)
	assert.Equal(t, "Test Folder", result.Name)
	assert.Equal(t, "Test Description", *result.Description)

	// Проверяем, что все методы были вызваны
	mockRepo.AssertExpectations(t)
}

func TestDocumentService_ListFolders(t *testing.T) {
	mockRepo := new(MockDocumentRepo)
	service := domain.NewDocumentService(mockRepo, "./test-storage")

	ctx := context.Background()
	tenantID := "test-tenant"
	parentID := stringPtr("parent-folder-id")

	expectedFolders := []repo.Folder{
		{
			ID:          "folder-1",
			TenantID:    tenantID,
			Name:        "Folder 1",
			Description: stringPtr("Description 1"),
			ParentID:    parentID,
			OwnerID:     "test-user",
			CreatedBy:   "test-user",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
			Metadata:    stringPtr("{}"),
		},
		{
			ID:          "folder-2",
			TenantID:    tenantID,
			Name:        "Folder 2",
			Description: stringPtr("Description 2"),
			ParentID:    parentID,
			OwnerID:     "test-user",
			CreatedBy:   "test-user",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
			Metadata:    stringPtr("{}"),
		},
	}

	// Настраиваем мок
	mockRepo.On("ListFolders", ctx, tenantID, parentID).Return(expectedFolders, nil)

	// Выполняем тест
	result, err := service.ListFolders(ctx, tenantID, parentID)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Folder 1", result[0].Name)
	assert.Equal(t, "Folder 2", result[1].Name)

	// Проверяем, что все методы были вызваны
	mockRepo.AssertExpectations(t)
}

func TestDocumentService_UpdateFolder(t *testing.T) {
	mockRepo := new(MockDocumentRepo)
	service := domain.NewDocumentService(mockRepo, "./test-storage")

	ctx := context.Background()
	tenantID := "test-tenant"
	folderID := "test-folder-id"
	userID := "test-user"

	req := dto.UpdateFolderDTO{
		Name:        "Updated Folder",
		Description: stringPtr("Updated Description"),
		Metadata:    stringPtr("{}"),
	}

	existingFolder := &repo.Folder{
		ID:          folderID,
		TenantID:    tenantID,
		Name:        "Original Folder",
		Description: stringPtr("Original Description"),
		ParentID:    nil,
		OwnerID:     userID,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
		Metadata:    stringPtr("{}"),
	}

	// Настраиваем мок
	mockRepo.On("GetFolderByID", ctx, folderID, tenantID).Return(existingFolder, nil)
	mockRepo.On("UpdateFolder", ctx, mock.AnythingOfType("repo.Folder")).Return(nil)
	mockRepo.On("CreateDocumentAuditLog", ctx, mock.AnythingOfType("repo.DocumentAuditLog")).Return(nil)

	// Выполняем тест
	err := service.UpdateFolder(ctx, folderID, tenantID, req, userID)

	// Проверяем результат
	assert.NoError(t, err)

	// Проверяем, что все методы были вызваны
	mockRepo.AssertExpectations(t)
}

func TestDocumentService_DeleteFolder(t *testing.T) {
	mockRepo := new(MockDocumentRepo)
	service := domain.NewDocumentService(mockRepo, "./test-storage")

	ctx := context.Background()
	tenantID := "test-tenant"
	folderID := "test-folder-id"
	userID := "test-user"

	// Настраиваем мок
	mockRepo.On("DeleteFolder", ctx, folderID, tenantID).Return(nil)
	mockRepo.On("CreateDocumentAuditLog", ctx, mock.AnythingOfType("repo.DocumentAuditLog")).Return(nil)

	// Выполняем тест
	err := service.DeleteFolder(ctx, folderID, tenantID, userID)

	// Проверяем результат
	assert.NoError(t, err)

	// Проверяем, что все методы были вызваны
	mockRepo.AssertExpectations(t)
}

func TestDocumentService_GetDocumentStats(t *testing.T) {
	mockRepo := new(MockDocumentRepo)
	service := domain.NewDocumentService(mockRepo, "./test-storage")

	ctx := context.Background()
	tenantID := "test-tenant"

	documents := []repo.Document{
		{
			ID:           "doc-1",
			TenantID:     tenantID,
			Name:         "Document 1",
			OriginalName: "doc1.pdf",
			FilePath:     "/path/to/doc1.pdf",
			FileSize:     1024,
			MimeType:     "application/pdf",
			FileHash:     "hash1",
			FolderID:     nil,
			OwnerID:      "test-user",
			CreatedBy:    "test-user",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			IsActive:     true,
			Version:      1,
			Metadata:     stringPtr("{}"),
		},
		{
			ID:           "doc-2",
			TenantID:     tenantID,
			Name:         "Document 2",
			OriginalName: "doc2.jpg",
			FilePath:     "/path/to/doc2.jpg",
			FileSize:     2048,
			MimeType:     "image/jpeg",
			FileHash:     "hash2",
			FolderID:     nil,
			OwnerID:      "test-user",
			CreatedBy:    "test-user",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			IsActive:     true,
			Version:      1,
			Metadata:     stringPtr("{}"),
		},
	}

	folders := []repo.Folder{
		{
			ID:          "folder-1",
			TenantID:    tenantID,
			Name:        "Folder 1",
			Description: stringPtr("Description 1"),
			ParentID:    nil,
			OwnerID:     "test-user",
			CreatedBy:   "test-user",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
			Metadata:    stringPtr("{}"),
		},
	}

	// Настраиваем мок
	mockRepo.On("ListDocuments", ctx, tenantID, map[string]interface{}{}).Return(documents, nil)
	mockRepo.On("ListFolders", ctx, tenantID, (*string)(nil)).Return(folders, nil)

	// Выполняем тест
	result, err := service.GetDocumentStats(ctx, tenantID)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalDocuments)
	assert.Equal(t, 1, result.TotalFolders)
	assert.Equal(t, int64(3072), result.TotalSize) // 1024 + 2048
	assert.Equal(t, int64(3072), result.StorageUsage)
	assert.Equal(t, 1, result.DocumentsByType["application"])
	assert.Equal(t, 1, result.DocumentsByType["image"])

	// Проверяем, что все методы были вызваны
	mockRepo.AssertExpectations(t)
}

func TestDocumentService_SearchDocuments(t *testing.T) {
	mockRepo := new(MockDocumentRepo)
	service := domain.NewDocumentService(mockRepo, "./test-storage")

	ctx := context.Background()
	tenantID := "test-tenant"
	searchTerm := "test search"

	documents := []repo.Document{
		{
			ID:           "doc-1",
			TenantID:     tenantID,
			Name:         "Test Document",
			OriginalName: "test.pdf",
			FilePath:     "/path/to/test.pdf",
			FileSize:     1024,
			MimeType:     "application/pdf",
			FileHash:     "hash1",
			FolderID:     nil,
			OwnerID:      "test-user",
			CreatedBy:    "test-user",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			IsActive:     true,
			Version:      1,
			Metadata:     stringPtr("{}"),
		},
	}

	ocrText := &repo.OCRText{
		ID:         "ocr-1",
		DocumentID: "doc-1",
		Content:    "This is test content",
		Language:   "ru",
		Confidence: floatPtr(0.95),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Настраиваем мок
	mockRepo.On("SearchDocuments", ctx, tenantID, searchTerm).Return(documents, nil)
	mockRepo.On("GetOCRText", ctx, "doc-1").Return(ocrText, nil)

	// Выполняем тест
	result, err := service.SearchDocuments(ctx, tenantID, searchTerm)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "doc-1", result[0].DocumentID)
	assert.Equal(t, "Test Document", result[0].Name)
	assert.Equal(t, "application/pdf", result[0].MimeType)
	assert.Equal(t, int64(1024), result[0].FileSize)
	assert.Equal(t, "This is test content", *result[0].OCRText)

	// Проверяем, что все методы были вызваны
	mockRepo.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
