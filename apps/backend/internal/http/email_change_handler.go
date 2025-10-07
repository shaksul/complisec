package http

import (
	"context"
	"strconv"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EmailChangeHandler struct {
	emailChangeService *domain.EmailChangeService
	validator          *validator.Validate
}

func NewEmailChangeHandler(emailChangeService *domain.EmailChangeService, validator *validator.Validate) *EmailChangeHandler {
	return &EmailChangeHandler{
		emailChangeService: emailChangeService,
		validator:          validator,
	}
}

func (h *EmailChangeHandler) Register(router fiber.Router) {
	api := router.Group("/email-change")

	api.Post("/request", h.requestEmailChange)
	api.Post("/verify-old", h.verifyOldEmail)
	api.Post("/verify-new", h.verifyNewEmail)
	api.Post("/complete", h.completeEmailChange)
	api.Post("/cancel", h.cancelEmailChange)
	api.Post("/resend", h.resendVerificationCode)
	api.Get("/status", h.getEmailChangeStatus)
	api.Get("/audit-logs", h.getAuditLogs)
}

// requestEmailChange создает запрос на смену email
func (h *EmailChangeHandler) requestEmailChange(c *fiber.Ctx) error {
	var req dto.RequestEmailChangeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id").(string)
	tenantID := c.Locals("tenant_id").(string)

	request, err := h.emailChangeService.RequestEmailChange(context.Background(), userID, tenantID, req.NewEmail)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	response := dto.RequestEmailChangeResponse{
		RequestID: request.ID,
		Message:   "Email change request created. Please check your current email for verification code.",
	}

	return c.Status(201).JSON(fiber.Map{"data": response})
}

// verifyOldEmail подтверждает старый email
func (h *EmailChangeHandler) verifyOldEmail(c *fiber.Ctx) error {
	var req dto.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.emailChangeService.VerifyOldEmail(context.Background(), req.RequestID, req.VerificationCode)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	response := dto.VerifyEmailResponse{
		Message: "Old email verified successfully. Please check your new email for verification code.",
		Status:  "old_email_verified",
	}

	return c.JSON(fiber.Map{"data": response})
}

// verifyNewEmail подтверждает новый email
func (h *EmailChangeHandler) verifyNewEmail(c *fiber.Ctx) error {
	var req dto.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.emailChangeService.VerifyNewEmail(context.Background(), req.RequestID, req.VerificationCode)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	response := dto.VerifyEmailResponse{
		Message: "New email verified successfully. You can now complete the email change.",
		Status:  "new_email_verified",
	}

	return c.JSON(fiber.Map{"data": response})
}

// completeEmailChange завершает смену email
func (h *EmailChangeHandler) completeEmailChange(c *fiber.Ctx) error {
	var req dto.CompleteEmailChangeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.emailChangeService.CompleteEmailChange(context.Background(), req.RequestID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	response := dto.CompleteEmailChangeResponse{
		Message: "Email changed successfully.",
	}

	return c.JSON(fiber.Map{"data": response})
}

// cancelEmailChange отменяет запрос на смену email
func (h *EmailChangeHandler) cancelEmailChange(c *fiber.Ctx) error {
	var req dto.CancelEmailChangeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.emailChangeService.CancelEmailChange(context.Background(), req.RequestID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	response := dto.CancelEmailChangeResponse{
		Message: "Email change request cancelled successfully.",
	}

	return c.JSON(fiber.Map{"data": response})
}

// resendVerificationCode повторно отправляет код подтверждения
func (h *EmailChangeHandler) resendVerificationCode(c *fiber.Ctx) error {
	var req dto.ResendVerificationCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.emailChangeService.ResendVerificationCode(context.Background(), req.RequestID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	response := dto.ResendVerificationCodeResponse{
		Message: "Verification code resent successfully.",
	}

	return c.JSON(fiber.Map{"data": response})
}

// getEmailChangeStatus получает статус активного запроса на смену email
func (h *EmailChangeHandler) getEmailChangeStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	tenantID := c.Locals("tenant_id").(string)

	request, err := h.emailChangeService.GetActiveEmailChangeRequest(context.Background(), userID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	response := dto.EmailChangeStatusResponse{
		HasActiveRequest: request != nil,
	}

	if request != nil {
		response.Request = &dto.EmailChangeRequestResponse{
			ID:        request.ID,
			OldEmail:  request.OldEmail,
			NewEmail:  request.NewEmail,
			Status:    request.Status,
			ExpiresAt: request.ExpiresAt,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
		}
	}

	return c.JSON(fiber.Map{"data": response})
}

// getAuditLogs получает аудит-лог изменений email
func (h *EmailChangeHandler) getAuditLogs(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	tenantID := c.Locals("tenant_id").(string)

	// Параметры пагинации
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	logs, err := h.emailChangeService.GetEmailChangeAuditLogs(context.Background(), userID, tenantID, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var responseLogs []dto.EmailChangeAuditLogResponse
	for _, log := range logs {
		responseLogs = append(responseLogs, dto.EmailChangeAuditLogResponse{
			ID:         log.ID,
			OldEmail:   log.OldEmail,
			NewEmail:   log.NewEmail,
			ChangeType: log.ChangeType,
			IPAddress:  log.IPAddress,
			UserAgent:  log.UserAgent,
			CreatedAt:  log.CreatedAt,
		})
	}

	response := dto.EmailChangeAuditLogsResponse{
		Logs:  responseLogs,
		Total: len(responseLogs),
	}

	return c.JSON(fiber.Map{"data": response})
}
