import { apiClient } from './client';

export interface Incident {
  id: string;
  tenant_id: string;
  title: string;
  description?: string;
  category: string;
  status: string;
  criticality: string;
  source: string;
  reported_by: string;
  assigned_to?: string;
  detected_at: string;
  resolved_at?: string;
  closed_at?: string;
  created_at: string;
  updated_at: string;
  assets?: AssetInfo[];
  risks?: RiskInfo[];
  reported_name?: string;
  assigned_name?: string;
}

export interface AssetInfo {
  id: string;
  name: string;
}

export interface RiskInfo {
  id: string;
  name: string;
}

export interface IncidentComment {
  id: string;
  incident_id: string;
  user_id: string;
  comment: string;
  is_internal: boolean;
  created_at: string;
  user_name?: string;
}

export interface IncidentAction {
  id: string;
  incident_id: string;
  action_type: string;
  title: string;
  description?: string;
  assigned_to?: string;
  due_date?: string;
  completed_at?: string;
  status: string;
  created_by: string;
  created_at: string;
  updated_at: string;
  assigned_name?: string;
  created_name?: string;
}

export interface IncidentMetrics {
  total_incidents: number;
  open_incidents: number;
  closed_incidents: number;
  average_mttr_hours: number;
  average_mttd_hours: number;
  by_criticality: Record<string, number>;
  by_category: Record<string, number>;
  by_status: Record<string, number>;
}

export interface CreateIncidentRequest {
  title: string;
  description?: string;
  category: string;
  criticality: string;
  source: string;
  asset_ids?: string[];
  risk_ids?: string[];
  assigned_to?: string;
  detected_at?: string;
}

export interface UpdateIncidentRequest {
  title?: string;
  description?: string;
  category?: string;
  criticality?: string;
  status?: string;
  asset_ids?: string[];
  risk_ids?: string[];
  assigned_to?: string;
  detected_at?: string;
}

export interface IncidentListRequest {
  page?: number;
  page_size?: number;
  status?: string;
  criticality?: string;
  category?: string;
  asset_id?: string;
  risk_id?: string;
  assigned_to?: string;
  search?: string;
}

export interface IncidentStatusUpdateRequest {
  status: string;
}

export interface IncidentCommentRequest {
  comment: string;
  is_internal: boolean;
}

export interface IncidentActionRequest {
  action_type: string;
  title: string;
  description?: string;
  assigned_to?: string;
  due_date?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

// Constants
export const INCIDENT_CATEGORIES = {
  TECHNICAL_FAILURE: 'technical_failure',
  DATA_BREACH: 'data_breach',
  UNAUTHORIZED_ACCESS: 'unauthorized_access',
  PHYSICAL: 'physical',
  MALWARE: 'malware',
  SOCIAL_ENGINEERING: 'social_engineering',
  OTHER: 'other',
} as const;

export const INCIDENT_CRITICALITY = {
  LOW: 'low',
  MEDIUM: 'medium',
  HIGH: 'high',
  CRITICAL: 'critical',
} as const;

export const INCIDENT_STATUS = {
  NEW: 'new',
  ASSIGNED: 'assigned',
  IN_PROGRESS: 'in_progress',
  RESOLVED: 'resolved',
  CLOSED: 'closed',
} as const;

export const INCIDENT_SOURCE = {
  USER_REPORT: 'user_report',
  AUTOMATIC_AGENT: 'automatic_agent',
  ADMIN_MANUAL: 'admin_manual',
  MONITORING: 'monitoring',
  SIEM: 'siem',
} as const;

export const ACTION_TYPES = {
  INVESTIGATION: 'investigation',
  CONTAINMENT: 'containment',
  ERADICATION: 'eradication',
  RECOVERY: 'recovery',
  PREVENTION: 'prevention',
} as const;

export const ACTION_STATUS = {
  PENDING: 'pending',
  IN_PROGRESS: 'in_progress',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled',
} as const;

// Helper function to clean incident data
function cleanIncidentData(data: any): any {
  const cleaned = { ...data };
  
  // Remove empty strings and convert to undefined
  if (cleaned.assigned_to === '' || cleaned.assigned_to === null) {
    delete cleaned.assigned_to;
  }
  
  if (cleaned.description === '' || cleaned.description === null) {
    delete cleaned.description;
  }
  
  if (cleaned.detected_at === '' || cleaned.detected_at === null) {
    delete cleaned.detected_at;
  }
  
  // Remove empty arrays
  if (cleaned.asset_ids && cleaned.asset_ids.length === 0) {
    delete cleaned.asset_ids;
  }
  
  if (cleaned.risk_ids && cleaned.risk_ids.length === 0) {
    delete cleaned.risk_ids;
  }
  
  return cleaned;
}

// API functions
export const incidentsApi = {

  // Get list of incidents
  async list(params: IncidentListRequest = {}): Promise<PaginatedResponse<Incident>> {
    const searchParams = new URLSearchParams();
    
    if (params.page) searchParams.append('page', params.page.toString());
    if (params.page_size) searchParams.append('page_size', params.page_size.toString());
    if (params.status) searchParams.append('status', params.status);
    if (params.criticality) searchParams.append('criticality', params.criticality);
    if (params.category) searchParams.append('category', params.category);
    if (params.asset_id) searchParams.append('asset_id', params.asset_id);
    if (params.risk_id) searchParams.append('risk_id', params.risk_id);
    if (params.assigned_to) searchParams.append('assigned_to', params.assigned_to);
    if (params.search) searchParams.append('search', params.search);

    const response = await apiClient.get(`/incidents?${searchParams.toString()}`);
    return response.data;
  },

  // Get single incident
  async get(id: string): Promise<Incident> {
    const response = await apiClient.get(`/incidents/${id}`);
    return response.data;
  },

  // Create incident
  async create(data: CreateIncidentRequest): Promise<Incident> {
    const cleanData = cleanIncidentData(data);
    const response = await apiClient.post('/incidents', cleanData);
    return response.data;
  },

  // Update incident
  async update(id: string, data: UpdateIncidentRequest): Promise<Incident> {
    const cleanData = cleanIncidentData(data);
    const response = await apiClient.put(`/incidents/${id}`, cleanData);
    return response.data;
  },

  // Delete incident
  async delete(id: string): Promise<void> {
    await apiClient.delete(`/incidents/${id}`);
  },

  // Update incident status
  async updateStatus(id: string, data: IncidentStatusUpdateRequest): Promise<Incident> {
    const response = await apiClient.put(`/incidents/${id}/status`, data);
    return response.data;
  },

  // Get incident metrics
  async getMetrics(): Promise<IncidentMetrics> {
    const response = await apiClient.get('/incidents/metrics');
    return response.data;
  },

  // Comments
  async addComment(incidentId: string, data: IncidentCommentRequest): Promise<IncidentComment> {
    const response = await apiClient.post(`/incidents/${incidentId}/comments`, data);
    return response.data;
  },

  async getComments(incidentId: string): Promise<IncidentComment[]> {
    const response = await apiClient.get(`/incidents/${incidentId}/comments`);
    return response.data;
  },

  // Actions
  async addAction(incidentId: string, data: IncidentActionRequest): Promise<IncidentAction> {
    const response = await apiClient.post(`/incidents/${incidentId}/actions`, data);
    return response.data;
  },

  async getActions(incidentId: string): Promise<IncidentAction[]> {
    const response = await apiClient.get(`/incidents/${incidentId}/actions`);
    return response.data;
  },

  async updateAction(incidentId: string, actionId: string, data: IncidentActionRequest): Promise<IncidentAction> {
    const response = await apiClient.put(`/incidents/${incidentId}/actions/${actionId}`, data);
    return response.data;
  },

  async deleteAction(incidentId: string, actionId: string): Promise<void> {
    await apiClient.delete(`/incidents/${incidentId}/actions/${actionId}`);
  },
};

// Helper functions
export const getCategoryLabel = (category: string): string => {
  const labels: Record<string, string> = {
    [INCIDENT_CATEGORIES.TECHNICAL_FAILURE]: 'Технический сбой',
    [INCIDENT_CATEGORIES.DATA_BREACH]: 'Утечка данных',
    [INCIDENT_CATEGORIES.UNAUTHORIZED_ACCESS]: 'Несанкционированный доступ',
    [INCIDENT_CATEGORIES.PHYSICAL]: 'Физический',
    [INCIDENT_CATEGORIES.MALWARE]: 'Вредоносное ПО',
    [INCIDENT_CATEGORIES.SOCIAL_ENGINEERING]: 'Социальная инженерия',
    [INCIDENT_CATEGORIES.OTHER]: 'Другое',
  };
  return labels[category] || category;
};

export const getCriticalityLabel = (criticality: string): string => {
  const labels: Record<string, string> = {
    [INCIDENT_CRITICALITY.LOW]: 'Низкий',
    [INCIDENT_CRITICALITY.MEDIUM]: 'Средний',
    [INCIDENT_CRITICALITY.HIGH]: 'Высокий',
    [INCIDENT_CRITICALITY.CRITICAL]: 'Критический',
  };
  return labels[criticality] || criticality;
};

export const getStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    [INCIDENT_STATUS.NEW]: 'Новый',
    [INCIDENT_STATUS.ASSIGNED]: 'Назначен',
    [INCIDENT_STATUS.IN_PROGRESS]: 'В работе',
    [INCIDENT_STATUS.RESOLVED]: 'Устранен',
    [INCIDENT_STATUS.CLOSED]: 'Закрыт',
  };
  return labels[status] || status;
};

export const getSourceLabel = (source: string): string => {
  const labels: Record<string, string> = {
    [INCIDENT_SOURCE.USER_REPORT]: 'Сообщение пользователя',
    [INCIDENT_SOURCE.AUTOMATIC_AGENT]: 'Автоматический агент',
    [INCIDENT_SOURCE.ADMIN_MANUAL]: 'Администратор',
    [INCIDENT_SOURCE.MONITORING]: 'Мониторинг',
    [INCIDENT_SOURCE.SIEM]: 'SIEM',
  };
  return labels[source] || source;
};

export const getActionTypeLabel = (actionType: string): string => {
  const labels: Record<string, string> = {
    [ACTION_TYPES.INVESTIGATION]: 'Расследование',
    [ACTION_TYPES.CONTAINMENT]: 'Сдерживание',
    [ACTION_TYPES.ERADICATION]: 'Устранение',
    [ACTION_TYPES.RECOVERY]: 'Восстановление',
    [ACTION_TYPES.PREVENTION]: 'Предотвращение',
  };
  return labels[actionType] || actionType;
};

export const getActionStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    [ACTION_STATUS.PENDING]: 'Ожидает',
    [ACTION_STATUS.IN_PROGRESS]: 'В работе',
    [ACTION_STATUS.COMPLETED]: 'Завершено',
    [ACTION_STATUS.CANCELLED]: 'Отменено',
  };
  return labels[status] || status;
};

export const getCriticalityColor = (criticality: string): string => {
  const colors: Record<string, string> = {
    [INCIDENT_CRITICALITY.LOW]: 'green',
    [INCIDENT_CRITICALITY.MEDIUM]: 'yellow',
    [INCIDENT_CRITICALITY.HIGH]: 'orange',
    [INCIDENT_CRITICALITY.CRITICAL]: 'red',
  };
  return colors[criticality] || 'gray';
};

export const getStatusColor = (status: string): string => {
  const colors: Record<string, string> = {
    [INCIDENT_STATUS.NEW]: 'blue',
    [INCIDENT_STATUS.ASSIGNED]: 'purple',
    [INCIDENT_STATUS.IN_PROGRESS]: 'orange',
    [INCIDENT_STATUS.RESOLVED]: 'green',
    [INCIDENT_STATUS.CLOSED]: 'gray',
  };
  return colors[status] || 'gray';
};