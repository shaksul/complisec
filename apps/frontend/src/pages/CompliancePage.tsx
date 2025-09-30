import { useEffect, useState } from "react"
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
  Tabs,
  Tab,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from "@mui/material"
import { Add, Gavel, Assessment, Warning } from "@mui/icons-material"
import { getStandards, getRequirements, getAssessments } from "../shared/api/compliance"
import { api } from "../shared/api/client"

interface ComplianceStandard {
  id: string
  name: string
  code: string
  description?: string
  version: string
  is_active: boolean
}

interface ComplianceRequirement {
  id: string
  code: string
  title: string
  description?: string
  category?: string
  priority: string
  is_mandatory: boolean
}

interface ComplianceAssessment {
  id: string
  requirement_id: string
  status: string
  evidence?: string
  assessor_id?: string
  assessed_at?: string
  next_review_date?: string
  notes?: string
  requirement_title: string
  standard_name: string
  assessor_name?: string
}

interface ComplianceGap {
  id: string
  title: string
  description?: string
  severity: string
  status: string
  remediation_plan?: string
  target_date?: string
  responsible_id?: string
  responsible_name?: string
}

export default function CompliancePage() {
  const [activeTab, setActiveTab] = useState(0)
  const [standards, setStandards] = useState<ComplianceStandard[]>([])
  const [requirements, setRequirements] = useState<ComplianceRequirement[]>([])
  const [assessments, setAssessments] = useState<ComplianceAssessment[]>([])
  const [gaps, setGaps] = useState<ComplianceGap[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  
  // Dialog states
  const [openStandard, setOpenStandard] = useState(false)
  const [openRequirement, setOpenRequirement] = useState(false)
  const [openAssessment, setOpenAssessment] = useState(false)
  const [openGap, setOpenGap] = useState(false)
  
  // Form states
  const [selectedStandard, setSelectedStandard] = useState("")
  const [selectedAssessment, setSelectedAssessment] = useState("")
  const [formData, setFormData] = useState<any>({})

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      setError(null)
      
      const [standardsData, assessmentsData] = await Promise.all([
        getStandards(),
        getAssessments()
      ])
      
      setStandards(standardsData || [])
      setAssessments(assessmentsData || [])
    } catch (err) {
      console.error('Error loading compliance data:', err)
      setError('Ошибка загрузки данных соответствия')
    } finally {
      setLoading(false)
    }
  }

  const loadRequirements = async (standardId: string) => {
    try {
      const data = await getRequirements(standardId)
      setRequirements(data || [])
    } catch (err) {
      console.error('Error loading requirements:', err)
    }
  }

  const loadGaps = async (assessmentId: string) => {
    try {
      const res = await api.get(`/compliance/assessments/${assessmentId}/gaps`)
      setGaps(res.data.data || [])
    } catch (err) {
      console.error('Error loading gaps:', err)
    }
  }

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue)
    if (newValue === 1 && selectedStandard) {
      loadRequirements(selectedStandard)
    } else if (newValue === 3 && selectedAssessment) {
      loadGaps(selectedAssessment)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'compliant': return 'success'
      case 'non-compliant': return 'error'
      case 'pending': return 'warning'
      default: return 'default'
    }
  }

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'high': return 'error'
      case 'medium': return 'warning'
      case 'low': return 'info'
      default: return 'default'
    }
  }

  const renderStandards = () => (
    <Paper>
      <Box display="flex" justifyContent="space-between" alignItems="center" p={2}>
        <Typography variant="h6">Стандарты соответствия</Typography>
        <Button variant="contained" startIcon={<Add />} onClick={() => setOpenStandard(true)}>
          Добавить стандарт
        </Button>
      </Box>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Название</TableCell>
              <TableCell>Код</TableCell>
              <TableCell>Версия</TableCell>
              <TableCell>Статус</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {standards.map((standard, index) => (
              <TableRow key={standard.id || `standard-${index}`}> 
                <TableCell>
                  <Box display="flex" alignItems="center">
                    <Gavel sx={{ mr: 1 }} />
                    {standard.name}
                  </Box>
                </TableCell>
                <TableCell>{standard.code}</TableCell>
                <TableCell>{standard.version}</TableCell>
                <TableCell>
                  <Chip
                    label={standard.is_active ? "Активен" : "Неактивен"}
                    color={standard.is_active ? "success" : "default"}
                    size="small"
                  />
                </TableCell>
                <TableCell>
                  <Button 
                    size="small" 
                    onClick={() => {
                      setSelectedStandard(standard.id)
                      loadRequirements(standard.id)
                    }}
                  >
                    Требования
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Paper>
  )

  const renderRequirements = () => (
    <Paper>
      <Box display="flex" justifyContent="space-between" alignItems="center" p={2}>
        <Typography variant="h6">Требования стандарта</Typography>
        <Button 
          variant="contained" 
          startIcon={<Add />} 
          onClick={() => setOpenRequirement(true)}
          disabled={!selectedStandard}
        >
          Добавить требование
        </Button>
      </Box>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Код</TableCell>
              <TableCell>Название</TableCell>
              <TableCell>Категория</TableCell>
              <TableCell>Приоритет</TableCell>
              <TableCell>Обязательность</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {requirements.map((requirement) => (
              <TableRow key={requirement.id}>
                <TableCell>{requirement.code}</TableCell>
                <TableCell>{requirement.title}</TableCell>
                <TableCell>{requirement.category || '-'}</TableCell>
                <TableCell>
                  <Chip
                    label={requirement.priority}
                    color={requirement.priority === 'high' ? 'error' : requirement.priority === 'medium' ? 'warning' : 'info'}
                    size="small"
                  />
                </TableCell>
                <TableCell>
                  <Chip
                    label={requirement.is_mandatory ? "Обязательное" : "Рекомендуемое"}
                    color={requirement.is_mandatory ? "error" : "default"}
                    size="small"
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Paper>
  )

  const renderAssessments = () => (
    <Paper>
      <Box display="flex" justifyContent="space-between" alignItems="center" p={2}>
        <Typography variant="h6">Оценки соответствия</Typography>
        <Button variant="contained" startIcon={<Add />} onClick={() => setOpenAssessment(true)}>
          Добавить оценку
        </Button>
      </Box>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Стандарт</TableCell>
              <TableCell>Требование</TableCell>
              <TableCell>Статус</TableCell>
              <TableCell>Оценщик</TableCell>
              <TableCell>Дата оценки</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {assessments.map((assessment) => (
              <TableRow key={assessment.id}>
                <TableCell>{assessment.standard_name}</TableCell>
                <TableCell>{assessment.requirement_title}</TableCell>
                <TableCell>
                  <Chip
                    label={assessment.status}
                    color={getStatusColor(assessment.status)}
                    size="small"
                  />
                </TableCell>
                <TableCell>{assessment.assessor_name || '-'}</TableCell>
                <TableCell>{assessment.assessed_at ? new Date(assessment.assessed_at).toLocaleDateString() : '-'}</TableCell>
                <TableCell>
                  <Button 
                    size="small" 
                    onClick={() => {
                      setSelectedAssessment(assessment.id)
                      loadGaps(assessment.id)
                    }}
                  >
                    Пробелы
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Paper>
  )

  const renderGaps = () => (
    <Paper>
      <Box display="flex" justifyContent="space-between" alignItems="center" p={2}>
        <Typography variant="h6">Пробелы соответствия</Typography>
        <Button 
          variant="contained" 
          startIcon={<Add />} 
          onClick={() => setOpenGap(true)}
          disabled={!selectedAssessment}
        >
          Добавить пробел
        </Button>
      </Box>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Название</TableCell>
              <TableCell>Критичность</TableCell>
              <TableCell>Статус</TableCell>
              <TableCell>Ответственный</TableCell>
              <TableCell>Целевая дата</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {gaps.map((gap) => (
              <TableRow key={gap.id}>
                <TableCell>{gap.title}</TableCell>
                <TableCell>
                  <Chip
                    label={gap.severity}
                    color={getSeverityColor(gap.severity)}
                    size="small"
                  />
                </TableCell>
                <TableCell>
                  <Chip
                    label={gap.status}
                    color={gap.status === 'closed' ? 'success' : 'warning'}
                    size="small"
                  />
                </TableCell>
                <TableCell>{gap.responsible_name || '-'}</TableCell>
                <TableCell>{gap.target_date ? new Date(gap.target_date).toLocaleDateString() : '-'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Paper>
  )

  if (loading) {
    return (
      <Container maxWidth="lg">
        <Typography>Загрузка...</Typography>
      </Container>
    )
  }

  if (error) {
    return (
      <Container maxWidth="lg">
        <Paper sx={{ p: 2, mb: 2, bgcolor: 'error.light', color: 'error.contrastText' }}>
          <Typography>{error}</Typography>
        </Paper>
      </Container>
    )
  }

  return (
    <Container maxWidth="lg">
      <Typography variant="h4" gutterBottom>
        Соответствие стандартам
      </Typography>

      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        <Tabs value={activeTab} onChange={handleTabChange}>
          <Tab label="Стандарты" icon={<Gavel />} />
          <Tab label="Требования" icon={<Assessment />} />
          <Tab label="Оценки" icon={<Assessment />} />
          <Tab label="Пробелы" icon={<Warning />} />
        </Tabs>
      </Box>

      {activeTab === 0 && renderStandards()}
      {activeTab === 1 && renderRequirements()}
      {activeTab === 2 && renderAssessments()}
      {activeTab === 3 && renderGaps()}
    </Container>
  )
}
