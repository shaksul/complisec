import React, { useEffect, useState } from 'react'
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
  Chip,
  LinearProgress,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
  IconButton,
  Tooltip,
  Alert,
} from '@mui/material'
import {
  Add,
  Edit,
  Delete,
  ContentCopy,
  Search,
  HelpOutline,
  Refresh,
} from '@mui/icons-material'
import { templatesApi, DocumentTemplate, TEMPLATE_TYPE_LABELS } from '../../shared/api/templates'
import { TemplateEditorModal } from '../../components/templates/TemplateEditorModal'
import { VariablesReference } from '../../components/templates/VariablesReference'

export const TemplatesPage: React.FC = () => {
  const [templates, setTemplates] = useState<DocumentTemplate[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [searchTerm, setSearchTerm] = useState('')
  const [typeFilter, setTypeFilter] = useState('')
  const [systemFilter, setSystemFilter] = useState<string>('')
  const [editorModalOpen, setEditorModalOpen] = useState(false)
  const [variablesModalOpen, setVariablesModalOpen] = useState(false)
  const [selectedTemplate, setSelectedTemplate] = useState<DocumentTemplate | null>(null)
  const [initializingDefaults, setInitializingDefaults] = useState(false)

  const loadTemplates = async () => {
    try {
      setLoading(true)
      const filters: any = {}
      if (searchTerm) filters.search = searchTerm
      if (typeFilter) filters.template_type = typeFilter
      if (systemFilter === 'system') filters.is_system = true
      if (systemFilter === 'user') filters.is_system = false

      const data = await templatesApi.listTemplates(filters)
      setTemplates(data)
      setError('')
    } catch (err: any) {
      setError(err.message || 'Не удалось загрузить шаблоны')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadTemplates()
  }, [searchTerm, typeFilter, systemFilter])

  const handleInitializeDefaults = async () => {
    if (!confirm('Инициализировать системные шаблоны? Это создаст базовые шаблоны паспортов.')) {
      return
    }

    try {
      setInitializingDefaults(true)
      await templatesApi.initializeDefaultTemplates()
      await loadTemplates()
      alert('Системные шаблоны успешно инициализированы')
    } catch (err: any) {
      setError(err.message || 'Не удалось инициализировать шаблоны')
    } finally {
      setInitializingDefaults(false)
    }
  }

  const handleCreateTemplate = () => {
    setSelectedTemplate(null)
    setEditorModalOpen(true)
  }

  const handleEditTemplate = (template: DocumentTemplate) => {
    if (template.is_system) {
      alert('Системные шаблоны нельзя редактировать. Вы можете скопировать шаблон и создать на его основе новый.')
      return
    }
    setSelectedTemplate(template)
    setEditorModalOpen(true)
  }

  const handleCopyTemplate = (template: DocumentTemplate) => {
    setSelectedTemplate({
      ...template,
      id: '',
      name: template.name + ' (копия)',
      is_system: false,
    } as DocumentTemplate)
    setEditorModalOpen(true)
  }

  const handleDeleteTemplate = async (template: DocumentTemplate) => {
    if (template.is_system) {
      alert('Системные шаблоны нельзя удалить')
      return
    }

    if (!confirm(`Удалить шаблон "${template.name}"?`)) {
      return
    }

    try {
      await templatesApi.deleteTemplate(template.id)
      await loadTemplates()
    } catch (err: any) {
      setError(err.message || 'Не удалось удалить шаблон')
    }
  }

  const handleSaveTemplate = async () => {
    await loadTemplates()
    setEditorModalOpen(false)
  }

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Typography variant="h4">Шаблоны документов</Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
            Управление шаблонами паспортов и других документов
          </Typography>
        </Box>
        <Box display="flex" gap={1}>
          <Button
            variant="outlined"
            startIcon={<HelpOutline />}
            onClick={() => setVariablesModalOpen(true)}
          >
            Справка
          </Button>
          <Button
            variant="outlined"
            startIcon={<Refresh />}
            onClick={handleInitializeDefaults}
            disabled={initializingDefaults}
          >
            {initializingDefaults ? 'Инициализация...' : 'Системные шаблоны'}
          </Button>
          <Button
            variant="contained"
            startIcon={<Add />}
            onClick={handleCreateTemplate}
          >
            Создать шаблон
          </Button>
        </Box>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
          {error}
        </Alert>
      )}

      {/* Фильтры */}
      <Paper sx={{ mb: 2, p: 2 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={4}>
            <TextField
              fullWidth
              size="small"
              placeholder="Поиск по названию..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              InputProps={{
                startAdornment: <Search sx={{ mr: 1, color: 'text.secondary' }} />,
              }}
            />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <FormControl fullWidth size="small">
              <InputLabel>Тип шаблона</InputLabel>
              <Select
                value={typeFilter}
                onChange={(e) => setTypeFilter(e.target.value)}
                label="Тип шаблона"
              >
                <MenuItem value="">Все типы</MenuItem>
                {Object.entries(TEMPLATE_TYPE_LABELS).map(([value, label]) => (
                  <MenuItem key={value} value={value}>
                    {label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <FormControl fullWidth size="small">
              <InputLabel>Категория</InputLabel>
              <Select
                value={systemFilter}
                onChange={(e) => setSystemFilter(e.target.value)}
                label="Категория"
              >
                <MenuItem value="">Все шаблоны</MenuItem>
                <MenuItem value="system">Системные</MenuItem>
                <MenuItem value="user">Пользовательские</MenuItem>
              </Select>
            </FormControl>
          </Grid>
        </Grid>
      </Paper>

      {/* Таблица шаблонов */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Название</TableCell>
                <TableCell>Тип</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Обновлено</TableCell>
                <TableCell align="right">Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={5} align="center">
                    <LinearProgress />
                    <Typography sx={{ mt: 1 }}>Загрузка шаблонов...</Typography>
                  </TableCell>
                </TableRow>
              ) : templates.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} align="center">
                    <Box py={6}>
                      <Typography color="text.secondary">
                        Шаблоны не найдены
                      </Typography>
                      <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                        Нажмите "Создать шаблон" или "Системные шаблоны" для начала работы
                      </Typography>
                    </Box>
                  </TableCell>
                </TableRow>
              ) : (
                templates.map((template) => (
                  <TableRow key={template.id} hover>
                    <TableCell>
                      <Box>
                        <Typography variant="body2" fontWeight="medium">
                          {template.name}
                        </Typography>
                        {template.description && (
                          <Typography variant="caption" color="text.secondary">
                            {template.description}
                          </Typography>
                        )}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {TEMPLATE_TYPE_LABELS[template.template_type] || template.template_type}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Box display="flex" gap={0.5}>
                        {template.is_system && (
                          <Chip
                            label="Системный"
                            color="primary"
                            size="small"
                            variant="outlined"
                          />
                        )}
                        <Chip
                          label={template.is_active ? 'Активный' : 'Неактивный'}
                          color={template.is_active ? 'success' : 'default'}
                          size="small"
                        />
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {new Date(template.updated_at).toLocaleString('ru-RU')}
                      </Typography>
                    </TableCell>
                    <TableCell align="right">
                      <Box display="flex" justifyContent="flex-end" gap={0.5}>
                        <Tooltip title={template.is_system ? 'Просмотр' : 'Редактировать'}>
                          <IconButton
                            size="small"
                            color="primary"
                            onClick={() => handleEditTemplate(template)}
                          >
                            <Edit />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Копировать">
                          <IconButton
                            size="small"
                            color="success"
                            onClick={() => handleCopyTemplate(template)}
                          >
                            <ContentCopy />
                          </IconButton>
                        </Tooltip>
                        {!template.is_system && (
                          <Tooltip title="Удалить">
                            <IconButton
                              size="small"
                              color="error"
                              onClick={() => handleDeleteTemplate(template)}
                            >
                              <Delete />
                            </IconButton>
                          </Tooltip>
                        )}
                      </Box>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {/* Модалки */}
      {editorModalOpen && (
        <TemplateEditorModal
          template={selectedTemplate}
          onClose={() => setEditorModalOpen(false)}
          onSave={handleSaveTemplate}
        />
      )}

      {variablesModalOpen && (
        <VariablesReference onClose={() => setVariablesModalOpen(false)} />
      )}
    </Container>
  )
}
