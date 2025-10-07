package main

import (
	"context"
	"testing"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTrainingRepo is a mock implementation of TrainingRepoInterface
type MockTrainingRepo struct {
	mock.Mock
}

func (m *MockTrainingRepo) CreateMaterial(ctx context.Context, material repo.Material) error {
	args := m.Called(ctx, material)
	return args.Error(0)
}

func (m *MockTrainingRepo) GetMaterialByID(ctx context.Context, id string) (*repo.Material, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repo.Material), args.Error(1)
}

func (m *MockTrainingRepo) ListMaterials(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.Material, error) {
	args := m.Called(ctx, tenantID, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repo.Material), args.Error(1)
}

func (m *MockTrainingRepo) UpdateMaterial(ctx context.Context, material repo.Material) error {
	args := m.Called(ctx, material)
	return args.Error(0)
}

func (m *MockTrainingRepo) DeleteMaterial(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTrainingRepo) CreateCourse(ctx context.Context, course repo.TrainingCourse) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockTrainingRepo) GetCourseByID(ctx context.Context, id string) (*repo.TrainingCourse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repo.TrainingCourse), args.Error(1)
}

func (m *MockTrainingRepo) ListCourses(ctx context.Context, tenantID string, filters map[string]interface{}) ([]repo.TrainingCourse, error) {
	args := m.Called(ctx, tenantID, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repo.TrainingCourse), args.Error(1)
}

func (m *MockTrainingRepo) UpdateCourse(ctx context.Context, course repo.TrainingCourse) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockTrainingRepo) DeleteCourse(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Stub methods for other interface methods
func (m *MockTrainingRepo) AddMaterialToCourse(ctx context.Context, courseMaterial repo.CourseMaterial) error {
	args := m.Called(ctx, courseMaterial)
	return args.Error(0)
}

func (m *MockTrainingRepo) RemoveMaterialFromCourse(ctx context.Context, courseID, materialID string) error {
	args := m.Called(ctx, courseID, materialID)
	return args.Error(0)
}

func (m *MockTrainingRepo) GetCourseMaterials(ctx context.Context, courseID string) ([]repo.CourseMaterial, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repo.CourseMaterial), args.Error(1)
}

// All other methods return nil for now
func (m *MockTrainingRepo) CreateAssignment(ctx context.Context, assignment repo.TrainingAssignment) error {
	return nil
}
func (m *MockTrainingRepo) GetAssignmentByID(ctx context.Context, id string) (*repo.TrainingAssignment, error) {
	return nil, nil
}
func (m *MockTrainingRepo) GetUserAssignments(ctx context.Context, userID string, filters map[string]interface{}) ([]repo.TrainingAssignment, error) {
	return nil, nil
}
func (m *MockTrainingRepo) UpdateAssignment(ctx context.Context, assignment repo.TrainingAssignment) error {
	return nil
}
func (m *MockTrainingRepo) DeleteAssignment(ctx context.Context, id string) error { return nil }
func (m *MockTrainingRepo) GetOverdueAssignments(ctx context.Context, tenantID string) ([]repo.TrainingAssignment, error) {
	return nil, nil
}
func (m *MockTrainingRepo) GetUpcomingDeadlines(ctx context.Context, tenantID string, days int) ([]repo.TrainingAssignment, error) {
	return nil, nil
}
func (m *MockTrainingRepo) CreateProgress(ctx context.Context, progress repo.TrainingProgress) error {
	return nil
}
func (m *MockTrainingRepo) UpdateProgress(ctx context.Context, progress repo.TrainingProgress) error {
	return nil
}
func (m *MockTrainingRepo) GetProgressByAssignment(ctx context.Context, assignmentID string) ([]repo.TrainingProgress, error) {
	return nil, nil
}
func (m *MockTrainingRepo) GetProgressByAssignmentAndMaterial(ctx context.Context, assignmentID, materialID string) (*repo.TrainingProgress, error) {
	return nil, nil
}
func (m *MockTrainingRepo) CreateQuizQuestion(ctx context.Context, question repo.QuizQuestion) error {
	return nil
}
func (m *MockTrainingRepo) GetQuizQuestionByID(ctx context.Context, id string) (*repo.QuizQuestion, error) {
	return nil, nil
}
func (m *MockTrainingRepo) ListQuizQuestions(ctx context.Context, materialID string) ([]repo.QuizQuestion, error) {
	return nil, nil
}
func (m *MockTrainingRepo) UpdateQuizQuestion(ctx context.Context, question repo.QuizQuestion) error {
	return nil
}
func (m *MockTrainingRepo) DeleteQuizQuestion(ctx context.Context, id string) error { return nil }
func (m *MockTrainingRepo) CreateQuizAttempt(ctx context.Context, attempt repo.QuizAttempt) error {
	return nil
}
func (m *MockTrainingRepo) GetQuizAttemptByID(ctx context.Context, id string) (*repo.QuizAttempt, error) {
	return nil, nil
}
func (m *MockTrainingRepo) GetQuizAttempts(ctx context.Context, assignmentID, materialID string) ([]repo.QuizAttempt, error) {
	return nil, nil
}
func (m *MockTrainingRepo) CreateCertificate(ctx context.Context, certificate repo.Certificate) error {
	return nil
}
func (m *MockTrainingRepo) GetCertificateByID(ctx context.Context, id string) (*repo.Certificate, error) {
	return nil, nil
}
func (m *MockTrainingRepo) GetCertificateByNumber(ctx context.Context, certificateNumber string) (*repo.Certificate, error) {
	return nil, nil
}
func (m *MockTrainingRepo) GetUserCertificates(ctx context.Context, userID string, filters map[string]interface{}) ([]repo.Certificate, error) {
	return nil, nil
}
func (m *MockTrainingRepo) CreateNotification(ctx context.Context, notification repo.TrainingNotification) error {
	return nil
}
func (m *MockTrainingRepo) GetUserNotifications(ctx context.Context, userID string, unreadOnly bool) ([]repo.TrainingNotification, error) {
	return nil, nil
}
func (m *MockTrainingRepo) MarkNotificationAsRead(ctx context.Context, notificationID, userID string) error {
	return nil
}
func (m *MockTrainingRepo) CreateAnalytics(ctx context.Context, analytics repo.TrainingAnalytics) error {
	return nil
}
func (m *MockTrainingRepo) GetUserAnalytics(ctx context.Context, userID string) (*repo.TrainingAnalytics, error) {
	return nil, nil
}
func (m *MockTrainingRepo) GetCourseAnalytics(ctx context.Context, courseID string) (*repo.CourseAnalytics, error) {
	return nil, nil
}
func (m *MockTrainingRepo) GetOrganizationAnalytics(ctx context.Context, tenantID string) (*repo.OrganizationAnalytics, error) {
	return nil, nil
}
func (m *MockTrainingRepo) CreateRoleAssignment(ctx context.Context, assignment repo.RoleTrainingAssignment) error {
	return nil
}
func (m *MockTrainingRepo) GetRoleAssignments(ctx context.Context, roleID string) ([]repo.RoleTrainingAssignment, error) {
	return nil, nil
}
func (m *MockTrainingRepo) DeleteRoleAssignment(ctx context.Context, id string) error { return nil }

func TestTrainingService_CreateMaterial(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	tenantID := "test-tenant"
	createdBy := "test-user"

	req := dto.CreateMaterialRequest{
		Title:        "Test Material",
		Description:  stringPtr("Test Description"),
		URI:          "test-uri",
		Type:         "file",
		MaterialType: "document",
		IsRequired:   false,
		PassingScore: 80,
		Tags:         []string{"test"},
	}

	expectedMaterial := repo.Material{
		ID:           "temp-id-123",
		TenantID:     tenantID,
		Title:        req.Title,
		Description:  req.Description,
		URI:          req.URI,
		Type:         req.Type,
		MaterialType: req.MaterialType,
		IsRequired:   req.IsRequired,
		PassingScore: req.PassingScore,
		Tags:         req.Tags,
		CreatedBy:    &createdBy,
	}

	mockRepo.On("CreateMaterial", ctx, mock.MatchedBy(func(m repo.Material) bool {
		return m.Title == req.Title && m.TenantID == tenantID && m.CreatedBy != nil && *m.CreatedBy == createdBy
	})).Return(nil)

	material, err := service.CreateMaterial(ctx, tenantID, req, createdBy)

	assert.NoError(t, err)
	assert.NotNil(t, material)
	assert.Equal(t, req.Title, material.Title)
	assert.Equal(t, tenantID, material.TenantID)
	assert.Equal(t, createdBy, *material.CreatedBy)

	mockRepo.AssertExpectations(t)
}

func TestTrainingService_GetMaterial(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	materialID := "test-material-id"

	expectedMaterial := &repo.Material{
		ID:       materialID,
		Title:    "Test Material",
		TenantID: "test-tenant",
	}

	mockRepo.On("GetMaterialByID", ctx, materialID).Return(expectedMaterial, nil)

	material, err := service.GetMaterial(ctx, materialID)

	assert.NoError(t, err)
	assert.NotNil(t, material)
	assert.Equal(t, materialID, material.ID)
	assert.Equal(t, "Test Material", material.Title)

	mockRepo.AssertExpectations(t)
}

func TestTrainingService_ListMaterials(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	tenantID := "test-tenant"
	filters := map[string]interface{}{
		"material_type": "document",
	}

	expectedMaterials := []repo.Material{
		{
			ID:           "material-1",
			Title:        "Material 1",
			TenantID:     tenantID,
			MaterialType: "document",
		},
		{
			ID:           "material-2",
			Title:        "Material 2",
			TenantID:     tenantID,
			MaterialType: "video",
		},
	}

	mockRepo.On("ListMaterials", ctx, tenantID, filters).Return(expectedMaterials, nil)

	materials, err := service.ListMaterials(ctx, tenantID, filters)

	assert.NoError(t, err)
	assert.Len(t, materials, 2)
	assert.Equal(t, "Material 1", materials[0].Title)
	assert.Equal(t, "Material 2", materials[1].Title)

	mockRepo.AssertExpectations(t)
}

func TestTrainingService_UpdateMaterial(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	materialID := "test-material-id"
	updatedBy := "test-user"

	existingMaterial := &repo.Material{
		ID:           materialID,
		Title:        "Original Title",
		Description:  stringPtr("Original Description"),
		URI:          "original-uri",
		Type:         "file",
		MaterialType: "document",
		IsRequired:   false,
		PassingScore: 80,
		Tags:         []string{"original"},
	}

	updateReq := dto.UpdateMaterialRequest{
		Title:       stringPtr("Updated Title"),
		Description: stringPtr("Updated Description"),
		IsRequired:  boolPtr(true),
	}

	expectedUpdatedMaterial := repo.Material{
		ID:           materialID,
		Title:        "Updated Title",
		Description:  stringPtr("Updated Description"),
		URI:          "original-uri",
		Type:         "file",
		MaterialType: "document",
		IsRequired:   true,
		PassingScore: 80,
		Tags:         []string{"original"},
	}

	mockRepo.On("GetMaterialByID", ctx, materialID).Return(existingMaterial, nil)
	mockRepo.On("UpdateMaterial", ctx, mock.MatchedBy(func(m repo.Material) bool {
		return m.ID == materialID && m.Title == "Updated Title" && m.Description != nil && *m.Description == "Updated Description" && m.IsRequired == true
	})).Return(nil)

	err := service.UpdateMaterial(ctx, materialID, updateReq, updatedBy)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTrainingService_DeleteMaterial(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	materialID := "test-material-id"
	deletedBy := "test-user"

	mockRepo.On("DeleteMaterial", ctx, materialID).Return(nil)

	err := service.DeleteMaterial(ctx, materialID, deletedBy)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTrainingService_CreateCourse(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	tenantID := "test-tenant"
	createdBy := "test-user"

	req := dto.CreateCourseRequest{
		Title:       "Test Course",
		Description: stringPtr("Test Course Description"),
		IsActive:    true,
	}

	mockRepo.On("CreateCourse", ctx, mock.MatchedBy(func(c repo.TrainingCourse) bool {
		return c.Title == req.Title && c.TenantID == tenantID && c.CreatedBy != nil && *c.CreatedBy == createdBy
	})).Return(nil)

	course, err := service.CreateCourse(ctx, tenantID, req, createdBy)

	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, req.Title, course.Title)
	assert.Equal(t, tenantID, course.TenantID)
	assert.Equal(t, createdBy, *course.CreatedBy)

	mockRepo.AssertExpectations(t)
}

func TestTrainingService_GetCourse(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	courseID := "test-course-id"

	expectedCourse := &repo.TrainingCourse{
		ID:       courseID,
		Title:    "Test Course",
		TenantID: "test-tenant",
		IsActive: true,
	}

	mockRepo.On("GetCourseByID", ctx, courseID).Return(expectedCourse, nil)

	course, err := service.GetCourse(ctx, courseID)

	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, courseID, course.ID)
	assert.Equal(t, "Test Course", course.Title)

	mockRepo.AssertExpectations(t)
}

func TestTrainingService_ListCourses(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	tenantID := "test-tenant"
	filters := map[string]interface{}{
		"is_active": true,
	}

	expectedCourses := []repo.TrainingCourse{
		{
			ID:       "course-1",
			Title:    "Course 1",
			TenantID: tenantID,
			IsActive: true,
		},
		{
			ID:       "course-2",
			Title:    "Course 2",
			TenantID: tenantID,
			IsActive: false,
		},
	}

	mockRepo.On("ListCourses", ctx, tenantID, filters).Return(expectedCourses, nil)

	courses, err := service.ListCourses(ctx, tenantID, filters)

	assert.NoError(t, err)
	assert.Len(t, courses, 2)
	assert.Equal(t, "Course 1", courses[0].Title)
	assert.Equal(t, "Course 2", courses[1].Title)

	mockRepo.AssertExpectations(t)
}

func TestTrainingService_UpdateCourse(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	courseID := "test-course-id"
	updatedBy := "test-user"

	existingCourse := &repo.TrainingCourse{
		ID:          courseID,
		Title:       "Original Course",
		Description: stringPtr("Original Description"),
		IsActive:    true,
	}

	updateReq := dto.UpdateCourseRequest{
		Title:    stringPtr("Updated Course"),
		IsActive: boolPtr(false),
	}

	mockRepo.On("GetCourseByID", ctx, courseID).Return(existingCourse, nil)
	mockRepo.On("UpdateCourse", ctx, mock.MatchedBy(func(c repo.TrainingCourse) bool {
		return c.ID == courseID && c.Title == "Updated Course" && c.IsActive == false
	})).Return(nil)

	err := service.UpdateCourse(ctx, courseID, updateReq, updatedBy)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTrainingService_DeleteCourse(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	courseID := "test-course-id"
	deletedBy := "test-user"

	mockRepo.On("DeleteCourse", ctx, courseID).Return(nil)

	err := service.DeleteCourse(ctx, courseID, deletedBy)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTrainingService_AddMaterialToCourse(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	courseID := "test-course-id"
	materialID := "test-material-id"
	addedBy := "test-user"

	req := dto.CourseMaterialRequest{
		OrderIndex: 1,
		IsRequired: true,
	}

	expectedCourseMaterial := repo.CourseMaterial{
		ID:         "temp-id-123",
		CourseID:   courseID,
		MaterialID: materialID,
		OrderIndex: 1,
		IsRequired: true,
	}

	mockRepo.On("AddMaterialToCourse", ctx, mock.MatchedBy(func(cm repo.CourseMaterial) bool {
		return cm.CourseID == courseID && cm.MaterialID == materialID && cm.OrderIndex == 1 && cm.IsRequired == true
	})).Return(nil)

	err := service.AddMaterialToCourse(ctx, courseID, materialID, req, addedBy)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTrainingService_RemoveMaterialFromCourse(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	courseID := "test-course-id"
	materialID := "test-material-id"
	removedBy := "test-user"

	mockRepo.On("RemoveMaterialFromCourse", ctx, courseID, materialID).Return(nil)

	err := service.RemoveMaterialFromCourse(ctx, courseID, materialID, removedBy)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTrainingService_GetCourseMaterials(t *testing.T) {
	mockRepo := new(MockTrainingRepo)
	service := domain.NewTrainingService(mockRepo)

	ctx := context.Background()
	courseID := "test-course-id"

	expectedMaterials := []repo.CourseMaterial{
		{
			ID:         "cm-1",
			CourseID:   courseID,
			MaterialID: "material-1",
			OrderIndex: 1,
			IsRequired: true,
		},
		{
			ID:         "cm-2",
			CourseID:   courseID,
			MaterialID: "material-2",
			OrderIndex: 2,
			IsRequired: false,
		},
	}

	mockRepo.On("GetCourseMaterials", ctx, courseID).Return(expectedMaterials, nil)

	materials, err := service.GetCourseMaterials(ctx, courseID)

	assert.NoError(t, err)
	assert.Len(t, materials, 2)
	assert.Equal(t, "material-1", materials[0].MaterialID)
	assert.Equal(t, "material-2", materials[1].MaterialID)

	mockRepo.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

