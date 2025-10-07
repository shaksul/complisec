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
import { createDocument, uploadDocumentVersion, submitDocumentForApproval } from '../../shared/api/documents'
// import { getUsers } from '../../shared/api/users'
import { type CreateDocumentDTO, type SubmitDocumentDTO } from '../../shared/api/documents'

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
  { value: 'policy', label: 'РџРѕР»РёС‚РёРєР°' },
  { value: 'standard', label: 'РЎС‚Р°РЅРґР°СЂС‚' },
  { value: 'procedure', label: 'РџСЂРѕС†РµРґСѓСЂР°' },
  { value: 'instruction', label: 'РРЅСЃС‚СЂСѓРєС†РёСЏ' },
  { value: 'act', label: 'РђРєС‚' },
  { value: 'other', label: 'Р”СЂСѓРіРѕРµ' },
]

const CLASSIFICATIONS = [
  { value: 'Public', label: 'РџСѓР±Р»РёС‡РЅС‹Р№' },
  { value: 'Internal', label: 'Р’РЅСѓС‚СЂРµРЅРЅРёР№' },
  { value: 'Confidential', label: 'РљРѕРЅС„РёРґРµРЅС†РёР°Р»СЊРЅС‹Р№' },
]

const COMMON_TAGS = [
  'Р‘РµР·РѕРїР°СЃРЅРѕСЃС‚СЊ',
  'РРЅС„РѕСЂРјР°С†РёРѕРЅРЅР°СЏ Р±РµР·РѕРїР°СЃРЅРѕСЃС‚СЊ',
  'РџРµСЂСЃРѕРЅР°Р»СЊРЅС‹Рµ РґР°РЅРЅС‹Рµ',
  'РЎРѕР±Р»СЋРґРµРЅРёРµ С‚СЂРµР±РѕРІР°РЅРёР№',
  'РЈРїСЂР°РІР»РµРЅРёРµ СЂРёСЃРєР°РјРё',
  'РћР±СѓС‡РµРЅРёРµ',
  'РџСЂРѕС†РµРґСѓСЂС‹',
  'РџРѕР»РёС‚РёРєРё',
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
      setError('РќР°Р·РІР°РЅРёРµ РґРѕРєСѓРјРµРЅС‚Р° РѕР±СЏР·Р°С‚РµР»СЊРЅРѕ')
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
      
      // First create the document
      const createData: CreateDocumentDTO = {
        title: step1Data.title,
        description: step1Data.description || undefined,
        type: step1Data.type as 'policy' | 'standard' | 'procedure' | 'instruction' | 'act' | 'other',
        category: step1Data.category || undefined,
        tags: step1Data.tags,
        // owner_id: step1Data.ownerId || undefined,
        // classification: step1Data.classification,
        // effectiveFrom: step1Data.effectiveFrom || undefined,
        // reviewPeriodMonths: step1Data.reviewPeriodMonths,
        // assetIds: step1Data.assetIds,
        // riskIds: step1Data.riskIds,
        // controlIds: step1Data.controlIds,
      }

      const document = await createDocument(createData)
      setDocumentId(document.id)
      
      // Simulate upload progress
      for (let i = 0; i <= 100; i += 10) {
        setStep2Data(prev => ({ ...prev, uploadProgress: i }))
        await new Promise(resolve => setTimeout(resolve, 100))
      }

      // Upload file
      await uploadDocumentVersion(document.id, file, { 
        title: step1Data.title,
        enableOCR: step2Data.enableOCR 
      })
      
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
      setError('РћС€РёР±РєР° Р·Р°РіСЂСѓР·РєРё С„Р°Р№Р»Р°: ' + (err as Error).message)
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
      setError('РћС€РёР±РєР° РѕС‚РїСЂР°РІРєРё РЅР° СЃРѕРіР»Р°СЃРѕРІР°РЅРёРµ: ' + (err as Error).message)
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
        РњРµС‚Р°РґР°РЅРЅС‹Рµ РґРѕРєСѓРјРµРЅС‚Р°
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <TextField
            label="РќР°Р·РІР°РЅРёРµ *"
            value={step1Data.title}
            onChange={(e) => setStep1Data(prev => ({ ...prev, title: e.target.value }))}
            fullWidth
            required
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <TextField
            label="РљРѕРґ РґРѕРєСѓРјРµРЅС‚Р°"
            value={step1Data.code}
            onChange={(e) => setStep1Data(prev => ({ ...prev, code: e.target.value }))}
            fullWidth
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <FormControl fullWidth>
            <InputLabel>РўРёРї РґРѕРєСѓРјРµРЅС‚Р° *</InputLabel>
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
            label="РћРїРёСЃР°РЅРёРµ"
            value={step1Data.description}
            onChange={(e) => setStep1Data(prev => ({ ...prev, description: e.target.value }))}
            fullWidth
            multiline
            rows={3}
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <TextField
            label="РљР°С‚РµРіРѕСЂРёСЏ"
            value={step1Data.category}
            onChange={(e) => setStep1Data(prev => ({ ...prev, category: e.target.value }))}
            fullWidth
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <FormControl fullWidth>
            <InputLabel>РљР»Р°СЃСЃРёС„РёРєР°С†РёСЏ *</InputLabel>
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
                label="РўРµРіРё"
                placeholder="Р”РѕР±Р°РІРёС‚СЊ С‚РµРі"
              />
            )}
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <FormControl fullWidth>
            <InputLabel>Р’Р»Р°РґРµР»РµС† РґРѕРєСѓРјРµРЅС‚Р°</InputLabel>
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
            label="Р”Р°С‚Р° РІСЃС‚СѓРїР»РµРЅРёСЏ РІ СЃРёР»Сѓ"
            type="date"
            value={step1Data.effectiveFrom}
            onChange={(e) => setStep1Data(prev => ({ ...prev, effectiveFrom: e.target.value }))}
            fullWidth
            InputLabelProps={{ shrink: true }}
          />
        </Grid>
        
        <Grid item xs={12} sm={6}>
          <TextField
            label="РџРµСЂРёРѕРґ РїРµСЂРµСЃРјРѕС‚СЂР° (РјРµСЃСЏС†С‹)"
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
        Р—Р°РіСЂСѓР·РєР° РІРµСЂСЃРёРё РґРѕРєСѓРјРµРЅС‚Р°
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
              РџРµСЂРµС‚Р°С‰РёС‚Рµ С„Р°Р№Р» СЃСЋРґР° РёР»Рё РЅР°Р¶РјРёС‚Рµ РґР»СЏ РІС‹Р±РѕСЂР°
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              РџРѕРґРґРµСЂР¶РёРІР°РµРјС‹Рµ С„РѕСЂРјР°С‚С‹: PDF, DOCX, TXT
            </Typography>
            <Typography variant="body2" color="text.secondary">
              РњР°РєСЃРёРјР°Р»СЊРЅС‹Р№ СЂР°Р·РјРµСЂ: 50 РњР‘
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
                Р’С‹Р±СЂР°С‚СЊ С„Р°Р№Р»
              </Button>
            </label>
          </Box>
        )}
        
        {step2Data.uploadStatus === 'uploading' && (
          <Box>
            <CircularProgress sx={{ mb: 2 }} />
            <Typography variant="h6" gutterBottom>
              Р—Р°РіСЂСѓР·РєР° С„Р°Р№Р»Р°...
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
              Р¤Р°Р№Р» СѓСЃРїРµС€РЅРѕ Р·Р°РіСЂСѓР¶РµРЅ
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
              РћС€РёР±РєР° Р·Р°РіСЂСѓР·РєРё
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
        label="Р’РєР»СЋС‡РёС‚СЊ OCR РґР»СЏ РёР·РІР»РµС‡РµРЅРёСЏ С‚РµРєСЃС‚Р°"
      />
    </Box>
  )

  const renderStep3 = () => (
    <Box>
      <Typography variant="h6" gutterBottom>
        РќР°СЃС‚СЂРѕР№РєРё СЃРѕРіР»Р°СЃРѕРІР°РЅРёСЏ
      </Typography>
      
      <FormControlLabel
        control={
          <Switch
            checked={step3Data.requiresApproval}
            onChange={(e) => setStep3Data(prev => ({ ...prev, requiresApproval: e.target.checked }))}
          />
        }
        label="РўСЂРµР±СѓРµС‚СЃСЏ СЃРѕРіР»Р°СЃРѕРІР°РЅРёРµ РґРѕРєСѓРјРµРЅС‚Р°"
        sx={{ mb: 3 }}
      />
      
      {step3Data.requiresApproval && (
        <>
          <FormControl fullWidth sx={{ mb: 3 }}>
            <InputLabel>РўРёРї РјР°СЂС€СЂСѓС‚Р°</InputLabel>
            <Select
              value={step3Data.workflowType}
              onChange={(e) => setStep3Data(prev => ({ 
                ...prev, 
                workflowType: e.target.value as 'sequential' | 'parallel' 
              }))}
            >
              <MenuItem value="sequential">РџРѕСЃР»РµРґРѕРІР°С‚РµР»СЊРЅС‹Р№</MenuItem>
              <MenuItem value="parallel">РџР°СЂР°Р»Р»РµР»СЊРЅС‹Р№</MenuItem>
            </Select>
          </FormControl>
          
          <Typography variant="subtitle1" gutterBottom>
            РЁР°РіРё СЃРѕРіР»Р°СЃРѕРІР°РЅРёСЏ:
          </Typography>
        </>
      )}
      
      {step3Data.steps.map((step, index) => (
        <Box key={index} sx={{ mb: 2, p: 2, border: 1, borderColor: 'divider', borderRadius: 1 }}>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={12} sm={3}>
              <TextField
                label="РџРѕСЂСЏРґРѕРє"
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
                <InputLabel>РЎРѕРіР»Р°СЃСѓСЋС‰РёР№</InputLabel>
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
                label="Р”РµРґР»Р°Р№РЅ"
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
                РЈРґР°Р»РёС‚СЊ
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
        Р”РѕР±Р°РІРёС‚СЊ С€Р°Рі СЃРѕРіР»Р°СЃРѕРІР°РЅРёСЏ
      </Button>
    </Box>
  )

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>РЎРѕР·РґР°РЅРёРµ РґРѕРєСѓРјРµРЅС‚Р°</DialogTitle>
      
      <DialogContent>
        <Stepper activeStep={activeStep} sx={{ mb: 4 }}>
          <Step>
            <StepLabel>РњРµС‚Р°РґР°РЅРЅС‹Рµ</StepLabel>
          </Step>
          <Step>
            <StepLabel>Р—Р°РіСЂСѓР·РєР° РІРµСЂСЃРёРё</StepLabel>
          </Step>
          <Step>
            <StepLabel>РЎРѕРіР»Р°СЃРѕРІР°РЅРёРµ</StepLabel>
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
          РћС‚РјРµРЅР°
        </Button>
        
        {activeStep === 0 && (
          <Button
            variant="contained"
            onClick={handleStep1Next}
            disabled={loading || !step1Data.title || !step1Data.type}
          >
            {loading ? <CircularProgress size={20} /> : 'Р”Р°Р»РµРµ'}
          </Button>
        )}
        
        {activeStep === 1 && (
          <>
            <Button onClick={() => setActiveStep(0)}>
              РќР°Р·Р°Рґ
            </Button>
            <Button
              variant="contained"
              onClick={handleStep2Next}
              disabled={step2Data.uploadStatus !== 'success'}
            >
              Р”Р°Р»РµРµ
            </Button>
          </>
        )}
        
        {activeStep === 2 && (
          <>
            <Button onClick={() => setActiveStep(1)}>
              РќР°Р·Р°Рґ
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
                'РћС‚РїСЂР°РІРёС‚СЊ РЅР° СЃРѕРіР»Р°СЃРѕРІР°РЅРёРµ'
              ) : (
                'РЎРѕС…СЂР°РЅРёС‚СЊ'
              )}
            </Button>
          </>
        )}
      </DialogActions>
    </Dialog>
  )
}

