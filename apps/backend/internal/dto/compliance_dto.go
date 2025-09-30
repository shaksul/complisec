package dto

type CreateComplianceStandardDTO struct {
	Name        string  `json:"name" validate:"required"`
	Code        string  `json:"code" validate:"required"`
	Description *string `json:"description"`
	Version     string  `json:"version"`
}

type CreateComplianceRequirementDTO struct {
	StandardID  string  `json:"standard_id" validate:"required"`
	Code        string  `json:"code" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Priority    string  `json:"priority"`
	IsMandatory bool    `json:"is_mandatory"`
}

type CreateComplianceAssessmentDTO struct {
	RequirementID  string  `json:"requirement_id" validate:"required"`
	Status         string  `json:"status"`
	Evidence       *string `json:"evidence"`
	AssessorID     *string `json:"assessor_id"`
	AssessedAt     *string `json:"assessed_at"`
	NextReviewDate *string `json:"next_review_date"`
	Notes          *string `json:"notes"`
}

type UpdateComplianceAssessmentDTO struct {
	Status         string  `json:"status"`
	Evidence       *string `json:"evidence"`
	AssessorID     *string `json:"assessor_id"`
	AssessedAt     *string `json:"assessed_at"`
	NextReviewDate *string `json:"next_review_date"`
	Notes          *string `json:"notes"`
}

type CreateComplianceGapDTO struct {
	AssessmentID    string  `json:"assessment_id" validate:"required"`
	Title           string  `json:"title" validate:"required"`
	Description     *string `json:"description"`
	Severity        string  `json:"severity" validate:"required"`
	Status          string  `json:"status"`
	RemediationPlan *string `json:"remediation_plan"`
	TargetDate      *string `json:"target_date"`
	ResponsibleID   *string `json:"responsible_id"`
}

type UpdateComplianceGapDTO struct {
	Status          string  `json:"status"`
	RemediationPlan *string `json:"remediation_plan"`
	TargetDate      *string `json:"target_date"`
	ResponsibleID   *string `json:"responsible_id"`
}
