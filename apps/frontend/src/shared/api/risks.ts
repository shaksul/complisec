import { apiClient } from './client';

export interface Risk {
  id: string;
  tenant_id: string;
  title: string;
  description?: string;
  category?: string;
  likelihood?: number;
  impact?: number;
  level?: number;
  status: string;
  owner_user_id?: string;
  asset_id?: string;
  methodology?: string;
  strategy?: string;
  due_date?: string;
  created_at: string;
  updated_at: string;
  level_label?: string;
}

export interface CreateRiskRequest {
  title: string;
  description?: string;
  category?: string;
  likelihood: number;
  impact: number;
  owner_user_id?: string;
  asset_id?: string;
  methodology?: string;
  strategy?: string;
  due_date?: string;
}

export interface UpdateRiskRequest {
  title?: string;
  description?: string;
  category?: string;
  likelihood?: number;
  impact?: number;
  status?: string;
  owner_user_id?: string;
  asset_id?: string;
  methodology?: string;
  strategy?: string;
  due_date?: string;
}

export interface RiskListParams {
  page?: number;
  page_size?: number;
  asset_id?: string;
  status?: string;
  level?: string;
  owner_user_id?: string;
  methodology?: string;
  strategy?: string;
  search?: string;
  category?: string;
  sort_field?: string;
  sort_direction?: string;
}

export const risksApi = {
  // List risks with pagination and filters
  list: async (params: RiskListParams = {}) => {
    const searchParams = new URLSearchParams();
    
    if (params.page) searchParams.append('page', params.page.toString());
    if (params.page_size) searchParams.append('page_size', params.page_size.toString());
    if (params.asset_id) searchParams.append('asset_id', params.asset_id);
    if (params.status) searchParams.append('status', params.status);
    if (params.level) searchParams.append('level', params.level);
    if (params.owner_user_id) searchParams.append('owner_user_id', params.owner_user_id);
    if (params.methodology) searchParams.append('methodology', params.methodology);
    if (params.strategy) searchParams.append('strategy', params.strategy);
    if (params.search) searchParams.append('search', params.search);
    if (params.category) searchParams.append('category', params.category);
    if (params.sort_field) searchParams.append('sort_field', params.sort_field);
    if (params.sort_direction) searchParams.append('sort_direction', params.sort_direction);
    
    const response = await apiClient.get(`/risks?${searchParams.toString()}`);
    return response.data;
  },

  // Get single risk
  get: async (id: string) => {
    const response = await apiClient.get(`/risks/${id}`);
    return response.data;
  },

  // Create new risk
  create: async (data: CreateRiskRequest) => {
    const response = await apiClient.post('/risks', data);
    return response.data;
  },

  // Update risk
  update: async (id: string, data: UpdateRiskRequest) => {
    const response = await apiClient.patch(`/risks/${id}`, data);
    return response.data;
  },

  // Delete risk
  delete: async (id: string) => {
    const response = await apiClient.delete(`/risks/${id}`);
    return response.data;
  },

  // Get risks by asset
  getByAsset: async (assetId: string) => {
    const response = await apiClient.get(`/risks?asset_id=${assetId}`);
    return response.data;
  }
};

// Risk level options (based on 1-4 scale)
export const RISK_LEVELS = [
  { value: 1, label: 'Low', color: 'success' },
  { value: 2, label: 'Medium', color: 'warning' },
  { value: 3, label: 'High', color: 'error' },
  { value: 4, label: 'Critical', color: 'error' }
];

// Likelihood scale (1-4)
export const LIKELIHOOD_LEVELS = [
  { value: 1, label: '1 - Очень низкая' },
  { value: 2, label: '2 - Низкая' },
  { value: 3, label: '3 - Средняя' },
  { value: 4, label: '4 - Высокая' }
];

// Impact scale (1-4)
export const IMPACT_LEVELS = [
  { value: 1, label: '1 - Низкое' },
  { value: 2, label: '2 - Среднее' },
  { value: 3, label: '3 - Высокое' },
  { value: 4, label: '4 - Критическое' }
];

// New risk statuses
export const RISK_STATUSES = [
  { value: 'new', label: 'Новый' },
  { value: 'in_analysis', label: 'В анализе' },
  { value: 'in_treatment', label: 'В обработке' },
  { value: 'accepted', label: 'Принят' },
  { value: 'transferred', label: 'Передан' },
  { value: 'mitigated', label: 'Смягчен' },
  { value: 'closed', label: 'Закрыт' }
];

// Risk methodologies
export const RISK_METHODOLOGIES = [
  { value: 'ISO27005', label: 'ISO 27005' },
  { value: 'NIST', label: 'NIST' },
  { value: 'COSO', label: 'COSO' },
  { value: 'Custom', label: 'Собственная' }
];

// Risk treatment strategies
export const RISK_STRATEGIES = [
  { value: 'accept', label: 'Принять' },
  { value: 'mitigate', label: 'Смягчить' },
  { value: 'transfer', label: 'Передать' },
  { value: 'avoid', label: 'Избежать' }
];

// Risk categories
export const RISK_CATEGORIES = [
  { value: 'security', label: 'Безопасность' },
  { value: 'operational', label: 'Операционные' },
  { value: 'financial', label: 'Финансовые' },
  { value: 'compliance', label: 'Соответствие' },
  { value: 'reputation', label: 'Репутационные' },
  { value: 'legal', label: 'Правовые' },
  { value: 'strategic', label: 'Стратегические' }
];

