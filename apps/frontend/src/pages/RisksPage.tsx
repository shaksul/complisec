import React, { useState, useEffect } from 'react'
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
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
  IconButton,
  Tooltip,
  TableSortLabel,
} from '@mui/material'
import {
  Add,
  Warning,
  Edit,
  Search,
  FilterList,
  Download,
  Clear,
  Delete,
  Visibility
} from '@mui/icons-material'
import { RiskModal } from '../components/risks/RiskModal'
import { RiskDetailsModal } from '../components/risks/RiskDetailsModal'
import { 
  risksApi, 
  Risk, 
  CreateRiskRequest, 
  UpdateRiskRequest,
  RISK_STATUSES,
  RISK_CATEGORIES,
  RISK_LEVELS,
  User
} from '../shared/api/risks'
import { getUsers } from '../shared/api/users'

type SortField = 'level' | 'created_at' | 'category' | 'title'
type SortDirection = 'asc' | 'desc'

interface Filters {
  status: string
  category: string
  owner_user_id: string
  level: string
  search: string
}

export const RisksPage: React.FC = () => {
  console.log('RisksPage component rendered')
  const [risks, setRisks] = useState<Risk[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingRisk, setEditingRisk] = useState<Risk | null>(null)
  const [detailsModalOpen, setDetailsModalOpen] = useState(false)
  const [selectedRisk, setSelectedRisk] = useState<Risk | null>(null)
  const [users, setUsers] = useState<User[]>([])
  const [filters, setFilters] = useState<Filters>({
    status: '',
    category: '',
    owner_user_id: '',
    level: '',
    search: ''
  })
  const [sortField, setSortField] = useState<SortField>('level')
  const [sortDirection, setSortDirection] = useState<SortDirection>('desc')
  const [showFilters, setShowFilters] = useState(false)

  useEffect(() => {
    console.log('RisksPage useEffect triggered')
    loadRisks()
    loadUsers()
  }, [filters, sortField, sortDirection])

  const loadUsers = async () => {
    try {
      const userData = await getUsers()
      console.log('RisksPage loadUsers - userData:', userData, 'isArray:', Array.isArray(userData))
      setUsers(Array.isArray(userData) ? userData : [])
    } catch (err) {
      console.error('Error loading users:', err)
      setUsers([])
    }
  }

  const loadRisks = async () => {
    try {
      setLoading(true)
      setError(null)
      const params = {
        ...filters,
        sort_field: sortField,
        sort_direction: sortDirection
      }
      const response = await risksApi.list(params)
      setRisks(response.data || [])
    } catch (err) {
      console.error('Error loading risks:', err)
      setError('Ошибка загрузки рисков')
      setRisks([])
    } finally {
      setLoading(false)
    }
  }

  const getRiskLevel = (level?: number) => {
    if (!level) return { color: 'default', label: 'Не определен' }
    if (level <= 2) return { color: 'success', label: 'Low', bgColor: '#4caf50' } // Зеленый
    if (level <= 4) return { color: 'warning', label: 'Medium', bgColor: '#ffeb3b' } // Желтый
    if (level <= 6) return { color: 'warning', label: 'High', bgColor: '#ff9800' } // Оранжевый
    return { color: 'error', label: 'Critical', bgColor: '#f44336' } // Красный
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

  const handleCreateRisk = async (data: CreateRiskRequest) => {
    try {
      await risksApi.create(data)
      await loadRisks()
    } catch (err) {
      console.error('Error creating risk:', err)
      setError('Ошибка создания риска')
    }
  }

  const handleUpdateRisk = async (data: UpdateRiskRequest) => {
    if (!editingRisk) return
    try {
      await risksApi.update(editingRisk.id, data)
      await loadRisks()
    } catch (err) {
      console.error('Error updating risk:', err)
      setError('Ошибка обновления риска')
    }
  }

  const handleDeleteRisk = async (risk: Risk) => {
    if (window.confirm(`Вы уверены, что хотите удалить риск "${risk.title}"?`)) {
      try {
        await risksApi.delete(risk.id)
        await loadRisks()
      } catch (err) {
        console.error('Error deleting risk:', err)
        setError('Ошибка удаления риска')
      }
    }
  }

  const handleEditRisk = (risk: Risk) => {
    setEditingRisk(risk)
    setModalOpen(true)
  }

  const handleViewRiskDetails = (risk: Risk) => {
    setSelectedRisk(risk)
    setDetailsModalOpen(true)
  }

  const handleModalClose = () => {
    setModalOpen(false)
    setEditingRisk(null)
  }

  const handleDetailsModalClose = () => {
    setDetailsModalOpen(false)
    setSelectedRisk(null)
  }

  const handleOpenCreateModal = () => {
    setEditingRisk(null)
    setModalOpen(true)
  }

  const handleFilterChange = (field: keyof Filters, value: string) => {
    setFilters(prev => ({ ...prev, [field]: value }))
  }

  const clearFilters = () => {
    setFilters({
      status: '',
      category: '',
      owner_user_id: '',
      level: '',
      search: ''
    })
  }

  const handleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection(prev => prev === 'asc' ? 'desc' : 'asc')
    } else {
      setSortField(field)
      setSortDirection('desc')
    }
  }

  const exportRisks = async (format: 'csv' | 'xlsx' | 'pdf') => {
    try {
      // TODO: Implement export functionality
      console.log(`Exporting risks in ${format} format`)
    } catch (err) {
      console.error('Error exporting risks:', err)
      setError('Ошибка экспорта рисков')
    }
  }

  const getUserName = (userId?: string) => {
    if (!userId) return 'Не назначен'
    const user = users.find(u => u.id === userId)
    return user ? `${user.first_name} ${user.last_name}` : 'Неизвестный пользователь'
  }

  const getCategoryLabel = (category?: string) => {
    if (!category) return 'Не указана'
    const cat = RISK_CATEGORIES.find(c => c.value === category)
    return cat ? cat.label : category
  }

  const getStatusLabel = (status: string) => {
    const stat = RISK_STATUSES.find(s => s.value === status)
    return stat ? stat.label : status
  }


  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Риски</Typography>
        <Box display="flex" gap={1}>
          <Button 
            variant="outlined" 
            startIcon={<FilterList />}
            onClick={() => setShowFilters(!showFilters)}
          >
            Фильтры
          </Button>
          <Button 
            variant="outlined" 
            startIcon={<Download />}
            onClick={() => exportRisks('csv')}
          >
            Экспорт
          </Button>
          <Button variant="contained" startIcon={<Add />} onClick={handleOpenCreateModal}>
            Добавить риск
          </Button>
        </Box>
      </Box>

      {error && (
        <Box mb={2} p={2} bgcolor="error.light" borderRadius={1}>
          <Typography color="error">{error}</Typography>
        </Box>
      )}

      {/* Search and Filters */}
      <Paper sx={{ mb: 2, p: 2 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={6}>
            <TextField
              fullWidth
              placeholder="Поиск по названию или описанию..."
              value={filters.search}
              onChange={(e) => handleFilterChange('search', e.target.value)}
              InputProps={{
                startAdornment: <Search sx={{ mr: 1, color: 'text.secondary' }} />
              }}
            />
          </Grid>
          <Grid item xs={12} md={2}>
            <Button
              variant="outlined"
              startIcon={<Clear />}
              onClick={clearFilters}
              fullWidth
            >
              Очистить
            </Button>
          </Grid>
        </Grid>

        {showFilters && (
          <Grid container spacing={2} sx={{ mt: 2 }}>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Статус</InputLabel>
                <Select
                  value={filters.status}
                  onChange={(e) => handleFilterChange('status', e.target.value)}
                  label="Статус"
                >
                  <MenuItem value="">Все</MenuItem>
                  {RISK_STATUSES.map((status) => (
                    <MenuItem key={status.value} value={status.value}>
                      {status.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Категория</InputLabel>
                <Select
                  value={filters.category}
                  onChange={(e) => handleFilterChange('category', e.target.value)}
                  label="Категория"
                >
                  <MenuItem value="">Все</MenuItem>
                  {RISK_CATEGORIES.map((category) => (
                    <MenuItem key={category.value} value={category.value}>
                      {category.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Владелец</InputLabel>
                <Select
                  value={filters.owner_user_id}
                  onChange={(e) => handleFilterChange('owner_user_id', e.target.value)}
                  label="Владелец"
                >
                  <MenuItem value="">Все</MenuItem>
                  {users.map((user) => (
                    <MenuItem key={user.id} value={user.id}>
                      {user.first_name} {user.last_name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Уровень риска</InputLabel>
                <Select
                  value={filters.level}
                  onChange={(e) => handleFilterChange('level', e.target.value)}
                  label="Уровень риска"
                >
                  <MenuItem value="">Все</MenuItem>
                  {RISK_LEVELS.map((level) => (
                    <MenuItem key={level.value} value={level.value.toString()}>
                      {level.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        )}
      </Paper>

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>
                  <TableSortLabel
                    active={sortField === 'title'}
                    direction={sortField === 'title' ? sortDirection : 'asc'}
                    onClick={() => handleSort('title')}
                  >
                    Название
                  </TableSortLabel>
                </TableCell>
                <TableCell>Категория</TableCell>
                <TableCell>Владелец</TableCell>
                <TableCell>Вероятность</TableCell>
                <TableCell>Воздействие</TableCell>
                <TableCell>
                  <TableSortLabel
                    active={sortField === 'level'}
                    direction={sortField === 'level' ? sortDirection : 'asc'}
                    onClick={() => handleSort('level')}
                  >
                    Уровень риска
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={sortField === 'status'}
                    direction={sortField === 'status' ? sortDirection : 'asc'}
                    onClick={() => handleSort('status')}
                  >
                    Статус
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={sortField === 'created_at'}
                    direction={sortField === 'created_at' ? sortDirection : 'asc'}
                    onClick={() => handleSort('created_at')}
                  >
                    Дата создания
                  </TableSortLabel>
                </TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={8} align="center">
                    <LinearProgress />
                    <Typography sx={{ mt: 1 }}>Загрузка рисков...</Typography>
                  </TableCell>
                </TableRow>
              ) : risks.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} align="center">
                    <Typography>Нет рисков для отображения.</Typography>
                  </TableCell>
                </TableRow>
              ) : (
                risks.map((risk) => {
                  const riskLevel = getRiskLevel(risk.level)
                  return (
                    <TableRow key={risk.id} hover>
                      <TableCell>
                        <Box display="flex" alignItems="center">
                          <Warning sx={{ mr: 1 }} />
                          <Box>
                            <Typography variant="body2" fontWeight="medium">
                              {risk.title}
                            </Typography>
                            {risk.description && (
                              <Typography variant="caption" color="text.secondary">
                                {risk.description.substring(0, 50)}...
                              </Typography>
                            )}
                          </Box>
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {getCategoryLabel(risk.category)}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {getUserName(risk.owner_user_id)}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Box display="flex" alignItems="center">
                          <LinearProgress
                            variant="determinate"
                            value={(risk.likelihood || 1) / 4 * 100}
                            sx={{ width: 60, mr: 1 }}
                          />
                          {risk.likelihood || 1}/4
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Box display="flex" alignItems="center">
                          <LinearProgress
                            variant="determinate"
                            value={(risk.impact || 1) / 4 * 100}
                            sx={{ width: 60, mr: 1 }}
                          />
                          {risk.impact || 1}/4
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={`${riskLevel.label} (${risk.level})`}
                          sx={{
                            backgroundColor: riskLevel.bgColor,
                            color: 'white',
                            fontWeight: 'bold'
                          }}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={getStatusLabel(risk.status)}
                          color={getStatusColor(risk.status) as any}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {new Date(risk.created_at).toLocaleDateString('ru-RU')}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Box display="flex" gap={0.5}>
                          <Tooltip title="Просмотр деталей">
                            <IconButton 
                              size="small"
                              color="primary"
                              onClick={() => handleViewRiskDetails(risk)}
                            >
                              <Visibility />
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Редактировать риск">
                            <IconButton 
                              size="small"
                              color="primary"
                              onClick={() => handleEditRisk(risk)}
                            >
                              <Edit />
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Удалить риск">
                            <IconButton 
                              size="small"
                              color="error"
                              onClick={() => handleDeleteRisk(risk)}
                            >
                              <Delete />
                            </IconButton>
                          </Tooltip>
                        </Box>
                      </TableCell>
                    </TableRow>
                  )
                })
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      <RiskModal
        open={modalOpen}
        onClose={handleModalClose}
        onSubmit={editingRisk ? handleUpdateRisk : handleCreateRisk}
        title={editingRisk ? 'Редактирование риска' : 'Создание нового риска'}
        initialData={editingRisk ? {
          title: editingRisk.title,
          description: editingRisk.description || '',
          category: editingRisk.category || '',
          likelihood: editingRisk.likelihood || 1,
          impact: editingRisk.impact || 1,
          status: editingRisk.status,
          owner_user_id: editingRisk.owner_user_id || '',
          methodology: editingRisk.methodology || '',
          strategy: editingRisk.strategy || '',
          due_date: editingRisk.due_date || '',
        } : undefined}
        users={users}
      />

      <RiskDetailsModal
        open={detailsModalOpen}
        onClose={handleDetailsModalClose}
        risk={selectedRisk}
      />
    </Container>
  )
}
