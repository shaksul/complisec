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
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  IconButton,
  InputAdornment,
  TableSortLabel,
  CircularProgress,
  Alert,
} from '@mui/material'
import { 
  Add, 
  Search, 
  Visibility, 
  Edit, 
  Refresh
} from '@mui/icons-material'
import { getUserCatalog, getUsersPaginated, UserCatalog, UserCatalogParams, PaginatedResponse, User } from '../shared/api/users'
import Pagination from '../components/Pagination'
import { UserDetailModal } from '../components/users/UserDetailModal'

export const UsersPage: React.FC = () => {
  const [users, setUsers] = useState<UserCatalog[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [pagination, setPagination] = useState({
    page: 1,
    page_size: 20,
    total: 0,
    total_pages: 0,
    has_next: false,
    has_prev: false,
  })
  
  // Фильтры
  const [filters, setFilters] = useState<UserCatalogParams>({
    page: 1,
    page_size: 20,
    search: '',
    role: '',
    is_active: undefined,
    sort_by: 'created_at',
    sort_dir: 'desc',
  })
  
  const [selectedUser, setSelectedUser] = useState<UserCatalog | null>(null)
  const [detailModalOpen, setDetailModalOpen] = useState(false)

  const loadUsers = async () => {
    try {
      setLoading(true)
      setError(null)
      const response: PaginatedResponse<UserCatalog> = await getUserCatalog(filters)
      setUsers(response.data)
      setPagination(response.pagination)
    } catch (err) {
      setError('Ошибка загрузки пользователей')
      console.error('Error loading users:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadUsers()
  }, [filters])

  const handleSearch = (value: string) => {
    setFilters(prev => ({ ...prev, search: value, page: 1 }))
  }

  const handleRoleFilter = (value: string) => {
    setFilters(prev => ({ ...prev, role: value, page: 1 }))
  }

  const handleStatusFilter = (value: string) => {
    const isActive = value === 'all' ? undefined : value === 'active'
    setFilters(prev => ({ ...prev, is_active: isActive, page: 1 }))
  }

  const handleSort = (field: string) => {
    const newSortDir = filters.sort_by === field && filters.sort_dir === 'asc' ? 'desc' : 'asc'
    setFilters(prev => ({ ...prev, sort_by: field, sort_dir: newSortDir, page: 1 }))
  }

  const handlePageChange = (page: number) => {
    setFilters(prev => ({ ...prev, page }))
  }

  const handleViewUser = (user: UserCatalog) => {
    setSelectedUser(user)
    setDetailModalOpen(true)
  }

  const handleRefresh = () => {
    loadUsers()
  }

  const getStatusChip = (isActive: boolean) => (
    <Chip
      label={isActive ? 'Активен' : 'Заблокирован'}
      color={isActive ? 'success' : 'error'}
      size="small"
    />
  )

  const getRoleChips = (roles: string[]) => (
    <Box display="flex" gap={0.5} flexWrap="wrap">
      {roles.map((role, index) => (
        <Chip key={index} label={role} size="small" variant="outlined" />
      ))}
    </Box>
  )

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Пользователи</Typography>
        <Box display="flex" gap={2}>
          <Button
            variant="outlined"
            startIcon={<Refresh />}
            onClick={handleRefresh}
            disabled={loading}
          >
            Обновить
          </Button>
          <Button variant="contained" startIcon={<Add />}>
            Добавить пользователя
          </Button>
        </Box>
      </Box>

      {/* Фильтры */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <Box display="flex" gap={2} alignItems="center" flexWrap="wrap">
          <TextField
            placeholder="Поиск по email, имени, фамилии..."
            value={filters.search}
            onChange={(e) => handleSearch(e.target.value)}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <Search />
                </InputAdornment>
              ),
            }}
            sx={{ minWidth: 300 }}
          />
          
          <FormControl sx={{ minWidth: 150 }}>
            <InputLabel>Роль</InputLabel>
            <Select
              value={filters.role}
              onChange={(e) => handleRoleFilter(e.target.value)}
              label="Роль"
            >
              <MenuItem value="">Все роли</MenuItem>
              <MenuItem value="Admin">Администратор</MenuItem>
              <MenuItem value="Manager">Менеджер</MenuItem>
              <MenuItem value="User">Пользователь</MenuItem>
            </Select>
          </FormControl>

          <FormControl sx={{ minWidth: 150 }}>
            <InputLabel>Статус</InputLabel>
            <Select
              value={filters.is_active === undefined ? 'all' : filters.is_active ? 'active' : 'inactive'}
              onChange={(e) => handleStatusFilter(e.target.value)}
              label="Статус"
            >
              <MenuItem value="all">Все</MenuItem>
              <MenuItem value="active">Активные</MenuItem>
              <MenuItem value="inactive">Заблокированные</MenuItem>
            </Select>
          </FormControl>
        </Box>
      </Paper>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>
                  <TableSortLabel
                    active={filters.sort_by === 'email'}
                    direction={filters.sort_by === 'email' ? filters.sort_dir : 'asc'}
                    onClick={() => handleSort('email')}
                  >
                    Email
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={filters.sort_by === 'first_name'}
                    direction={filters.sort_by === 'first_name' ? filters.sort_dir : 'asc'}
                    onClick={() => handleSort('first_name')}
                  >
                    Имя
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={filters.sort_by === 'last_name'}
                    direction={filters.sort_by === 'last_name' ? filters.sort_dir : 'asc'}
                    onClick={() => handleSort('last_name')}
                  >
                    Фамилия
                  </TableSortLabel>
                </TableCell>
                <TableCell>Роли</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>
                  <TableSortLabel
                    active={filters.sort_by === 'created_at'}
                    direction={filters.sort_by === 'created_at' ? filters.sort_dir : 'asc'}
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
                  <TableCell colSpan={7} align="center">
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : users.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center">
                    Пользователи не найдены
                  </TableCell>
                </TableRow>
              ) : (
                users.map((user) => (
                  <TableRow key={user.id} hover>
                    <TableCell>{user.email}</TableCell>
                    <TableCell>{user.first_name || '-'}</TableCell>
                    <TableCell>{user.last_name || '-'}</TableCell>
                    <TableCell>{getRoleChips(user.roles)}</TableCell>
                    <TableCell>{getStatusChip(user.is_active)}</TableCell>
                    <TableCell>
                      {new Date(user.created_at).toLocaleDateString('ru-RU')}
                    </TableCell>
                    <TableCell>
                      <Box display="flex" gap={1}>
                        <IconButton
                          size="small"
                          onClick={() => handleViewUser(user)}
                          title="Просмотр"
                        >
                          <Visibility />
                        </IconButton>
                        <IconButton
                          size="small"
                          title="Редактировать"
                        >
                          <Edit />
                        </IconButton>
                      </Box>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
        
        {pagination.total_pages > 1 && (
          <Box p={2}>
            <Pagination
              currentPage={pagination.page}
              totalPages={pagination.total_pages}
              hasNext={pagination.has_next}
              hasPrev={pagination.has_prev}
              onPageChange={handlePageChange}
            />
          </Box>
        )}
      </Paper>

      {selectedUser && (
        <UserDetailModal
          open={detailModalOpen}
          onClose={() => setDetailModalOpen(false)}
          user={selectedUser}
        />
      )}
    </Container>
  )
}
