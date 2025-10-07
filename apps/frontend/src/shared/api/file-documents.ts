import { api as apiClient } from './client'

// File Document Types
export interface FileDocument {
  id: string
  tenant_id: string
  name: string
  original_name: string
  description?: string
  file_path: string
  file_size: number
  mime_type: string
  file_hash: string
  folder_id?: string
  owner_id: string
  created_by: string
  created_at: string
  updated_at: string
  is_active: boolean
  version: number
  metadata?: string
  tags: string[]
  links: DocumentLink[]
  ocr_text?: string
}

export interface DocumentLink {
  module: 'risk' | 'asset' | 'incident' | 'training' | 'compliance'
  entity_id: string
}

export interface Folder {
  id: string
  tenant_id: string
  name: string
  description?: string
  parent_id?: string
  owner_id: string
  created_by: string
  created_at: string
  updated_at: string
  is_active: boolean
  metadata?: string
  children?: Folder[]
}

export interface DocumentPermission {
  id: string
  subject_type: 'user' | 'role'
  subject_id: string
  object_type: 'document' | 'folder'
  object_id: string
  permission: 'view' | 'edit' | 'delete' | 'share'
  granted_by: string
  granted_at: string
  expires_at?: string
  is_active: boolean
}

export interface DocumentVersion {
  id: string
  document_id: string
  version_number: number
  file_path: string
  file_size: number
  file_hash: string
  created_by: string
  created_at: string
  change_description?: string
}

export interface DocumentAuditLog {
  id: string
  tenant_id: string
  document_id?: string
  folder_id?: string
  user_id: string
  action: string
  details?: string
  ip_address?: string
  user_agent?: string
  created_at: string
}

export interface DocumentSearchResult {
  document_id: string
  name: string
  description?: string
  mime_type: string
  file_size: number
  created_at: string
  relevance_score?: number
  ocr_text?: string
}

export interface DocumentStats {
  total_documents: number
  total_folders: number
  total_size: number
  documents_by_type: Record<string, number>
  recent_documents: FileDocument[]
  storage_usage: number
}

// DTOs
export interface CreateFolderDTO {
  name: string
  description?: string
  parent_id?: string
  metadata?: string
}

export interface UpdateFolderDTO {
  name: string
  description?: string
  metadata?: string
}

export interface UploadDocumentDTO {
  name: string
  description?: string
  folder_id?: string
  tags: string[]
  linked_to?: DocumentLink
  enable_ocr: boolean
  metadata?: string
}

export interface UpdateDocumentDTO {
  name: string
  description?: string
  folder_id?: string
  tags: string[]
  metadata?: string
}

export interface DocumentFilters {
  folder_id?: string
  tags?: string[]
  mime_type?: string
  owner_id?: string
  search?: string
  date_from?: string
  date_to?: string
  page: number
  limit: number
  sort_by?: string
  sort_order?: string
}

export interface CreateDocumentPermissionDTO {
  subject_type: 'user' | 'role'
  subject_id: string
  object_type: 'document' | 'folder'
  object_id: string
  permission: 'view' | 'edit' | 'delete' | 'share'
  expires_at?: string
}

// API Functions

// Folders API
export const createFolder = async (data: CreateFolderDTO): Promise<Folder> => {
  const response = await apiClient.post('/folders', data)
  return response.data
}

export const getFolder = async (id: string): Promise<Folder> => {
  const response = await apiClient.get(`/folders/${id}`)
  return response.data
}

export const listFolders = async (parentId?: string): Promise<Folder[]> => {
  const params = parentId ? `?parent_id=${parentId}` : ''
  const response = await apiClient.get(`/folders${params}`)
  return Array.isArray(response.data) ? response.data : response.data.data || []
}

export const updateFolder = async (id: string, data: UpdateFolderDTO): Promise<void> => {
  await apiClient.put(`/folders/${id}`, data)
}

export const deleteFolder = async (id: string): Promise<void> => {
  await apiClient.delete(`/folders/${id}`)
}

// Documents API
export const uploadDocument = async (
  file: File,
  data: UploadDocumentDTO
): Promise<FileDocument> => {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('name', data.name)
  if (data.description) formData.append('description', data.description)
  if (data.folder_id) formData.append('folder_id', data.folder_id)
  if (data.tags.length > 0) formData.append('tags', data.tags.join(','))
  if (data.linked_to) formData.append('linked_to', JSON.stringify(data.linked_to))
  formData.append('enable_ocr', data.enable_ocr.toString())
  if (data.metadata) formData.append('metadata', data.metadata)

  const response = await apiClient.post('/documents/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  return response.data
}

export const getDocument = async (id: string): Promise<FileDocument> => {
  const response = await apiClient.get(`/documents/${id}`)
  return response.data
}

export const listDocuments = async (filters: DocumentFilters): Promise<FileDocument[]> => {
  const params = new URLSearchParams()
  
  if (filters.folder_id) params.append('folder_id', filters.folder_id)
  if (filters.mime_type) params.append('mime_type', filters.mime_type)
  if (filters.owner_id) params.append('owner_id', filters.owner_id)
  if (filters.search) params.append('search', filters.search)
  if (filters.date_from) params.append('date_from', filters.date_from)
  if (filters.date_to) params.append('date_to', filters.date_to)
  if (filters.sort_by) params.append('sort_by', filters.sort_by)
  if (filters.sort_order) params.append('sort_order', filters.sort_order)
  
  params.append('page', filters.page.toString())
  params.append('limit', filters.limit.toString())

  const response = await apiClient.get(`/documents?${params.toString()}`)
  return Array.isArray(response.data) ? response.data : response.data.data || []
}

export const listStructuredDocuments = async (): Promise<any> => {
  const response = await apiClient.get('/documents/structured')
  return response.data.data
}

export const updateDocument = async (id: string, data: UpdateDocumentDTO): Promise<void> => {
  await apiClient.put(`/documents/${id}`, data)
}

export const deleteDocument = async (id: string): Promise<void> => {
  await apiClient.delete(`/documents/${id}`)
}

export const downloadDocument = async (id: string): Promise<Blob> => {
  const response = await apiClient.get(`/documents/${id}/download`, {
    responseType: 'blob'
  })
  return response.data
}

// Search and Stats API
export const searchDocuments = async (query: string): Promise<DocumentSearchResult[]> => {
  const response = await apiClient.get(`/documents/search?q=${encodeURIComponent(query)}`)
  return response.data
}

export const getDocumentStats = async (): Promise<DocumentStats> => {
  const response = await apiClient.get('/documents/stats')
  return response.data
}

// Helper functions
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

export const getFileIcon = (mimeType: string): string => {
  if (mimeType.startsWith('image/')) return 'image'
  if (mimeType.startsWith('video/')) return 'video'
  if (mimeType.startsWith('audio/')) return 'audio'
  if (mimeType.includes('pdf')) return 'picture_as_pdf'
  if (mimeType.includes('word')) return 'description'
  if (mimeType.includes('excel') || mimeType.includes('spreadsheet')) return 'table_chart'
  if (mimeType.includes('powerpoint') || mimeType.includes('presentation')) return 'slideshow'
  if (mimeType.includes('text/')) return 'text_snippet'
  if (mimeType.includes('zip') || mimeType.includes('rar')) return 'folder_zip'
  return 'insert_drive_file'
}

export const getMimeTypeLabel = (mimeType: string): string => {
  const labels: Record<string, string> = {
    'application/pdf': 'PDF',
    'application/msword': 'Word',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document': 'Word',
    'application/vnd.ms-excel': 'Excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': 'Excel',
    'application/vnd.ms-powerpoint': 'PowerPoint',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation': 'PowerPoint',
    'text/plain': 'Текст',
    'text/csv': 'CSV',
    'image/jpeg': 'JPEG',
    'image/png': 'PNG',
    'image/gif': 'GIF',
    'video/mp4': 'MP4',
    'audio/mp3': 'MP3',
  }
  return labels[mimeType] || mimeType.split('/')[1]?.toUpperCase() || 'Файл'
}

