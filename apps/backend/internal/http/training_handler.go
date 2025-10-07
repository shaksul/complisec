package http

import (
	"strconv"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type TrainingHandler struct {
	trainingService domain.TrainingServiceInterface
}

func NewTrainingHandler(trainingService domain.TrainingServiceInterface) *TrainingHandler {
	return &TrainingHandler{
		trainingService: trainingService,
	}
}

// Materials handlers

// CreateMaterial godoc
// @Summary Create training material
// @Description Create a new training material
// @Tags training
// @Accept json
// @Produce json
// @Param request body dto.CreateMaterialRequest true "Material data"
// @Success 201 {object} dto.MaterialResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/materials [post]
func (h *TrainingHandler) CreateMaterial(c *fiber.Ctx) error {
	var req dto.CreateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	material, err := h.trainingService.CreateMaterial(c.Context(), tenantID, req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create material",
		})
	}

	return c.Status(201).JSON(materialToResponse(*material))
}

// GetMaterial godoc
// @Summary Get training material
// @Description Get training material by ID
// @Tags training
// @Produce json
// @Param id path string true "Material ID"
// @Success 200 {object} dto.MaterialResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/materials/{id} [get]
func (h *TrainingHandler) GetMaterial(c *fiber.Ctx) error {
	materialID := c.Params("id")
	if materialID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Material ID is required",
		})
	}

	material, err := h.trainingService.GetMaterial(c.Context(), materialID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Material not found",
		})
	}

	return c.JSON(materialToResponse(*material))
}

// ListMaterials godoc
// @Summary List training materials
// @Description List training materials with filters
// @Tags training
// @Produce json
// @Param material_type query string false "Material type filter"
// @Param is_required query bool false "Required filter"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.TrainingListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/materials [get]
func (h *TrainingHandler) ListMaterials(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	filters := make(map[string]interface{})

	if materialType := c.Query("material_type"); materialType != "" {
		filters["material_type"] = materialType
	}

	if isRequired := c.Query("is_required"); isRequired != "" {
		if parsed, err := strconv.ParseBool(isRequired); err == nil {
			filters["is_required"] = parsed
		}
	}

	materials, err := h.trainingService.ListMaterials(c.Context(), tenantID, filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to list materials",
		})
	}

	// Convert to response format
	var responses []dto.MaterialResponse
	for _, material := range materials {
		responses = append(responses, materialToResponse(material))
	}

	return c.JSON(dto.TrainingListResponse{
		Items:      interfaceSlice(responses),
		Total:      int64(len(responses)),
		Page:       1,
		PageSize:   len(responses),
		TotalPages: 1,
	})
}

// UpdateMaterial godoc
// @Summary Update training material
// @Description Update training material by ID
// @Tags training
// @Accept json
// @Produce json
// @Param id path string true "Material ID"
// @Param request body dto.UpdateMaterialRequest true "Material data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/materials/{id} [put]
func (h *TrainingHandler) UpdateMaterial(c *fiber.Ctx) error {
	materialID := c.Params("id")
	if materialID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Material ID is required",
		})
	}

	var req dto.UpdateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	userID := c.Locals("user_id").(string)

	err := h.trainingService.UpdateMaterial(c.Context(), materialID, req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update material",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Material updated successfully",
	})
}

// DeleteMaterial godoc
// @Summary Delete training material
// @Description Delete training material by ID
// @Tags training
// @Param id path string true "Material ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/materials/{id} [delete]
func (h *TrainingHandler) DeleteMaterial(c *fiber.Ctx) error {
	materialID := c.Params("id")
	if materialID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Material ID is required",
		})
	}

	userID := c.Locals("user_id").(string)

	err := h.trainingService.DeleteMaterial(c.Context(), materialID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete material",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Material deleted successfully",
	})
}

// Courses handlers

// CreateCourse godoc
// @Summary Create training course
// @Description Create a new training course
// @Tags training
// @Accept json
// @Produce json
// @Param request body dto.CreateCourseRequest true "Course data"
// @Success 201 {object} dto.CourseResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/courses [post]
func (h *TrainingHandler) CreateCourse(c *fiber.Ctx) error {
	var req dto.CreateCourseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	course, err := h.trainingService.CreateCourse(c.Context(), tenantID, req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create course",
		})
	}

	return c.Status(201).JSON(courseToResponse(*course))
}

// GetCourse godoc
// @Summary Get training course
// @Description Get training course by ID
// @Tags training
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} dto.CourseResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/courses/{id} [get]
func (h *TrainingHandler) GetCourse(c *fiber.Ctx) error {
	courseID := c.Params("id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Course ID is required",
		})
	}

	course, err := h.trainingService.GetCourse(c.Context(), courseID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Course not found",
		})
	}

	return c.JSON(courseToResponse(*course))
}

// ListCourses godoc
// @Summary List training courses
// @Description List training courses with filters
// @Tags training
// @Produce json
// @Param is_active query bool false "Active filter"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.TrainingListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/courses [get]
func (h *TrainingHandler) ListCourses(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	filters := make(map[string]interface{})

	if isActive := c.Query("is_active"); isActive != "" {
		if parsed, err := strconv.ParseBool(isActive); err == nil {
			filters["is_active"] = parsed
		}
	}

	courses, err := h.trainingService.ListCourses(c.Context(), tenantID, filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to list courses",
		})
	}

	// Convert to response format
	var responses []dto.CourseResponse
	for _, course := range courses {
		responses = append(responses, courseToResponse(course))
	}

	return c.JSON(dto.TrainingListResponse{
		Items:      interfaceSlice(responses),
		Total:      int64(len(responses)),
		Page:       1,
		PageSize:   len(responses),
		TotalPages: 1,
	})
}

// UpdateCourse godoc
// @Summary Update training course
// @Description Update training course by ID
// @Tags training
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param request body dto.UpdateCourseRequest true "Course data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/courses/{id} [put]
func (h *TrainingHandler) UpdateCourse(c *fiber.Ctx) error {
	courseID := c.Params("id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Course ID is required",
		})
	}

	var req dto.UpdateCourseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	userID := c.Locals("user_id").(string)

	err := h.trainingService.UpdateCourse(c.Context(), courseID, req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update course",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Course updated successfully",
	})
}

// DeleteCourse godoc
// @Summary Delete training course
// @Description Delete training course by ID
// @Tags training
// @Param id path string true "Course ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/courses/{id} [delete]
func (h *TrainingHandler) DeleteCourse(c *fiber.Ctx) error {
	courseID := c.Params("id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Course ID is required",
		})
	}

	userID := c.Locals("user_id").(string)

	err := h.trainingService.DeleteCourse(c.Context(), courseID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete course",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Course deleted successfully",
	})
}

// Course Materials handlers

// AddMaterialToCourse godoc
// @Summary Add material to course
// @Description Add a material to a training course
// @Tags training
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param material_id path string true "Material ID"
// @Param request body dto.CourseMaterialRequest true "Course material data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/courses/{id}/materials/{material_id} [post]
func (h *TrainingHandler) AddMaterialToCourse(c *fiber.Ctx) error {
	courseID := c.Params("id")
	materialID := c.Params("material_id")

	if courseID == "" || materialID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Course ID and Material ID are required",
		})
	}

	var req dto.CourseMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	userID := c.Locals("user_id").(string)

	err := h.trainingService.AddMaterialToCourse(c.Context(), courseID, materialID, req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to add material to course",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Material added to course successfully",
	})
}

// RemoveMaterialFromCourse godoc
// @Summary Remove material from course
// @Description Remove a material from a training course
// @Tags training
// @Param id path string true "Course ID"
// @Param material_id path string true "Material ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/courses/{id}/materials/{material_id} [delete]
func (h *TrainingHandler) RemoveMaterialFromCourse(c *fiber.Ctx) error {
	courseID := c.Params("id")
	materialID := c.Params("material_id")

	if courseID == "" || materialID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Course ID and Material ID are required",
		})
	}

	userID := c.Locals("user_id").(string)

	err := h.trainingService.RemoveMaterialFromCourse(c.Context(), courseID, materialID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to remove material from course",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Material removed from course successfully",
	})
}

// GetCourseMaterials godoc
// @Summary Get course materials
// @Description Get all materials in a training course
// @Tags training
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {array} dto.CourseMaterialResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/training/courses/{id}/materials [get]
func (h *TrainingHandler) GetCourseMaterials(c *fiber.Ctx) error {
	courseID := c.Params("id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Course ID is required",
		})
	}

	materials, err := h.trainingService.GetCourseMaterials(c.Context(), courseID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get course materials",
		})
	}

	// Convert to response format
	var responses []dto.CourseMaterialResponse
	for _, material := range materials {
		responses = append(responses, courseMaterialToResponse(material))
	}

	return c.JSON(responses)
}

// Helper functions for converting to response DTOs

func materialToResponse(material repo.Material) dto.MaterialResponse {
	return dto.MaterialResponse{
		ID:              material.ID,
		TenantID:        material.TenantID,
		Title:           material.Title,
		Description:     material.Description,
		URI:             material.URI,
		Type:            material.Type,
		MaterialType:    material.MaterialType,
		DurationMinutes: material.DurationMinutes,
		Tags:            material.Tags,
		IsRequired:      material.IsRequired,
		PassingScore:    material.PassingScore,
		AttemptsLimit:   material.AttemptsLimit,
		Metadata:        material.Metadata,
		CreatedBy:       material.CreatedBy,
		CreatedAt:       material.CreatedAt,
		UpdatedAt:       material.UpdatedAt,
	}
}

func courseToResponse(course repo.TrainingCourse) dto.CourseResponse {
	return dto.CourseResponse{
		ID:          course.ID,
		TenantID:    course.TenantID,
		Title:       course.Title,
		Description: course.Description,
		IsActive:    course.IsActive,
		CreatedBy:   course.CreatedBy,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}
}

func courseMaterialToResponse(courseMaterial repo.CourseMaterial) dto.CourseMaterialResponse {
	return dto.CourseMaterialResponse{
		ID:         courseMaterial.ID,
		CourseID:   courseMaterial.CourseID,
		MaterialID: courseMaterial.MaterialID,
		OrderIndex: courseMaterial.OrderIndex,
		IsRequired: courseMaterial.IsRequired,
		CreatedAt:  courseMaterial.CreatedAt,
	}
}

func interfaceSlice(slice interface{}) []interface{} {
	switch v := slice.(type) {
	case []dto.MaterialResponse:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result
	case []dto.CourseResponse:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result
	default:
		return []interface{}{}
	}
}

// Register registers training routes
func (h *TrainingHandler) Register(router fiber.Router) {
	training := router.Group("/training")

	// Materials routes
	materials := training.Group("/materials")
	materials.Post("/", h.CreateMaterial)
	materials.Get("/:id", h.GetMaterial)
	materials.Get("/", h.ListMaterials)
	materials.Put("/:id", h.UpdateMaterial)
	materials.Delete("/:id", h.DeleteMaterial)

	// Courses routes
	courses := training.Group("/courses")
	courses.Post("/", h.CreateCourse)
	courses.Get("/:id", h.GetCourse)
	courses.Get("/", h.ListCourses)
	courses.Put("/:id", h.UpdateCourse)
	courses.Delete("/:id", h.DeleteCourse)

	// Course materials routes
	courses.Post("/:id/materials/:material_id", h.AddMaterialToCourse)
	courses.Delete("/:id/materials/:material_id", h.RemoveMaterialFromCourse)
	courses.Get("/:id/materials", h.GetCourseMaterials)
}
