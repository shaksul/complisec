import React, { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Button,
  Card,
  CardContent,
  TextField,
  Avatar,
  Divider,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material'
import { Add } from '@mui/icons-material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { risksApi, RiskComment } from '../../../shared/api/risks'

const commentSchema = z.object({
  comment: z.string().min(1, 'Комментарий не может быть пустым'),
})

type CommentFormData = z.infer<typeof commentSchema>

interface RiskCommentsTabProps {
  riskId: string
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays === 0) {
    return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
  }
  if (diffDays === 1) {
    return 'Вчера'
  }
  if (diffDays < 7) {
    return `${diffDays} дней назад`
  }
  return date.toLocaleDateString('ru-RU')
}

export const RiskCommentsTab: React.FC<RiskCommentsTabProps> = ({ riskId }) => {
  const [comments, setComments] = useState<RiskComment[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CommentFormData>({
    resolver: zodResolver(commentSchema),
    defaultValues: {
      comment: '',
    },
  })

  useEffect(() => {
    void loadComments()
  }, [riskId])

  const loadComments = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await risksApi.getComments(riskId)
      setComments(response)
    } catch (err) {
      console.error('Error loading comments:', err)
      setError('Не удалось загрузить комментарии')
      setComments([])
    } finally {
      setLoading(false)
    }
  }

  const handleCreateNew = () => {
    reset({ comment: '' })
    setModalOpen(true)
  }

  const onSubmit = async (data: CommentFormData) => {
    try {
      await risksApi.createComment(riskId, { comment: data.comment })
      setModalOpen(false)
      reset({ comment: '' })
      await loadComments()
    } catch (err) {
      console.error('Error saving comment:', err)
      setError('Не удалось добавить комментарий')
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
        <Typography variant="h6">Комментарии ({comments.length})</Typography>
        <Button variant="contained" startIcon={<Add />} onClick={handleCreateNew}>
          Добавить комментарий
        </Button>
      </Box>

      {error && (
        <Box mb={2}>
          <Typography color="error">{error}</Typography>
        </Box>
      )}

      {comments.length === 0 ? (
        <Card>
          <CardContent>
            <Box textAlign="center" py={4}>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                Комментариев пока нет
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Оставьте первый комментарий, чтобы задокументировать решение по риску.
              </Typography>
            </Box>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardContent>
            <Box display="flex" flexDirection="column" gap={2}>
              {comments.map((comment) => (
                <Box key={comment.id}>
                  <Box display="flex" alignItems="center" gap={2} mb={1}>
                    <Avatar sx={{ width: 32, height: 32 }}>
                      {(comment.user_name ?? comment.user_id).slice(0, 2).toUpperCase()}
                    </Avatar>
                    <Box>
                      <Typography variant="subtitle2">
                        {comment.user_name ?? 'Неизвестный пользователь'}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {formatDate(comment.created_at)}
                      </Typography>
                    </Box>
                  </Box>
                  <Typography variant="body1" paragraph>
                    {comment.comment}
                  </Typography>
                  <Divider />
                </Box>
              ))}
            </Box>
          </CardContent>
        </Card>
      )}

      <Dialog open={modalOpen} onClose={() => setModalOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Добавление комментария</DialogTitle>
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogContent>
            <Controller
              name="comment"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  label="Комментарий"
                  fullWidth
                  multiline
                  rows={4}
                  error={!!errors.comment}
                  helperText={errors.comment?.message}
                />
              )}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setModalOpen(false)}>Отмена</Button>
            <Button type="submit" variant="contained">
              Сохранить
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  )
}

