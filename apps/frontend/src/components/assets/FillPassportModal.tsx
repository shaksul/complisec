import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography,
  Alert,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  Tabs,
  Tab,
  Paper,
} from '@mui/material';
import {
  Close as CloseIcon,
  Visibility as PreviewIcon,
  Download as DownloadIcon,
  Save as SaveIcon,
} from '@mui/icons-material';
import { templatesApi, DocumentTemplate, TEMPLATE_TYPE_LABELS } from '../../shared/api/templates';
import { Asset } from '../../shared/api/assets';

interface FillPassportModalProps {
  open: boolean;
  onClose: () => void;
  asset: Asset;
  onSuccess: () => void;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index } = props;
  return (
    <div role="tabpanel" hidden={value !== index} style={{ flex: 1, overflow: 'auto' }}>
      {value === index && <Box sx={{ pt: 2 }}>{children}</Box>}
    </div>
  );
}

export const FillPassportModal: React.FC<FillPassportModalProps> = ({
  open,
  onClose,
  asset,
  onSuccess: _onSuccess
}) => {
  const [templates, setTemplates] = useState<DocumentTemplate[]>([]);
  const [selectedTemplateId, setSelectedTemplateId] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [tabValue, setTabValue] = useState(0);
  const [preview, setPreview] = useState<string>('');
  const [additionalData, setAdditionalData] = useState<Record<string, string>>({});

  useEffect(() => {
    if (open) {
      loadTemplates();
    }
  }, [open]);

  const loadTemplates = async () => {
    try {
      setLoading(true);
      const data = await templatesApi.listTemplates({ is_active: true });
      setTemplates(data.filter(t => t.template_type.startsWith('passport')));
    } catch (err: any) {
      console.error('Error loading templates:', err);
      setError(err.message || 'Ошибка загрузки шаблонов');
    } finally {
      setLoading(false);
    }
  };

  const handleGeneratePreview = async () => {
    if (!selectedTemplateId) {
      setError('Выберите шаблон');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const response = await templatesApi.fillTemplate(asset.id, {
        template_id: selectedTemplateId,
        asset_id: asset.id,
        additional_data: additionalData,
        save_as_document: false,
        generate_pdf: false,
      });
      
      if (response.html) {
        setPreview(response.html);
        setTabValue(1); // Switch to preview tab
      }
    } catch (err: any) {
      console.error('Error generating preview:', err);
      setError(err.message || 'Ошибка генерации предпросмотра');
    } finally {
      setLoading(false);
    }
  };

  const handleDownloadPDF = async () => {
    if (!selectedTemplateId) {
      setError('Выберите шаблон');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const response = await templatesApi.fillTemplate(asset.id, {
        template_id: selectedTemplateId,
        asset_id: asset.id,
        additional_data: additionalData,
        save_as_document: false,
        generate_pdf: true,
      });
      
      if (response.pdf_base64) {
        // Decode base64 and download
        const binaryString = atob(response.pdf_base64);
        const bytes = new Uint8Array(binaryString.length);
        for (let i = 0; i < binaryString.length; i++) {
          bytes[i] = binaryString.charCodeAt(i);
        }
        const blob = new Blob([bytes], { type: 'application/pdf' });
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `Паспорт_${asset.inventory_number}_${new Date().toISOString().split('T')[0]}.pdf`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);
      }
    } catch (err: any) {
      console.error('Error generating PDF:', err);
      setError(err.message || 'Ошибка генерации PDF');
    } finally {
      setLoading(false);
    }
  };

  const handleSaveAsDocument = async () => {
    if (!selectedTemplateId) {
      setError('Выберите шаблон');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      await templatesApi.fillTemplate(asset.id, {
        template_id: selectedTemplateId,
        asset_id: asset.id,
        additional_data: additionalData,
        save_as_document: true,
        document_title: `Паспорт ${asset.name}`,
        generate_pdf: true, // ✅ Save as PDF instead of HTML
      });
      
      alert('Документ успешно сохранён');
      _onSuccess();
      handleClose();
    } catch (err: any) {
      console.error('Error saving document:', err);
      setError(err.message || 'Ошибка сохранения документа');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setSelectedTemplateId('');
    setPreview('');
    setTabValue(0);
    setAdditionalData({});
    setError(null);
    onClose();
  };

  return (
    <Dialog 
      open={open} 
      onClose={handleClose}
      maxWidth="lg"
      fullWidth
      fullScreen
    >
      <DialogTitle>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h6">Создать паспорт для актива: {asset.name}</Typography>
          <Button onClick={handleClose} startIcon={<CloseIcon />}>
            Закрыть
          </Button>
        </Box>
      </DialogTitle>

      <DialogContent sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        <Tabs value={tabValue} onChange={(_, val) => setTabValue(val)} sx={{ mb: 2 }}>
          <Tab label="Настройки" />
          <Tab label="Предпросмотр" disabled={!preview} />
        </Tabs>

        <TabPanel value={tabValue} index={0}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
            <FormControl fullWidth>
              <InputLabel>Тип шаблона</InputLabel>
              <Select
                value={selectedTemplateId}
                onChange={(e) => setSelectedTemplateId(e.target.value)}
                label="Тип шаблона"
              >
                {templates.map((template) => (
                  <MenuItem key={template.id} value={template.id}>
                    {template.name} - {TEMPLATE_TYPE_LABELS[template.template_type]}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            <Alert severity="info">
              <Typography variant="subtitle2" gutterBottom>
                Автоматически заполненные данные:
              </Typography>
              <ul style={{ margin: 0, paddingLeft: '20px' }}>
                <li>Название: {asset.name}</li>
                <li>Инвентарный номер: {asset.inventory_number}</li>
                {asset.serial_number && <li>Серийный номер: {asset.serial_number}</li>}
                {asset.pc_number && <li>Номер ПК: {asset.pc_number}</li>}
                {asset.manufacturer && <li>Производитель: {asset.manufacturer}</li>}
                {asset.model && <li>Модель: {asset.model}</li>}
                {asset.cpu && <li>Процессор: {asset.cpu}</li>}
                {asset.ram && <li>Оперативная память: {asset.ram}</li>}
                {asset.hdd_info && <li>Жесткий диск: {asset.hdd_info}</li>}
                {asset.network_card && <li>Сетевая карта: {asset.network_card}</li>}
                {asset.optical_drive && <li>Оптический привод: {asset.optical_drive}</li>}
                {asset.ip_address && <li>IP адрес: {asset.ip_address}</li>}
                {asset.mac_address && <li>MAC адрес: {asset.mac_address}</li>}
                {asset.purchase_year && <li>Год покупки: {asset.purchase_year}</li>}
                {asset.warranty_until && <li>Гарантия до: {new Date(asset.warranty_until).toLocaleDateString('ru-RU')}</li>}
              </ul>
            </Alert>

            <Paper sx={{ p: 2, bgcolor: 'grey.50' }}>
              <Typography variant="subtitle1" gutterBottom fontWeight="bold">
                Дополнительные данные (опционально)
              </Typography>
              <Typography variant="caption" color="text.secondary" display="block" mb={2}>
                Здесь вы можете ввести дополнительные данные, которых нет в системе
              </Typography>
              
              <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 2 }}>
                <TextField
                  label="Диагональ монитора"
                  value={additionalData.monitor_diagonal || ''}
                  onChange={(e) => setAdditionalData({ ...additionalData, monitor_diagonal: e.target.value })}
                  size="small"
                />
                <TextField
                  label="Тип монитора (CRT/LCD)"
                  value={additionalData.monitor_type || ''}
                  onChange={(e) => setAdditionalData({ ...additionalData, monitor_type: e.target.value })}
                  size="small"
                />
                <TextField
                  label="Операционная система"
                  value={additionalData.os || ''}
                  onChange={(e) => setAdditionalData({ ...additionalData, os: e.target.value })}
                  size="small"
                />
                <TextField
                  label="Офисный пакет"
                  value={additionalData.office_suite || ''}
                  onChange={(e) => setAdditionalData({ ...additionalData, office_suite: e.target.value })}
                  size="small"
                />
                <TextField
                  label="Антивирус"
                  value={additionalData.antivirus || ''}
                  onChange={(e) => setAdditionalData({ ...additionalData, antivirus: e.target.value })}
                  size="small"
                />
              </Box>
            </Paper>

            <Box display="flex" justifyContent="flex-end">
              <Button
                variant="contained"
                startIcon={<PreviewIcon />}
                onClick={handleGeneratePreview}
                disabled={!selectedTemplateId || loading}
              >
                {loading ? 'Генерация...' : 'Сгенерировать предпросмотр'}
              </Button>
            </Box>
          </Box>
        </TabPanel>

        <TabPanel value={tabValue} index={1}>
          {preview && (
            <Box sx={{ flex: 1, overflow: 'auto' }}>
              <Paper sx={{ p: 3 }}>
                <div dangerouslySetInnerHTML={{ __html: preview }} />
              </Paper>
            </Box>
          )}
        </TabPanel>
      </DialogContent>

      <DialogActions>
        <Button onClick={handleClose}>
          Отмена
        </Button>
        <Button
          startIcon={<DownloadIcon />}
          onClick={handleDownloadPDF}
          disabled={!preview}
        >
          Скачать PDF
        </Button>
        <Button
          variant="contained"
          startIcon={<SaveIcon />}
          onClick={handleSaveAsDocument}
          disabled={!preview}
        >
          Сохранить как документ
        </Button>
      </DialogActions>
    </Dialog>
  );
};

