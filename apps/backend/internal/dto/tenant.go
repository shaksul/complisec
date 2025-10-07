package dto

import "time"

// CreateTenantDTO - DTO для создания организации
type CreateTenantDTO struct {
	Name   string  `json:"name" validate:"required,min=2,max=255"`
	Domain *string `json:"domain,omitempty" validate:"omitempty,min=3,max=255"`
}

// UpdateTenantDTO - DTO для обновления организации
type UpdateTenantDTO struct {
	Name   string  `json:"name" validate:"required,min=2,max=255"`
	Domain *string `json:"domain,omitempty" validate:"omitempty,min=3,max=255"`
}

// TenantResponse - DTO для ответа с данными организации
type TenantResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Domain    *string   `json:"domain,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TenantListResponse - DTO для списка организаций
type TenantListResponse struct {
	Data       []TenantResponse `json:"data"`
	Pagination PaginationInfo   `json:"pagination"`
}

// PaginationInfo - информация о пагинации
type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

