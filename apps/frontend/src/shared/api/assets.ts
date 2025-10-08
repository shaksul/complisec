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

export interface Asset {
  id: string
  tenant_id: string
  inventory_number: string
  name: string
  type: string
  class: string
  owner_id: string | null
  owner_name?: string | null
  responsible_user_id: string | null
  responsible_user_name?: string | null
  location: string | null
  criticality: string
  confidentiality: string
  integrity: string
  availability: string
  status: string
  created_at: string
  updated_at: string
  deleted_at?: string | null
}

export interface AssetDocument {
  id: string
  asset_id: string
  title: string
  document_type: string
  mime: string
  size_bytes: number
  download_url: string
  created_by: string
  created_by_name?: string | null  // ФИО пользователя
  created_by_email?: string | null // Email пользователя
  created_at: string
}

export interface StorageDocument {
  id: string
  title: string
  document_type: string
  version: string
  size_bytes: number
  mime: string
  created_by: string
  created_at: string
}

export interface AssetSoftware {
  id: string
  asset_id: string
  software_name: string
  version?: string | null
  installed_at?: string | null
  updated_at: string
}

export interface AssetHistory {
  id: string
  asset_id: string
  field_changed: string
  old_value?: string | null
  new_value: string
  changed_by: string
  changed_by_name?: string | null  // ФИО пользователя
  changed_by_email?: string | null // Email пользователя
  changed_at: string
}

export interface AssetWithDetails extends Asset {
  documents: AssetDocument[]
  software: AssetSoftware[]
  history: AssetHistory[]
}

export interface CreateAssetRequest {
  name: string
  type: string
  class: string
  owner_id?: string
  responsible_user_id?: string
  location?: string
  criticality: string
  confidentiality: string
  integrity: string
  availability: string
  status?: string
}

export interface UpdateAssetRequest {
  name?: string
  type?: string
  class?: string
  owner_id?: string
  responsible_user_id?: string
  location?: string
  criticality?: string
  confidentiality?: string
  integrity?: string
  availability?: string
  status?: string
}

export interface AssetListParams {
  page?: number
  page_size?: number
  type?: string
  class?: string
  status?: string
  criticality?: string
  owner_id?: string
  search?: string
}

export interface AssetDocumentRequest {
  document_type: string
  file_path: string
}

export interface AssetSoftwareRequest {
  software_name: string
  version?: string
  installed_at?: string
}

export interface AssetInventoryRequest {
  asset_ids: string[]
  action: 'verify' | 'update_status'
  status?: string
  notes?: string
}

export interface AssetHistoryFilters {
  changed_by?: string
  from_date?: string
  to_date?: string
}

export interface BulkUpdateStatusRequest {
  asset_ids: string[]
  status: string
}

export interface BulkUpdateOwnerRequest {
  asset_ids: string[]
  owner_id: string
}

export const assetsApi = {
  async list(params: AssetListParams = {}): Promise<PaginatedResponse<Asset>> {
    const response = await api.get('/assets', { params })
    return {
      data: response.data.data ?? [],
      pagination: response.data.pagination,
    }
  },

  async get(id: string): Promise<Asset> {
    const response = await api.get(`/assets/${id}`)
    return response.data.data
  },

  async getDetails(id: string): Promise<AssetWithDetails> {
    const response = await api.get(`/assets/${id}/details`)
    const asset = response.data.data
    if (!asset) {
      throw new Error(`Asset ${id} not found`)
    }
    return {
      ...asset,
      documents: asset.documents ?? [],
      software: asset.software ?? [],
      history: asset.history ?? [],
    }
  },

  async create(payload: CreateAssetRequest): Promise<Asset> {
    const response = await api.post('/assets', payload)
    return response.data.data
  },

  async update(id: string, payload: UpdateAssetRequest): Promise<void> {
    await api.put(`/assets/${id}`, payload)
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/assets/${id}`)
  },

  async getDocuments(id: string): Promise<AssetDocument[]> {
    const response = await api.get(`/assets/${id}/documents`)
    return response.data.data ?? []
  },

  async addDocument(id: string, payload: AssetDocumentRequest): Promise<void> {
    await api.post(`/assets/${id}/documents`, payload)
  },

  async uploadDocument(assetId: string, formData: FormData): Promise<AssetDocument> {
    const response = await api.post(`/assets/${assetId}/documents/upload`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data.data
  },

  async linkDocument(assetId: string, data: { document_id: string; document_type: string }): Promise<AssetDocument> {
    const response = await api.post(`/assets/${assetId}/documents/link`, data)
    return response.data.data
  },

  async getDocumentStorage(params: {
    query?: string;
    type?: string;
    page?: number;
    page_size?: number;
  }): Promise<{ data: StorageDocument[]; pagination: any }> {
    const searchParams = new URLSearchParams()
    if (params.query) searchParams.append('query', params.query)
    if (params.type) searchParams.append('type', params.type)
    if (params.page) searchParams.append('page', params.page.toString())
    if (params.page_size) searchParams.append('page_size', params.page_size.toString())
    
    const response = await api.get(`/assets/documents/storage?${searchParams}`)
    return response.data
  },

  async downloadDocument(documentId: string): Promise<Blob> {
    const response = await api.get(`/assets/documents/${documentId}/download`, {
      responseType: 'blob'
    })
    return response.data
  },

  async deleteDocument(documentId: string): Promise<void> {
    await api.delete(`/assets/documents/${documentId}`)
  },

  async getDocument(documentId: string): Promise<AssetDocument> {
    const response = await api.get(`/assets/documents/${documentId}`)
    return response.data.data
  },

  async getSoftware(id: string): Promise<AssetSoftware[]> {
    const response = await api.get(`/assets/${id}/software`)
    return response.data.data ?? []
  },

  async addSoftware(id: string, payload: AssetSoftwareRequest): Promise<void> {
    await api.post(`/assets/${id}/software`, payload)
  },

  async getHistory(id: string): Promise<AssetHistory[]> {
    const response = await api.get(`/assets/${id}/history`)
    return response.data.data ?? []
  },

  async getHistoryWithFilters(id: string, filters: AssetHistoryFilters): Promise<AssetHistory[]> {
    const response = await api.get(`/assets/${id}/history/filtered`, { params: filters })
    return response.data.data ?? []
  },

  async getRisks(id: string) {
    const response = await api.get(`/assets/${id}/risks`)
    return response.data.data ?? []
  },

  async getIncidents(id: string) {
    const response = await api.get(`/assets/${id}/incidents`)
    return response.data.data ?? []
  },

  async canAddRisk(id: string): Promise<boolean> {
    const response = await api.get(`/assets/${id}/can-add-risk`)
    return Boolean(response.data?.allowed ?? response.data?.data)
  },

  async canAddIncident(id: string): Promise<boolean> {
    const response = await api.get(`/assets/${id}/can-add-incident`)
    return Boolean(response.data?.allowed ?? response.data?.data)
  },

  async getAssetsWithoutOwner(): Promise<Asset[]> {
    const response = await api.get('/assets/inventory/without-owner')
    return response.data.data ?? []
  },

  async getAssetsWithoutPassport(): Promise<Asset[]> {
    const response = await api.get('/assets/inventory/without-passport')
    return response.data.data ?? []
  },

  async getAssetsWithoutCriticality(): Promise<Asset[]> {
    const response = await api.get('/assets/inventory/without-criticality')
    return response.data.data ?? []
  },

  async bulkUpdateStatus(payload: BulkUpdateStatusRequest): Promise<void> {
    await api.post('/assets/bulk/update-status', payload)
  },

  async bulkUpdateOwner(payload: BulkUpdateOwnerRequest): Promise<void> {
    await api.post('/assets/bulk/update-owner', payload)
  },

  async performInventory(payload: AssetInventoryRequest): Promise<void> {
    await api.post('/assets/inventory', payload)
  },

  async export(params: AssetListParams = {}): Promise<Blob> {
    const response = await api.get('/assets/export', {
      params,
      responseType: 'blob',
    })
    return response.data
  },
}

export const ASSET_TYPES = [
  { value: 'server', label: 'Сервер' },
  { value: 'workstation', label: 'Рабочая станция' },
  { value: 'application', label: 'Приложение' },
  { value: 'database', label: 'База данных' },
  { value: 'document', label: 'Документ' },
  { value: 'network_device', label: 'Сетевое устройство' },
  { value: 'other', label: 'Другое' },
]

export const ASSET_CLASSES = [
  { value: 'hardware', label: 'Оборудование' },
  { value: 'software', label: 'Программное обеспечение' },
  { value: 'data', label: 'Данные' },
  { value: 'service', label: 'Сервис' },
]

export const CRITICALITY_LEVELS = [
  { value: 'low', label: 'Низкая' },
  { value: 'medium', label: 'Средняя' },
  { value: 'high', label: 'Высокая' },
]

export const ASSET_STATUSES = [
  { value: 'active', label: 'Активен' },
  { value: 'in_repair', label: 'В ремонте' },
  { value: 'storage', label: 'На хранении' },
  { value: 'decommissioned', label: 'Списан' },
]

export const DOCUMENT_TYPES = [
  { value: 'passport', label: 'Паспорт актива' },
  { value: 'transfer_act', label: 'Акт передачи' },
  { value: 'writeoff_act', label: 'Акт списания' },
  { value: 'repair_log', label: 'Журнал ремонта' },
  { value: 'other', label: 'Другое' },
]
