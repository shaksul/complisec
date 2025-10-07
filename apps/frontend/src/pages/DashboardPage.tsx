import React from 'react'
import { Grid, Box, Typography, Paper, Stack } from '@mui/material'
import { PageContainer, PageHeader, SectionCard } from '../components/common/Page'
import { RiskHeatmap } from '../components/dashboard/RiskHeatmap'
import { UsersIcon, AssetsIcon, RisksIcon, DocumentsIcon } from '../shared/icons'

const KPI_CARDS = [
  { label: 'Активные пользователи', value: '128', icon: <UsersIcon fontSize="large" /> },
  { label: 'Активы под контролем', value: '864', icon: <AssetsIcon fontSize="large" /> },
  { label: 'Открытые риски', value: '32', icon: <RisksIcon fontSize="large" /> },
  { label: 'Документы в работе', value: '56', icon: <DocumentsIcon fontSize="large" /> },
]

const DashboardStatCard: React.FC<{ label: string; value: string; icon: React.ReactNode }> = ({ label, value, icon }) => (
  <Paper
    elevation={0}
    sx={{
      p: 3,
      borderRadius: 3,
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      gap: 2,
    }}
  >
    <Stack spacing={0.5}>
      <Typography variant="body2" color="text.secondary">
        {label}
      </Typography>
      <Typography variant="h4" fontWeight={700}>
        {value}
      </Typography>
    </Stack>
    <Box
      sx={{
        width: 48,
        height: 48,
        borderRadius: 3,
        backgroundColor: (theme) => theme.palette.primary.main,
        color: 'common.white',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        boxShadow: (theme) => `0 12px 32px ${theme.palette.primary.main}33`,
      }}
    >
      {icon}
    </Box>
  </Paper>
)

export const DashboardPage: React.FC = () => {
  return (
    <PageContainer>
      <PageHeader
        title="Панель мониторинга"
        subtitle="Ключевые показатели по активам, рискам, инцидентам и документации в одном окне"
      />

      <Grid container spacing={3}>
        {KPI_CARDS.map((card) => (
          <Grid item key={card.label} xs={12} sm={6} md={3}>
            <DashboardStatCard {...card} />
          </Grid>
        ))}
      </Grid>

      <SectionCard
        title="Тепловая карта рисков"
        description="Приоритеты обработки по сочетанию вероятности и влияния"
      >
        <RiskHeatmap />
      </SectionCard>

      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <SectionCard
            title="Оперативные задачи"
            description="Фокус команды безопасности на ближайшие сутки"
          >
            <Typography variant="body2" color="text.secondary">
              • Завершить анализ последнего инцидента и обновить план реагирования.
            </Typography>
            <Typography variant="body2" color="text.secondary">
              • Проверить выполнение корректирующих мероприятий по высокорисковым приложениям.
            </Typography>
          </SectionCard>
        </Grid>
        <Grid item xs={12} md={6}>
          <SectionCard
            title="Предстоящие события"
            description="Запланированные обзоры, тренировки и контрольные точки"
          >
            <Typography variant="body2" color="text.secondary">
              • 14:00 — Брифинг по комплаенсу перед инспекцией ISO 27001.
            </Typography>
            <Typography variant="body2" color="text.secondary">
              • 16:30 — Учебная отработка сценария фишинга для службы поддержки.
            </Typography>
          </SectionCard>
        </Grid>
      </Grid>
    </PageContainer>
  )
}
