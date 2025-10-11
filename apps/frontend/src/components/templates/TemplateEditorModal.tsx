import React, { useState, useEffect } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Alert,
  Box,
  Tabs,
  Tab,
  Typography,
  Checkbox,
  FormControlLabel,
  Chip,
} from '@mui/material'
import Editor from '@monaco-editor/react'
import { templatesApi, DocumentTemplate, CreateTemplateRequest, TEMPLATE_TYPE_OPTIONS } from '../../shared/api/templates'

interface TemplateEditorModalProps {
  template: DocumentTemplate | null
  onClose: () => void
  onSave: () => void
}

export const TemplateEditorModal: React.FC<TemplateEditorModalProps> = ({
  template,
  onClose,
  onSave,
}) => {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [templateType, setTemplateType] = useState('passport_device')
  const [content, setContent] = useState('')
  const [isActive, setIsActive] = useState(true)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [currentTab, setCurrentTab] = useState(0)

  const isEditMode = template && template.id
  const isReadOnly = template?.is_system || false

  useEffect(() => {
    if (template) {
      setName(template.name)
      setDescription(template.description || '')
      setTemplateType(template.template_type)
      setContent(template.content)
      setIsActive(template.is_active)
    }
  }, [template])

  const handleSave = async () => {
    if (!name.trim()) {
      setError('Введите название шаблона')
      return
    }

    if (!content.trim()) {
      setError('Введите содержимое шаблона')
      return
    }

    try {
      setLoading(true)
      setError('')

      if (isEditMode && !isReadOnly) {
        await templatesApi.updateTemplate(template.id, {
          name,
          description: description || undefined,
          template_type: templateType,
          content,
          is_active: isActive,
        })
      } else {
        const request: CreateTemplateRequest = {
          name,
          description: description || undefined,
          template_type: templateType,
          content,
        }
        await templatesApi.createTemplate(request)
      }

      onSave()
    } catch (err: any) {
      setError(err.message || 'Не удалось сохранить шаблон')
    } finally {
      setLoading(false)
    }
  }

  const handleEditorChange = (value: string | undefined) => {
    setContent(value || '')
  }

  const renderPreview = () => {
    const testData: Record<string, string> = {
      asset_name: 'Тестовый актив',
      inventory_number: 'РСП-2025-0001',
      serial_number: 'ABC123456',
      model: 'Dell OptiPlex 7090',
      manufacturer: 'Dell Inc.',
      cpu: 'Intel Core i7-11700',
      ram: '16 GB DDR4',
      hdd_info: 'SSD 512GB',
      network_card: 'Intel I219-V',
      ip_address: '192.168.1.100',
      mac_address: '00:1A:2B:3C:4D:5E',
      pc_number: 'PC-101',
      location: 'Кабинет 101',
      owner_name: 'Иванов И.И.',
      responsible_user_name: 'Петров П.П.',
      current_date: new Date().toLocaleDateString('ru-RU'),
      purchase_year: '2023',
      warranty_until: '31.12.2026',
    }

    let previewHtml = content
    Object.entries(testData).forEach(([key, value]) => {
      const regex = new RegExp(`\\{\\{${key}\\}\\}`, 'g')
      previewHtml = previewHtml.replace(regex, value)
    })

    return (
      <Box
        sx={{
          border: 1,
          borderColor: 'divider',
          borderRadius: 1,
          p: 3,
          bgcolor: 'background.paper',
          height: '100%',
          overflow: 'auto',
        }}
        dangerouslySetInnerHTML={{ __html: previewHtml }}
      />
    )
  }

  return (
    <Dialog open fullScreen onClose={onClose}>
      <DialogTitle>
        <Box>
          <Typography variant="h6">
            {isReadOnly
              ? 'Просмотр шаблона'
              : isEditMode
              ? 'Редактирование шаблона'
              : 'Создание шаблона'}
          </Typography>
          {isReadOnly && (
            <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
              Системный шаблон нельзя редактировать. Используйте "Копировать" для создания нового.
            </Typography>
          )}
        </Box>
      </DialogTitle>

      <DialogContent dividers>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
            {error}
          </Alert>
        )}

        <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 2, mb: 2 }}>
          <TextField
            label="Название"
            required
            fullWidth
            value={name}
            onChange={(e) => setName(e.target.value)}
            disabled={isReadOnly}
            placeholder="Например: Паспорт ПК"
          />

          <FormControl fullWidth required>
            <InputLabel>Тип шаблона</InputLabel>
            <Select
              value={templateType}
              onChange={(e) => setTemplateType(e.target.value)}
              disabled={isReadOnly}
              label="Тип шаблона"
            >
              {TEMPLATE_TYPE_OPTIONS.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <Box sx={{ gridColumn: '1 / -1' }}>
            <TextField
              label="Описание"
              fullWidth
              multiline
              rows={2}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={isReadOnly}
              placeholder="Краткое описание шаблона"
            />
          </Box>

          {isEditMode && !isReadOnly && (
            <Box sx={{ gridColumn: '1 / -1' }}>
              <FormControlLabel
                control={
                  <Checkbox
                    checked={isActive}
                    onChange={(e) => setIsActive(e.target.checked)}
                  />
                }
                label="Активный шаблон"
              />
            </Box>
          )}
        </Box>

        <Tabs value={currentTab} onChange={(_, v) => setCurrentTab(v)} sx={{ mb: 2 }}>
          <Tab label="Редактор" />
          <Tab label="Предпросмотр" />
        </Tabs>

        <Box sx={{ height: 500 }}>
          {currentTab === 0 ? (
            <Box sx={{ border: 1, borderColor: 'divider', borderRadius: 1, overflow: 'hidden', height: '100%' }}>
              <Editor
                height="100%"
                defaultLanguage="html"
                value={content}
                onChange={handleEditorChange}
                options={{
                  readOnly: isReadOnly,
                  minimap: { enabled: false },
                  fontSize: 13,
                  lineNumbers: 'on',
                  scrollBeyondLastLine: false,
                  wordWrap: 'on',
                }}
                theme="vs-light"
              />
            </Box>
          ) : (
            renderPreview()
          )}
        </Box>

        <Box sx={{ mt: 2 }}>
          <Typography variant="caption" color="text.secondary">
            Используйте переменные в формате <Chip label="{{имя_переменной}}" size="small" sx={{ mx: 0.5 }} />
          </Typography>
          <Box sx={{ mt: 0.5 }}>
            <Typography variant="caption" color="text.secondary">
              Примеры:
              {[' {{asset_name}}', ' {{serial_number}}', ' {{current_date}}'].map((example) => (
                <Chip key={example} label={example} size="small" sx={{ mx: 0.25 }} />
              ))}
            </Typography>
          </Box>
        </Box>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} color="inherit">
          {isReadOnly ? 'Закрыть' : 'Отмена'}
        </Button>
        {!isReadOnly && (
          <Button onClick={handleSave} variant="contained" disabled={loading}>
            {loading ? 'Сохранение...' : 'Сохранить'}
          </Button>
        )}
      </DialogActions>
    </Dialog>
  )
}
