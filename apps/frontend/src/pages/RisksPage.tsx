import React, { useEffect, useMemo, useState } from 'react'
import {
  Box,
  Button,
  Chip,
  FormControl,
  Grid,
  IconButton,
  InputAdornment,
  InputLabel,
  LinearProgress,
  MenuItem,
  Paper,
  Select,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TableSortLabel,
  TextField,
  Tooltip,
  Typography,
} from '@mui/material'
import { alpha, useTheme } from '@mui/material/styles'
import type { SxProps, Theme } from '@mui/material'
import {
  Add,
  Clear,
  Delete,
  Download,
  Edit,
  FilterList,
  Search,
  Visibility,
} from '@mui/icons-material'
import Pagination from '../components/Pagination'
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
} from '../shared/api/risks'
import { getUsers, type User } from '../shared/api/users'
import type { CorporateTheme } from '../shared/theme'
import { PageContainer, PageHeader, SectionCard } from '../components/common/Page'

type SortField = 'level' | 'created_at' | 'category' | 'title' | 'status'
type SortDirection = 'asc' | 'desc'

type Filters = {
  status: string
  category: string
  owner_user_id: string
  level: string
  search: string
}

type RiskBadgeConfig = {
  label: string
  sx: SxProps<Theme>
}

const STATUS_LABELS: Record<string, string> = {
  new: 'Новый',
  in_analysis: 'На анализе',
  in_treatment: 'В обработке',
  accepted: 'Принят',
  transferred: 'Передан',
  mitigated: 'Снижен',
  closed: 'Закрыт',
}

const LEVEL_LABELS: Record<'low' | 'medium' | 'high' | 'critical', string> = {
  low: 'Низкий',
  medium: 'Средний',
  high: 'Высокий',
  critical: 'Критический',
}

const formatUserName = (user: User) => {
  const names = [user.first_name, user.last_name].filter(Boolean).join(' ').trim()
  return names ? (user.email ? `${names} (${user.email})` : names) : user.email
}

export const RisksPage: React.FC = () => {
  const [risks, setRisks] = useState<Risk[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingRisk, setEditingRisk] = useState<Risk | null>(null)
  const [detailsOpen, setDetailsOpen] = useState(false)
  const [selectedRisk, setSelectedRisk] = useState<Risk | null>(null)
  const [users, setUsers] = useState<User[]>([])
  const [filters, setFilters] = useState<Filters>({ status: '', category: '', owner_user_id: '', level: '', search: '' })
  const [sortField, setSortField] = useState<SortField>('level')
  const [sortDirection, setSortDirection] = useState<SortDirection>('desc')
  const [showFilters, setShowFilters] = useState(false)
  const [pagination, setPagination] = useState({
    page: 1,
    page_size: 20,
    total: 0,
    total_pages: 0,
    has_next: false,
    has_prev: false,
  })

  const theme = useTheme<CorporateTheme>()

  const riskLevelTokens = useMemo<Record<'low' | 'medium' | 'high' | 'critical' | 'unknown', RiskBadgeConfig>>(() => ({
    low: {
      label: LEVEL_LABELS.low,
      sx: {
        bgcolor: alpha(theme.palette.success.main, 0.16),
        color: theme.palette.success.dark,
      },
    },
    medium: {
      label: LEVEL_LABELS.medium,
      sx: {
        bgcolor: alpha(theme.palette.warning.light ?? theme.palette.warning.main, 0.18),
        color: theme.palette.warning.dark ?? theme.palette.warning.main,
      },
    },
    high: {
      label: LEVEL_LABELS.high,
      sx: {
        bgcolor: alpha(theme.palette.warning.main, 0.22),
        color: theme.palette.common.white,
      },
    },
    critical: {
      label: LEVEL_LABELS.critical,
      sx: {
        bgcolor: alpha(theme.palette.error.main, 0.28),
        color: theme.palette.common.white,
      },
    },
    unknown: {
      label: 'Не определён',
      sx: {
        bgcolor: theme.palette.background.default,
        color: theme.palette.text.secondary,
        borderWidth: 1,
        borderStyle: 'dashed',
        borderColor: theme.palette.divider,
      },
    },
  }), [theme])

  useEffect(() => {
    void loadRisks()
  }, [filters, sortField, sortDirection, pagination.page])

  useEffect(() => {
    void loadUsers()
  }, [])

  const loadUsers = async () => {
    try {
      const userData = await getUsers()
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
        sort_direction: sortDirection,
        page: pagination.page,
        page_size: pagination.page_size,
      }
      const response = await risksApi.list(params)
      setRisks(response.data || [])
      if (response.pagination) {
        setPagination(response.pagination)
      }
    } catch (err) {
      console.error('Error loading risks:', err)
      setError('Не удалось загрузить риски')
      setRisks([])
    } finally {
      setLoading(false)
    }
  }

  const getRiskLevelBadge = (risk: Risk) => {
    const raw = risk.level_label ?? (() => {
      if (!risk.level) return 'unknown'
      if (risk.level <= 2) return 'low'
      if (risk.level <= 4) return 'medium'
      if (risk.level <= 6) return 'high'
      return 'critical'
    })()
    const key = raw.toLowerCase() as keyof typeof riskLevelTokens
    return riskLevelTokens[key] || riskLevelTokens.unknown
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'new':
        return 'info'
      case 'in_analysis':
        return 'warning'
      case 'in_treatment':
        return 'primary'
      case 'accepted':
        return 'success'
      case 'transferred':
        return 'secondary'
      case 'mitigated':
        return 'success'
      case 'closed':
        return 'default'
      default:
        return 'default'
    }
  }

  const handleCreateRisk = async (payload: CreateRiskRequest) => {
    try {
      await risksApi.create(payload)
      await loadRisks()
      setModalOpen(false)
      setEditingRisk(null)
      setError(null)
    } catch (err) {
      console.error('Error creating risk:', err)
      setError('Не удалось создать риск')
    }
  }

  const handleUpdateRisk = async (payload: UpdateRiskRequest) => {
    if (!editingRisk) return
    try {
      await risksApi.update(editingRisk.id, payload)
      await loadRisks()
      setModalOpen(false)
      setEditingRisk(null)
      setError(null)
    } catch (err) {
      console.error('Error updating risk:', err)
      setError('Не удалось обновить риск')
    }
  }

  const handleDeleteRisk = async (risk: Risk) => {
    if (!window.confirm(`Удалить риск "${risk.title}"?`)) return
    try {
      await risksApi.delete(risk.id)
      await loadRisks()
    } catch (err) {
      console.error('Error deleting risk:', err)
      setError('Не удалось удалить риск')
    }
  }

  const handleFilterChange = <T extends keyof Filters>(key: T, value: Filters[T]) => {
    setFilters((prev) => ({ ...prev, [key]: value }))
    setPagination((prev) => ({ ...prev, page: 1 }))
  }

  const handleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection((prev) => (prev === 'asc' ? 'desc' : 'asc'))
    } else {
      setSortField(field)
      setSortDirection('asc')
    }
  }

  const hasActiveFilters = Object.values(filters).some(Boolean)

  const getUserName = (id?: string | null) => {
    if (!id) return 'Не назначен'
    const user = users.find((u) => u.id === id)
    return user ? formatUserName(user) : 'Не назначен'
  }

  return (
    <PageContainer maxWidth="xl">
      <PageHeader
        title="Реестр рисков"
        subtitle="Контроль владельцев, уровней и статусов обработки в режиме реального времени"
        actions={
          <Button variant="contained" startIcon={<Add />} onClick={() => setModalOpen(true)}>
            Добавить риск
          </Button>
        }
      />

      <SectionCard
        title="Фильтры"
        description="Ограничьте список по статусу, категории, ответственному или ключевым словам"
        action={
          <Button startIcon={<FilterList />} onClick={() => setShowFilters((prev) => !prev)}>
            {showFilters ? 'Скрыть фильтры' : 'Показать фильтры'}
          </Button>
        }
      >
        {showFilters && (
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel id="filter-status">Статус</InputLabel>
                <Select
                  labelId="filter-status"
                  value={filters.status}
                  label="Статус"
                  onChange={(event) => handleFilterChange('status', event.target.value)}
                >
                  <MenuItem value="">Все статусы</MenuItem>
                  {RISK_STATUSES.map((item) => (
                    <MenuItem key={item.value} value={item.value}>
                      {STATUS_LABELS[item.value] ?? item.label ?? item.value}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel id="filter-category">Категория</InputLabel>
                <Select
                  labelId="filter-category"
                  value={filters.category}
                  label="Категория"
                  onChange={(event) => handleFilterChange('category', event.target.value)}
                >
                  <MenuItem value="">Все категории</MenuItem>
                  {RISK_CATEGORIES.map((item) => (
                    <MenuItem key={item.value} value={item.value}>
                      {item.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel id="filter-owner">Ответственный</InputLabel>
                <Select
                  labelId="filter-owner"
                  value={filters.owner_user_id}
                  label="Ответственный"
                  onChange={(event) => handleFilterChange('owner_user_id', event.target.value)}
                >
                  <MenuItem value="">Все ответственные</MenuItem>
                  {users.map((user, index) => (
                    <MenuItem key={user.id || index} value={user.id}>
                      {formatUserName(user)}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel id="filter-level">Уровень</InputLabel>
                <Select
                  labelId="filter-level"
                  value={filters.level}
                  label="Уровень"
                  onChange={(event) => handleFilterChange('level', event.target.value)}
                >
                  <MenuItem value="">Все уровни</MenuItem>
                  {RISK_LEVELS.map((item) => (
                    <MenuItem key={item.value} value={item.value}>
                      {item.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                size="small"
                label="Поиск по названию или описанию"
                value={filters.search}
                onChange={(event) => handleFilterChange('search', event.target.value)}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <Search fontSize="small" />
                    </InputAdornment>
                  ),
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <Button
                variant="text"
                startIcon={<Clear />}
                onClick={() => setFilters({ status: '', category: '', owner_user_id: '', level: '', search: '' })}
                disabled={!hasActiveFilters}
              >
                Сбросить фильтры
              </Button>
            </Grid>
          </Grid>
        )}
      </SectionCard>

      <SectionCard title="Список рисков" description="Экспозиция, ответственные и статусы обработки">
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="body2" color="text.secondary">
            Всего: {pagination.total} рисков
          </Typography>
          <Button variant="outlined" startIcon={<Download />}>
            Экспорт CSV
          </Button>
        </Box>

        <TableContainer component={Paper}>
          <Table size="medium">
            <TableHead>
              <TableRow>
                <TableCell>
                  <TableSortLabel
                    active={sortField === 'title'}
                    direction={sortDirection}
                    onClick={() => handleSort('title')}
                  >
                    Риск
                  </TableSortLabel>
                </TableCell>
                <TableCell>Категория</TableCell>
                <TableCell>Ответственный</TableCell>
                <TableCell>Вероятность</TableCell>
                <TableCell>Влияние</TableCell>
                <TableCell>Уровень</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Создан</TableCell>
                <TableCell align="right">Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={9}>
                    <Box display="flex" justifyContent="center" py={6}>
                      <LinearProgress sx={{ width: 320 }} />
                    </Box>
                  </TableCell>
                </TableRow>
              ) : risks.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={9}>
                    <Box textAlign="center" py={6}>
                      <Typography variant="body1" fontWeight={600}>
                        Риски не найдены
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Измените фильтры или добавьте новую запись.
                      </Typography>
                    </Box>
                  </TableCell>
                </TableRow>
              ) : (
                risks.map((risk) => {
                  const levelBadge = getRiskLevelBadge(risk)

                  return (
                    <TableRow key={risk.id} hover>
                      <TableCell>
                        <Box display="flex" flexDirection="column">
                          <Typography variant="body1" fontWeight={600}>
                            {risk.title}
                          </Typography>
                          {risk.description && (
                            <Typography variant="body2" color="text.secondary">
                              {risk.description.slice(0, 96)}{risk.description.length > 96 ? '…' : ''}
                            </Typography>
                          )}
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {RISK_CATEGORIES.find((item) => item.value === risk.category)?.label ?? '—'}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">{getUserName(risk.owner_user_id)}</Typography>
                      </TableCell>
                      <TableCell>
                        <Box display="flex" alignItems="center">
                          <LinearProgress
                            variant="determinate"
                            value={Math.min(((risk.likelihood ?? 0) / 4) * 100, 100)}
                            sx={{ width: 60, mr: 1 }}
                          />
                          {(risk.likelihood ?? 0)}/4
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Box display="flex" alignItems="center">
                          <LinearProgress
                            variant="determinate"
                            value={Math.min(((risk.impact ?? 0) / 4) * 100, 100)}
                            sx={{ width: 60, mr: 1 }}
                          />
                          {(risk.impact ?? 0)}/4
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={`${levelBadge.label}${risk.level ? ` (${risk.level})` : ''}`}
                          size="small"
                          sx={{
                            fontWeight: 600,
                            letterSpacing: '0.01em',
                            borderRadius: 2,
                            ...levelBadge.sx,
                          }}
                        />
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={STATUS_LABELS[risk.status] ?? risk.status}
                          color={getStatusColor(risk.status) as any}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">
                          {new Date(risk.created_at).toLocaleDateString('ru-RU')}
                        </Typography>
                      </TableCell>
                      <TableCell align="right">
                        <Box display="flex" justifyContent="flex-end" gap={0.5}>
                          <Tooltip title="Просмотреть">
                            <IconButton size="small" color="primary" onClick={() => {
                              setSelectedRisk(risk)
                              setDetailsOpen(true)
                            }}>
                              <Visibility fontSize="small" />
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Редактировать">
                            <IconButton size="small" color="primary" onClick={() => {
                              setEditingRisk(risk)
                              setModalOpen(true)
                            }}>
                              <Edit fontSize="small" />
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Удалить">
                            <IconButton size="small" color="error" onClick={() => handleDeleteRisk(risk)}>
                              <Delete fontSize="small" />
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

        <Box display="flex" justifyContent="center" mt={3}>
          <Pagination
            currentPage={pagination.page}
            totalPages={pagination.total_pages}
            hasNext={pagination.has_next}
            hasPrev={pagination.has_prev}
            onPageChange={(page) => setPagination((prev) => ({ ...prev, page }))}
          />
        </Box>
      </SectionCard>

      <RiskModal
        open={modalOpen}
        onClose={() => {
          setModalOpen(false)
          setEditingRisk(null)
        }}
        onSubmit={editingRisk ? (payload) => handleUpdateRisk(payload) : handleCreateRisk}
        title={editingRisk ? 'Редактирование риска' : 'Создание риска'}
        initialData={editingRisk ? {
          title: editingRisk.title,
          description: editingRisk.description ?? undefined,
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
        open={detailsOpen}
        onClose={() => {
          setDetailsOpen(false)
          setSelectedRisk(null)
        }}
        risk={selectedRisk}
      />

      {error && (
        <Box mt={2}>
          <Typography color="error">{error}</Typography>
        </Box>
      )}
    </PageContainer>
  )
}
