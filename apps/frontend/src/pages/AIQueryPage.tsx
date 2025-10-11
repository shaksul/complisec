import { useState } from 'react'
import {
  TextField,
  Button,
  Box,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Chip,
  Stack,
  FormControlLabel,
  Switch,
  Typography,
  Card,
  CardContent,
  Divider,
} from '@mui/material'
import { SendRounded, PsychologyAltRounded } from '@mui/icons-material'
import { queryAI, type RAGSource } from '@/shared/api/ai'
import { PageContainer, PageHeader, SectionCard } from '@/components/common/Page'

type RoleKey = 'docs' | 'risks' | 'incidents' | 'compliance'

const ROLES: Array<{ value: RoleKey; label: string }> = [
  { value: 'docs', label: 'Документы и политики' },
  { value: 'risks', label: 'Риски и оценка' },
  { value: 'incidents', label: 'Инциденты и реагирование' },
  { value: 'compliance', label: 'Комплаенс и аудит' },
]

const ROLE_HINTS: Record<RoleKey, string> = {
  docs: 'docs — поиск и анализ регламентов, политик и процедур',
  risks: 'risks — оценка рисков и рекомендации по обработке',
  incidents: 'incidents — реагирование на инциденты и уроки, извлечённые из них',
  compliance: 'compliance — соответствие стандартам и подготовка к проверкам',
}

export default function AIQueryPage() {
  const [input, setInput] = useState('')
  const [output, setOutput] = useState('')
  const [role, setRole] = useState<RoleKey>('docs')
  const [isLoading, setIsLoading] = useState(false)
  const [useRAG, setUseRAG] = useState(true)
  const [sources, setSources] = useState<RAGSource[]>([])

  async function handleSend() {
    setIsLoading(true)
    setSources([])
    try {
      // Используем единый endpoint /api/ai/query с флагом use_rag
      // provider_id берем из настроек (по умолчанию первый активный провайдер)
      const response = await queryAI({
        provider_id: 'b88bc531-cb8d-4b4d-9298-b3b093da8f03', // nitek провайдер (можно сделать динамическим)
        role,
        input,
        context: {},
        use_rag: useRAG,
      })
      
      if (useRAG && response.sources) {
        // Если RAG включен, response будет содержать answer и sources
        setOutput(response.answer || response.output || '')
        setSources(response.sources || [])
      } else {
        // Обычный ответ без RAG
        setOutput(response.output || response.answer || '')
      }
    } catch (error) {
      const message = (error as any)?.message ?? 'произошла ошибка'
      setOutput(`Ошибка запроса: ${message}`)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <PageContainer>
      <PageHeader
        title="AI-аналитика"
        subtitle="Сформулируйте задачу для корпоративного ассистента: политика, риск, инцидент или комплаенс"
        actions={<PsychologyAltRounded color="primary" fontSize="large" />}
      />

      <SectionCard
        title="Формулировка запроса"
        description="Выберите контекст и опишите задачу, которую нужно решить"
      >
        <Stack spacing={2.5}>
          <FormControlLabel
            control={
              <Switch
                checked={useRAG}
                onChange={(e) => setUseRAG(e.target.checked)}
                color="primary"
              />
            }
            label={
              <Box>
                <Typography variant="body1">
                  Использовать контекст из документов (GraphRAG)
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  При включении ИИ будет использовать информацию из ваших документов
                </Typography>
              </Box>
            }
          />

          {!useRAG && (
            <FormControl fullWidth>
              <InputLabel id="ai-role">Контекст</InputLabel>
              <Select<RoleKey>
                labelId="ai-role"
                value={role}
                label="Контекст"
                onChange={(event) => setRole(event.target.value as RoleKey)}
              >
                {ROLES.map((item) => (
                  <MenuItem key={item.value} value={item.value}>
                    {item.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          )}

          <TextField
            label="Опишите задачу"
            multiline
            minRows={5}
            value={input}
            placeholder={
              useRAG
                ? 'Например: какие требования к паролям указаны в политике безопасности?'
                : 'Например: подготовь краткий обзор требований ISO 27001 для отчёта руководству'
            }
            onChange={(event) => setInput(event.target.value)}
          />

          <Button
            variant="contained"
            size="large"
            startIcon={<SendRounded />}
            onClick={handleSend}
            disabled={!input.trim() || isLoading}
          >
            {isLoading ? 'Запрос обрабатывается…' : 'Отправить ассистенту'}
          </Button>
        </Stack>
      </SectionCard>

      {output && (
        <>
          <SectionCard
            title="Ответ ассистента"
            description="Результат можно скопировать или прикрепить к задаче"
          >
            <Box
              component="pre"
              sx={{
                m: 0,
                fontFamily: 'JetBrains Mono, Menlo, monospace',
                whiteSpace: 'pre-wrap',
                p: 2,
                borderRadius: 2,
                backgroundColor: (theme) => theme.palette.background.default,
                border: (theme) => `1px solid ${theme.palette.divider}`,
              }}
            >
              {output}
            </Box>
          </SectionCard>

          {useRAG && sources.length > 0 && (
            <SectionCard
              title="Источники информации"
              description="Документы, использованные для формирования ответа"
            >
              <Stack spacing={2}>
                {sources.map((source, index) => (
                  <Card key={index} variant="outlined">
                    <CardContent>
                      <Stack direction="row" spacing={2} alignItems="center" mb={1}>
                        <Chip
                          label={`Релевантность: ${(source.score * 100).toFixed(0)}%`}
                          size="small"
                          color="primary"
                          variant="outlined"
                        />
                        <Typography variant="caption" color="text.secondary">
                          ID: {source.document_id.substring(0, 8)}...
                        </Typography>
                      </Stack>
                      <Typography variant="subtitle2" fontWeight="bold" mb={1}>
                        📄 {source.title}
                      </Typography>
                      <Divider sx={{ my: 1 }} />
                      <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{
                          maxHeight: '150px',
                          overflow: 'auto',
                          fontFamily: 'monospace',
                          fontSize: '0.85rem',
                        }}
                      >
                        {source.content}
                      </Typography>
                    </CardContent>
                  </Card>
                ))}
              </Stack>
            </SectionCard>
          )}
        </>
      )}

      <SectionCard title="Подсказки по профилям">
        <Stack direction="row" spacing={1} flexWrap="wrap">
          {ROLES.map((item) => (
            <Chip key={item.value} label={ROLE_HINTS[item.value]} size="small" />
          ))}
        </Stack>
      </SectionCard>
    </PageContainer>
  )
}
