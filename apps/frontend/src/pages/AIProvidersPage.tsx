import { useEffect, useState } from "react"
import {
  Box,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  IconButton,
  Alert,
  Switch,
  FormControlLabel,
  Tooltip,
} from "@mui/material"
import { Add, Edit, Delete, Visibility, Psychology } from "@mui/icons-material"
import { PageContainer, PageHeader, SectionCard } from '@/components/common/Page'
import type { AIProvider, CreateAIProviderDTO, UpdateAIProviderDTO } from "@/shared/api/ai"
import { getProviders, addProvider, updateProvider, deleteProvider, getProvider } from "@/shared/api/ai"

export default function AIProvidersPage() {
  const [items, setItems] = useState<AIProvider[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [openDialog, setOpenDialog] = useState(false)
  const [editingProvider, setEditingProvider] = useState<AIProvider | null>(null)
  const [viewingProvider, setViewingProvider] = useState<AIProvider | null>(null)

  // Form fields
  const [name, setName] = useState("")
  const [url, setUrl] = useState("")
  const [apiKey, setApiKey] = useState("")
  const [roles, setRoles] = useState("")
  const [models, setModels] = useState("")
  const [defaultModel, setDefaultModel] = useState("")
  const [isActive, setIsActive] = useState(true)

  useEffect(() => {
    loadProviders()
  }, [])

  async function loadProviders() {
    try {
      setLoading(true)
      setError(null)
      const data = await getProviders()
      setItems(Array.isArray(data) ? data : [])
    } catch (err: any) {
      console.error('Error loading AI providers:', err)
      setError(err?.response?.data?.error || 'Ошибка загрузки провайдеров')
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  function handleOpenCreateDialog() {
    setEditingProvider(null)
    setName("")
    setUrl("")
    setApiKey("")
    setRoles("")
    setModels("llama3.2, llama3.1, mistral, gemma2")
    setDefaultModel("llama3.2")
    setIsActive(true)
    setOpenDialog(true)
  }

  async function handleOpenEditDialog(id: string) {
    try {
      const provider = await getProvider(id)
      setEditingProvider(provider)
      setName(provider.name)
      setUrl(provider.base_url)
      setApiKey(provider.api_key || "")
      setRoles(provider.roles.join(", "))
      setModels(provider.models.join(", "))
      setDefaultModel(provider.default_model)
      setIsActive(provider.is_active)
      setOpenDialog(true)
    } catch (err: any) {
      setError(err?.response?.data?.error || 'Ошибка загрузки провайдера')
    }
  }

  async function handleOpenViewDialog(id: string) {
    try {
      const provider = await getProvider(id)
      setViewingProvider(provider)
    } catch (err: any) {
      setError(err?.response?.data?.error || 'Ошибка загрузки провайдера')
    }
  }

  async function handleSave() {
    try {
      const rolesArray = roles.split(",").map(r => r.trim()).filter(r => r !== "")
      const modelsArray = models.split(",").map(m => m.trim()).filter(m => m !== "")

      if (editingProvider) {
        // Обновление
        const updateData: UpdateAIProviderDTO = {
          name,
          base_url: url,
          api_key: apiKey,
          roles: rolesArray,
          models: modelsArray,
          default_model: defaultModel,
          is_active: isActive,
        }
        await updateProvider(editingProvider.id, updateData)
      } else {
        // Создание
        const createData: CreateAIProviderDTO = {
          name,
          base_url: url,
          api_key: apiKey,
          roles: rolesArray,
          models: modelsArray,
          default_model: defaultModel,
        }
        await addProvider(createData)
      }

      await loadProviders()
      setOpenDialog(false)
    } catch (err: any) {
      console.error('Error saving provider:', err)
      setError(err?.response?.data?.error || 'Ошибка сохранения провайдера')
    }
  }

  async function handleDelete(id: string) {
    if (!confirm('Вы уверены, что хотите удалить этого провайдера?')) {
      return
    }

    try {
      await deleteProvider(id)
      await loadProviders()
    } catch (err: any) {
      setError(err?.response?.data?.error || 'Ошибка удаления провайдера')
    }
  }

  return (
    <PageContainer>
      <PageHeader
        title="AI Провайдеры"
        subtitle="Управление подключениями к AI сервисам и моделям"
        actions={
          <Button
            variant="contained"
            startIcon={<Add />}
            onClick={handleOpenCreateDialog}
          >
            Добавить провайдера
          </Button>
        }
      />

      {error && (
        <Alert severity="error" onClose={() => setError(null)} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <SectionCard
        title="Список провайдеров"
        description="Настроенные AI провайдеры и модели"
      >
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Название</TableCell>
                <TableCell>URL</TableCell>
                <TableCell>Модели</TableCell>
                <TableCell>Роли</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell align="right">Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={6} align="center">
                    Загрузка...
                  </TableCell>
                </TableRow>
              ) : items.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} align="center">
                    Нет провайдеров
                  </TableCell>
                </TableRow>
              ) : (
                items.map((provider) => (
                  <TableRow key={provider.id}>
                    <TableCell>
                      <Box display="flex" alignItems="center" gap={1}>
                        <Psychology />
                        {provider.name}
                      </Box>
                    </TableCell>
                    <TableCell>{provider.base_url}</TableCell>
                    <TableCell>
                      <Box display="flex" gap={0.5} flexWrap="wrap">
                        {provider.models?.map((model: string) => (
                          <Chip
                            key={model}
                            label={model}
                            size="small"
                            color={model === provider.default_model ? "primary" : "default"}
                          />
                        ))}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Box display="flex" gap={0.5} flexWrap="wrap">
                        {provider.roles?.map((role: string) => (
                          <Chip key={role} label={role} size="small" variant="outlined" />
                        ))}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={provider.is_active ? "Активен" : "Неактивен"}
                        color={provider.is_active ? "success" : "default"}
                        size="small"
                      />
                    </TableCell>
                    <TableCell align="right">
                      <Tooltip title="Просмотр">
                        <IconButton
                          size="small"
                          onClick={() => handleOpenViewDialog(provider.id)}
                        >
                          <Visibility />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Редактировать">
                        <IconButton
                          size="small"
                          onClick={() => handleOpenEditDialog(provider.id)}
                        >
                          <Edit />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Удалить">
                        <IconButton
                          size="small"
                          color="error"
                          onClick={() => handleDelete(provider.id)}
                        >
                          <Delete />
                        </IconButton>
                      </Tooltip>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </SectionCard>

      {/* Диалог создания/редактирования */}
      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          {editingProvider ? "Редактировать провайдера" : "Добавить AI провайдера"}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
            <TextField
              label="Название"
              fullWidth
              variant="outlined"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
            <TextField
              label="Base URL"
              fullWidth
              variant="outlined"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="https://api.openai.com или http://localhost:11434"
              required
            />
            <TextField
              label="API Key"
              fullWidth
              variant="outlined"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              placeholder="Оставьте пустым, если не требуется"
            />
            <TextField
              label="Роли (через запятую)"
              fullWidth
              variant="outlined"
              value={roles}
              onChange={(e) => setRoles(e.target.value)}
              placeholder="docs, risks, incidents, chat"
              required
            />
            <TextField
              label="Модели (через запятую)"
              fullWidth
              variant="outlined"
              value={models}
              onChange={(e) => setModels(e.target.value)}
              placeholder="llama3.2, llama3.1, mistral, gemma2"
              helperText="Список доступных моделей для этого провайдера"
            />
            <TextField
              label="Модель по умолчанию"
              fullWidth
              variant="outlined"
              value={defaultModel}
              onChange={(e) => setDefaultModel(e.target.value)}
              placeholder="llama3.2"
              helperText="Модель, которая будет использоваться по умолчанию"
            />
            {editingProvider && (
              <FormControlLabel
                control={
                  <Switch
                    checked={isActive}
                    onChange={(e) => setIsActive(e.target.checked)}
                  />
                }
                label="Активен"
              />
            )}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Отмена</Button>
          <Button onClick={handleSave} variant="contained">
            {editingProvider ? "Сохранить" : "Добавить"}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Диалог просмотра */}
      <Dialog open={!!viewingProvider} onClose={() => setViewingProvider(null)} maxWidth="sm" fullWidth>
        <DialogTitle>Информация о провайдере</DialogTitle>
        <DialogContent>
          {viewingProvider && (
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
              <Box>
                <strong>Название:</strong> {viewingProvider.name}
              </Box>
              <Box>
                <strong>URL:</strong> {viewingProvider.base_url}
              </Box>
              <Box>
                <strong>Модели:</strong>
                <Box sx={{ mt: 1, display: 'flex', gap: 0.5, flexWrap: 'wrap' }}>
                  {viewingProvider.models?.map((model) => (
                    <Chip
                      key={model}
                      label={model}
                      size="small"
                      color={model === viewingProvider.default_model ? "primary" : "default"}
                    />
                  ))}
                </Box>
              </Box>
              <Box>
                <strong>Модель по умолчанию:</strong> {viewingProvider.default_model}
              </Box>
              <Box>
                <strong>Роли:</strong>
                <Box sx={{ mt: 1, display: 'flex', gap: 0.5, flexWrap: 'wrap' }}>
                  {viewingProvider.roles?.map((role) => (
                    <Chip key={role} label={role} size="small" variant="outlined" />
                  ))}
                </Box>
              </Box>
              <Box>
                <strong>Статус:</strong>{' '}
                <Chip
                  label={viewingProvider.is_active ? "Активен" : "Неактивен"}
                  color={viewingProvider.is_active ? "success" : "default"}
                  size="small"
                />
              </Box>
              {viewingProvider.api_key && (
                <Box>
                  <strong>API Key:</strong> ••••••••
                </Box>
              )}
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setViewingProvider(null)}>Закрыть</Button>
        </DialogActions>
      </Dialog>
    </PageContainer>
  )
}
