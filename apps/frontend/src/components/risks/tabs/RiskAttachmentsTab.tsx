import React, { useEffect, useState, useRef } from 'react'
import {
  Box,
  Typography,
  Card,
  CardContent,
  IconButton,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  ListItemSecondaryAction,
  Chip,
  CircularProgress,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  LinearProgress,
} from '@mui/material'
import {
  Delete,
  InsertDriveFile,
  Image,
  PictureAsPdf,
  Description,
  VideoFile,
  AudioFile,
  Archive,
  Add,
  Upload,
} from '@mui/icons-material'
import { risksApi, RiskAttachment } from '../../../shared/api/risks'

interface RiskAttachmentsTabProps {
  riskId: string
}

const getFileIcon = (mimeType: string) => {
  if (mimeType.startsWith('image/')) return <Image />
  if (mimeType.includes('pdf')) return <PictureAsPdf />
  if (mimeType.includes('video')) return <VideoFile />
  if (mimeType.includes('audio')) return <AudioFile />
  if (mimeType.includes('zip') || mimeType.includes('rar')) return <Archive />
  if (mimeType.includes('document') || mimeType.includes('text')) return <Description />
  return <InsertDriveFile />
}

const formatFileSize = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

export const RiskAttachmentsTab: React.FC<RiskAttachmentsTabProps> = ({ riskId }) => {
  const [attachments, setAttachments] = useState<RiskAttachment[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState(0)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [fileDescription, setFileDescription] = useState('')
  const fileInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    void loadAttachments()
  }, [riskId])

  const loadAttachments = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await risksApi.getAttachments(riskId)
      setAttachments(response)
    } catch (err) {
      console.error('Error loading attachments:', err)
      setError('Не удалось загрузить вложения')
      setAttachments([])
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (attachmentId: string) => {
    try {
      await risksApi.deleteAttachment(riskId, attachmentId)
      await loadAttachments()
    } catch (err) {
      console.error('Error deleting attachment:', err)
      setError('Не удалось удалить вложение')
    }
  }

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      setSelectedFile(file)
      setFileDescription('')
    }
  }

  const handleUpload = async () => {
    if (!selectedFile) return

    try {
      setUploading(true)
      setUploadProgress(0)
      setError(null)

      // Загружаем файл через API рисков
      await risksApi.createAttachment(riskId, selectedFile, fileDescription)
      
      setUploadProgress(100)

      // Обновляем список вложений
      await loadAttachments()
      
      // Закрываем диалог и сбрасываем состояние
      setUploadDialogOpen(false)
      setSelectedFile(null)
      setFileDescription('')
      setUploadProgress(0)
    } catch (err) {
      console.error('Error uploading file:', err)
      setError('Не удалось загрузить файл')
    } finally {
      setUploading(false)
    }
  }

  const handleCancelUpload = () => {
    setUploadDialogOpen(false)
    setSelectedFile(null)
    setFileDescription('')
    setUploadProgress(0)
    setError(null)
  }

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
        <Typography variant="h6">
          Вложения ({attachments.length})
        </Typography>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={() => setUploadDialogOpen(true)}
          disabled={uploading}
        >
          Добавить файл
        </Button>
      </Box>

      {error && (
        <Box mb={2}>
          <Typography color="error">{error}</Typography>
        </Box>
      )}

      {attachments.length === 0 ? (
        <Card>
          <CardContent>
            <Box textAlign="center" py={4}>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                Вложения отсутствуют
              </Typography>
              <Typography variant="body2" color="text.secondary" mb={2}>
                Нажмите "Добавить файл" чтобы загрузить документы для этого риска.
              </Typography>
              <Button
                variant="outlined"
                startIcon={<Upload />}
                onClick={() => setUploadDialogOpen(true)}
                disabled={uploading}
              >
                Загрузить файл
              </Button>
            </Box>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardContent>
            <List>
              {attachments.map((attachment) => (
                <ListItem key={attachment.id} divider>
                  <ListItemIcon>{getFileIcon(attachment.mime_type)}</ListItemIcon>
                  <ListItemText
                    primary={attachment.original_name}
                    secondary={
                      <Box display="flex" alignItems="center" gap={1} flexWrap="wrap">
                        <Typography variant="body2" color="text.secondary">
                          {formatFileSize(attachment.file_size)}
                        </Typography>
                        <Chip label={attachment.mime_type} size="small" variant="outlined" />
                        <Typography variant="body2" color="text.secondary">
                          Загрузил: {attachment.created_by}
                        </Typography>
                        {attachment.description && (
                          <Typography variant="body2" color="text.secondary">
                            • {attachment.description}
                          </Typography>
                        )}
                      </Box>
                    }
                  />
                  <ListItemSecondaryAction>
                    <IconButton edge="end" color="error" onClick={() => handleDelete(attachment.id)}>
                      <Delete />
                    </IconButton>
                  </ListItemSecondaryAction>
                </ListItem>
              ))}
            </List>
          </CardContent>
        </Card>
      )}

      {/* Диалог загрузки файла */}
      <Dialog open={uploadDialogOpen} onClose={handleCancelUpload} maxWidth="sm" fullWidth>
        <DialogTitle>Загрузить файл</DialogTitle>
        <DialogContent>
          <Box py={2}>
            <input
              ref={fileInputRef}
              type="file"
              onChange={handleFileSelect}
              style={{ display: 'none' }}
              accept="*/*"
            />
            
            <Button
              variant="outlined"
              startIcon={<Upload />}
              onClick={() => fileInputRef.current?.click()}
              fullWidth
              sx={{ mb: 2 }}
            >
              {selectedFile ? 'Выбрать другой файл' : 'Выбрать файл'}
            </Button>

            {selectedFile && (
              <Box mb={2}>
                <Typography variant="body2" color="text.secondary">
                  Выбранный файл: {selectedFile.name} ({formatFileSize(selectedFile.size)})
                </Typography>
              </Box>
            )}

            <TextField
              fullWidth
              label="Описание файла (необязательно)"
              value={fileDescription}
              onChange={(e) => setFileDescription(e.target.value)}
              multiline
              rows={2}
              disabled={uploading}
            />

            {uploading && (
              <Box mt={2}>
                <LinearProgress variant="determinate" value={uploadProgress} />
                <Typography variant="body2" color="text.secondary" mt={1}>
                  Загрузка... {uploadProgress}%
                </Typography>
              </Box>
            )}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCancelUpload} disabled={uploading}>
            Отмена
          </Button>
          <Button
            onClick={handleUpload}
            variant="contained"
            disabled={!selectedFile || uploading}
          >
            {uploading ? 'Загрузка...' : 'Загрузить'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}


