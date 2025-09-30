import React, { useState } from 'react'
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
  FormControlLabel,
  Switch,
  LinearProgress,
} from '@mui/material'
import { CloudUpload, Description } from '@mui/icons-material'
import { uploadDocumentVersion, type CreateDocumentVersionDTO } from '../../shared/api/documents'

interface UploadNewVersionDialogProps {
  open: boolean
  onClose: () => void
  onSuccess: () => void
  documentId: string
  documentTitle: string
}

export default function UploadNewVersionDialog({ 
  open, 
  onClose, 
  onSuccess, 
  documentId, 
  documentTitle 
}: UploadNewVersionDialogProps) {
  const [file, setFile] = useState<File | null>(null)
  const [enableOCR, setEnableOCR] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [uploadProgress, setUploadProgress] = useState(0)

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = event.target.files?.[0]
    if (selectedFile) {
      setFile(selectedFile)
      setError(null)
    }
  }

  const handleSubmit = async () => {
    if (!file) {
      setError('Пожалуйста, выберите файл')
      return
    }

    try {
      setLoading(true)
      setError(null)
      setUploadProgress(0)

      const options: CreateDocumentVersionDTO = {
        title: documentTitle,
        enableOCR,
      }

      // Симуляция прогресса загрузки
      const progressInterval = setInterval(() => {
        setUploadProgress(prev => {
          if (prev >= 90) {
            clearInterval(progressInterval)
            return 90
          }
          return prev + 10
        })
      }, 200)

      await uploadDocumentVersion(documentId, file, options)
      
      clearInterval(progressInterval)
      setUploadProgress(100)
      
      setTimeout(() => {
        setFile(null)
        setEnableOCR(false)
        setUploadProgress(0)
        onSuccess()
        onClose()
      }, 500)

    } catch (err) {
      console.error('Error uploading new version:', err)
      setError('Ошибка загрузки новой версии: ' + (err as Error).message)
    } finally {
      setLoading(false)
    }
  }

  const handleClose = () => {
    if (!loading) {
      setFile(null)
      setEnableOCR(false)
      setUploadProgress(0)
      setError(null)
      onClose()
    }
  }

  const formatFileSize = (bytes: number): string => {
    const sizes = ['Б', 'КБ', 'МБ', 'ГБ']
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Box display="flex" alignItems="center" gap={1}>
          <CloudUpload />
          Загрузить новую версию: {documentTitle}
        </Box>
      </DialogTitle>
      
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
          <Box>
            <input
              accept=".pdf,.doc,.docx,.txt,.jpg,.jpeg,.png"
              style={{ display: 'none' }}
              id="upload-new-version"
              type="file"
              onChange={handleFileChange}
              disabled={loading}
            />
            <label htmlFor="upload-new-version">
              <Button
                variant="outlined"
                component="span"
                startIcon={<Description />}
                disabled={loading}
                fullWidth
                sx={{ py: 2 }}
              >
                {file ? `Выбран файл: ${file.name}` : 'Выберите файл'}
              </Button>
            </label>
          </Box>

          {file && (
            <Box sx={{ p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
              <Typography variant="body2" color="text.secondary">
                <strong>Имя файла:</strong> {file.name}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                <strong>Размер:</strong> {formatFileSize(file.size)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                <strong>Тип:</strong> {file.type || 'Неизвестно'}
              </Typography>
            </Box>
          )}

          <FormControlLabel
            control={
              <Switch
                checked={enableOCR}
                onChange={(e) => setEnableOCR(e.target.checked)}
                disabled={loading}
              />
            }
            label="Включить OCR для извлечения текста"
          />

          {loading && (
            <Box>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Загрузка файла...
              </Typography>
              <LinearProgress variant="determinate" value={uploadProgress} />
              <Typography variant="caption" color="text.secondary">
                {uploadProgress}%
              </Typography>
            </Box>
          )}
        </Box>
      </DialogContent>
      
      <DialogActions>
        <Button onClick={handleClose} disabled={loading}>
          Отмена
        </Button>
        <Button 
          onClick={handleSubmit} 
          variant="contained" 
          disabled={!file || loading}
          startIcon={loading ? <CircularProgress size={20} /> : <CloudUpload />}
        >
          {loading ? 'Загрузка...' : 'Загрузить версию'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
