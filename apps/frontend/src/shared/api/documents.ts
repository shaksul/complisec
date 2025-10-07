import { api as apiClient } from './client'

export interface Document {
  id: string
  tenant_id: string
  title: string
  code?: string
  description?: string
  type: 'policy' | 'standard' | 'procedure' | 'instruction' | 'act' | 'other'
  category?: string
  tags: string[]
  status: 'draft' | 'in_review' | 'approved' | 'obsolete'
  current_version: number
  owner_id?: string
  classification: 'Public' | 'Internal' | 'Confidential'
  effective_from?: string
  review_period_months: number
  asset_ids: string[]
  risk_ids: string[]
  control_ids: string[]
  storage_key?: string
  mime_type?: string
  size_bytes?: number
  checksum_sha256?: string
  ocr_text?: string
  av_scan_status: 'pending' | 'clean' | 'infected' | 'error'
  av_scan_result?: string
  created_by: string
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface DocumentVersion {
  id: string
  document_id: string
  version_number: number
  storage_key: string
  mime_type?: string
  size_bytes?: number
  checksum_sha256?: string
  ocr_text?: string
  av_scan_status: 'pending' | 'clean' | 'infected' | 'error'
  av_scan_result?: string
  created_by: string
  created_at: string
  deleted_at?: string
}

export interface DocumentAcknowledgment {
  id: string
  document_id: string
  version_id?: string
  user_id: string
  status: 'pending' | 'completed' | 'failed'
  quiz_score?: number
  quiz_passed: boolean
  acknowledged_at?: string
  deadline?: string
  created_at: string
  updated_at: string
}

export interface DocumentQuiz {
  id: string
  document_id: string
  version_id?: string
  question: string
  question_order: number
  options?: string
  correct_answer?: string
  is_active: boolean
  created_at: string
}

export interface CreateDocumentDTO {
  title: string
  description?: string
  type: 'policy' | 'standard' | 'procedure' | 'instruction' | 'act' | 'other'
  category?: string
  tags: string[]
}

export interface UpdateDocumentDTO {
  title: string
  description?: string
  type: 'policy' | 'standard' | 'procedure' | 'instruction' | 'act' | 'other'
  category?: string
  tags: string[]
  status: 'draft' | 'in_review' | 'approved' | 'obsolete'
}

export interface CreateDocumentVersionDTO {
  title: string
  content?: string
  file_path?: string
  file_size?: number
  mime_type?: string
  checksum_sha256?: string
}

export interface CreateDocumentAcknowledgmentDTO {
  user_id: string
  version_id?: string
  deadline?: string
}

export interface UpdateDocumentAcknowledgmentDTO {
  status: 'pending' | 'completed' | 'failed'
  quiz_score?: number
  quiz_passed: boolean
  acknowledged_at?: string
}

export interface CreateDocumentQuizDTO {
  question: string
  question_order: number
  options: string[]
  correct_answer: string
}

export interface DocumentFilters {
  status?: string
  type?: string
  category?: string
  search?: string
  page?: number
  limit?: number
}

// Document CRUD API
export const getDocuments = async (filters?: DocumentFilters): Promise<Document[]> => {
  const params = new URLSearchParams()
  
  if (filters?.status) params.append('status', filters.status)
  if (filters?.type) params.append('type', filters.type)
  if (filters?.category) params.append('category', filters.category)
  if (filters?.search) params.append('search', filters.search)
  if (filters?.page) params.append('page', filters.page.toString())
  if (filters?.limit) params.append('limit', filters.limit.toString())

  const url = `/documents${params.toString() ? `?${params.toString()}` : ''}`
  const response = await apiClient.get(url)
  return response.data.data || []
}

export const getDocument = async (id: string): Promise<Document> => {
  const response = await apiClient.get(`/documents/${id}`)
  return response.data.data
}

export const createDocument = async (data: CreateDocumentDTO): Promise<Document> => {
  const response = await apiClient.post('/documents', data)
  return response.data.data
}

export const updateDocument = async (id: string, data: UpdateDocumentDTO): Promise<Document> => {
  const response = await apiClient.put(`/documents/${id}`, data)
  return response.data.data
}

export const deleteDocument = async (id: string): Promise<void> => {
  await apiClient.delete(`/documents/${id}`)
}

// Document Versions API
export const getDocumentVersions = async (documentId: string): Promise<DocumentVersion[]> => {
  const response = await apiClient.get(`/documents/${documentId}/versions`)
  return response.data.data || []
}

export const getDocumentVersion = async (versionId: string): Promise<DocumentVersion> => {
  const response = await apiClient.get(`/documents/versions/${versionId}`)
  return response.data.data
}

export const createDocumentVersion = async (documentId: string, data: CreateDocumentVersionDTO): Promise<DocumentVersion> => {
  const response = await apiClient.post(`/documents/${documentId}/versions`, data)
  return response.data.data
}

// Document Acknowledgments API
export const getDocumentAcknowledgment = async (documentId: string): Promise<DocumentAcknowledgment[]> => {
  const response = await apiClient.get(`/documents/${documentId}/acknowledgments`)
  return response.data.data || []
}

export const createDocumentAcknowledgment = async (documentId: string, data: CreateDocumentAcknowledgmentDTO): Promise<DocumentAcknowledgment> => {
  const response = await apiClient.post(`/documents/${documentId}/acknowledgments`, data)
  return response.data.data
}

export const updateDocumentAcknowledgment = async (ackId: string, data: UpdateDocumentAcknowledgmentDTO): Promise<DocumentAcknowledgment> => {
  const response = await apiClient.put(`/documents/acknowledgments/${ackId}`, data)
  return response.data.data
}

// Document Quizzes API
export const getDocumentQuizzes = async (documentId: string): Promise<DocumentQuiz[]> => {
  const response = await apiClient.get(`/documents/${documentId}/quizzes`)
  return response.data.data || []
}

export const createDocumentQuiz = async (documentId: string, data: CreateDocumentQuizDTO): Promise<DocumentQuiz> => {
  const response = await apiClient.post(`/documents/${documentId}/quizzes`, data)
  return response.data.data
}

// User-specific API
export const getUserPendingAcknowledgment = async (): Promise<DocumentAcknowledgment[]> => {
  const response = await apiClient.get('/users/me/pending-acknowledgments')
  return response.data.data || []
}

// Document type labels
export const getDocumentTypeLabel = (type: string): string => {
  const labels: Record<string, string> = {
    policy: 'Политика',
    standard: 'Стандарт',
    procedure: 'Процедура',
    instruction: 'Инструкция',
    act: 'Акт',
    other: 'Другое'
  }
  return labels[type] || type
}

// Document status labels
export const getDocumentStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    draft: 'Черновик',
    in_review: 'На согласовании',
    approved: 'Утвержден',
    obsolete: 'Устарел'
  }
  return labels[status] || status
}

// Document status colors
export const getDocumentStatusColor = (status: string): 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning' => {
  const colors: Record<string, 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning'> = {
    draft: 'default',
    in_review: 'warning',
    approved: 'success',
    obsolete: 'error'
  }
  return colors[status] || 'default'
}

// New interfaces for enhanced functionality
export interface ApprovalWorkflow {
  id: string
  document_id: string
  workflow_type: 'sequential' | 'parallel'
  status: 'pending' | 'in_progress' | 'approved' | 'rejected' | 'cancelled'
  created_by: string
  created_at: string
  completed_at?: string
}

export interface ApprovalStep {
  id: string
  workflow_id: string
  step_order: number
  approver_id: string
  status: 'pending' | 'approved' | 'rejected' | 'skipped'
  comments?: string
  deadline?: string
  completed_at?: string
  created_at: string
}

export interface ACKCampaign {
  id: string
  document_id: string
  title: string
  description?: string
  audience_type: 'all' | 'role' | 'department' | 'custom'
  audience_ids: string[]
  deadline?: string
  quiz_id?: string
  status: 'draft' | 'active' | 'completed' | 'cancelled'
  created_by: string
  created_at: string
  completed_at?: string
}

export interface ACKAssignment {
  id: string
  campaign_id: string
  user_id: string
  status: 'pending' | 'completed' | 'overdue'
  quiz_score?: number
  quiz_passed: boolean
  completed_at?: string
  created_at: string
}

export interface TrainingMaterial {
  id: string
  tenant_id: string
  title: string
  description?: string
  type: 'document' | 'video' | 'presentation' | 'other'
  storage_key: string
  mime_type?: string
  size_bytes?: number
  checksum_sha256?: string
  created_by: string
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface TrainingAssignment {
  id: string
  material_id: string
  user_id: string
  deadline?: string
  quiz_id?: string
  quiz_passed: boolean
  quiz_score?: number
  status: 'assigned' | 'in_progress' | 'completed' | 'overdue'
  completed_at?: string
  created_at: string
}

export interface Quiz {
  id: string
  tenant_id: string
  title: string
  description?: string
  questions: QuizQuestion[]
  passing_score: number
  time_limit_minutes?: number
  created_by: string
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface QuizQuestion {
  id: string
  question: string
  options: string[]
  correct_answer: number
  explanation?: string
}

// Enhanced DTOs
export interface CreateDocumentVersionDTO {
  enableOCR: boolean
}

export interface SubmitDocumentDTO {
  workflow_type: 'sequential' | 'parallel'
  steps: ApprovalStepDTO[]
}

export interface ApprovalStepDTO {
  step_order: number
  approver_id: string
  deadline?: string
}

export interface ApprovalActionDTO {
  action: 'approve' | 'reject'
  comment?: string
}

export interface CreateACKCampaignDTO {
  title: string
  description?: string
  audience_type: 'all' | 'role' | 'department' | 'custom'
  audience_ids: string[]
  deadline?: string
  quiz_id?: string
}

export interface CreateTrainingMaterialDTO {
  title: string
  description?: string
  type: 'document' | 'video' | 'presentation' | 'other'
}

export interface CreateTrainingAssignmentDTO {
  user_ids: string[]
  deadline?: string
  quiz_id?: string
}

export interface CreateQuizDTO {
  title: string
  description?: string
  questions: QuizQuestionDTO[]
  passing_score: number
  time_limit_minutes?: number
}

export interface QuizQuestionDTO {
  question: string
  options: string[]
  correct_answer: number
  explanation?: string
}

export interface SubmitQuizAnswerDTO {
  answers: QuizAnswerDTO[]
}

export interface QuizAnswerDTO {
  question_id: string
  answer: number
}

// Enhanced API functions
export const uploadDocument = async (
  file: File,
  name: string,
  description?: string,
  tags?: string[],
  linkedTo?: { module: string; entity_id: string }
): Promise<Document> => {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('name', name)
  if (description) {
    formData.append('description', description)
  }
  if (tags && tags.length > 0) {
    formData.append('tags', JSON.stringify(tags))
  }
  if (linkedTo) {
    formData.append('linked_to', JSON.stringify(linkedTo))
  }

  const response = await apiClient.post('/documents/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  return response.data.data
}

export const uploadDocumentVersion = async (
  documentId: string, 
  file: File, 
  options: CreateDocumentVersionDTO
): Promise<DocumentVersion> => {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('title', options.title)
  formData.append('enableOCR', options.enableOCR.toString())

  const response = await apiClient.post(`/documents/${documentId}/versions`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  return response.data
}

export const submitDocumentForApproval = async (
  documentId: string, 
  data: SubmitDocumentDTO
): Promise<ApprovalWorkflow> => {
  const response = await apiClient.post(`/documents/${documentId}/submit`, data)
  return response.data
}

export const approveDocument = async (
  documentId: string, 
  stepId: string, 
  action: ApprovalActionDTO
): Promise<void> => {
  await apiClient.post(`/documents/${documentId}/approval/${stepId}`, action)
}

export const publishDocument = async (documentId: string): Promise<void> => {
  await apiClient.post(`/documents/${documentId}/publish`)
}

export const createACKCampaign = async (
  documentId: string, 
  data: CreateACKCampaignDTO
): Promise<ACKCampaign> => {
  const response = await apiClient.post(`/documents/${documentId}/ack-campaigns`, data)
  return response.data
}

// Training materials API
export const getTrainingMaterials = async (): Promise<TrainingMaterial[]> => {
  const response = await apiClient.get('/training-materials')
  return response.data
}

export const createTrainingMaterial = async (data: CreateTrainingMaterialDTO): Promise<TrainingMaterial> => {
  const response = await apiClient.post('/training-materials', data)
  return response.data
}

export const uploadTrainingMaterial = async (
  materialId: string, 
  file: File
): Promise<TrainingMaterial> => {
  const formData = new FormData()
  formData.append('file', file)

  const response = await apiClient.post(`/training-materials/${materialId}/upload`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  return response.data
}

export const assignTrainingMaterial = async (
  materialId: string, 
  data: CreateTrainingAssignmentDTO
): Promise<TrainingAssignment[]> => {
  const response = await apiClient.post(`/training-materials/${materialId}/assign`, data)
  return response.data
}

// Quizzes API
export const getQuizzes = async (): Promise<Quiz[]> => {
  const response = await apiClient.get('/quizzes')
  return response.data
}

export const createQuiz = async (data: CreateQuizDTO): Promise<Quiz> => {
  const response = await apiClient.post('/quizzes', data)
  return response.data
}

export const submitQuizAnswers = async (
  quizId: string, 
  data: SubmitQuizAnswerDTO
): Promise<{ score: number; passed: boolean }> => {
  const response = await apiClient.post(`/quizzes/${quizId}/submit`, data)
  return response.data
}

// File operations API
export const downloadDocument = async (documentId: string): Promise<Blob> => {
  const response = await apiClient.get(`/documents/${documentId}/download`, {
    responseType: 'blob'
  })
  return response.data
}

export const downloadDocumentVersion = async (versionId: string): Promise<Blob> => {
  const response = await apiClient.get(`/documents/versions/${versionId}/download`, {
    responseType: 'blob'
  })
  return response.data
}

export const getDocumentVersionPreview = async (versionId: string): Promise<string> => {
  const response = await apiClient.get(`/documents/versions/${versionId}/preview`)
  return response.data.url || response.data
}
