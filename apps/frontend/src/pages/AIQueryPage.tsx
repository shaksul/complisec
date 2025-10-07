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
} from '@mui/material'
import { SendRounded, PsychologyAltRounded } from '@mui/icons-material'
import { queryAI } from '@/shared/api/ai'
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

  async function handleSend() {
    setIsLoading(true)
    try {
      const response = await queryAI({
        provider_id: 'demo',
        role,
        input,
        context: {},
      })
      setOutput(response.output)
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

          <TextField
            label="Опишите задачу"
            multiline
            minRows={5}
            value={input}
            placeholder="Например: подготовь краткий обзор требований ISO 27001 для отчёта руководству"
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
        <SectionCard title="Ответ ассистента" description="Результат можно скопировать или прикрепить к задаче">
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
