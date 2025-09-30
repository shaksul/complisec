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
} from '@mui/material'
import { Add, Report } from '@mui/icons-material'

export const IncidentsPage: React.FC = () => {
  const incidents = [
    { id: '1', title: 'Подозрительная активность в сети', severity: 'high', status: 'investigating' },
    { id: '2', title: 'Недоступность веб-сайта', severity: 'medium', status: 'resolved' },
    { id: '3', title: 'Попытка несанкционированного доступа', severity: 'critical', status: 'new' },
  ]

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'low': return 'success'
      case 'medium': return 'warning'
      case 'high': return 'error'
      case 'critical': return 'error'
      default: return 'default'
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'new': return 'info'
      case 'triage': return 'warning'
      case 'investigating': return 'primary'
      case 'contained': return 'warning'
      case 'resolved': return 'success'
      case 'closed': return 'success'
      case 'false_positive': return 'default'
      default: return 'default'
    }
  }

  const getSeverityLabel = (severity: string) => {
    switch (severity) {
      case 'low': return 'Низкий'
      case 'medium': return 'Средний'
      case 'high': return 'Высокий'
      case 'critical': return 'Критический'
      default: return severity
    }
  }

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'new': return 'Новый'
      case 'triage': return 'Классификация'
      case 'investigating': return 'Расследование'
      case 'contained': return 'Локализован'
      case 'resolved': return 'Решен'
      case 'closed': return 'Закрыт'
      case 'false_positive': return 'Ложное срабатывание'
      default: return status
    }
  }

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Инциденты</Typography>
        <Button variant="contained" startIcon={<Add />}>
          Создать инцидент
        </Button>
      </Box>

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Название</TableCell>
                <TableCell>Критичность</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Дата создания</TableCell>
                <TableCell>Назначен</TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {incidents.map((incident) => (
                <TableRow key={incident.id}>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      <Report sx={{ mr: 1 }} />
                      {incident.title}
                    </Box>
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={getSeverityLabel(incident.severity)}
                      color={getSeverityColor(incident.severity) as any}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={getStatusLabel(incident.status)}
                      color={getStatusColor(incident.status) as any}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>2024-01-15 14:30</TableCell>
                  <TableCell>Иванов И.И.</TableCell>
                  <TableCell>
                    <Button size="small">Редактировать</Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
    </Container>
  )
}
