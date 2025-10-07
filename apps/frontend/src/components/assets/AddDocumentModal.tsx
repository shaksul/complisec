import React, { useState } from 'react';
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
  CircularProgress
} from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { assetsApi } from '../../shared/api/assets';

interface AddDocumentModalProps {
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

export const AddDocumentModal: React.FC<AddDocumentModalProps> = ({
  open,
  onClose,
  assetId,
  onSuccess
}) => {
  const [formData, setFormData] = useState({
    document_type: 'passport',
    file_path: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await assetsApi.addDocument(assetId, formData);
      onSuccess();
      onClose();
      setFormData({ document_type: 'passport', file_path: '' });
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка при добавлении документа');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    if (!loading) {
      onClose();
      setFormData({ document_type: 'passport', file_path: '' });
      setError(null);
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Box display="flex" alignItems="center" gap={1}>
          <AddIcon />
          <Typography variant="h6">Добавить документ</Typography>
        </Box>
      </DialogTitle>
      
      <form onSubmit={handleSubmit}>
        <DialogContent>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          
          <FormControl fullWidth margin="normal" required>
            <InputLabel>Тип документа</InputLabel>
            <Select
              value={formData.document_type}
              onChange={(e) => setFormData({ ...formData, document_type: e.target.value })}
              label="Тип документа"
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
            margin="normal"
            label="Путь к файлу"
            value={formData.file_path}
            onChange={(e) => setFormData({ ...formData, file_path: e.target.value })}
            required
            placeholder="Введите путь к файлу документа"
            helperText="Укажите полный путь к файлу документа"
          />
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose} disabled={loading}>
            Отмена
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={loading}
            startIcon={loading ? <CircularProgress size={20} /> : <AddIcon />}
          >
            {loading ? 'Добавление...' : 'Добавить'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};
