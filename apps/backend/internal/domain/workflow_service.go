package domain

import (
	"context"
	"fmt"
	"time"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"
)

type WorkflowService struct {
	workflowRepo *repo.WorkflowRepo
	documentRepo *repo.DocumentRepo
	userRepo     *repo.UserRepo
}

func NewWorkflowService(workflowRepo *repo.WorkflowRepo, documentRepo *repo.DocumentRepo, userRepo *repo.UserRepo) *WorkflowService {
	return &WorkflowService{
		workflowRepo: workflowRepo,
		documentRepo: documentRepo,
		userRepo:     userRepo,
	}
}

// SubmitDocumentForApproval submits a document for approval workflow
func (s *WorkflowService) SubmitDocumentForApproval(ctx context.Context, tenantID, userID, documentID string, req dto.SubmitDocumentDTO) (*repo.ApprovalWorkflow, error) {
	// Check if document exists and is in draft status
	document, err := s.documentRepo.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	if document == nil {
		return nil, fmt.Errorf("document not found")
	}
	if document.Status != "draft" {
		return nil, fmt.Errorf("document must be in draft status to submit for approval")
	}

	// Create workflow
	workflow := repo.ApprovalWorkflow{
		ID:           generateUUID(),
		DocumentID:   documentID,
		WorkflowType: req.WorkflowType,
		Status:       "pending",
		CreatedBy:    userID,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	// Note: Steps are created separately via CreateApprovalStep

	// Save workflow
	err = s.workflowRepo.CreateWorkflow(ctx, workflow)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	// Update document status
	document.Status = "in_review"
	err = s.documentRepo.UpdateDocument(ctx, *document)
	if err != nil {
		return nil, fmt.Errorf("failed to update document status: %w", err)
	}

	// Note: Notifications are handled separately

	return &workflow, nil
}

// ProcessApprovalAction processes an approval or rejection action
func (s *WorkflowService) ProcessApprovalAction(ctx context.Context, tenantID, userID, documentID, stepID string, action dto.ApprovalActionDTO) error {
	// Get workflow and step
	workflow, err := s.workflowRepo.GetWorkflowByDocumentID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("failed to get workflow: %w", err)
	}
	if workflow == nil {
		return fmt.Errorf("workflow not found")
	}

	// Find the step
	step, err := s.workflowRepo.GetApprovalStep(ctx, stepID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get approval step: %w", err)
	}
	if step == nil {
		return fmt.Errorf("approval step not found")
	}

	// Check if user is authorized to approve this step
	if step.ApproverID != userID {
		return fmt.Errorf("user not authorized to approve this step")
	}

	// Update step status
	step.Status = action.Action
	step.Comments = action.Comment
	completedAt := time.Now().Format(time.RFC3339)
	step.CompletedAt = &completedAt

	err = s.workflowRepo.UpdateApprovalStep(ctx, *step)
	if err != nil {
		return fmt.Errorf("failed to update approval step: %w", err)
	}

	// Check if workflow is complete
	if action.Action == "approved" {
		if s.isWorkflowComplete(workflow) {
			// All steps approved, complete workflow
			workflow.Status = "approved"
			completedAt := time.Now().Format(time.RFC3339)
			workflow.CompletedAt = &completedAt
			err = s.workflowRepo.UpdateWorkflow(ctx, *workflow)
			if err != nil {
				return fmt.Errorf("failed to update workflow: %w", err)
			}

			// Update document status
			document, err := s.documentRepo.GetDocument(ctx, documentID, tenantID)
			if err != nil {
				return fmt.Errorf("failed to get document: %w", err)
			}
			document.Status = "approved"
			err = s.documentRepo.UpdateDocument(ctx, *document)
			if err != nil {
				return fmt.Errorf("failed to update document status: %w", err)
			}
		} else {
			// Move to next step
			nextStep := s.getNextStep(workflow)
			if nextStep != nil {
				s.sendApprovalNotification(ctx, tenantID, nextStep.ApproverID, documentID, workflow.ID, nextStep.ID)
			}
		}
	} else if action.Action == "rejected" {
		// Reject workflow
		workflow.Status = "rejected"
		completedAt := time.Now().Format(time.RFC3339)
		workflow.CompletedAt = &completedAt
		err = s.workflowRepo.UpdateWorkflow(ctx, *workflow)
		if err != nil {
			return fmt.Errorf("failed to update workflow: %w", err)
		}

		// Update document status back to draft
		document, err := s.documentRepo.GetDocument(ctx, documentID, tenantID)
		if err != nil {
			return fmt.Errorf("failed to get document: %w", err)
		}
		document.Status = "draft"
		err = s.documentRepo.UpdateDocument(ctx, *document)
		if err != nil {
			return fmt.Errorf("failed to update document status: %w", err)
		}
	}

	return nil
}

// CreateACKCampaign creates an acknowledgment campaign
func (s *WorkflowService) CreateACKCampaign(ctx context.Context, tenantID, userID, documentID string, req dto.CreateACKCampaignDTO) (*repo.ACKCampaign, error) {
	// Check if document exists and is approved
	document, err := s.documentRepo.GetDocument(ctx, documentID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	if document == nil {
		return nil, fmt.Errorf("document not found")
	}
	if document.Status != "approved" {
		return nil, fmt.Errorf("document must be approved to create acknowledgment campaign")
	}

	// Create campaign
	campaign := repo.ACKCampaign{
		ID:           generateUUID(),
		DocumentID:   documentID,
		Title:        req.Title,
		Description:  req.Description,
		AudienceType: req.AudienceType,
		AudienceIDs:  req.AudienceIDs,
		Deadline:     req.Deadline,
		QuizID:       req.QuizID,
		Status:       "draft",
		CreatedBy:    userID,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	// Save campaign
	err = s.workflowRepo.CreateACKCampaign(ctx, campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACK campaign: %w", err)
	}

	// Create assignments based on audience
	assignments, err := s.createACKAssignments(ctx, tenantID, campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACK assignments: %w", err)
	}

	// Send notifications
	for _, assignment := range assignments {
		s.sendACKNotification(ctx, tenantID, assignment.UserID, documentID, campaign.ID)
	}

	return &campaign, nil
}

// Helper functions
func (s *WorkflowService) isWorkflowComplete(workflow *repo.ApprovalWorkflow) bool {
	// Get all steps for this workflow
	steps, err := s.workflowRepo.ListApprovalSteps(context.Background(), workflow.ID, "00000000-0000-0000-0000-000000000001")
	if err != nil {
		return false
	}

	if workflow.WorkflowType == "sequential" {
		// Sequential: all steps must be approved in order
		for i, step := range steps {
			if step.Status != "approved" {
				return false
			}
			if i > 0 && steps[i-1].Status != "approved" {
				return false
			}
		}
		return true
	} else {
		// Parallel: all steps must be approved
		for _, step := range steps {
			if step.Status != "approved" {
				return false
			}
		}
		return true
	}
}

func (s *WorkflowService) getNextStep(workflow *repo.ApprovalWorkflow) *repo.ApprovalStep {
	// Get all steps for this workflow
	steps, err := s.workflowRepo.ListApprovalSteps(context.Background(), workflow.ID, "00000000-0000-0000-0000-000000000001")
	if err != nil {
		return nil
	}

	if workflow.WorkflowType == "sequential" {
		// Find first pending step
		for i := range steps {
			if steps[i].Status == "pending" {
				return &steps[i]
			}
		}
	} else {
		// Parallel: all steps are independent
		for i := range steps {
			if steps[i].Status == "pending" {
				return &steps[i]
			}
		}
	}
	return nil
}

func (s *WorkflowService) createACKAssignments(ctx context.Context, tenantID string, campaign repo.ACKCampaign) ([]repo.ACKAssignment, error) {
	var assignments []repo.ACKAssignment

	// Get users based on audience type
	var userIDs []string
	switch campaign.AudienceType {
	case "all":
		users, err := s.userRepo.GetUsersByTenant(ctx, tenantID)
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}
	case "custom":
		userIDs = campaign.AudienceIDs
	default:
		// For role/department, we would need additional logic
		userIDs = campaign.AudienceIDs
	}

	// Create assignments
	for _, userID := range userIDs {
		assignment := repo.ACKAssignment{
			ID:         generateUUID(),
			CampaignID: campaign.ID,
			UserID:     userID,
			Status:     "pending",
			CreatedAt:  time.Now().Format(time.RFC3339),
		}
		assignments = append(assignments, assignment)
	}

	// Save assignments
	for _, assignment := range assignments {
		err := s.workflowRepo.CreateACKAssignment(ctx, assignment)
		if err != nil {
			return nil, err
		}
	}

	return assignments, nil
}

func (s *WorkflowService) sendApprovalNotification(ctx context.Context, tenantID, userID, documentID, workflowID, stepID string) {
	// Implementation would send notification to user
	fmt.Printf("Sending approval notification to user %s for document %s\n", userID, documentID)
}

func (s *WorkflowService) sendACKNotification(ctx context.Context, tenantID, userID, documentID, campaignID string) {
	// Implementation would send notification to user
	fmt.Printf("Sending ACK notification to user %s for document %s\n", userID, documentID)
}
