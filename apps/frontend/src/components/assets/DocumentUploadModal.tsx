import React, { useState, useRef, useCallback } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Box,
  Typography,
  Alert,
  CircularProgress,
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  IconButton,
  LinearProgress,
  Radio,
  InputAdornment,
  Tooltip
} from '@mui/material';
import {
  Close as CloseIcon,
  CloudUpload as UploadIcon,
  Search as SearchIcon,
  AttachFile as AttachIcon,
  Delete as DeleteIcon,
  Download as DownloadIcon
} from '@mui/icons-material';
import { api } from '../../shared/api/client';

interface DocumentUploadModalProps {
  open: boolean;
  onClose: () => void;
  assetId: string;
  onSuccess: () => void;
}

const DOCUMENT_TYPES = [
  { value: 'passport', label: 'Паспорт' },
  { value: 'transfer_act', label: 'Акт передачи' },
  { value: 'writeoff_act', label: 'Акт списания' },
  { value: 'repair_log', label: 'Журнал ремонтов' },
  { value: 'other', label: 'Другое' }
];

const ALLOWED_FILE_TYPES = [
  'application/pdf',
  'image/jpeg',
  'image/png',
  'application/msword',
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  'application/vnd.ms-excel',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
];

const MAX_FILE_SIZE = 50 * 1024 * 1024; // 50MB

interface DocumentFile {
  file: File;
  preview: string;
  size: number;
  type: string;
}

interface StorageDocument {
  id: string;
  title: string;
  document_type: string;
  version: string;
  size_bytes: number;
  mime: string;
  created_by: string;
  created_at: string;
}

export const DocumentUploadModal: React.FC<DocumentUploadModalProps> = ({
  open,
  onClose,
  assetId,
  onSuccess
}) => {
  const [tabValue, setTabValue] = useState(0);
  const [formData, setFormData] = useState({
    document_type: 'passport',
    title: ''
  });
  const [selectedFile, setSelectedFile] = useState<DocumentFile | null>(null);
  const [dragActive, setDragActive] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  
  // Storage tab state
  const [storageDocuments, setStorageDocuments] = useState<StorageDocument[]>([]);
  const [selectedDocumentId, setSelectedDocumentId] = useState<string>('');
  const [storageLoading, setStorageLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterType, setFilterType] = useState('');
  const [pagination, setPagination] = useState({
    page: 1,
    pageSize: 25,
    total: 0,
    totalPages: 0
  });

  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = useCallback((file: File) => {
    // Validate file type by MIME type OR file extension
    const fileExtension = file.name.split('.').pop()?.toLowerCase() || '';
    const allowedExtensions = ['pdf', 'jpg', 'jpeg', 'png', 'doc', 'docx', 'xls', 'xlsx'];
    
    const isValidType = ALLOWED_FILE_TYPES.includes(file.type) || allowedExtensions.includes(fileExtension);
    
    if (!isValidType) {
      setError('Неподдерживаемый тип файла. Разрешены: PDF, JPG, PNG, DOC, DOCX, XLS, XLSX');
      return;
    }

    // Validate file size
    if (file.size > MAX_FILE_SIZE) {
      setError('Размер файла превышает 50MB');
      return;
    }

    setError(null);
    setSelectedFile({
      file,
      preview: URL.createObjectURL(file),
      size: file.size,
      type: file.type
    });

    // Auto-fill title if empty
    if (!formData.title) {
      setFormData(prev => ({
        ...prev,
        title: file.name.replace(/\.[^/.]+$/, '') // Remove extension
      }));
    }
  }, [formData.title]);

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFileSelect(e.dataTransfer.files[0]);
    }
  }, [handleFileSelect]);

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      handleFileSelect(e.target.files[0]);
    }
  };

  const removeFile = () => {
    if (selectedFile) {
      URL.revokeObjectURL(selectedFile.preview);
    }
    setSelectedFile(null);
    setError(null);
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setError('Выберите файл для загрузки');
      return;
    }

    setLoading(true);
    setError(null);
    setUploadProgress(0);

    try {
      const uploadFormData = new FormData();
      uploadFormData.append('document_type', formData.document_type);
      uploadFormData.append('title', formData.title);
      uploadFormData.append('file', selectedFile.file);

      // Simulate upload progress
      const progressInterval = setInterval(() => {
        setUploadProgress(prev => {
          if (prev >= 90) {
            clearInterval(progressInterval);
            return prev;
          }
          return prev + 10;
        });
      }, 200);

      await api.post(`/assets/${assetId}/documents/upload`, uploadFormData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        onUploadProgress: (progressEvent) => {
          if (progressEvent.total) {
            const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
            setUploadProgress(progress);
          }
        }
      });

      clearInterval(progressInterval);
      setUploadProgress(100);

      onSuccess();
      onClose();
      setFormData({ document_type: 'passport', title: '' });
      removeFile();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка при загрузке документа');
    } finally {
      setLoading(false);
      setUploadProgress(0);
    }
  };

  const loadStorageDocuments = async () => {
    setStorageLoading(true);
    try {
      const params = new URLSearchParams({
        page: pagination.page.toString(),
        page_size: pagination.pageSize.toString(),
        ...(searchQuery && { query: searchQuery }),
        ...(filterType && { type: filterType })
      });

      const response = await api.get(`/assets/documents/storage?${params}`);
      setStorageDocuments(response.data.data || []);
      setPagination(prev => ({
        ...prev,
        total: response.data.pagination.total,
        totalPages: response.data.pagination.total_pages
      }));
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка загрузки документов из хранилища');
    } finally {
      setStorageLoading(false);
    }
  };

  const handleLinkDocument = async () => {
    if (!selectedDocumentId) {
      setError('Выберите документ из хранилища');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const selectedDoc = storageDocuments.find(doc => doc.id === selectedDocumentId);
      if (!selectedDoc) {
        throw new Error('Выбранный документ не найден');
      }

      await api.post(`/assets/${assetId}/documents/link`, {
        document_id: selectedDocumentId,
        document_type: selectedDoc.document_type
      });

      onSuccess();
      onClose();
      setSelectedDocumentId('');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка при привязке документа');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    setPagination(prev => ({ ...prev, page: 1 }));
    loadStorageDocuments();
  };

  React.useEffect(() => {
    if (open && tabValue === 1) {
      loadStorageDocuments();
    }
  }, [open, tabValue]);

  const handleClose = () => {
    removeFile();
    setError(null);
    setSelectedDocumentId('');
    setSearchQuery('');
    setFilterType('');
    setTabValue(0);
    onClose();
  };

  return (
    <Dialog 
      open={open} 
      onClose={handleClose}
      maxWidth="md"
      fullWidth
      aria-labelledby="document-upload-dialog"
    >
      <DialogTitle id="document-upload-dialog">
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h6">Добавить документ к активу</Typography>
          <IconButton onClick={handleClose} size="small">
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>

      <DialogContent>
        <Tabs 
          value={tabValue} 
          onChange={(_, newValue) => setTabValue(newValue)}
          sx={{ mb: 2 }}
        >
          <Tab label="Загрузить новый" />
          <Tab label="Из хранилища" />
        </Tabs>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {tabValue === 0 && (
          <Box>
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Тип документа *</InputLabel>
              <Select
                value={formData.document_type}
                onChange={(e) => setFormData(prev => ({ ...prev, document_type: e.target.value }))}
                label="Тип документа *"
              >
                {DOCUMENT_TYPES.map((type) => (
                  <MenuItem key={type.value} value={type.value}>
                    {type.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            <TextField
              fullWidth
              label="Название документа"
              value={formData.title}
              onChange={(e) => setFormData(prev => ({ ...prev, title: e.target.value }))}
              sx={{ mb: 2 }}
              placeholder="Введите название документа (необязательно)"
            />

            {/* File Upload Area */}
            <Box
              sx={{
                border: '2px dashed',
                borderColor: dragActive ? 'primary.main' : 'grey.300',
                borderRadius: 2,
                p: 3,
                textAlign: 'center',
                cursor: 'pointer',
                transition: 'all 0.2s',
                bgcolor: dragActive ? 'action.hover' : 'background.paper',
                '&:hover': {
                  borderColor: 'primary.main',
                  bgcolor: 'action.hover'
                }
              }}
              onDragEnter={handleDrag}
              onDragLeave={handleDrag}
              onDragOver={handleDrag}
              onDrop={handleDrop}
              onClick={() => fileInputRef.current?.click()}
            >
              <input
                ref={fileInputRef}
                type="file"
                hidden
                accept=".pdf,.jpg,.jpeg,.png,.doc,.docx,.xls,.xlsx"
                onChange={handleFileInputChange}
              />
              
              {selectedFile ? (
                <Box>
                  <AttachIcon sx={{ fontSize: 48, color: 'primary.main', mb: 1 }} />
                  <Typography variant="h6" gutterBottom>
                    {selectedFile.file.name}
                  </Typography>
                  <Box display="flex" justifyContent="center" gap={2} mb={2}>
                    <Chip 
                      label={formatFileSize(selectedFile.size)} 
                      size="small" 
                      color="primary" 
                      variant="outlined" 
                    />
                    <Chip 
                      label={selectedFile.type} 
                      size="small" 
                      color="secondary" 
                      variant="outlined" 
                    />
                  </Box>
                  <Button
                    variant="outlined"
                    color="error"
                    startIcon={<DeleteIcon />}
                    onClick={(e) => {
                      e.stopPropagation();
                      removeFile();
                    }}
                  >
                    Удалить файл
                  </Button>
                </Box>
              ) : (
                <Box>
                  <UploadIcon sx={{ fontSize: 48, color: 'grey.400', mb: 1 }} />
                  <Typography variant="h6" gutterBottom>
                    Перетащите файл сюда или нажмите для выбора
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Поддерживаемые форматы: PDF, JPG, PNG, DOC, DOCX, XLS, XLSX
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    Максимальный размер: 50MB
                  </Typography>
                </Box>
              )}
            </Box>

            {loading && (
              <Box sx={{ mt: 2 }}>
                <Typography variant="body2" gutterBottom>
                  Загрузка файла... {uploadProgress}%
                </Typography>
                <LinearProgress variant="determinate" value={uploadProgress} />
              </Box>
            )}
          </Box>
        )}

        {tabValue === 1 && (
          <Box>
            <Box display="flex" gap={2} mb={2}>
              <TextField
                fullWidth
                label="Поиск документов"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <SearchIcon />
                    </InputAdornment>
                  ),
                }}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
              />
              <FormControl sx={{ minWidth: 150 }}>
                <InputLabel>Тип</InputLabel>
                <Select
                  value={filterType}
                  onChange={(e) => setFilterType(e.target.value)}
                  label="Тип"
                >
                  <MenuItem value="">Все типы</MenuItem>
                  {DOCUMENT_TYPES.map((type) => (
                    <MenuItem key={type.value} value={type.value}>
                      {type.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
              <Button variant="contained" onClick={handleSearch} disabled={storageLoading}>
                Поиск
              </Button>
            </Box>

            <TableContainer component={Paper} sx={{ maxHeight: 400 }}>
              <Table stickyHeader>
                <TableHead>
                  <TableRow>
                    <TableCell padding="checkbox"></TableCell>
                    <TableCell>Название</TableCell>
                    <TableCell>Тип</TableCell>
                    <TableCell>Размер</TableCell>
                    <TableCell>Дата загрузки</TableCell>
                    <TableCell>Кем добавлен</TableCell>
                    <TableCell>Действия</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {storageLoading ? (
                    <TableRow>
                      <TableCell colSpan={7} align="center">
                        <CircularProgress size={24} />
                      </TableCell>
                    </TableRow>
                  ) : storageDocuments.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} align="center">
                        <Typography color="text.secondary">
                          Документы не найдены
                        </Typography>
                      </TableCell>
                    </TableRow>
                  ) : (
                    storageDocuments.map((doc) => (
                      <TableRow key={doc.id} hover>
                        <TableCell padding="checkbox">
                          <Radio
                            checked={selectedDocumentId === doc.id}
                            onChange={() => setSelectedDocumentId(doc.id)}
                          />
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2" fontWeight="medium">
                            {doc.title}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={DOCUMENT_TYPES.find(t => t.value === doc.document_type)?.label || doc.document_type}
                            size="small"
                            color="primary"
                            variant="outlined"
                          />
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2">
                            {formatFileSize(doc.size_bytes)}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2">
                            {new Date(doc.created_at).toLocaleDateString('ru-RU')}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2">
                            {doc.created_by}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Tooltip title="Скачать">
                            <IconButton size="small">
                              <DownloadIcon />
                            </IconButton>
                          </Tooltip>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </TableContainer>
          </Box>
        )}
      </DialogContent>

      <DialogActions>
        <Button onClick={handleClose} disabled={loading}>
          Отмена
        </Button>
        {tabValue === 0 ? (
          <Button
            variant="contained"
            onClick={handleUpload}
            disabled={loading || !selectedFile}
            startIcon={loading ? <CircularProgress size={20} /> : <UploadIcon />}
          >
            {loading ? 'Загрузка...' : 'Загрузить'}
          </Button>
        ) : (
          <Button
            variant="contained"
            onClick={handleLinkDocument}
            disabled={loading || !selectedDocumentId}
            startIcon={loading ? <CircularProgress size={20} /> : <AttachIcon />}
          >
            {loading ? 'Привязка...' : 'Привязать'}
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
};
