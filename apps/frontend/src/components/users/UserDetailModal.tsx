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
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
} from 'recharts'
import { UserCatalog, getUserDetail, UserDetail, getUserActivity, getUserActivityStats, getUserAssets, UserActivity, UserActivityStats, UserAsset } from '../../shared/api/users'

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
  const [userActivity, setUserActivity] = useState<UserActivity[]>([])
  const [activityStats, setActivityStats] = useState<UserActivityStats | null>(null)
  const [userAssets, setUserAssets] = useState<UserAsset[]>([])
  const [loading, setLoading] = useState(false)
  const [activityLoading, setActivityLoading] = useState(false)
  const [assetsLoading, setAssetsLoading] = useState(false)
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

  const loadUserActivity = async () => {
    try {
      setActivityLoading(true)
      const [activityResponse, stats] = await Promise.all([
        getUserActivity(user.id, 1, 20),
        getUserActivityStats(user.id)
      ])
      setUserActivity(activityResponse.data)
      setActivityStats(stats)
    } catch (err) {
      console.error('Error loading user activity:', err)
    } finally {
      setActivityLoading(false)
    }
  }

  const loadUserAssets = async () => {
    try {
      setAssetsLoading(true)
      const assets = await getUserAssets(user.id)
      setUserAssets(assets)
    } catch (err) {
      console.error('Error loading user assets:', err)
    } finally {
      setAssetsLoading(false)
    }
  }

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue)
    if (newValue === 1 && !activityStats) {
      loadUserActivity()
    }
    if (newValue === 2 && userActivity.length === 0) {
      loadUserActivity()
    }
    if (newValue === 3 && userAssets.length === 0) {
      loadUserAssets()
    }
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

  const formatDateCompact = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      day: 'numeric',
      month: 'short',
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
                <Tab label="Активы" />
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

                  {/* Дополнительная статистика */}
                  <Grid item xs={12} sm={6} md={3}>
                    <Card>
                      <CardContent sx={{ textAlign: 'center' }}>
                        <Person color="success" sx={{ fontSize: 40, mb: 1 }} />
                        <Typography variant="h4" color="success.main">
                          {userDetail.stats.sessions_count}
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                          Сессий
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} sm={6} md={3}>
                    <Card>
                      <CardContent sx={{ textAlign: 'center' }}>
                        <CalendarToday color="secondary" sx={{ fontSize: 40, mb: 1 }} />
                        <Typography variant="h4" color="secondary.main">
                          {userDetail.stats.login_count}
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                          Входов
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} sm={6} md={3}>
                    <Card>
                      <CardContent sx={{ textAlign: 'center' }}>
                        <Security color="info" sx={{ fontSize: 40, mb: 1 }} />
                        <Typography variant="h4" color="info.main">
                          {userDetail.stats.activity_score}%
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                          Активность
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} sm={6} md={3}>
                    <Card>
                      <CardContent sx={{ textAlign: 'center' }}>
                        <CalendarToday color="action" sx={{ fontSize: 40, mb: 1 }} />
                        <Typography variant="body1" color="textSecondary">
                          {userDetail.last_login 
                            ? formatDate(userDetail.last_login)
                            : 'Никогда'
                          }
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                          Последний вход
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  {/* Графики активности */}
                  {activityStats && activityStats.daily_activity.length > 0 && (
                    <Grid item xs={12}>
                      <Card>
                        <CardContent>
                          <Typography variant="h6" gutterBottom>
                            Активность по дням
                          </Typography>
                          <ResponsiveContainer width="100%" height={300}>
                            <LineChart data={activityStats.daily_activity}>
                              <CartesianGrid strokeDasharray="3 3" />
                              <XAxis 
                                dataKey="date" 
                                tickFormatter={(value) => new Date(value).toLocaleDateString('ru-RU', { month: 'short', day: 'numeric' })}
                              />
                              <YAxis />
                              <Tooltip 
                                labelFormatter={(value) => new Date(value).toLocaleDateString('ru-RU')}
                              />
                              <Line 
                                type="monotone" 
                                dataKey="count" 
                                stroke="#1976d2" 
                                strokeWidth={2}
                                dot={{ fill: '#1976d2', strokeWidth: 2, r: 4 }}
                              />
                            </LineChart>
                          </ResponsiveContainer>
                        </CardContent>
                      </Card>
                    </Grid>
                  )}

                  {/* График популярных действий */}
                  {activityStats && activityStats.top_actions.length > 0 && (
                    <Grid item xs={12} md={6}>
                      <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                        <CardContent sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
                          <Typography variant="h6" gutterBottom>
                            Популярные действия
                          </Typography>
                          <Box sx={{ flexGrow: 1, minHeight: 280 }}>
                            <ResponsiveContainer width="100%" height="100%">
                              <BarChart data={activityStats.top_actions} layout="horizontal">
                                <CartesianGrid strokeDasharray="3 3" />
                                <XAxis type="number" />
                                <YAxis 
                                  type="category" 
                                  dataKey="action" 
                                  width={120}
                                  tick={{ fontSize: 12 }}
                                />
                                <Tooltip 
                                  formatter={(value) => [value, 'Количество']}
                                  labelFormatter={(label) => `Действие: ${label}`}
                                />
                                <Bar dataKey="count" fill="#1976d2" />
                              </BarChart>
                            </ResponsiveContainer>
                          </Box>
                        </CardContent>
                      </Card>
                    </Grid>
                  )}

                  {/* Круговая диаграмма статистики */}
                  <Grid item xs={12} md={6}>
                    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                      <CardContent sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
                        <Typography variant="h6" gutterBottom>
                          Распределение активности
                        </Typography>
                        <Box sx={{ flexGrow: 1, minHeight: 280, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                          <ResponsiveContainer width="100%" height="100%">
                            <PieChart>
                              <Pie
                                data={[
                                  { name: 'Документы', value: userDetail.stats.documents_count, color: '#1976d2' },
                                  { name: 'Риски', value: userDetail.stats.risks_count, color: '#ed6c02' },
                                  { name: 'Инциденты', value: userDetail.stats.incidents_count, color: '#d32f2f' },
                                  { name: 'Активы', value: userDetail.stats.assets_count, color: '#0288d1' },
                                ].filter(item => item.value > 0)}
                                cx="40%"
                                cy="50%"
                                labelLine={false}
                                label={(entry: any) => entry.value > 0 ? `${entry.name}: ${entry.value}` : ''}
                                outerRadius={70}
                                fill="#8884d8"
                                dataKey="value"
                              >
                                {[
                                  { name: 'Документы', value: userDetail.stats.documents_count, color: '#1976d2' },
                                  { name: 'Риски', value: userDetail.stats.risks_count, color: '#ed6c02' },
                                  { name: 'Инциденты', value: userDetail.stats.incidents_count, color: '#d32f2f' },
                                  { name: 'Активы', value: userDetail.stats.assets_count, color: '#0288d1' },
                                ].filter(item => item.value > 0).map((entry, index) => (
                                  <Cell key={`cell-${index}`} fill={entry.color} />
                                ))}
                              </Pie>
                              <Tooltip 
                                formatter={(value, name) => [value, name]}
                                labelFormatter={(label) => `${label}`}
                              />
                            </PieChart>
                          </ResponsiveContainer>
                          {/* Легенда справа */}
                          <Box sx={{ ml: 2, minWidth: 120 }}>
                            {[
                              { name: 'Документы', value: userDetail.stats.documents_count, color: '#1976d2' },
                              { name: 'Риски', value: userDetail.stats.risks_count, color: '#ed6c02' },
                              { name: 'Инциденты', value: userDetail.stats.incidents_count, color: '#d32f2f' },
                              { name: 'Активы', value: userDetail.stats.assets_count, color: '#0288d1' },
                            ].map((item, index) => (
                              <Box key={index} display="flex" alignItems="center" sx={{ mb: 1 }}>
                                <Box
                                  sx={{
                                    width: 12,
                                    height: 12,
                                    backgroundColor: item.color,
                                    borderRadius: '50%',
                                    mr: 1
                                  }}
                                />
                                <Typography variant="body2" sx={{ fontSize: '0.75rem' }}>
                                  {item.name}: {item.value}
                                </Typography>
                              </Box>
                            ))}
                          </Box>
                        </Box>
                      </CardContent>
                    </Card>
                  </Grid>
                </Grid>
              )}
            </TabPanel>

            <TabPanel value={tabValue} index={2}>
              {activityLoading ? (
                <Box display="flex" justifyContent="center" p={4}>
                  <CircularProgress />
                </Box>
              ) : (
                <Grid container spacing={3}>
                  {/* История активности */}
                  <Grid item xs={12} md={8}>
                    <Card>
                      <CardContent>
                        <Typography variant="h6" gutterBottom>
                          История активности
                        </Typography>
                        {userActivity.length === 0 ? (
                          <Typography variant="body2" color="textSecondary">
                            Активность не найдена
                          </Typography>
                        ) : (
                          <Box>
                            {userActivity.map((activity) => (
                              <Box key={activity.id} sx={{ mb: 2, p: 2, border: '1px solid', borderColor: 'divider', borderRadius: 1 }}>
                                <Box display="flex" justifyContent="space-between" alignItems="flex-start">
                                  <Box>
                                    <Typography variant="body1" fontWeight="medium">
                                      {activity.action}
                                    </Typography>
                                    <Typography variant="body2" color="textSecondary">
                                      {activity.description}
                                    </Typography>
                                    {activity.ip_address && (
                                      <Typography variant="caption" color="textSecondary">
                                        IP: {activity.ip_address}
                                      </Typography>
                                    )}
                                  </Box>
                                  <Typography variant="caption" color="textSecondary">
                                    {formatDate(activity.created_at)}
                                  </Typography>
                                </Box>
                              </Box>
                            ))}
                          </Box>
                        )}
                      </CardContent>
                    </Card>
                  </Grid>

                  {/* Статистика активности */}
                  <Grid item xs={12} md={4}>
                    <Card>
                      <CardContent>
                        <Typography variant="h6" gutterBottom>
                          Статистика активности
                        </Typography>
                        {activityStats ? (
                          <Box>
                            <Typography variant="subtitle2" gutterBottom>
                              Популярные действия:
                            </Typography>
                            {activityStats.top_actions.map((action, index) => (
                              <Box key={index} display="flex" justifyContent="space-between" sx={{ mb: 1 }}>
                                <Typography variant="body2">{action.action}</Typography>
                                <Typography variant="body2" color="textSecondary">
                                  {action.count}
                                </Typography>
                              </Box>
                            ))}
                          </Box>
                        ) : (
                          <Typography variant="body2" color="textSecondary">
                            Статистика недоступна
                          </Typography>
                        )}
                      </CardContent>
                    </Card>

                    {/* История входов */}
                    <Card sx={{ mt: 2 }}>
                      <CardContent>
                        <Typography variant="h6" gutterBottom>
                          История входов
                        </Typography>
                        {activityStats?.login_history && activityStats.login_history.length > 0 ? (
                          <Box>
                            {activityStats.login_history.slice(0, 5).map((login, index) => (
                              <Box key={index} sx={{ mb: 2, p: 2, border: '1px solid', borderColor: 'divider', borderRadius: 1 }}>
                                <Box display="flex" flexDirection="column" gap={1}>
                                  {/* IP адрес и статус в одной строке */}
                                  <Box display="flex" justifyContent="space-between" alignItems="center">
                                    <Typography variant="body2" fontWeight="medium">
                                      {login.ip_address}
                                    </Typography>
                                    <Chip
                                      label={login.success ? 'Успешно' : 'Ошибка'}
                                      color={login.success ? 'success' : 'error'}
                                      size="small"
                                    />
                                  </Box>
                                  
                                  {/* User Agent */}
                                  <Typography variant="caption" color="textSecondary" sx={{ wordBreak: 'break-all' }}>
                                    {login.user_agent}
                                  </Typography>
                                  
                                  {/* Дата и время в отдельной строке */}
                                  <Typography variant="caption" color="textSecondary" sx={{ fontWeight: 'medium' }}>
                                    {formatDateCompact(login.created_at)}
                                  </Typography>
                                </Box>
                              </Box>
                            ))}
                          </Box>
                        ) : (
                          <Typography variant="body2" color="textSecondary">
                            История входов недоступна
                          </Typography>
                        )}
                      </CardContent>
                    </Card>
                  </Grid>
                </Grid>
              )}
            </TabPanel>

            <TabPanel value={tabValue} index={3}>
              {assetsLoading ? (
                <Box display="flex" justifyContent="center" p={4}>
                  <CircularProgress />
                </Box>
              ) : (
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Активы под ответственностью ({userAssets.length})
                    </Typography>
                    {userAssets.length === 0 ? (
                      <Typography variant="body2" color="textSecondary">
                        Нет активов под ответственностью
                      </Typography>
                    ) : (
                      <Box sx={{ mt: 2 }}>
                        {userAssets.map((asset) => (
                          <Box
                            key={asset.id}
                            sx={{
                              mb: 2,
                              p: 2,
                              border: '1px solid',
                              borderColor: 'divider',
                              borderRadius: 1,
                              '&:hover': {
                                backgroundColor: 'action.hover',
                              },
                            }}
                          >
                            <Box display="flex" justifyContent="space-between" alignItems="flex-start" mb={1}>
                              <Box flex={1}>
                                <Typography variant="h6" gutterBottom>
                                  {asset.name}
                                </Typography>
                                <Box display="flex" gap={1} flexWrap="wrap" mb={1}>
                                  <Chip label={asset.inventory_number} size="small" variant="outlined" />
                                  <Chip label={asset.type} size="small" color="primary" />
                                  <Chip 
                                    label={asset.criticality} 
                                    size="small" 
                                    color={
                                      asset.criticality === 'Высокая' ? 'error' :
                                      asset.criticality === 'Средняя' ? 'warning' :
                                      'success'
                                    }
                                  />
                                  <Chip 
                                    label={asset.status} 
                                    size="small" 
                                    color={
                                      asset.status === 'Активен' ? 'success' :
                                      asset.status === 'В ремонте' ? 'warning' :
                                      'default'
                                    }
                                  />
                                </Box>
                                <Typography variant="caption" color="textSecondary">
                                  Создан: {formatDate(asset.created_at)}
                                </Typography>
                              </Box>
                            </Box>
                          </Box>
                        ))}
                      </Box>
                    )}
                  </CardContent>
                </Card>
              )}
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
