package dto

import "time"

type CreateUserRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=6"`
	FirstName string   `json:"first_name" validate:"required"`
	LastName  string   `json:"last_name" validate:"required"`
	RoleIDs   []string `json:"role_ids" validate:"dive,uuid"`
}

type UpdateUserRequest struct {
	FirstName *string  `json:"first_name,omitempty"`
	LastName  *string  `json:"last_name,omitempty"`
	IsActive  *bool    `json:"is_active,omitempty"`
	RoleIDs   []string `json:"role_ids,omitempty"`
}

// Role DTOs moved to dto/role.go

type UserResponse struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	FirstName   *string  `json:"first_name"`
	LastName    *string  `json:"last_name"`
	IsActive    bool     `json:"is_active"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

// RoleResponse and PermissionResponse moved to dto/role.go

type UserRoleRequest struct {
	UserID string `json:"user_id" validate:"required"`
	RoleID string `json:"role_id" validate:"required"`
}

type UserCatalogRequest struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	Search   string `json:"search" query:"search"`
	Role     string `json:"role" query:"role"`
	IsActive *bool  `json:"is_active" query:"is_active"`
	SortBy   string `json:"sort_by" query:"sort_by"`
	SortDir  string `json:"sort_dir" query:"sort_dir"`
}

type UserCatalogResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserDetailResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Stats     UserStats `json:"stats"`
}

type UserStats struct {
	DocumentsCount int `json:"documents_count"`
	RisksCount     int `json:"risks_count"`
	IncidentsCount int `json:"incidents_count"`
	AssetsCount    int `json:"assets_count"`
}
