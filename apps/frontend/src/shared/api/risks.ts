import { api } from './client'

export interface PaginationMeta {
  page: number
  page_size: number
  total: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: PaginationMeta
}

export interface Risk {
  id: string
  tenant_id: string
  title: string
  description: string | null
  category: string | null
  likelihood: number | null
  impact: number | null
  level: number | null
  status: string
  owner_user_id: string | null
  asset_id: string | null
  methodology: string | null
  strategy: string | null
  due_date: string | null
  created_at: string
  updated_at: string
  level_label?: string | null
}

export interface CreateRiskRequest {
  title: string
  description?: string
  category?: string
  likelihood: number
  impact: number
  owner_user_id?: string
  asset_id?: string
  methodology?: string
  strategy?: string
  due_date?: string
}

export interface UpdateRiskRequest {
  title?: string
  description?: string
  category?: string
  likelihood?: number
  impact?: number
  status?: string
  owner_user_id?: string
  asset_id?: string
  methodology?: string
  strategy?: string
  due_date?: string
}

export interface RiskListParams {
  page?: number
  page_size?: number
  asset_id?: string
  status?: string
  level?: string
  owner_user_id?: string
  methodology?: string
  strategy?: string
  search?: string
  category?: string
  sort_field?: string
  sort_direction?: string
}

export interface RiskControl {
  id: string
  risk_id: string
  control_id: string
  control_name: string
  control_type: string
  implementation_status: string
  effectiveness?: string | null
  description?: string | null
  created_by?: string | null
  created_at: string
  updated_at: string
}

export interface CreateRiskControlRequest {
  control_id: string
  control_name: string
  control_type: string
  implementation_status: string
  effectiveness?: string
  description?: string
}

export interface UpdateRiskControlRequest {
  control_name: string
  control_type: string
  implementation_status: string
  effectiveness?: string
  description?: string
}

export interface RiskComment {
  id: string
  risk_id: string
  user_id: string
  comment: string
  is_internal: boolean
  user_name?: string | null
  created_at: string
  updated_at: string
}

export interface CreateRiskCommentRequest {
  comment: string
  is_internal?: boolean
}

export interface RiskHistoryEntry {
  id: string
  risk_id: string
  field_changed: string
  old_value?: string | null
  new_value?: string | null
  change_reason?: string | null
  changed_by: string
  changed_at: string
  changed_by_name?: string | null
}

export interface RiskAttachment {
  id: string
  risk_id: string
  file_name: string
  file_path: string
  file_size: number
  mime_type: string
  file_hash?: string | null
  description?: string | null
  uploaded_by: string
  uploaded_at: string
  uploaded_by_name?: string | null
}

export interface CreateRiskAttachmentRequest {
  file_name: string
  file_path: string
  file_size: number
  mime_type: string
  file_hash?: string
  description?: string
}

export const RISK_LEVELS = [
  { value: 'Low', label: 'Низкий (1-2)' },
  { value: 'Medium', label: 'Средний (3-4)' },
  { value: 'High', label: 'Высокий (5-6)' },
  { value: 'Critical', label: 'Критический (7+)' },
]

export const LIKELIHOOD_LEVELS = [
  { value: 1, label: '1 - Очень низкая' },
  { value: 2, label: '2 - Низкая' },
  { value: 3, label: '3 - Средняя' },
  { value: 4, label: '4 - Высокая' },
]

export const IMPACT_LEVELS = [
  { value: 1, label: '1 - Низкое' },
  { value: 2, label: '2 - Среднее' },
  { value: 3, label: '3 - Высокое' },
  { value: 4, label: '4 - Критическое' },
]

export const RISK_STATUSES = [
  { value: 'new', label: 'Новый' },
  { value: 'in_analysis', label: 'В анализе' },
  { value: 'in_treatment', label: 'В обработке' },
  { value: 'accepted', label: 'Принят' },
  { value: 'transferred', label: 'Передан' },
  { value: 'mitigated', label: 'Смягчен' },
  { value: 'closed', label: 'Закрыт' },
]

export const RISK_METHODOLOGIES = [
  { value: 'ISO27005', label: 'ISO 27005' },
  { value: 'NIST', label: 'NIST' },
  { value: 'COSO', label: 'COSO' },
  { value: 'Custom', label: 'Собственная' },
]

export const RISK_STRATEGIES = [
  { value: 'accept', label: 'Принять' },
  { value: 'mitigate', label: 'Смягчить' },
  { value: 'transfer', label: 'Передать' },
  { value: 'avoid', label: 'Избежать' },
]

export const RISK_CATEGORIES = [
  { value: 'security', label: 'Безопасность' },
  { value: 'operational', label: 'Операционные' },
  { value: 'financial', label: 'Финансовые' },
  { value: 'compliance', label: 'Соответствие' },
  { value: 'reputation', label: 'Репутационные' },
  { value: 'legal', label: 'Правовые' },
  { value: 'strategic', label: 'Стратегические' },
]

export const CONTROL_TYPES = [
  { value: 'preventive', label: 'Предотвращающий' },
  { value: 'detective', label: 'Обнаруживающий' },
  { value: 'corrective', label: 'Корректирующий' },
]

export const CONTROL_IMPLEMENTATION_STATUSES = [
  { value: 'planned', label: 'Запланирован' },
  { value: 'in_progress', label: 'В процессе' },
  { value: 'implemented', label: 'Реализован' },
  { value: 'not_applicable', label: 'Не применимо' },
]

export const CONTROL_EFFECTIVENESS = [
  { value: 'high', label: 'Высокая' },
  { value: 'medium', label: 'Средняя' },
  { value: 'low', label: 'Низкая' },
]

const mapListParamsToSearch = (params: RiskListParams): URLSearchParams => {
  const searchParams = new URLSearchParams()

  if (params.page) searchParams.append('page', params.page.toString())
  if (params.page_size) searchParams.append('page_size', params.page_size.toString())
  if (params.asset_id) searchParams.append('asset_id', params.asset_id)
  if (params.status) searchParams.append('status', params.status)
  if (params.level) searchParams.append('level', params.level)
  if (params.owner_user_id) searchParams.append('owner_user_id', params.owner_user_id)
  if (params.methodology) searchParams.append('methodology', params.methodology)
  if (params.strategy) searchParams.append('strategy', params.strategy)
  if (params.search) searchParams.append('search', params.search)
  if (params.category) searchParams.append('category', params.category)
  if (params.sort_field) searchParams.append('sort_field', params.sort_field)
  if (params.sort_direction) searchParams.append('sort_direction', params.sort_direction)

  return searchParams
}

export const risksApi = {
  async list(params: RiskListParams = {}): Promise<PaginatedResponse<Risk>> {
    const searchParams = mapListParamsToSearch(params)
    const response = await api.get(`/risks?${searchParams.toString()}`)
    return {
      data: response.data.data ?? [],
      pagination: response.data.pagination,
    }
  },

  async get(id: string): Promise<Risk> {
    const response = await api.get(`/risks/${id}`)
    return response.data.data
  },

  async create(payload: CreateRiskRequest): Promise<Risk> {
    const response = await api.post('/risks', payload)
    return response.data.data
  },

  async update(id: string, payload: UpdateRiskRequest): Promise<Risk> {
    const response = await api.patch(`/risks/${id}`, payload)
    return response.data.data
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/risks/${id}`)
  },


  async getByAsset(assetId: string): Promise<Risk[]> {
    const response = await api.get(`/risks/asset/${assetId}`)
    return response.data.data ?? []
  },
  async getControls(riskId: string): Promise<RiskControl[]> {
    const response = await api.get(`/risks/${riskId}/controls`)
    return response.data.data ?? []
  },

  async createControl(riskId: string, payload: CreateRiskControlRequest): Promise<void> {
    await api.post(`/risks/${riskId}/controls`, payload)
  },

  async updateControl(riskId: string, controlId: string, payload: UpdateRiskControlRequest): Promise<void> {
    await api.put(`/risks/${riskId}/controls/${controlId}`, payload)
  },

  async deleteControl(riskId: string, controlId: string): Promise<void> {
    await api.delete(`/risks/${riskId}/controls/${controlId}`)
  },

  async getComments(riskId: string, includeInternal = false): Promise<RiskComment[]> {
    const response = await api.get(`/risks/${riskId}/comments`, {
      params: { include_internal: includeInternal },
    })
    return response.data.data ?? []
  },

  async createComment(riskId: string, payload: CreateRiskCommentRequest): Promise<void> {
    await api.post(`/risks/${riskId}/comments`, payload)
  },

  async getHistory(riskId: string): Promise<RiskHistoryEntry[]> {
    const response = await api.get(`/risks/${riskId}/history`)
    return response.data.data ?? []
  },

  async getAttachments(riskId: string): Promise<RiskAttachment[]> {
    const response = await api.get(`/risks/${riskId}/attachments`)
    return response.data.data ?? []
  },

  async createAttachment(riskId: string, payload: CreateRiskAttachmentRequest): Promise<void> {
    await api.post(`/risks/${riskId}/attachments`, payload)
  },

  async deleteAttachment(riskId: string, attachmentId: string): Promise<void> {
    await api.delete(`/risks/${riskId}/attachments/${attachmentId}`)
  },
}
