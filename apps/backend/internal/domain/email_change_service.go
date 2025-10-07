package domain

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"time"

	"risknexus/backend/internal/repo"
)

type EmailChangeService struct {
	emailChangeRepo *repo.EmailChangeRepo
	userRepo        *repo.UserRepo
}

func NewEmailChangeService(emailChangeRepo *repo.EmailChangeRepo, userRepo *repo.UserRepo) *EmailChangeService {
	return &EmailChangeService{
		emailChangeRepo: emailChangeRepo,
		userRepo:        userRepo,
	}
}

// generateVerificationCode генерирует 6-значный код подтверждения
func (s *EmailChangeService) generateVerificationCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 6)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}

// RequestEmailChange создает запрос на смену email
func (s *EmailChangeService) RequestEmailChange(ctx context.Context, userID, tenantID, newEmail string) (*repo.EmailChangeRequest, error) {
	// Проверяем, что пользователь существует
	user, err := s.userRepo.GetByIDAndTenant(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Проверяем, что новый email отличается от текущего
	if user.Email == newEmail {
		return nil, errors.New("new email must be different from current email")
	}

	// Проверяем, что новый email не занят другим пользователем
	existingUser, err := s.userRepo.GetByEmailAndTenant(ctx, newEmail, tenantID)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email is already in use")
	}

	// Генерируем код подтверждения
	verificationCode, err := s.generateVerificationCode()
	if err != nil {
		return nil, err
	}

	// Создаем запрос (действителен 24 часа)
	expiresAt := time.Now().Add(24 * time.Hour)
	request, err := s.emailChangeRepo.CreateEmailChangeRequest(ctx, userID, tenantID, user.Email, newEmail, verificationCode, expiresAt)
	if err != nil {
		return nil, err
	}

	// Добавляем запись в аудит-лог
	err = s.emailChangeRepo.AddAuditLog(ctx, userID, tenantID, user.Email, newEmail, "requested", nil, nil)
	if err != nil {
		log.Printf("WARNING: Failed to add audit log for email change request: %v", err)
	}

	// TODO: Отправить email с кодом подтверждения на старый email
	log.Printf("Email change requested for user %s: %s -> %s, verification code: %s", userID, user.Email, newEmail, verificationCode)

	return request, nil
}

// VerifyOldEmail подтверждает старый email
func (s *EmailChangeService) VerifyOldEmail(ctx context.Context, requestID, verificationCode string) error {
	// Получаем запрос
	request, err := s.emailChangeRepo.GetEmailChangeRequestByID(ctx, requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("email change request not found")
	}

	// Проверяем, что запрос не истек
	if time.Now().After(request.ExpiresAt) {
		// Помечаем как истекший
		s.emailChangeRepo.ExpireEmailChangeRequests(ctx)
		return errors.New("verification code has expired")
	}

	// Проверяем статус
	if request.Status != "pending" {
		return errors.New("invalid request status")
	}

	// Подтверждаем старый email
	err = s.emailChangeRepo.VerifyOldEmail(ctx, requestID, verificationCode)
	if err != nil {
		return err
	}

	// Добавляем запись в аудит-лог
	err = s.emailChangeRepo.AddAuditLog(ctx, request.UserID, request.TenantID, request.OldEmail, request.NewEmail, "old_email_verified", nil, nil)
	if err != nil {
		log.Printf("WARNING: Failed to add audit log for old email verification: %v", err)
	}

	// TODO: Отправить email с кодом подтверждения на новый email
	log.Printf("Old email verified for user %s, sending verification to new email: %s", request.UserID, request.NewEmail)

	return nil
}

// VerifyNewEmail подтверждает новый email
func (s *EmailChangeService) VerifyNewEmail(ctx context.Context, requestID, verificationCode string) error {
	// Получаем запрос
	request, err := s.emailChangeRepo.GetEmailChangeRequestByID(ctx, requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("email change request not found")
	}

	// Проверяем, что запрос не истек
	if time.Now().After(request.ExpiresAt) {
		// Помечаем как истекший
		s.emailChangeRepo.ExpireEmailChangeRequests(ctx)
		return errors.New("verification code has expired")
	}

	// Проверяем статус
	if request.Status != "old_email_verified" {
		return errors.New("invalid request status")
	}

	// Подтверждаем новый email
	err = s.emailChangeRepo.VerifyNewEmail(ctx, requestID, verificationCode)
	if err != nil {
		return err
	}

	// Добавляем запись в аудит-лог
	err = s.emailChangeRepo.AddAuditLog(ctx, request.UserID, request.TenantID, request.OldEmail, request.NewEmail, "new_email_verified", nil, nil)
	if err != nil {
		log.Printf("WARNING: Failed to add audit log for new email verification: %v", err)
	}

	log.Printf("New email verified for user %s, email change can now be completed", request.UserID)

	return nil
}

// CompleteEmailChange завершает смену email
func (s *EmailChangeService) CompleteEmailChange(ctx context.Context, requestID string) error {
	// Получаем запрос
	request, err := s.emailChangeRepo.GetEmailChangeRequestByID(ctx, requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("email change request not found")
	}

	// Проверяем, что запрос не истек
	if time.Now().After(request.ExpiresAt) {
		// Помечаем как истекший
		s.emailChangeRepo.ExpireEmailChangeRequests(ctx)
		return errors.New("verification code has expired")
	}

	// Проверяем статус
	if request.Status != "new_email_verified" {
		return errors.New("invalid request status")
	}

	// Завершаем смену email
	err = s.emailChangeRepo.CompleteEmailChange(ctx, requestID)
	if err != nil {
		return err
	}

	// Добавляем запись в аудит-лог
	err = s.emailChangeRepo.AddAuditLog(ctx, request.UserID, request.TenantID, request.OldEmail, request.NewEmail, "completed", nil, nil)
	if err != nil {
		log.Printf("WARNING: Failed to add audit log for email change completion: %v", err)
	}

	log.Printf("Email change completed for user %s: %s -> %s", request.UserID, request.OldEmail, request.NewEmail)

	return nil
}

// CancelEmailChange отменяет запрос на смену email
func (s *EmailChangeService) CancelEmailChange(ctx context.Context, requestID string) error {
	// Получаем запрос
	request, err := s.emailChangeRepo.GetEmailChangeRequestByID(ctx, requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("email change request not found")
	}

	// Отменяем запрос
	err = s.emailChangeRepo.CancelEmailChange(ctx, requestID)
	if err != nil {
		return err
	}

	// Добавляем запись в аудит-лог
	err = s.emailChangeRepo.AddAuditLog(ctx, request.UserID, request.TenantID, request.OldEmail, request.NewEmail, "cancelled", nil, nil)
	if err != nil {
		log.Printf("WARNING: Failed to add audit log for email change cancellation: %v", err)
	}

	log.Printf("Email change cancelled for user %s", request.UserID)

	return nil
}

// GetActiveEmailChangeRequest получает активный запрос на смену email
func (s *EmailChangeService) GetActiveEmailChangeRequest(ctx context.Context, userID, tenantID string) (*repo.EmailChangeRequest, error) {
	// Сначала помечаем истекшие запросы
	err := s.emailChangeRepo.ExpireEmailChangeRequests(ctx)
	if err != nil {
		log.Printf("WARNING: Failed to expire email change requests: %v", err)
	}

	return s.emailChangeRepo.GetActiveEmailChangeRequest(ctx, userID, tenantID)
}

// GetEmailChangeAuditLogs получает аудит-лог изменений email
func (s *EmailChangeService) GetEmailChangeAuditLogs(ctx context.Context, userID, tenantID string, limit, offset int) ([]repo.EmailChangeAuditLog, error) {
	return s.emailChangeRepo.GetAuditLogs(ctx, userID, tenantID, limit, offset)
}

// ResendVerificationCode повторно отправляет код подтверждения
func (s *EmailChangeService) ResendVerificationCode(ctx context.Context, requestID string) error {
	// Получаем запрос
	request, err := s.emailChangeRepo.GetEmailChangeRequestByID(ctx, requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("email change request not found")
	}

	// Проверяем, что запрос не истек
	if time.Now().After(request.ExpiresAt) {
		// Помечаем как истекший
		s.emailChangeRepo.ExpireEmailChangeRequests(ctx)
		return errors.New("verification code has expired")
	}

	// Проверяем статус
	if request.Status != "pending" && request.Status != "old_email_verified" {
		return errors.New("invalid request status for resending verification code")
	}

	// Генерируем новый код
	newCode, err := s.generateVerificationCode()
	if err != nil {
		return err
	}

	// Обновляем код в базе данных
	err = s.emailChangeRepo.UpdateVerificationCode(ctx, requestID, newCode)
	if err != nil {
		return err
	}

	// TODO: Отправить новый код на соответствующий email
	if request.Status == "pending" {
		log.Printf("Resending verification code to old email for user %s: %s", request.UserID, request.OldEmail)
	} else {
		log.Printf("Resending verification code to new email for user %s: %s", request.UserID, request.NewEmail)
	}

	return nil
}
