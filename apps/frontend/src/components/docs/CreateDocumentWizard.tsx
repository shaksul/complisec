import React, { useState } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Stepper,
  Step,
  StepLabel,
  Box,
  Typography,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Autocomplete,
  Grid,
  FormControlLabel,
  Switch,
  Alert,
  CircularProgress,
} from '@mui/material'
import {
  CloudUpload,
  CheckCircle,
  Error as ErrorIcon,
} from '@mui/icons-material'
import { uploadDocument, submitDocumentForApproval } from '../../shared/api/documents'
// import { getUsers } from '../../shared/api/users'
import { type SubmitDocumentDTO } from '../../shared/api/documents'

interface CreateDocumentWizardProps {
  open: boolean
  onClose: () => void
  onSuccess: (documentId: string) => void
}

interface User {
  id: string
  first_name?: string
  last_name?: string
  name: string
  email: string
}

interface Step1Data {
  title: string
  code: string
  description: string
  type: string
  category: string
  tags: string[]
  ownerId: string
  classification: string
  effectiveFrom: string
  reviewPeriodMonths: number
  assetIds: string[]
  riskIds: string[]
  controlIds: string[]
}

interface Step2Data {
  file: File | null
  enableOCR: boolean
  uploadProgress: number
  uploadStatus: 'idle' | 'uploading' | 'success' | 'error'
  uploadError: string | null
}

interface Step3Data {
  requiresApproval: boolean
  workflowType: 'sequential' | 'parallel'
  steps: Array<{
    stepOrder: number
    approverId: string
    deadline: string
  }>
}

const DOCUMENT_TYPES = [
  { value: 'policy', label: 'Политика' },
  { value: 'standard', label: 'Стандарт' },
  { value: 'procedure', label: 'Процедура' },
  { value: 'instruction', label: 'Инструкция' },
  { value: 'act', label: 'Акт' },
  { value: 'other', label: 'Другое' },
]

const CLASSIFICATIONS = [
  { value: 'Public', label: 'Публичный' },
  { value: 'Internal', label: 'Внутренний' },
  { value: 'Confidential', label: 'Конфиденциальный' },
]

const COMMON_TAGS = [
  'Безопасность',
  'Информационная безопасность',
  'Персональные данные',
  'Соблюдение требований',
  'Управление рисками',
  'Обучение',
  'Процедуры',
  'Политики',
]

export default function CreateDocumentWizard({ open, onClose, onSuccess }: CreateDocumentWizardProps) {
  const [activeStep, setActiveStep] = useState(0)
  const [documentId, setDocumentId] = useState<string | null>(null)
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Step 1 data
  const [step1Data, setStep1Data] = useState<Step1Data>({
    title: '',
    code: '',
    description: '',
    type: 'policy',
    category: '',
    tags: [],
    ownerId: '',
    classification: 'Internal',
    effectiveFrom: '',
    reviewPeriodMonths: 12,
    assetIds: [],
    riskIds: [],
    controlIds: [],
  })

  // Step 2 data
  const [step2Data, setStep2Data] = useState<Step2Data>({
    file: null,
    enableOCR: true,
    uploadProgress: 0,
    uploadStatus: 'idle',
    uploadError: null,
  })

  // Step 3 data
  const [step3Data, setStep3Data] = useState<Step3Data>({
    requiresApproval: false,
    workflowType: 'sequential',
    steps: [{ stepOrder: 1, approverId: '', deadline: '' }],
  })

  React.useEffect(() => {
    if (open) {
      loadUsers()
    }
  }, [open])

  const loadUsers = async () => {
    try {
      // const userData = await getUsers()
      setUsers([])
    } catch (err) {
      console.error('Error loading users:', err)
      // Set mock users for now
      setUsers([
        { id: '1', first_name: 'Admin', last_name: 'User', name: 'Admin User', email: 'admin@demo.local' },
        { id: '2', first_name: 'John', last_name: 'Doe', name: 'John Doe', email: 'john@demo.local' },
        { id: '3', first_name: 'Jane', last_name: 'Smith', name: 'Jane Smith', email: 'jane@demo.local' }
      ])
    }
  }

  const handleStep1Next = () => {
    // Just validate and move to next step, don't create document yet
    if (!step1Data.title.trim()) {
      setError('Название документа обязательно')
      return
    }
    setError(null)
    setActiveStep(1)
  }

  const handleFileUpload = async (file: File) => {
    try {
      setStep2Data(prev => ({ ...prev, uploadStatus: 'uploading', uploadError: null }))
      setLoading(true)
      setError(null)
      
      // Simulate upload progress
      for (let i = 0; i <= 100; i += 10) {
        setStep2Data(prev => ({ ...prev, uploadProgress: i }))
        await new Promise(resolve => setTimeout(resolve, 100))
      }

      // Upload file directly using the main upload API
      const document = await uploadDocument(
        file,
        step1Data.title,
        step1Data.description || undefined,
        step1Data.tags,
        { module: 'documents', entity_id: 'general' }
      )
      setDocumentId(document.id)
      
      setStep2Data(prev => ({ 
        ...prev, 
        file, 
        uploadStatus: 'success',
        uploadProgress: 100 
      }))
    } catch (err) {
      setStep2Data(prev => ({ 
        ...prev, 
        uploadStatus: 'error',
        uploadError: (err as Error).message 
      }))
      setError('Ошибка загрузки файла: ' + (err as Error).message)
    } finally {
      setLoading(false)
    }
  }

  const handleStep2Next = () => {
    if (step2Data.uploadStatus === 'success') {
      setActiveStep(2)
    }
  }

  const handleAddApprovalStep = () => {
    setStep3Data(prev => ({
      ...prev,
      steps: [
        ...prev.steps,
        { 
          stepOrder: prev.steps.length + 1, 
          approverId: '', 
          deadline: '' 
        }
      ]
    }))
  }

  const handleRemoveApprovalStep = (index: number) => {
    setStep3Data(prev => ({
      ...prev,
      steps: prev.steps.filter((_, i) => i !== index)
    }))
  }

  const handleStep3Submit = async () => {
    if (!documentId) return

    try {
      setLoading(true)
      setError(null)

      if (step3Data.requiresApproval) {
        const submitData: SubmitDocumentDTO = {
          workflow_type: step3Data.workflowType,
          steps: step3Data.steps.map(step => ({
            step_order: step.stepOrder,
            approver_id: step.approverId,
            deadline: step.deadline || undefined,
          }))
        }

        await submitDocumentForApproval(documentId, submitData)
      }
      
      onSuccess(documentId)
      onClose()
    } catch (err) {
      setError('Ошибка отправки на согласование: ' + (err as Error).message)
    } finally {
      setLoading(false)
    }
  }

  const handleClose = () => {
    setActiveStep(0)
    setDocumentId(null)
    setError(null)
    setStep1Data({
      title: '',
      code: '',
      description: '',
      type: 'policy',
      category: '',
      tags: [],
      ownerId: '',
      classification: 'Internal',
      effectiveFrom: '',
      reviewPeriodMonths: 12,
      assetIds: [],
      riskIds: [],
      controlIds: [],
    })
    setStep2Data({
      file: null,
      enableOCR: true,
      uploadProgress: 0,
      uploadStatus: 'idle',
      uploadError: null,
    })
    setStep3Data({
      requiresApproval: false,
      workflowType: 'sequential',
      steps: [{ stepOrder: 1, approverId: '', deadline: '' }],
    })
    onClose()
  }

  const renderStep1 = () => (
    <Box>
      <Typography variant="h6" gutterBottom>
        Метаданные документа
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <TextField
            label="Название *"
            value={step1Data.title}
            onChange={(e) => setStep1Data(prev => ({ ...prev, title: e.target.value }))}
            fullWidth
            required
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <TextField
            label="Код документа"
            value={step1Data.code}
            onChange={(e) => setStep1Data(prev => ({ ...prev, code: e.target.value }))}
            fullWidth
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <FormControl fullWidth>
            <InputLabel>Тип документа *</InputLabel>
            <Select
              value={step1Data.type}
              onChange={(e) => setStep1Data(prev => ({ ...prev, type: e.target.value }))}
            >
              {DOCUMENT_TYPES.map(type => (
                <MenuItem key={type.value} value={type.value}>
                  {type.label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12}>
          <TextField
            label="Описание"
            value={step1Data.description}
            onChange={(e) => setStep1Data(prev => ({ ...prev, description: e.target.value }))}
            fullWidth
            multiline
            rows={3}
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <TextField
            label="Категория"
            value={step1Data.category}
            onChange={(e) => setStep1Data(prev => ({ ...prev, category: e.target.value }))}
            fullWidth
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <FormControl fullWidth>
            <InputLabel>Классификация *</InputLabel>
            <Select
              value={step1Data.classification}
              onChange={(e) => setStep1Data(prev => ({ ...prev, classification: e.target.value }))}
            >
              {CLASSIFICATIONS.map(classification => (
                <MenuItem key={classification.value} value={classification.value}>
                  {classification.label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12}>
          <Autocomplete
            multiple
            options={COMMON_TAGS}
            value={step1Data.tags}
            onChange={(_, newValue) => setStep1Data(prev => ({ ...prev, tags: newValue }))}
            freeSolo
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip variant="outlined" label={option} {...getTagProps({ index })} />
              ))
            }
            renderInput={(params) => (
              <TextField
                {...params}
                label="Теги"
                placeholder="Добавить тег"
              />
            )}
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <FormControl fullWidth>
            <InputLabel>Владелец документа</InputLabel>
            <Select
              value={step1Data.ownerId}
              onChange={(e) => setStep1Data(prev => ({ ...prev, ownerId: e.target.value }))}
            >
              {users.map(user => (
                <MenuItem key={user.id} value={user.id}>
                  {user.name} ({user.email})
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <TextField
            label="Дата вступления в силу"
            type="date"
            value={step1Data.effectiveFrom}
            onChange={(e) => setStep1Data(prev => ({ ...prev, effectiveFrom: e.target.value }))}
            fullWidth
            InputLabelProps={{ shrink: true }}
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <TextField
            label="Период пересмотра (месяцы)"
            type="number"
            value={step1Data.reviewPeriodMonths}
            onChange={(e) => setStep1Data(prev => ({ ...prev, reviewPeriodMonths: parseInt(e.target.value) || 12 }))}
            fullWidth
            inputProps={{ min: 1, max: 120 }}
          />
        </Grid>
      </Grid>
    </Box>
  )

  const renderStep2 = () => (
    <Box>
      <Typography variant="h6" gutterBottom>
        Загрузка версии документа
      </Typography>
      
      <Box
        sx={{
          border: (theme) => `2px dashed ${theme.palette.divider}`,
          borderRadius: 2,
          p: 4,
          textAlign: 'center',
          mb: 3,
        }}
      >
        {step2Data.uploadStatus === 'idle' && (
          <Box>
            <CloudUpload sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" gutterBottom>
              Перетащите файл сюда или нажмите для выбора
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              Поддерживаемые форматы: PDF, DOCX, TXT
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Максимальный размер: 50 МБ
            </Typography>
            <input
              type="file"
              accept=".pdf,.docx"
              onChange={(e) => {
                const file = e.target.files?.[0]
                if (file) handleFileUpload(file)
              }}
              style={{ display: 'none' }}
              id="file-upload"
            />
            <label htmlFor="file-upload">
              <Button variant="contained" component="span" sx={{ mt: 2 }}>
                Выбрать файл
              </Button>
            </label>
          </Box>
        )}
        
        {step2Data.uploadStatus === 'uploading' && (
          <Box>
            <CircularProgress sx={{ mb: 2 }} />
            <Typography variant="h6" gutterBottom>
              Загрузка файла...
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {step2Data.uploadProgress}%
            </Typography>
          </Box>
        )}
        
        {step2Data.uploadStatus === 'success' && (
          <Box>
            <CheckCircle sx={{ fontSize: 48, color: 'success.main', mb: 2 }} />
            <Typography variant="h6" gutterBottom>
              Файл успешно загружен
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {step2Data.file?.name}
            </Typography>
          </Box>
        )}
        
        {step2Data.uploadStatus === 'error' && (
          <Box>
            <ErrorIcon sx={{ fontSize: 48, color: 'error.main', mb: 2 }} />
            <Typography variant="h6" gutterBottom color="error">
              Ошибка загрузки
            </Typography>
            <Typography variant="body2" color="error">
              {step2Data.uploadError}
            </Typography>
          </Box>
        )}
      </Box>
      
      <FormControlLabel
        control={
          <Switch
            checked={step2Data.enableOCR}
            onChange={(e) => setStep2Data(prev => ({ ...prev, enableOCR: e.target.checked }))}
          />
        }
        label="Включить OCR для извлечения текста"
      />
    </Box>
  )

  const renderStep3 = () => (
    <Box>
      <Typography variant="h6" gutterBottom>
        Настройки согласования
      </Typography>
      
      <FormControlLabel
        control={
          <Switch
            checked={step3Data.requiresApproval}
            onChange={(e) => setStep3Data(prev => ({ ...prev, requiresApproval: e.target.checked }))}
          />
        }
        label="Требуется согласование документа"
        sx={{ mb: 3 }}
      />
      
      {step3Data.requiresApproval && (
        <>
          <FormControl fullWidth sx={{ mb: 3 }}>
            <InputLabel>Тип маршрута</InputLabel>
            <Select
              value={step3Data.workflowType}
              onChange={(e) => setStep3Data(prev => ({ 
                ...prev, 
                workflowType: e.target.value as 'sequential' | 'parallel' 
              }))}
            >
              <MenuItem value="sequential">Последовательный</MenuItem>
              <MenuItem value="parallel">Параллельный</MenuItem>
            </Select>
          </FormControl>
          
          <Typography variant="subtitle1" gutterBottom>
            Шаги согласования:
          </Typography>
        </>
      )}
      
      {step3Data.steps.map((step, index) => (
        <Box key={index} sx={{ mb: 2, p: 2, border: 1, borderColor: 'divider', borderRadius: 1 }}>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={12} sm={3}>
              <TextField
                label="Порядок"
                type="number"
                value={step.stepOrder}
                onChange={(e) => {
                  const newSteps = [...step3Data.steps]
                  newSteps[index].stepOrder = parseInt(e.target.value) || 1
                  setStep3Data(prev => ({ ...prev, steps: newSteps }))
                }}
                fullWidth
                inputProps={{ min: 1 }}
              />
            </Grid>
            
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Согласующий</InputLabel>
                <Select
                  value={step.approverId}
                  onChange={(e) => {
                    const newSteps = [...step3Data.steps]
                    newSteps[index].approverId = e.target.value
                    setStep3Data(prev => ({ ...prev, steps: newSteps }))
                  }}
                >
                  {users.map(user => (
                    <MenuItem key={user.id} value={user.id}>
                      {user.name} ({user.email})
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            
            <Grid item xs={12} sm={2}>
              <TextField
                label="Дедлайн"
                type="date"
                value={step.deadline}
                onChange={(e) => {
                  const newSteps = [...step3Data.steps]
                  newSteps[index].deadline = e.target.value
                  setStep3Data(prev => ({ ...prev, steps: newSteps }))
                }}
                fullWidth
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            
            <Grid item xs={12} sm={1}>
              <Button
                color="error"
                onClick={() => handleRemoveApprovalStep(index)}
                disabled={step3Data.steps.length === 1}
              >
                Удалить
              </Button>
            </Grid>
          </Grid>
        </Box>
      ))}
      
      <Button
        variant="outlined"
        onClick={handleAddApprovalStep}
        sx={{ mb: 2 }}
      >
        Добавить шаг согласования
      </Button>
    </Box>
  )

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>Создание документа</DialogTitle>
      
      <DialogContent>
        <Stepper activeStep={activeStep} sx={{ mb: 4 }}>
          <Step>
            <StepLabel>Метаданные</StepLabel>
          </Step>
          <Step>
            <StepLabel>Загрузка версии</StepLabel>
          </Step>
          <Step>
            <StepLabel>Согласование</StepLabel>
          </Step>
        </Stepper>
        
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        
        {activeStep === 0 && renderStep1()}
        {activeStep === 1 && renderStep2()}
        {activeStep === 2 && renderStep3()}
      </DialogContent>
      
      <DialogActions>
        <Button onClick={handleClose}>
          Отмена
        </Button>
        
        {activeStep === 0 && (
          <Button
            variant="contained"
            onClick={handleStep1Next}
            disabled={loading || !step1Data.title || !step1Data.type}
          >
            {loading ? <CircularProgress size={20} /> : 'Далее'}
          </Button>
        )}
        
        {activeStep === 1 && (
          <>
            <Button onClick={() => setActiveStep(0)}>
              Назад
            </Button>
            <Button
              variant="contained"
              onClick={handleStep2Next}
              disabled={step2Data.uploadStatus !== 'success'}
            >
              Далее
            </Button>
          </>
        )}
        
        {activeStep === 2 && (
          <>
            <Button onClick={() => setActiveStep(1)}>
              Назад
            </Button>
            <Button
              variant="contained"
              onClick={handleStep3Submit}
              disabled={
                loading || (step3Data.requiresApproval && step3Data.steps.some(step => !step.approverId))
              }
            >
              {loading ? (
                <CircularProgress size={20} />
              ) : step3Data.requiresApproval ? (
                'Отправить на согласование'
              ) : (
                'Сохранить'
              )}
            </Button>
          </>
        )}
      </DialogActions>
    </Dialog>
  )
}
