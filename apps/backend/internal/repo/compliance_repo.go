package repo

import (
	"context"
)

type ComplianceStandard struct {
	ID          string
	TenantID    string
	Name        string
	Code        string
	Description *string
	Version     string
	IsActive    bool
	CreatedAt   string
	UpdatedAt   string
}

type ComplianceRequirement struct {
	ID          string
	StandardID  string
	Code        string
	Title       string
	Description *string
	Category    *string
	Priority    string
	IsMandatory bool
	CreatedAt   string
	UpdatedAt   string
}

type ComplianceAssessment struct {
	ID             string
	TenantID       string
	RequirementID  string
	Status         string
	Evidence       *string
	AssessorID     *string
	AssessedAt     *string
	NextReviewDate *string
	Notes          *string
	CreatedAt      string
	UpdatedAt      string
	// Joined fields
	RequirementTitle string
	StandardName     string
	AssessorName     *string
}

type ComplianceGap struct {
	ID              string
	AssessmentID    string
	Title           string
	Description     *string
	Severity        string
	Status          string
	RemediationPlan *string
	TargetDate      *string
	ResponsibleID   *string
	CreatedAt       string
	UpdatedAt       string
	// Joined fields
	ResponsibleName *string
}

type ComplianceRepo struct {
	db *DB
}

func NewComplianceRepo(db *DB) *ComplianceRepo {
	return &ComplianceRepo{db: db}
}

// Standards
func (r *ComplianceRepo) ListStandards(ctx context.Context, tenantID string) ([]ComplianceStandard, error) {
	rows, err := r.db.Query(`SELECT id, tenant_id, name, code, description, version, is_active, created_at, updated_at FROM compliance_standards WHERE tenant_id=$1 ORDER BY name`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ComplianceStandard
	for rows.Next() {
		var s ComplianceStandard
		if err := rows.Scan(&s.ID, &s.TenantID, &s.Name, &s.Code, &s.Description, &s.Version, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

func (r *ComplianceRepo) CreateStandard(ctx context.Context, s ComplianceStandard) error {
	_, err := r.db.Exec(`INSERT INTO compliance_standards(id,tenant_id,name,code,description,version,is_active) VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6)`, s.TenantID, s.Name, s.Code, s.Description, s.Version, s.IsActive)
	return err
}

// Requirements
func (r *ComplianceRepo) ListRequirements(ctx context.Context, standardID string) ([]ComplianceRequirement, error) {
	rows, err := r.db.Query(`SELECT id, standard_id, code, title, description, category, priority, is_mandatory, created_at, updated_at FROM compliance_requirements WHERE standard_id=$1 ORDER BY code`, standardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ComplianceRequirement
	for rows.Next() {
		var req ComplianceRequirement
		if err := rows.Scan(&req.ID, &req.StandardID, &req.Code, &req.Title, &req.Description, &req.Category, &req.Priority, &req.IsMandatory, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, req)
	}
	return items, nil
}

func (r *ComplianceRepo) CreateRequirement(ctx context.Context, req ComplianceRequirement) error {
	_, err := r.db.Exec(`INSERT INTO compliance_requirements(id,standard_id,code,title,description,category,priority,is_mandatory) VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7)`, req.StandardID, req.Code, req.Title, req.Description, req.Category, req.Priority, req.IsMandatory)
	return err
}

// Assessments
func (r *ComplianceRepo) ListAssessments(ctx context.Context, tenantID string) ([]ComplianceAssessment, error) {
	query := `
		SELECT a.id, a.tenant_id, a.requirement_id, a.status, a.evidence, a.assessor_id, 
		       a.assessed_at, a.next_review_date, a.notes, a.created_at, a.updated_at,
		       r.title as requirement_title, s.name as standard_name,
		       u.first_name || ' ' || u.last_name as assessor_name
		FROM compliance_assessments a
		JOIN compliance_requirements r ON a.requirement_id = r.id
		JOIN compliance_standards s ON r.standard_id = s.id
		LEFT JOIN users u ON a.assessor_id = u.id
		WHERE a.tenant_id = $1
		ORDER BY a.created_at DESC
	`
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ComplianceAssessment
	for rows.Next() {
		var a ComplianceAssessment
		if err := rows.Scan(&a.ID, &a.TenantID, &a.RequirementID, &a.Status, &a.Evidence, &a.AssessorID, &a.AssessedAt, &a.NextReviewDate, &a.Notes, &a.CreatedAt, &a.UpdatedAt, &a.RequirementTitle, &a.StandardName, &a.AssessorName); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, nil
}

func (r *ComplianceRepo) CreateAssessment(ctx context.Context, a ComplianceAssessment) error {
	_, err := r.db.Exec(`INSERT INTO compliance_assessments(id,tenant_id,requirement_id,status,evidence,assessor_id,assessed_at,next_review_date,notes) VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8)`, a.TenantID, a.RequirementID, a.Status, a.Evidence, a.AssessorID, a.AssessedAt, a.NextReviewDate, a.Notes)
	return err
}

func (r *ComplianceRepo) UpdateAssessment(ctx context.Context, id string, a ComplianceAssessment) error {
	_, err := r.db.Exec(`UPDATE compliance_assessments SET status=$1, evidence=$2, assessor_id=$3, assessed_at=$4, next_review_date=$5, notes=$6, updated_at=CURRENT_TIMESTAMP WHERE id=$7`, a.Status, a.Evidence, a.AssessorID, a.AssessedAt, a.NextReviewDate, a.Notes, id)
	return err
}

// Gaps
func (r *ComplianceRepo) ListGaps(ctx context.Context, assessmentID string) ([]ComplianceGap, error) {
	query := `
		SELECT g.id, g.assessment_id, g.title, g.description, g.severity, g.status,
		       g.remediation_plan, g.target_date, g.responsible_id, g.created_at, g.updated_at,
		       u.first_name || ' ' || u.last_name as responsible_name
		FROM compliance_gaps g
		LEFT JOIN users u ON g.responsible_id = u.id
		WHERE g.assessment_id = $1
		ORDER BY g.created_at DESC
	`
	rows, err := r.db.Query(query, assessmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ComplianceGap
	for rows.Next() {
		var gap ComplianceGap
		if err := rows.Scan(&gap.ID, &gap.AssessmentID, &gap.Title, &gap.Description, &gap.Severity, &gap.Status, &gap.RemediationPlan, &gap.TargetDate, &gap.ResponsibleID, &gap.CreatedAt, &gap.UpdatedAt, &gap.ResponsibleName); err != nil {
			return nil, err
		}
		items = append(items, gap)
	}
	return items, nil
}

func (r *ComplianceRepo) CreateGap(ctx context.Context, gap ComplianceGap) error {
	_, err := r.db.Exec(`INSERT INTO compliance_gaps(id,assessment_id,title,description,severity,status,remediation_plan,target_date,responsible_id) VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8)`, gap.AssessmentID, gap.Title, gap.Description, gap.Severity, gap.Status, gap.RemediationPlan, gap.TargetDate, gap.ResponsibleID)
	return err
}

func (r *ComplianceRepo) UpdateGap(ctx context.Context, id string, gap ComplianceGap) error {
	_, err := r.db.Exec(`UPDATE compliance_gaps SET status=$1, remediation_plan=$2, target_date=$3, responsible_id=$4, updated_at=CURRENT_TIMESTAMP WHERE id=$5`, gap.Status, gap.RemediationPlan, gap.TargetDate, gap.ResponsibleID, id)
	return err
}
