import React, { useEffect, useMemo, useState } from 'react'
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
  Chip,
} from '@mui/material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  RISK_STATUSES,
  RISK_METHODOLOGIES,
  RISK_STRATEGIES,
  RISK_CATEGORIES,
} from '../../shared/api/risks'
import type { User } from '../../shared/api/users'

const riskSchema = z.object({
  title: z.string().min(1, 'Название обязательно'),
  description: z.string().optional(),
  category: z.string().min(1, 'Выберите категорию'),
  likelihood: z.number().min(1).max(4),
  impact: z.number().min(1).max(4),
  status: z.string().min(1, 'Статус обязателен'),
  owner_user_id: z.string().optional(),
  methodology: z.string().optional(),
  strategy: z.string().optional(),
  due_date: z.string().optional(),
})

export type RiskFormData = z.infer<typeof riskSchema>

interface RiskModalProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: RiskFormData) => void
  title: string
  initialData?: Partial<RiskFormData>
  users?: User[]
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

const buildUserLabel = (user: User) => {
  const fullName = [user.first_name, user.last_name].filter(Boolean).join(' ').trim()
  return fullName ? (user.email ? `${fullName} (${user.email})` : fullName) : user.email
}

const buildUserKey = (user: User, index: number) => {
  if (user.id) return `owner-${user.id}`
  if (user.email) return `owner-email-${user.email}`
  return `owner-index-${index}`
}

const normalizeDate = (value?: string | null) => {
  if (!value) return ''
  try {
    return new Date(value).toISOString().split('T')[0]
  } catch (err) {
    return ''
  }
}

export const RiskModal: React.FC<RiskModalProps> = ({
  open,
  onClose,
  onSubmit,
  title,
  initialData,
  users = [],
}) => {
  const [isSubmitting, setIsSubmitting] = useState(false)
  
  const {
    control,
    handleSubmit,
    register,
    reset,
    watch,
    formState: { errors },
  } = useForm<RiskFormData>({
    resolver: zodResolver(riskSchema),
    defaultValues: {
      title: initialData?.title ?? '',
      description: initialData?.description ?? '',
      category: initialData?.category ?? '',
      likelihood: initialData?.likelihood ?? 1,
      impact: initialData?.impact ?? 1,
      status: initialData?.status ?? 'new',
      owner_user_id: initialData?.owner_user_id ?? '',
      methodology: initialData?.methodology ?? '',
      strategy: initialData?.strategy ?? '',
      due_date: normalizeDate(initialData?.due_date),
    },
  })

  useEffect(() => {
    reset({
      title: initialData?.title ?? '',
      description: initialData?.description ?? '',
      category: initialData?.category ?? '',
      likelihood: initialData?.likelihood ?? 1,
      impact: initialData?.impact ?? 1,
      status: initialData?.status ?? 'new',
      owner_user_id: initialData?.owner_user_id ?? '',
      methodology: initialData?.methodology ?? '',
      strategy: initialData?.strategy ?? '',
      due_date: normalizeDate(initialData?.due_date),
    })
  }, [initialData, reset])

  const sliderMarks = useMemo(
    () => [
      { value: 1, label: '1' },
      { value: 2, label: '2' },
      { value: 3, label: '3' },
      { value: 4, label: '4' },
    ],
    [],
  )

  const handleFormSubmit = async (data: RiskFormData) => {
    try {
      setIsSubmitting(true)
      await onSubmit({
        ...data,
        description: data.description?.trim() || undefined,
        category: data.category?.trim() || '',
        owner_user_id: data.owner_user_id?.trim() || undefined,
        methodology: data.methodology?.trim() || undefined,
        strategy: data.strategy?.trim() || undefined,
        due_date: data.due_date?.trim() || undefined,
      })
    } catch (error) {
      console.error('Error submitting form:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const likelihood = watch('likelihood')
  const impact = watch('impact')
  const level = likelihood * impact

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
      <DialogTitle>{title}</DialogTitle>
      <DialogContent dividers>
        <Grid container spacing={3} sx={{ mt: 0.5 }}>
          <Grid item xs={12}>
            <TextField
              label="Название риска"
              fullWidth
              {...register('title')}
              error={!!errors.title}
              helperText={errors.title?.message}
            />
          </Grid>

          <Grid item xs={12}>
            <TextField
              label="Описание"
              fullWidth
              multiline
              minRows={3}
              {...register('description')}
              error={!!errors.description}
              helperText={errors.description?.message}
            />
          </Grid>

          <Grid item xs={12} sm={6}>
            <Controller
              name="category"
              control={control}
              render={({ field }) => (
                <FormControl fullWidth error={!!errors.category}>
                  <InputLabel>Категория</InputLabel>
                  <Select
                    {...field}
                    label="Категория"
                    value={field.value ?? ''}
                  >
                    {RISK_CATEGORIES.map((option) => (
                      <MenuItem key={option.value} value={option.value}>
                        {option.label}
                      </MenuItem>
                    ))}
                  </Select>
                  {errors.category && (
                    <Typography variant="caption" color="error">
                      {errors.category.message}
                    </Typography>
                  )}
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
                  <Select
                    {...field}
                    label="Статус"
                    value={field.value ?? 'new'}
                  >
                    {RISK_STATUSES.map((option) => (
                      <MenuItem key={option.value} value={option.value}>
                        {STATUS_LABELS[option.value] ?? option.label ?? option.value}
                      </MenuItem>
                    ))}
                  </Select>
                  {errors.status && (
                    <Typography variant="caption" color="error">
                      {errors.status.message}
                    </Typography>
                  )}
                </FormControl>
              )}
            />
          </Grid>

          <Grid item xs={12} md={6}>
            <Typography variant="subtitle2" gutterBottom>
              Вероятность
            </Typography>
            <Controller
              name="likelihood"
              control={control}
              render={({ field }) => (
                <Slider
                  {...field}
                  value={field.value ?? 1}
                  min={1}
                  max={4}
                  step={1}
                  marks={sliderMarks}
                  valueLabelDisplay="auto"
                />
              )}
            />
          </Grid>

          <Grid item xs={12} md={6}>
            <Typography variant="subtitle2" gutterBottom>
              Влияние
            </Typography>
            <Controller
              name="impact"
              control={control}
              render={({ field }) => (
                <Slider
                  {...field}
                  value={field.value ?? 1}
                  min={1}
                  max={4}
                  step={1}
                  marks={sliderMarks}
                  valueLabelDisplay="auto"
                />
              )}
            />
          </Grid>

          <Grid item xs={12}>
            <Box display="flex" justifyContent="space-between" alignItems="center">
              <Typography variant="subtitle2">Уровень риска: {level}</Typography>
              <Chip
                label={
                  level <= 4
                    ? 'Низкий'
                    : level <= 8
                    ? 'Средний'
                    : level <= 12
                    ? 'Высокий'
                    : 'Критический'
                }
                color={level <= 4 ? 'success' : level <= 8 ? 'warning' : 'error'}
                size="small"
              />
            </Box>
          </Grid>

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
                    value={field.value ?? ''}
                  >
                    <MenuItem value="">
                      <em>Не назначен</em>
                    </MenuItem>
                    {users.map((user, index) => (
                      <MenuItem key={buildUserKey(user, index)} value={user.id}>
                        {buildUserLabel(user)}
                      </MenuItem>
                    ))}
                  </Select>
                  {errors.owner_user_id && (
                    <Typography variant="caption" color="error">
                      {errors.owner_user_id.message}
                    </Typography>
                  )}
                </FormControl>
              )}
            />
          </Grid>

          <Grid item xs={12} sm={6}>
            <Controller
              name="methodology"
              control={control}
              render={({ field }) => (
                <FormControl fullWidth>
                  <InputLabel>Метод оценки</InputLabel>
                  <Select
                    {...field}
                    label="Метод оценки"
                    value={field.value ?? ''}
                  >
                    <MenuItem value="">
                      <em>Не выбран</em>
                    </MenuItem>
                    {RISK_METHODOLOGIES.map((option) => (
                      <MenuItem key={option.value} value={option.value}>
                        {option.label}
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
                <FormControl fullWidth>
                  <InputLabel>Стратегия обработки</InputLabel>
                  <Select
                    {...field}
                    label="Стратегия обработки"
                    value={field.value ?? ''}
                  >
                    <MenuItem value="">
                      <em>Не выбрана</em>
                    </MenuItem>
                    {RISK_STRATEGIES.map((option) => (
                      <MenuItem key={option.value} value={option.value}>
                        {option.label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              )}
            />
          </Grid>

          <Grid item xs={12} sm={6}>
            <TextField
              label="Срок завершения"
              type="date"
              fullWidth
              InputLabelProps={{ shrink: true }}
              {...register('due_date')}
              defaultValue={normalizeDate(initialData?.due_date)}
            />
          </Grid>
        </Grid>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} disabled={isSubmitting}>
          Отмена
        </Button>
        <Button 
          onClick={handleSubmit(handleFormSubmit)} 
          variant="contained"
          disabled={isSubmitting}
        >
          {isSubmitting ? 'Сохранение...' : 'Сохранить'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
