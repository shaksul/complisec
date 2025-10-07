package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Risk struct {
	ID          string
	TenantID    string
	Title       string
	Description *string
	Category    *string
	Likelihood  *int
	Impact      *int
	Level       *int
	Status      string
	OwnerUserID *string
	AssetID     *string
	Methodology *string
	Strategy    *string
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RiskControl represents a control associated with a risk
type RiskControl struct {
	ID                   string
	RiskID               string
	ControlID            string
	ControlName          string
	ControlType          string
	ImplementationStatus string
	Effectiveness        *string
	Description          *string
	CreatedBy            *string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// RiskComment represents a comment on a risk
type RiskComment struct {
	ID         string
	RiskID     string
	UserID     string
	Comment    string
	IsInternal bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserName   *string // joined from users table
}

// RiskHistory represents a history entry for a risk
type RiskHistory struct {
	ID            string
	RiskID        string
	FieldChanged  string
	OldValue      *string
	NewValue      *string
	ChangeReason  *string
	ChangedBy     string
	ChangedAt     time.Time
	ChangedByName *string // joined from users table
}

// RiskAttachment represents an attachment to a risk
type RiskAttachment struct {
	ID             string
	RiskID         string
	FileName       string
	FilePath       string
	FileSize       int64
	MimeType       string
	FileHash       *string
	Description    *string
	UploadedBy     string
	UploadedAt     time.Time
	UploadedByName *string // joined from users table
}

// RiskTag represents a tag for a risk
type RiskTag struct {
	ID        string
	RiskID    string
	TagName   string
	TagColor  string
	CreatedBy *string
	CreatedAt time.Time
}

type RiskRepo struct {
	db *DB
}

func NewRiskRepo(db *DB) *RiskRepo {
	return &RiskRepo{db: db}
}

func (r *RiskRepo) Create(ctx context.Context, risk Risk) error {
	_, err := r.db.Exec(`
		INSERT INTO risks (id, tenant_id, title, description, category, likelihood, impact, status, owner_user_id, asset_id, methodology, strategy, due_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, risk.ID, risk.TenantID, risk.Title, risk.Description, risk.Category, risk.Likelihood, risk.Impact, risk.Status, risk.OwnerUserID, risk.AssetID, risk.Methodology, risk.Strategy, risk.DueDate)
	return err
}

func (r *RiskRepo) GetByID(ctx context.Context, id string) (*Risk, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, title, description, category, likelihood, impact, level, status, owner_user_id, asset_id, methodology, strategy, due_date, created_at, updated_at
		FROM risks WHERE id = $1
	`, id)

	var risk Risk
	err := row.Scan(&risk.ID, &risk.TenantID, &risk.Title, &risk.Description, &risk.Category, &risk.Likelihood, &risk.Impact, &risk.Level, &risk.Status, &risk.OwnerUserID, &risk.AssetID, &risk.Methodology, &risk.Strategy, &risk.DueDate, &risk.CreatedAt, &risk.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &risk, nil
}

func (r *RiskRepo) GetByIDWithTenant(ctx context.Context, id, tenantID string) (*Risk, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, title, description, category, likelihood, impact, level, status, owner_user_id, asset_id, methodology, strategy, due_date, created_at, updated_at
		FROM risks WHERE id = $1 AND tenant_id = $2
	`, id, tenantID)

	var risk Risk
	err := row.Scan(&risk.ID, &risk.TenantID, &risk.Title, &risk.Description, &risk.Category, &risk.Likelihood, &risk.Impact, &risk.Level, &risk.Status, &risk.OwnerUserID, &risk.AssetID, &risk.Methodology, &risk.Strategy, &risk.DueDate, &risk.CreatedAt, &risk.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &risk, nil
}

func (r *RiskRepo) List(ctx context.Context, tenantID string) ([]Risk, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, title, description, category, likelihood, impact, level, status, owner_user_id, asset_id, methodology, strategy, due_date, created_at, updated_at
		FROM risks WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var risks []Risk
	for rows.Next() {
		var risk Risk
		err := rows.Scan(&risk.ID, &risk.TenantID, &risk.Title, &risk.Description, &risk.Category, &risk.Likelihood, &risk.Impact, &risk.Level, &risk.Status, &risk.OwnerUserID, &risk.AssetID, &risk.Methodology, &risk.Strategy, &risk.DueDate, &risk.CreatedAt, &risk.UpdatedAt)
		if err != nil {
			return nil, err
		}
		risks = append(risks, risk)
	}
	return risks, nil
}

func (r *RiskRepo) ListWithFilters(ctx context.Context, tenantID string, filters map[string]interface{}, sortField, sortDirection string) ([]Risk, error) {
	query := `
		SELECT id, tenant_id, title, description, category, likelihood, impact, level, status, owner_user_id, asset_id, methodology, strategy, due_date, created_at, updated_at
		FROM risks WHERE tenant_id = $1
	`
	args := []interface{}{tenantID}
	argIndex := 2

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}
	if levelRange, ok := filters["level_range"].([]int); ok && len(levelRange) == 2 {
		query += fmt.Sprintf(" AND level BETWEEN $%d AND $%d", argIndex, argIndex+1)
		args = append(args, levelRange[0], levelRange[1])
		argIndex += 2
	}
	if levelExact, ok := filters["level_exact"].(int); ok {
		query += fmt.Sprintf(" AND level = $%d", argIndex)
		args = append(args, levelExact)
		argIndex++
	}
	if ownerUserID, ok := filters["owner_user_id"].(string); ok && ownerUserID != "" {
		query += fmt.Sprintf(" AND owner_user_id = $%d", argIndex)
		args = append(args, ownerUserID)
		argIndex++
	}
	if methodology, ok := filters["methodology"].(string); ok && methodology != "" {
		query += fmt.Sprintf(" AND methodology = $%d", argIndex)
		args = append(args, methodology)
		argIndex++
	}
	if strategy, ok := filters["strategy"].(string); ok && strategy != "" {
		query += fmt.Sprintf(" AND strategy = $%d", argIndex)
		args = append(args, strategy)
		argIndex++
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	// Apply sorting
	validSortFields := map[string]bool{
		"level":      true,
		"created_at": true,
		"category":   true,
		"title":      true,
		"status":     true,
	}
	if !validSortFields[sortField] {
		sortField = "level"
	}
	if sortDirection != "asc" && sortDirection != "desc" {
		sortDirection = "desc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortDirection)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var risks []Risk
	for rows.Next() {
		var risk Risk
		err := rows.Scan(&risk.ID, &risk.TenantID, &risk.Title, &risk.Description, &risk.Category, &risk.Likelihood, &risk.Impact, &risk.Level, &risk.Status, &risk.OwnerUserID, &risk.AssetID, &risk.Methodology, &risk.Strategy, &risk.DueDate, &risk.CreatedAt, &risk.UpdatedAt)
		if err != nil {
			return nil, err
		}
		risks = append(risks, risk)
	}
	return risks, nil
}

func (r *RiskRepo) Update(ctx context.Context, risk Risk) error {
	_, err := r.db.Exec(`
		UPDATE risks SET title = $1, description = $2, category = $3, likelihood = $4, impact = $5, status = $6, owner_user_id = $7, asset_id = $8, methodology = $9, strategy = $10, due_date = $11, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
	`, risk.Title, risk.Description, risk.Category, risk.Likelihood, risk.Impact, risk.Status, risk.OwnerUserID, risk.AssetID, risk.Methodology, risk.Strategy, risk.DueDate, risk.ID)
	return err
}

func (r *RiskRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(`
		UPDATE risks SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, status, id)
	return err
}

func (r *RiskRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec("DELETE FROM risks WHERE id = $1", id)
	return err
}

// Risk Controls methods
func (r *RiskRepo) AddControl(ctx context.Context, control RiskControl) error {
	_, err := r.db.Exec(`
		INSERT INTO risk_controls (id, risk_id, control_id, control_name, control_type, implementation_status, effectiveness, description, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, control.ID, control.RiskID, control.ControlID, control.ControlName, control.ControlType, control.ImplementationStatus, control.Effectiveness, control.Description, control.CreatedBy)
	return err
}

func (r *RiskRepo) GetControls(ctx context.Context, riskID string) ([]RiskControl, error) {
	rows, err := r.db.Query(`
		SELECT id, risk_id, control_id, control_name, control_type, implementation_status, effectiveness, description, created_by, created_at, updated_at
		FROM risk_controls WHERE risk_id = $1 ORDER BY created_at DESC
	`, riskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var controls []RiskControl
	for rows.Next() {
		var control RiskControl
		err := rows.Scan(&control.ID, &control.RiskID, &control.ControlID, &control.ControlName, &control.ControlType, &control.ImplementationStatus, &control.Effectiveness, &control.Description, &control.CreatedBy, &control.CreatedAt, &control.UpdatedAt)
		if err != nil {
			return nil, err
		}
		controls = append(controls, control)
	}
	return controls, nil
}

func (r *RiskRepo) UpdateControl(ctx context.Context, control RiskControl) error {
	_, err := r.db.Exec(`
		UPDATE risk_controls SET control_name = $1, control_type = $2, implementation_status = $3, effectiveness = $4, description = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`, control.ControlName, control.ControlType, control.ImplementationStatus, control.Effectiveness, control.Description, control.ID)
	return err
}

func (r *RiskRepo) DeleteControl(ctx context.Context, controlID string) error {
	_, err := r.db.Exec("DELETE FROM risk_controls WHERE id = $1", controlID)
	return err
}

// Risk Comments methods
func (r *RiskRepo) AddComment(ctx context.Context, comment RiskComment) error {
	_, err := r.db.Exec(`
		INSERT INTO risk_comments (id, risk_id, user_id, comment, is_internal)
		VALUES ($1, $2, $3, $4, $5)
	`, comment.ID, comment.RiskID, comment.UserID, comment.Comment, comment.IsInternal)
	return err
}

func (r *RiskRepo) GetComments(ctx context.Context, riskID string, includeInternal bool) ([]RiskComment, error) {
	query := `
		SELECT rc.id, rc.risk_id, rc.user_id, rc.comment, rc.is_internal, rc.created_at, rc.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, u.email) as user_name
		FROM risk_comments rc
		LEFT JOIN users u ON rc.user_id = u.id
		WHERE rc.risk_id = $1`

	args := []interface{}{riskID}
	if !includeInternal {
		query += " AND rc.is_internal = false"
	}
	query += " ORDER BY rc.created_at ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []RiskComment
	for rows.Next() {
		var comment RiskComment
		err := rows.Scan(&comment.ID, &comment.RiskID, &comment.UserID, &comment.Comment, &comment.IsInternal, &comment.CreatedAt, &comment.UpdatedAt, &comment.UserName)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// Risk History methods
func (r *RiskRepo) AddHistory(ctx context.Context, history RiskHistory) error {
	_, err := r.db.Exec(`
		INSERT INTO risk_history (id, risk_id, field_changed, old_value, new_value, change_reason, changed_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, history.ID, history.RiskID, history.FieldChanged, history.OldValue, history.NewValue, history.ChangeReason, history.ChangedBy)
	return err
}

func (r *RiskRepo) GetHistory(ctx context.Context, riskID string) ([]RiskHistory, error) {
	rows, err := r.db.Query(`
		SELECT rh.id, rh.risk_id, rh.field_changed, rh.old_value, rh.new_value, rh.change_reason, rh.changed_by, rh.changed_at,
		       COALESCE(u.first_name || ' ' || u.last_name, u.email) as changed_by_name
		FROM risk_history rh
		LEFT JOIN users u ON rh.changed_by = u.id
		WHERE rh.risk_id = $1 ORDER BY rh.changed_at DESC
	`, riskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []RiskHistory
	for rows.Next() {
		var h RiskHistory
		err := rows.Scan(&h.ID, &h.RiskID, &h.FieldChanged, &h.OldValue, &h.NewValue, &h.ChangeReason, &h.ChangedBy, &h.ChangedAt, &h.ChangedByName)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

// Risk Attachments methods
func (r *RiskRepo) AddAttachment(ctx context.Context, attachment RiskAttachment) error {
	_, err := r.db.Exec(`
		INSERT INTO risk_attachments (id, risk_id, file_name, file_path, file_size, mime_type, file_hash, description, uploaded_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, attachment.ID, attachment.RiskID, attachment.FileName, attachment.FilePath, attachment.FileSize, attachment.MimeType, attachment.FileHash, attachment.Description, attachment.UploadedBy)
	return err
}

func (r *RiskRepo) GetAttachments(ctx context.Context, riskID string) ([]RiskAttachment, error) {
	rows, err := r.db.Query(`
		SELECT ra.id, ra.risk_id, ra.file_name, ra.file_path, ra.file_size, ra.mime_type, ra.file_hash, ra.description, ra.uploaded_by, ra.uploaded_at,
		       COALESCE(u.first_name || ' ' || u.last_name, u.email) as uploaded_by_name
		FROM risk_attachments ra
		LEFT JOIN users u ON ra.uploaded_by = u.id
		WHERE ra.risk_id = $1 ORDER BY ra.uploaded_at DESC
	`, riskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []RiskAttachment
	for rows.Next() {
		var attachment RiskAttachment
		err := rows.Scan(&attachment.ID, &attachment.RiskID, &attachment.FileName, &attachment.FilePath, &attachment.FileSize, &attachment.MimeType, &attachment.FileHash, &attachment.Description, &attachment.UploadedBy, &attachment.UploadedAt, &attachment.UploadedByName)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}
	return attachments, nil
}

func (r *RiskRepo) DeleteAttachment(ctx context.Context, attachmentID string) error {
	_, err := r.db.Exec("DELETE FROM risk_attachments WHERE id = $1", attachmentID)
	return err
}

// Risk Tags methods
func (r *RiskRepo) AddTag(ctx context.Context, tag RiskTag) error {
	_, err := r.db.Exec(`
		INSERT INTO risk_tags (id, risk_id, tag_name, tag_color, created_by)
		VALUES ($1, $2, $3, $4, $5)
	`, tag.ID, tag.RiskID, tag.TagName, tag.TagColor, tag.CreatedBy)
	return err
}

func (r *RiskRepo) GetTags(ctx context.Context, riskID string) ([]RiskTag, error) {
	rows, err := r.db.Query(`
		SELECT id, risk_id, tag_name, tag_color, created_by, created_at
		FROM risk_tags WHERE risk_id = $1 ORDER BY created_at ASC
	`, riskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []RiskTag
	for rows.Next() {
		var tag RiskTag
		err := rows.Scan(&tag.ID, &tag.RiskID, &tag.TagName, &tag.TagColor, &tag.CreatedBy, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *RiskRepo) DeleteTag(ctx context.Context, tagID string) error {
	_, err := r.db.Exec("DELETE FROM risk_tags WHERE id = $1", tagID)
	return err
}

func (r *RiskRepo) DeleteTagByName(ctx context.Context, riskID, tagName string) error {
	_, err := r.db.Exec("DELETE FROM risk_tags WHERE risk_id = $1 AND tag_name = $2", riskID, tagName)
	return err
}
