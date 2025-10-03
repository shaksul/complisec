import React, { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Button,
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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from '@mui/material'
import { 
  Add, 
  Download, 
  Delete, 
  InsertDriveFile,
  Image,
  Description,
  PictureAsPdf,
  VideoFile,
  AudioFile,
  Archive,
} from '@mui/icons-material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

const attachmentSchema = z.object({
  name: z.string().min(1, 'Название обязательно'),
  description: z.string().optional(),
})

type AttachmentFormData = z.infer<typeof attachmentSchema>

interface RiskAttachment {
  id: string
  risk_id: string
  name: string
  description?: string
  filename: string
  file_size: number
  file_type: string
  uploaded_by: string
  uploaded_by_name: string
  created_at: string
}

interface RiskAttachmentsTabProps {
  riskId: string
}

export const RiskAttachmentsTab: React.FC<RiskAttachmentsTabProps> = ({ riskId }) => {
  const [attachments, setAttachments] = useState<RiskAttachment[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<AttachmentFormData>({
    resolver: zodResolver(attachmentSchema),
    defaultValues: {
      name: '',
      description: '',
    },
  })

  useEffect(() => {
    loadAttachments()
  }, [riskId])

  const loadAttachments = async () => {
    try {
      setLoading(true)
      // TODO: Implement API call to load attachments
      // const response = await riskAttachmentsApi.list(riskId)
      // setAttachments(response.data || [])
      
      // Mock data for now
      setAttachments([
        {
          id: '1',
          risk_id: riskId,
          name: 'Анализ риска',
          description: 'Подробный анализ риска безопасности',
          filename: 'risk_analysis.pdf',
          file_size: 1024000, // 1MB
          file_type: 'application/pdf',
          uploaded_by: 'user1',
          uploaded_by_name: 'Иван Петров',
          created_at: new Date(Date.now() - 86400000).toISOString(),
        },
        {
          id: '2',
          risk_id: riskId,
          name: 'Схема архитектуры',
          description: 'Диаграмма архитектуры системы',
          filename: 'architecture_diagram.png',
          file_size: 512000, // 512KB
          file_type: 'image/png',
          uploaded_by: 'user2',
          uploaded_by_name: 'Мария Сидорова',
          created_at: new Date(Date.now() - 172800000).toISOString(),
        },
        {
          id: '3',
          risk_id: riskId,
          name: 'Отчет по аудиту',
          description: 'Результаты аудита безопасности',
          filename: 'audit_report.docx',
          file_size: 2048000, // 2MB
          file_type: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
          uploaded_by: 'user3',
          uploaded_by_name: 'Алексей Козлов',
          created_at: new Date(Date.now() - 259200000).toISOString(),
        },
      ])
    } catch (err) {
      console.error('Error loading attachments:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateNew = () => {
    reset({ name: '', description: '' })
    setModalOpen(true)
  }

  const onSubmit = async (data: AttachmentFormData) => {
    try {
      // TODO: Implement file upload and API call to create attachment
      // const formData = new FormData()
      // formData.append('file', file)
      // formData.append('name', data.name)
      // formData.append('description', data.description || '')
      // formData.append('risk_id', riskId)
      // const response = await riskAttachmentsApi.create(formData)
      
      console.log('Upload attachment:', data)
      setModalOpen(false)
    } catch (err) {
      console.error('Error uploading attachment:', err)
    }
  }

  const handleDownload = async (attachment: RiskAttachment) => {
    try {
      // TODO: Implement API call to download attachment
      // await riskAttachmentsApi.download(attachment.id)
      console.log('Download attachment:', attachment.filename)
    } catch (err) {
      console.error('Error downloading attachment:', err)
    }
  }

  const handleDelete = async (attachmentId: string) => {
    try {
      // TODO: Implement API call to delete attachment
      // await riskAttachmentsApi.delete(attachmentId)
      setAttachments(prev => prev.filter(a => a.id !== attachmentId))
    } catch (err) {
      console.error('Error deleting attachment:', err)
    }
  }

  const getFileIcon = (fileType: string) => {
    if (fileType.startsWith('image/')) return <Image />
    if (fileType.includes('pdf')) return <PictureAsPdf />
    if (fileType.includes('video')) return <VideoFile />
    if (fileType.includes('audio')) return <AudioFile />
    if (fileType.includes('zip') || fileType.includes('rar')) return <Archive />
    if (fileType.includes('document') || fileType.includes('text')) return <Description />
    return <InsertDriveFile />
  }

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
    
    if (diffDays === 0) {
      return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
    } else if (diffDays === 1) {
      return 'Вчера'
    } else if (diffDays < 7) {
      return `${diffDays} дней назад`
    } else {
      return date.toLocaleDateString('ru-RU')
    }
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
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h6">
          Вложения ({attachments.length})
        </Typography>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={handleCreateNew}
        >
          Добавить вложение
        </Button>
      </Box>

      {attachments.length === 0 ? (
        <Card>
          <CardContent>
            <Box textAlign="center" py={4}>
              <InsertDriveFile sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" color="text.secondary" gutterBottom>
                Вложения не найдены
              </Typography>
              <Typography variant="body2" color="text.secondary" mb={3}>
                Добавьте файлы для документирования риска
              </Typography>
              <Button
                variant="contained"
                startIcon={<Add />}
                onClick={handleCreateNew}
              >
                Добавить вложение
              </Button>
            </Box>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <List>
            {attachments.map((attachment, index) => (
              <ListItem key={attachment.id} divider={index < attachments.length - 1}>
                <ListItemIcon>
                  {getFileIcon(attachment.file_type)}
                </ListItemIcon>
                
                <ListItemText
                  primary={
                    <Box>
                      <Typography variant="subtitle1" fontWeight="bold">
                        {attachment.name}
                      </Typography>
                      {attachment.description && (
                        <Typography variant="body2" color="text.secondary">
                          {attachment.description}
                        </Typography>
                      )}
                    </Box>
                  }
                  secondary={
                    <Box display="flex" alignItems="center" gap={1} mt={1}>
                      <Typography variant="caption" color="text.secondary">
                        {attachment.uploaded_by_name}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        •
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {formatDate(attachment.created_at)}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        •
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {formatFileSize(attachment.file_size)}
                      </Typography>
                      <Chip
                        label={attachment.filename.split('.').pop()?.toUpperCase() || 'FILE'}
                        size="small"
                        variant="outlined"
                      />
                    </Box>
                  }
                />
                
                <ListItemSecondaryAction>
                  <Box display="flex" gap={1}>
                    <IconButton
                      size="small"
                      onClick={() => handleDownload(attachment)}
                      title="Скачать"
                    >
                      <Download />
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => handleDelete(attachment.id)}
                      color="error"
                      title="Удалить"
                    >
                      <Delete />
                    </IconButton>
                  </Box>
                </ListItemSecondaryAction>
              </ListItem>
            ))}
          </List>
        </Card>
      )}

      {/* Add Attachment Modal */}
      <Dialog open={modalOpen} onClose={() => setModalOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          Добавление вложения
        </DialogTitle>
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogContent>
            <Box mb={2}>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Выберите файл для загрузки
              </Typography>
              <Button
                variant="outlined"
                component="label"
                fullWidth
                sx={{ py: 2 }}
              >
                <Add sx={{ mr: 1 }} />
                Выбрать файл
                <input
                  type="file"
                  hidden
                  // onChange={(e) => setFile(e.target.files?.[0])}
                />
              </Button>
            </Box>

            <Controller
              name="name"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  label="Название вложения"
                  fullWidth
                  error={!!errors.name}
                  helperText={errors.name?.message}
                  sx={{ mb: 2 }}
                />
              )}
            />

            <Controller
              name="description"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  label="Описание (необязательно)"
                  fullWidth
                  multiline
                  rows={3}
                  error={!!errors.description}
                  helperText={errors.description?.message}
                />
              )}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setModalOpen(false)}>
              Отмена
            </Button>
            <Button type="submit" variant="contained">
              Загрузить
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  )
}

