import React, { useState } from 'react'
import {
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  Box,
  Alert,
  Stack,
} from '@mui/material'
import { useAuth } from '../contexts/AuthContext'
import { useNavigate } from 'react-router-dom'

export const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    setIsLoading(true)
    setError('')

    try {
      await login(email, password)
      navigate('/dashboard')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Не удалось выполнить вход. Проверьте данные и попробуйте снова.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Container component="main" maxWidth="sm" sx={{ display: 'flex', alignItems: 'center', minHeight: '100vh' }}>
      <Paper elevation={0} sx={{ p: 5, borderRadius: 4, width: '100%' }}>
        <Stack spacing={3}>
          <Box textAlign="center">
            <Typography component="h1" variant="h4" fontWeight={700} gutterBottom>
              RiskNexus
            </Typography>
            <Typography variant="body1" color="text.secondary">
              Единый контур управления рисками, инцидентами и комплаенсом
            </Typography>
          </Box>

          {error && (
            <Alert severity="error">{error}</Alert>
          )}

          <Box component="form" onSubmit={handleSubmit}>
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="Рабочий e-mail"
              name="email"
              autoComplete="email"
              autoFocus
              value={email}
              onChange={(event) => setEmail(event.target.value)}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="Пароль"
              type="password"
              id="password"
              autoComplete="current-password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              size="large"
              sx={{ mt: 3 }}
              disabled={isLoading}
            >
              {isLoading ? 'Выполняем вход…' : 'Войти в систему'}
            </Button>
          </Box>

          <Typography variant="body2" color="text.secondary" textAlign="center">
            Демо-доступ: admin@demo.local / admin123
          </Typography>
        </Stack>
      </Paper>
    </Container>
  )
}
