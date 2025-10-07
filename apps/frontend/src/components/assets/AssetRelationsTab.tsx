import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Tabs,
  Tab,
  Button,
  List,
  ListItem,
  ListItemText,
  Chip,
  CircularProgress,
  Alert,
  Card,
  Divider,
} from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { risksApi, Risk, RISK_STATUSES } from '../../shared/api/risks';
import { incidentsApi, Incident, INCIDENT_CRITICALITY, INCIDENT_STATUS } from '../../shared/api/incidents';

interface AssetRelationsTabProps {
  assetId: string;
  assetName: string;
}

const AssetRelationsTab: React.FC<AssetRelationsTabProps> = ({ assetId, assetName }) => {
  const [risks, setRisks] = useState<Risk[]>([]);
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'risks' | 'incidents'>('risks');

  useEffect(() => {
    loadRelations();
  }, [assetId]);

  const loadRelations = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const [risksResponse, incidentsResponse] = await Promise.all([
        risksApi.getByAsset(assetId),
        incidentsApi.list({ asset_id: assetId, page_size: 100 })
      ]);
      
      setRisks(risksResponse || []);
      setIncidents(incidentsResponse.data || []);
    } catch (err) {
      setError('Ошибка загрузки связанных данных');
      console.error('Error loading relations:', err);
    } finally {
      setLoading(false);
    }
  };

  const getRiskLevelColor = (level: string) => {
    switch (level) {
      case 'critical': return 'error';
      case 'high': return 'warning';
      case 'medium': return 'info';
      case 'low': return 'success';
      default: return 'default';
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'error';
      case 'high': return 'warning';
      case 'medium': return 'info';
      case 'low': return 'success';
      default: return 'default';
    }
  };

  const getStatusColor = (status: string, isRisk: boolean = true) => {
    if (isRisk) {
      switch (status) {
        case 'open': return 'error';
        case 'mitigated': return 'warning';
        case 'accepted': return 'info';
        case 'closed': return 'success';
        default: return 'default';
      }
    } else {
      switch (status) {
        case 'open': return 'error';
        case 'in_progress': return 'warning';
        case 'resolved': return 'info';
        case 'closed': return 'success';
        default: return 'default';
      }
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU');
  };

  const handleTabChange = (_event: React.SyntheticEvent, newValue: 'risks' | 'incidents') => {
    setActiveTab(newValue);
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ m: 2 }}>
        {error}
      </Alert>
    );
  }

  return (
    <Box sx={{ width: '100%' }}>
      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={activeTab} onChange={handleTabChange} aria-label="relations tabs">
          <Tab 
            label={`Риски (${risks.length})`} 
            value="risks"
            sx={{ textTransform: 'none' }}
          />
          <Tab 
            label={`Инциденты (${incidents.length})`} 
            value="incidents"
            sx={{ textTransform: 'none' }}
          />
        </Tabs>
      </Box>

      {/* Risks Tab */}
      {activeTab === 'risks' && (
        <Box sx={{ mt: 3 }}>
          <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
            <Typography variant="h6" component="h3">
              Риски для актива "{assetName}"
            </Typography>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              size="small"
            >
              Добавить риск
            </Button>
          </Box>
          
          {risks.length > 0 ? (
            <Card>
              <List disablePadding>
                {risks.map((risk, index) => (
                  <React.Fragment key={risk.id}>
                    <ListItem sx={{ py: 2 }}>
                      <ListItemText
                        primary={
                          <Box display="flex" alignItems="center" gap={1} mb={1}>
                            <Typography variant="subtitle1" fontWeight="medium">
                              {risk.title}
                            </Typography>
                            <Chip
                              label={risk.level_label ?? 'Не определен'}
                              color={getRiskLevelColor(String(risk.level_label ?? '').toLowerCase()) as any}
                              size="small"
                            />
                            <Chip
                              label={RISK_STATUSES.find(s => s.value === risk.status)?.label || risk.status}
                              color={getStatusColor(risk.status, true) as any}
                              size="small"
                              variant="outlined"
                            />
                          </Box>
                        }
                        secondary={
                          <Box>
                            {risk.description && (
                              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                                {risk.description}
                              </Typography>
                            )}
                            <Typography variant="caption" color="text.disabled">
                              Создан: {formatDate(risk.created_at)}
                            </Typography>
                          </Box>
                        }
                      />
                      <Box display="flex" gap={1}>
                        <Button size="small" color="primary">
                          Просмотр
                        </Button>
                        <Button size="small" color="success">
                          Редактировать
                        </Button>
                      </Box>
                    </ListItem>
                    {index < risks.length - 1 && <Divider />}
                  </React.Fragment>
                ))}
              </List>
            </Card>
          ) : (
            <Box textAlign="center" py={4}>
              <Typography color="text.secondary">Риски не найдены</Typography>
            </Box>
          )}
        </Box>
      )}

      {/* Incidents Tab */}
      {activeTab === 'incidents' && (
        <Box sx={{ mt: 3 }}>
          <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
            <Typography variant="h6" component="h3">
              Инциденты для актива "{assetName}"
            </Typography>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              size="small"
            >
              Добавить инцидент
            </Button>
          </Box>
          
          {incidents.length > 0 ? (
            <Card>
              <List disablePadding>
                {incidents.map((incident, index) => (
                  <React.Fragment key={incident.id}>
                    <ListItem sx={{ py: 2 }}>
                      <ListItemText
                        primary={
                          <Box display="flex" alignItems="center" gap={1} mb={1}>
                            <Typography variant="subtitle1" fontWeight="medium">
                              {incident.title}
                            </Typography>
                            <Chip
                              label={Object.values(INCIDENT_CRITICALITY).find(s => s === incident.criticality) || incident.criticality}
                              color={getSeverityColor(incident.criticality) as any}
                              size="small"
                            />
                            <Chip
                              label={Object.values(INCIDENT_STATUS).find(s => s === incident.status) || incident.status}
                              color={getStatusColor(incident.status, false) as any}
                              size="small"
                              variant="outlined"
                            />
                          </Box>
                        }
                        secondary={
                          <Box>
                            {incident.description && (
                              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                                {incident.description}
                              </Typography>
                            )}
                            <Typography variant="caption" color="text.disabled">
                              Создан: {formatDate(incident.created_at)}
                              {incident.assigned_name && ` • Назначен: ${incident.assigned_name}`}
                              {incident.resolved_at && ` • Решен: ${formatDate(incident.resolved_at)}`}
                            </Typography>
                          </Box>
                        }
                      />
                      <Box display="flex" gap={1}>
                        <Button size="small" color="primary">
                          Просмотр
                        </Button>
                        <Button size="small" color="success">
                          Редактировать
                        </Button>
                      </Box>
                    </ListItem>
                    {index < incidents.length - 1 && <Divider />}
                  </React.Fragment>
                ))}
              </List>
            </Card>
          ) : (
            <Box textAlign="center" py={4}>
              <Typography color="text.secondary">Инциденты не найдены</Typography>
            </Box>
          )}
        </Box>
      )}
    </Box>
  );
};

export default AssetRelationsTab;

