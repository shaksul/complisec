import React from 'react'
import {
  Container,
  Grid,
  Paper,
  Typography,
  Box,
  Card,
  CardContent,
} from '@mui/material'
import {
  People,
  Computer,
  Warning,
  Description,
  Report,
  School,
} from '@mui/icons-material'

const StatCard: React.FC<{
  title: string
  value: string | number
  icon: React.ReactNode
  color: string
}> = ({ title, value, icon, color }) => (
  <Card>
    <CardContent>
      <Box display="flex" alignItems="center" justifyContent="space-between">
        <Box>
          <Typography color="textSecondary" gutterBottom variant="h6">
            {title}
          </Typography>
          <Typography variant="h4" component="div">
            {value}
          </Typography>
        </Box>
        <Box color={color} fontSize="3rem">
          {icon}
        </Box>
      </Box>
    </CardContent>
  </Card>
)

export const DashboardPage: React.FC = () => {
  return (
    <Container maxWidth="lg">
      <Typography variant="h4" gutterBottom>
        Дашборд
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Пользователи"
            value="12"
            icon={<People />}
            color="primary.main"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Активы"
            value="45"
            icon={<Computer />}
            color="success.main"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Риски"
            value="8"
            icon={<Warning />}
            color="warning.main"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Инциденты"
            value="3"
            icon={<Report />}
            color="error.main"
          />
        </Grid>
        
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              Последние риски
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Здесь будет список последних рисков...
            </Typography>
          </Paper>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              Активные инциденты
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Здесь будет список активных инцидентов...
            </Typography>
          </Paper>
        </Grid>
      </Grid>
    </Container>
  )
}
