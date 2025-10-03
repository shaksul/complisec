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
  IconButton,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material'
import { Add, Edit, Delete, Person } from '@mui/icons-material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

const commentSchema = z.object({
  content: z.string().min(1, 'Комментарий не может быть пустым'),
})

type CommentFormData = z.infer<typeof commentSchema>

interface RiskComment {
  id: string
  risk_id: string
  user_id: string
  user_name: string
  content: string
  created_at: string
  updated_at: string
}

interface RiskCommentsTabProps {
  riskId: string
}

export const RiskCommentsTab: React.FC<RiskCommentsTabProps> = ({ riskId }) => {
  const [comments, setComments] = useState<RiskComment[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingComment, setEditingComment] = useState<RiskComment | null>(null)

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CommentFormData>({
    resolver: zodResolver(commentSchema),
    defaultValues: {
      content: '',
    },
  })

  useEffect(() => {
    loadComments()
  }, [riskId])

  const loadComments = async () => {
    try {
      setLoading(true)
      // TODO: Implement API call to load comments
      // const response = await riskCommentsApi.list(riskId)
      // setComments(response.data || [])
      
      // Mock data for now
      setComments([
        {
          id: '1',
          risk_id: riskId,
          user_id: 'user1',
          user_name: 'Иван Петров',
          content: 'Необходимо срочно принять меры по снижению данного риска. Рекомендую внедрить дополнительные контроли.',
          created_at: new Date(Date.now() - 86400000).toISOString(), // 1 day ago
          updated_at: new Date(Date.now() - 86400000).toISOString(),
        },
        {
          id: '2',
          risk_id: riskId,
          user_id: 'user2',
          user_name: 'Мария Сидорова',
          content: 'Согласен с оценкой риска. Уровень Critical требует немедленного внимания руководства.',
          created_at: new Date(Date.now() - 172800000).toISOString(), // 2 days ago
          updated_at: new Date(Date.now() - 172800000).toISOString(),
        },
      ])
    } catch (err) {
      console.error('Error loading comments:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateNew = () => {
    setEditingComment(null)
    reset({ content: '' })
    setModalOpen(true)
  }

  const handleEdit = (comment: RiskComment) => {
    setEditingComment(comment)
    reset({ content: comment.content })
    setModalOpen(true)
  }

  const onSubmit = async (data: CommentFormData) => {
    try {
      if (editingComment) {
        // TODO: Implement API call to update comment
        // await riskCommentsApi.update(editingComment.id, data)
        setComments(prev => prev.map(c => 
          c.id === editingComment.id 
            ? { ...c, ...data, updated_at: new Date().toISOString() }
            : c
        ))
      } else {
        // TODO: Implement API call to create comment
        // const response = await riskCommentsApi.create({ ...data, risk_id: riskId })
        const newComment: RiskComment = {
          id: Date.now().toString(),
          risk_id: riskId,
          user_id: 'current_user',
          user_name: 'Текущий пользователь',
          content: data.content,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }
        setComments(prev => [newComment, ...prev])
      }
      setModalOpen(false)
      setEditingComment(null)
    } catch (err) {
      console.error('Error saving comment:', err)
    }
  }

  const handleDelete = async (commentId: string) => {
    try {
      // TODO: Implement API call to delete comment
      // await riskCommentsApi.delete(commentId)
      setComments(prev => prev.filter(c => c.id !== commentId))
    } catch (err) {
      console.error('Error deleting comment:', err)
    }
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
          Комментарии ({comments.length})
        </Typography>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={handleCreateNew}
        >
          Добавить комментарий
        </Button>
      </Box>

      {comments.length === 0 ? (
        <Card>
          <CardContent>
            <Box textAlign="center" py={4}>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                Комментариев пока нет
              </Typography>
              <Typography variant="body2" color="text.secondary" mb={3}>
                Добавьте первый комментарий для обсуждения риска
              </Typography>
              <Button
                variant="contained"
                startIcon={<Add />}
                onClick={handleCreateNew}
              >
                Добавить комментарий
              </Button>
            </Box>
          </CardContent>
        </Card>
      ) : (
        <Box>
          {comments.map((comment, index) => (
            <Card key={comment.id} sx={{ mb: 2 }}>
              <CardContent>
                <Box display="flex" gap={2}>
                  <Avatar>
                    <Person />
                  </Avatar>
                  
                  <Box flex={1}>
                    <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
                      <Typography variant="subtitle2" fontWeight="bold">
                        {comment.user_name}
                      </Typography>
                      <Box display="flex" alignItems="center" gap={1}>
                        <Typography variant="caption" color="text.secondary">
                          {formatDate(comment.created_at)}
                        </Typography>
                        <IconButton
                          size="small"
                          onClick={() => handleEdit(comment)}
                        >
                          <Edit fontSize="small" />
                        </IconButton>
                        <IconButton
                          size="small"
                          onClick={() => handleDelete(comment.id)}
                          color="error"
                        >
                          <Delete fontSize="small" />
                        </IconButton>
                      </Box>
                    </Box>
                    
                    <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap' }}>
                      {comment.content}
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          ))}
        </Box>
      )}

      {/* Add/Edit Comment Modal */}
      <Dialog open={modalOpen} onClose={() => setModalOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          {editingComment ? 'Редактирование комментария' : 'Добавление комментария'}
        </DialogTitle>
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogContent>
            <Controller
              name="content"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  label="Комментарий"
                  fullWidth
                  multiline
                  rows={4}
                  error={!!errors.content}
                  helperText={errors.content?.message}
                  placeholder="Введите ваш комментарий..."
                />
              )}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setModalOpen(false)}>
              Отмена
            </Button>
            <Button type="submit" variant="contained">
              {editingComment ? 'Обновить' : 'Добавить'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  )
}

