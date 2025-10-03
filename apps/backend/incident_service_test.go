package main

import (
	"context"
	"testing"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockIncidentRepo - мок для IncidentRepo
type MockIncidentRepo struct {
	mock.Mock
}

// Убеждаемся, что MockIncidentRepo реализует интерфейс
var _ domain.IncidentRepoInterface = (*MockIncidentRepo)(nil)

func (m *MockIncidentRepo) Create(ctx context.Context, incident *repo.Incident) error {
	arguments := m.Called(ctx, incident)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) GetByID(ctx context.Context, id, tenantID string) (*repo.Incident, error) {
	arguments := m.Called(ctx, id, tenantID)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.Incident), arguments.Error(1)
}

func (m *MockIncidentRepo) Update(ctx context.Context, incident *repo.Incident) error {
	arguments := m.Called(ctx, incident)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) Delete(ctx context.Context, id, tenantID string) error {
	arguments := m.Called(ctx, id, tenantID)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) List(ctx context.Context, tenantID string, filters map[string]interface{}, limit, offset int) ([]*repo.Incident, int, error) {
	arguments := m.Called(ctx, tenantID, filters, limit, offset)
	return arguments.Get(0).([]*repo.Incident), arguments.Get(1).(int), arguments.Error(2)
}

func (m *MockIncidentRepo) AddAsset(ctx context.Context, incidentID, assetID string) error {
	arguments := m.Called(ctx, incidentID, assetID)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) RemoveAsset(ctx context.Context, incidentID, assetID string) error {
	arguments := m.Called(ctx, incidentID, assetID)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) GetAssets(ctx context.Context, incidentID string) ([]*repo.Asset, error) {
	arguments := m.Called(ctx, incidentID)
	return arguments.Get(0).([]*repo.Asset), arguments.Error(1)
}

func (m *MockIncidentRepo) AddRisk(ctx context.Context, incidentID, riskID string) error {
	arguments := m.Called(ctx, incidentID, riskID)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) RemoveRisk(ctx context.Context, incidentID, riskID string) error {
	arguments := m.Called(ctx, incidentID, riskID)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) GetRisks(ctx context.Context, incidentID string) ([]*repo.Risk, error) {
	arguments := m.Called(ctx, incidentID)
	return arguments.Get(0).([]*repo.Risk), arguments.Error(1)
}

func (m *MockIncidentRepo) AddComment(ctx context.Context, comment *repo.IncidentComment) error {
	arguments := m.Called(ctx, comment)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) GetComments(ctx context.Context, incidentID string) ([]*repo.IncidentComment, error) {
	arguments := m.Called(ctx, incidentID)
	return arguments.Get(0).([]*repo.IncidentComment), arguments.Error(1)
}

func (m *MockIncidentRepo) AddAttachment(ctx context.Context, attachment *repo.IncidentAttachment) error {
	arguments := m.Called(ctx, attachment)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) GetAttachments(ctx context.Context, incidentID string) ([]*repo.IncidentAttachment, error) {
	arguments := m.Called(ctx, incidentID)
	return arguments.Get(0).([]*repo.IncidentAttachment), arguments.Error(1)
}

func (m *MockIncidentRepo) DeleteAttachment(ctx context.Context, attachmentID string) error {
	arguments := m.Called(ctx, attachmentID)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) AddAction(ctx context.Context, action *repo.IncidentAction) error {
	arguments := m.Called(ctx, action)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) UpdateAction(ctx context.Context, action *repo.IncidentAction) error {
	arguments := m.Called(ctx, action)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) GetActions(ctx context.Context, incidentID string) ([]*repo.IncidentAction, error) {
	arguments := m.Called(ctx, incidentID)
	return arguments.Get(0).([]*repo.IncidentAction), arguments.Error(1)
}

func (m *MockIncidentRepo) DeleteAction(ctx context.Context, actionID string) error {
	arguments := m.Called(ctx, actionID)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) AddMetric(ctx context.Context, metric *repo.IncidentMetrics) error {
	arguments := m.Called(ctx, metric)
	return arguments.Error(0)
}

func (m *MockIncidentRepo) GetMetrics(ctx context.Context, incidentID string) ([]*repo.IncidentMetrics, error) {
	arguments := m.Called(ctx, incidentID)
	return arguments.Get(0).([]*repo.IncidentMetrics), arguments.Error(1)
}

func (m *MockIncidentRepo) GetIncidentMetrics(ctx context.Context, tenantID string) (*repo.IncidentMetricsSummary, error) {
	arguments := m.Called(ctx, tenantID)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.IncidentMetricsSummary), arguments.Error(1)
}

// MockUserRepo - мок для UserRepo
type MockUserRepo struct {
	mock.Mock
}

var _ domain.UserRepoInterface = (*MockUserRepo)(nil)

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*repo.User, error) {
	arguments := m.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.User), arguments.Error(1)
}

// MockAssetRepo - мок для AssetRepo
type MockAssetRepo struct {
	mock.Mock
}

var _ domain.AssetRepoInterface = (*MockAssetRepo)(nil)

func (m *MockAssetRepo) GetByID(ctx context.Context, id string) (*repo.Asset, error) {
	arguments := m.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.Asset), arguments.Error(1)
}

// MockRiskRepo - мок для RiskRepo
type MockRiskRepo struct {
	mock.Mock
}

var _ domain.RiskRepoInterface = (*MockRiskRepo)(nil)

func (m *MockRiskRepo) GetByID(ctx context.Context, id, tenantID string) (*repo.Risk, error) {
	arguments := m.Called(ctx, id, tenantID)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(*repo.Risk), arguments.Error(1)
}

func TestIncidentService_CreateIncident(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New().String()
	userID := uuid.New().String()

	mockIncidentRepo := &MockIncidentRepo{}
	mockUserRepo := &MockUserRepo{}
	mockAssetRepo := &MockAssetRepo{}
	mockRiskRepo := &MockRiskRepo{}

	service := domain.NewIncidentService(mockIncidentRepo, mockUserRepo, mockAssetRepo, mockRiskRepo)

	t.Run("successful creation", func(t *testing.T) {
		req := dto.CreateIncidentRequest{
			IncidentRequest: dto.IncidentRequest{
				Title:       "Test Incident",
				Description: stringPtr("Test Description"),
				Category:    dto.IncidentCategoryTechnicalFailure,
				Criticality: dto.IncidentCriticalityHigh,
				Source:      dto.IncidentSourceUserReport,
				AssetIDs:    []string{uuid.New().String()},
				RiskIDs:     []string{uuid.New().String()},
			},
		}

		// Mock user exists
		mockUserRepo.On("GetByID", ctx, userID).Return(&repo.User{ID: userID}, nil)

		// Mock asset exists
		assetID := req.AssetIDs[0]
		mockAssetRepo.On("GetByID", ctx, assetID).Return(&repo.Asset{ID: assetID}, nil)

		// Mock risk exists
		riskID := req.RiskIDs[0]
		mockRiskRepo.On("GetByID", ctx, riskID, tenantID).Return(&repo.Risk{ID: riskID}, nil)

		// Mock incident creation
		mockIncidentRepo.On("Create", ctx, mock.AnythingOfType("*repo.Incident")).Return(nil)
		mockIncidentRepo.On("AddAsset", ctx, mock.AnythingOfType("string"), assetID).Return(nil)
		mockIncidentRepo.On("AddRisk", ctx, mock.AnythingOfType("string"), riskID).Return(nil)

		incident, err := service.CreateIncident(ctx, tenantID, req, userID)

		assert.NoError(t, err)
		assert.NotNil(t, incident)
		assert.Equal(t, req.Title, incident.Title)
		assert.Equal(t, dto.IncidentStatusNew, incident.Status)

		mockIncidentRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockAssetRepo.AssertExpectations(t)
		mockRiskRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		req := dto.CreateIncidentRequest{
			IncidentRequest: dto.IncidentRequest{
				Title:       "Test Incident",
				Category:    dto.IncidentCategoryTechnicalFailure,
				Criticality: dto.IncidentCriticalityHigh,
				Source:      dto.IncidentSourceUserReport,
				AssignedTo:  stringPtr(uuid.New().String()),
			},
		}

		// Mock user not found
		mockUserRepo.On("GetByID", ctx, *req.AssignedTo).Return(nil, assert.AnError)

		incident, err := service.CreateIncident(ctx, tenantID, req, userID)

		assert.Error(t, err)
		assert.Nil(t, incident)
		assert.Contains(t, err.Error(), "assigned user not found")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("asset not found", func(t *testing.T) {
		req := dto.CreateIncidentRequest{
			IncidentRequest: dto.IncidentRequest{
				Title:       "Test Incident",
				Category:    dto.IncidentCategoryTechnicalFailure,
				Criticality: dto.IncidentCriticalityHigh,
				Source:      dto.IncidentSourceUserReport,
				AssetIDs:    []string{uuid.New().String()},
			},
		}

		// Mock asset not found
		assetID := req.AssetIDs[0]
		mockAssetRepo.On("GetByID", ctx, assetID).Return(nil, assert.AnError)

		incident, err := service.CreateIncident(ctx, tenantID, req, userID)

		assert.Error(t, err)
		assert.Nil(t, incident)
		assert.Contains(t, err.Error(), "asset not found")

		mockAssetRepo.AssertExpectations(t)
	})
}

func TestIncidentService_GetIncident(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New().String()
	incidentID := uuid.New().String()

	mockIncidentRepo := &MockIncidentRepo{}
	mockUserRepo := &MockUserRepo{}
	mockAssetRepo := &MockAssetRepo{}
	mockRiskRepo := &MockRiskRepo{}

	service := domain.NewIncidentService(mockIncidentRepo, mockUserRepo, mockAssetRepo, mockRiskRepo)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedIncident := &repo.Incident{
			ID:          incidentID,
			TenantID:    tenantID,
			Title:       "Test Incident",
			Category:    dto.IncidentCategoryTechnicalFailure,
			Status:      dto.IncidentStatusNew,
			Criticality: dto.IncidentCriticalityHigh,
			Source:      dto.IncidentSourceUserReport,
			ReportedBy:  uuid.New().String(),
		}

		mockIncidentRepo.On("GetByID", ctx, incidentID, tenantID).Return(expectedIncident, nil)

		incident, err := service.GetIncident(ctx, incidentID, tenantID)

		assert.NoError(t, err)
		assert.NotNil(t, incident)
		assert.Equal(t, expectedIncident.ID, incident.ID)
		assert.Equal(t, expectedIncident.Title, incident.Title)

		mockIncidentRepo.AssertExpectations(t)
	})

	t.Run("incident not found", func(t *testing.T) {
		mockIncidentRepo.On("GetByID", ctx, incidentID, tenantID).Return(nil, assert.AnError)

		incident, err := service.GetIncident(ctx, incidentID, tenantID)

		assert.Error(t, err)
		assert.Nil(t, incident)

		mockIncidentRepo.AssertExpectations(t)
	})
}

func TestIncidentService_UpdateIncident(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New().String()
	userID := uuid.New().String()
	incidentID := uuid.New().String()

	mockIncidentRepo := &MockIncidentRepo{}
	mockUserRepo := &MockUserRepo{}
	mockAssetRepo := &MockAssetRepo{}
	mockRiskRepo := &MockRiskRepo{}

	service := domain.NewIncidentService(mockIncidentRepo, mockUserRepo, mockAssetRepo, mockRiskRepo)

	t.Run("successful update", func(t *testing.T) {
		existingIncident := &repo.Incident{
			ID:          incidentID,
			TenantID:    tenantID,
			Title:       "Original Title",
			Category:    dto.IncidentCategoryTechnicalFailure,
			Status:      dto.IncidentStatusNew,
			Criticality: dto.IncidentCriticalityHigh,
			Source:      dto.IncidentSourceUserReport,
			ReportedBy:  userID,
		}

		req := dto.UpdateIncidentRequest{
			Title:       stringPtr("Updated Title"),
			Description: stringPtr("Updated Description"),
			Status:      stringPtr(dto.IncidentStatusInProgress),
		}

		mockIncidentRepo.On("GetByID", ctx, incidentID, tenantID).Return(existingIncident, nil)
		mockIncidentRepo.On("Update", ctx, mock.AnythingOfType("*repo.Incident")).Return(nil)

		incident, err := service.UpdateIncident(ctx, incidentID, tenantID, req, userID)

		assert.NoError(t, err)
		assert.NotNil(t, incident)
		assert.Equal(t, "Updated Title", incident.Title)
		assert.Equal(t, dto.IncidentStatusInProgress, incident.Status)

		mockIncidentRepo.AssertExpectations(t)
	})
}

func TestIncidentService_DeleteIncident(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New().String()
	incidentID := uuid.New().String()

	mockIncidentRepo := &MockIncidentRepo{}
	mockUserRepo := &MockUserRepo{}
	mockAssetRepo := &MockAssetRepo{}
	mockRiskRepo := &MockRiskRepo{}

	service := domain.NewIncidentService(mockIncidentRepo, mockUserRepo, mockAssetRepo, mockRiskRepo)

	t.Run("successful deletion", func(t *testing.T) {
		mockIncidentRepo.On("Delete", ctx, incidentID, tenantID).Return(nil)

		err := service.DeleteIncident(ctx, incidentID, tenantID)

		assert.NoError(t, err)

		mockIncidentRepo.AssertExpectations(t)
	})
}

func TestIncidentService_ListIncidents(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New().String()

	mockIncidentRepo := &MockIncidentRepo{}
	mockUserRepo := &MockUserRepo{}
	mockAssetRepo := &MockAssetRepo{}
	mockRiskRepo := &MockRiskRepo{}

	service := domain.NewIncidentService(mockIncidentRepo, mockUserRepo, mockAssetRepo, mockRiskRepo)

	t.Run("successful listing", func(t *testing.T) {
		req := dto.IncidentListRequest{
			Page:     1,
			PageSize: 20,
			Status:   dto.IncidentStatusNew,
		}

		expectedIncidents := []*repo.Incident{
			{
				ID:          uuid.New().String(),
				TenantID:    tenantID,
				Title:       "Test Incident 1",
				Category:    dto.IncidentCategoryTechnicalFailure,
				Status:      dto.IncidentStatusNew,
				Criticality: dto.IncidentCriticalityHigh,
				Source:      dto.IncidentSourceUserReport,
			},
		}

		mockIncidentRepo.On("List", ctx, tenantID, mock.AnythingOfType("map[string]interface {}"), 20, 0).Return(expectedIncidents, 1, nil)

		incidents, total, err := service.ListIncidents(ctx, tenantID, req)

		assert.NoError(t, err)
		assert.Len(t, incidents, 1)
		assert.Equal(t, 1, total)

		mockIncidentRepo.AssertExpectations(t)
	})
}

func TestIncidentService_AddComment(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New().String()
	userID := uuid.New().String()
	incidentID := uuid.New().String()

	mockIncidentRepo := &MockIncidentRepo{}
	mockUserRepo := &MockUserRepo{}
	mockAssetRepo := &MockAssetRepo{}
	mockRiskRepo := &MockRiskRepo{}

	service := domain.NewIncidentService(mockIncidentRepo, mockUserRepo, mockAssetRepo, mockRiskRepo)

	t.Run("successful comment addition", func(t *testing.T) {
		req := dto.IncidentCommentRequest{
			Comment:    "Test comment",
			IsInternal: false,
		}

		existingIncident := &repo.Incident{
			ID:       incidentID,
			TenantID: tenantID,
		}

		mockIncidentRepo.On("GetByID", ctx, incidentID, tenantID).Return(existingIncident, nil)
		mockIncidentRepo.On("AddComment", ctx, mock.AnythingOfType("*repo.IncidentComment")).Return(nil)

		comment, err := service.AddComment(ctx, incidentID, tenantID, req, userID)

		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, req.Comment, comment.Comment)
		assert.Equal(t, userID, comment.UserID)

		mockIncidentRepo.AssertExpectations(t)
	})
}

func TestIncidentService_GetIncidentMetrics(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New().String()

	mockIncidentRepo := &MockIncidentRepo{}
	mockUserRepo := &MockUserRepo{}
	mockAssetRepo := &MockAssetRepo{}
	mockRiskRepo := &MockRiskRepo{}

	service := domain.NewIncidentService(mockIncidentRepo, mockUserRepo, mockAssetRepo, mockRiskRepo)

	t.Run("successful metrics retrieval", func(t *testing.T) {
		expectedMetrics := &repo.IncidentMetricsSummary{
			TotalIncidents:  10,
			OpenIncidents:   5,
			ClosedIncidents: 5,
			AverageMTTR:     2.5,
			AverageMTTD:     1.0,
			ByCriticality: map[string]int{
				"high": 3,
				"low":  7,
			},
			ByCategory: map[string]int{
				"technical_failure": 8,
				"data_breach":       2,
			},
			ByStatus: map[string]int{
				"new":    3,
				"closed": 5,
			},
		}

		mockIncidentRepo.On("GetIncidentMetrics", ctx, tenantID).Return(expectedMetrics, nil)

		metrics, err := service.GetIncidentMetrics(ctx, tenantID)

		assert.NoError(t, err)
		assert.NotNil(t, metrics)
		assert.Equal(t, expectedMetrics.TotalIncidents, metrics.TotalIncidents)
		assert.Equal(t, expectedMetrics.OpenIncidents, metrics.OpenIncidents)

		mockIncidentRepo.AssertExpectations(t)
	})
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
