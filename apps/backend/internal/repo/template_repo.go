package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"risknexus/backend/internal/dto"
)

// DocumentTemplate entity
type DocumentTemplate struct {
	ID           string
	TenantID     string
	Name         string
	Description  *string
	TemplateType string
	Content      string
	IsSystem     bool
	IsActive     bool
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// InventoryNumberRule entity
type InventoryNumberRule struct {
	ID              string
	TenantID        string
	AssetType       string
	AssetClass      *string
	Pattern         string
	CurrentSequence int
	Prefix          *string
	Description     *string
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// TemplateRepo handles document template operations
type TemplateRepo struct {
	db DBInterface
}

// NewTemplateRepo creates a new template repository
func NewTemplateRepo(db DBInterface) *TemplateRepo {
	return &TemplateRepo{db: db}
}

// ListTemplates returns templates for a tenant with optional filters
func (r *TemplateRepo) ListTemplates(ctx context.Context, tenantID string, filters map[string]interface{}) ([]DocumentTemplate, error) {
	query := `
		SELECT id, tenant_id, name, description, template_type, content, is_system, is_active,
		       created_by, created_at, updated_at, deleted_at
		FROM document_templates
		WHERE tenant_id = $1 AND deleted_at IS NULL`

	args := []interface{}{tenantID}
	argIndex := 2

	// Apply filters
	if templateType, ok := filters["template_type"].(string); ok && templateType != "" {
		query += fmt.Sprintf(" AND template_type = $%d", argIndex)
		args = append(args, templateType)
		argIndex++
	}

	if isSystem, ok := filters["is_system"].(bool); ok {
		query += fmt.Sprintf(" AND is_system = $%d", argIndex)
		args = append(args, isSystem)
		argIndex++
	}

	if isActive, ok := filters["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	query += " ORDER BY is_system DESC, name ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	var templates []DocumentTemplate
	for rows.Next() {
		var t DocumentTemplate
		err := rows.Scan(&t.ID, &t.TenantID, &t.Name, &t.Description, &t.TemplateType,
			&t.Content, &t.IsSystem, &t.IsActive, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}
		templates = append(templates, t)
	}

	return templates, nil
}

// GetTemplateByID retrieves a template by ID
func (r *TemplateRepo) GetTemplateByID(ctx context.Context, id, tenantID string) (*DocumentTemplate, error) {
	query := `
		SELECT id, tenant_id, name, description, template_type, content, is_system, is_active,
		       created_by, created_at, updated_at, deleted_at
		FROM document_templates
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`

	var t DocumentTemplate
	err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&t.ID, &t.TenantID, &t.Name, &t.Description, &t.TemplateType,
		&t.Content, &t.IsSystem, &t.IsActive, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return &t, nil
}

// CreateTemplate creates a new template
func (r *TemplateRepo) CreateTemplate(ctx context.Context, template *DocumentTemplate) error {
	query := `
		INSERT INTO document_templates (id, tenant_id, name, description, template_type, content, 
		                                is_system, is_active, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		template.ID, template.TenantID, template.Name, template.Description,
		template.TemplateType, template.Content, template.IsSystem, template.IsActive,
		template.CreatedBy, template.CreatedAt, template.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	return nil
}

// UpdateTemplate updates an existing template
func (r *TemplateRepo) UpdateTemplate(ctx context.Context, template *DocumentTemplate) error {
	query := `
		UPDATE document_templates
		SET name = $1, description = $2, template_type = $3, content = $4,
		    is_active = $5, updated_at = $6
		WHERE id = $7 AND tenant_id = $8 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		template.Name, template.Description, template.TemplateType, template.Content,
		template.IsActive, time.Now(), template.ID, template.TenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("template not found or already deleted")
	}

	return nil
}

// DeleteTemplate soft deletes a template
func (r *TemplateRepo) DeleteTemplate(ctx context.Context, id, tenantID string) error {
	query := `
		UPDATE document_templates
		SET deleted_at = $1, updated_at = $1
		WHERE id = $2 AND tenant_id = $3 AND deleted_at IS NULL AND is_system = false`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("template not found, already deleted, or is a system template")
	}

	return nil
}

// ListInventoryRules returns inventory number rules for a tenant
func (r *TemplateRepo) ListInventoryRules(ctx context.Context, tenantID string) ([]InventoryNumberRule, error) {
	query := `
		SELECT id, tenant_id, asset_type, asset_class, pattern, current_sequence, prefix,
		       description, is_active, created_at, updated_at
		FROM inventory_number_rules
		WHERE tenant_id = $1 AND is_active = true
		ORDER BY asset_type, asset_class`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory rules: %w", err)
	}
	defer rows.Close()

	var rules []InventoryNumberRule
	for rows.Next() {
		var rule InventoryNumberRule
		err := rows.Scan(&rule.ID, &rule.TenantID, &rule.AssetType, &rule.AssetClass,
			&rule.Pattern, &rule.CurrentSequence, &rule.Prefix, &rule.Description,
			&rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory rule: %w", err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// GetInventoryRuleByType retrieves a rule by asset type and class
func (r *TemplateRepo) GetInventoryRuleByType(ctx context.Context, tenantID, assetType string, assetClass *string) (*InventoryNumberRule, error) {
	var query string
	var args []interface{}

	if assetClass != nil && *assetClass != "" {
		query = `
			SELECT id, tenant_id, asset_type, asset_class, pattern, current_sequence, prefix,
			       description, is_active, created_at, updated_at
			FROM inventory_number_rules
			WHERE tenant_id = $1 AND asset_type = $2 AND asset_class = $3 AND is_active = true`
		args = []interface{}{tenantID, assetType, *assetClass}
	} else {
		query = `
			SELECT id, tenant_id, asset_type, asset_class, pattern, current_sequence, prefix,
			       description, is_active, created_at, updated_at
			FROM inventory_number_rules
			WHERE tenant_id = $1 AND asset_type = $2 AND (asset_class IS NULL OR asset_class = '') AND is_active = true`
		args = []interface{}{tenantID, assetType}
	}

	var rule InventoryNumberRule
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&rule.ID, &rule.TenantID, &rule.AssetType, &rule.AssetClass,
		&rule.Pattern, &rule.CurrentSequence, &rule.Prefix, &rule.Description,
		&rule.IsActive, &rule.CreatedAt, &rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory rule: %w", err)
	}

	return &rule, nil
}

// CreateInventoryRule creates a new inventory number rule
func (r *TemplateRepo) CreateInventoryRule(ctx context.Context, rule *InventoryNumberRule) error {
	query := `
		INSERT INTO inventory_number_rules (id, tenant_id, asset_type, asset_class, pattern,
		                                    current_sequence, prefix, description, is_active,
		                                    created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		rule.ID, rule.TenantID, rule.AssetType, rule.AssetClass, rule.Pattern,
		rule.CurrentSequence, rule.Prefix, rule.Description, rule.IsActive,
		rule.CreatedAt, rule.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create inventory rule: %w", err)
	}

	return nil
}

// UpdateInventoryRule updates an existing rule
func (r *TemplateRepo) UpdateInventoryRule(ctx context.Context, rule *InventoryNumberRule) error {
	query := `
		UPDATE inventory_number_rules
		SET pattern = $1, current_sequence = $2, prefix = $3, description = $4,
		    is_active = $5, updated_at = $6
		WHERE id = $7 AND tenant_id = $8`

	result, err := r.db.ExecContext(ctx, query,
		rule.Pattern, rule.CurrentSequence, rule.Prefix, rule.Description,
		rule.IsActive, time.Now(), rule.ID, rule.TenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update inventory rule: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("inventory rule not found")
	}

	return nil
}

// IncrementSequence increments the current sequence and returns the new value
func (r *TemplateRepo) IncrementSequence(ctx context.Context, ruleID string) (int, error) {
	// Increment and return new value atomically
	query := `
		UPDATE inventory_number_rules
		SET current_sequence = current_sequence + 1, updated_at = $1
		WHERE id = $2
		RETURNING current_sequence`

	var newSequence int
	err := r.db.QueryRowContext(ctx, query, time.Now(), ruleID).Scan(&newSequence)
	if err != nil {
		return 0, fmt.Errorf("failed to increment sequence: %w", err)
	}

	return newSequence, nil
}

// ToDTO converts DocumentTemplate entity to DTO
func (t *DocumentTemplate) ToDTO() dto.DocumentTemplateDTO {
	return dto.DocumentTemplateDTO{
		ID:           t.ID,
		TenantID:     t.TenantID,
		Name:         t.Name,
		Description:  t.Description,
		TemplateType: t.TemplateType,
		Content:      t.Content,
		IsSystem:     t.IsSystem,
		IsActive:     t.IsActive,
		CreatedBy:    t.CreatedBy,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

// ToDTO converts InventoryNumberRule entity to DTO
func (r *InventoryNumberRule) ToDTO() dto.InventoryNumberRuleDTO {
	return dto.InventoryNumberRuleDTO{
		ID:              r.ID,
		TenantID:        r.TenantID,
		AssetType:       r.AssetType,
		AssetClass:      r.AssetClass,
		Pattern:         r.Pattern,
		CurrentSequence: r.CurrentSequence,
		Prefix:          r.Prefix,
		Description:     r.Description,
		IsActive:        r.IsActive,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
