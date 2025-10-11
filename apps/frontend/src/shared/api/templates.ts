import { api } from './client'

// Template Types
export interface DocumentTemplate {
  id: string
  tenant_id: string
  name: string
  description?: string
  template_type: 'passport_pc' | 'passport_monitor' | 'passport_device' | 'transfer_act' | 'writeoff_act' | 'repair_log' | 'other'
  content: string
  is_system: boolean
  is_active: boolean
  created_by: string
  created_at: string
  updated_at: string
}

export interface CreateTemplateRequest {
  name: string
  description?: string
  template_type: string
  content: string
}

export interface UpdateTemplateRequest {
  name?: string
  description?: string
  template_type?: string
  content?: string
  is_active?: boolean
}

export interface FillTemplateRequest {
  template_id: string
  asset_id: string
  additional_data?: Record<string, any>
  save_as_document?: boolean
  document_title?: string
  generate_pdf?: boolean
}

export interface FillTemplateResponse {
  html?: string
  pdf_base64?: string
  document_id?: string
}

// Inventory Number Rules
export interface InventoryNumberRule {
  id: string
  tenant_id: string
  asset_type: string
  asset_class?: string
  pattern: string
  current_sequence: number
  prefix?: string
  description?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateInventoryRuleRequest {
  asset_type: string
  asset_class?: string
  pattern: string
  prefix?: string
  description?: string
}

export interface UpdateInventoryRuleRequest {
  pattern?: string
  current_sequence?: number
  prefix?: string
  description?: string
  is_active?: boolean
}

export interface GenerateInventoryNumberRequest {
  asset_type: string
  asset_class?: string
}

export interface GenerateInventoryNumberResponse {
  inventory_number: string
  pattern: string
  sequence: number
}

// Template Variables
export interface TemplateVariable {
  name: string
  placeholder: string
  description: string
  example: string
  category: string
}

export interface TemplateVariablesResponse {
  variables: TemplateVariable[]
}

// API Functions
export const templatesApi = {
  // Templates
  async listTemplates(filters?: {
    template_type?: string
    is_system?: boolean
    is_active?: boolean
    search?: string
  }): Promise<DocumentTemplate[]> {
    const response = await api.get('/admin/templates', { params: filters })
    return response.data.data
  },

  async getTemplate(id: string): Promise<DocumentTemplate> {
    const response = await api.get(`/admin/templates/${id}`)
    return response.data.data
  },

  async createTemplate(payload: CreateTemplateRequest): Promise<DocumentTemplate> {
    const response = await api.post('/admin/templates', payload)
    return response.data.data
  },

  async updateTemplate(id: string, payload: UpdateTemplateRequest): Promise<void> {
    await api.put(`/admin/templates/${id}`, payload)
  },

  async deleteTemplate(id: string): Promise<void> {
    await api.delete(`/admin/templates/${id}`)
  },

  async initializeDefaultTemplates(): Promise<void> {
    await api.post('/admin/templates/initialize-defaults')
  },

  async getTemplateVariables(): Promise<TemplateVariablesResponse> {
    const response = await api.get('/admin/templates/variables')
    return response.data.data
  },

  // Inventory Rules
  async listInventoryRules(): Promise<InventoryNumberRule[]> {
    const response = await api.get('/admin/inventory-rules')
    return response.data.data
  },

  async createInventoryRule(payload: CreateInventoryRuleRequest): Promise<InventoryNumberRule> {
    const response = await api.post('/admin/inventory-rules', payload)
    return response.data.data
  },

  async updateInventoryRule(id: string, payload: UpdateInventoryRuleRequest): Promise<void> {
    await api.put(`/admin/inventory-rules/${id}`, payload)
  },

  async generateInventoryNumber(assetId: string, payload: GenerateInventoryNumberRequest): Promise<GenerateInventoryNumberResponse> {
    const response = await api.post(`/assets/${assetId}/generate-inventory-number`, payload)
    return response.data.data
  },

  // Asset Template Operations
  async fillTemplate(assetId: string, payload: FillTemplateRequest): Promise<FillTemplateResponse> {
    const response = await api.post(`/assets/${assetId}/fill-template`, payload)
    return response.data.data
  },
}

// Template Type Labels
export const TEMPLATE_TYPE_LABELS: Record<string, string> = {
  passport_pc: 'Паспорт ПК',
  passport_monitor: 'Паспорт монитора',
  passport_device: 'Паспорт устройства',
  transfer_act: 'Акт передачи',
  writeoff_act: 'Акт списания',
  repair_log: 'Журнал ремонтов',
  other: 'Другое',
}

// Template Type Options
export const TEMPLATE_TYPE_OPTIONS = [
  { value: 'passport_pc', label: 'Паспорт ПК' },
  { value: 'passport_monitor', label: 'Паспорт монитора' },
  { value: 'passport_device', label: 'Паспорт устройства' },
  { value: 'transfer_act', label: 'Акт передачи' },
  { value: 'writeoff_act', label: 'Акт списания' },
  { value: 'repair_log', label: 'Журнал ремонтов' },
  { value: 'other', label: 'Другое' },
]


