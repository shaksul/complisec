export interface Document {
  id: string
  tenant_id: string
  name: string
  original_name?: string
  description?: string
  file_path?: string
  file_size?: number
  mime_type?: string
  file_hash?: string
  folder_id?: string
  owner_id?: string
  created_by: string
  created_at: string
  updated_at: string
  is_active: boolean
  version: string
  metadata?: string
  tags?: string[]
  links?: DocumentLink[]
  ocr_text?: string
}

export interface DocumentLink {
  id: string
  document_id: string
  module: string
  entity_id: string
  link_type: string
  description?: string
  created_by: string
  created_at: string
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
  document_count?: number
}

export interface UploadDocumentRequest {
  name: string
  description?: string
  tags?: string[]
  linkedTo?: {
    module: string
    entityId: string
  }
  folderId?: string
  metadata?: Record<string, any>
}

export interface DocumentFilters {
  module?: string
  entityId?: string
  folderId?: string
  tags?: string[]
  mimeType?: string
  search?: string
  ownerId?: string
  createdBy?: string
  dateFrom?: string
  dateTo?: string
  page?: number
  pageSize?: number
}

export interface DocumentListResponse {
  documents: Document[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

export interface DocumentStats {
  totalDocuments: number
  totalFolders: number
  totalSize: number
  documentsByType: Record<string, number>
  recentDocuments: Document[]
  storageUsage: number
}

export interface CreateFolderRequest {
  name: string
  description?: string
  parentId?: string
  metadata?: Record<string, any>
}

export interface UpdateFolderRequest {
  name?: string
  description?: string
  metadata?: Record<string, any>
}

export interface CreateDocumentLinkRequest {
  documentId: string
  module: string
  entityId: string
  linkType: string
  description?: string
}

export interface DocumentSearchResult {
  documents: Document[]
  total: number
  query: string
  filters?: DocumentFilters
}

// Типы для различных модулей
export type DocumentModule = 'assets' | 'risks' | 'incidents' | 'training' | 'compliance'

export interface AssetDocumentRequest {
  documentType: string
  title?: string
  description?: string
}

export interface RiskDocumentRequest {
  documentType: string
  title?: string
  description?: string
  isInternal?: boolean
}

export interface IncidentDocumentRequest {
  documentType: string
  title?: string
  description?: string
}

export interface TrainingDocumentRequest {
  documentType: string
  title?: string
  description?: string
  materialType?: string
}

// Типы для файлов
export interface FileUploadProgress {
  loaded: number
  total: number
  percentage: number
}

export interface FileUploadStatus {
  status: 'idle' | 'uploading' | 'success' | 'error'
  progress?: FileUploadProgress
  error?: string
}

// Типы для OCR
export interface OCRResult {
  text: string
  confidence: number
  language: string
  processingTime: number
}

// Типы для версионирования
export interface DocumentVersion {
  id: string
  document_id: string
  version: string
  file_path: string
  file_size: number
  mime_type: string
  file_hash: string
  created_by: string
  created_at: string
  change_notes?: string
}

export interface CreateVersionRequest {
  changeNotes?: string
}

// Типы для комментариев
export interface DocumentComment {
  id: string
  document_id: string
  user_id: string
  user_name: string
  comment: string
  is_internal: boolean
  created_at: string
  updated_at: string
}

export interface CreateCommentRequest {
  comment: string
  isInternal?: boolean
}

// Типы для разрешений
export interface DocumentPermission {
  id: string
  document_id: string
  user_id?: string
  role_id?: string
  permission: 'read' | 'write' | 'delete' | 'share'
  granted_by: string
  granted_at: string
  expires_at?: string
}

export interface GrantPermissionRequest {
  userId?: string
  roleId?: string
  permission: 'read' | 'write' | 'delete' | 'share'
  expiresAt?: string
}
