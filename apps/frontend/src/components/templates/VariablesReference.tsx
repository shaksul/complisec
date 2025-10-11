import React, { useEffect, useState } from 'react'
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
  Box,
  Typography,
  LinearProgress,
  Card,
  CardContent,
  Chip,
  IconButton,
  Grid,
} from '@mui/material'
import { ContentCopy, Search } from '@mui/icons-material'
import { templatesApi, TemplateVariable } from '../../shared/api/templates'

interface VariablesReferenceProps {
  onClose: () => void
}

export const VariablesReference: React.FC<VariablesReferenceProps> = ({ onClose }) => {
  const [variables, setVariables] = useState<TemplateVariable[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')

  useEffect(() => {
    loadVariables()
  }, [])

  const loadVariables = async () => {
    try {
      setLoading(true)
      const response = await templatesApi.getTemplateVariables()
      setVariables(response.variables)
    } catch (err) {
      console.error('Failed to load variables:', err)
    } finally {
      setLoading(false)
    }
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    alert(`Скопировано: ${text}`)
  }

  const filteredVariables = variables.filter((v) => {
    const matchesSearch = searchTerm
      ? v.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        v.description.toLowerCase().includes(searchTerm.toLowerCase())
      : true
    const matchesCategory = categoryFilter ? v.category === categoryFilter : true
    return matchesSearch && matchesCategory
  })

  const groupedVariables = filteredVariables.reduce((acc, variable) => {
    if (!acc[variable.category]) {
      acc[variable.category] = []
    }
    acc[variable.category].push(variable)
    return acc
  }, {} as Record<string, TemplateVariable[]>)

  const categoryLabels: Record<string, string> = {
    asset: 'Активы',
    passport: 'Паспортные данные',
    user: 'Пользователи',
    date: 'Даты',
  }

  const categories = Object.keys(groupedVariables).sort()

  return (
    <Dialog open onClose={onClose} maxWidth="lg" fullWidth>
      <DialogTitle>
        <Box>
          <Typography variant="h6">Справка по переменным шаблонов</Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
            Используйте эти переменные в шаблонах для автоматической подстановки данных
          </Typography>
        </Box>
      </DialogTitle>

      <DialogContent dividers>
        {loading ? (
          <Box textAlign="center" py={6}>
            <LinearProgress />
            <Typography sx={{ mt: 2 }} color="text.secondary">
              Загрузка переменных...
            </Typography>
          </Box>
        ) : (
          <>
            {/* Filters */}
            <Grid container spacing={2} sx={{ mb: 3 }}>
              <Grid item xs={12} md={8}>
                <TextField
                  fullWidth
                  size="small"
                  placeholder="Поиск переменных..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  InputProps={{
                    startAdornment: <Search sx={{ mr: 1, color: 'text.secondary' }} />,
                  }}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <FormControl fullWidth size="small">
                  <InputLabel>Категория</InputLabel>
                  <Select
                    value={categoryFilter}
                    onChange={(e) => setCategoryFilter(e.target.value)}
                    label="Категория"
                  >
                    <MenuItem value="">Все категории</MenuItem>
                    {categories.map((cat) => (
                      <MenuItem key={cat} value={cat}>
                        {categoryLabels[cat] || cat}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
            </Grid>

            {/* Variables by category */}
            {categories.map((category) => (
              <Box key={category} sx={{ mb: 3 }}>
                <Typography variant="h6" gutterBottom>
                  {categoryLabels[category] || category}
                </Typography>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
                  {groupedVariables[category].map((variable) => (
                    <Card key={variable.name} variant="outlined" sx={{ '&:hover': { borderColor: 'primary.main' } }}>
                      <CardContent>
                        <Box display="flex" alignItems="flex-start" justifyContent="space-between">
                          <Box flex={1}>
                            <Box display="flex" alignItems="center" gap={1} mb={1}>
                              <Chip
                                label={variable.placeholder}
                                color="primary"
                                variant="outlined"
                                size="small"
                                sx={{ fontFamily: 'monospace' }}
                              />
                              <IconButton
                                size="small"
                                onClick={() => copyToClipboard(variable.placeholder)}
                                title="Копировать в буфер обмена"
                              >
                                <ContentCopy fontSize="small" />
                              </IconButton>
                            </Box>
                            <Typography variant="body2" color="text.primary" gutterBottom>
                              {variable.description}
                            </Typography>
                            {variable.example && (
                              <Typography variant="caption" color="text.secondary">
                                <strong>Пример:</strong> {variable.example}
                              </Typography>
                            )}
                          </Box>
                        </Box>
                      </CardContent>
                    </Card>
                  ))}
                </Box>
              </Box>
            ))}

            {filteredVariables.length === 0 && (
              <Box textAlign="center" py={6}>
                <Typography color="text.secondary">Переменные не найдены</Typography>
              </Box>
            )}

            {/* Usage examples */}
            <Card sx={{ mt: 4, bgcolor: 'info.lighter' }}>
              <CardContent>
                <Typography variant="subtitle2" color="info.dark" gutterBottom>
                  Примеры использования:
                </Typography>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                  <Chip
                    label={`<p>Инвентарный номер: {{inventory_number}}</p>`}
                    sx={{ fontFamily: 'monospace', justifyContent: 'flex-start' }}
                  />
                  <Chip
                    label={`<td>{{model}}</td>`}
                    sx={{ fontFamily: 'monospace', justifyContent: 'flex-start' }}
                  />
                  <Chip
                    label={`<span>Дата: {{current_date}}</span>`}
                    sx={{ fontFamily: 'monospace', justifyContent: 'flex-start' }}
                  />
                </Box>
              </CardContent>
            </Card>
          </>
        )}
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} variant="contained">
          Закрыть
        </Button>
      </DialogActions>
    </Dialog>
  )
}
