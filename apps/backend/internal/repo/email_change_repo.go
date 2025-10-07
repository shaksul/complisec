package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type EmailChangeRequest struct {
	ID               string
	UserID           string
	TenantID         string
	OldEmail         string
	NewEmail         string
	VerificationCode string
	OldEmailVerified bool
	NewEmailVerified bool
	ExpiresAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Status           string
}

type EmailChangeAuditLog struct {
	ID         string
	UserID     string
	TenantID   string
	OldEmail   string
	NewEmail   string
	ChangeType string
	IPAddress  *string
	UserAgent  *string
	CreatedAt  time.Time
}

type EmailChangeRepo struct {
	db *DB
}

func NewEmailChangeRepo(db *DB) *EmailChangeRepo {
	return &EmailChangeRepo{db: db}
}

// CreateEmailChangeRequest создает новый запрос на смену email
func (r *EmailChangeRepo) CreateEmailChangeRequest(ctx context.Context, userID, tenantID, oldEmail, newEmail, verificationCode string, expiresAt time.Time) (*EmailChangeRequest, error) {
	// Сначала отменяем все активные запросы для этого пользователя
	_, err := r.db.ExecContext(ctx, `
		UPDATE email_change_requests 
		SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND tenant_id = $2 
		AND status IN ('pending', 'old_email_verified', 'new_email_verified')
	`, userID, tenantID)
	if err != nil {
		return nil, err
	}

	// Создаем новый запрос
	request := &EmailChangeRequest{
		ID:               generateUUID(),
		UserID:           userID,
		TenantID:         tenantID,
		OldEmail:         oldEmail,
		NewEmail:         newEmail,
		VerificationCode: verificationCode,
		ExpiresAt:        expiresAt,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           "pending",
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO email_change_requests 
		(id, user_id, tenant_id, old_email, new_email, verification_code, expires_at, created_at, updated_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, request.ID, request.UserID, request.TenantID, request.OldEmail, request.NewEmail,
		request.VerificationCode, request.ExpiresAt, request.CreatedAt, request.UpdatedAt, request.Status)

	if err != nil {
		return nil, err
	}

	return request, nil
}

// GetEmailChangeRequestByID получает запрос по ID
func (r *EmailChangeRepo) GetEmailChangeRequestByID(ctx context.Context, id string) (*EmailChangeRequest, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, tenant_id, old_email, new_email, verification_code, 
		       old_email_verified, new_email_verified, expires_at, created_at, updated_at, status
		FROM email_change_requests WHERE id = $1
	`, id)

	var request EmailChangeRequest
	err := row.Scan(&request.ID, &request.UserID, &request.TenantID, &request.OldEmail,
		&request.NewEmail, &request.VerificationCode, &request.OldEmailVerified,
		&request.NewEmailVerified, &request.ExpiresAt, &request.CreatedAt, &request.UpdatedAt, &request.Status)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &request, nil
}

// GetActiveEmailChangeRequest получает активный запрос для пользователя
func (r *EmailChangeRepo) GetActiveEmailChangeRequest(ctx context.Context, userID, tenantID string) (*EmailChangeRequest, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, tenant_id, old_email, new_email, verification_code, 
		       old_email_verified, new_email_verified, expires_at, created_at, updated_at, status
		FROM email_change_requests 
		WHERE user_id = $1 AND tenant_id = $2 
		AND status IN ('pending', 'old_email_verified', 'new_email_verified')
		ORDER BY created_at DESC
		LIMIT 1
	`, userID, tenantID)

	var request EmailChangeRequest
	err := row.Scan(&request.ID, &request.UserID, &request.TenantID, &request.OldEmail,
		&request.NewEmail, &request.VerificationCode, &request.OldEmailVerified,
		&request.NewEmailVerified, &request.ExpiresAt, &request.CreatedAt, &request.UpdatedAt, &request.Status)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &request, nil
}

// VerifyOldEmail подтверждает старый email
func (r *EmailChangeRepo) VerifyOldEmail(ctx context.Context, requestID, verificationCode string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE email_change_requests 
		SET old_email_verified = TRUE, status = 'old_email_verified', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND verification_code = $2 AND status = 'pending'
	`, requestID, verificationCode)
	return err
}

// VerifyNewEmail подтверждает новый email
func (r *EmailChangeRepo) VerifyNewEmail(ctx context.Context, requestID, verificationCode string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE email_change_requests 
		SET new_email_verified = TRUE, status = 'new_email_verified', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND verification_code = $2 AND status = 'old_email_verified'
	`, requestID, verificationCode)
	return err
}

// CompleteEmailChange завершает смену email
func (r *EmailChangeRepo) CompleteEmailChange(ctx context.Context, requestID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Получаем данные запроса
	var userID, newEmail string
	err = tx.QueryRowContext(ctx, `
		SELECT user_id, new_email FROM email_change_requests 
		WHERE id = $1 AND status = 'new_email_verified'
	`, requestID).Scan(&userID, &newEmail)
	if err != nil {
		return err
	}

	// Обновляем email пользователя
	_, err = tx.ExecContext(ctx, `
		UPDATE users SET email = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2
	`, newEmail, userID)
	if err != nil {
		return err
	}

	// Помечаем запрос как завершенный
	_, err = tx.ExecContext(ctx, `
		UPDATE email_change_requests 
		SET status = 'completed', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, requestID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CancelEmailChange отменяет запрос на смену email
func (r *EmailChangeRepo) CancelEmailChange(ctx context.Context, requestID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE email_change_requests 
		SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, requestID)
	return err
}

// ExpireEmailChangeRequests помечает истекшие запросы как expired
func (r *EmailChangeRepo) ExpireEmailChangeRequests(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE email_change_requests 
		SET status = 'expired', updated_at = CURRENT_TIMESTAMP
		WHERE expires_at < CURRENT_TIMESTAMP 
		AND status IN ('pending', 'old_email_verified', 'new_email_verified')
	`)
	return err
}

// AddAuditLog добавляет запись в аудит-лог
func (r *EmailChangeRepo) AddAuditLog(ctx context.Context, userID, tenantID, oldEmail, newEmail, changeType string, ipAddress, userAgent *string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO email_change_audit_log 
		(id, user_id, tenant_id, old_email, new_email, change_type, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
	`, generateUUID(), userID, tenantID, oldEmail, newEmail, changeType, ipAddress, userAgent)
	return err
}

// GetAuditLogs получает аудит-лог для пользователя
func (r *EmailChangeRepo) GetAuditLogs(ctx context.Context, userID, tenantID string, limit, offset int) ([]EmailChangeAuditLog, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, tenant_id, old_email, new_email, change_type, ip_address, user_agent, created_at
		FROM email_change_audit_log 
		WHERE user_id = $1 AND tenant_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, userID, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []EmailChangeAuditLog
	for rows.Next() {
		var log EmailChangeAuditLog
		err := rows.Scan(&log.ID, &log.UserID, &log.TenantID, &log.OldEmail,
			&log.NewEmail, &log.ChangeType, &log.IPAddress, &log.UserAgent, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// UpdateVerificationCode обновляет код подтверждения
func (r *EmailChangeRepo) UpdateVerificationCode(ctx context.Context, requestID, newCode string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE email_change_requests 
		SET verification_code = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, newCode, requestID)
	return err
}
