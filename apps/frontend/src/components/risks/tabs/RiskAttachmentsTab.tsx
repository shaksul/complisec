import React, { useEffect, useState } from 'react'
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

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box>
      <Typography variant="h6" mb={2}>
        Вложения ({attachments.length})
      </Typography>

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
              <Typography variant="body2" color="text.secondary">
                Добавьте файлы через backend API или интеграцию с хранилищем, чтобы они появились в списке.
              </Typography>
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
                    primary={attachment.file_name}
                    secondary={
                      <Box display="flex" alignItems="center" gap={1} flexWrap="wrap">
                        <Typography variant="body2" color="text.secondary">
                          {formatFileSize(attachment.file_size)}
                        </Typography>
                        <Chip label={attachment.mime_type} size="small" variant="outlined" />
                        <Typography variant="body2" color="text.secondary">
                          Загрузил: {attachment.uploaded_by_name ?? attachment.uploaded_by}
                        </Typography>
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
    </Box>
  )
}


