import React, { useState, useEffect } from 'react';
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
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
  IconButton,
  Tooltip,
  Dialog,
  DialogTitle,
  DialogContent,
} from '@mui/material';
import {
  Add,
  Warning,
  Edit,
  FilterList,
  Download,
  Delete,
  Visibility
} from '@mui/icons-material';
import { incidentsApi, Incident, CreateIncidentRequest, UpdateIncidentRequest } from '../shared/api/incidents';
import { usersApi, User } from '../shared/api/users';

const IncidentsPage: React.FC = () => {
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showViewModal, setShowViewModal] = useState(false);
  const [selectedIncident, setSelectedIncident] = useState<Incident | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [incidentsResponse, usersResponse] = await Promise.all([
        incidentsApi.list({ page: 1, page_size: 20 }),
        usersApi.getUsers()
      ]);
      
      setIncidents(incidentsResponse.data || []);
      setUsers(Array.isArray(usersResponse) ? usersResponse : []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateIncident = async (data: CreateIncidentRequest) => {
    try {
      await incidentsApi.create(data);
      setShowCreateModal(false);
      loadData();
    } catch (err: any) {
      throw new Error(err.response?.data?.error || 'Failed to create incident');
    }
  };

  const handleUpdateIncident = async (data: UpdateIncidentRequest) => {
    try {
      if (selectedIncident) {
        await incidentsApi.update(selectedIncident.id, data);
        setShowEditModal(false);
        setSelectedIncident(null);
        loadData();
      }
    } catch (err: any) {
      throw new Error(err.response?.data?.error || 'Failed to update incident');
    }
  };

  const handleDeleteIncident = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этот инцидент?')) {
      try {
        await incidentsApi.delete(id);
        loadData();
      } catch (err: any) {
        setError(err.response?.data?.error || 'Failed to delete incident');
      }
    }
  };

  const handleViewIncident = (incident: Incident) => {
    setSelectedIncident(incident);
    setShowViewModal(true);
  };

  const handleEditIncident = (incident: Incident) => {
    setSelectedIncident(incident);
    setShowEditModal(true);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'new': return 'info';
      case 'assigned': return 'warning';
      case 'in_progress': return 'warning';
      case 'resolved': return 'success';
      case 'closed': return 'default';
      default: return 'default';
    }
  };

  const getCriticalityColor = (criticality: string) => {
    switch (criticality) {
      case 'low': return 'success';
      case 'medium': return 'warning';
      case 'high': return 'warning';
      case 'critical': return 'error';
      default: return 'default';
    }
  };

  if (loading) {
    return (
      <Container maxWidth="lg">
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Typography variant="h4">Инциденты</Typography>
        </Box>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
          <Box textAlign="center">
            <LinearProgress sx={{ width: 200, mb: 2 }} />
            <Typography>Загрузка инцидентов...</Typography>
          </Box>
        </Box>
      </Container>
    );
  }

  if (error) {
    return (
      <Container maxWidth="lg">
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Typography variant="h4">Инциденты</Typography>
        </Box>
        <Box mb={2} p={2} bgcolor="error.light" borderRadius={1}>
          <Typography color="error">Ошибка: {error}</Typography>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Инциденты</Typography>
        <Box display="flex" gap={1}>
          <Button 
            variant="outlined" 
            startIcon={<FilterList />}
          >
            Фильтры
          </Button>
          <Button 
            variant="outlined" 
            startIcon={<Download />}
          >
            Экспорт
          </Button>
          <Button 
            variant="contained" 
            startIcon={<Add />} 
          onClick={() => setShowCreateModal(true)}
          >
            Добавить инцидент
          </Button>
        </Box>
      </Box>

      {error && (
        <Box mb={2} p={2} bgcolor="error.light" borderRadius={1}>
          <Typography color="error">{error}</Typography>
        </Box>
      )}

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Название</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Критичность</TableCell>
                <TableCell>Категория</TableCell>
                <TableCell>Ответственный</TableCell>
                <TableCell>Дата создания</TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={7} align="center">
                    <LinearProgress />
                    <Typography sx={{ mt: 1 }}>Загрузка инцидентов...</Typography>
                  </TableCell>
                </TableRow>
              ) : incidents.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center">
                    <Typography>Нет инцидентов для отображения.</Typography>
                  </TableCell>
                </TableRow>
              ) : (
                incidents.map((incident) => (
                  <TableRow key={incident.id} hover>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <Warning sx={{ mr: 1 }} />
                        <Box>
                          <Typography variant="body2" fontWeight="medium">
                      {incident.title}
                          </Typography>
                      {incident.description && (
                            <Typography variant="caption" color="text.secondary" sx={{ display: 'block', maxWidth: 300, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                          {incident.description}
                            </Typography>
                          )}
                        </Box>
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={incident.status}
                        color={getStatusColor(incident.status) as any}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={incident.criticality}
                        color={getCriticalityColor(incident.criticality) as any}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                      {incident.category}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                      {incident.assigned_name || 'Не назначен'}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                      {new Date(incident.created_at).toLocaleDateString('ru-RU')}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Box display="flex" gap={0.5}>
                        <Tooltip title="Просмотр деталей">
                          <IconButton 
                            size="small"
                            color="primary"
                          onClick={() => handleViewIncident(incident)}
                          >
                            <Visibility />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Редактировать инцидент">
                          <IconButton 
                            size="small"
                            color="primary"
                          onClick={() => handleEditIncident(incident)}
                          >
                            <Edit />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Удалить инцидент">
                          <IconButton 
                            size="small"
                            color="error"
                          onClick={() => handleDeleteIncident(incident.id)}
                          >
                            <Delete />
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

      {/* Create Modal */}
      {showCreateModal && (
        <CreateIncidentModal
          onClose={() => setShowCreateModal(false)}
          onSave={handleCreateIncident}
          users={users}
        />
      )}

      {/* Edit Modal */}
      {showEditModal && selectedIncident && (
        <EditIncidentModal
          incident={selectedIncident}
          onClose={() => {
            setShowEditModal(false);
            setSelectedIncident(null);
          }}
          onSave={handleUpdateIncident}
          users={users}
        />
      )}

      {/* View Modal */}
      {showViewModal && selectedIncident && (
        <ViewIncidentModal
          incident={selectedIncident}
          onClose={() => {
            setShowViewModal(false);
            setSelectedIncident(null);
          }}
        />
      )}
    </Container>
  );
};

// Create Incident Modal Component
const CreateIncidentModal: React.FC<{
  onClose: () => void;
  onSave: (data: CreateIncidentRequest) => Promise<void>;
  users: User[];
}> = ({ onClose, onSave, users }) => {
  const [formData, setFormData] = useState<CreateIncidentRequest>({
    title: '',
    description: '',
    category: 'technical_failure',
    criticality: 'medium',
    source: 'user_report',
    assigned_to: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await onSave(formData);
    } catch (error: any) {
      alert(error.message);
    }
  };

  return (
    <Dialog open={true} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Typography variant="h6">Создать инцидент</Typography>
      </DialogTitle>
      <DialogContent>
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                required
                label="Название"
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                multiline
                rows={3}
                label="Описание"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Категория</InputLabel>
                <Select
                value={formData.category}
                onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                  label="Категория"
                >
                  <MenuItem value="technical_failure">Технический сбой</MenuItem>
                  <MenuItem value="data_breach">Утечка данных</MenuItem>
                  <MenuItem value="unauthorized_access">Несанкционированный доступ</MenuItem>
                  <MenuItem value="physical">Физический инцидент</MenuItem>
                  <MenuItem value="malware">Вредоносное ПО</MenuItem>
                  <MenuItem value="social_engineering">Социальная инженерия</MenuItem>
                  <MenuItem value="other">Другое</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Критичность</InputLabel>
                <Select
                value={formData.criticality}
                onChange={(e) => setFormData({ ...formData, criticality: e.target.value })}
                  label="Критичность"
                >
                  <MenuItem value="low">Низкая</MenuItem>
                  <MenuItem value="medium">Средняя</MenuItem>
                  <MenuItem value="high">Высокая</MenuItem>
                  <MenuItem value="critical">Критическая</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Источник</InputLabel>
                <Select
                value={formData.source}
                onChange={(e) => setFormData({ ...formData, source: e.target.value })}
                  label="Источник"
                >
                  <MenuItem value="user_report">Сообщение пользователя</MenuItem>
                  <MenuItem value="automatic_agent">Автоматический агент</MenuItem>
                  <MenuItem value="admin_manual">Ручное создание админом</MenuItem>
                  <MenuItem value="monitoring">Мониторинг</MenuItem>
                  <MenuItem value="siem">SIEM</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Ответственный</InputLabel>
                <Select
                value={formData.assigned_to}
                onChange={(e) => setFormData({ ...formData, assigned_to: e.target.value })}
                  label="Ответственный"
              >
                  <MenuItem value="">Не назначен</MenuItem>
                {users.map(user => (
                    <MenuItem key={user.id} value={user.id}>
                    {user.first_name} {user.last_name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
          </Grid>
          <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end', gap: 1 }}>
            <Button onClick={onClose}>
                Отмена
            </Button>
            <Button type="submit" variant="contained">
                Создать
            </Button>
          </Box>
        </Box>
      </DialogContent>
    </Dialog>
  );
};

// Edit Incident Modal Component
const EditIncidentModal: React.FC<{
  incident: Incident;
  onClose: () => void;
  onSave: (data: UpdateIncidentRequest) => Promise<void>;
  users: User[];
}> = ({ incident, onClose, onSave, users }) => {
  const [formData, setFormData] = useState<UpdateIncidentRequest>({
    title: incident.title,
    description: incident.description || '',
    category: incident.category,
    criticality: incident.criticality,
    assigned_to: incident.assigned_to || '',
    status: incident.status,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await onSave(formData);
    } catch (error: any) {
      alert(error.message);
    }
  };

  return (
    <Dialog open={true} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Typography variant="h6">Редактировать инцидент</Typography>
      </DialogTitle>
      <DialogContent>
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                required
                label="Название"
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                multiline
                rows={3}
                label="Описание"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Статус</InputLabel>
                <Select
                value={formData.status}
                onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                  label="Статус"
                >
                  <MenuItem value="new">Новый</MenuItem>
                  <MenuItem value="assigned">Назначен</MenuItem>
                  <MenuItem value="in_progress">В работе</MenuItem>
                  <MenuItem value="resolved">Решен</MenuItem>
                  <MenuItem value="closed">Закрыт</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Критичность</InputLabel>
                <Select
                value={formData.criticality}
                onChange={(e) => setFormData({ ...formData, criticality: e.target.value })}
                  label="Критичность"
                >
                  <MenuItem value="low">Низкая</MenuItem>
                  <MenuItem value="medium">Средняя</MenuItem>
                  <MenuItem value="high">Высокая</MenuItem>
                  <MenuItem value="critical">Критическая</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <FormControl fullWidth>
                <InputLabel>Ответственный</InputLabel>
                <Select
                value={formData.assigned_to}
                onChange={(e) => setFormData({ ...formData, assigned_to: e.target.value })}
                  label="Ответственный"
              >
                  <MenuItem value="">Не назначен</MenuItem>
                {users.map(user => (
                    <MenuItem key={user.id} value={user.id}>
                    {user.first_name} {user.last_name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
          </Grid>
          <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end', gap: 1 }}>
            <Button onClick={onClose}>
                Отмена
            </Button>
            <Button type="submit" variant="contained">
                Сохранить
            </Button>
          </Box>
        </Box>
      </DialogContent>
    </Dialog>
  );
};

// View Incident Modal Component
const ViewIncidentModal: React.FC<{
  incident: Incident;
  onClose: () => void;
}> = ({ incident, onClose }) => {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'new': return 'info';
      case 'assigned': return 'warning';
      case 'in_progress': return 'warning';
      case 'resolved': return 'success';
      case 'closed': return 'default';
      default: return 'default';
    }
  };

  const getCriticalityColor = (criticality: string) => {
    switch (criticality) {
      case 'low': return 'success';
      case 'medium': return 'warning';
      case 'high': return 'warning';
      case 'critical': return 'error';
      default: return 'default';
    }
  };

  return (
    <Dialog open={true} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        <Typography variant="h6">Детали инцидента</Typography>
      </DialogTitle>
      <DialogContent>
        <Grid container spacing={2} sx={{ mt: 1 }}>
          <Grid item xs={12} sm={6}>
            <Typography variant="caption" color="text.secondary">Название</Typography>
            <Typography variant="body2">{incident.title}</Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Typography variant="caption" color="text.secondary">Статус</Typography>
            <Box>
              <Chip
                label={incident.status}
                color={getStatusColor(incident.status) as any}
                size="small"
              />
            </Box>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Typography variant="caption" color="text.secondary">Критичность</Typography>
            <Box>
              <Chip
                label={incident.criticality}
                color={getCriticalityColor(incident.criticality) as any}
                size="small"
              />
            </Box>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Typography variant="caption" color="text.secondary">Категория</Typography>
            <Typography variant="body2">{incident.category}</Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Typography variant="caption" color="text.secondary">Источник</Typography>
            <Typography variant="body2">{incident.source}</Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Typography variant="caption" color="text.secondary">Ответственный</Typography>
            <Typography variant="body2">{incident.assigned_name || 'Не назначен'}</Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Typography variant="caption" color="text.secondary">Дата обнаружения</Typography>
            <Typography variant="body2">
                {new Date(incident.detected_at).toLocaleString('ru-RU')}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Typography variant="caption" color="text.secondary">Дата создания</Typography>
            <Typography variant="body2">
                {new Date(incident.created_at).toLocaleString('ru-RU')}
            </Typography>
          </Grid>
            {incident.resolved_at && (
            <Grid item xs={12} sm={6}>
              <Typography variant="caption" color="text.secondary">Дата решения</Typography>
              <Typography variant="body2">
                  {new Date(incident.resolved_at).toLocaleString('ru-RU')}
              </Typography>
            </Grid>
            )}
            {incident.closed_at && (
            <Grid item xs={12} sm={6}>
              <Typography variant="caption" color="text.secondary">Дата закрытия</Typography>
              <Typography variant="body2">
                  {new Date(incident.closed_at).toLocaleString('ru-RU')}
              </Typography>
            </Grid>
            )}
          {incident.description && (
            <Grid item xs={12}>
              <Typography variant="caption" color="text.secondary">Описание</Typography>
              <Typography variant="body2">{incident.description}</Typography>
            </Grid>
          )}
        </Grid>
        <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end' }}>
          <Button onClick={onClose}>
              Закрыть
          </Button>
        </Box>
      </DialogContent>
    </Dialog>
  );
};

export default IncidentsPage;