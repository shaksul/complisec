import React, { useState, useEffect } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Typography,
  Box,
  Chip,
  IconButton,
  Tooltip,
  CircularProgress,
  Alert,
} from '@mui/material'
import {
  Download,
  Visibility,
  CloudDownload,
  Security,
  Description,
  CloudUpload,
  PictureAsPdf,
  Description as WordIcon,
  TextSnippet,
  Image,
} from '@mui/icons-material'
import { getDocumentVersions, downloadDocumentVersion, getDocumentVersionPreview, type DocumentVersion } from '../../shared/api/documents'
import UploadNewVersionDialog from './UploadNewVersionDialog'
import DocumentViewer from './DocumentViewer'

interface DocumentVersionsDialogProps {
  open: boolean
  onClose: () => void
  documentId: string
  documentTitle: string
}

export default function DocumentVersionsDialog({ 
  open, 
  onClose, 
  documentId, 
  documentTitle 
}: DocumentVersionsDialogProps) {
  const [versions, setVersions] = useState<DocumentVersion[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [openUploadNewVersion, setOpenUploadNewVersion] = useState(false)
  const [viewingVersion, setViewingVersion] = useState<DocumentVersion | null>(null)

  useEffect(() => {
    if (open && documentId) {
      loadVersions()
    }
  }, [open, documentId])

  const loadVersions = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await getDocumentVersions(documentId)
      setVersions(data)
    } catch (err) {
      console.error('Error loading document versions:', err)
      setError('Ошибка загрузки версий документа')
    } finally {
      setLoading(false)
    }
  }

  const formatFileSize = (bytes?: number): string => {
    if (!bytes) return '-'
    const sizes = ['Б', 'КБ', 'МБ', 'ГБ']
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`
  }

  const getScanStatusColor = (status: string): 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning' => {
    const colors: Record<string, 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning'> = {
      pending: 'warning',
      clean: 'success',
      infected: 'error',
      error: 'error'
    }
    return colors[status] || 'default'
  }

  const getScanStatusLabel = (status: string): string => {
    const labels: Record<string, string> = {
      pending: 'Проверяется',
      clean: 'Безопасен',
      infected: 'Заражен',
      error: 'Ошибка проверки'
    }
    return labels[status] || status
  }

  const handleDownload = async (version: DocumentVersion) => {
    try {
      setLoading(true)
      const blob = await downloadDocumentVersion(version.id)
      
      // Создаем ссылку для скачивания
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      
      // Определяем имя файла на основе версии
      const fileName = `${documentTitle}_v${version.version_number}.${getFileExtension(version.mime_type)}`
      link.download = fileName
      
      // Добавляем ссылку в DOM, кликаем и удаляем
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      
      // Освобождаем память
      window.URL.revokeObjectURL(url)
    } catch (err) {
      console.error('Error downloading file:', err)
      setError('Ошибка скачивания файла')
    } finally {
      setLoading(false)
    }
  }

  const handlePreview = (version: DocumentVersion) => {
    setViewingVersion(version)
  }

  const getFileExtension = (mimeType?: string): string => {
    if (!mimeType) return 'bin'
    
    const extensions: Record<string, string> = {
      'application/pdf': 'pdf',
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document': 'docx',
      'application/msword': 'doc',
      'text/plain': 'txt',
      'image/jpeg': 'jpg',
      'image/png': 'png',
    }
    
    return extensions[mimeType] || 'bin'
  }

  const getFileIcon = (mimeType?: string) => {
    if (!mimeType) return <Description fontSize="small" />
    
    const mimeTypeLower = mimeType.toLowerCase()
    if (mimeTypeLower.includes('pdf')) return <PictureAsPdf fontSize="small" color="error" />
    if (mimeTypeLower.includes('msword') || mimeTypeLower.includes('openxmlformats')) {
      return <WordIcon fontSize="small" color="primary" />
    }
    if (mimeTypeLower.includes('text/plain')) return <TextSnippet fontSize="small" color="info" />
    if (mimeTypeLower.includes('image/')) return <Image fontSize="small" color="secondary" />
    
    return <Description fontSize="small" />
  }

  return (
    <Dialog open={open} onClose={onClose} maxWidth="lg" fullWidth>
      <DialogTitle>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Box display="flex" alignItems="center" gap={1}>
            <Description />
            Версии документа: {documentTitle}
          </Box>
          <Button
            variant="outlined"
            startIcon={<CloudUpload />}
            onClick={() => setOpenUploadNewVersion(true)}
            size="small"
          >
            Новая версия
          </Button>
        </Box>
      </DialogTitle>
      
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {loading ? (
          <Box display="flex" justifyContent="center" p={3}>
            <CircularProgress />
          </Box>
        ) : versions.length === 0 ? (
          <Box textAlign="center" p={3}>
            <Typography variant="body1" color="text.secondary">
              Нет версий документа
            </Typography>
          </Box>
        ) : (
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Версия</TableCell>
                  <TableCell>Размер</TableCell>
                  <TableCell>Тип файла</TableCell>
                  <TableCell>Проверка безопасности</TableCell>
                  <TableCell>Создано</TableCell>
                  <TableCell>Создал</TableCell>
                  <TableCell>Действия</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {versions.map((version) => (
                  <TableRow key={version.id}>
                    <TableCell>
                      <Typography variant="subtitle2">
                        v{version.version_number}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      {formatFileSize(version.size_bytes)}
                    </TableCell>
                    <TableCell>
                      <Box display="flex" alignItems="center" gap={1}>
                        {getFileIcon(version.mime_type)}
                        <Typography variant="body2">
                          {version.mime_type ? getFileExtension(version.mime_type).toUpperCase() : 'Неизвестно'}
                        </Typography>
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={getScanStatusLabel(version.av_scan_status)}
                        color={getScanStatusColor(version.av_scan_status)}
                        size="small"
                        icon={<Security />}
                      />
                    </TableCell>
                    <TableCell>
                      {new Date(version.created_at).toLocaleString()}
                    </TableCell>
                    <TableCell>
                      {version.created_by}
                    </TableCell>
                    <TableCell>
                      <Box display="flex" gap={1}>
                        <Tooltip title={
                          version.mime_type?.toLowerCase().includes('msword') || 
                          version.mime_type?.toLowerCase().includes('openxmlformats') 
                            ? "Предварительный просмотр" 
                            : "Предварительный просмотр"
                        }>
                          <IconButton 
                            size="small" 
                            onClick={() => handlePreview(version)}
                          >
                            <Visibility />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Скачать">
                          <IconButton 
                            size="small" 
                            onClick={() => handleDownload(version)}
                          >
                            <CloudDownload />
                          </IconButton>
                        </Tooltip>
                      </Box>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </DialogContent>
      
      <DialogActions>
        <Button onClick={onClose}>
          Закрыть
        </Button>
      </DialogActions>

      {/* Upload New Version Dialog */}
      <UploadNewVersionDialog
        open={openUploadNewVersion}
        onClose={() => setOpenUploadNewVersion(false)}
        onSuccess={loadVersions}
        documentId={documentId}
        documentTitle={documentTitle}
      />

      {/* Document Viewer */}
      {viewingVersion && (
        <DocumentViewer
          open={!!viewingVersion}
          onClose={() => setViewingVersion(null)}
          versionId={viewingVersion.id}
          fileName={`${documentTitle}_v${viewingVersion.version_number}.${getFileExtension(viewingVersion.mime_type)}`}
          mimeType={viewingVersion.mime_type}
        />
      )}
    </Dialog>
  )
}
