import React, { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Card,
  CardContent,
  Chip,
  Timeline,
  TimelineItem,
  TimelineSeparator,
  TimelineConnector,
  TimelineContent,
  TimelineDot,
  Avatar,
  CircularProgress,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from '@mui/material'
import { 
  Edit, 
  Add, 
  Delete, 
  Person, 
  ExpandMore,
  History,
  Security,
  Assessment,
} from '@mui/icons-material'

interface RiskHistoryEntry {
  id: string
  risk_id: string
  user_id: string
  user_name: string
  action: string
  field: string
  old_value?: string
  new_value?: string
  created_at: string
}

interface RiskHistoryTabProps {
  riskId: string
}

export const RiskHistoryTab: React.FC<RiskHistoryTabProps> = ({ riskId }) => {
  const [history, setHistory] = useState<RiskHistoryEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [expandedItems, setExpandedItems] = useState<Set<string>>(new Set())

  useEffect(() => {
    loadHistory()
  }, [riskId])

  const loadHistory = async () => {
    try {
      setLoading(true)
      // TODO: Implement API call to load history
      // const response = await riskHistoryApi.list(riskId)
      // setHistory(response.data || [])
      
      // Mock data for now
      setHistory([
        {
          id: '1',
          risk_id: riskId,
          user_id: 'user1',
          user_name: 'Иван Петров',
          action: 'created',
          field: 'risk',
          created_at: new Date(Date.now() - 259200000).toISOString(), // 3 days ago
        },
        {
          id: '2',
          risk_id: riskId,
          user_id: 'user2',
          user_name: 'Мария Сидорова',
          action: 'updated',
          field: 'likelihood',
          old_value: '2',
          new_value: '3',
          created_at: new Date(Date.now() - 172800000).toISOString(), // 2 days ago
        },
        {
          id: '3',
          risk_id: riskId,
          user_id: 'user2',
          user_name: 'Мария Сидорова',
          action: 'updated',
          field: 'impact',
          old_value: '3',
          new_value: '4',
          created_at: new Date(Date.now() - 172800000).toISOString(), // 2 days ago
        },
        {
          id: '4',
          risk_id: riskId,
          user_id: 'user1',
          user_name: 'Иван Петров',
          action: 'updated',
          field: 'status',
          old_value: 'new',
          new_value: 'in_analysis',
          created_at: new Date(Date.now() - 86400000).toISOString(), // 1 day ago
        },
        {
          id: '5',
          risk_id: riskId,
          user_id: 'user3',
          user_name: 'Алексей Козлов',
          action: 'added',
          field: 'control',
          new_value: 'Многофакторная аутентификация',
          created_at: new Date(Date.now() - 43200000).toISOString(), // 12 hours ago
        },
      ])
    } catch (err) {
      console.error('Error loading history:', err)
    } finally {
      setLoading(false)
    }
  }

  const getActionIcon = (action: string) => {
    switch (action) {
      case 'created': return <Add />
      case 'updated': return <Edit />
      case 'deleted': return <Delete />
      case 'added': return <Add />
      case 'removed': return <Delete />
      default: return <Edit />
    }
  }

  const getActionColor = (action: string) => {
    switch (action) {
      case 'created': return 'success'
      case 'updated': return 'primary'
      case 'deleted': return 'error'
      case 'added': return 'success'
      case 'removed': return 'error'
      default: return 'default'
    }
  }

  const getActionLabel = (action: string) => {
    switch (action) {
      case 'created': return 'Создан'
      case 'updated': return 'Обновлен'
      case 'deleted': return 'Удален'
      case 'added': return 'Добавлен'
      case 'removed': return 'Удален'
      default: return action
    }
  }

  const getFieldLabel = (field: string) => {
    switch (field) {
      case 'risk': return 'Риск'
      case 'title': return 'Название'
      case 'description': return 'Описание'
      case 'category': return 'Категория'
      case 'likelihood': return 'Вероятность'
      case 'impact': return 'Воздействие'
      case 'status': return 'Статус'
      case 'owner_user_id': return 'Ответственный'
      case 'methodology': return 'Методология'
      case 'strategy': return 'Стратегия'
      case 'due_date': return 'Срок обработки'
      case 'control': return 'Контроль'
      case 'comment': return 'Комментарий'
      default: return field
    }
  }

  const getValueLabel = (field: string, value?: string) => {
    if (!value) return ''
    
    switch (field) {
      case 'status':
        switch (value) {
          case 'new': return 'Новый'
          case 'in_analysis': return 'В анализе'
          case 'in_treatment': return 'В обработке'
          case 'accepted': return 'Принят'
          case 'transferred': return 'Передан'
          case 'mitigated': return 'Смягчен'
          case 'closed': return 'Закрыт'
          default: return value
        }
      case 'category':
        switch (value) {
          case 'security': return 'Безопасность'
          case 'operational': return 'Операционные'
          case 'financial': return 'Финансовые'
          case 'compliance': return 'Соответствие'
          case 'reputation': return 'Репутационные'
          case 'legal': return 'Правовые'
          case 'strategic': return 'Стратегические'
          default: return value
        }
      default: return value
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
    
    if (diffHours < 1) {
      const diffMinutes = Math.floor(diffMs / (1000 * 60))
      return `${diffMinutes} мин назад`
    } else if (diffHours < 24) {
      return `${diffHours} ч назад`
    } else if (diffDays === 1) {
      return 'Вчера'
    } else if (diffDays < 7) {
      return `${diffDays} дней назад`
    } else {
      return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }
  }

  const handleAccordionChange = (itemId: string) => {
    const newExpanded = new Set(expandedItems)
    if (newExpanded.has(itemId)) {
      newExpanded.delete(itemId)
    } else {
      newExpanded.add(itemId)
    }
    setExpandedItems(newExpanded)
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
      <Box display="flex" alignItems="center" gap={1} mb={3}>
        <History color="primary" />
        <Typography variant="h6">
          История изменений ({history.length})
        </Typography>
      </Box>

      {history.length === 0 ? (
        <Card>
          <CardContent>
            <Box textAlign="center" py={4}>
              <History sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" color="text.secondary" gutterBottom>
                История изменений пуста
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Здесь будет отображаться история всех изменений риска
              </Typography>
            </Box>
          </CardContent>
        </Card>
      ) : (
        <Box>
          {history.map((entry, index) => (
            <Accordion 
              key={entry.id}
              expanded={expandedItems.has(entry.id)}
              onChange={() => handleAccordionChange(entry.id)}
            >
              <AccordionSummary expandIcon={<ExpandMore />}>
                <Box display="flex" alignItems="center" gap={2} width="100%">
                  <Avatar sx={{ width: 32, height: 32 }}>
                    <Person fontSize="small" />
                  </Avatar>
                  
                  <Box flex={1}>
                    <Box display="flex" alignItems="center" gap={1} mb={0.5}>
                      <Typography variant="subtitle2" fontWeight="bold">
                        {entry.user_name}
                      </Typography>
                      <Chip
                        icon={getActionIcon(entry.action)}
                        label={getActionLabel(entry.action)}
                        size="small"
                        color={getActionColor(entry.action) as any}
                      />
                      <Typography variant="caption" color="text.secondary">
                        {getFieldLabel(entry.field)}
                      </Typography>
                    </Box>
                    <Typography variant="caption" color="text.secondary">
                      {formatDate(entry.created_at)}
                    </Typography>
                  </Box>
                </Box>
              </AccordionSummary>
              
              <AccordionDetails>
                <Box>
                  {entry.old_value && entry.new_value ? (
                    <Box>
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        Изменение значения:
                      </Typography>
                      <Box display="flex" gap={1} flexWrap="wrap">
                        <Chip
                          label={getValueLabel(entry.field, entry.old_value)}
                          size="small"
                          color="error"
                          variant="outlined"
                        />
                        <Typography variant="body2" sx={{ alignSelf: 'center' }}>
                          →
                        </Typography>
                        <Chip
                          label={getValueLabel(entry.field, entry.new_value)}
                          size="small"
                          color="success"
                        />
                      </Box>
                    </Box>
                  ) : entry.new_value ? (
                    <Box>
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        Добавлено:
                      </Typography>
                      <Chip
                        label={getValueLabel(entry.field, entry.new_value)}
                        size="small"
                        color="success"
                      />
                    </Box>
                  ) : entry.old_value ? (
                    <Box>
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        Удалено:
                      </Typography>
                      <Chip
                        label={getValueLabel(entry.field, entry.old_value)}
                        size="small"
                        color="error"
                        variant="outlined"
                      />
                    </Box>
                  ) : (
                    <Typography variant="body2" color="text.secondary">
                      {getActionLabel(entry.action)} {getFieldLabel(entry.field)}
                    </Typography>
                  )}
                </Box>
              </AccordionDetails>
            </Accordion>
          ))}
        </Box>
      )}
    </Box>
  )
}

