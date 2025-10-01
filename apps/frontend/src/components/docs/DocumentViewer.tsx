import React, { useState, useEffect } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography,
  Alert,
  CircularProgress,
  Paper,
} from '@mui/material'
import {
  Download,
  OpenInNew,
  Close,
  TextSnippet,
  PictureAsPdf,
  Description as WordIcon,
} from '@mui/icons-material'

interface DocumentViewerProps {
  open: boolean
  onClose: () => void
  versionId: string
  fileName: string
  mimeType?: string
}

export default function DocumentViewer({ 
  open, 
  onClose, 
  versionId, 
  fileName, 
  mimeType 
}: DocumentViewerProps) {
  const [loading, setLoading] = useState(false)
  const [fileContent, setFileContent] = useState<string>('')
  const [error, setError] = useState<string | null>(null)

  const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'
  const token = localStorage.getItem('access_token')
  const query = new URLSearchParams({ preview: 'true' })
  if (token) query.append('access_token', token)
  const downloadUrl = `${apiBaseUrl}/documents/versions/${versionId}/download?${query.toString()}`

  const isOfficeDoc = mimeType?.toLowerCase().includes('msword') || 
                     mimeType?.toLowerCase().includes('openxmlformats') ||
                     mimeType?.toLowerCase().includes('application/vnd.ms')

  const isPdf = mimeType?.toLowerCase().includes('pdf')
  const isText = mimeType?.toLowerCase().includes('text/plain')

  useEffect(() => {
    if (open && (isText || isOfficeDoc)) {
      loadFileContent()
    }
  }, [open, versionId, isText, isOfficeDoc])

  const loadFileContent = async () => {
    try {
      setLoading(true)
      setError(null)
      
      if (isText) {
        const response = await fetch(downloadUrl, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        
        if (!response.ok) {
          throw new Error('Ошибка загрузки файла')
        }
        
        // Исправляем кодировку для текстовых файлов
        const arrayBuffer = await response.arrayBuffer()
        const decoder = new TextDecoder('utf-8')
        const text = decoder.decode(arrayBuffer)
        setFileContent(text)
      } else if (isOfficeDoc) {
        // For Office documents, get HTML version
        const htmlUrl = `${apiBaseUrl}/documents/versions/${versionId}/html?${query.toString()}`
        const response = await fetch(htmlUrl, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        
        if (!response.ok) {
          throw new Error('Ошибка загрузки HTML версии файла')
        }
        
        const html = await response.text()
        setFileContent(html)
      }
    } catch (err) {
      setError('Ошибка загрузки содержимого файла: ' + (err as Error).message)
    } finally {
      setLoading(false)
    }
  }

  const handleDownload = () => {
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  const handleOpenInNewTab = () => {
    window.open(downloadUrl, '_blank')
  }

  const renderViewer = () => {
    if (loading) {
      return (
        <Box display="flex" justifyContent="center" alignItems="center" height="400px">
          <CircularProgress />
        </Box>
      )
    }

    if (error) {
      return (
        <Box sx={{ p: 3, textAlign: 'center' }}>
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
          <Button
            variant="contained"
            startIcon={<Download />}
            onClick={handleDownload}
          >
            Скачать файл
          </Button>
        </Box>
      )
    }

    if (isPdf) {
      return (
        <Box sx={{ width: '100%', height: '80vh' }}>
          <iframe
            src={downloadUrl}
            width="100%"
            height="100%"
            style={{ border: 'none' }}
            title={fileName}
          />
        </Box>
      )
    }

    if (isText) {
      return (
        <Box sx={{ width: '100%', height: '80vh', p: 2 }}>
          <Paper sx={{ p: 2, height: '100%', overflow: 'auto' }}>
            <Typography variant="h6" gutterBottom sx={{ mb: 2 }}>
              <TextSnippet sx={{ mr: 1, verticalAlign: 'middle' }} />
              Содержимое файла
            </Typography>
            <Box
              component="pre"
              sx={{
                whiteSpace: 'pre-wrap',
                wordWrap: 'break-word',
                fontFamily: 'monospace',
                fontSize: '0.875rem',
                lineHeight: 1.5,
                margin: 0,
                padding: 0,
              }}
            >
              {fileContent}
            </Box>
          </Paper>
        </Box>
      )
    }

    if (isOfficeDoc) {
      return (
        <Box sx={{ width: '100%', height: '80vh', p: 3 }}>
          <Paper sx={{ p: 3, height: '100%', textAlign: 'center', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
            <WordIcon sx={{ fontSize: 80, color: 'primary.main', mb: 3 }} />
            <Typography variant="h5" gutterBottom>
              Документ Microsoft Office
            </Typography>
            <Typography variant="h6" color="text.secondary" gutterBottom sx={{ mb: 2 }}>
              {fileName}
            </Typography>
            <Typography variant="body1" color="text.secondary" sx={{ mb: 4, maxWidth: 600, mx: 'auto' }}>
              Этот документ не может быть отображен в браузере. Скачайте файл и откройте его в Microsoft Office или совместимом приложении для просмотра содержимого.
            </Typography>
            <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center', flexWrap: 'wrap' }}>
              <Button
                variant="contained"
                size="large"
                startIcon={<Download />}
                onClick={handleDownload}
                sx={{ minWidth: 200 }}
              >
                Скачать файл
              </Button>
              <Button
                variant="outlined"
                size="large"
                startIcon={<OpenInNew />}
                onClick={handleOpenInNewTab}
                sx={{ minWidth: 200 }}
              >
                Открыть в новой вкладке
              </Button>
            </Box>
          </Paper>
        </Box>
      )
    }

    return (
      <Box sx={{ p: 3, textAlign: 'center' }}>
        <Typography variant="h6" gutterBottom>
          Предварительный просмотр недоступен
        </Typography>
        <Typography variant="body2" color="text.secondary" gutterBottom>
          Этот тип файла не поддерживает встроенный просмотр.
        </Typography>
        <Button
          variant="contained"
          startIcon={<Download />}
          onClick={handleDownload}
          sx={{ mt: 2 }}
        >
          Скачать файл
        </Button>
      </Box>
    )
  }

  return (
    <Dialog open={open} onClose={onClose} maxWidth="lg" fullWidth>
      <DialogTitle>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Box display="flex" alignItems="center" gap={1}>
            {isPdf && <PictureAsPdf color="error" />}
            {isText && <TextSnippet color="info" />}
            {isOfficeDoc && <WordIcon color="primary" />}
            <Typography variant="h6">
              {fileName}
            </Typography>
          </Box>
          <Button
            onClick={onClose}
            size="small"
            startIcon={<Close />}
          >
            Закрыть
          </Button>
        </Box>
      </DialogTitle>
      
      <DialogContent sx={{ p: 0 }}>
        {renderViewer()}
      </DialogContent>
      
      <DialogActions>
        <Button
          startIcon={<OpenInNew />}
          onClick={handleOpenInNewTab}
        >
          Открыть в новой вкладке
        </Button>
        <Button
          startIcon={<Download />}
          onClick={handleDownload}
          variant="outlined"
        >
          Скачать
        </Button>
      </DialogActions>
    </Dialog>
  )
}
