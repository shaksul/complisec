import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Tabs,
  Tab,
  Box,
  Typography,
  IconButton,
  CircularProgress,
  Button,
} from '@mui/material';
import { Close, Add as AddIcon } from '@mui/icons-material';
import { AssetWithDetails, AssetDocument, AssetSoftware, AssetHistory, DOCUMENT_TYPES } from '../../shared/api/assets';
import { assetsApi } from '../../shared/api/assets';
import AssetRelationsTab from './AssetRelationsTab';
import { DocumentUploadModal } from './DocumentUploadModal';
import { AddSoftwareModal } from './AddSoftwareModal';

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
      id={`asset-tabpanel-${index}`}
      aria-labelledby={`asset-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          {children}
        </Box>
      )}
    </div>
  )
}

interface AssetDetailsModalProps {
  assetId: string;
  onClose: () => void;
}

const AssetDetailsModal: React.FC<AssetDetailsModalProps> = ({ assetId, onClose }) => {
  const [asset, setAsset] = useState<AssetWithDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [tabValue, setTabValue] = useState(0);
  const [documentUploadModalOpen, setDocumentUploadModalOpen] = useState(false);
  const [softwareModalOpen, setSoftwareModalOpen] = useState(false);
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
  const [documentToDelete, setDocumentToDelete] = useState<AssetDocument | null>(null);

  useEffect(() => {
    loadAssetDetails();
  }, [assetId]);

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const handleClose = () => {
    setTabValue(0);
    onClose();
  };

  const loadAssetDetails = async () => {
    try {
      setLoading(true);
      setError(null);
      const details = await assetsApi.getDetails(assetId);
      setAsset(details);
    } catch (err) {
      setError('Ошибка загрузки деталей актива');
      console.error('Error loading asset details:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleDocumentAdded = () => {
    loadAssetDetails();
  };

  const handleSoftwareAdded = () => {
    loadAssetDetails();
  };

  const handleDownloadDocument = async (doc: AssetDocument) => {
    try {
      const blob = await assetsApi.downloadDocument(doc.id);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = doc.title || 'document';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Error downloading document:', err);
      setError('Ошибка скачивания документа');
    }
  };

  const handleDeleteDocument = (doc: AssetDocument) => {
    setDocumentToDelete(doc);
    setDeleteConfirmOpen(true);
  };

  const confirmDeleteDocument = async () => {
    if (!documentToDelete) return;
    
    try {
      await assetsApi.deleteDocument(documentToDelete.id);
      setDeleteConfirmOpen(false);
      setDocumentToDelete(null);
      await loadAssetDetails();
    } catch (err) {
      console.error('Error deleting document:', err);
      setError('Ошибка удаления документа');
      setDeleteConfirmOpen(false);
    }
  };

  const cancelDeleteDocument = () => {
    setDeleteConfirmOpen(false);
    setDocumentToDelete(null);
  };

  const getCriticalityColor = (criticality: string) => {
    switch (criticality) {
      case 'high': return 'text-red-600 bg-red-100';
      case 'medium': return 'text-yellow-600 bg-yellow-100';
      case 'low': return 'text-green-600 bg-green-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'text-green-600 bg-green-100';
      case 'in_repair': return 'text-yellow-600 bg-yellow-100';
      case 'storage': return 'text-blue-600 bg-blue-100';
      case 'decommissioned': return 'text-red-600 bg-red-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getDocumentTypeLabel = (type: string) => {
    return DOCUMENT_TYPES.find(dt => dt.value === type)?.label || type;
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU');
  };

  if (!asset) {
    return null;
  }

  return (
    <Dialog 
      open={true} 
      onClose={handleClose} 
      maxWidth="lg" 
      fullWidth
      fullScreen
    >
      <DialogTitle>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Box>
            <Typography variant="h5" component="div">
              {asset.name}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Инв. номер: {asset.inventory_number}
            </Typography>
          </Box>
          <IconButton onClick={handleClose}>
            <Close />
          </IconButton>
        </Box>
      </DialogTitle>
      
      <DialogContent sx={{ p: 0 }}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs 
            value={tabValue} 
            onChange={handleTabChange} 
            aria-label="asset details tabs"
            variant="scrollable"
            scrollButtons="auto"
          >
            <Tab label="Обзор" />
            <Tab label={`Документы (${asset.documents?.length || 0})`} />
            <Tab label={`ПО (${asset.software?.length || 0})`} />
            <Tab label={`История (${asset.history?.length || 0})`} />
            <Tab label="Связи" />
          </Tabs>
        </Box>

        <TabPanel value={tabValue} index={0}>
          {loading ? (
            <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
              <CircularProgress />
            </Box>
          ) : error ? (
            <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
              <Typography color="error">{error}</Typography>
            </Box>
          ) : (
            <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 3 }}>
              {/* Основная информация */}
              <Box>
                <Typography variant="h6" sx={{ mb: 2 }}>Основная информация</Typography>
                <Box sx={{ bgcolor: 'grey.50', p: 2, borderRadius: 1 }}>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary">Тип</Typography>
                    <Typography variant="body2" sx={{ textTransform: 'capitalize' }}>{asset.type}</Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary">Класс</Typography>
                    <Typography variant="body2" sx={{ textTransform: 'capitalize' }}>{asset.class}</Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary">Владелец</Typography>
                    <Typography variant="body2">{asset.owner_name || 'Не назначен'}</Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary">Местоположение</Typography>
                    <Typography variant="body2">{asset.location || 'Не указано'}</Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary">Статус</Typography>
                    <Box>
                      <Typography 
                        variant="caption" 
                        sx={{ 
                          px: 1, 
                          py: 0.5, 
                          borderRadius: 1, 
                          bgcolor: getStatusColor(asset.status).includes('green') ? 'success.light' : 
                                   getStatusColor(asset.status).includes('yellow') ? 'warning.light' : 
                                   getStatusColor(asset.status).includes('red') ? 'error.light' : 'grey.300',
                          color: getStatusColor(asset.status).includes('green') ? 'success.dark' : 
                                 getStatusColor(asset.status).includes('yellow') ? 'warning.dark' : 
                                 getStatusColor(asset.status).includes('red') ? 'error.dark' : 'grey.700'
                        }}
                      >
                      {asset.status}
                      </Typography>
                    </Box>
                  </Box>
                </Box>
              </Box>

              {/* CIA Оценка */}
              <Box>
                <Typography variant="h6" sx={{ mb: 2 }}>CIA Оценка</Typography>
                <Box sx={{ bgcolor: 'grey.50', p: 2, borderRadius: 1 }}>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary">Критичность</Typography>
                    <Box>
                      <Typography 
                        variant="caption" 
                        sx={{ 
                          px: 1, 
                          py: 0.5, 
                          borderRadius: 1, 
                          bgcolor: getCriticalityColor(asset.criticality).includes('red') ? 'error.light' : 
                                   getCriticalityColor(asset.criticality).includes('yellow') ? 'warning.light' : 
                                   getCriticalityColor(asset.criticality).includes('green') ? 'success.light' : 'grey.300',
                          color: getCriticalityColor(asset.criticality).includes('red') ? 'error.dark' : 
                                 getCriticalityColor(asset.criticality).includes('yellow') ? 'warning.dark' : 
                                 getCriticalityColor(asset.criticality).includes('green') ? 'success.dark' : 'grey.700'
                        }}
                      >
                      {asset.criticality}
                      </Typography>
                    </Box>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary">Конфиденциальность</Typography>
                    <Box>
                      <Typography 
                        variant="caption" 
                        sx={{ 
                          px: 1, 
                          py: 0.5, 
                          borderRadius: 1, 
                          bgcolor: getCriticalityColor(asset.confidentiality).includes('red') ? 'error.light' : 
                                   getCriticalityColor(asset.confidentiality).includes('yellow') ? 'warning.light' : 
                                   getCriticalityColor(asset.confidentiality).includes('green') ? 'success.light' : 'grey.300',
                          color: getCriticalityColor(asset.confidentiality).includes('red') ? 'error.dark' : 
                                 getCriticalityColor(asset.confidentiality).includes('yellow') ? 'warning.dark' : 
                                 getCriticalityColor(asset.confidentiality).includes('green') ? 'success.dark' : 'grey.700'
                        }}
                      >
                      {asset.confidentiality}
                      </Typography>
                    </Box>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary">Целостность</Typography>
                    <Box>
                      <Typography 
                        variant="caption" 
                        sx={{ 
                          px: 1, 
                          py: 0.5, 
                          borderRadius: 1, 
                          bgcolor: getCriticalityColor(asset.integrity).includes('red') ? 'error.light' : 
                                   getCriticalityColor(asset.integrity).includes('yellow') ? 'warning.light' : 
                                   getCriticalityColor(asset.integrity).includes('green') ? 'success.light' : 'grey.300',
                          color: getCriticalityColor(asset.integrity).includes('red') ? 'error.dark' : 
                                 getCriticalityColor(asset.integrity).includes('yellow') ? 'warning.dark' : 
                                 getCriticalityColor(asset.integrity).includes('green') ? 'success.dark' : 'grey.700'
                        }}
                      >
                      {asset.integrity}
                      </Typography>
                    </Box>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary">Доступность</Typography>
                    <Box>
                      <Typography 
                        variant="caption" 
                        sx={{ 
                          px: 1, 
                          py: 0.5, 
                          borderRadius: 1, 
                          bgcolor: getCriticalityColor(asset.availability).includes('red') ? 'error.light' : 
                                   getCriticalityColor(asset.availability).includes('yellow') ? 'warning.light' : 
                                   getCriticalityColor(asset.availability).includes('green') ? 'success.light' : 'grey.300',
                          color: getCriticalityColor(asset.availability).includes('red') ? 'error.dark' : 
                                 getCriticalityColor(asset.availability).includes('yellow') ? 'warning.dark' : 
                                 getCriticalityColor(asset.availability).includes('green') ? 'success.dark' : 'grey.700'
                        }}
                      >
                      {asset.availability}
                      </Typography>
                    </Box>
                  </Box>
                </Box>
              </Box>

              {/* Метаданные */}
              <Box sx={{ gridColumn: { xs: '1', md: '1 / -1' } }}>
                <Typography variant="h6" sx={{ mb: 2 }}>Метаданные</Typography>
                <Box sx={{ bgcolor: 'grey.50', p: 2, borderRadius: 1, display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 2 }}>
                  <Box>
                    <Typography variant="caption" color="text.secondary">Создан</Typography>
                    <Typography variant="body2">{formatDate(asset.created_at)}</Typography>
                  </Box>
                  <Box>
                    <Typography variant="caption" color="text.secondary">Обновлен</Typography>
                    <Typography variant="body2">{formatDate(asset.updated_at)}</Typography>
                  </Box>
                </Box>
              </Box>
            </Box>
          )}
        </TabPanel>

        <TabPanel value={tabValue} index={1}>
          <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
            <Typography variant="h6">Документы актива</Typography>
            <Button
              variant="contained"
              size="small"
              startIcon={<AddIcon />}
              onClick={() => setDocumentUploadModalOpen(true)}
            >
              Добавить документ
            </Button>
          </Box>
          {asset.documents && asset.documents.length > 0 ? (
            <Box sx={{ bgcolor: 'background.paper', borderRadius: 1, overflow: 'hidden' }}>
              {asset.documents.map((doc: AssetDocument) => (
                <Box key={doc.id} sx={{ p: 2, borderBottom: 1, borderColor: 'divider', '&:last-child': { borderBottom: 0 } }}>
                  <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Box display="flex" alignItems="center">
                      <Box sx={{ width: 32, height: 32, bgcolor: 'grey.200', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', mr: 2 }}>
                        <Typography variant="caption" color="text.secondary">DOC</Typography>
                      </Box>
                      <Box>
                        <Typography variant="body2" fontWeight="medium">
                          {getDocumentTypeLabel(doc.document_type)}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          {formatDate(doc.created_at)}
                        </Typography>
                      </Box>
                    </Box>
                    <Box display="flex" gap={1}>
                      <Button 
                        size="small" 
                        color="primary"
                        onClick={() => handleDownloadDocument(doc)}
                      >
                        Скачать
                      </Button>
                      <Button 
                        size="small" 
                        color="error"
                        onClick={() => handleDeleteDocument(doc)}
                      >
                        Удалить
                      </Button>
                    </Box>
                  </Box>
                </Box>
              ))}
            </Box>
          ) : (
            <Box textAlign="center" py={4}>
              <Typography color="text.secondary">Документы не найдены</Typography>
            </Box>
          )}
        </TabPanel>

        <TabPanel value={tabValue} index={2}>
          <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
            <Typography variant="h6">Установленное ПО</Typography>
            <Button
              variant="contained"
              size="small"
              startIcon={<AddIcon />}
              onClick={() => setSoftwareModalOpen(true)}
            >
              Добавить ПО
            </Button>
          </Box>
          {asset.software && asset.software.length > 0 ? (
            <Box sx={{ bgcolor: 'background.paper', borderRadius: 1, overflow: 'hidden' }}>
              {asset.software.map((sw: AssetSoftware) => (
                <Box key={sw.id} sx={{ p: 2, borderBottom: 1, borderColor: 'divider', '&:last-child': { borderBottom: 0 } }}>
                  <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Box>
                      <Typography variant="body2" fontWeight="medium">
                        {sw.software_name}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {sw.version && `Версия: ${sw.version}`}
                        {sw.installed_at && ` • Установлено: ${formatDate(sw.installed_at)}`}
                      </Typography>
                    </Box>
                    <Box display="flex" gap={1}>
                      <Button size="small" color="primary">
                        Редактировать
                      </Button>
                      <Button size="small" color="error">
                        Удалить
                      </Button>
                    </Box>
                  </Box>
                </Box>
              ))}
            </Box>
          ) : (
            <Box textAlign="center" py={4}>
              <Typography color="text.secondary">ПО не найдено</Typography>
            </Box>
          )}
        </TabPanel>

        <TabPanel value={tabValue} index={3}>
          <Typography variant="h6" sx={{ mb: 2 }}>История изменений</Typography>
              {asset.history && asset.history.length > 0 ? (
            <Box sx={{ bgcolor: 'background.paper', borderRadius: 1, overflow: 'hidden' }}>
                    {asset.history.map((entry: AssetHistory) => (
                <Box key={entry.id} sx={{ p: 2, borderBottom: 1, borderColor: 'divider', '&:last-child': { borderBottom: 0 } }}>
                  <Box display="flex" alignItems="flex-start">
                    <Box sx={{ width: 32, height: 32, bgcolor: 'grey.200', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', mr: 2 }}>
                      <Typography variant="caption" color="text.secondary">H</Typography>
                    </Box>
                    <Box flex={1}>
                      <Typography variant="body2" fontWeight="medium">
                              {entry.field_changed}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                              {entry.old_value && `Было: ${entry.old_value}`}
                              {entry.old_value && entry.new_value && ' → '}
                              {entry.new_value && `Стало: ${entry.new_value}`}
                      </Typography>
                      <Typography variant="caption" color="text.disabled" display="block" sx={{ mt: 0.5 }}>
                              {formatDate(entry.changed_at)} • {entry.changed_by}
                      </Typography>
                    </Box>
                  </Box>
                </Box>
              ))}
            </Box>
          ) : (
            <Box textAlign="center" py={4}>
              <Typography color="text.secondary">История изменений пуста</Typography>
            </Box>
          )}
        </TabPanel>

        <TabPanel value={tabValue} index={4}>
            <AssetRelationsTab
              assetId={asset.id}
              assetName={asset.name}
            />
        </TabPanel>
      </DialogContent>

      {/* Модальные окна */}
      <DocumentUploadModal
        open={documentUploadModalOpen}
        onClose={() => setDocumentUploadModalOpen(false)}
        assetId={assetId}
        onSuccess={handleDocumentAdded}
      />
      
      <AddSoftwareModal
        open={softwareModalOpen}
        onClose={() => setSoftwareModalOpen(false)}
        assetId={assetId}
        onSuccess={handleSoftwareAdded}
      />

      {/* Диалог подтверждения удаления документа */}
      <Dialog
        open={deleteConfirmOpen}
        onClose={cancelDeleteDocument}
        maxWidth="xs"
        fullWidth
      >
        <DialogTitle>Подтверждение удаления</DialogTitle>
        <DialogContent>
          <Typography>
            Вы уверены, что хотите удалить документ "{documentToDelete ? getDocumentTypeLabel(documentToDelete.document_type) : ''}"?
          </Typography>
          <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
            Это действие удалит только связь документа с активом. Сам документ останется в хранилище.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={cancelDeleteDocument} color="inherit">
            Отмена
          </Button>
          <Button onClick={confirmDeleteDocument} color="error" variant="contained">
            Удалить
          </Button>
        </DialogActions>
      </Dialog>
    </Dialog>
  );
};

export default AssetDetailsModal;
