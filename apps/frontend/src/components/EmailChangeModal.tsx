import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Typography,
  Box,
  Alert,
  CircularProgress,
  Stepper,
  Step,
  StepLabel,
  StepContent,
  Chip,
} from '@mui/material';
import { emailChangeApi, EmailChangeStatusResponse } from '../shared/api/emailChange';

interface EmailChangeModalProps {
  open: boolean;
  onClose: () => void;
  currentEmail: string;
}

const EmailChangeModal: React.FC<EmailChangeModalProps> = ({
  open,
  onClose,
  currentEmail,
}) => {
  const [activeStep, setActiveStep] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  
  // Форма для запроса смены email
  const [newEmail, setNewEmail] = useState('');
  const [requestId, setRequestId] = useState<string | null>(null);
  
  // Форма для подтверждения
  const [verificationCode, setVerificationCode] = useState('');
  
  // Статус запроса
  const [, setStatus] = useState<EmailChangeStatusResponse | null>(null);

  const steps = [
    'Запрос на смену email',
    'Подтверждение старого email',
    'Подтверждение нового email',
    'Завершение смены email',
  ];

  // Загружаем статус при открытии модального окна
  useEffect(() => {
    if (open) {
      loadStatus();
    }
  }, [open]);

  const loadStatus = async () => {
    try {
      const statusData = await emailChangeApi.getEmailChangeStatus();
      setStatus(statusData);
      
      if (statusData.has_active_request && statusData.request) {
        setRequestId(statusData.request.id);
        setNewEmail(statusData.request.new_email);
        
        // Определяем текущий шаг на основе статуса
        switch (statusData.request.status) {
          case 'pending':
            setActiveStep(1);
            break;
          case 'old_email_verified':
            setActiveStep(2);
            break;
          case 'new_email_verified':
            setActiveStep(3);
            break;
          default:
            setActiveStep(0);
        }
      } else {
        setActiveStep(0);
        setRequestId(null);
        setNewEmail('');
      }
    } catch (error) {
      console.error('Ошибка загрузки статуса:', error);
    }
  };

  const handleRequestEmailChange = async () => {
    if (!newEmail || newEmail === currentEmail) {
      setError('Новый email должен отличаться от текущего');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await emailChangeApi.requestEmailChange({
        new_email: newEmail,
      });
      
      setRequestId(response.request_id);
      setActiveStep(1);
      setSuccess('Запрос создан. Проверьте текущий email для получения кода подтверждения.');
    } catch (error: any) {
      setError(error.response?.data?.error || 'Ошибка создания запроса');
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyOldEmail = async () => {
    if (!verificationCode || !requestId) {
      setError('Введите код подтверждения');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      await emailChangeApi.verifyOldEmail({
        request_id: requestId,
        verification_code: verificationCode,
      });
      
      setActiveStep(2);
      setVerificationCode('');
      setSuccess('Старый email подтвержден. Проверьте новый email для получения кода подтверждения.');
    } catch (error: any) {
      setError(error.response?.data?.error || 'Ошибка подтверждения старого email');
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyNewEmail = async () => {
    if (!verificationCode || !requestId) {
      setError('Введите код подтверждения');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      await emailChangeApi.verifyNewEmail({
        request_id: requestId,
        verification_code: verificationCode,
      });
      
      setActiveStep(3);
      setVerificationCode('');
      setSuccess('Новый email подтвержден. Теперь можно завершить смену email.');
    } catch (error: any) {
      setError(error.response?.data?.error || 'Ошибка подтверждения нового email');
    } finally {
      setLoading(false);
    }
  };

  const handleCompleteEmailChange = async () => {
    if (!requestId) {
      setError('Отсутствует ID запроса');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      await emailChangeApi.completeEmailChange({
        request_id: requestId,
      });
      
      setSuccess('Email успешно изменен!');
      setTimeout(() => {
        onClose();
        window.location.reload(); // Перезагружаем страницу для обновления данных
      }, 2000);
    } catch (error: any) {
      setError(error.response?.data?.error || 'Ошибка завершения смены email');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = async () => {
    if (!requestId) {
      onClose();
      return;
    }

    setLoading(true);
    try {
      await emailChangeApi.cancelEmailChange({
        request_id: requestId,
      });
      onClose();
    } catch (error: any) {
      setError(error.response?.data?.error || 'Ошибка отмены запроса');
    } finally {
      setLoading(false);
    }
  };

  const handleResendCode = async () => {
    if (!requestId) {
      setError('Отсутствует ID запроса');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      await emailChangeApi.resendVerificationCode({
        request_id: requestId,
      });
      setSuccess('Код подтверждения отправлен повторно');
    } catch (error: any) {
      setError(error.response?.data?.error || 'Ошибка повторной отправки кода');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setActiveStep(0);
    setError(null);
    setSuccess(null);
    setNewEmail('');
    setVerificationCode('');
    setRequestId(null);
    setStatus(null);
    onClose();
  };

  const renderStepContent = (step: number) => {
    switch (step) {
      case 0:
        return (
          <Box sx={{ mt: 2 }}>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Текущий email: <strong>{currentEmail}</strong>
            </Typography>
            <TextField
              fullWidth
              label="Новый email"
              type="email"
              value={newEmail}
              onChange={(e) => setNewEmail(e.target.value)}
              disabled={loading}
              sx={{ mb: 2 }}
            />
            <Button
              variant="contained"
              onClick={handleRequestEmailChange}
              disabled={loading || !newEmail}
              startIcon={loading && <CircularProgress size={20} />}
            >
              Создать запрос
            </Button>
          </Box>
        );

      case 1:
        return (
          <Box sx={{ mt: 2 }}>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Код подтверждения отправлен на: <strong>{currentEmail}</strong>
            </Typography>
            <TextField
              fullWidth
              label="Код подтверждения"
              value={verificationCode}
              onChange={(e) => setVerificationCode(e.target.value)}
              disabled={loading}
              sx={{ mb: 2 }}
              inputProps={{ maxLength: 6 }}
            />
            <Box sx={{ display: 'flex', gap: 1 }}>
              <Button
                variant="contained"
                onClick={handleVerifyOldEmail}
                disabled={loading || !verificationCode}
                startIcon={loading && <CircularProgress size={20} />}
              >
                Подтвердить
              </Button>
              <Button
                variant="outlined"
                onClick={handleResendCode}
                disabled={loading}
              >
                Отправить повторно
              </Button>
            </Box>
          </Box>
        );

      case 2:
        return (
          <Box sx={{ mt: 2 }}>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Код подтверждения отправлен на: <strong>{newEmail}</strong>
            </Typography>
            <TextField
              fullWidth
              label="Код подтверждения"
              value={verificationCode}
              onChange={(e) => setVerificationCode(e.target.value)}
              disabled={loading}
              sx={{ mb: 2 }}
              inputProps={{ maxLength: 6 }}
            />
            <Box sx={{ display: 'flex', gap: 1 }}>
              <Button
                variant="contained"
                onClick={handleVerifyNewEmail}
                disabled={loading || !verificationCode}
                startIcon={loading && <CircularProgress size={20} />}
              >
                Подтвердить
              </Button>
              <Button
                variant="outlined"
                onClick={handleResendCode}
                disabled={loading}
              >
                Отправить повторно
              </Button>
            </Box>
          </Box>
        );

      case 3:
        return (
          <Box sx={{ mt: 2 }}>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Оба email подтверждены. Нажмите кнопку ниже для завершения смены email.
            </Typography>
            <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
              <Chip label={currentEmail} color="default" />
              <Typography variant="body2" sx={{ alignSelf: 'center' }}>→</Typography>
              <Chip label={newEmail} color="primary" />
            </Box>
            <Button
              variant="contained"
              color="success"
              onClick={handleCompleteEmailChange}
              disabled={loading}
              startIcon={loading && <CircularProgress size={20} />}
            >
              Завершить смену email
            </Button>
          </Box>
        );

      default:
        return null;
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>Смена email адреса</DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}
        
        {success && (
          <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccess(null)}>
            {success}
          </Alert>
        )}

        <Stepper activeStep={activeStep} orientation="vertical">
          {steps.map((label, index) => (
            <Step key={label}>
              <StepLabel>{label}</StepLabel>
              <StepContent>
                {renderStepContent(index)}
              </StepContent>
            </Step>
          ))}
        </Stepper>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} disabled={loading}>
          Закрыть
        </Button>
        {requestId && (
          <Button onClick={handleCancel} color="error" disabled={loading}>
            Отменить запрос
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
};

export default EmailChangeModal;
