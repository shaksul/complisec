import { useState, useCallback } from 'react'
import { api } from '../shared/api/client'
import { Document } from '../shared/types/documents'

interface UploadDocumentData {
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

interface DocumentFilters {
  module?: string
  entityId?: string
  folderId?: string
  tags?: string[]
  mimeType?: string
  search?: string
}

export function useDocuments() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const uploadDocument = useCallback(async (
    file: File, 
    data: UploadDocumentData
  ): Promise<Document> => {
    try {
      setLoading(true)
      setError(null)

      const formData = new FormData()
      formData.append('file', file)
      formData.append('name', data.name)
      
      if (data.description) {
        formData.append('description', data.description)
      }
      
      if (data.tags) {
        formData.append('tags', JSON.stringify(data.tags))
      }
      
      if (data.linkedTo) {
        formData.append('linkedTo', JSON.stringify(data.linkedTo))
      }
      
      if (data.folderId) {
        formData.append('folderId', data.folderId)
      }
      
      if (data.metadata) {
        formData.append('metadata', JSON.stringify(data.metadata))
      }

      const response = await api.post('/documents/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })

      return response.data
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to upload document'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const getDocument = useCallback(async (documentId: string): Promise<Document> => {
    try {
      setLoading(true)
      setError(null)

      const response = await api.get(`/documents/${documentId}`)
      return response.data
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to get document'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const getModuleDocuments = useCallback(async (
    module: string, 
    entityId: string
  ): Promise<Document[]> => {
    try {
      setLoading(true)
      setError(null)

      const response = await api.get('/documents', {
        params: {
          module,
          entityId,
        },
      })

      return response.data.documents || []
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to get documents'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const listDocuments = useCallback(async (filters: DocumentFilters = {}): Promise<{
    documents: Document[]
    total: number
  }> => {
    try {
      setLoading(true)
      setError(null)

      const params = new URLSearchParams()
      
      if (filters.module) params.append('module', filters.module)
      if (filters.entityId) params.append('entityId', filters.entityId)
      if (filters.folderId) params.append('folderId', filters.folderId)
      if (filters.mimeType) params.append('mimeType', filters.mimeType)
      if (filters.search) params.append('search', filters.search)
      if (filters.tags) params.append('tags', JSON.stringify(filters.tags))

      const response = await api.get(`/documents?${params.toString()}`)
      return response.data
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to list documents'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const deleteDocument = useCallback(async (documentId: string): Promise<void> => {
    try {
      setLoading(true)
      setError(null)

      await api.delete(`/documents/${documentId}`)
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to delete document'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const linkDocument = useCallback(async (
    documentId: string, 
    module: string, 
    entityId: string
  ): Promise<void> => {
    try {
      setLoading(true)
      setError(null)

      await api.post(`/documents/${documentId}/links`, {
        module,
        entityId,
        linkType: 'attachment',
        description: `Linked to ${module}: ${entityId}`,
      })
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to link document'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const unlinkDocument = useCallback(async (
    documentId: string, 
    module: string, 
    entityId: string
  ): Promise<void> => {
    try {
      setLoading(true)
      setError(null)

      await api.delete(`/documents/${documentId}/links`, {
        data: {
          module,
          entityId,
        },
      })
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to unlink document'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const downloadDocument = useCallback(async (
    documentId: string, 
    fileName: string
  ): Promise<void> => {
    try {
      setLoading(true)
      setError(null)

      const response = await api.get(`/documents/${documentId}/download`, {
        responseType: 'blob',
      })

      // Создаем ссылку для скачивания
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', fileName)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to download document'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const previewDocument = useCallback(async (documentId: string): Promise<void> => {
    try {
      setLoading(true)
      setError(null)

      // Открываем документ в новом окне для предварительного просмотра
      const previewUrl = `/api/documents/${documentId}/preview`
      window.open(previewUrl, '_blank')
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to preview document'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const searchDocuments = useCallback(async (query: string): Promise<Document[]> => {
    try {
      setLoading(true)
      setError(null)

      const response = await api.get('/documents/search', {
        params: { q: query },
      })

      return response.data.documents || []
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to search documents'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const getDocumentStats = useCallback(async (): Promise<{
    totalDocuments: number
    totalFolders: number
    totalSize: number
    documentsByType: Record<string, number>
  }> => {
    try {
      setLoading(true)
      setError(null)

      const response = await api.get('/documents/stats')
      return response.data
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to get document stats'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  // Методы для работы с папками
  const createFolder = useCallback(async (data: {
    name: string
    description?: string
    parentId?: string
    metadata?: Record<string, any>
  }) => {
    try {
      setLoading(true)
      setError(null)

      const response = await api.post('/folders', data)
      return response.data
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to create folder'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  const getFolders = useCallback(async (parentId?: string) => {
    try {
      setLoading(true)
      setError(null)

      const params = parentId ? { parentId } : {}
      const response = await api.get('/folders', { params })
      return response.data.folders || []
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to get folders'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    // State
    loading,
    error,
    
    // Document operations
    uploadDocument,
    getDocument,
    getModuleDocuments,
    listDocuments,
    deleteDocument,
    linkDocument,
    unlinkDocument,
    downloadDocument,
    previewDocument,
    searchDocuments,
    getDocumentStats,
    
    // Folder operations
    createFolder,
    getFolders,
    
    // Utility
    clearError: () => setError(null),
  }
}
