import React, { useState } from 'react'
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
  Tabs,
  Tab,
} from '@mui/material'
import { Add, School, VideoLibrary, Article, Assignment, Quiz } from '@mui/icons-material'

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`training-tabpanel-${index}`}
      aria-labelledby={`training-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  )
}

export const TrainingPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState(0)
  
  const materials = [
    { id: '1', title: 'Основы информационной безопасности', type: 'video', progress: 100, status: 'completed' },
    { id: '2', title: 'Работа с персональными данными', type: 'document', progress: 60, status: 'in_progress' },
    { id: '3', title: 'Антифишинговая защита', type: 'quiz', progress: 0, status: 'assigned' },
  ]

  const acknowledgments = [
    { id: '1', title: 'Политика информационной безопасности', type: 'acknowledgment', progress: 100, status: 'completed' },
    { id: '2', title: 'Правила работы с персональными данными', type: 'acknowledgment', progress: 0, status: 'assigned' },
    { id: '3', title: 'Инструкция по антивирусной защите', type: 'acknowledgment', progress: 75, status: 'in_progress' },
  ]

  const quizzes = [
    { id: '1', title: 'Тест по основам ИБ', type: 'quiz', progress: 100, status: 'completed', score: 85 },
    { id: '2', title: 'Квиз по работе с ПД', type: 'quiz', progress: 0, status: 'assigned', score: null },
    { id: '3', title: 'Проверка знаний по антифишингу', type: 'quiz', progress: 50, status: 'in_progress', score: null },
  ]

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue)
  }

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'video': return <VideoLibrary />
      case 'document': return <Article />
      case 'quiz': return <Quiz />
      case 'acknowledgment': return <Assignment />
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

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'video': return 'Видео'
      case 'document': return 'Документ'
      case 'quiz': return 'Квиз'
      case 'acknowledgment': return 'Ознакомление'
      default: return type
    }
  }

  const renderMaterialsTable = (items: any[], showScore = false) => (
    <Paper>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Материал</TableCell>
              <TableCell>Тип</TableCell>
              <TableCell>Прогресс</TableCell>
              <TableCell>Статус</TableCell>
              {showScore && <TableCell>Оценка</TableCell>}
              <TableCell>Срок сдачи</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {items.map((item) => (
              <TableRow key={item.id}>
                <TableCell>
                  <Box display="flex" alignItems="center">
                    {getTypeIcon(item.type)}
                    <Typography sx={{ ml: 1 }}>{item.title}</Typography>
                  </Box>
                </TableCell>
                <TableCell>
                  <Chip
                    label={getTypeLabel(item.type)}
                    size="small"
                    variant="outlined"
                  />
                </TableCell>
                <TableCell>
                  <Box display="flex" alignItems="center">
                    <LinearProgress
                      variant="determinate"
                      value={item.progress}
                      sx={{ width: 100, mr: 1 }}
                    />
                    <Typography variant="body2">{item.progress}%</Typography>
                  </Box>
                </TableCell>
                <TableCell>
                  <Chip
                    label={getStatusLabel(item.status)}
                    color={getStatusColor(item.status) as any}
                    size="small"
                  />
                </TableCell>
                {showScore && (
                  <TableCell>
                    {item.score ? `${item.score}%` : '-'}
                  </TableCell>
                )}
                <TableCell>2024-02-01</TableCell>
                <TableCell>
                  <Button size="small">
                    {item.status === 'completed' ? 'Просмотреть' : 'Начать'}
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Paper>
  )

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Обучение</Typography>
        <Button variant="contained" startIcon={<Add />}>
          Добавить материал
        </Button>
      </Box>

      <Paper>
        <Tabs value={activeTab} onChange={handleTabChange}>
          <Tab label="Материалы" icon={<School />} />
          <Tab label="Ознакомления" icon={<Assignment />} />
          <Tab label="Квизы" icon={<Quiz />} />
        </Tabs>

        <TabPanel value={activeTab} index={0}>
          {renderMaterialsTable(materials)}
        </TabPanel>
        <TabPanel value={activeTab} index={1}>
          {renderMaterialsTable(acknowledgments)}
        </TabPanel>
        <TabPanel value={activeTab} index={2}>
          {renderMaterialsTable(quizzes, true)}
        </TabPanel>
      </Paper>
    </Container>
  )
}
