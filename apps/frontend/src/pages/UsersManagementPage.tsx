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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  Tooltip,
  CircularProgress,
  Chip,
  Pagination,
  Checkbox,
  FormControlLabel,
  FormGroup,
} from '@mui/material'
import {
  Add,
  Edit,
  Delete,
  Person,
  Security,
} from '@mui/icons-material'
import { usersApi, User } from '../shared/api/users'
import { rolesApi, Role } from '../shared/api/roles'
import { useAuth } from '../contexts/AuthContext'
import EmailChangeModal from '../components/EmailChangeModal'

const UsersManagementPage: React.FC = () => {
  const [users, setUsers] = useState<User[]>([])
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [showRolesModal, setShowRolesModal] = useState(false)
  const [showEmailChangeModal, setShowEmailChangeModal] = useState(false)

  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [, setHasNext] = useState(false)
  const [, setHasPrev] = useState(false)

  const { user } = useAuth()

  useEffect(() => {
    if (user) {
      loadData()
    } else {
      setLoading(false)
    }
  }, [user])

  const loadData = async (page: number = currentPage) => {
    try {
      setLoading(true)
      const [usersResponse, rolesData] = await Promise.all([
        usersApi.getUsersPaginated(page, 20),
        rolesApi.getRoles()
      ])

      const normalizedUsers = usersResponse.data.map((user) => ({
        ...user,
        roles: user.roles || []
      }))

      setUsers(normalizedUsers)
      setRoles(rolesData)
      setCurrentPage(usersResponse.pagination.page)
      setTotalPages(usersResponse.pagination.total_pages)
      setHasNext(usersResponse.pagination.has_next)
      setHasPrev(usersResponse.pagination.has_prev)
    } catch (error) {
      console.error('Ошибка загрузки данных:', error)
    } finally {
      setLoading(false)
    }
  }

  const handlePageChange = (page: number) => {
    setCurrentPage(page)
    loadData(page)
  }

  const handleCreateUser = async (userData: any) => {
    try {
      await usersApi.createUser(userData)
      await loadData(1)
      setCurrentPage(1)
      setShowCreateModal(false)
    } catch (error) {
      console.error('Ошибка создания пользователя:', error)
      // Ошибка будет показана в модальном окне
      throw error
    }
  }

  const handleUpdateUser = async (id: string, userData: any) => {
    try {
      await usersApi.updateUser(id, userData)
      await loadData(currentPage)
      setShowEditModal(false)
      setSelectedUser(null)
    } catch (error) {
      console.error('Ошибка обновления пользователя:', error)
      // Ошибка будет показана в модальном окне
      throw error
    }
  }

  const handleDeleteUser = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить пользователя?')) {
      try {
        await usersApi.deleteUser(id)
        await loadData(currentPage)
      } catch (error) {
        console.error('Ошибка удаления пользователя:', error)
        alert('Ошибка удаления пользователя: ' + (error as Error).message)
      }
    }
  }

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
          <CircularProgress />
        </Box>
      </Container>
    )
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1" gutterBottom>
          Управление пользователями
        </Typography>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={() => setShowCreateModal(true)}
          sx={{ ml: 2 }}
        >
          Добавить пользователя
        </Button>
      </Box>

      <Paper sx={{ width: '100%', overflow: 'hidden' }}>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Пользователь</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Дата создания</TableCell>
                <TableCell align="center">Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {users.map((user, index) => (
                <TableRow key={user.id || `user-${index}`} hover>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      <Person sx={{ mr: 1, color: 'primary.main' }} />
                      <Box>
                        <Typography variant="body2" fontWeight="medium">
                          {[user.first_name, user.last_name].filter(Boolean).join(' ') || '—'}
                        </Typography>
                        {user.roles && user.roles.length > 0 && (
                          <Typography variant="caption" color="text.secondary">
                            {user.roles.join(', ')}
                          </Typography>
                        )}
                      </Box>
                    </Box>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2">
                      {user.email}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={user.is_active ? 'Активен' : 'Заблокирован'}
                      color={user.is_active ? 'success' : 'error'}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2">
                      {user.created_at ? new Date(user.created_at).toLocaleDateString('ru-RU') : '—'}
                    </Typography>
                  </TableCell>
                  <TableCell align="center">
                    <Tooltip title="Редактировать">
                      <IconButton
                        size="small"
                        onClick={() => {
                          setSelectedUser(user)
                          setShowEditModal(true)
                        }}
                        color="primary"
                      >
                        <Edit />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Роли">
                      <IconButton
                        size="small"
                        onClick={() => {
                          setSelectedUser(user)
                          setShowRolesModal(true)
                        }}
                        color="secondary"
                      >
                        <Security />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Удалить">
                      <IconButton
                        size="small"
                        onClick={() => handleDeleteUser(user.id)}
                        color="error"
                      >
                        <Delete />
                      </IconButton>
                    </Tooltip>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {totalPages > 1 && (
        <Box display="flex" justifyContent="center" p={2}>
          <Pagination
            count={totalPages}
            page={currentPage}
            onChange={(_, page) => handlePageChange(page)}
            color="primary"
          />
        </Box>
      )}

      {showCreateModal && (
        <CreateUserModal
          onClose={() => setShowCreateModal(false)}
          onSubmit={handleCreateUser}
          roles={roles}
        />
      )}

      {showEditModal && selectedUser && (
        <EditUserModal
          user={selectedUser}
          onClose={() => {
            setShowEditModal(false)
            setSelectedUser(null)
          }}
          onSubmit={handleUpdateUser}
          roles={roles}
          onEmailChange={() => setShowEmailChangeModal(true)}
        />
      )}

      {showRolesModal && selectedUser && (
        <UserRolesModal
          user={selectedUser}
          onClose={() => {
            setShowRolesModal(false)
            setSelectedUser(null)
          }}
          roles={roles}
        />
      )}

      {showEmailChangeModal && selectedUser && (
        <EmailChangeModal
          open={showEmailChangeModal}
          onClose={() => {
            setShowEmailChangeModal(false)
            setSelectedUser(null)
          }}
          currentEmail={selectedUser.email}
        />
      )}
    </Container>
  )
}

const CreateUserModal: React.FC<{ onClose: () => void; onSubmit: (data: any) => void; roles: Role[] }> = ({ onClose, onSubmit, roles }) => {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    role_ids: [] as string[]
  })
  const [formErrors, setFormErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)

  const validateForm = (): boolean => {
    const errors: Record<string, string> = {}

    if (!formData.email.trim()) {
      errors.email = 'Email обязателен'
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      errors.email = 'Некорректный email'
    }

    if (!formData.password) {
      errors.password = 'Пароль обязателен'
    } else if (formData.password.length < 6) {
      errors.password = 'Пароль должен содержать минимум 6 символов'
    }

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
      const submitData = {
        ...formData,
        role_ids: filteredRoleIds
      }
      await onSubmit(submitData)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Dialog open onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Создать пользователя</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Email"
            type="email"
            fullWidth
            variant="outlined"
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            error={!!formErrors.email}
            helperText={formErrors.email}
            autoComplete="email"
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Пароль"
            type="password"
            fullWidth
            variant="outlined"
            value={formData.password}
            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            error={!!formErrors.password}
            helperText={formErrors.password}
            autoComplete="new-password"
            sx={{ mb: 2 }}
          />
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
              {roles.map((role, index) => (
                <FormControlLabel
                  key={role.id || `role-${index}`}
                  control={
                    <Checkbox
                      checked={formData.role_ids.includes(role.id)}
                      onChange={(e) => {
                        if (e.target.checked && role.id) {
                          setFormData({ ...formData, role_ids: [...formData.role_ids, role.id] })
                        } else if (role.id) {
                          setFormData({
                            ...formData,
                            role_ids: formData.role_ids.filter((id) => id !== role.id)
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
            {submitting ? <CircularProgress size={20} /> : 'Создать'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

const EditUserModal: React.FC<{
  user: User;
  onClose: () => void;
  onSubmit: (id: string, data: any) => void;
  roles: Role[];
  onEmailChange: () => void;
}> = ({ user, onClose, onSubmit, roles, onEmailChange }) => {
  const [formData, setFormData] = useState({
    first_name: user.first_name || '',
    last_name: user.last_name || '',
    is_active: user.is_active,
    role_ids: [] as string[]
  });
  const [formErrors, setFormErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    // Загружаем роли пользователя
    const loadUserRoles = async () => {
      try {
        const userRoles = await rolesApi.getUserRoles(user.id);
        console.log('EditUserModal loadUserRoles - userRoles:', userRoles);
        // userRoles теперь содержит ID ролей, а не имена
        setFormData(prev => ({ ...prev, role_ids: userRoles }));
      } catch (error) {
        console.error('Ошибка загрузки ролей пользователя:', error);
      }
    };
    loadUserRoles();
  }, [user.id]);

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
    e.preventDefault();
    
    if (!validateForm()) return

    try {
      setSubmitting(true)
      // Фильтруем null значения из role_ids
      const filteredRoleIds = formData.role_ids.filter(id => id !== null && id !== undefined && id !== '')
      
      console.log('EditUserModal handleSubmit - formData.role_ids:', formData.role_ids);
      console.log('EditUserModal handleSubmit - filteredRoleIds:', filteredRoleIds);
      
      await onSubmit(user.id, {
        ...formData,
        role_ids: filteredRoleIds
      });
    } finally {
      setSubmitting(false)
    }
  };

  return (
    <Dialog open onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Редактировать пользователя</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <Box sx={{ mb: 2 }}>
            <Typography variant="subtitle2" gutterBottom>
              Email
            </Typography>
            <Box display="flex" alignItems="center" gap={1}>
              <TextField
                fullWidth
                value={user.email}
                disabled
                variant="outlined"
                size="small"
              />
              <Button
                variant="outlined"
                size="small"
                onClick={() => {
                  onEmailChange();
                  onClose();
                }}
              >
                Сменить
              </Button>
            </Box>
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
                          });
                        } else if (role.id) {
                          setFormData({
                            ...formData,
                            role_ids: formData.role_ids.filter(id => id !== role.id)
                          });
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
  );
};

// Компонент управления ролями пользователя
const UserRolesModal: React.FC<{
  user: User;
  onClose: () => void;
  roles: Role[];
}> = ({ user, onClose, roles }) => {
  const [userRoles, setUserRoles] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadUserRoles = async () => {
      try {
        const roles = await rolesApi.getUserRoles(user.id);
        setUserRoles(roles);
      } catch (error) {
        console.error('Ошибка загрузки ролей пользователя:', error);
      } finally {
        setLoading(false);
      }
    };
    loadUserRoles();
  }, [user.id]);

  const handleRoleToggle = async (roleId: string, assigned: boolean) => {
    try {
      if (assigned) {
        await rolesApi.removeRoleFromUser(user.id, roleId);
        setUserRoles(prev => prev.filter(id => id !== roleId));
      } else {
        await rolesApi.assignRoleToUser(user.id, roleId);
        setUserRoles(prev => [...prev, roleId]);
      }
    } catch (error) {
      console.error('Ошибка изменения роли:', error);
    }
  };

  if (loading) {
    return (
      <Dialog open onClose={onClose} maxWidth="sm" fullWidth>
        <DialogContent>
          <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
            <CircularProgress />
          </Box>
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Dialog open onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        Роли пользователя: {user.first_name} {user.last_name}
      </DialogTitle>
      <DialogContent>
        <FormGroup>
          {roles.map((role) => {
            const isAssigned = userRoles.includes(role.id);
            return (
              <FormControlLabel
                key={role.id}
                control={
                  <Checkbox
                    checked={isAssigned}
                    onChange={(_e) => handleRoleToggle(role.id, isAssigned)}
                  />
                }
                label={
                  <Box>
                    <Typography variant="body2" fontWeight="medium">
                      {role.name}
                    </Typography>
                    {role.description && (
                      <Typography variant="caption" color="text.secondary">
                        {role.description}
                      </Typography>
                    )}
                  </Box>
                }
                sx={{ mb: 1 }}
              />
            );
          })}
        </FormGroup>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} variant="contained">
          Закрыть
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default UsersManagementPage;
