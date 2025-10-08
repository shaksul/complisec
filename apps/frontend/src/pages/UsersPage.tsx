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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Checkbox,
  FormControlLabel,
  FormGroup,
} from '@mui/material'
import { 
  Add, 
  Search, 
  Visibility, 
  Edit, 
  Refresh
} from '@mui/icons-material'
import { getUserCatalog, UserCatalog, UserCatalogParams, PaginatedResponse, usersApi } from '../shared/api/users'
import { rolesApi, Role } from '../shared/api/roles'
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
  const [editModalOpen, setEditModalOpen] = useState(false)
  const [roles, setRoles] = useState<Role[]>([])

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
    loadRoles()
  }, [filters])

  const loadRoles = async () => {
    try {
      const rolesData = await rolesApi.getRoles()
      setRoles(rolesData)
    } catch (err) {
      console.error('Error loading roles:', err)
    }
  }

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

  const handleEditUser = (user: UserCatalog) => {
    setSelectedUser(user)
    setEditModalOpen(true)
  }

  const handleUpdateUser = async (id: string, userData: any) => {
    try {
      await usersApi.updateUser(id, userData)
      await loadUsers()
      setEditModalOpen(false)
      setSelectedUser(null)
    } catch (err) {
      console.error('Error updating user:', err)
      throw err
    }
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
                          onClick={() => handleEditUser(user)}
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

      {editModalOpen && selectedUser && (
        <EditUserModal
          user={selectedUser}
          roles={roles}
          onClose={() => {
            setEditModalOpen(false)
            setSelectedUser(null)
          }}
          onSubmit={handleUpdateUser}
        />
      )}
    </Container>
  )
}

// Компонент модального окна редактирования пользователя
const EditUserModal: React.FC<{
  user: UserCatalog
  roles: Role[]
  onClose: () => void
  onSubmit: (id: string, data: any) => void
}> = ({ user, roles, onClose, onSubmit }) => {
  const [formData, setFormData] = useState({
    first_name: user.first_name || '',
    last_name: user.last_name || '',
    is_active: user.is_active,
    role_ids: [] as string[]
  })
  const [formErrors, setFormErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    const loadUserRoles = async () => {
      try {
        const userRoles = await rolesApi.getUserRoles(user.id)
        setFormData(prev => ({ ...prev, role_ids: userRoles }))
      } catch (error) {
        console.error('Ошибка загрузки ролей пользователя:', error)
      }
    }
    loadUserRoles()
  }, [user.id])

  const validateForm = (): boolean => {
    const errors: Record<string, string> = {}

    if (!formData.first_name.trim()) {
      errors.first_name = 'Имя обязательно'
    }

    if (!formData.last_name.trim()) {
      errors.last_name = 'Фамилия обязательна'
    }

    const filteredRoleIds = formData.role_ids.filter(id => id !== null && id !== undefined && id !== '')
    if (filteredRoleIds.length === 0) {
      errors.roles = 'Выберите хотя бы одну роль'
    }

    setFormErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) return

    try {
      setSubmitting(true)
      const filteredRoleIds = formData.role_ids.filter(id => id !== null && id !== undefined && id !== '')
      
      await onSubmit(user.id, {
        ...formData,
        role_ids: filteredRoleIds
      })
    } catch (error) {
      console.error('Ошибка сохранения:', error)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Dialog open onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Редактировать пользователя</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <Box sx={{ mb: 2 }}>
            <Typography variant="subtitle2" gutterBottom>
              Email
            </Typography>
            <TextField
              fullWidth
              value={user.email}
              disabled
              variant="outlined"
              size="small"
            />
          </Box>
          <TextField
            margin="dense"
            label="Имя"
            fullWidth
            variant="outlined"
            value={formData.first_name}
            onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
            error={!!formErrors.first_name}
            helperText={formErrors.first_name}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Фамилия"
            fullWidth
            variant="outlined"
            value={formData.last_name}
            onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
            error={!!formErrors.last_name}
            helperText={formErrors.last_name}
            sx={{ mb: 2 }}
          />
          <FormControlLabel
            control={
              <Checkbox
                checked={formData.is_active}
                onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
              />
            }
            label="Активен"
            sx={{ mb: 2 }}
          />
          <Box sx={{ mb: 2 }}>
            <Typography variant="subtitle2" gutterBottom>
              Роли
            </Typography>
            {formErrors.roles && (
              <Typography variant="caption" color="error">
                {formErrors.roles}
              </Typography>
            )}
            <FormGroup>
              {roles.map((role) => (
                <FormControlLabel
                  key={role.id}
                  control={
                    <Checkbox
                      checked={formData.role_ids.includes(role.id)}
                      onChange={(e) => {
                        if (e.target.checked && role.id) {
                          setFormData({
                            ...formData,
                            role_ids: [...formData.role_ids, role.id]
                          })
                        } else if (role.id) {
                          setFormData({
                            ...formData,
                            role_ids: formData.role_ids.filter(id => id !== role.id)
                          })
                        }
                      }}
                    />
                  }
                  label={role.name}
                />
              ))}
            </FormGroup>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Отмена</Button>
          <Button
            type="submit"
            variant="contained"
            disabled={submitting}
          >
            {submitting ? <CircularProgress size={20} /> : 'Сохранить'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}
