import React from 'react'
import {
  Container,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  Box,
  Chip,
  LinearProgress,
} from '@mui/material'
import { Add, Warning } from '@mui/icons-material'

export const RisksPage: React.FC = () => {
  const risks = [
    { id: '1', title: 'Утечка данных', likelihood: 3, impact: 4, level: 12, status: 'draft' },
    { id: '2', title: 'Отказ сервера', likelihood: 2, impact: 5, level: 10, status: 'registered' },
    { id: '3', title: 'Фишинговая атака', likelihood: 4, impact: 3, level: 12, status: 'analysis' },
  ]

  const getRiskLevel = (level: number) => {
    if (level <= 5) return { color: 'success', label: 'Низкий' }
    if (level <= 10) return { color: 'warning', label: 'Средний' }
    if (level <= 15) return { color: 'error', label: 'Высокий' }
    return { color: 'error', label: 'Критический' }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'draft': return 'default'
      case 'registered': return 'info'
      case 'analysis': return 'warning'
      case 'treatment': return 'primary'
      case 'monitoring': return 'success'
      case 'closed': return 'success'
      default: return 'default'
    }
  }

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Риски</Typography>
        <Button variant="contained" startIcon={<Add />}>
          Добавить риск
        </Button>
      </Box>

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Название</TableCell>
                <TableCell>Вероятность</TableCell>
                <TableCell>Воздействие</TableCell>
                <TableCell>Уровень риска</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {risks.map((risk) => {
                const riskLevel = getRiskLevel(risk.level)
                return (
                  <TableRow key={risk.id}>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <Warning sx={{ mr: 1 }} />
                        {risk.title}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <LinearProgress
                          variant="determinate"
                          value={(risk.likelihood / 5) * 100}
                          sx={{ width: 60, mr: 1 }}
                        />
                        {risk.likelihood}/5
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <LinearProgress
                          variant="determinate"
                          value={(risk.impact / 5) * 100}
                          sx={{ width: 60, mr: 1 }}
                        />
                        {risk.impact}/5
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={`${riskLevel.label} (${risk.level})`}
                        color={riskLevel.color as any}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={risk.status}
                        color={getStatusColor(risk.status) as any}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <Button size="small">Редактировать</Button>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
    </Container>
  )
}
