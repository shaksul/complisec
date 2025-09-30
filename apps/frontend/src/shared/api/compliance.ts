import { api as apiClient } from './client'

export interface ComplianceStandard {
  id: string
  tenant_id: string
  name: string
  code: string
  description: string
  version: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface ComplianceRequirement {
  id: string
  tenant_id: string
  standard_id: string
  code: string
  description: string
  category: string
  created_at: string
  updated_at: string
}

export interface ComplianceAssessment {
  id: string
  tenant_id: string
  requirement_id: string
  asset_id?: string
  status: string
  notes?: string
  assessed_by?: string
  assessed_at: string
  created_at: string
  updated_at: string
}

export interface ComplianceRemediationPlan {
  id: string
  tenant_id: string
  assessment_id: string
  title: string
  description?: string
  status: string
  due_date?: string
  assigned_to?: string
  created_at: string
  updated_at: string
}

export interface CreateComplianceStandardDTO {
  name: string
  description?: string
}

export interface CreateComplianceRequirementDTO {
  standard_id: string
  code: string
  description: string
  category?: string
}

export interface CreateComplianceAssessmentDTO {
  requirement_id: string
  asset_id?: string
  status: string
  notes?: string
}

export interface CreateComplianceRemediationPlanDTO {
  assessment_id: string
  title: string
  description?: string
  status: string
  due_date?: string
  assigned_to?: string
}

// Standards API
export const getStandards = async (): Promise<ComplianceStandard[]> => {
  const response = await apiClient.get('/compliance/standards')
  return response.data.data || []
}

export const createStandard = async (data: CreateComplianceStandardDTO): Promise<ComplianceStandard> => {
  const response = await apiClient.post('/compliance/standards', data)
  return response.data
}

export const getStandard = async (id: string): Promise<ComplianceStandard> => {
  const response = await apiClient.get(`/compliance/standards/${id}`)
  return response.data
}

// Requirements API
export const getRequirements = async (standardId?: string): Promise<ComplianceRequirement[]> => {
  const url = standardId ? `/compliance/requirements?standard_id=${standardId}` : '/compliance/requirements'
  const response = await apiClient.get(url)
  return response.data
}

export const createRequirement = async (data: CreateComplianceRequirementDTO): Promise<ComplianceRequirement> => {
  const response = await apiClient.post('/compliance/requirements', data)
  return response.data
}

export const getRequirement = async (id: string): Promise<ComplianceRequirement> => {
  const response = await apiClient.get(`/compliance/requirements/${id}`)
  return response.data
}

// Assessments API
export const getAssessments = async (requirementId?: string): Promise<ComplianceAssessment[]> => {
  const url = requirementId ? `/compliance/assessments?requirement_id=${requirementId}` : '/compliance/assessments'
  const response = await apiClient.get(url)
  return response.data.data || []
}

export const createAssessment = async (data: CreateComplianceAssessmentDTO): Promise<ComplianceAssessment> => {
  const response = await apiClient.post('/compliance/assessments', data)
  return response.data
}

export const getAssessment = async (id: string): Promise<ComplianceAssessment> => {
  const response = await apiClient.get(`/compliance/assessments/${id}`)
  return response.data
}

// Remediation Plans API
export const getRemediationPlans = async (assessmentId?: string): Promise<ComplianceRemediationPlan[]> => {
  const url = assessmentId ? `/compliance/remediation-plans?assessment_id=${assessmentId}` : '/compliance/remediation-plans'
  const response = await apiClient.get(url)
  return response.data
}

export const createRemediationPlan = async (data: CreateComplianceRemediationPlanDTO): Promise<ComplianceRemediationPlan> => {
  const response = await apiClient.post('/compliance/remediation-plans', data)
  return response.data
}

export const getRemediationPlan = async (id: string): Promise<ComplianceRemediationPlan> => {
  const response = await apiClient.get(`/compliance/remediation-plans/${id}`)
  return response.data
}
