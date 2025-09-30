package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Document represents a document in the system
type Document struct {
	ID                 string   `json:"id"`
	TenantID           string   `json:"tenant_id"`
	Title              string   `json:"title"`
	Code               *string  `json:"code"`
	Description        *string  `json:"description"`
	Type               string   `json:"type"`
	Category           *string  `json:"category"`
	Tags               []string `json:"tags"`
	Status             string   `json:"status"`
	CurrentVersion     int      `json:"current_version"`
	OwnerID            *string  `json:"owner_id"`
	Classification     string   `json:"classification"`
	EffectiveFrom      *string  `json:"effective_from"`
	ReviewPeriodMonths int      `json:"review_period_months"`
	AssetIDs           []string `json:"asset_ids"`
	RiskIDs            []string `json:"risk_ids"`
	ControlIDs         []string `json:"control_ids"`
	StorageKey         *string  `json:"storage_key"`
	MimeType           *string  `json:"mime_type"`
	SizeBytes          *int64   `json:"size_bytes"`
	ChecksumSHA256     *string  `json:"checksum_sha256"`
	OCRText            *string  `json:"ocr_text"`
	AVScanStatus       string   `json:"av_scan_status"`
	AVScanResult       *string  `json:"av_scan_result"`
	CreatedBy          string   `json:"created_by"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
	DeletedAt          *string  `json:"deleted_at,omitempty"`
}

// DocumentVersion represents a version of a document
type DocumentVersion struct {
	ID             string  `json:"id"`
	DocumentID     string  `json:"document_id"`
	VersionNumber  int     `json:"version_number"`
	StorageKey     string  `json:"storage_key"`
	MimeType       *string `json:"mime_type"`
	SizeBytes      *int64  `json:"size_bytes"`
	ChecksumSHA256 *string `json:"checksum_sha256"`
	OCRText        *string `json:"ocr_text"`
	AVScanStatus   string  `json:"av_scan_status"`
	AVScanResult   *string `json:"av_scan_result"`
	CreatedBy      string  `json:"created_by"`
	CreatedAt      string  `json:"created_at"`
	DeletedAt      *string `json:"deleted_at,omitempty"`
}

// ApprovalWorkflow represents a document approval workflow
type ApprovalWorkflow struct {
	ID           string  `json:"id"`
	DocumentID   string  `json:"document_id"`
	WorkflowType string  `json:"workflow_type"` // sequential or parallel
	Status       string  `json:"status"`
	CreatedBy    string  `json:"created_by"`
	CreatedAt    string  `json:"created_at"`
	CompletedAt  *string `json:"completed_at,omitempty"`
}

// ApprovalStep represents a step in an approval workflow
type ApprovalStep struct {
	ID          string  `json:"id"`
	WorkflowID  string  `json:"workflow_id"`
	StepOrder   int     `json:"step_order"`
	ApproverID  string  `json:"approver_id"`
	Status      string  `json:"status"`
	Comments    *string `json:"comments"`
	Deadline    *string `json:"deadline"`
	CompletedAt *string `json:"completed_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// ACKCampaign represents an acknowledgment campaign
type ACKCampaign struct {
	ID           string   `json:"id"`
	DocumentID   string   `json:"document_id"`
	Title        string   `json:"title"`
	Description  *string  `json:"description"`
	AudienceType string   `json:"audience_type"`
	AudienceIDs  []string `json:"audience_ids"`
	Deadline     *string  `json:"deadline"`
	QuizID       *string  `json:"quiz_id"`
	Status       string   `json:"status"`
	CreatedBy    string   `json:"created_by"`
	CreatedAt    string   `json:"created_at"`
	CompletedAt  *string  `json:"completed_at,omitempty"`
}

// ACKAssignment represents a user assignment for acknowledgment
type ACKAssignment struct {
	ID          string  `json:"id"`
	CampaignID  string  `json:"campaign_id"`
	UserID      string  `json:"user_id"`
	Status      string  `json:"status"`
	QuizScore   *int    `json:"quiz_score"`
	QuizPassed  bool    `json:"quiz_passed"`
	CompletedAt *string `json:"completed_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// TrainingMaterial represents a training material
type TrainingMaterial struct {
	ID             string  `json:"id"`
	TenantID       string  `json:"tenant_id"`
	Title          string  `json:"title"`
	Description    *string `json:"description"`
	Type           string  `json:"type"`
	StorageKey     string  `json:"storage_key"`
	MimeType       *string `json:"mime_type"`
	SizeBytes      *int64  `json:"size_bytes"`
	ChecksumSHA256 *string `json:"checksum_sha256"`
	CreatedBy      string  `json:"created_by"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at,omitempty"`
}

// Quiz represents a quiz
type Quiz struct {
	ID               string  `json:"id"`
	TenantID         string  `json:"tenant_id"`
	Title            string  `json:"title"`
	Description      *string `json:"description"`
	Questions        string  `json:"questions"` // JSONB
	PassingScore     int     `json:"passing_score"`
	TimeLimitMinutes *int    `json:"time_limit_minutes"`
	CreatedBy        string  `json:"created_by"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	DeletedAt        *string `json:"deleted_at,omitempty"`
}

// DocumentApprovalRoute represents an approval route for a document
type DocumentApprovalRoute struct {
	ID         string  `json:"id"`
	DocumentID string  `json:"document_id"`
	VersionID  *string `json:"version_id"`
	RouteName  string  `json:"route_name"`
	IsActive   bool    `json:"is_active"`
	CreatedBy  string  `json:"created_by"`
	CreatedAt  string  `json:"created_at"`
}

// DocumentApprovalStep represents a step in an approval route
type DocumentApprovalStep struct {
	ID             string  `json:"id"`
	RouteID        string  `json:"route_id"`
	StepOrder      int     `json:"step_order"`
	ApproverRoleID *string `json:"approver_role_id"`
	ApproverUserID *string `json:"approver_user_id"`
	IsRequired     bool    `json:"is_required"`
	CreatedAt      string  `json:"created_at"`
}

// DocumentApproval represents an approval action
type DocumentApproval struct {
	ID         string  `json:"id"`
	VersionID  string  `json:"version_id"`
	StepID     string  `json:"step_id"`
	ApproverID string  `json:"approver_id"`
	Status     string  `json:"status"`
	Comment    *string `json:"comment"`
	ApprovedAt *string `json:"approved_at"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// DocumentAcknowledgment represents user acknowledgment of a document
type DocumentAcknowledgment struct {
	ID             string  `json:"id"`
	DocumentID     string  `json:"document_id"`
	VersionID      *string `json:"version_id"`
	UserID         string  `json:"user_id"`
	Status         string  `json:"status"`
	QuizScore      *int    `json:"quiz_score"`
	QuizPassed     bool    `json:"quiz_passed"`
	AcknowledgedAt *string `json:"acknowledged_at"`
	Deadline       *string `json:"deadline"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// DocumentQuiz represents a quiz question for document acknowledgment
type DocumentQuiz struct {
	ID            string  `json:"id"`
	DocumentID    string  `json:"document_id"`
	VersionID     *string `json:"version_id"`
	Question      string  `json:"question"`
	QuestionOrder int     `json:"question_order"`
	Options       *string `json:"options"` // JSON string
	CorrectAnswer *string `json:"correct_answer"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
}

// DocumentQuizAnswer represents an answer to a quiz question
type DocumentQuizAnswer struct {
	ID               string `json:"id"`
	QuizID           string `json:"quiz_id"`
	AcknowledgmentID string `json:"acknowledgment_id"`
	UserID           string `json:"user_id"`
	Answer           string `json:"answer"`
	IsCorrect        *bool  `json:"is_correct"`
	AnsweredAt       string `json:"answered_at"`
}

// DocumentRepo handles database operations for documents
type DocumentRepo struct {
	db *DB
}

// NewDocumentRepo creates a new document repository
func NewDocumentRepo(db *DB) *DocumentRepo {
	return &DocumentRepo{db: db}
}

// ListDocuments retrieves documents for a tenant with optional filtering
func (r *DocumentRepo) ListDocuments(ctx context.Context, tenantID string, filters map[string]interface{}) ([]Document, error) {
	query := `
		SELECT id, tenant_id, title, code, description, type, category, 
		       array_to_string(tags, ',') as tags_str, status, 1 as current_version, 
		       owner_id, classification, effective_from, review_period_months,
		       asset_ids, risk_ids, control_ids, av_scan_status,
		       created_by, created_at, updated_at, deleted_at
		FROM documents 
		WHERE tenant_id = $1 AND deleted_at IS NULL`

	args := []interface{}{tenantID}
	argIndex := 2

	// Apply filters
	if status, ok := filters["status"]; ok && status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if docType, ok := filters["type"]; ok && docType != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, docType)
		argIndex++
	}

	if category, ok := filters["category"]; ok && category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if search, ok := filters["search"]; ok && search != "" {
		query += fmt.Sprintf(" AND to_tsvector('russian', title || ' ' || COALESCE(description, '')) @@ plainto_tsquery('russian', $%d)", argIndex)
		args = append(args, search)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var doc Document
		var tagsStr sql.NullString
		var assetIDsStr, riskIDsStr, controlIDsStr sql.NullString

		err := rows.Scan(
			&doc.ID, &doc.TenantID, &doc.Title, &doc.Code, &doc.Description, &doc.Type,
			&doc.Category, &tagsStr, &doc.Status, &doc.CurrentVersion,
			&doc.OwnerID, &doc.Classification, &doc.EffectiveFrom, &doc.ReviewPeriodMonths,
			&assetIDsStr, &riskIDsStr, &controlIDsStr, &doc.AVScanStatus,
			&doc.CreatedBy, &doc.CreatedAt, &doc.UpdatedAt, &doc.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		if tagsStr.Valid && tagsStr.String != "" {
			doc.Tags = strings.Split(tagsStr.String, ",")
		}

		if assetIDsStr.Valid && assetIDsStr.String != "" {
			doc.AssetIDs = strings.Split(assetIDsStr.String, ",")
		}

		if riskIDsStr.Valid && riskIDsStr.String != "" {
			doc.RiskIDs = strings.Split(riskIDsStr.String, ",")
		}

		if controlIDsStr.Valid && controlIDsStr.String != "" {
			doc.ControlIDs = strings.Split(controlIDsStr.String, ",")
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

// GetDocument retrieves a document by ID
func (r *DocumentRepo) GetDocument(ctx context.Context, id, tenantID string) (*Document, error) {
	query := `
		SELECT id, tenant_id, title, code, description, type, category, 
		       array_to_string(tags, ',') as tags_str, status, 1 as current_version, 
		       owner_id, classification, effective_from, review_period_months,
		       asset_ids, risk_ids, control_ids, av_scan_status,
		       created_by, created_at, updated_at, deleted_at
		FROM documents 
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`

	var doc Document
	var tagsStr sql.NullString
	var assetIDsStr, riskIDsStr, controlIDsStr sql.NullString

	err := r.db.DB.QueryRowContext(ctx, query, id, tenantID).Scan(
		&doc.ID, &doc.TenantID, &doc.Title, &doc.Code, &doc.Description, &doc.Type,
		&doc.Category, &tagsStr, &doc.Status, &doc.CurrentVersion,
		&doc.OwnerID, &doc.Classification, &doc.EffectiveFrom, &doc.ReviewPeriodMonths,
		&assetIDsStr, &riskIDsStr, &controlIDsStr, &doc.AVScanStatus,
		&doc.CreatedBy, &doc.CreatedAt, &doc.UpdatedAt, &doc.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if tagsStr.Valid && tagsStr.String != "" {
		doc.Tags = strings.Split(tagsStr.String, ",")
	}

	if assetIDsStr.Valid && assetIDsStr.String != "" {
		doc.AssetIDs = strings.Split(assetIDsStr.String, ",")
	}

	if riskIDsStr.Valid && riskIDsStr.String != "" {
		doc.RiskIDs = strings.Split(riskIDsStr.String, ",")
	}

	if controlIDsStr.Valid && controlIDsStr.String != "" {
		doc.ControlIDs = strings.Split(controlIDsStr.String, ",")
	}

	return &doc, nil
}

// CreateDocument creates a new document
func (r *DocumentRepo) CreateDocument(ctx context.Context, doc Document) error {
	fmt.Printf("DEBUG: CreateDocument repo called with: %+v\n", doc)
	query := `
		INSERT INTO documents (id, tenant_id, title, code, description, type, category, tags, 
		                      status, version, owner_id, classification, effective_from, 
		                      review_period_months, asset_ids, risk_ids, control_ids, 
		                      av_scan_status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	fmt.Printf("DEBUG: Executing query: %s\n", query)

	// Convert slices to PostgreSQL arrays
	var tagsArray interface{}
	if len(doc.Tags) > 0 {
		tagsArray = "{" + strings.Join(doc.Tags, ",") + "}"
	} else {
		tagsArray = "{}"
	}

	var assetIDsArray interface{}
	if len(doc.AssetIDs) > 0 {
		assetIDsArray = "{" + strings.Join(doc.AssetIDs, ",") + "}"
	} else {
		assetIDsArray = "{}"
	}

	var riskIDsArray interface{}
	if len(doc.RiskIDs) > 0 {
		riskIDsArray = "{" + strings.Join(doc.RiskIDs, ",") + "}"
	} else {
		riskIDsArray = "{}"
	}

	var controlIDsArray interface{}
	if len(doc.ControlIDs) > 0 {
		controlIDsArray = "{" + strings.Join(doc.ControlIDs, ",") + "}"
	} else {
		controlIDsArray = "{}"
	}

	_, err := r.db.DB.ExecContext(ctx, query,
		doc.ID, doc.TenantID, doc.Title, doc.Code, doc.Description, doc.Type, doc.Category,
		tagsArray, doc.Status, "1.0", doc.OwnerID, doc.Classification, doc.EffectiveFrom,
		doc.ReviewPeriodMonths, assetIDsArray, riskIDsArray, controlIDsArray,
		doc.AVScanStatus, doc.CreatedBy,
	)
	if err != nil {
		fmt.Printf("DEBUG: Database error: %v\n", err)
		return err
	}
	fmt.Printf("DEBUG: Document created successfully in database\n")
	return nil
}

// UpdateDocument updates an existing document
func (r *DocumentRepo) UpdateDocument(ctx context.Context, doc Document) error {
	query := `
		UPDATE documents 
		SET title = $1, description = $2, type = $3, category = $4, 
		    status = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6 AND tenant_id = $7`

	_, err := r.db.DB.ExecContext(ctx, query,
		doc.Title, doc.Description, doc.Type, doc.Category,
		doc.Status, doc.ID, doc.TenantID,
	)
	return err
}

// DeleteDocument soft deletes a document
func (r *DocumentRepo) DeleteDocument(ctx context.Context, id, tenantID string) error {
	query := `UPDATE documents SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.DB.ExecContext(ctx, query, id, tenantID)
	return err
}

// ListDocumentVersions retrieves versions for a document
func (r *DocumentRepo) ListDocumentVersions(ctx context.Context, documentID, tenantID string) ([]DocumentVersion, error) {
	query := `
		SELECT dv.id, dv.document_id, dv.version_number, dv.storage_key, dv.mime_type,
		       dv.size_bytes, dv.checksum_sha256, dv.created_by, dv.created_at
		FROM document_versions dv
		JOIN documents d ON dv.document_id = d.id
		WHERE dv.document_id = $1 AND d.tenant_id = $2
		ORDER BY dv.version_number DESC`

	rows, err := r.db.DB.QueryContext(ctx, query, documentID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []DocumentVersion
	for rows.Next() {
		var version DocumentVersion
		err := rows.Scan(
			&version.ID, &version.DocumentID, &version.VersionNumber, &version.StorageKey,
			&version.MimeType, &version.SizeBytes, &version.ChecksumSHA256, &version.CreatedBy, &version.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}

// GetDocumentVersion retrieves a specific version of a document
func (r *DocumentRepo) GetDocumentVersion(ctx context.Context, versionID, tenantID string) (*DocumentVersion, error) {
	query := `
		SELECT dv.id, dv.document_id, dv.version_number, dv.storage_key, dv.mime_type,
		       dv.size_bytes, dv.checksum_sha256, dv.created_by, dv.created_at
		FROM document_versions dv
		JOIN documents d ON dv.document_id = d.id
		WHERE dv.id = $1 AND d.tenant_id = $2`

	var version DocumentVersion
	err := r.db.DB.QueryRowContext(ctx, query, versionID, tenantID).Scan(
		&version.ID, &version.DocumentID, &version.VersionNumber, &version.StorageKey,
		&version.MimeType, &version.SizeBytes, &version.ChecksumSHA256, &version.CreatedBy, &version.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &version, nil
}

// CreateDocumentVersion creates a new version of a document
func (r *DocumentRepo) CreateDocumentVersion(ctx context.Context, version DocumentVersion) error {
	query := `
		INSERT INTO document_versions (id, document_id, version_number, storage_key, 
		                             size_bytes, mime_type, checksum_sha256, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.DB.ExecContext(ctx, query,
		version.ID, version.DocumentID, version.VersionNumber, version.StorageKey,
		version.SizeBytes, version.MimeType, version.ChecksumSHA256, version.CreatedBy,
	)
	return err
}

// ListDocumentAcknowledgment retrieves acknowledgments for a document
func (r *DocumentRepo) ListDocumentAcknowledgment(ctx context.Context, documentID, tenantID string) ([]DocumentAcknowledgment, error) {
	query := `
		SELECT da.id, da.document_id, da.version_id, da.user_id, da.status, da.quiz_score,
		       da.quiz_passed, da.acknowledged_at, da.deadline, da.created_at, da.updated_at
		FROM document_acknowledgments da
		JOIN documents d ON da.document_id = d.id
		WHERE da.document_id = $1 AND d.tenant_id = $2
		ORDER BY da.created_at DESC`

	rows, err := r.db.DB.QueryContext(ctx, query, documentID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acknowledgments []DocumentAcknowledgment
	for rows.Next() {
		var ack DocumentAcknowledgment
		err := rows.Scan(
			&ack.ID, &ack.DocumentID, &ack.VersionID, &ack.UserID, &ack.Status,
			&ack.QuizScore, &ack.QuizPassed, &ack.AcknowledgedAt, &ack.Deadline,
			&ack.CreatedAt, &ack.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		acknowledgments = append(acknowledgments, ack)
	}

	return acknowledgments, nil
}

// CreateDocumentAcknowledgment creates a new acknowledgment
func (r *DocumentRepo) CreateDocumentAcknowledgment(ctx context.Context, ack DocumentAcknowledgment) error {
	query := `
		INSERT INTO document_acknowledgments (id, document_id, version_id, user_id, status, deadline)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.DB.ExecContext(ctx, query,
		ack.ID, ack.DocumentID, ack.VersionID, ack.UserID, ack.Status, ack.Deadline,
	)
	return err
}

// UpdateDocumentAcknowledgment updates an acknowledgment
func (r *DocumentRepo) UpdateDocumentAcknowledgment(ctx context.Context, ack DocumentAcknowledgment) error {
	query := `
		UPDATE document_acknowledgments 
		SET status = $1, quiz_score = $2, quiz_passed = $3, acknowledged_at = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5`

	_, err := r.db.DB.ExecContext(ctx, query,
		ack.Status, ack.QuizScore, ack.QuizPassed, ack.AcknowledgedAt, ack.ID,
	)
	return err
}

// ListDocumentQuizzes retrieves quizzes for a document
func (r *DocumentRepo) ListDocumentQuizzes(ctx context.Context, documentID, tenantID string) ([]DocumentQuiz, error) {
	query := `
		SELECT dq.id, dq.document_id, dq.version_id, dq.question, dq.question_order,
		       dq.options, dq.correct_answer, dq.is_active, dq.created_at
		FROM document_quizzes dq
		JOIN documents d ON dq.document_id = d.id
		WHERE dq.document_id = $1 AND d.tenant_id = $2 AND dq.is_active = true
		ORDER BY dq.question_order`

	rows, err := r.db.DB.QueryContext(ctx, query, documentID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []DocumentQuiz
	for rows.Next() {
		var quiz DocumentQuiz
		err := rows.Scan(
			&quiz.ID, &quiz.DocumentID, &quiz.VersionID, &quiz.Question,
			&quiz.QuestionOrder, &quiz.Options, &quiz.CorrectAnswer, &quiz.IsActive, &quiz.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

// CreateDocumentQuiz creates a new quiz question
func (r *DocumentRepo) CreateDocumentQuiz(ctx context.Context, quiz DocumentQuiz) error {
	query := `
		INSERT INTO document_quizzes (id, document_id, version_id, question, question_order, options, correct_answer, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.DB.ExecContext(ctx, query,
		quiz.ID, quiz.DocumentID, quiz.VersionID, quiz.Question, quiz.QuestionOrder,
		quiz.Options, quiz.CorrectAnswer, quiz.IsActive,
	)
	return err
}

// GetUserPendingAcknowledgment retrieves pending acknowledgments for a user
func (r *DocumentRepo) GetUserPendingAcknowledgment(ctx context.Context, userID, tenantID string) ([]DocumentAcknowledgment, error) {
	query := `
		SELECT da.id, da.document_id, da.version_id, da.user_id, da.status, da.quiz_score,
		       da.quiz_passed, da.acknowledged_at, da.deadline, da.created_at, da.updated_at
		FROM document_acknowledgments da
		JOIN documents d ON da.document_id = d.id
		WHERE da.user_id = $1 AND d.tenant_id = $2 AND da.status = 'pending'
		ORDER BY da.deadline ASC NULLS LAST, da.created_at ASC`

	rows, err := r.db.DB.QueryContext(ctx, query, userID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acknowledgments []DocumentAcknowledgment
	for rows.Next() {
		var ack DocumentAcknowledgment
		err := rows.Scan(
			&ack.ID, &ack.DocumentID, &ack.VersionID, &ack.UserID, &ack.Status,
			&ack.QuizScore, &ack.QuizPassed, &ack.AcknowledgedAt, &ack.Deadline,
			&ack.CreatedAt, &ack.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		acknowledgments = append(acknowledgments, ack)
	}

	return acknowledgments, nil
}
