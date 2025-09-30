package domain

import "errors"

// Предопределенные ошибки для ролей и прав
var (
	// Ошибки ролей
	ErrRoleNotFound       = errors.New("role not found")
	ErrRoleAlreadyExists  = errors.New("role with this name already exists")
	ErrRoleInUse          = errors.New("cannot delete role: it is assigned to users")
	ErrInvalidRoleName    = errors.New("invalid role name")
	ErrInvalidDescription = errors.New("invalid role description")

	// Ошибки прав
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrInvalidPermissionID = errors.New("invalid permission ID")

	// Ошибки пользователей
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyHasRole  = errors.New("user already has this role")
	ErrUserDoesNotHaveRole = errors.New("user does not have this role")

	// Ошибки валидации
	ErrValidationFailed = errors.New("validation failed")
	ErrEmptyField       = errors.New("field cannot be empty")
	ErrFieldTooLong     = errors.New("field exceeds maximum length")
	ErrFieldTooShort    = errors.New("field is too short")
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// BusinessError представляет бизнес-ошибку
type BusinessError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

func (e BusinessError) Error() string {
	return e.Message
}

// NewValidationError создает новую ошибку валидации
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// NewBusinessError создает новую бизнес-ошибку
func NewBusinessError(code, message string, details map[string]interface{}) BusinessError {
	return BusinessError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
