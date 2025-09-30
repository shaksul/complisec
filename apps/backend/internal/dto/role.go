package dto

import "time"

// CreateRoleRequest DTO для создания роли
type CreateRoleRequest struct {
	Name          string   `json:"name" validate:"required,min=1,max=100"`
	Description   string   `json:"description" validate:"max=500"`
	PermissionIDs []string `json:"permission_ids" validate:"dive,uuid"`
}

// UpdateRoleRequest DTO для обновления роли
type UpdateRoleRequest struct {
	Name          *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description   *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	PermissionIDs []string `json:"permission_ids,omitempty" validate:"dive,uuid"`
}

// RoleResponse DTO для ответа с информацией о роли
type RoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RoleWithPermissionsResponse DTO для роли с правами
type RoleWithPermissionsResponse struct {
	RoleResponse
	Permissions []string `json:"permissions"`
}

// PermissionResponse DTO для информации о праве
type PermissionResponse struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Module      string  `json:"module"`
	Description *string `json:"description"`
}

// RoleListResponse DTO для списка ролей с пагинацией
type RoleListResponse struct {
	Roles      []RoleResponse `json:"roles"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// AssignRoleRequest DTO для назначения роли пользователю
type AssignRoleRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	RoleID string `json:"role_id" validate:"required,uuid"`
}

// RoleFilter DTO для фильтрации ролей
type RoleFilter struct {
	Name      string `json:"name,omitempty"`
	Module    string `json:"module,omitempty"`
	HasUsers  *bool  `json:"has_users,omitempty"`
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}
