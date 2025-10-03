import React, { useState, useEffect } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography,
  Grid,
  Chip,
  Card,
  CardContent,
  CircularProgress,
  Alert,
  Tabs,
  Tab,
  IconButton,
} from '@mui/material'
import {
  Close,
  Person,
  Email,
  CalendarToday,
  Security,
  Description,
  Warning,
  BugReport,
  Storage,
} from '@mui/icons-material'
import { UserCatalog, getUserDetail, UserDetail } from '../../shared/api/users'

interface UserDetailModalProps {
  open: boolean
  onClose: () => void
  user: UserCatalog
}

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`user-tabpanel-${index}`}
      aria-labelledby={`user-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  )
}

export const UserDetailModal: React.FC<UserDetailModalProps> = ({
  open,
  onClose,
  user,
}) => {
  const [tabValue, setTabValue] = useState(0)
  const [userDetail, setUserDetail] = useState<UserDetail | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (open && user) {
      loadUserDetail()
    }
  }, [open, user])

  const loadUserDetail = async () => {
    try {
      setLoading(true)
      setError(null)
      const detail = await getUserDetail(user.id)
      setUserDetail(detail)
    } catch (err) {
      setError('Ошибка загрузки детальной информации')
      console.error('Error loading user detail:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue)
  }

  const getStatusChip = (isActive: boolean) => (
    <Chip
      label={isActive ? 'Активен' : 'Заблокирован'}
      color={isActive ? 'success' : 'error'}
      icon={<Person />}
    />
  )

  const getRoleChips = (roles: string[]) => (
    <Box display="flex" gap={1} flexWrap="wrap">
      {roles.map((role, index) => (
        <Chip
          key={index}
          label={role}
          color="primary"
          variant="outlined"
          icon={<Security />}
        />
      ))}
    </Box>
  )

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="md"
      fullWidth
      PaperProps={{
        sx: { height: '80vh' }
      }}
    >
      <DialogTitle>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h6">
            {user.first_name && user.last_name 
              ? `${user.first_name} ${user.last_name}` 
              : user.email
            }
          </Typography>
          <IconButton onClick={onClose} size="small">
            <Close />
          </IconButton>
        </Box>
      </DialogTitle>

      <DialogContent dividers>
        {loading ? (
          <Box display="flex" justifyContent="center" p={4}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Alert severity="error">{error}</Alert>
        ) : (
          <>
            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
              <Tabs value={tabValue} onChange={handleTabChange}>
                <Tab label="Общая информация" />
                <Tab label="Статистика" />
                <Tab label="Активность" />
              </Tabs>
            </Box>

            <TabPanel value={tabValue} index={0}>
              <Grid container spacing={3}>
                <Grid item xs={12} md={6}>
                  <Card>
                    <CardContent>
                      <Typography variant="h6" gutterBottom>
                        Основная информация
                      </Typography>
                      <Box display="flex" flexDirection="column" gap={2}>
                        <Box display="flex" alignItems="center" gap={1}>
                          <Email color="action" />
                          <Typography variant="body2" color="textSecondary">
                            Email:
                          </Typography>
                          <Typography variant="body1">{user.email}</Typography>
                        </Box>
                        
                        {user.first_name && (
                          <Box display="flex" alignItems="center" gap={1}>
                            <Person color="action" />
                            <Typography variant="body2" color="textSecondary">
                              Имя:
                            </Typography>
                            <Typography variant="body1">{user.first_name}</Typography>
                          </Box>
                        )}
                        
                        {user.last_name && (
                          <Box display="flex" alignItems="center" gap={1}>
                            <Person color="action" />
                            <Typography variant="body2" color="textSecondary">
                              Фамилия:
                            </Typography>
                            <Typography variant="body1">{user.last_name}</Typography>
                          </Box>
                        )}
                        
                        <Box display="flex" alignItems="center" gap={1}>
                          <Typography variant="body2" color="textSecondary">
                            Статус:
                          </Typography>
                          {getStatusChip(user.is_active)}
                        </Box>
                      </Box>
                    </CardContent>
                  </Card>
                </Grid>

                <Grid item xs={12} md={6}>
                  <Card>
                    <CardContent>
                      <Typography variant="h6" gutterBottom>
                        Роли и права
                      </Typography>
                      <Box>
                        <Typography variant="body2" color="textSecondary" gutterBottom>
                          Роли:
                        </Typography>
                        {getRoleChips(user.roles)}
                      </Box>
                    </CardContent>
                  </Card>
                </Grid>

                <Grid item xs={12}>
                  <Card>
                    <CardContent>
                      <Typography variant="h6" gutterBottom>
                        Временные метки
                      </Typography>
                      <Grid container spacing={2}>
                        <Grid item xs={12} sm={6}>
                          <Box display="flex" alignItems="center" gap={1}>
                            <CalendarToday color="action" />
                            <Typography variant="body2" color="textSecondary">
                              Создан:
                            </Typography>
                            <Typography variant="body1">
                              {formatDate(user.created_at)}
                            </Typography>
                          </Box>
                        </Grid>
                        <Grid item xs={12} sm={6}>
                          <Box display="flex" alignItems="center" gap={1}>
                            <CalendarToday color="action" />
                            <Typography variant="body2" color="textSecondary">
                              Обновлен:
                            </Typography>
                            <Typography variant="body1">
                              {formatDate(user.updated_at)}
                            </Typography>
                          </Box>
                        </Grid>
                      </Grid>
                    </CardContent>
                  </Card>
                </Grid>
              </Grid>
            </TabPanel>

            <TabPanel value={tabValue} index={1}>
              {userDetail && (
                <Grid container spacing={3}>
                  <Grid item xs={12} sm={6} md={3}>
                    <Card>
                      <CardContent sx={{ textAlign: 'center' }}>
                        <Description color="primary" sx={{ fontSize: 40, mb: 1 }} />
                        <Typography variant="h4" color="primary">
                          {userDetail.stats.documents_count}
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                          Документов
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                  
                  <Grid item xs={12} sm={6} md={3}>
                    <Card>
                      <CardContent sx={{ textAlign: 'center' }}>
                        <Warning color="warning" sx={{ fontSize: 40, mb: 1 }} />
                        <Typography variant="h4" color="warning.main">
                          {userDetail.stats.risks_count}
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                          Рисков
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                  
                  <Grid item xs={12} sm={6} md={3}>
                    <Card>
                      <CardContent sx={{ textAlign: 'center' }}>
                        <BugReport color="error" sx={{ fontSize: 40, mb: 1 }} />
                        <Typography variant="h4" color="error.main">
                          {userDetail.stats.incidents_count}
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                          Инцидентов
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                  
                  <Grid item xs={12} sm={6} md={3}>
                    <Card>
                      <CardContent sx={{ textAlign: 'center' }}>
                        <Storage color="info" sx={{ fontSize: 40, mb: 1 }} />
                        <Typography variant="h4" color="info.main">
                          {userDetail.stats.assets_count}
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                          Активов
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                </Grid>
              )}
            </TabPanel>

            <TabPanel value={tabValue} index={2}>
              <Typography variant="body1" color="textSecondary">
                История активности пользователя будет отображаться здесь.
              </Typography>
              <Typography variant="body2" color="textSecondary" sx={{ mt: 1 }}>
                В будущих версиях здесь будет показана информация о последних действиях,
                изменениях статуса, входе в систему и других активностях.
              </Typography>
            </TabPanel>
          </>
        )}
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose}>Закрыть</Button>
      </DialogActions>
    </Dialog>
  )
}
