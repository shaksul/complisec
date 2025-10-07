import { apiClient } from './client';

// Types
export interface Material {
  id: string;
  tenant_id: string;
  title: string;
  description?: string;
  uri: string;
  type: 'file' | 'link' | 'video';
  material_type: 'document' | 'video' | 'quiz' | 'simulation' | 'acknowledgment';
  duration_minutes?: number;
  tags: string[];
  is_required: boolean;
  passing_score: number;
  attempts_limit?: number;
  metadata: Record<string, any>;
  created_by?: string;
  created_at: string;
  updated_at: string;
}

export interface TrainingCourse {
  id: string;
  tenant_id: string;
  title: string;
  description?: string;
  is_active: boolean;
  created_by?: string;
  created_at: string;
  updated_at: string;
  materials?: Material[];
}

export interface CourseMaterial {
  id: string;
  course_id: string;
  material_id: string;
  order_index: number;
  is_required: boolean;
  created_at: string;
  material?: Material;
}

export interface CreateMaterialRequest {
  title: string;
  description?: string;
  uri: string;
  type: 'file' | 'link' | 'video';
  material_type: 'document' | 'video' | 'quiz' | 'simulation' | 'acknowledgment';
  duration_minutes?: number;
  tags?: string[];
  is_required?: boolean;
  passing_score?: number;
  attempts_limit?: number;
  metadata?: Record<string, any>;
}

export interface UpdateMaterialRequest {
  title?: string;
  description?: string;
  uri?: string;
  type?: 'file' | 'link' | 'video';
  material_type?: 'document' | 'video' | 'quiz' | 'simulation' | 'acknowledgment';
  duration_minutes?: number;
  tags?: string[];
  is_required?: boolean;
  passing_score?: number;
  attempts_limit?: number;
  metadata?: Record<string, any>;
}

export interface CreateCourseRequest {
  title: string;
  description?: string;
  is_active?: boolean;
}

export interface UpdateCourseRequest {
  title?: string;
  description?: string;
  is_active?: boolean;
}

export interface CourseMaterialRequest {
  order_index: number;
  is_required: boolean;
}

export interface TrainingListRequest {
  page?: number;
  page_size?: number;
  filters?: Record<string, any>;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

export interface TrainingListResponse<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

// Materials API
export const materialsApi = {
  // Create material
  create: (data: CreateMaterialRequest): Promise<Material> =>
    apiClient.post('/training/materials', data).then((res) => res.data),

  // Get material by ID
  getById: (id: string): Promise<Material> =>
    apiClient.get(`/training/materials/${id}`).then((res) => res.data),

  // List materials
  list: (params?: TrainingListRequest): Promise<TrainingListResponse<Material>> =>
    apiClient.get('/training/materials', { params }).then((res) => res.data),

  // Update material
  update: (id: string, data: UpdateMaterialRequest): Promise<void> =>
    apiClient.put(`/training/materials/${id}`, data).then(() => undefined),

  // Delete material
  delete: (id: string): Promise<void> =>
    apiClient.delete(`/training/materials/${id}`).then(() => undefined),
};

// Courses API
export const coursesApi = {
  // Create course
  create: (data: CreateCourseRequest): Promise<TrainingCourse> =>
    apiClient.post('/training/courses', data).then((res) => res.data),

  // Get course by ID
  getById: (id: string): Promise<TrainingCourse> =>
    apiClient.get(`/training/courses/${id}`).then((res) => res.data),

  // List courses
  list: (params?: TrainingListRequest): Promise<TrainingListResponse<TrainingCourse>> =>
    apiClient.get('/training/courses', { params }).then((res) => res.data),

  // Update course
  update: (id: string, data: UpdateCourseRequest): Promise<void> =>
    apiClient.put(`/training/courses/${id}`, data).then(() => undefined),

  // Delete course
  delete: (id: string): Promise<void> =>
    apiClient.delete(`/training/courses/${id}`).then(() => undefined),
};

// Course Materials API
export const courseMaterialsApi = {
  // Add material to course
  add: (courseId: string, materialId: string, data: CourseMaterialRequest): Promise<void> =>
    apiClient.post(`/training/courses/${courseId}/materials/${materialId}`, data).then(() => undefined),

  // Remove material from course
  remove: (courseId: string, materialId: string): Promise<void> =>
    apiClient.delete(`/training/courses/${courseId}/materials/${materialId}`).then(() => undefined),

  // Get course materials
  getByCourse: (courseId: string): Promise<CourseMaterial[]> =>
    apiClient.get(`/training/courses/${courseId}/materials`).then((res) => res.data),
};

// Constants
export const MATERIAL_TYPES = {
  document: 'Документ',
  video: 'Видео',
  quiz: 'Квиз',
  simulation: 'Симуляция',
  acknowledgment: 'Ознакомление',
} as const;

export const MATERIAL_SOURCES = {
  file: 'Файл',
  link: 'Ссылка',
  video: 'Видео',
} as const;

export const MATERIAL_TYPE_OPTIONS = Object.entries(MATERIAL_TYPES).map(([value, label]) => ({
  value,
  label,
}));

export const MATERIAL_SOURCE_OPTIONS = Object.entries(MATERIAL_SOURCES).map(([value, label]) => ({
  value,
  label,
}));






