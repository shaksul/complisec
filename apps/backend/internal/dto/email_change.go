package dto

import "time"

// RequestEmailChangeRequest запрос на создание запроса смены email
type RequestEmailChangeRequest struct {
	NewEmail string `json:"new_email" validate:"required,email"`
}

// RequestEmailChangeResponse ответ на создание запроса смены email
type RequestEmailChangeResponse struct {
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

// VerifyEmailRequest запрос на подтверждение email
type VerifyEmailRequest struct {
	RequestID        string `json:"request_id" validate:"required"`
	VerificationCode string `json:"verification_code" validate:"required,len=6"`
}

// VerifyEmailResponse ответ на подтверждение email
type VerifyEmailResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// CompleteEmailChangeRequest запрос на завершение смены email
type CompleteEmailChangeRequest struct {
	RequestID string `json:"request_id" validate:"required"`
}

// CompleteEmailChangeResponse ответ на завершение смены email
type CompleteEmailChangeResponse struct {
	Message string `json:"message"`
}

// CancelEmailChangeRequest запрос на отмену смены email
type CancelEmailChangeRequest struct {
	RequestID string `json:"request_id" validate:"required"`
}

// CancelEmailChangeResponse ответ на отмену смены email
type CancelEmailChangeResponse struct {
	Message string `json:"message"`
}

// ResendVerificationCodeRequest запрос на повторную отправку кода
type ResendVerificationCodeRequest struct {
	RequestID string `json:"request_id" validate:"required"`
}

// ResendVerificationCodeResponse ответ на повторную отправку кода
type ResendVerificationCodeResponse struct {
	Message string `json:"message"`
}

// EmailChangeRequestResponse информация о запросе на смену email
type EmailChangeRequestResponse struct {
	ID        string    `json:"id"`
	OldEmail  string    `json:"old_email"`
	NewEmail  string    `json:"new_email"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EmailChangeAuditLogResponse запись аудит-лога
type EmailChangeAuditLogResponse struct {
	ID         string    `json:"id"`
	OldEmail   string    `json:"old_email"`
	NewEmail   string    `json:"new_email"`
	ChangeType string    `json:"change_type"`
	IPAddress  *string   `json:"ip_address,omitempty"`
	UserAgent  *string   `json:"user_agent,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// EmailChangeStatusResponse статус запроса на смену email
type EmailChangeStatusResponse struct {
	HasActiveRequest bool                        `json:"has_active_request"`
	Request          *EmailChangeRequestResponse `json:"request,omitempty"`
}

// EmailChangeAuditLogsResponse список записей аудит-лога
type EmailChangeAuditLogsResponse struct {
	Logs  []EmailChangeAuditLogResponse `json:"logs"`
	Total int                           `json:"total"`
}

