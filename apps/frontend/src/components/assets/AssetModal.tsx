import React, { useEffect, useMemo } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
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
  Typography,
  Grid,
  CircularProgress,
} from '@mui/material'
import {
  Asset,
  ASSET_TYPES,
  ASSET_CLASSES,
  CRITICALITY_LEVELS,
  ASSET_STATUSES,
} from '../../shared/api/assets'
import type { UserCatalog } from '../../shared/api/users'
import {
  createAssetSchema,
  updateAssetSchema,
  type CreateAssetFormData,
  type UpdateAssetFormData,
} from '../../shared/validation/assets'

interface AssetModalProps {
  open: boolean
  asset?: Asset | null
  users: UserCatalog[]
  onSave: (data: CreateAssetFormData | Partial<UpdateAssetFormData>) => void
  onClose: () => void
}

type AssetFormData = CreateAssetFormData

type SelectOption = { value: string; label: string }

const typeOptions: SelectOption[] = ASSET_TYPES
const classOptions: SelectOption[] = ASSET_CLASSES
const statusOptions: SelectOption[] = ASSET_STATUSES
const ciaOptions: SelectOption[] = CRITICALITY_LEVELS

const buildUserLabel = (user: UserCatalog) => {
  const fullName = [user.first_name, user.last_name].filter(Boolean).join(' ').trim()
  return fullName ? `${fullName} (${user.email})` : user.email
}

const buildUserKey = (user: UserCatalog, index: number) => user.id || `user-${index}`

const AssetModal: React.FC<AssetModalProps> = ({ open, asset, users, onSave, onClose }) => {
  const isEdit = Boolean(asset)
  const schema = useMemo(() => (isEdit ? updateAssetSchema : createAssetSchema), [isEdit])

  const {
    control,
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<AssetFormData>({
    resolver: zodResolver(schema as any),
    defaultValues: {
      name: asset?.name ?? '',
      type: asset?.type ?? '',
      class: asset?.class ?? '',
      owner_id: asset?.owner_id ?? '',
      responsible_user_id: asset?.responsible_user_id ?? '',
      location: asset?.location ?? '',
      criticality: asset?.criticality ?? '',
      confidentiality: asset?.confidentiality ?? '',
      integrity: asset?.integrity ?? '',
      availability: asset?.availability ?? '',
      status: asset?.status ?? 'active',
    },
  })

  useEffect(() => {
    reset({
      name: asset?.name ?? '',
      type: asset?.type ?? '',
      class: asset?.class ?? '',
      owner_id: asset?.owner_id ?? '',
      responsible_user_id: asset?.responsible_user_id ?? '',
      location: asset?.location ?? '',
      criticality: asset?.criticality ?? '',
      confidentiality: asset?.confidentiality ?? '',
      integrity: asset?.integrity ?? '',
      availability: asset?.availability ?? '',
      status: asset?.status ?? 'active',
    })
  }, [asset, reset])

  const onSubmit = (formData: AssetFormData) => {
    if (isEdit) {
      const payload: Partial<UpdateAssetFormData> = {}
      Object.entries(formData).forEach(([key, value]) => {
        if (value !== undefined && value !== null && `${value}`.trim() !== '') {
          payload[key as keyof UpdateAssetFormData] = value as never
        }
      })
      onSave(payload)
    } else {
      onSave(formData)
    }
  }

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>{isEdit ? 'Редактирование актива' : 'Создание актива'}</DialogTitle>

      <form onSubmit={handleSubmit(onSubmit)}>
        <DialogContent>
          <Grid container spacing={3}>
            <Grid item xs={12} sm={6}>
              <TextField
                label={`Название${isEdit ? '' : ' *'}`}
                placeholder="Введите название актива"
                fullWidth
                {...register('name')}
                error={!!errors.name}
                helperText={errors.name?.message}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="type"
                control={control}
                defaultValue={asset?.type ?? ''}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.type}>
                    <InputLabel>Тип актива</InputLabel>
                    <Select {...field} label="Тип актива" value={field.value ?? ''}>
                      <MenuItem value="">
                        <em>Не выбран</em>
                      </MenuItem>
                      {typeOptions.map((option) => (
                        <MenuItem key={option.value} value={option.value}>
                          {option.label}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.type && (
                      <Typography variant="caption" color="error">
                        {errors.type.message as string}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="class"
                control={control}
                defaultValue={asset?.class ?? ''}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.class}>
                    <InputLabel>Класс актива</InputLabel>
                    <Select {...field} label="Класс актива" value={field.value ?? ''}>
                      <MenuItem value="">
                        <em>Не выбран</em>
                      </MenuItem>
                      {classOptions.map((option) => (
                        <MenuItem key={option.value} value={option.value}>
                          {option.label}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.class && (
                      <Typography variant="caption" color="error">
                        {errors.class.message as string}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="owner_id"
                control={control}
                defaultValue={asset?.owner_id ?? ''}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.owner_id}>
                    <InputLabel>Владелец (бизнес)</InputLabel>
                    <Select {...field} label="Владелец (бизнес)" value={field.value ?? ''}>
                      <MenuItem value="">
                        <em>Не назначен</em>
                      </MenuItem>
                      {users.map((user, index) => (
                        <MenuItem key={buildUserKey(user, index)} value={user.id}>
                          {buildUserLabel(user)}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.owner_id && (
                      <Typography variant="caption" color="error">
                        {errors.owner_id.message as string}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="responsible_user_id"
                control={control}
                defaultValue={asset?.responsible_user_id ?? ''}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.responsible_user_id}>
                    <InputLabel>Ответственный (ИТ)</InputLabel>
                    <Select {...field} label="Ответственный (ИТ)" value={field.value ?? ''}>
                      <MenuItem value="">
                        <em>Не назначен</em>
                      </MenuItem>
                      {users.map((user, index) => (
                        <MenuItem key={`responsible-${buildUserKey(user, index)}`} value={user.id}>
                          {buildUserLabel(user)}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.responsible_user_id && (
                      <Typography variant="caption" color="error">
                        {errors.responsible_user_id.message as string}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Расположение"
                placeholder="Например: Москва, дата-центр"
                fullWidth
                {...register('location')}
                error={!!errors.location}
                helperText={errors.location?.message}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="criticality"
                control={control}
                defaultValue={asset?.criticality ?? ''}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.criticality}>
                    <InputLabel>Критичность</InputLabel>
                    <Select {...field} label="Критичность" value={field.value ?? ''}>
                      <MenuItem value="">
                        <em>Не выбрано</em>
                      </MenuItem>
                      {ciaOptions.map((option) => (
                        <MenuItem key={`criticality-${option.value}`} value={option.value}>
                          {option.label}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.criticality && (
                      <Typography variant="caption" color="error">
                        {errors.criticality.message as string}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="confidentiality"
                control={control}
                defaultValue={asset?.confidentiality ?? ''}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.confidentiality}>
                    <InputLabel>Конфиденциальность</InputLabel>
                    <Select {...field} label="Конфиденциальность" value={field.value ?? ''}>
                      <MenuItem value="">
                        <em>Не выбрано</em>
                      </MenuItem>
                      {ciaOptions.map((option) => (
                        <MenuItem key={`conf-${option.value}`} value={option.value}>
                          {option.label}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.confidentiality && (
                      <Typography variant="caption" color="error">
                        {errors.confidentiality.message as string}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="integrity"
                control={control}
                defaultValue={asset?.integrity ?? ''}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.integrity}>
                    <InputLabel>Целостность</InputLabel>
                    <Select {...field} label="Целостность" value={field.value ?? ''}>
                      <MenuItem value="">
                        <em>Не выбрано</em>
                      </MenuItem>
                      {ciaOptions.map((option) => (
                        <MenuItem key={`integrity-${option.value}`} value={option.value}>
                          {option.label}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.integrity && (
                      <Typography variant="caption" color="error">
                        {errors.integrity.message as string}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <Controller
                name="availability"
                control={control}
                defaultValue={asset?.availability ?? ''}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.availability}>
                    <InputLabel>Доступность</InputLabel>
                    <Select {...field} label="Доступность" value={field.value ?? ''}>
                      <MenuItem value="">
                        <em>Не выбрано</em>
                      </MenuItem>
                      {ciaOptions.map((option) => (
                        <MenuItem key={`availability-${option.value}`} value={option.value}>
                          {option.label}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.availability && (
                      <Typography variant="caption" color="error">
                        {errors.availability.message as string}
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
                defaultValue={asset?.status ?? 'active'}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.status}>
                    <InputLabel>Статус</InputLabel>
                    <Select {...field} label="Статус" value={field.value ?? 'active'}>
                      {statusOptions.map((option) => (
                        <MenuItem key={option.value} value={option.value}>
                          {option.label}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.status && (
                      <Typography variant="caption" color="error">
                        {errors.status.message as string}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>
          </Grid>
        </DialogContent>

        <DialogActions>
          <Button onClick={onClose} disabled={isSubmitting}>
            Отмена
          </Button>
          <Button type="submit" variant="contained" disabled={isSubmitting}
            startIcon={isSubmitting ? <CircularProgress size={20} /> : null}
          >
            {isSubmitting ? 'Сохраняем…' : isEdit ? 'Сохранить изменения' : 'Создать актив'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

export default AssetModal
