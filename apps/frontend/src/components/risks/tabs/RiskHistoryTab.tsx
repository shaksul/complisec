import React, { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Card,
  CardContent,
  Chip,
  Divider,
  CircularProgress,
} from '@mui/material'
import { History } from '@mui/icons-material'
import { risksApi, RiskHistoryEntry } from '../../../shared/api/risks'

interface RiskHistoryTabProps {
  riskId: string
}

const getFieldLabel = (field: string) => {
  switch (field) {
    case 'risk':
      return 'Риск'
    case 'title':
      return 'Название'
    case 'description':
      return 'Описание'
    case 'category':
      return 'Категория'
    case 'likelihood':
      return 'Вероятность'
    case 'impact':
      return 'Воздействие'
    case 'status':
      return 'Статус'
    case 'owner_user_id':
      return 'Ответственный'
    case 'methodology':
      return 'Методология'
    case 'strategy':
      return 'Стратегия'
    case 'due_date':
      return 'Срок обработки'
    default:
      return field
  }
}

const formatDateTime = (dateString: string) =>
  new Date(dateString).toLocaleString('ru-RU')

export const RiskHistoryTab: React.FC<RiskHistoryTabProps> = ({ riskId }) => {
  const [history, setHistory] = useState<RiskHistoryEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    void loadHistory()
  }, [riskId])

  const loadHistory = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await risksApi.getHistory(riskId)
      setHistory(response)
    } catch (err) {
      console.error('Error loading history:', err)
      setError('Не удалось загрузить историю изменений')
      setHistory([])
    } finally {
      setLoading(false)
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
      <Typography variant="h6" mb={2} display="flex" alignItems="center" gap={1}>
        <History /> История изменений
      </Typography>

      {error && (
        <Box mb={2}>
          <Typography color="error">{error}</Typography>
        </Box>
      )}

      {history.length === 0 ? (
        <Card>
          <CardContent>
            <Box textAlign="center" py={4}>
              <Typography variant="h6" color="text.secondary">
                История изменений отсутствует
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Изменения рисков появятся здесь автоматически после первых операций.
              </Typography>
            </Box>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardContent>
            <Box display="flex" flexDirection="column" gap={2}>
              {history.map((entry) => (
                <Box key={entry.id}>
                  <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
                    <Typography variant="subtitle1">
                      {getFieldLabel(entry.field_changed)}
                    </Typography>
                    <Chip
                      label={entry.changed_by_name ?? entry.changed_by}
                      variant="outlined"
                      size="small"
                    />
                  </Box>
                  <Typography variant="body2" color="text.secondary" mb={1}>
                    Изменено: {formatDateTime(entry.changed_at)}
                  </Typography>
                  <Box display="flex" flexDirection="column" gap={0.5}>
                    {entry.old_value !== undefined && (
                      <Typography variant="body2" color="text.secondary">
                        Было: {entry.old_value ?? '—'}
                      </Typography>
                    )}
                    {entry.new_value !== undefined && (
                      <Typography variant="body2">
                        Стало: {entry.new_value ?? '—'}
                      </Typography>
                    )}
                    {entry.change_reason && (
                      <Typography variant="body2" color="text.secondary">
                        Причина: {entry.change_reason}
                      </Typography>
                    )}
                  </Box>
                  <Divider sx={{ mt: 2 }} />
                </Box>
              ))}
            </Box>
          </CardContent>
        </Card>
      )}
    </Box>
  )
}


