package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Incident struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Category    string     `json:"category"`
	Status      string     `json:"status"`
	Severity    string     `json:"severity"`
	Criticality string     `json:"criticality"`
	Source      string     `json:"source"`
	ReportedBy  string     `json:"reported_by"`
	AssignedTo  *string    `json:"assigned_to"`
	AssetID     *string    `json:"asset_id"`
	RiskID      *string    `json:"risk_id"`
	CreatedBy   string     `json:"created_by"`
	DetectedAt  time.Time  `json:"detected_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	ClosedAt    *time.Time `json:"closed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type IncidentAsset struct {
	ID         string    `json:"id"`
	IncidentID string    `json:"incident_id"`
	AssetID    string    `json:"asset_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type IncidentRisk struct {
	ID         string    `json:"id"`
	IncidentID string    `json:"incident_id"`
	RiskID     string    `json:"risk_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type IncidentComment struct {
	ID         string    `json:"id"`
	IncidentID string    `json:"incident_id"`
	UserID     string    `json:"user_id"`
	Comment    string    `json:"comment"`
	IsInternal bool      `json:"is_internal"`
	CreatedAt  time.Time `json:"created_at"`
}

type IncidentAttachment struct {
	ID         string    `json:"id"`
	IncidentID string    `json:"incident_id"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `json:"mime_type"`
	UploadedBy string    `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type IncidentAction struct {
	ID          string     `json:"id"`
	IncidentID  string     `json:"incident_id"`
	ActionType  string     `json:"action_type"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	AssignedTo  *string    `json:"assigned_to"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	Status      string     `json:"status"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type IncidentMetrics struct {
	ID           string    `json:"id"`
	IncidentID   string    `json:"incident_id"`
	MetricType   string    `json:"metric_type"`
	ValueMinutes int       `json:"value_minutes"`
	CalculatedAt time.Time `json:"calculated_at"`
}

type IncidentRepository interface {
	Create(ctx context.Context, incident *Incident) error
	GetByID(ctx context.Context, id, tenantID string) (*Incident, error)
	Update(ctx context.Context, incident *Incident) error
	Delete(ctx context.Context, id, tenantID string) error
	List(ctx context.Context, tenantID string, filters map[string]interface{}, limit, offset int) ([]*Incident, int, error)

	// Asset relations
	AddAsset(ctx context.Context, incidentID, assetID string) error
	RemoveAsset(ctx context.Context, incidentID, assetID string) error
	GetAssets(ctx context.Context, incidentID string) ([]*Asset, error)

	// Risk relations
	AddRisk(ctx context.Context, incidentID, riskID string) error
	RemoveRisk(ctx context.Context, incidentID, riskID string) error
	GetRisks(ctx context.Context, incidentID string) ([]*Risk, error)

	// Comments
	AddComment(ctx context.Context, comment *IncidentComment) error
	GetComments(ctx context.Context, incidentID string) ([]*IncidentComment, error)

	// Attachments
	AddAttachment(ctx context.Context, attachment *IncidentAttachment) error
	GetAttachments(ctx context.Context, incidentID string) ([]*IncidentAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID string) error

	// Actions
	AddAction(ctx context.Context, action *IncidentAction) error
	UpdateAction(ctx context.Context, action *IncidentAction) error
	GetActions(ctx context.Context, incidentID string) ([]*IncidentAction, error)
	DeleteAction(ctx context.Context, actionID string) error

	// Metrics
	AddMetric(ctx context.Context, metric *IncidentMetrics) error
	GetMetrics(ctx context.Context, incidentID string) ([]*IncidentMetrics, error)
	GetIncidentMetrics(ctx context.Context, tenantID string) (*IncidentMetricsSummary, error)
}

type IncidentMetricsSummary struct {
	TotalIncidents  int            `json:"total_incidents"`
	OpenIncidents   int            `json:"open_incidents"`
	ClosedIncidents int            `json:"closed_incidents"`
	AverageMTTR     float64        `json:"average_mttr_hours"`
	AverageMTTD     float64        `json:"average_mttd_hours"`
	ByCriticality   map[string]int `json:"by_criticality"`
	ByCategory      map[string]int `json:"by_category"`
	ByStatus        map[string]int `json:"by_status"`
}

type incidentRepository struct {
	db *sql.DB
}

func NewIncidentRepository(db *sql.DB) IncidentRepository {
	return &incidentRepository{db: db}
}

func (r *incidentRepository) Create(ctx context.Context, incident *Incident) error {
	query := `
		INSERT INTO incidents (id, tenant_id, title, description, category, status, criticality, source, reported_by, assigned_to, detected_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(ctx, query,
		incident.ID, incident.TenantID, incident.Title, incident.Description,
		incident.Category, incident.Status, incident.Criticality, incident.Source,
		incident.ReportedBy, incident.AssignedTo, incident.DetectedAt,
		incident.CreatedAt, incident.UpdatedAt)

	return err
}

func (r *incidentRepository) GetByID(ctx context.Context, id, tenantID string) (*Incident, error) {
	query := `
		SELECT id, tenant_id, title, description, category, status, criticality, source, 
		       reported_by, assigned_to, detected_at, resolved_at, closed_at, created_at, updated_at
		FROM incidents 
		WHERE id = $1 AND tenant_id = $2
	`

	var incident Incident
	err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&incident.ID, &incident.TenantID, &incident.Title, &incident.Description,
		&incident.Category, &incident.Status, &incident.Criticality, &incident.Source,
		&incident.ReportedBy, &incident.AssignedTo, &incident.DetectedAt,
		&incident.ResolvedAt, &incident.ClosedAt, &incident.CreatedAt, &incident.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("incident not found")
		}
		return nil, err
	}

	return &incident, nil
}

func (r *incidentRepository) Update(ctx context.Context, incident *Incident) error {
	query := `
		UPDATE incidents 
		SET title = $1, description = $2, category = $3, status = $4, criticality = $5, 
		    source = $6, assigned_to = $7, detected_at = $8, resolved_at = $9, closed_at = $10, 
		    updated_at = $11
		WHERE id = $12 AND tenant_id = $13
	`

	_, err := r.db.ExecContext(ctx, query,
		incident.Title, incident.Description, incident.Category, incident.Status,
		incident.Criticality, incident.Source, incident.AssignedTo, incident.DetectedAt,
		incident.ResolvedAt, incident.ClosedAt, incident.UpdatedAt,
		incident.ID, incident.TenantID)

	return err
}

func (r *incidentRepository) Delete(ctx context.Context, id, tenantID string) error {
	query := `DELETE FROM incidents WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, tenantID)
	return err
}

func (r *incidentRepository) List(ctx context.Context, tenantID string, filters map[string]interface{}, limit, offset int) ([]*Incident, int, error) {
	whereClause := "WHERE tenant_id = $1"
	args := []interface{}{tenantID}
	argIndex := 2

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if criticality, ok := filters["criticality"].(string); ok && criticality != "" {
		whereClause += fmt.Sprintf(" AND criticality = $%d", argIndex)
		args = append(args, criticality)
		argIndex++
	}

	if category, ok := filters["category"].(string); ok && category != "" {
		whereClause += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if assignedTo, ok := filters["assigned_to"].(string); ok && assignedTo != "" {
		whereClause += fmt.Sprintf(" AND assigned_to = $%d", argIndex)
		args = append(args, assignedTo)
		argIndex++
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		whereClause += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM incidents %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// List query
	listQuery := fmt.Sprintf(`
		SELECT id, tenant_id, title, description, category, status, criticality, source, 
		       reported_by, assigned_to, detected_at, resolved_at, closed_at, created_at, updated_at
		FROM incidents %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var incidents []*Incident
	for rows.Next() {
		var incident Incident
		err := rows.Scan(
			&incident.ID, &incident.TenantID, &incident.Title, &incident.Description,
			&incident.Category, &incident.Status, &incident.Criticality, &incident.Source,
			&incident.ReportedBy, &incident.AssignedTo, &incident.DetectedAt,
			&incident.ResolvedAt, &incident.ClosedAt, &incident.CreatedAt, &incident.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		incidents = append(incidents, &incident)
	}

	return incidents, total, nil
}

// Asset relations
func (r *incidentRepository) AddAsset(ctx context.Context, incidentID, assetID string) error {
	query := `INSERT INTO incident_assets (id, incident_id, asset_id) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, uuid.New().String(), incidentID, assetID)
	return err
}

func (r *incidentRepository) RemoveAsset(ctx context.Context, incidentID, assetID string) error {
	query := `DELETE FROM incident_assets WHERE incident_id = $1 AND asset_id = $2`
	_, err := r.db.ExecContext(ctx, query, incidentID, assetID)
	return err
}

func (r *incidentRepository) GetAssets(ctx context.Context, incidentID string) ([]*Asset, error) {
	query := `
		SELECT a.id, a.tenant_id, a.inventory_number, a.name, a.type, a.class, a.owner_id, 
		       a.location, a.criticality, a.confidentiality, a.integrity, a.availability, 
		       a.status, a.created_at, a.updated_at
		FROM assets a
		JOIN incident_assets ia ON a.id = ia.asset_id
		WHERE ia.incident_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*Asset
	for rows.Next() {
		var asset Asset
		err := rows.Scan(
			&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name, &asset.Type,
			&asset.Class, &asset.OwnerID, &asset.Location, &asset.Criticality,
			&asset.Confidentiality, &asset.Integrity, &asset.Availability,
			&asset.Status, &asset.CreatedAt, &asset.UpdatedAt)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// Risk relations
func (r *incidentRepository) AddRisk(ctx context.Context, incidentID, riskID string) error {
	query := `INSERT INTO incident_risks (id, incident_id, risk_id) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, uuid.New().String(), incidentID, riskID)
	return err
}

func (r *incidentRepository) RemoveRisk(ctx context.Context, incidentID, riskID string) error {
	query := `DELETE FROM incident_risks WHERE incident_id = $1 AND risk_id = $2`
	_, err := r.db.ExecContext(ctx, query, incidentID, riskID)
	return err
}

func (r *incidentRepository) GetRisks(ctx context.Context, incidentID string) ([]*Risk, error) {
	query := `
		SELECT r.id, r.tenant_id, r.title, r.description, r.category, r.likelihood, 
		       r.impact, r.level, r.status, r.owner_user_id, r.created_at, r.updated_at
		FROM risks r
		JOIN incident_risks ir ON r.id = ir.risk_id
		WHERE ir.incident_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var risks []*Risk
	for rows.Next() {
		var risk Risk
		err := rows.Scan(
			&risk.ID, &risk.TenantID, &risk.Title, &risk.Description, &risk.Category,
			&risk.Likelihood, &risk.Impact, &risk.Level, &risk.Status,
			&risk.OwnerUserID, &risk.CreatedAt, &risk.UpdatedAt)
		if err != nil {
			return nil, err
		}
		risks = append(risks, &risk)
	}

	return risks, nil
}

// Comments
func (r *incidentRepository) AddComment(ctx context.Context, comment *IncidentComment) error {
	query := `
		INSERT INTO incident_comments (id, incident_id, user_id, comment, is_internal, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		comment.ID, comment.IncidentID, comment.UserID, comment.Comment,
		comment.IsInternal, comment.CreatedAt)
	return err
}

func (r *incidentRepository) GetComments(ctx context.Context, incidentID string) ([]*IncidentComment, error) {
	query := `
		SELECT id, incident_id, user_id, comment, is_internal, created_at
		FROM incident_comments
		WHERE incident_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*IncidentComment
	for rows.Next() {
		var comment IncidentComment
		err := rows.Scan(
			&comment.ID, &comment.IncidentID, &comment.UserID, &comment.Comment,
			&comment.IsInternal, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

// Attachments
func (r *incidentRepository) AddAttachment(ctx context.Context, attachment *IncidentAttachment) error {
	query := `
		INSERT INTO incident_attachments (id, incident_id, file_name, file_path, file_size, mime_type, uploaded_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		attachment.ID, attachment.IncidentID, attachment.FileName, attachment.FilePath,
		attachment.FileSize, attachment.MimeType, attachment.UploadedBy, attachment.CreatedAt)
	return err
}

func (r *incidentRepository) GetAttachments(ctx context.Context, incidentID string) ([]*IncidentAttachment, error) {
	query := `
		SELECT id, incident_id, file_name, file_path, file_size, mime_type, uploaded_by, created_at
		FROM incident_attachments
		WHERE incident_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []*IncidentAttachment
	for rows.Next() {
		var attachment IncidentAttachment
		err := rows.Scan(
			&attachment.ID, &attachment.IncidentID, &attachment.FileName, &attachment.FilePath,
			&attachment.FileSize, &attachment.MimeType, &attachment.UploadedBy, &attachment.CreatedAt)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, &attachment)
	}

	return attachments, nil
}

func (r *incidentRepository) DeleteAttachment(ctx context.Context, attachmentID string) error {
	query := `DELETE FROM incident_attachments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, attachmentID)
	return err
}

// Actions
func (r *incidentRepository) AddAction(ctx context.Context, action *IncidentAction) error {
	query := `
		INSERT INTO incident_actions (id, incident_id, action_type, title, description, assigned_to, due_date, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		action.ID, action.IncidentID, action.ActionType, action.Title, action.Description,
		action.AssignedTo, action.DueDate, action.Status, action.CreatedBy,
		action.CreatedAt, action.UpdatedAt)
	return err
}

func (r *incidentRepository) UpdateAction(ctx context.Context, action *IncidentAction) error {
	query := `
		UPDATE incident_actions 
		SET action_type = $1, title = $2, description = $3, assigned_to = $4, due_date = $5, 
		    completed_at = $6, status = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.ExecContext(ctx, query,
		action.ActionType, action.Title, action.Description, action.AssignedTo,
		action.DueDate, action.CompletedAt, action.Status, action.UpdatedAt, action.ID)
	return err
}

func (r *incidentRepository) GetActions(ctx context.Context, incidentID string) ([]*IncidentAction, error) {
	query := `
		SELECT id, incident_id, action_type, title, description, assigned_to, due_date, 
		       completed_at, status, created_by, created_at, updated_at
		FROM incident_actions
		WHERE incident_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []*IncidentAction
	for rows.Next() {
		var action IncidentAction
		err := rows.Scan(
			&action.ID, &action.IncidentID, &action.ActionType, &action.Title,
			&action.Description, &action.AssignedTo, &action.DueDate, &action.CompletedAt,
			&action.Status, &action.CreatedBy, &action.CreatedAt, &action.UpdatedAt)
		if err != nil {
			return nil, err
		}
		actions = append(actions, &action)
	}

	return actions, nil
}

func (r *incidentRepository) DeleteAction(ctx context.Context, actionID string) error {
	query := `DELETE FROM incident_actions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, actionID)
	return err
}

// Metrics
func (r *incidentRepository) AddMetric(ctx context.Context, metric *IncidentMetrics) error {
	query := `
		INSERT INTO incident_metrics (id, incident_id, metric_type, value_minutes, calculated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		metric.ID, metric.IncidentID, metric.MetricType, metric.ValueMinutes, metric.CalculatedAt)
	return err
}

func (r *incidentRepository) GetMetrics(ctx context.Context, incidentID string) ([]*IncidentMetrics, error) {
	query := `
		SELECT id, incident_id, metric_type, value_minutes, calculated_at
		FROM incident_metrics
		WHERE incident_id = $1
		ORDER BY calculated_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*IncidentMetrics
	for rows.Next() {
		var metric IncidentMetrics
		err := rows.Scan(
			&metric.ID, &metric.IncidentID, &metric.MetricType, &metric.ValueMinutes, &metric.CalculatedAt)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

func (r *incidentRepository) GetIncidentMetrics(ctx context.Context, tenantID string) (*IncidentMetricsSummary, error) {
	// Get basic counts
	countQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status IN ('new', 'assigned', 'in_progress') THEN 1 END) as open,
			COUNT(CASE WHEN status = 'closed' THEN 1 END) as closed
		FROM incidents 
		WHERE tenant_id = $1
	`

	var total, open, closed int
	err := r.db.QueryRowContext(ctx, countQuery, tenantID).Scan(&total, &open, &closed)
	if err != nil {
		return nil, err
	}

	// Get metrics by criticality
	criticalityQuery := `
		SELECT criticality, COUNT(*) 
		FROM incidents 
		WHERE tenant_id = $1
		GROUP BY criticality
	`

	byCriticality := make(map[string]int)
	rows, err := r.db.QueryContext(ctx, criticalityQuery, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var criticality string
		var count int
		err := rows.Scan(&criticality, &count)
		if err != nil {
			return nil, err
		}
		byCriticality[criticality] = count
	}

	// Get metrics by category
	categoryQuery := `
		SELECT category, COUNT(*) 
		FROM incidents 
		WHERE tenant_id = $1
		GROUP BY category
	`

	byCategory := make(map[string]int)
	rows, err = r.db.QueryContext(ctx, categoryQuery, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var count int
		err := rows.Scan(&category, &count)
		if err != nil {
			return nil, err
		}
		byCategory[category] = count
	}

	// Get metrics by status
	statusQuery := `
		SELECT status, COUNT(*) 
		FROM incidents 
		WHERE tenant_id = $1
		GROUP BY status
	`

	byStatus := make(map[string]int)
	rows, err = r.db.QueryContext(ctx, statusQuery, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, err
		}
		byStatus[status] = count
	}

	// Calculate average MTTR and MTTD
	metricsQuery := `
		SELECT 
			AVG(CASE WHEN metric_type = 'mttr' THEN value_minutes END) as avg_mttr,
			AVG(CASE WHEN metric_type = 'mttd' THEN value_minutes END) as avg_mttd
		FROM incident_metrics im
		JOIN incidents i ON im.incident_id = i.id
		WHERE i.tenant_id = $1
	`

	var avgMTTR, avgMTTD sql.NullFloat64
	err = r.db.QueryRowContext(ctx, metricsQuery, tenantID).Scan(&avgMTTR, &avgMTTD)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	summary := &IncidentMetricsSummary{
		TotalIncidents:  total,
		OpenIncidents:   open,
		ClosedIncidents: closed,
		AverageMTTR:     0,
		AverageMTTD:     0,
		ByCriticality:   byCriticality,
		ByCategory:      byCategory,
		ByStatus:        byStatus,
	}

	if avgMTTR.Valid {
		summary.AverageMTTR = avgMTTR.Float64 / 60.0 // Convert to hours
	}
	if avgMTTD.Valid {
		summary.AverageMTTD = avgMTTD.Float64 / 60.0 // Convert to hours
	}

	return summary, nil
}
