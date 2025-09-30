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
import { Add, School, VideoLibrary, Article } from '@mui/icons-material'

export const TrainingPage: React.FC = () => {
  const materials = [
    { id: '1', title: 'Основы информационной безопасности', type: 'video', progress: 100, status: 'completed' },
    { id: '2', title: 'Работа с персональными данными', type: 'document', progress: 60, status: 'in_progress' },
    { id: '3', title: 'Антифишинговая защита', type: 'quiz', progress: 0, status: 'assigned' },
  ]

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'video': return <VideoLibrary />
      case 'document': return <Article />
      case 'quiz': return <School />
      default: return <Article />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'assigned': return 'default'
      case 'in_progress': return 'primary'
      case 'completed': return 'success'
      case 'failed': return 'error'
      case 'overdue': return 'error'
      default: return 'default'
    }
  }

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'assigned': return 'Назначено'
      case 'in_progress': return 'В процессе'
      case 'completed': return 'Завершено'
      case 'failed': return 'Не пройдено'
      case 'overdue': return 'Просрочено'
      default: return status
    }
  }

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Обучение</Typography>
        <Button variant="contained" startIcon={<Add />}>
          Добавить материал
        </Button>
      </Box>

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Материал</TableCell>
                <TableCell>Тип</TableCell>
                <TableCell>Прогресс</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Срок сдачи</TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {materials.map((material) => (
                <TableRow key={material.id}>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      {getTypeIcon(material.type)}
                      <Typography sx={{ ml: 1 }}>{material.title}</Typography>
                    </Box>
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={material.type}
                      size="small"
                      variant="outlined"
                    />
                  </TableCell>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      <LinearProgress
                        variant="determinate"
                        value={material.progress}
                        sx={{ width: 100, mr: 1 }}
                      />
                      <Typography variant="body2">{material.progress}%</Typography>
                    </Box>
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={getStatusLabel(material.status)}
                      color={getStatusColor(material.status) as any}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>2024-02-01</TableCell>
                  <TableCell>
                    <Button size="small">
                      {material.status === 'completed' ? 'Просмотреть' : 'Начать'}
                    </Button>
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
