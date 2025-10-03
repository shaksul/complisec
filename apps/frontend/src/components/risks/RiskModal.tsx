import React from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Box,
  Typography,
  Slider,
  Grid,
} from '@mui/material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { 
  RISK_STATUSES, 
  RISK_METHODOLOGIES, 
  RISK_STRATEGIES, 
  RISK_CATEGORIES 
} from '../../shared/api/risks'
import { User } from '../../shared/api/users'

const riskSchema = z.object({
  title: z.string().min(1, 'Название обязательно'),
  description: z.string().optional(),
  category: z.string().min(1, 'Категория обязательна'),
  likelihood: z.number().min(1).max(4),
  impact: z.number().min(1).max(4),
  status: z.string().min(1, 'Статус обязателен'),
  owner_user_id: z.string().optional(),
  methodology: z.string().optional(),
  strategy: z.string().optional(),
  due_date: z.string().optional(),
})

type RiskFormData = z.infer<typeof riskSchema>

interface RiskModalProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: RiskFormData) => void
  title: string
  initialData?: Partial<RiskFormData>
  users?: User[]
}

export const RiskModal: React.FC<RiskModalProps> = ({
  open,
  onClose,
  onSubmit,
  title,
  initialData,
  users,
}) => {
  console.log('RiskModal rendered, open:', open, 'users:', users, 'users length:', users?.length, 'users type:', typeof users, 'isArray:', Array.isArray(users))

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
    watch,
  } = useForm<RiskFormData>({
    resolver: zodResolver(riskSchema),
    defaultValues: {
      title: '',
      description: '',
      category: '',
      likelihood: 1,
      impact: 1,
      status: 'new',
      owner_user_id: '',
      methodology: '',
      strategy: '',
      due_date: '',
      ...initialData,
    },
  })

  const likelihood = watch('likelihood')
  const impact = watch('impact')
  const riskLevel = likelihood * impact


  const getRiskLevelLabel = (level: number) => {
    if (level <= 2) return { label: 'Low', color: 'success' }
    if (level <= 4) return { label: 'Medium', color: 'warning' }
    if (level <= 6) return { label: 'High', color: 'error' }
    return { label: 'Critical', color: 'error' }
  }

  const riskLevelInfo = getRiskLevelLabel(riskLevel)

  const handleFormSubmit = (data: RiskFormData) => {
    // Преобразуем пустые строки в undefined для опциональных полей
    const submitData = {
      ...data,
      description: data.description && data.description.trim() !== '' ? data.description : undefined,
      category: data.category && data.category.trim() !== '' ? data.category : undefined,
      owner_user_id: data.owner_user_id && data.owner_user_id.trim() !== '' ? data.owner_user_id : undefined,
      methodology: data.methodology && data.methodology.trim() !== '' ? data.methodology : undefined,
      strategy: data.strategy && data.strategy.trim() !== '' ? data.strategy : undefined,
      due_date: data.due_date && data.due_date.trim() !== '' ? data.due_date : undefined,
    }
    
    onSubmit(submitData)
    reset()
    onClose()
  }

  const handleClose = () => {
    reset()
    onClose()
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>{title}</DialogTitle>
      <form onSubmit={handleSubmit(handleFormSubmit)}>
        <DialogContent>
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <Controller
                name="title"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    label="Название риска"
                    fullWidth
                    error={!!errors.title}
                    helperText={errors.title?.message}
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
                name="category"
                control={control}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.category}>
                    <InputLabel>Категория</InputLabel>
                    <Select {...field} label="Категория">
                      {RISK_CATEGORIES.map((category) => (
                        <MenuItem key={category.value} value={category.value}>
                          {category.label}
                        </MenuItem>
                      ))}
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
                      {RISK_STATUSES.map((status) => (
                        <MenuItem key={status.value} value={status.value}>
                          {status.label}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Box>
                <Typography gutterBottom>
                  Вероятность: {likelihood}/4
                </Typography>
                <Controller
                  name="likelihood"
                  control={control}
                  render={({ field }) => (
                    <Slider
                      {...field}
                      min={1}
                      max={4}
                      step={1}
                      marks={[
                        { value: 1, label: '1' },
                        { value: 2, label: '2' },
                        { value: 3, label: '3' },
                        { value: 4, label: '4' }
                      ]}
                      valueLabelDisplay="auto"
                    />
                  )}
                />
              </Box>
            </Grid>

            <Grid item xs={12} sm={6}>
              <Box>
                <Typography gutterBottom>
                  Воздействие: {impact}/4
                </Typography>
                <Controller
                  name="impact"
                  control={control}
                  render={({ field }) => (
                    <Slider
                      {...field}
                      min={1}
                      max={4}
                      step={1}
                      marks={[
                        { value: 1, label: '1' },
                        { value: 2, label: '2' },
                        { value: 3, label: '3' },
                        { value: 4, label: '4' }
                      ]}
                      valueLabelDisplay="auto"
                    />
                  )}
                />
              </Box>
            </Grid>

            <Grid item xs={12}>
              <Box p={2} bgcolor="grey.100" borderRadius={1}>
                <Typography variant="h6" gutterBottom>
                  Уровень риска
                </Typography>
                <Typography variant="h4" color={`${riskLevelInfo.color}.main`}>
                  {riskLevelInfo.label} ({riskLevel})
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Вероятность: {likelihood} × Воздействие: {impact} = {riskLevel}
                </Typography>
              </Box>
            </Grid>

            {/* New fields */}
                <Grid item xs={12} sm={6}>
                  <Controller
                    name="owner_user_id"
                    control={control}
                    render={({ field }) => (
                      <FormControl fullWidth error={!!errors.owner_user_id}>
                        <InputLabel>Ответственный</InputLabel>
                        <Select
                          {...field}
                          label="Ответственный"
                        >
                          <MenuItem key="empty" value="">
                            <em>Не выбран</em>
                          </MenuItem>
                          {users && users.length > 0 ? (
                            users.map((user) => (
                              <MenuItem key={user.id} value={user.id}>
                                {user.first_name && user.last_name 
                                  ? `${user.first_name} ${user.last_name} (${user.email})`
                                  : user.email
                                }
                              </MenuItem>
                            ))
                          ) : (
                            // Fallback если пользователи не загружены
                            <>
                              <MenuItem key="1" value="1">
                                Admin User (admin@demo.local)
                              </MenuItem>
                              <MenuItem key="2" value="2">
                                John Doe (john@demo.local)
                              </MenuItem>
                            </>
                          )}
                        </Select>
                      </FormControl>
                    )}
                  />
                </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="methodology"
                control={control}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.methodology}>
                    <InputLabel>Методология</InputLabel>
                    <Select {...field} label="Методология">
                      <MenuItem value="">
                        <em>Не выбрана</em>
                      </MenuItem>
                      {RISK_METHODOLOGIES.map((methodology) => (
                        <MenuItem key={methodology.value} value={methodology.value}>
                          {methodology.label}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="strategy"
                control={control}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.strategy}>
                    <InputLabel>Стратегия обработки</InputLabel>
                    <Select {...field} label="Стратегия обработки">
                      <MenuItem value="">
                        <em>Не выбрана</em>
                      </MenuItem>
                      {RISK_STRATEGIES.map((strategy) => (
                        <MenuItem key={strategy.value} value={strategy.value}>
                          {strategy.label}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="due_date"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    label="Срок обработки"
                    type="date"
                    fullWidth
                    InputLabelProps={{
                      shrink: true,
                    }}
                    error={!!errors.due_date}
                    helperText={errors.due_date?.message}
                  />
                )}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Отмена</Button>
          <Button type="submit" variant="contained">
            Сохранить
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}
