package repo

import (
	"context"
	"database/sql"
	"strings"
)

type WorkflowRepo struct {
	db *DB
}

func NewWorkflowRepo(db *DB) *WorkflowRepo {
	return &WorkflowRepo{db: db}
}

// CreateWorkflow creates a new approval workflow
func (r *WorkflowRepo) CreateWorkflow(ctx context.Context, workflow ApprovalWorkflow) error {
	query := `
		INSERT INTO approval_workflows (id, document_id, workflow_type, status, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.DB.ExecContext(ctx, query,
		workflow.ID, workflow.DocumentID, workflow.WorkflowType, workflow.Status,
		workflow.CreatedBy, workflow.CreatedAt,
	)
	if err != nil {
		return err
	}

	// Note: Approval steps are created separately via CreateApprovalStep

	return nil
}

// CreateApprovalStep creates a new approval step
func (r *WorkflowRepo) CreateApprovalStep(ctx context.Context, step ApprovalStep) error {
	query := `
		INSERT INTO approval_steps (id, workflow_id, step_order, approver_id, status, deadline, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.DB.ExecContext(ctx, query,
		step.ID, step.WorkflowID, step.StepOrder, step.ApproverID, step.Status,
		step.Deadline, step.CreatedAt,
	)
	return err
}

// GetWorkflowByDocumentID retrieves workflow by document ID
func (r *WorkflowRepo) GetWorkflowByDocumentID(ctx context.Context, documentID string) (*ApprovalWorkflow, error) {
	// Get workflow
	workflowQuery := `
		SELECT id, document_id, workflow_type, status, created_by, created_at, completed_at
		FROM approval_workflows
		WHERE document_id = $1`

	var workflow ApprovalWorkflow
	err := r.db.DB.QueryRowContext(ctx, workflowQuery, documentID).Scan(
		&workflow.ID, &workflow.DocumentID, &workflow.WorkflowType, &workflow.Status,
		&workflow.CreatedBy, &workflow.CreatedAt, &workflow.CompletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get steps
	stepsQuery := `
		SELECT id, workflow_id, step_order, approver_id, status, comments, deadline, completed_at, created_at
		FROM approval_steps
		WHERE workflow_id = $1
		ORDER BY step_order`

	rows, err := r.db.DB.QueryContext(ctx, stepsQuery, workflow.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []ApprovalStep
	for rows.Next() {
		var step ApprovalStep
		err := rows.Scan(
			&step.ID, &step.WorkflowID, &step.StepOrder, &step.ApproverID,
			&step.Status, &step.Comments, &step.Deadline, &step.CompletedAt, &step.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}

	// Note: Steps are retrieved separately via ListApprovalSteps
	return &workflow, nil
}

// GetApprovalStep retrieves a specific approval step
func (r *WorkflowRepo) GetApprovalStep(ctx context.Context, stepID, tenantID string) (*ApprovalStep, error) {
	query := `
		SELECT s.id, s.workflow_id, s.step_order, s.approver_id, s.status, 
		       s.comments, s.deadline, s.completed_at, s.created_at
		FROM approval_steps s
		JOIN approval_workflows w ON s.workflow_id = w.id
		WHERE s.id = $1 AND w.tenant_id = $2`

	var step ApprovalStep
	err := r.db.DB.QueryRowContext(ctx, query, stepID, tenantID).Scan(
		&step.ID, &step.WorkflowID, &step.StepOrder, &step.ApproverID, &step.Status,
		&step.Comments, &step.Deadline, &step.CompletedAt, &step.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &step, nil
}

// ListApprovalSteps retrieves all steps for a workflow
func (r *WorkflowRepo) ListApprovalSteps(ctx context.Context, workflowID, tenantID string) ([]ApprovalStep, error) {
	query := `
		SELECT s.id, s.workflow_id, s.step_order, s.approver_id, s.status, 
		       s.comments, s.deadline, s.completed_at, s.created_at
		FROM approval_steps s
		JOIN approval_workflows w ON s.workflow_id = w.id
		WHERE s.workflow_id = $1 AND w.tenant_id = $2
		ORDER BY s.step_order`

	rows, err := r.db.DB.QueryContext(ctx, query, workflowID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []ApprovalStep
	for rows.Next() {
		var step ApprovalStep
		err := rows.Scan(
			&step.ID, &step.WorkflowID, &step.StepOrder, &step.ApproverID, &step.Status,
			&step.Comments, &step.Deadline, &step.CompletedAt, &step.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}

	return steps, nil
}

// UpdateWorkflow updates an approval workflow
func (r *WorkflowRepo) UpdateWorkflow(ctx context.Context, workflow ApprovalWorkflow) error {
	query := `
		UPDATE approval_workflows
		SET status = $1, completed_at = $2
		WHERE id = $3`

	_, err := r.db.DB.ExecContext(ctx, query,
		workflow.Status, workflow.CompletedAt, workflow.ID,
	)
	return err
}

// UpdateApprovalStep updates an approval step
func (r *WorkflowRepo) UpdateApprovalStep(ctx context.Context, step ApprovalStep) error {
	query := `
		UPDATE approval_steps
		SET status = $1, comments = $2, completed_at = $3
		WHERE id = $4`

	_, err := r.db.DB.ExecContext(ctx, query,
		step.Status, step.Comments, step.CompletedAt, step.ID,
	)
	return err
}

// CreateACKCampaign creates an acknowledgment campaign
func (r *WorkflowRepo) CreateACKCampaign(ctx context.Context, campaign ACKCampaign) error {
	query := `
		INSERT INTO ack_campaigns (id, document_id, title, description, audience_type, audience_ids, deadline, quiz_id, status, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	// Convert audience IDs to PostgreSQL array
	var audienceIDsArray interface{}
	if len(campaign.AudienceIDs) > 0 {
		audienceIDsArray = "{" + strings.Join(campaign.AudienceIDs, ",") + "}"
	} else {
		audienceIDsArray = "{}"
	}

	_, err := r.db.DB.ExecContext(ctx, query,
		campaign.ID, campaign.DocumentID, campaign.Title, campaign.Description,
		campaign.AudienceType, audienceIDsArray, campaign.Deadline, campaign.QuizID,
		campaign.Status, campaign.CreatedBy, campaign.CreatedAt,
	)
	return err
}

// CreateACKAssignment creates an acknowledgment assignment
func (r *WorkflowRepo) CreateACKAssignment(ctx context.Context, assignment ACKAssignment) error {
	query := `
		INSERT INTO ack_assignments (id, campaign_id, user_id, status, quiz_score, quiz_passed, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.DB.ExecContext(ctx, query,
		assignment.ID, assignment.CampaignID, assignment.UserID, assignment.Status,
		assignment.QuizScore, assignment.QuizPassed, assignment.CompletedAt, assignment.CreatedAt,
	)
	return err
}

// GetACKCampaignsByDocumentID retrieves ACK campaigns for a document
func (r *WorkflowRepo) GetACKCampaignsByDocumentID(ctx context.Context, documentID string) ([]ACKCampaign, error) {
	query := `
		SELECT id, document_id, title, description, audience_type, audience_ids, deadline, quiz_id, status, created_by, created_at, completed_at
		FROM ack_campaigns
		WHERE document_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.DB.QueryContext(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []ACKCampaign
	for rows.Next() {
		var campaign ACKCampaign
		var audienceIDsStr sql.NullString

		err := rows.Scan(
			&campaign.ID, &campaign.DocumentID, &campaign.Title, &campaign.Description,
			&campaign.AudienceType, &audienceIDsStr, &campaign.Deadline, &campaign.QuizID,
			&campaign.Status, &campaign.CreatedBy, &campaign.CreatedAt, &campaign.CompletedAt,
		)
		if err != nil {
			return nil, err
		}

		if audienceIDsStr.Valid && audienceIDsStr.String != "" {
			campaign.AudienceIDs = strings.Split(audienceIDsStr.String, ",")
		}

		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

// GetACKAssignmentsByCampaignID retrieves ACK assignments for a campaign
func (r *WorkflowRepo) GetACKAssignmentsByCampaignID(ctx context.Context, campaignID string) ([]ACKAssignment, error) {
	query := `
		SELECT id, campaign_id, user_id, status, quiz_score, quiz_passed, completed_at, created_at
		FROM ack_assignments
		WHERE campaign_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.DB.QueryContext(ctx, query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []ACKAssignment
	for rows.Next() {
		var assignment ACKAssignment
		err := rows.Scan(
			&assignment.ID, &assignment.CampaignID, &assignment.UserID, &assignment.Status,
			&assignment.QuizScore, &assignment.QuizPassed, &assignment.CompletedAt, &assignment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}
