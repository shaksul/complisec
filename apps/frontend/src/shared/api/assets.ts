import { apiClient } from './client';

export interface Asset {
  id: string;
  inventory_number: string;
  name: string;
  type: string;
  class: string;
  owner_id?: string;
  owner_name?: string;
  location?: string;
  criticality: string;
  confidentiality: string;
  integrity: string;
  availability: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface AssetWithDetails extends Asset {
  documents: AssetDocument[];
  software: AssetSoftware[];
  history: AssetHistory[];
}

export interface AssetDocument {
  id: string;
  asset_id: string;
  document_type: string;
  file_path: string;
  created_by: string;
  created_at: string;
}

export interface AssetSoftware {
  id: string;
  asset_id: string;
  software_name: string;
  version?: string;
  installed_at?: string;
  updated_at: string;
}

export interface AssetHistory {
  id: string;
  asset_id: string;
  field_changed: string;
  old_value?: string;
  new_value: string;
  changed_by: string;
  changed_at: string;
}

export interface CreateAssetRequest {
  name: string;
  type: string;
  class: string;
  owner_id?: string;
  location?: string;
  criticality: string;
  confidentiality: string;
  integrity: string;
  availability: string;
  status?: string;
}

export interface UpdateAssetRequest {
  name?: string;
  type?: string;
  class?: string;
  owner_id?: string;
  location?: string;
  criticality?: string;
  confidentiality?: string;
  integrity?: string;
  availability?: string;
  status?: string;
}

export interface AssetListParams {
  page?: number;
  page_size?: number;
  type?: string;
  class?: string;
  status?: string;
  criticality?: string;
  owner_id?: string;
  search?: string;
}

export interface AssetDocumentRequest {
  document_type: string;
  file_path: string;
}

export interface AssetSoftwareRequest {
  software_name: string;
  version?: string;
  installed_at?: string;
}

export interface AssetInventoryRequest {
  asset_ids: string[];
  action: 'verify' | 'update_status';
  status?: string;
  notes?: string;
}

export const assetsApi = {
  // List assets with pagination and filters
  list: async (params: AssetListParams = {}) => {
    const searchParams = new URLSearchParams();
    
    if (params.page) searchParams.append('page', params.page.toString());
    if (params.page_size) searchParams.append('page_size', params.page_size.toString());
    if (params.type) searchParams.append('type', params.type);
    if (params.class) searchParams.append('class', params.class);
    if (params.status) searchParams.append('status', params.status);
    if (params.criticality) searchParams.append('criticality', params.criticality);
    if (params.owner_id) searchParams.append('owner_id', params.owner_id);
    if (params.search) searchParams.append('search', params.search);

    const response = await apiClient.get(`/assets?${searchParams.toString()}`);
    return response.data;
  },

  // Get single asset
  get: async (id: string) => {
    const response = await apiClient.get(`/assets/${id}`);
    return response.data;
  },

  // Get asset with full details
  getDetails: async (id: string) => {
    const response = await apiClient.get(`/assets/${id}/details`);
    return response.data;
  },

  // Create new asset
  create: async (data: CreateAssetRequest) => {
    const response = await apiClient.post('/assets', data);
    return response.data;
  },

  // Update asset
  update: async (id: string, data: UpdateAssetRequest) => {
    const response = await apiClient.put(`/assets/${id}`, data);
    return response.data;
  },

  // Delete asset
  delete: async (id: string) => {
    const response = await apiClient.delete(`/assets/${id}`);
    return response.data;
  },

  // Asset documents
  getDocuments: async (id: string) => {
    const response = await apiClient.get(`/assets/${id}/documents`);
    return response.data;
  },

  addDocument: async (id: string, data: AssetDocumentRequest) => {
    const response = await apiClient.post(`/assets/${id}/documents`, data);
    return response.data;
  },

  // Asset software
  getSoftware: async (id: string) => {
    const response = await apiClient.get(`/assets/${id}/software`);
    return response.data;
  },

  addSoftware: async (id: string, data: AssetSoftwareRequest) => {
    const response = await apiClient.post(`/assets/${id}/software`, data);
    return response.data;
  },

  // Asset history
  getHistory: async (id: string) => {
    const response = await apiClient.get(`/assets/${id}/history`);
    return response.data;
  },

  // Inventory operations
  performInventory: async (data: AssetInventoryRequest) => {
    const response = await apiClient.post('/assets/inventory', data);
    return response.data;
  },

  // Export assets
  export: async (params: AssetListParams = {}) => {
    const searchParams = new URLSearchParams();
    
    if (params.type) searchParams.append('type', params.type);
    if (params.class) searchParams.append('class', params.class);
    if (params.status) searchParams.append('status', params.status);
    if (params.criticality) searchParams.append('criticality', params.criticality);
    if (params.owner_id) searchParams.append('owner_id', params.owner_id);
    if (params.search) searchParams.append('search', params.search);

    const response = await apiClient.get(`/assets/export?${searchParams.toString()}`, {
      responseType: 'blob'
    });
    return response.data;
  }
};

// Asset type options
export const ASSET_TYPES = [
  { value: 'server', label: 'Сервер' },
  { value: 'workstation', label: 'Рабочая станция' },
  { value: 'application', label: 'Приложение' },
  { value: 'database', label: 'База данных' },
  { value: 'document', label: 'Документ' },
  { value: 'network_device', label: 'Сетевое устройство' },
  { value: 'other', label: 'Другое' }
];

export const ASSET_CLASSES = [
  { value: 'hardware', label: 'Оборудование' },
  { value: 'software', label: 'Программное обеспечение' },
  { value: 'data', label: 'Данные' },
  { value: 'service', label: 'Сервис' }
];

export const CRITICALITY_LEVELS = [
  { value: 'low', label: 'Низкая' },
  { value: 'medium', label: 'Средняя' },
  { value: 'high', label: 'Высокая' }
];

export const ASSET_STATUSES = [
  { value: 'active', label: 'Активен' },
  { value: 'in_repair', label: 'В ремонте' },
  { value: 'storage', label: 'На хранении' },
  { value: 'decommissioned', label: 'Списан' }
];

export const DOCUMENT_TYPES = [
  { value: 'passport', label: 'Паспорт актива' },
  { value: 'transfer_act', label: 'Акт передачи' },
  { value: 'writeoff_act', label: 'Акт списания' },
  { value: 'repair_log', label: 'Журнал ремонта' },
  { value: 'other', label: 'Другое' }
];
