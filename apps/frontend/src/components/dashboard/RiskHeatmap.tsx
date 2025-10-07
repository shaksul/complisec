import React, { useEffect, useState } from 'react'
import {
  Box,
  Typography,
  Paper,
  Tooltip,
  Grid,
  Chip,
  Skeleton,
} from '@mui/material'
import { alpha, useTheme } from '@mui/material/styles'
import type { CorporateTheme } from '../../shared/theme'
import { risksApi, type Risk } from '../../shared/api/risks'
import { useAuth } from '../../contexts/AuthContext'

interface HeatmapCell {
  likelihood: number
  impact: number
  count: number
  risks: Risk[]
  level: number
}

const RISK_LEVEL_LABEL: Record<'low' | 'medium' | 'high' | 'critical', string> = {
  low: 'Низкий',
  medium: 'Средний',
  high: 'Высокий',
  critical: 'Критический',
}

export const RiskHeatmap: React.FC = () => {
  const [heatmapData, setHeatmapData] = useState<HeatmapCell[][]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const theme = useTheme<CorporateTheme>()
  const { user } = useAuth()

  useEffect(() => {
    if (user) {
      void loadHeatmapData()
    } else {
      setLoading(false)
    }
  }, [user])

  const loadHeatmapData = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await risksApi.list({ page_size: 1000 })
      const risks: Risk[] = response.data || []

      const matrix: HeatmapCell[][] = []

      for (let impact = 4; impact >= 1; impact--) {
        const row: HeatmapCell[] = []

        for (let likelihood = 1; likelihood <= 4; likelihood++) {
          const cellRisks = risks.filter(
            (risk) => (risk.likelihood || 1) === likelihood && (risk.impact || 1) === impact,
          )

          row.push({
            likelihood,
            impact,
            count: cellRisks.length,
            risks: cellRisks,
            level: likelihood * impact,
          })
        }

        matrix.push(row)
      }

      setHeatmapData(matrix)
    } catch (err) {
      console.error('Error loading heatmap data:', err)
      setError('Не удалось загрузить тепловую карту. Попробуйте обновить страницу позже.')
    } finally {
      setLoading(false)
    }
  }

  const getCellColor = (level: number, count: number) => {
    if (count === 0) {
      return theme.palette.background.default
    }

    if (level <= 2) {
      return alpha(theme.palette.success.main, 0.85)
    }

    if (level <= 4) {
      return alpha(theme.palette.warning.light ?? theme.palette.warning.main, 0.9)
    }

    if (level <= 6) {
      return alpha(theme.palette.warning.main, 0.95)
    }

    return alpha(theme.palette.error.main, 0.92)
  }

  const getCellOpacity = (count: number) => {
    if (count === 0) return 0.32
    if (count === 1) return 0.72
    if (count === 2) return 0.84
    return 1
  }

  const getLevelKey = (level: number): keyof typeof RISK_LEVEL_LABEL => {
    if (level <= 2) return 'low'
    if (level <= 4) return 'medium'
    if (level <= 6) return 'high'
    return 'critical'
  }

  const getLevelChipColor = (level: number) => {
    if (level <= 2) return 'success'
    if (level <= 6) return 'warning'
    return 'error'
  }

  const renderRiskList = (risks: Risk[]) => (
    <Box mt={1}>
      <Typography variant="caption" display="block" color="text.secondary">
        Примеры рисков:
      </Typography>
      {risks.slice(0, 3).map((risk) => (
        <Typography key={risk.id} variant="caption" display="block">
          • {risk.title}
        </Typography>
      ))}
      {risks.length > 3 && (
        <Typography variant="caption" color="text.secondary" display="block">
          Ещё {risks.length - 3}
        </Typography>
      )}
    </Box>
  )

  if (loading) {
    return (
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Тепловая карта рисков
        </Typography>
        <Skeleton variant="rounded" height={260} sx={{ borderRadius: 3 }} />
      </Paper>
    )
  }

  if (error) {
    return (
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Тепловая карта рисков
        </Typography>
        <Typography color="error" variant="body2">
          {error}
        </Typography>
      </Paper>
    )
  }

  return (
    <Paper sx={{ p: 3, display: 'flex', flexDirection: 'column', gap: 3 }}>
      <Box>
        <Typography variant="h6" gutterBottom>
          Тепловая карта рисков
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Матрица вероятности и влияния помогает увидеть концентрацию рисков и определить приоритеты обработки.
        </Typography>
      </Box>

      <Box display="flex" flexDirection="column" gap={2}>
        {heatmapData.map((row, rowIndex) => (
          <Box key={`row-${rowIndex}`} display="flex" gap={1.5}>
            {row.map((cell) => (
              <Tooltip
                key={`${cell.likelihood}-${cell.impact}`}
                arrow
                title={
                  <Box>
                    <Typography variant="subtitle2" gutterBottom>
                      Вероятность: {cell.likelihood}, Влияние: {cell.impact}
                    </Typography>
                    <Typography variant="body2">
                      Уровень: {RISK_LEVEL_LABEL[getLevelKey(cell.level)]} ({cell.level})
                    </Typography>
                    <Typography variant="body2">Количество рисков: {cell.count}</Typography>
                    {cell.risks.length > 0 && renderRiskList(cell.risks)}
                  </Box>
                }
              >
                <Box
                  width={88}
                  height={64}
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  bgcolor={getCellColor(cell.level, cell.count)}
                  sx={{
                    border: `1px solid ${theme.palette.divider}`,
                    borderRadius: 2,
                    opacity: getCellOpacity(cell.count),
                    cursor: cell.count > 0 ? 'pointer' : 'default',
                    transition: 'all 0.2s ease',
                    '&:hover': {
                      opacity: 1,
                      transform: cell.count > 0 ? 'translateY(-4px)' : 'none',
                    },
                  }}
                >
                  <Typography
                    variant="body1"
                    fontWeight={700}
                    sx={{
                      color: cell.count > 0 ? theme.palette.common.white : theme.palette.text.secondary,
                      textShadow: cell.count > 0 ? `0 1px 2px ${alpha(theme.palette.common.black, 0.35)}` : 'none',
                    }}
                  >
                    {cell.count}
                  </Typography>
                </Box>
              </Tooltip>
            ))}
          </Box>
        ))}
      </Box>

      <Box>
        <Typography variant="subtitle2" gutterBottom>
          Обозначения
        </Typography>
        <Grid container spacing={2}>
          <Grid item>
            <Box display="flex" alignItems="center" gap={1}>
              <Box
                width={22}
                height={22}
                bgcolor={alpha(theme.palette.success.main, 0.85)}
                borderRadius={2}
                border={`1px solid ${theme.palette.divider}`}
              />
              <Typography variant="caption">Низкий (1–2)</Typography>
            </Box>
          </Grid>
          <Grid item>
            <Box display="flex" alignItems="center" gap={1}>
              <Box
                width={22}
                height={22}
                bgcolor={alpha(theme.palette.warning.light ?? theme.palette.warning.main, 0.9)}
                borderRadius={2}
                border={`1px solid ${theme.palette.divider}`}
              />
              <Typography variant="caption">Средний (3–4)</Typography>
            </Box>
          </Grid>
          <Grid item>
            <Box display="flex" alignItems="center" gap={1}>
              <Box
                width={22}
                height={22}
                bgcolor={alpha(theme.palette.warning.main, 0.95)}
                borderRadius={2}
                border={`1px solid ${theme.palette.divider}`}
              />
              <Typography variant="caption">Высокий (5–6)</Typography>
            </Box>
          </Grid>
          <Grid item>
            <Box display="flex" alignItems="center" gap={1}>
              <Box
                width={22}
                height={22}
                bgcolor={alpha(theme.palette.error.main, 0.92)}
                borderRadius={2}
                border={`1px solid ${theme.palette.divider}`}
              />
              <Typography variant="caption">Критический (7+)</Typography>
            </Box>
          </Grid>
        </Grid>
      </Box>

      <Box>
        <Typography variant="subtitle2" gutterBottom>
          Распределение по уровням
        </Typography>
        <Grid container spacing={1.5}>
          {[...Array(16)].map((_, index) => {
            const level = index + 1
            const levelCount = heatmapData.flat().reduce((total, cell) => {
              if (cell.level === level) {
                return total + cell.count
              }
              return total
            }, 0)

            if (levelCount === 0) {
              return null
            }

            const key = getLevelKey(level)

            return (
              <Grid item key={`chip-${level}`}>
                <Chip
                  label={`${RISK_LEVEL_LABEL[key]} (${level}) — ${levelCount}`}
                  size="small"
                  color={getLevelChipColor(level) as any}
                  variant="outlined"
                />
              </Grid>
            )
          })}
        </Grid>
      </Box>
    </Paper>
  )
}
