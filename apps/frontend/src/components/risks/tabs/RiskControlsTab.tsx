import React, { useEffect, useState } from 'react'
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
import {
  risksApi,
  RiskControl,
  CONTROL_TYPES,
  CONTROL_IMPLEMENTATION_STATUSES,
  CONTROL_EFFECTIVENESS,
} from '../../../shared/api/risks'

const controlSchema = z.object({
  control_id: z.string().uuid().optional(),
  control_name: z.string().min(1, 'Название обязательно'),
  control_type: z.enum(['preventive', 'detective', 'corrective'], {
    errorMap: () => ({ message: 'Тип обязателен' }),
  }),
  implementation_status: z.enum(['planned', 'in_progress', 'implemented', 'not_applicable'], {
    errorMap: () => ({ message: 'Статус обязателен' }),
  }),
  effectiveness: z.enum(['high', 'medium', 'low']).optional(),
  description: z.string().optional(),
})

type ControlFormData = z.infer<typeof controlSchema>

interface RiskControlsTabProps {
  riskId: string
}

const getControlTypeLabel = (type: string) =>
  CONTROL_TYPES.find((item) => item.value === type)?.label ?? type

const getImplementationStatusLabel = (status: string) =>
  CONTROL_IMPLEMENTATION_STATUSES.find((item) => item.value === status)?.label ?? status

const getEffectivenessLabel = (value?: string | null) =>
  CONTROL_EFFECTIVENESS.find((item) => item.value === value)?.label ?? 'Не указана'

export const RiskControlsTab: React.FC<RiskControlsTabProps> = ({ riskId }) => {
  const [controls, setControls] = useState<RiskControl[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
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
      control_name: '',
      control_type: 'preventive',
      implementation_status: 'planned',
      effectiveness: undefined,
      description: '',
    },
  })

  useEffect(() => {
    void loadControls()
  }, [riskId])

  const loadControls = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await risksApi.getControls(riskId)
      setControls(response)
    } catch (err) {
      console.error('Error loading controls:', err)
      setError('Не удалось загрузить контроли')
      setControls([])
    } finally {
      setLoading(false)
    }
  }

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, controlItem: RiskControl) => {
    setAnchorEl(event.currentTarget)
    setSelectedControl(controlItem)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
    setSelectedControl(null)
  }

  const handleEdit = () => {
    if (!selectedControl) return
    setEditingControl(selectedControl)
    reset({
      control_id: selectedControl.control_id,
      control_name: selectedControl.control_name,
          control_type: selectedControl.control_type as 'preventive' | 'detective' | 'corrective',
          implementation_status: selectedControl.implementation_status as 'planned' | 'in_progress' | 'implemented' | 'not_applicable',
      effectiveness: selectedControl.effectiveness as 'high' | 'medium' | 'low' | undefined,
      description: selectedControl.description ?? '',
    })
    setModalOpen(true)
    handleMenuClose()
  }

  const handleDelete = async () => {
    if (!selectedControl) return
    try {
      await risksApi.deleteControl(riskId, selectedControl.id)
      await loadControls()
    } catch (err) {
      console.error('Error deleting control:', err)
      setError('Не удалось удалить контроль')
    } finally {
      handleMenuClose()
    }
  }

  const handleCreateNew = () => {
    setEditingControl(null)
    reset({
      control_id: undefined,
      control_name: '',
      control_type: 'preventive',
      implementation_status: 'planned',
      effectiveness: undefined,
      description: '',
    })
    setModalOpen(true)
  }

  const onSubmit = async (data: ControlFormData) => {
    try {
      setError(null)
      if (editingControl) {
        await risksApi.updateControl(riskId, editingControl.id, {
          control_name: data.control_name,
          control_type: data.control_type,
          implementation_status: data.implementation_status,
          effectiveness: data.effectiveness || undefined,
          description: data.description?.trim() || undefined,
        })
      } else {
        const controlId = data.control_id ?? (globalThis.crypto?.randomUUID?.() ?? Math.random().toString(36).slice(2))
        await risksApi.createControl(riskId, {
          control_id: controlId,
          control_name: data.control_name,
          control_type: data.control_type,
          implementation_status: data.implementation_status,
          effectiveness: data.effectiveness || undefined,
          description: data.description?.trim() || undefined,
        })
      }

      setModalOpen(false)
      setEditingControl(null)
      await loadControls()
    } catch (err) {
      console.error('Error saving control:', err)
      setError('Не удалось сохранить контроль')
    }
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
          Контроли ({controls.length})
        </Typography>
        <Button variant="contained" startIcon={<Add />} onClick={handleCreateNew}>
          Добавить контроль
        </Button>
      </Box>

      {error && (
        <Box mb={2}>
          <Typography color="error">{error}</Typography>
        </Box>
      )}

      {controls.length === 0 ? (
        <Card>
          <CardContent>
            <Box textAlign="center" py={4}>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                Контроли ещё не добавлены
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Используйте кнопку «Добавить контроль», чтобы связать контроль с риском.
              </Typography>
            </Box>
          </CardContent>
        </Card>
      ) : (
        <Grid container spacing={2}>
          {controls.map((controlItem) => (
            <Grid item xs={12} md={6} key={controlItem.id}>
              <Card>
                <CardContent>
                  <Box display="flex" justifyContent="space-between" alignItems="flex-start" mb={2}>
                    <Box display="flex" alignItems="center" gap={1}>
                      <Security color="primary" />
                      <Typography variant="h6">{controlItem.control_name}</Typography>
                    </Box>
                    <IconButton onClick={(event) => handleMenuOpen(event, controlItem)}>
                      <MoreVert />
                    </IconButton>
                  </Box>

                  {controlItem.description && (
                    <Typography variant="body2" color="text.secondary" mb={2}>
                      {controlItem.description}
                    </Typography>
                  )}

                  <Grid container spacing={1}>
                    <Grid item>
                      <Chip label={getControlTypeLabel(controlItem.control_type)} variant="outlined" />
                    </Grid>
                    <Grid item>
                      <Chip
                        label={getImplementationStatusLabel(controlItem.implementation_status)}
                        color={controlItem.implementation_status === 'implemented' ? 'success' : 'default'}
                        variant="outlined"
                      />
                    </Grid>
                    <Grid item>
                      <Chip label={`Эффективность: ${getEffectivenessLabel(controlItem.effectiveness)}`} variant="outlined" />
                    </Grid>
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleMenuClose}>
        <MenuItem onClick={handleEdit}>
          <Edit sx={{ mr: 1 }} /> Редактировать
        </MenuItem>
        <MenuItem onClick={handleDelete}>
          <Delete sx={{ mr: 1 }} /> Удалить
        </MenuItem>
      </Menu>

      <Dialog open={modalOpen} onClose={() => setModalOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>{editingControl ? 'Редактирование контроля' : 'Добавление контроля'}</DialogTitle>
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogContent>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Controller
                  name="control_name"
                  control={control}
                  render={({ field }) => (
                    <TextField
                      {...field}
                      label="Название контроля"
                      fullWidth
                      error={!!errors.control_name}
                      helperText={errors.control_name?.message}
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
                  name="control_type"
                  control={control}
                  render={({ field }) => (
                    <FormControl fullWidth error={!!errors.control_type}>
                      <InputLabel>Тип контроля</InputLabel>
                      <Select {...field} label="Тип контроля">
                        {CONTROL_TYPES.map((item) => (
                          <MenuItem key={item.value} value={item.value}>
                            {item.label}
                          </MenuItem>
                        ))}
                      </Select>
                    </FormControl>
                  )}
                />
              </Grid>

              <Grid item xs={12} sm={6}>
                <Controller
                  name="implementation_status"
                  control={control}
                  render={({ field }) => (
                    <FormControl fullWidth error={!!errors.implementation_status}>
                      <InputLabel>Статус внедрения</InputLabel>
                      <Select {...field} label="Статус внедрения">
                        {CONTROL_IMPLEMENTATION_STATUSES.map((item) => (
                          <MenuItem key={item.value} value={item.value}>
                            {item.label}
                          </MenuItem>
                        ))}
                      </Select>
                    </FormControl>
                  )}
                />
              </Grid>

              <Grid item xs={12} sm={6}>
                <Controller
                  name="effectiveness"
                  control={control}
                  render={({ field }) => (
                    <FormControl fullWidth>
                      <InputLabel>Эффективность</InputLabel>
                      <Select
                        {...field}
                        label="Эффективность"
                        value={field.value ?? ''}
                        onChange={(event) => field.onChange(event.target.value || undefined)}
                      >
                        <MenuItem value="">Не указана</MenuItem>
                        {CONTROL_EFFECTIVENESS.map((item) => (
                          <MenuItem key={item.value} value={item.value}>
                            {item.label}
                          </MenuItem>
                        ))}
                      </Select>
                    </FormControl>
                  )}
                />
              </Grid>
            </Grid>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setModalOpen(false)}>Отмена</Button>
            <Button type="submit" variant="contained">
              {editingControl ? 'Обновить' : 'Создать'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  )
}

