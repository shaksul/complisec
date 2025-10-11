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
  IconButton,
  Tooltip,
  Alert,
  Card,
  CardContent,
  Grid,
} from '@mui/material'
import {
  Add,
  Edit,
  Refresh as RefreshIcon,
} from '@mui/icons-material'
import { templatesApi, InventoryNumberRule } from '../../shared/api/templates'
import { InventoryRuleModal } from '../../components/inventory/InventoryRuleModal'

export const InventoryRulesPage: React.FC = () => {
  const [rules, setRules] = useState<InventoryNumberRule[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [selectedRule, setSelectedRule] = useState<InventoryNumberRule | null>(null)

  const loadRules = async () => {
    try {
      setLoading(true)
      const data = await templatesApi.listInventoryRules()
      setRules(data)
      setError('')
    } catch (err: any) {
      setError(err.message || 'Не удалось загрузить правила')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadRules()
  }, [])

  const handleCreateRule = () => {
    setSelectedRule(null)
    setModalOpen(true)
  }

  const handleEditRule = (rule: InventoryNumberRule) => {
    setSelectedRule(rule)
    setModalOpen(true)
  }

  const handleResetSequence = async (rule: InventoryNumberRule) => {
    if (!confirm(`Сбросить счетчик для "${rule.asset_type}"? Следующий номер будет 1.`)) {
      return
    }

    try {
      await templatesApi.updateInventoryRule(rule.id, { current_sequence: 0 })
      await loadRules()
      alert('Счетчик успешно сброшен')
    } catch (err: any) {
      setError(err.message || 'Не удалось сбросить счетчик')
    }
  }

  const handleSaveRule = async () => {
    await loadRules()
    setModalOpen(false)
  }

  const renderPatternPreview = (rule: InventoryNumberRule) => {
    const now = new Date()
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const nextSeq = rule.current_sequence + 1

    let preview = rule.pattern
    preview = preview.replace(/\{\{type_code\}\}/g, rule.asset_type)
    preview = preview.replace(/\{\{year\}\}/g, String(year))
    preview = preview.replace(/\{\{month\}\}/g, month)
    preview = preview.replace(/\{\{class\}\}/g, rule.asset_class || 'CLASS')
    
    // Handle sequence with format
    preview = preview.replace(/\{\{sequence:?(\d+)?\}\}/g, (_match, width) => {
      if (width) {
        return String(nextSeq).padStart(parseInt(width), '0')
      }
      return String(nextSeq)
    })

    return preview
  }

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Typography variant="h4">Правила генерации инвентарных номеров</Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
            Настройка автоматической генерации инвентарных номеров для активов
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={handleCreateRule}
        >
          Добавить правило
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
          {error}
        </Alert>
      )}

      {/* Справка по переменным паттерна */}
      <Card sx={{ mb: 3, bgcolor: 'info.lighter' }}>
        <CardContent>
          <Typography variant="h6" gutterBottom color="info.dark">
            Переменные для паттерна:
          </Typography>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <Box display="flex" alignItems="center" gap={1}>
                <Chip label="{{type_code}}" size="small" />
                <Typography variant="body2">- код типа актива</Typography>
              </Box>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Box display="flex" alignItems="center" gap={1}>
                <Chip label="{{year}}" size="small" />
                <Typography variant="body2">- год (4 цифры)</Typography>
              </Box>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Box display="flex" alignItems="center" gap={1}>
                <Chip label="{{month}}" size="small" />
                <Typography variant="body2">- месяц (2 цифры)</Typography>
              </Box>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Box display="flex" alignItems="center" gap={1}>
                <Chip label="{{sequence:0000}}" size="small" />
                <Typography variant="body2">- последовательность с форматом</Typography>
              </Box>
            </Grid>
          </Grid>
          <Box mt={2} p={1.5} bgcolor="background.paper" borderRadius={1}>
            <Typography variant="caption" color="text.secondary">
              Пример: <Chip label="РСП-{{year}}-{{sequence:0000}}" size="small" sx={{ mx: 0.5 }} /> → РСП-2025-0001
            </Typography>
          </Box>
        </CardContent>
      </Card>

      {/* Таблица правил */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Тип актива</TableCell>
                <TableCell>Паттерн</TableCell>
                <TableCell>Текущий счетчик</TableCell>
                <TableCell>Следующий номер</TableCell>
                <TableCell align="right">Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={5} align="center">
                    <LinearProgress />
                    <Typography sx={{ mt: 1 }}>Загрузка правил...</Typography>
                  </TableCell>
                </TableRow>
              ) : rules.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} align="center">
                    <Box py={6}>
                      <Typography color="text.secondary">
                        Правила не найдены
                      </Typography>
                      <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                        Нажмите "Добавить правило" для создания первого правила
                      </Typography>
                    </Box>
                  </TableCell>
                </TableRow>
              ) : (
                rules.map((rule) => (
                  <TableRow key={rule.id} hover>
                    <TableCell>
                      <Box>
                        <Typography variant="body2" fontWeight="medium">
                          {rule.asset_type}
                          {rule.asset_class && (
                            <Chip
                              label={rule.asset_class}
                              size="small"
                              sx={{ ml: 1 }}
                              variant="outlined"
                            />
                          )}
                        </Typography>
                        {rule.prefix && (
                          <Typography variant="caption" color="text.secondary" display="block">
                            Префикс: {rule.prefix}
                          </Typography>
                        )}
                        {rule.description && (
                          <Typography variant="caption" color="text.secondary" display="block">
                            {rule.description}
                          </Typography>
                        )}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={rule.pattern}
                        size="small"
                        sx={{ fontFamily: 'monospace' }}
                      />
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" fontWeight="medium">
                        {rule.current_sequence}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={renderPatternPreview(rule)}
                        color="success"
                        size="small"
                        sx={{ fontFamily: 'monospace' }}
                      />
                    </TableCell>
                    <TableCell align="right">
                      <Box display="flex" justifyContent="flex-end" gap={0.5}>
                        <Tooltip title="Редактировать">
                          <IconButton
                            size="small"
                            color="primary"
                            onClick={() => handleEditRule(rule)}
                          >
                            <Edit />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Сбросить счетчик">
                          <IconButton
                            size="small"
                            color="warning"
                            onClick={() => handleResetSequence(rule)}
                          >
                            <RefreshIcon />
                          </IconButton>
                        </Tooltip>
                      </Box>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {/* Модалка */}
      {modalOpen && (
        <InventoryRuleModal
          rule={selectedRule}
          onClose={() => setModalOpen(false)}
          onSave={handleSaveRule}
        />
      )}
    </Container>
  )
}
