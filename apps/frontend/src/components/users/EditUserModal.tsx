import React, { useState, useEffect } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  CircularProgress,
  Checkbox,
  FormControlLabel,
  FormGroup,
  Typography,
  Box,
  IconButton,
  InputAdornment,
  Alert,
  Snackbar,
} from '@mui/material'
import { Refresh, ContentCopy } from '@mui/icons-material'
import { rolesApi, Role } from '../../shared/api/roles'

interface EditUserModalProps {
  user: {
    id: string
    email: string
    first_name?: string
    last_name?: string
    is_active: boolean
  }
  roles: Role[]
  onClose: () => void
  onSubmit: (id: string, data: any) => void
}

// Функция генерации случайного пароля
const generatePassword = (): string => {
  const uppercase = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  const lowercase = 'abcdefghijklmnopqrstuvwxyz'
  const numbers = '0123456789'
  const special = '!@#$%^&*'
  const allChars = uppercase + lowercase + numbers + special

  let password = ''
  // Гарантируем хотя бы по одному символу каждого типа
  password += uppercase[Math.floor(Math.random() * uppercase.length)]
  password += lowercase[Math.floor(Math.random() * lowercase.length)]
  password += numbers[Math.floor(Math.random() * numbers.length)]
  password += special[Math.floor(Math.random() * special.length)]

  // Заполняем остаток случайными символами (до 12 символов)
  for (let i = password.length; i < 12; i++) {
    password += allChars[Math.floor(Math.random() * allChars.length)]
  }

  // Перемешиваем символы
  return password.split('').sort(() => Math.random() - 0.5).join('')
}

export const EditUserModal: React.FC<EditUserModalProps> = ({ user, roles, onClose, onSubmit }) => {
  const [formData, setFormData] = useState({
    first_name: user.first_name || '',
    last_name: user.last_name || '',
    password: '',
    is_active: user.is_active,
    role_ids: [] as string[]
  })
  const [formErrors, setFormErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [copySnackbarOpen, setCopySnackbarOpen] = useState(false)

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

  const handleGeneratePassword = () => {
    const newPassword = generatePassword()
    setFormData(prev => ({ ...prev, password: newPassword }))
    setShowPassword(true)
  }

  const handleCopyPassword = () => {
    if (formData.password) {
      navigator.clipboard.writeText(formData.password)
      setCopySnackbarOpen(true)
    }
  }

  const validateForm = (): boolean => {
    const errors: Record<string, string> = {}

    if (!formData.first_name.trim()) {
      errors.first_name = 'Имя обязательно'
    }

    if (!formData.last_name.trim()) {
      errors.last_name = 'Фамилия обязательна'
    }

    // Проверка пароля только если он введен
    if (formData.password && formData.password.length < 6) {
      errors.password = 'Пароль должен содержать минимум 6 символов'
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
      
      const submitData: any = {
        first_name: formData.first_name,
        last_name: formData.last_name,
        is_active: formData.is_active,
        role_ids: filteredRoleIds
      }

      // Добавляем пароль только если он был изменен
      if (formData.password && formData.password.trim() !== '') {
        submitData.password = formData.password
      }

      await onSubmit(user.id, submitData)
    } catch (error) {
      console.error('Ошибка сохранения:', error)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <>
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

            <Box sx={{ mb: 2 }}>
              <Typography variant="subtitle2" gutterBottom>
                Новый пароль (опционально)
              </Typography>
              <TextField
                fullWidth
                type={showPassword ? 'text' : 'password'}
                placeholder="Оставьте пустым, чтобы не менять пароль"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                error={!!formErrors.password}
                helperText={formErrors.password || 'Минимум 6 символов'}
                variant="outlined"
                size="small"
                InputProps={{
                  endAdornment: formData.password && (
                    <InputAdornment position="end">
                      <IconButton
                        onClick={handleCopyPassword}
                        edge="end"
                        title="Копировать пароль"
                      >
                        <ContentCopy fontSize="small" />
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />
              <Button
                startIcon={<Refresh />}
                onClick={handleGeneratePassword}
                size="small"
                sx={{ mt: 1 }}
              >
                Сгенерировать пароль
              </Button>
              {formData.password && (
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={showPassword}
                      onChange={(e) => setShowPassword(e.target.checked)}
                    />
                  }
                  label="Показать пароль"
                  sx={{ ml: 2 }}
                />
              )}
            </Box>
            
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

            {formData.password && (
              <Alert severity="warning" sx={{ mt: 2 }}>
                Обязательно сохраните новый пароль! Он будет показан только один раз.
              </Alert>
            )}
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

      <Snackbar
        open={copySnackbarOpen}
        autoHideDuration={2000}
        onClose={() => setCopySnackbarOpen(false)}
        message="Пароль скопирован в буфер обмена"
      />
    </>
  )
}

