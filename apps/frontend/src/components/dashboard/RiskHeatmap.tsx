import React, { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Paper,
  Tooltip,
  Card,
  CardContent,
  Grid,
  Chip,
} from '@mui/material'
import { risksApi, Risk } from '../../shared/api/risks'

interface HeatmapCell {
  likelihood: number
  impact: number
  count: number
  risks: Risk[]
  level: number
}

export const RiskHeatmap: React.FC = () => {
  const [heatmapData, setHeatmapData] = useState<HeatmapCell[][]>([])
  const [loading, setLoading] = useState(true)
  const [totalRisks, setTotalRisks] = useState(0)

  useEffect(() => {
    loadHeatmapData()
  }, [])

  const loadHeatmapData = async () => {
    try {
      setLoading(true)
      const response = await risksApi.list()
      const risks: Risk[] = response.data || []
      
      setTotalRisks(risks.length)
      
      // Create 4x4 heatmap matrix
      const matrix: HeatmapCell[][] = []
      
      for (let impact = 4; impact >= 1; impact--) {
        const row: HeatmapCell[] = []
        for (let likelihood = 1; likelihood <= 4; likelihood++) {
          const cellRisks = risks.filter(risk => 
            (risk.likelihood || 1) === likelihood && (risk.impact || 1) === impact
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
    } finally {
      setLoading(false)
    }
  }

  const getCellColor = (level: number, count: number) => {
    if (count === 0) return '#f5f5f5'
    
    // Color intensity based on risk level
    switch (level) {
      case 1:
      case 2:
        return '#4caf50' // Green for Low
      case 3:
      case 4:
        return '#ffeb3b' // Yellow for Medium
      case 5:
      case 6:
        return '#ff9800' // Orange for High
      case 7:
      case 8:
      case 9:
      case 10:
      case 11:
      case 12:
      case 13:
      case 14:
      case 15:
      case 16:
        return '#f44336' // Red for Critical
      default:
        return '#e0e0e0'
    }
  }

  const getCellOpacity = (count: number) => {
    if (count === 0) return 0.3
    if (count === 1) return 0.6
    if (count === 2) return 0.8
    return 1.0
  }

  const getRiskLevelLabel = (level: number) => {
    if (level <= 2) return 'Low'
    if (level <= 4) return 'Medium'
    if (level <= 6) return 'High'
    return 'Critical'
  }

  const getRiskLevelColor = (level: number) => {
    if (level <= 2) return 'success'
    if (level <= 4) return 'warning'
    if (level <= 6) return 'warning'
    return 'error'
  }

  if (loading) {
    return (
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Матрица рисков
        </Typography>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
          <Typography>Загрузка данных...</Typography>
        </Box>
      </Paper>
    )
  }

  return (
    <Paper sx={{ p: 3 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h6">
          Матрица рисков (4×4)
        </Typography>
        <Chip 
          label={`Всего рисков: ${totalRisks}`}
          color="primary"
          variant="outlined"
        />
      </Box>

      {/* Heatmap Grid */}
      <Box>
        {/* Headers */}
        <Box display="flex" mb={1}>
          <Box width={60} />
          {[1, 2, 3, 4].map(likelihood => (
            <Box key={likelihood} width={80} textAlign="center">
              <Typography variant="caption" fontWeight="bold">
                {likelihood}
              </Typography>
            </Box>
          ))}
        </Box>

        {/* Heatmap rows */}
        {heatmapData.map((row, impactIndex) => {
          const impact = 4 - impactIndex // Reverse order for display
          return (
            <Box key={impact} display="flex" mb={1}>
              {/* Impact label */}
              <Box width={60} display="flex" alignItems="center">
                <Typography variant="caption" fontWeight="bold">
                  {impact}
                </Typography>
              </Box>

              {/* Cells */}
              {row.map((cell, likelihoodIndex) => (
                <Tooltip
                  key={`${cell.likelihood}-${cell.impact}`}
                  title={
                    <Box>
                      <Typography variant="subtitle2" gutterBottom>
                        Вероятность: {cell.likelihood}, Воздействие: {cell.impact}
                      </Typography>
                      <Typography variant="body2">
                        Уровень: {getRiskLevelLabel(cell.level)} ({cell.level})
                      </Typography>
                      <Typography variant="body2">
                        Количество рисков: {cell.count}
                      </Typography>
                      {cell.risks.length > 0 && (
                        <Box mt={1}>
                          <Typography variant="caption" display="block">
                            Риски:
                          </Typography>
                          {cell.risks.slice(0, 3).map(risk => (
                            <Typography key={risk.id} variant="caption" display="block">
                              • {risk.title}
                            </Typography>
                          ))}
                          {cell.risks.length > 3 && (
                            <Typography variant="caption" display="block">
                              ... и еще {cell.risks.length - 3}
                            </Typography>
                          )}
                        </Box>
                      )}
                    </Box>
                  }
                  arrow
                >
                  <Box
                    width={80}
                    height={60}
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    bgcolor={getCellColor(cell.level, cell.count)}
                    sx={{
                      opacity: getCellOpacity(cell.count),
                      border: '1px solid #ddd',
                      cursor: cell.count > 0 ? 'pointer' : 'default',
                      transition: 'opacity 0.2s',
                      '&:hover': {
                        opacity: 1,
                        transform: 'scale(1.05)',
                      },
                    }}
                    borderRadius={1}
                  >
                    <Typography
                      variant="body2"
                      fontWeight="bold"
                      color={cell.count > 0 ? 'white' : 'text.secondary'}
                      sx={{ textShadow: cell.count > 0 ? '1px 1px 2px rgba(0,0,0,0.5)' : 'none' }}
                    >
                      {cell.count}
                    </Typography>
                  </Box>
                </Tooltip>
              ))}
            </Box>
          )
        })}
      </Box>

      {/* Legend */}
      <Box mt={3}>
        <Typography variant="subtitle2" gutterBottom>
          Легенда:
        </Typography>
        <Grid container spacing={2}>
          <Grid item>
            <Box display="flex" alignItems="center" gap={1}>
              <Box
                width={20}
                height={20}
                bgcolor="#4caf50"
                borderRadius={1}
                border="1px solid #ddd"
              />
              <Typography variant="caption">Low (1-2)</Typography>
            </Box>
          </Grid>
          <Grid item>
            <Box display="flex" alignItems="center" gap={1}>
              <Box
                width={20}
                height={20}
                bgcolor="#ffeb3b"
                borderRadius={1}
                border="1px solid #ddd"
              />
              <Typography variant="caption">Medium (3-4)</Typography>
            </Box>
          </Grid>
          <Grid item>
            <Box display="flex" alignItems="center" gap={1}>
              <Box
                width={20}
                height={20}
                bgcolor="#ff9800"
                borderRadius={1}
                border="1px solid #ddd"
              />
              <Typography variant="caption">High (5-6)</Typography>
            </Box>
          </Grid>
          <Grid item>
            <Box display="flex" alignItems="center" gap={1}>
              <Box
                width={20}
                height={20}
                bgcolor="#f44336"
                borderRadius={1}
                border="1px solid #ddd"
              />
              <Typography variant="caption">Critical (7+)</Typography>
            </Box>
          </Grid>
        </Grid>
        
        <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
          Интенсивность цвета показывает количество рисков в каждой ячейке
        </Typography>
      </Box>

      {/* Summary Statistics */}
      <Box mt={3}>
        <Typography variant="subtitle2" gutterBottom>
          Статистика по уровням риска:
        </Typography>
        <Grid container spacing={2}>
          {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16].map(level => {
            const count = heatmapData.flat().reduce((sum, cell) => {
              if (cell.level === level) {
                return sum + cell.count
              }
              return sum
            }, 0)
            
            if (count === 0) return null
            
            return (
              <Grid item key={level}>
                <Chip
                  label={`${getRiskLevelLabel(level)} (${level}): ${count}`}
                  size="small"
                  color={getRiskLevelColor(level) as any}
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

