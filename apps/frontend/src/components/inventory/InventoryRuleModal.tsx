import React, { useState, useEffect } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Alert,
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Chip,
} from '@mui/material'
import { templatesApi, InventoryNumberRule, CreateInventoryRuleRequest } from '../../shared/api/templates'

interface InventoryRuleModalProps {
  rule: InventoryNumberRule | null
  onClose: () => void
  onSave: () => void
}

export const InventoryRuleModal: React.FC<InventoryRuleModalProps> = ({
  rule,
  onClose,
  onSave,
}) => {
  const [assetType, setAssetType] = useState('')
  const [assetClass, setAssetClass] = useState('')
  const [pattern, setPattern] = useState('{{type_code}}-{{year}}-{{sequence:0000}}')
  const [prefix, setPrefix] = useState('')
  const [description, setDescription] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const isEditMode = rule && rule.id

  useEffect(() => {
    if (rule) {
      setAssetType(rule.asset_type)
      setAssetClass(rule.asset_class || '')
      setPattern(rule.pattern)
      setPrefix(rule.prefix || '')
      setDescription(rule.description || '')
    }
  }, [rule])

  const renderLivePreview = () => {
    const now = new Date()
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const sequence = (rule?.current_sequence || 0) + 1

    let preview = pattern
    preview = preview.replace(/\{\{type_code\}\}/g, assetType || 'TYPE')
    preview = preview.replace(/\{\{year\}\}/g, String(year))
    preview = preview.replace(/\{\{month\}\}/g, month)
    preview = preview.replace(/\{\{class\}\}/g, assetClass || 'CLASS')
    
    preview = preview.replace(/\{\{sequence:?(\d+)?\}\}/g, (_match, width) => {
      if (width) {
        return String(sequence).padStart(parseInt(width), '0')
      }
      return String(sequence)
    })

    return preview
  }

  const handleSave = async () => {
    if (!assetType.trim()) {
      setError('Введите тип актива')
      return
    }

    if (!pattern.trim()) {
      setError('Введите паттерн')
      return
    }

    try {
      setLoading(true)
      setError('')

      if (isEditMode) {
        await templatesApi.updateInventoryRule(rule.id, {
          pattern,
          prefix: prefix || undefined,
          description: description || undefined,
        })
      } else {
        const request: CreateInventoryRuleRequest = {
          asset_type: assetType,
          asset_class: assetClass || undefined,
          pattern,
          prefix: prefix || undefined,
          description: description || undefined,
        }
        await templatesApi.createInventoryRule(request)
      }

      onSave()
    } catch (err: any) {
      setError(err.message || 'Не удалось сохранить правило')
    } finally {
      setLoading(false)
    }
  }

  const patternTemplates = [
    { label: 'Префикс-Год-Счетчик', value: '{{type_code}}-{{year}}-{{sequence:0000}}' },
    { label: 'Префикс-Год-Месяц-Счетчик', value: '{{type_code}}-{{year}}-{{month}}-{{sequence:000}}' },
    { label: 'Только счетчик', value: '{{sequence:00000}}' },
    { label: 'Кастомный префикс', value: 'CUSTOM-{{year}}-{{sequence:0000}}' },
  ]

  return (
    <Dialog open onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        <Box>
          <Typography variant="h6">
            {isEditMode ? 'Редактирование правила' : 'Создание правила'}
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
            Настройте автоматическую генерацию инвентарных номеров для типа актива
          </Typography>
        </Box>
      </DialogTitle>

      <DialogContent dividers>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
            {error}
          </Alert>
        )}

        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <TextField
            label="Тип актива"
            required
            fullWidth
            value={assetType}
            onChange={(e) => setAssetType(e.target.value)}
            disabled={!!isEditMode}
            placeholder="hardware, software, network..."
            helperText={isEditMode && 'Тип актива нельзя изменить после создания'}
          />

          <TextField
            label="Класс актива (опционально)"
            fullWidth
            value={assetClass}
            onChange={(e) => setAssetClass(e.target.value)}
            disabled={!!isEditMode}
            placeholder="workstation, server, router..."
          />

          {!isEditMode && (
            <Box>
              <Typography variant="body2" fontWeight="medium" gutterBottom>
                Шаблоны паттернов
              </Typography>
              <Grid container spacing={1}>
                {patternTemplates.map((template) => (
                  <Grid item xs={6} key={template.value}>
                    <Button
                      fullWidth
                      variant="outlined"
                      onClick={() => setPattern(template.value)}
                      sx={{ textAlign: 'left', display: 'flex', flexDirection: 'column', alignItems: 'flex-start', py: 1.5 }}
                    >
                      <Typography variant="body2" fontWeight="medium">
                        {template.label}
                      </Typography>
                      <Typography variant="caption" color="text.secondary" sx={{ fontFamily: 'monospace' }}>
                        {template.value}
                      </Typography>
                    </Button>
                  </Grid>
                ))}
              </Grid>
            </Box>
          )}

          <TextField
            label="Паттерн"
            required
            fullWidth
            value={pattern}
            onChange={(e) => setPattern(e.target.value)}
            placeholder="РСП-{{year}}-{{sequence:0000}}"
            helperText="Доступные переменные: {{type_code}}, {{year}}, {{month}}, {{sequence:0000}}, {{class}}"
            InputProps={{ sx: { fontFamily: 'monospace', fontSize: '0.9rem' } }}
          />

          <TextField
            label="Префикс (опционально)"
            fullWidth
            value={prefix}
            onChange={(e) => setPrefix(e.target.value)}
            placeholder="РСП, СО, СР..."
          />

          <TextField
            label="Описание"
            fullWidth
            multiline
            rows={2}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Краткое описание правила"
          />

          <Card sx={{ bgcolor: 'success.lighter' }}>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
                <Typography variant="subtitle2" color="success.dark">
                  Live Preview
                </Typography>
                <Typography variant="caption" color="success.dark">
                  Следующий номер: {(rule?.current_sequence || 0) + 1}
                </Typography>
              </Box>
              <Box
                sx={{
                  bgcolor: 'background.paper',
                  p: 2,
                  borderRadius: 1,
                  border: 1,
                  borderColor: 'success.main',
                }}
              >
                <Chip
                  label={renderLivePreview()}
                  color="success"
                  sx={{ fontFamily: 'monospace', fontSize: '1.1rem', height: 'auto', py: 1 }}
                />
              </Box>
              <Typography variant="caption" color="success.dark" display="block" mt={1}>
                Так будет выглядеть следующий сгенерированный номер
              </Typography>
            </CardContent>
          </Card>
        </Box>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} color="inherit">
          Отмена
        </Button>
        <Button onClick={handleSave} variant="contained" disabled={loading}>
          {loading ? 'Сохранение...' : 'Сохранить'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
