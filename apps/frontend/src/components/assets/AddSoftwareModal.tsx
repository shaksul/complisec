import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Box,
  Typography,
  Alert,
  CircularProgress
} from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { assetsApi } from '../../shared/api/assets';

interface AddSoftwareModalProps {
  open: boolean;
  onClose: () => void;
  assetId: string;
  onSuccess: () => void;
}

export const AddSoftwareModal: React.FC<AddSoftwareModalProps> = ({
  open,
  onClose,
  assetId,
  onSuccess
}) => {
  const [formData, setFormData] = useState({
    software_name: '',
    version: '',
    installed_at: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const data = {
        software_name: formData.software_name,
        version: formData.version || undefined,
        installed_at: formData.installed_at ? new Date(formData.installed_at).toISOString() : undefined
      };
      
      await assetsApi.addSoftware(assetId, data);
      onSuccess();
      onClose();
      setFormData({ software_name: '', version: '', installed_at: '' });
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка при добавлении ПО');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    if (!loading) {
      onClose();
      setFormData({ software_name: '', version: '', installed_at: '' });
      setError(null);
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Box display="flex" alignItems="center" gap={1}>
          <AddIcon />
          <Typography variant="h6">Добавить ПО</Typography>
        </Box>
      </DialogTitle>
      
      <form onSubmit={handleSubmit}>
        <DialogContent>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          
          <TextField
            fullWidth
            margin="normal"
            label="Название ПО"
            value={formData.software_name}
            onChange={(e) => setFormData({ ...formData, software_name: e.target.value })}
            required
            placeholder="Введите название программного обеспечения"
          />

          <TextField
            fullWidth
            margin="normal"
            label="Версия"
            value={formData.version}
            onChange={(e) => setFormData({ ...formData, version: e.target.value })}
            placeholder="Введите версию ПО (необязательно)"
          />

          <TextField
            fullWidth
            margin="normal"
            label="Дата установки"
            type="datetime-local"
            value={formData.installed_at}
            onChange={(e) => setFormData({ ...formData, installed_at: e.target.value })}
            InputLabelProps={{
              shrink: true,
            }}
            helperText="Укажите дату и время установки ПО (необязательно)"
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
