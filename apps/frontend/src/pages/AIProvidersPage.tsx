import { useEffect, useState } from "react"
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
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Chip,
} from "@mui/material"
import { Add, Psychology } from "@mui/icons-material"
import { getProviders, addProvider } from "@/shared/api/ai"

export default function AIProvidersPage() {
  const [items, setItems] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [open, setOpen] = useState(false)
  const [name, setName] = useState("")
  const [url, setUrl] = useState("")
  const [apiKey, setApiKey] = useState("")
  const [roles, setRoles] = useState("")

  useEffect(() => {
    const loadProviders = async () => {
      try {
        setLoading(true)
        setError(null)
        const data = await getProviders()
        setItems(Array.isArray(data) ? data : [])
      } catch (err) {
        console.error('Error loading AI providers:', err)
        setError('Ошибка загрузки провайдеров')
        setItems([])
      } finally {
        setLoading(false)
      }
    }
    loadProviders()
  }, [])

  async function handleAdd() {
    try {
      await addProvider({
        name,
        base_url: url,
        api_key: apiKey,
        roles: roles.split(",").map(r => r.trim()),
      })
      const data = await getProviders()
      setItems(Array.isArray(data) ? data : [])
      setOpen(false)
      setName("")
      setUrl("")
      setApiKey("")
      setRoles("")
    } catch (err) {
      console.error('Error adding provider:', err)
      setError('Ошибка добавления провайдера')
    }
  }

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">AI Провайдеры</Typography>
        <Button variant="contained" startIcon={<Add />} onClick={() => setOpen(true)}>
          Добавить провайдера
        </Button>
      </Box>

      {error && (
        <Paper sx={{ p: 2, mb: 2, bgcolor: 'error.light', color: 'error.contrastText' }}>
          <Typography>{error}</Typography>
        </Paper>
      )}

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Название</TableCell>
                <TableCell>URL</TableCell>
                <TableCell>Роли</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={5} align="center">
                    <Typography>Загрузка...</Typography>
                  </TableCell>
                </TableRow>
              ) : items.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} align="center">
                    <Typography>Нет провайдеров</Typography>
                  </TableCell>
                </TableRow>
              ) : (
                items.map((provider) => (
                  <TableRow key={provider.id}>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <Psychology sx={{ mr: 1 }} />
                        {provider.name}
                      </Box>
                    </TableCell>
                    <TableCell>{provider.base_url}</TableCell>
                    <TableCell>
                      <Box display="flex" gap={0.5} flexWrap="wrap">
                        {provider.roles?.map((role: string) => (
                          <Chip key={role} label={role} size="small" />
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
                    <TableCell>
                      <Button size="small">Редактировать</Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      <Dialog open={open} onClose={() => setOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Добавить AI провайдера</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Название"
            fullWidth
            variant="outlined"
            value={name}
            onChange={(e) => setName(e.target.value)}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Base URL"
            fullWidth
            variant="outlined"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="API Key"
            fullWidth
            variant="outlined"
            value={apiKey}
            onChange={(e) => setApiKey(e.target.value)}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Роли (через запятую)"
            fullWidth
            variant="outlined"
            value={roles}
            onChange={(e) => setRoles(e.target.value)}
            placeholder="docs, risks, incidents"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Отмена</Button>
          <Button onClick={handleAdd} variant="contained">
            Добавить
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  )
}
