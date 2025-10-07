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
  Alert,
  CircularProgress,
  Chip,
  Pagination,
} from '@mui/material'
import {
  Add,
  Edit,
  Delete,
  Business,
} from '@mui/icons-material'
import { tenantsApi, Tenant, CreateTenantDTO, UpdateTenantDTO } from '../shared/api/tenants'
import { useAuth } from '../contexts/AuthContext'

const OrganizationsPage: React.FC = () => {
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [, setTotal] = useState(0)
  const [pageSize] = useState(20)

  // Modal states
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [selectedTenant, setSelectedTenant] = useState<Tenant | null>(null)

  // Form states
  const [formData, setFormData] = useState<CreateTenantDTO>({
    name: '',
    domain: '',
  })
  const [formErrors, setFormErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)

  const { } = useAuth()

  useEffect(() => {
    loadTenants()
  }, [currentPage])

  const loadTenants = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await tenantsApi.getTenants(currentPage, pageSize)
      setTenants(response.data)
      setTotalPages(response.pagination.total_pages)
      setTotal(response.pagination.total)
    } catch (err) {
      setError('Ошибка загрузки организаций: ' + (err as Error).message)
    } finally {
      setLoading(false)
    }
  }

  const handlePageChange = (_: React.ChangeEvent<unknown>, page: number) => {
    setCurrentPage(page)
  }

  const handleCreateTenant = () => {
    setFormData({ name: '', domain: '' })
    setFormErrors({})
    setShowCreateModal(true)
  }

  const handleEditTenant = (tenant: Tenant) => {
    setSelectedTenant(tenant)
    setFormData({
      name: tenant.name,
      domain: tenant.domain || '',
    })
    setFormErrors({})
    setShowEditModal(true)
  }

  const handleDeleteTenant = async (tenant: Tenant) => {
    if (!window.confirm(`Вы уверены, что хотите удалить организацию "${tenant.name}"?`)) {
      return
    }

    try {
      await tenantsApi.deleteTenant(tenant.id)
      await loadTenants()
    } catch (err) {
      setError('Ошибка удаления организации: ' + (err as Error).message)
    }
  }

  const validateForm = (): boolean => {
    const errors: Record<string, string> = {}

    if (!formData.name.trim()) {
      errors.name = 'Название организации обязательно'
    } else if (formData.name.length < 2) {
      errors.name = 'Название должно содержать минимум 2 символа'
    }

    if (formData.domain && formData.domain.length > 0) {
      if (formData.domain.length < 3) {
        errors.domain = 'Домен должен содержать минимум 3 символа'
      }
    }

    setFormErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async () => {
    if (!validateForm()) return

    try {
      setSubmitting(true)
      
      if (showCreateModal) {
        const createData: CreateTenantDTO = {
          name: formData.name.trim(),
          domain: formData.domain?.trim() || undefined,
        }
        await tenantsApi.createTenant(createData)
        setShowCreateModal(false)
      } else if (showEditModal && selectedTenant) {
        const updateData: UpdateTenantDTO = {
          name: formData.name.trim(),
          domain: formData.domain?.trim() || undefined,
        }
        await tenantsApi.updateTenant(selectedTenant.id, updateData)
        setShowEditModal(false)
        setSelectedTenant(null)
      }

      await loadTenants()
    } catch (err) {
      setError('Ошибка сохранения: ' + (err as Error).message)
    } finally {
      setSubmitting(false)
    }
  }

  const handleCloseModal = () => {
    setShowCreateModal(false)
    setShowEditModal(false)
    setSelectedTenant(null)
    setFormData({ name: '', domain: '' })
    setFormErrors({})
  }

  if (loading && tenants.length === 0) {
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
          Управление организациями
        </Typography>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={handleCreateTenant}
          sx={{ ml: 2 }}
        >
          Добавить организацию
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <Paper sx={{ width: '100%', overflow: 'hidden' }}>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Название</TableCell>
                <TableCell>Домен</TableCell>
                <TableCell>Дата создания</TableCell>
                <TableCell>Дата обновления</TableCell>
                <TableCell align="center">Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {tenants.map((tenant) => (
                <TableRow key={tenant.id} hover>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      <Business sx={{ mr: 1, color: 'primary.main' }} />
                      <Typography variant="body2" fontWeight="medium">
                        {tenant.name}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>
                    {tenant.domain ? (
                      <Chip label={tenant.domain} size="small" variant="outlined" />
                    ) : (
                      <Typography variant="body2" color="text.secondary">
                        Не указан
                      </Typography>
                    )}
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2">
                      {new Date(tenant.created_at).toLocaleDateString('ru-RU')}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2">
                      {new Date(tenant.updated_at).toLocaleDateString('ru-RU')}
                    </Typography>
                  </TableCell>
                  <TableCell align="center">
                    <Tooltip title="Редактировать">
                      <IconButton
                        size="small"
                        onClick={() => handleEditTenant(tenant)}
                        color="primary"
                      >
                        <Edit />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Удалить">
                      <IconButton
                        size="small"
                        onClick={() => handleDeleteTenant(tenant)}
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

        {totalPages > 1 && (
          <Box display="flex" justifyContent="center" p={2}>
            <Pagination
              count={totalPages}
              page={currentPage}
              onChange={handlePageChange}
              color="primary"
            />
          </Box>
        )}
      </Paper>

      {/* Create Modal */}
      <Dialog open={showCreateModal} onClose={handleCloseModal} maxWidth="sm" fullWidth>
        <DialogTitle>Создать организацию</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Название организации"
            fullWidth
            variant="outlined"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            error={!!formErrors.name}
            helperText={formErrors.name}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Домен (необязательно)"
            fullWidth
            variant="outlined"
            value={formData.domain}
            onChange={(e) => setFormData({ ...formData, domain: e.target.value })}
            error={!!formErrors.domain}
            helperText={formErrors.domain || 'Например: company.com'}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseModal}>Отмена</Button>
          <Button
            onClick={handleSubmit}
            variant="contained"
            disabled={submitting}
          >
            {submitting ? <CircularProgress size={20} /> : 'Создать'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Edit Modal */}
      <Dialog open={showEditModal} onClose={handleCloseModal} maxWidth="sm" fullWidth>
        <DialogTitle>Редактировать организацию</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Название организации"
            fullWidth
            variant="outlined"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            error={!!formErrors.name}
            helperText={formErrors.name}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Домен (необязательно)"
            fullWidth
            variant="outlined"
            value={formData.domain}
            onChange={(e) => setFormData({ ...formData, domain: e.target.value })}
            error={!!formErrors.domain}
            helperText={formErrors.domain || 'Например: company.com'}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseModal}>Отмена</Button>
          <Button
            onClick={handleSubmit}
            variant="contained"
            disabled={submitting}
          >
            {submitting ? <CircularProgress size={20} /> : 'Сохранить'}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  )
}

export default OrganizationsPage
