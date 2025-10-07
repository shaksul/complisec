import React from 'react'
import {
  Box,
  Typography,
  Grid,
  Paper,
  Chip,
  Button,
  Divider,
  Card,
  CardContent,
  LinearProgress,
} from '@mui/material'
import { Edit, Warning, Person, Category, Assessment, Security } from '@mui/icons-material'
import { Risk } from '../../../shared/api/risks'
import { RISK_STATUSES, RISK_CATEGORIES, RISK_METHODOLOGIES, RISK_STRATEGIES } from '../../../shared/api/risks'

interface RiskGeneralTabProps {
  risk: Risk
  onUpdate: () => void
}

export const RiskGeneralTab: React.FC<RiskGeneralTabProps> = ({ risk }) => {

  const getRiskLevel = () => {
    const label = risk.level_label ?? (() => {
      if (!risk.level) return 'Не определен'
      if (risk.level <= 2) return 'Низкий'
      if (risk.level <= 4) return 'Средний'
      if (risk.level <= 6) return 'Высокий'
      return 'Критический'
    })()

    switch (label) {
      case 'Низкий':
        return { color: 'success', label, progressColor: 'success' as const }
      case 'Средний':
        return { color: 'warning', label, progressColor: 'warning' as const }
      case 'Высокий':
        return { color: 'warning', label, progressColor: 'warning' as const }
      case 'Критический':
        return { color: 'error', label, progressColor: 'error' as const }
      default:
        return { color: 'default', label: 'Не определен', progressColor: 'primary' as const }
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'new': return 'info'
      case 'in_analysis': return 'warning'
      case 'in_treatment': return 'primary'
      case 'accepted': return 'success'
      case 'transferred': return 'secondary'
      case 'mitigated': return 'success'
      case 'closed': return 'default'
      default: return 'default'
    }
  }

  const getStatusLabel = (status: string) => {
    const stat = RISK_STATUSES.find(s => s.value === status)
    return stat ? stat.label : status
  }

  const getCategoryLabel = (category?: string) => {
    if (!category) return 'Не указана'
    const cat = RISK_CATEGORIES.find(c => c.value === category)
    return cat ? cat.label : category
  }

  const getMethodologyLabel = (methodology?: string) => {
    if (!methodology) return 'Не указана'
    const meth = RISK_METHODOLOGIES.find(m => m.value === methodology)
    return meth ? meth.label : methodology
  }

  const getStrategyLabel = (strategy?: string) => {
    if (!strategy) return 'Не указана'
    const strat = RISK_STRATEGIES.find(s => s.value === strategy)
    return strat ? strat.label : strategy
  }

  const riskLevel = getRiskLevel()

  return (
    <Box>
      {/* Header with risk level and status */}
      <Paper sx={{ p: 3, mb: 3, bgcolor: 'grey.50' }}>
        <Grid container spacing={3} alignItems="center">
          <Grid item xs={12} md={6}>
            <Box display="flex" alignItems="center" gap={2}>
              <Warning color="warning" sx={{ fontSize: 40 }} />
              <Box>
                <Typography variant="h4" gutterBottom>
                  {risk.title}
                </Typography>
                <Box display="flex" gap={1} flexWrap="wrap">
                  <Chip
                    label={`${riskLevel.label}${risk.level ? ` (${risk.level})` : ''}`}
                    color={riskLevel.color as any}
                    size="medium"
                    sx={{ fontWeight: 'bold' }}
                  />
                  <Chip
                    label={getStatusLabel(risk.status)}
                    color={getStatusColor(risk.status) as any}
                    size="medium"
                  />
                </Box>
              </Box>
            </Box>
          </Grid>
          <Grid item xs={12} md={6}>
            <Box display="flex" justifyContent="flex-end">
              <Button
                variant="outlined"
                startIcon={<Edit />}
                onClick={() => console.log('Edit risk')}
              >
                Редактировать
              </Button>
            </Box>
          </Grid>
        </Grid>
      </Paper>

      <Grid container spacing={3}>
        {/* Left column - Basic information */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom display="flex" alignItems="center" gap={1}>
                <Assessment color="primary" />
                Оценка риска
              </Typography>
              <Divider sx={{ mb: 2 }} />
              
              <Box mb={2}>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Вероятность: {risk.likelihood ?? 0}/4
                </Typography>
                <LinearProgress
                  variant="determinate"
                  value={Math.min(((risk.likelihood ?? 0) / 4) * 100, 100)}
                  sx={{ mb: 1 }}
                />
              </Box>

              <Box mb={2}>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Воздействие: {risk.impact ?? 0}/4
                </Typography>
                <LinearProgress
                  variant="determinate"
                  value={Math.min(((risk.impact ?? 0) / 4) * 100, 100)}
                  sx={{ mb: 1 }}
                />
              </Box>

              <Box>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Уровень риска: {riskLevel.label}{risk.level ? ` (${risk.level})` : ''}
                </Typography>
                <LinearProgress
                  variant="determinate"
                  value={Math.min(((risk.level ?? 0) / 16) * 100, 100)}
                  color={riskLevel.progressColor}
                  sx={{ height: 8, borderRadius: 4 }}
                />
              </Box>
            </CardContent>
          </Card>

          <Card sx={{ mt: 2 }}>
            <CardContent>
              <Typography variant="h6" gutterBottom display="flex" alignItems="center" gap={1}>
                <Category color="primary" />
                Классификация
              </Typography>
              <Divider sx={{ mb: 2 }} />
              
              <Box mb={1}>
                <Typography variant="body2" color="text.secondary">
                  Категория:
                </Typography>
                <Typography variant="body1">
                  {getCategoryLabel(risk.category ?? undefined)}
                </Typography>
              </Box>

              <Box mb={1}>
                <Typography variant="body2" color="text.secondary">
                  Методология:
                </Typography>
                <Typography variant="body1">
                  {getMethodologyLabel(risk.methodology ?? undefined)}
                </Typography>
              </Box>

              <Box>
                <Typography variant="body2" color="text.secondary">
                  Стратегия обработки:
                </Typography>
                <Typography variant="body1">
                  {getStrategyLabel(risk.strategy ?? undefined)}
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Right column - Management information */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom display="flex" alignItems="center" gap={1}>
                <Person color="primary" />
                Управление
              </Typography>
              <Divider sx={{ mb: 2 }} />
              
              <Box mb={1}>
                <Typography variant="body2" color="text.secondary">
                  Ответственный:
                </Typography>
                <Typography variant="body1">
                  {risk.owner_user_id ? 'Назначен' : 'Не назначен'}
                </Typography>
              </Box>

              {risk.due_date && (
                <Box mb={1}>
                  <Typography variant="body2" color="text.secondary">
                    Срок обработки:
                  </Typography>
                  <Typography variant="body1">
                    {new Date(risk.due_date).toLocaleDateString('ru-RU')}
                  </Typography>
                </Box>
              )}

              <Box mb={1}>
                <Typography variant="body2" color="text.secondary">
                  Дата создания:
                </Typography>
                <Typography variant="body1">
                  {new Date(risk.created_at).toLocaleDateString('ru-RU')}
                </Typography>
              </Box>

              <Box>
                <Typography variant="body2" color="text.secondary">
                  Последнее обновление:
                </Typography>
                <Typography variant="body1">
                  {new Date(risk.updated_at).toLocaleDateString('ru-RU')}
                </Typography>
              </Box>
            </CardContent>
          </Card>

          {risk.asset_id && (
            <Card sx={{ mt: 2 }}>
              <CardContent>
                <Typography variant="h6" gutterBottom display="flex" alignItems="center" gap={1}>
                  <Security color="primary" />
                  Связанный актив
                </Typography>
                <Divider sx={{ mb: 2 }} />
                
                <Typography variant="body2" color="text.secondary">
                  ID актива:
                </Typography>
                <Typography variant="body1" fontFamily="monospace">
                  {risk.asset_id}
                </Typography>
              </CardContent>
            </Card>
          )}
        </Grid>

        {/* Description */}
        {risk.description && (
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Описание
                </Typography>
                <Divider sx={{ mb: 2 }} />
                <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
                  {risk.description}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>
    </Box>
  )
}
