import { useState, useEffect } from "react"
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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Select,
  Tabs,
  Tab,
  MenuItem,
  FormControl,
  InputLabel,
  IconButton,
  Tooltip,
  Alert,
  CircularProgress,
} from "@mui/material"
import {
  Add,
  Description,
  Edit,
  Delete,
  Visibility,
  History,
  Search,
  Download,
  Psychology,
} from "@mui/icons-material"
import DocumentPreview from "../components/docs/DocumentPreview"
import {
  getDocuments,
  createDocument,
  updateDocument,
  deleteDocument,
  downloadDocument,
  getDocumentTypeLabel,
  getDocumentStatusLabel,
  getDocumentStatusColor,
  type Document,
  type CreateDocumentDTO,
  type DocumentFilters,
} from "../shared/api/documents"
import { indexDocument } from "../shared/api/rag"
import { useAuth } from "../contexts/AuthContext"
import CreateDocumentWizard from "../components/docs/CreateDocumentWizard"
import DocumentVersionsDialog from "../components/docs/DocumentVersionsDialog"
import { formatTextForDisplay } from "../shared/utils/textNormalization"


function DocumentsPage() {
  const [documents, setDocuments] = useState<Document[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState("")
  const [statusFilter, setStatusFilter] = useState("")
  const [typeFilter, setTypeFilter] = useState("")
  
  // Dialog states
  const [openCreate, setOpenCreate] = useState(false)
  const [openCreateWizard, setOpenCreateWizard] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [openView, setOpenView] = useState(false)
  const [openVersions, setOpenVersions] = useState(false)
  const [selectedDocument, setSelectedDocument] = useState<Document | null>(null)
  const [viewTab, setViewTab] = useState(0)
  
  // Form states
  const [formData, setFormData] = useState<CreateDocumentDTO>({
    title: "",
    description: "",
    type: "policy",
    category: "",
    tags: [],
  })

  const { user } = useAuth()

  useEffect(() => {
    if (user) {
      loadDocuments()
    } else {
      setLoading(false)
    }
  }, [searchTerm, statusFilter, typeFilter, user])

  const loadDocuments = async () => {
    try {
      setLoading(true)
      setError(null)
      
      const filters: DocumentFilters = {}
      if (searchTerm) filters.search = searchTerm
      if (statusFilter) filters.status = statusFilter
      if (typeFilter) filters.type = typeFilter
      
      const data = await getDocuments(filters)
      setDocuments(Array.isArray(data) ? data : [])
    } catch (err) {
      console.error('Error loading documents:', err)
      setError('Ошибка загрузки документов')
      setDocuments([])
    } finally {
      setLoading(false)
    }
  }


  const handleCreateDocument = async () => {
    try {
      await createDocument(formData)
      setOpenCreate(false)
      resetForm()
      loadDocuments()
    } catch (err) {
      console.error('Error creating document:', err)
      setError('Ошибка создания документа')
    }
  }

  const handleCreateWizardSuccess = (documentId: string) => {
    setOpenCreateWizard(false)
    loadDocuments()
    // Можно добавить редирект на страницу документа
    console.log('Document created with ID:', documentId)
  }

  const handleUpdateDocument = async () => {
    if (!selectedDocument) return
    
    try {
      await updateDocument(selectedDocument.id, {
        ...formData,
        status: selectedDocument.status,
      })
      setOpenEdit(false)
      resetForm()
      loadDocuments()
    } catch (err) {
      console.error('Error updating document:', err)
      setError('Ошибка обновления документа')
    }
  }

  const handleDeleteDocument = async () => {
    if (!selectedDocument) return
    
    try {
      await deleteDocument(selectedDocument.id)
      setOpenDelete(false)
      loadDocuments()
    } catch (err) {
      console.error('Error deleting document:', err)
      setError('Ошибка удаления документа')
    }
  }

  const resetForm = () => {
    setFormData({
      title: "",
      description: "",
      type: "policy",
      category: "",
      tags: [],
    })
    setSelectedDocument(null)
  }

  const openViewDialog = (document: Document) => {
    setSelectedDocument(document)
    setOpenView(true)
  }

  const openEditDialog = (document: Document) => {
    setSelectedDocument(document)
    setFormData({
      title: document.title,
      description: document.description || "",
      type: document.type,
      category: document.category || "",
      tags: document.tags,
    })
    setOpenEdit(true)
  }

  const openDeleteDialog = (document: Document) => {
    setSelectedDocument(document)
    setOpenDelete(true)
  }

  const openVersionsDialog = (document: Document) => {
    setSelectedDocument(document)
    setOpenVersions(true)
  }

  const handleDownloadDocument = async (doc: Document) => {
    try {
      const blob = await downloadDocument(doc.id)
      const url = window.URL.createObjectURL(blob)
      const link = window.document.createElement('a')
      link.href = url
      link.download = `${doc.title}.${getFileExtension(doc.storage_key || '')}`
      window.document.body.appendChild(link)
      link.click()
      window.document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    } catch (err) {
      console.error('Error downloading document:', err)
      setError('Ошибка скачивания документа')
    }
  }

  const getFileExtension = (filename: string): string => {
    const ext = filename.split('.').pop()
    return ext || 'txt'
  }

  const handleIndexToRAG = async (doc: Document) => {
    try {
      await indexDocument(doc.id)
      setError(null)
      // Показываем уведомление об успехе
      alert(`Документ "${doc.title}" отправлен на индексацию в GraphRAG`)
    } catch (err) {
      console.error('Error indexing document to RAG:', err)
      setError('Ошибка индексации документа в RAG')
    }
  }

  const renderDocumentsList = () => (
    <Paper>
      <Box display="flex" justifyContent="space-between" alignItems="center" p={2}>
        <Typography variant="h6">Документы</Typography>
        <Button variant="contained" startIcon={<Add />} onClick={() => setOpenCreateWizard(true)}>
          Создать документ
        </Button>
      </Box>
      
      <Box p={2} display="flex" gap={2} alignItems="center">
        <TextField
          size="small"
          placeholder="Поиск документов..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          InputProps={{
            startAdornment: <Search sx={{ mr: 1, color: 'text.secondary' }} />
          }}
        />
        
        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel>Статус</InputLabel>
          <Select
            value={statusFilter}
            label="Статус"
            onChange={(e) => setStatusFilter(e.target.value)}
          >
            <MenuItem value="">Все</MenuItem>
            <MenuItem value="draft">Черновик</MenuItem>
            <MenuItem value="in_review">На согласовании</MenuItem>
            <MenuItem value="approved">Утвержден</MenuItem>
            <MenuItem value="obsolete">Устарел</MenuItem>
          </Select>
        </FormControl>
        
        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel>Тип</InputLabel>
          <Select
            value={typeFilter}
            label="Тип"
            onChange={(e) => setTypeFilter(e.target.value)}
          >
            <MenuItem value="">Все</MenuItem>
            <MenuItem value="policy">Политика</MenuItem>
            <MenuItem value="standard">Стандарт</MenuItem>
            <MenuItem value="procedure">Процедура</MenuItem>
            <MenuItem value="instruction">Инструкция</MenuItem>
            <MenuItem value="act">Акт</MenuItem>
            <MenuItem value="other">Другое</MenuItem>
          </Select>
        </FormControl>
      </Box>
      
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Название</TableCell>
              <TableCell>Тип</TableCell>
              <TableCell>Категория</TableCell>
              <TableCell>Статус</TableCell>
              <TableCell>Версия</TableCell>
              <TableCell>Создан</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={7} align="center">
                  <CircularProgress />
                </TableCell>
              </TableRow>
            ) : documents.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} align="center">
                  <Typography>Нет документов</Typography>
                </TableCell>
              </TableRow>
            ) : (
              documents.map((document) => (
                <TableRow key={document.id}>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      <Description sx={{ mr: 1 }} />
                      {formatTextForDisplay(document.title, 50)}
                    </Box>
                  </TableCell>
                  <TableCell>{getDocumentTypeLabel(document.type)}</TableCell>
                  <TableCell>{formatTextForDisplay(document.category) || '-'}</TableCell>
                  <TableCell>
                    <Chip
                      label={getDocumentStatusLabel(document.status)}
                      color={getDocumentStatusColor(document.status)}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>v{document.current_version}</TableCell>
                  <TableCell>
                    {new Date(document.created_at).toLocaleDateString()}
                  </TableCell>
                  <TableCell>
                    <Box display="flex" gap={1}>
                      <Tooltip title="Скачать">
                        <IconButton size="small" onClick={() => handleDownloadDocument(document)}>
                          <Download />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Просмотр">
                        <IconButton size="small" onClick={() => openViewDialog(document)}>
                          <Visibility />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Редактировать">
                        <IconButton size="small" onClick={() => openEditDialog(document)}>
                          <Edit />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Версии">
                        <IconButton size="small" onClick={() => openVersionsDialog(document)}>
                          <History />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Индексировать в RAG">
                        <IconButton 
                          size="small" 
                          color="primary"
                          onClick={() => handleIndexToRAG(document)}
                        >
                          <Psychology />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Удалить">
                        <IconButton 
                          size="small" 
                          color="error"
                          onClick={() => openDeleteDialog(document)}
                        >
                          <Delete />
                        </IconButton>
                      </Tooltip>
                    </Box>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>
    </Paper>
  )


  if (error) {
    return (
      <Container maxWidth="lg">
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      </Container>
    )
  }

  return (
    <Container maxWidth="lg">
      <Typography variant="h4" gutterBottom>
        Управление документами
      </Typography>

      {renderDocumentsList()}

      {/* Create Document Dialog */}
      <Dialog open={openCreate} onClose={() => setOpenCreate(false)} maxWidth="md" fullWidth>
        <DialogTitle>Создать документ</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Название"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              fullWidth
              required
            />
            <TextField
              label="Описание"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              fullWidth
              multiline
              rows={3}
            />
            <FormControl fullWidth>
              <InputLabel>Тип документа</InputLabel>
              <Select
                value={formData.type}
                label="Тип документа"
                onChange={(e) => setFormData({ ...formData, type: e.target.value as any })}
              >
                <MenuItem value="policy">Политика</MenuItem>
                <MenuItem value="standard">Стандарт</MenuItem>
                <MenuItem value="procedure">Процедура</MenuItem>
                <MenuItem value="instruction">Инструкция</MenuItem>
                <MenuItem value="act">Акт</MenuItem>
                <MenuItem value="other">Другое</MenuItem>
              </Select>
            </FormControl>
            <TextField
              label="Категория"
              value={formData.category}
              onChange={(e) => setFormData({ ...formData, category: e.target.value })}
              fullWidth
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenCreate(false)}>Отмена</Button>
          <Button onClick={handleCreateDocument} variant="contained">
            Создать
          </Button>
        </DialogActions>
      </Dialog>

      {/* Edit Document Dialog */}
      <Dialog open={openEdit} onClose={() => setOpenEdit(false)} maxWidth="md" fullWidth>
        <DialogTitle>Редактировать документ</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Название"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              fullWidth
              required
            />
            <TextField
              label="Описание"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              fullWidth
              multiline
              rows={3}
            />
            <FormControl fullWidth>
              <InputLabel>Тип документа</InputLabel>
              <Select
                value={formData.type}
                label="Тип документа"
                onChange={(e) => setFormData({ ...formData, type: e.target.value as any })}
              >
                <MenuItem value="policy">Политика</MenuItem>
                <MenuItem value="standard">Стандарт</MenuItem>
                <MenuItem value="procedure">Процедура</MenuItem>
                <MenuItem value="instruction">Инструкция</MenuItem>
                <MenuItem value="act">Акт</MenuItem>
                <MenuItem value="other">Другое</MenuItem>
              </Select>
            </FormControl>
            <TextField
              label="Категория"
              value={formData.category}
              onChange={(e) => setFormData({ ...formData, category: e.target.value })}
              fullWidth
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenEdit(false)}>Отмена</Button>
          <Button onClick={handleUpdateDocument} variant="contained">
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>

      {/* View Document Dialog */}
      <Dialog open={openView} onClose={() => setOpenView(false)} maxWidth="lg" fullWidth>
        <DialogTitle>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="h6">
              {selectedDocument ? formatTextForDisplay(selectedDocument.title) : "Просмотр документа"}
            </Typography>
            <Button onClick={() => setOpenView(false)} size="small">
              Закрыть
            </Button>
          </Box>
        </DialogTitle>
        <DialogContent sx={{ p: 0 }}>
          {selectedDocument && (
            <Box>
              <Tabs 
                value={viewTab} 
                onChange={(_, newValue) => setViewTab(newValue)}
                sx={{ borderBottom: 1, borderColor: 'divider', px: 2 }}
              >
                <Tab label="Метаданные" />
                <Tab label="Предпросмотр" />
              </Tabs>
              
              {viewTab === 0 && (
                <Box p={2}>
                  <Box mb={2}>
                    <Typography variant="body2" color="textSecondary">
                      <strong>Тип:</strong> {getDocumentTypeLabel(selectedDocument.type)}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      <strong>Статус:</strong> {getDocumentStatusLabel(selectedDocument.status)}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      <strong>Категория:</strong> {formatTextForDisplay(selectedDocument.category) || "Не указана"}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      <strong>Версия:</strong> v{selectedDocument.current_version}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      <strong>Создан:</strong> {new Date(selectedDocument.created_at).toLocaleString()}
                    </Typography>
                  </Box>
                  {selectedDocument.description && (
                    <Box mb={2}>
                      <Typography variant="subtitle2" gutterBottom>
                        Описание:
                      </Typography>
                      <Typography variant="body2">
                        {formatTextForDisplay(selectedDocument.description)}
                      </Typography>
                    </Box>
                  )}
                  {selectedDocument.tags && selectedDocument.tags.length > 0 && (
                    <Box mb={2}>
                      <Typography variant="subtitle2" gutterBottom>
                        Теги:
                      </Typography>
                      <Box display="flex" gap={1} flexWrap="wrap">
                        {selectedDocument.tags.map((tag, index) => (
                          <Chip key={index} label={tag} size="small" />
                        ))}
                      </Box>
                    </Box>
                  )}
                </Box>
              )}
              
              {viewTab === 1 && (
                <Box>
                  {selectedDocument.id && (
                    <DocumentPreview
                      documentId={selectedDocument.id}
                      fileName={selectedDocument.title}
                      mimeType={selectedDocument.mime_type}
                    />
                  )}
                </Box>
              )}
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenView(false)}>Закрыть</Button>
          {selectedDocument && (
            <Button 
              variant="contained" 
              startIcon={<Download />}
              onClick={() => handleDownloadDocument(selectedDocument)}
            >
              Скачать
            </Button>
          )}
        </DialogActions>
      </Dialog>

      {/* Delete Document Dialog */}
      <Dialog open={openDelete} onClose={() => setOpenDelete(false)}>
        <DialogTitle>Удалить документ</DialogTitle>
        <DialogContent>
          <Typography>
            Вы уверены, что хотите удалить документ "{selectedDocument?.title}"?
            Это действие нельзя отменить.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDelete(false)}>Отмена</Button>
          <Button onClick={handleDeleteDocument} variant="contained" color="error">
            Удалить
          </Button>
        </DialogActions>
      </Dialog>

      {/* Create Document Wizard */}
      <CreateDocumentWizard
        open={openCreateWizard}
        onClose={() => setOpenCreateWizard(false)}
        onSuccess={handleCreateWizardSuccess}
      />

      {/* Document Versions Dialog */}
      <DocumentVersionsDialog
        open={openVersions}
        onClose={() => setOpenVersions(false)}
        documentId={selectedDocument?.id || ''}
        documentTitle={selectedDocument?.title || ''}
      />
    </Container>
  )
}

export { DocumentsPage }
export default DocumentsPage