import React, { useEffect, useMemo, useState } from 'react'
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
  Tabs,
  Tab,
  Divider,
  IconButton,
  Tooltip,
} from '@mui/material'
import AutorenewIcon from '@mui/icons-material/Autorenew'
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
import { templatesApi } from '../../shared/api/templates'

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
  const [tabValue, setTabValue] = useState(0);
  const [generatingInventoryNumber, setGeneratingInventoryNumber] = useState(false);

  const {
    control,
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
    setValue,
    watch,
  } = useForm<AssetFormData>({
    resolver: zodResolver(schema as any),
    defaultValues: {
      name: asset?.name || '',
      inventory_number: asset?.inventory_number || '',
      type: asset?.type || '',
      class: asset?.class || '',
      owner_id: asset?.owner_id || '',
      responsible_user_id: asset?.responsible_user_id || '',
      location: asset?.location || '',
      criticality: asset?.criticality || '',
      confidentiality: asset?.confidentiality || '',
      integrity: asset?.integrity || '',
      availability: asset?.availability || '',
      status: asset?.status || 'active',
      template_id: '',
      // Passport fields
      serial_number: asset?.serial_number || '',
      pc_number: asset?.pc_number || '',
      model: asset?.model || '',
      manufacturer: asset?.manufacturer || '',
      cpu: asset?.cpu || '',
      ram: asset?.ram || '',
      hdd_info: asset?.hdd_info || '',
      network_card: asset?.network_card || '',
      optical_drive: asset?.optical_drive || '',
      ip_address: asset?.ip_address || '',
      mac_address: asset?.mac_address || '',
      purchase_year: asset?.purchase_year || undefined,
      warranty_until: asset?.warranty_until || '',
    },
  })

  const assetType = watch('type')
  const assetClass = watch('class')

  // Сбрасываем вкладку на первую, если класс актива изменился на не-hardware
  useEffect(() => {
    if (assetClass !== 'hardware' && tabValue === 1) {
      setTabValue(0)
    }
  }, [assetClass, tabValue])

  useEffect(() => {
    console.log('=== AssetModal useEffect ===');
    console.log('Asset object:', asset);
    console.log('Passport fields:', {
      serial_number: asset?.serial_number,
      pc_number: asset?.pc_number,
      model: asset?.model,
      cpu: asset?.cpu,
      ram: asset?.ram,
    });
    
    reset({
      name: asset?.name || '',
      inventory_number: asset?.inventory_number || '',
      type: asset?.type || '',
      class: asset?.class || '',
      owner_id: asset?.owner_id || '',
      responsible_user_id: asset?.responsible_user_id || '',
      location: asset?.location || '',
      criticality: asset?.criticality || '',
      confidentiality: asset?.confidentiality || '',
      integrity: asset?.integrity || '',
      availability: asset?.availability || '',
      status: asset?.status || 'active',
      template_id: '',
      // Passport fields
      serial_number: asset?.serial_number || '',
      pc_number: asset?.pc_number || '',
      model: asset?.model || '',
      manufacturer: asset?.manufacturer || '',
      cpu: asset?.cpu || '',
      ram: asset?.ram || '',
      hdd_info: asset?.hdd_info || '',
      network_card: asset?.network_card || '',
      optical_drive: asset?.optical_drive || '',
      ip_address: asset?.ip_address || '',
      mac_address: asset?.mac_address || '',
      purchase_year: asset?.purchase_year || undefined,
      warranty_until: asset?.warranty_until || '',
    })
    setTabValue(0); // Reset to first tab when modal opens
  }, [asset, reset, open])

  const handleGenerateInventoryNumber = async () => {
    if (!assetType) {
      alert('Сначала выберите тип актива')
      return
    }

    try {
      setGeneratingInventoryNumber(true)
      const response = await templatesApi.generateInventoryNumber('temp-id', {
        asset_type: assetType,
      })
      setValue('inventory_number', response.inventory_number)
    } catch (error: any) {
      alert(error.message || 'Не удалось сгенерировать инвентарный номер. Возможно, правило для этого типа актива не настроено.')
    } finally {
      setGeneratingInventoryNumber(false)
    }
  }


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
    <Dialog open={open} onClose={onClose} maxWidth="lg" fullWidth>
      <DialogTitle>{isEdit ? 'Редактирование актива' : 'Создание актива'}</DialogTitle>

      <form onSubmit={handleSubmit(onSubmit)}>
        <DialogContent sx={{ maxHeight: '70vh', overflowY: 'auto' }}>
          {assetClass === 'hardware' && (
            <Tabs value={tabValue} onChange={(_, val) => setTabValue(val)} sx={{ mb: 3 }}>
              <Tab label="Основные данные" />
              <Tab label="Данные паспорта" />
            </Tabs>
          )}

          {(assetClass !== 'hardware' || tabValue === 0) && (
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
              <TextField
                label={`Инвентарный номер${isEdit ? '' : ' *'}`}
                placeholder="Введите или сгенерируйте"
                fullWidth
                {...register('inventory_number')}
                error={!!errors.inventory_number}
                helperText={errors.inventory_number?.message}
                InputProps={{
                  endAdornment: !isEdit && (
                    <Tooltip title="Сгенерировать номер">
                      <IconButton
                        onClick={handleGenerateInventoryNumber}
                        disabled={generatingInventoryNumber || !assetType}
                        size="small"
                        color="primary"
                      >
                        <AutorenewIcon />
                      </IconButton>
                    </Tooltip>
                  ),
                }}
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
          )}

          {assetClass === 'hardware' && tabValue === 1 && (
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold', mb: 2 }}>
                Технические характеристики
              </Typography>
              <Divider sx={{ mb: 2 }} />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Серийный номер (S/N)"
                fullWidth
                {...register('serial_number')}
                placeholder="Введите серийный номер"
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Номер ПК (НО)"
                fullWidth
                {...register('pc_number')}
                placeholder="Введите номер ПК"
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Модель"
                fullWidth
                {...register('model')}
                placeholder="Например: Dell OptiPlex 7090"
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Производитель"
                fullWidth
                {...register('manufacturer')}
                placeholder="Например: Dell Inc."
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Процессор (CPU)"
                fullWidth
                {...register('cpu')}
                placeholder="Например: Intel Core i7-11700"
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Оперативная память (RAM)"
                fullWidth
                {...register('ram')}
                placeholder="Например: 16 GB DDR4"
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="HDD информация"
                fullWidth
                {...register('hdd_info')}
                placeholder="Например: SSD 512GB"
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Сетевая карта"
                fullWidth
                {...register('network_card')}
                placeholder="Например: Intel I219-V"
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Оптический привод"
                fullWidth
                {...register('optical_drive')}
                placeholder="Например: DVD-RW"
              />
            </Grid>

            <Grid item xs={12}>
              <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold', mb: 2, mt: 2 }}>
                Сетевые параметры
              </Typography>
              <Divider sx={{ mb: 2 }} />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="IP адрес"
                fullWidth
                {...register('ip_address')}
                placeholder="Например: 192.168.1.100"
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="MAC адрес"
                fullWidth
                {...register('mac_address')}
                placeholder="Например: 00:1A:2B:3C:4D:5E"
              />
            </Grid>

            <Grid item xs={12}>
              <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold', mb: 2, mt: 2 }}>
                Производственно-экономические параметры
              </Typography>
              <Divider sx={{ mb: 2 }} />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Год приобретения"
                type="number"
                fullWidth
                {...register('purchase_year', { valueAsNumber: true })}
                placeholder="Например: 2023"
                inputProps={{ min: 1900, max: 2100 }}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                label="Гарантия до"
                type="date"
                fullWidth
                {...register('warranty_until')}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
          </Grid>
          )}
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
