import React, { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Button,
  Card,
  CardContent,
  Chip,
  Grid,
  IconButton,
  Menu,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  Select,
  CircularProgress,
} from '@mui/material'
import { Add, MoreVert, Edit, Delete, Security } from '@mui/icons-material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

const controlSchema = z.object({
  name: z.string().min(1, 'Название обязательно'),
  description: z.string().optional(),
  type: z.string().min(1, 'Тип обязателен'),
  effectiveness: z.number().min(1).max(4),
  status: z.string().min(1, 'Статус обязателен'),
})

type ControlFormData = z.infer<typeof controlSchema>

interface RiskControl {
  id: string
  risk_id: string
  name: string
  description?: string
  type: string
  effectiveness: number
  status: string
  created_at: string
  updated_at: string
}

interface RiskControlsTabProps {
  riskId: string
}

export const RiskControlsTab: React.FC<RiskControlsTabProps> = ({ riskId }) => {
  const [controls, setControls] = useState<RiskControl[]>([])
  const [loading, setLoading] = useState(true)
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [selectedControl, setSelectedControl] = useState<RiskControl | null>(null)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingControl, setEditingControl] = useState<RiskControl | null>(null)

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<ControlFormData>({
    resolver: zodResolver(controlSchema),
    defaultValues: {
      name: '',
      description: '',
      type: '',
      effectiveness: 1,
      status: 'planned',
    },
  })

  useEffect(() => {
    loadControls()
  }, [riskId])

  const loadControls = async () => {
    try {
      setLoading(true)
      // TODO: Implement API call to load controls
      // const response = await riskControlsApi.list(riskId)
      // setControls(response.data || [])
      
      // Mock data for now
      setControls([
        {
          id: '1',
          risk_id: riskId,
          name: 'Многофакторная аутентификация',
          description: 'Внедрение MFA для всех пользователей',
          type: 'preventive',
          effectiveness: 4,
          status: 'implemented',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: '2',
          risk_id: riskId,
          name: 'Мониторинг безопасности',
          description: 'Система непрерывного мониторинга',
          type: 'detective',
          effectiveness: 3,
          status: 'planned',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
      ])
    } catch (err) {
      console.error('Error loading controls:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, control: RiskControl) => {
    setAnchorEl(event.currentTarget)
    setSelectedControl(control)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
    setSelectedControl(null)
  }

  const handleEdit = () => {
    if (selectedControl) {
      setEditingControl(selectedControl)
      reset({
        name: selectedControl.name,
        description: selectedControl.description || '',
        type: selectedControl.type,
        effectiveness: selectedControl.effectiveness,
        status: selectedControl.status,
      })
      setModalOpen(true)
    }
    handleMenuClose()
  }

  const handleDelete = async () => {
    if (selectedControl) {
      try {
        // TODO: Implement API call to delete control
        // await riskControlsApi.delete(selectedControl.id)
        setControls(prev => prev.filter(c => c.id !== selectedControl.id))
      } catch (err) {
        console.error('Error deleting control:', err)
      }
    }
    handleMenuClose()
  }

  const handleCreateNew = () => {
    setEditingControl(null)
    reset({
      name: '',
      description: '',
      type: '',
      effectiveness: 1,
      status: 'planned',
    })
    setModalOpen(true)
  }

  const onSubmit = async (data: ControlFormData) => {
    try {
      if (editingControl) {
        // TODO: Implement API call to update control
        // await riskControlsApi.update(editingControl.id, data)
        setControls(prev => prev.map(c => 
          c.id === editingControl.id 
            ? { ...c, ...data, updated_at: new Date().toISOString() }
            : c
        ))
      } else {
        // TODO: Implement API call to create control
        // const response = await riskControlsApi.create({ ...data, risk_id: riskId })
        const newControl: RiskControl = {
          id: Date.now().toString(),
          risk_id: riskId,
          ...data,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }
        setControls(prev => [...prev, newControl])
      }
      setModalOpen(false)
      setEditingControl(null)
    } catch (err) {
      console.error('Error saving control:', err)
    }
  }

  const getControlTypeLabel = (type: string) => {
    switch (type) {
      case 'preventive': return 'Предупреждающий'
      case 'detective': return 'Обнаруживающий'
      case 'corrective': return 'Корректирующий'
      case 'compensating': return 'Компенсирующий'
      default: return type
    }
  }

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'planned': return 'Запланирован'
      case 'implemented': return 'Реализован'
      case 'testing': return 'Тестирование'
      case 'operational': return 'Эксплуатация'
      default: return status
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'planned': return 'warning'
      case 'implemented': return 'success'
      case 'testing': return 'info'
      case 'operational': return 'primary'
      default: return 'default'
    }
  }

  const getEffectivenessColor = (effectiveness: number) => {
    if (effectiveness >= 4) return 'success'
    if (effectiveness >= 3) return 'warning'
    if (effectiveness >= 2) return 'error'
    return 'default'
  }

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h6">
          Контроли риска ({controls.length})
        </Typography>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={handleCreateNew}
        >
          Добавить контроль
        </Button>
      </Box>

      {controls.length === 0 ? (
        <Card>
          <CardContent>
            <Box textAlign="center" py={4}>
              <Security sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" color="text.secondary" gutterBottom>
                Контроли не найдены
              </Typography>
              <Typography variant="body2" color="text.secondary" mb={3}>
                Добавьте контроли для управления данным риском
              </Typography>
              <Button
                variant="contained"
                startIcon={<Add />}
                onClick={handleCreateNew}
              >
                Добавить первый контроль
              </Button>
            </Box>
          </CardContent>
        </Card>
      ) : (
        <Grid container spacing={2}>
          {controls.map((control) => (
            <Grid item xs={12} key={control.id}>
              <Card>
                <CardContent>
                  <Box display="flex" justifyContent="space-between" alignItems="flex-start">
                    <Box flex={1}>
                      <Typography variant="h6" gutterBottom>
                        {control.name}
                      </Typography>
                      
                      {control.description && (
                        <Typography variant="body2" color="text.secondary" paragraph>
                          {control.description}
                        </Typography>
                      )}

                      <Box display="flex" gap={1} flexWrap="wrap" mt={2}>
                        <Chip
                          label={getControlTypeLabel(control.type)}
                          size="small"
                          color="primary"
                          variant="outlined"
                        />
                        <Chip
                          label={getStatusLabel(control.status)}
                          size="small"
                          color={getStatusColor(control.status) as any}
                        />
                        <Chip
                          label={`Эффективность: ${control.effectiveness}/4`}
                          size="small"
                          color={getEffectivenessColor(control.effectiveness) as any}
                        />
                      </Box>
                    </Box>

                    <IconButton
                      onClick={(e) => handleMenuOpen(e, control)}
                    >
                      <MoreVert />
                    </IconButton>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      {/* Context Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={handleEdit}>
          <Edit sx={{ mr: 1 }} />
          Редактировать
        </MenuItem>
        <MenuItem onClick={handleDelete}>
          <Delete sx={{ mr: 1 }} />
          Удалить
        </MenuItem>
      </Menu>

      {/* Add/Edit Control Modal */}
      <Dialog open={modalOpen} onClose={() => setModalOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          {editingControl ? 'Редактирование контроля' : 'Добавление контроля'}
        </DialogTitle>
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogContent>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Controller
                  name="name"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      label="Название контроля"
                      fullWidth
                      error={!!errors.name}
                      helperText={errors.name?.message}
                    />
                  )}
                />
              </Grid>

              <Grid item xs={12}>
                <Controller
                  name="description"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      label="Описание"
                      fullWidth
                      multiline
                      rows={3}
                      error={!!errors.description}
                      helperText={errors.description?.message}
                    />
                  )}
                />
              </Grid>

              <Grid item xs={12} sm={6}>
                <Controller
                  name="type"
                  control={control}
                  render={({ field }) => (
                    <FormControl fullWidth error={!!errors.type}>
                      <InputLabel>Тип контроля</InputLabel>
                      <Select {...field} label="Тип контроля">
                        <MenuItem value="preventive">Предупреждающий</MenuItem>
                        <MenuItem value="detective">Обнаруживающий</MenuItem>
                        <MenuItem value="corrective">Корректирующий</MenuItem>
                        <MenuItem value="compensating">Компенсирующий</MenuItem>
                      </Select>
                    </FormControl>
                  )}
                />
              </Grid>

              <Grid item xs={12} sm={6}>
                <Controller
                  name="status"
                  control={control}
                  render={({ field }) => (
                    <FormControl fullWidth error={!!errors.status}>
                      <InputLabel>Статус</InputLabel>
                      <Select {...field} label="Статус">
                        <MenuItem value="planned">Запланирован</MenuItem>
                        <MenuItem value="testing">Тестирование</MenuItem>
                        <MenuItem value="implemented">Реализован</MenuItem>
                        <MenuItem value="operational">Эксплуатация</MenuItem>
                      </Select>
                    </FormControl>
                  )}
                />
              </Grid>

              <Grid item xs={12}>
                <Controller
                  name="effectiveness"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      label="Эффективность (1-4)"
                      type="number"
                      inputProps={{ min: 1, max: 4 }}
                      fullWidth
                      error={!!errors.effectiveness}
                      helperText={errors.effectiveness?.message}
                    />
                  )}
                />
              </Grid>
            </Grid>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setModalOpen(false)}>
              Отмена
            </Button>
            <Button type="submit" variant="contained">
              {editingControl ? 'Обновить' : 'Создать'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  )
}

