import { useState, useEffect } from "react"
import {
  Container,
  Typography,
  Paper,
  Box,
  Button,
  TextField,
  IconButton,
  Tooltip,
  Alert,
  CircularProgress,
  Breadcrumbs,
  Link,
  Chip,
  Grid,
  Card,
  CardContent,
  CardActions,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Tabs,
  Tab,
} from "@mui/material"
import {
  Upload,
  Folder,
  FolderOpen,
  InsertDriveFile,
  Search,
  Download,
  Edit,
  Delete,
  Visibility,
  Home,
  CloudUpload,
  Description,
  Image,
  VideoFile,
  AudioFile,
  PictureAsPdf,
  TableChart,
  Slideshow,
  TextSnippet,
  FolderZip,
  ArrowBack,
} from "@mui/icons-material"
import {
  FileDocument,
  Folder as FolderType,
  DocumentStats,
  createFolder,
  listFolders,
  uploadDocument,
  listDocuments,
  listStructuredDocuments,
  deleteDocument,
  deleteFolder,
  downloadDocument,
  searchDocuments,
  getDocumentStats,
  formatFileSize,
  getFileIcon,
  getMimeTypeLabel,
} from "../shared/api/file-documents"
import { useAuth } from "../contexts/AuthContext"

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
      id={`file-documents-tabpanel-${index}`}
      aria-labelledby={`file-documents-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  )
}

function FileDocumentsPage() {
  const { } = useAuth()
  const [currentTab, setCurrentTab] = useState(0)
  const [folders, setFolders] = useState<FolderType[]>([])
  const [documents, setDocuments] = useState<FileDocument[]>([])
  const [structuredData, setStructuredData] = useState<any>(null)
  const [currentFolder, setCurrentFolder] = useState<FolderType | null>(null)
  const [currentModule, setCurrentModule] = useState<string | null>(null)
  const [currentCategory, setCurrentCategory] = useState<string | null>(null)
  const [currentModuleData, setCurrentModuleData] = useState<any>(null)
  const [currentCategoryData, setCurrentCategoryData] = useState<any[]>([])
  const [breadcrumbs, setBreadcrumbs] = useState<FolderType[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState("")
  const [stats, setStats] = useState<DocumentStats | null>(null)
  
  // Dialog states
  const [openCreateFolder, setOpenCreateFolder] = useState(false)
  const [openUpload, setOpenUpload] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [selectedItem, setSelectedItem] = useState<FileDocument | FolderType | null>(null)
  
  // Form states
  const [folderName, setFolderName] = useState("")
  const [folderDescription, setFolderDescription] = useState("")
  const [uploadFile, setUploadFile] = useState<File | null>(null)
  const [documentName, setDocumentName] = useState("")
  const [documentDescription, setDocumentDescription] = useState("")
  const [documentTags, setDocumentTags] = useState<string[]>([])
  const [tagInput, setTagInput] = useState("")

  useEffect(() => {
    loadData()
  }, [currentFolder])

  useEffect(() => {
    if (currentTab === 2) {
      loadStats()
    }
  }, [currentTab])

  const loadData = async () => {
    try {
      setLoading(true)
      setError(null)
      
      if (currentTab === 0) {
        // Загружаем структурированные данные для файлового хранилища
        const structuredData = await listStructuredDocuments()
        setStructuredData(structuredData)
        
        // Также загружаем обычные папки и документы для совместимости
        const [foldersData, documentsData] = await Promise.all([
          listFolders(currentFolder?.id),
          listDocuments({
            folder_id: currentFolder?.id,
            page: 1,
            limit: 100,
          })
        ])
        
        setFolders(foldersData)
        setDocuments(documentsData)
      } else {
        // Для других вкладок используем обычную загрузку
        const [foldersData, documentsData] = await Promise.all([
          listFolders(currentFolder?.id),
          listDocuments({
            folder_id: currentFolder?.id,
            page: 1,
            limit: 100,
          })
        ])
        
        setFolders(foldersData)
        setDocuments(documentsData)
      }
    } catch (err) {
      console.error('Error loading data:', err)
      setError('Ошибка загрузки данных')
    } finally {
      setLoading(false)
    }
  }

  const loadStats = async () => {
    try {
      const statsData = await getDocumentStats()
      setStats(statsData)
    } catch (err) {
      console.error('Error loading stats:', err)
    }
  }

  const handleCreateFolder = async () => {
    try {
      await createFolder({
        name: folderName,
        description: folderDescription || undefined,
        parent_id: currentFolder?.id,
      })
      setOpenCreateFolder(false)
      setFolderName("")
      setFolderDescription("")
      loadData()
    } catch (err) {
      console.error('Error creating folder:', err)
      setError('Ошибка создания папки')
    }
  }

  const handleUploadDocument = async () => {
    if (!uploadFile) return
    
    try {
      await uploadDocument(uploadFile, {
        name: documentName || uploadFile.name,
        description: documentDescription || undefined,
        folder_id: currentFolder?.id,
        tags: documentTags,
        enable_ocr: true,
      })
      setOpenUpload(false)
      setUploadFile(null)
      setDocumentName("")
      setDocumentDescription("")
      setDocumentTags([])
      setTagInput("")
      loadData()
    } catch (err) {
      console.error('Error uploading document:', err)
      setError('Ошибка загрузки документа')
    }
  }

  const handleDeleteItem = async () => {
    if (!selectedItem) return
    
    try {
      if ('file_path' in selectedItem) {
        // It's a document
        await deleteDocument(selectedItem.id)
      } else {
        // It's a folder
        await deleteFolder(selectedItem.id)
      }
      setOpenDelete(false)
      setSelectedItem(null)
      loadData()
    } catch (err) {
      console.error('Error deleting item:', err)
      setError('Ошибка удаления')
    }
  }

  const handleDownloadDocument = async (doc: FileDocument) => {
    try {
      const blob = await downloadDocument(doc.id)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = doc.original_name
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    } catch (err) {
      console.error('Error downloading document:', err)
      setError('Ошибка скачивания документа')
    }
  }

  const handleFolderClick = (folder: FolderType) => {
    setCurrentFolder(folder)
    setBreadcrumbs(prev => [...prev, folder])
  }

  const handleBreadcrumbClick = (index: number) => {
    const newBreadcrumbs = breadcrumbs.slice(0, index + 1)
    setBreadcrumbs(newBreadcrumbs)
    setCurrentFolder(newBreadcrumbs[newBreadcrumbs.length - 1] || null)
  }

  const handleModuleClick = (moduleName: string, moduleData: any) => {
    setCurrentModule(moduleName)
    setCurrentModuleData(moduleData)
    setCurrentCategory(null)
    setCurrentCategoryData([])
  }

  const handleCategoryClick = (_moduleName: string, categoryName: string, documents: any[]) => {
    setCurrentCategory(categoryName)
    setCurrentCategoryData(Array.isArray(documents) ? documents : [])
  }

  const handleSearch = async () => {
    if (!searchTerm.trim()) return
    
    try {
      setLoading(true)
      const results = await searchDocuments(searchTerm)
      // Convert search results to documents format
      const searchResults = results.map(result => ({
        id: result.document_id,
        name: result.name,
        description: result.description,
        mime_type: result.mime_type,
        file_size: result.file_size,
        created_at: result.created_at,
        // Add other required fields with defaults
        tenant_id: '',
        original_name: result.name,
        file_path: '',
        file_hash: '',
        owner_id: '',
        created_by: '',
        updated_at: result.created_at,
        is_active: true,
        version: 1,
        tags: [],
        links: [],
        ocr_text: result.ocr_text,
      }))
      setDocuments(searchResults)
    } catch (err) {
      console.error('Error searching documents:', err)
      setError('Ошибка поиска')
    } finally {
      setLoading(false)
    }
  }

  const addTag = () => {
    if (tagInput.trim() && !documentTags.includes(tagInput.trim())) {
      setDocumentTags([...documentTags, tagInput.trim()])
      setTagInput("")
    }
  }

  const removeTag = (tag: string) => {
    setDocumentTags(documentTags.filter(t => t !== tag))
  }

  const getFileIconComponent = (mimeType: string) => {
    const iconName = getFileIcon(mimeType)
    const iconMap: Record<string, React.ReactElement> = {
      image: <Image />,
      video: <VideoFile />,
      audio: <AudioFile />,
      picture_as_pdf: <PictureAsPdf />,
      description: <Description />,
      table_chart: <TableChart />,
      slideshow: <Slideshow />,
      text_snippet: <TextSnippet />,
      folder_zip: <FolderZip />,
    }
    return iconMap[iconName] || <InsertDriveFile />
  }

  const renderStructuredDocuments = () => {
    if (!structuredData?.modules) return null

    return (
      <Grid container spacing={2}>
        {Object.entries(structuredData.modules).map(([moduleName, moduleData]: [string, any]) => (
          <Grid item xs={12} sm={6} md={4} lg={3} key={moduleName}>
            <Card
              sx={{
                '&:hover': { boxShadow: 3 },
                cursor: 'pointer',
              }}
              onClick={() => handleModuleClick(moduleName, moduleData)}
            >
              <CardContent>
                <Box display="flex" alignItems="center" mb={1}>
                  <FolderOpen color="primary" sx={{ mr: 1 }} />
                  <Typography variant="h6" noWrap sx={{ flex: 1 }}>
                    {moduleName === 'documents' ? 'Документы' : 
                     moduleName === 'risks' ? 'Риски' :
                     moduleName === 'assets' ? 'Активы' :
                     moduleName === 'general' ? 'Общие' : moduleName}
                  </Typography>
                </Box>
                <Typography variant="body2" color="text.secondary">
                  {moduleData.categories ? Object.keys(moduleData.categories).length : 0} категорий
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {moduleData.categories ? 
                    Object.values(moduleData.categories).reduce((total: number, docs: any) => 
                      total + (Array.isArray(docs) ? docs.length : 0), 0) : 0} документов
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    )
  }

  const renderModuleDocuments = (moduleName: string, moduleData: any) => {
    return (
      <Box>
        <Box display="flex" alignItems="center" mb={2}>
          <IconButton onClick={() => setCurrentModule(null)} sx={{ mr: 1 }}>
            <ArrowBack />
          </IconButton>
          <Typography variant="h6">
            {moduleName === 'documents' ? 'Документы' : 
             moduleName === 'risks' ? 'Риски' :
             moduleName === 'assets' ? 'Активы' :
             moduleName === 'general' ? 'Общие' : moduleName}
          </Typography>
        </Box>

        <Grid container spacing={2}>
          {moduleData.categories && Object.entries(moduleData.categories).map(([categoryName, documents]: [string, any]) => (
            <Grid item xs={12} sm={6} md={4} lg={3} key={categoryName}>
              <Card
                sx={{
                  '&:hover': { boxShadow: 3 },
                  cursor: 'pointer',
                }}
                onClick={() => handleCategoryClick(moduleName, categoryName, documents)}
              >
                <CardContent>
                  <Box display="flex" alignItems="center" mb={1}>
                    <FolderOpen color="primary" sx={{ mr: 1 }} />
                    <Typography variant="h6" noWrap sx={{ flex: 1 }}>
                      {categoryName === 'uncategorized' ? 'Без категории' : categoryName}
                    </Typography>
                  </Box>
                  <Typography variant="body2" color="text.secondary">
                    {Array.isArray(documents) ? documents.length : 0} документов
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Box>
    )
  }

  const renderCategoryDocuments = (moduleName: string, categoryName: string, documents: any[]) => {
    return (
      <Box>
        <Box display="flex" alignItems="center" mb={2}>
          <IconButton onClick={() => setCurrentCategory(null)} sx={{ mr: 1 }}>
            <ArrowBack />
          </IconButton>
          <Typography variant="h6">
            {moduleName === 'documents' ? 'Документы' : 
             moduleName === 'risks' ? 'Риски' :
             moduleName === 'assets' ? 'Активы' :
             moduleName === 'general' ? 'Общие' : moduleName} → {categoryName === 'uncategorized' ? 'Без категории' : categoryName}
          </Typography>
        </Box>

        <Grid container spacing={2}>
          {Array.isArray(documents) && documents.map((document: any) => (
            <Grid item xs={12} sm={6} md={4} lg={3} key={document.id}>
              <Card
                sx={{
                  '&:hover': { boxShadow: 3 },
                }}
              >
                <CardContent>
                  <Box display="flex" alignItems="center" mb={1}>
                    {getFileIconComponent(document.mime_type)}
                    <Typography variant="h6" noWrap sx={{ ml: 1, flex: 1 }}>
                      {document.title}
                    </Typography>
                  </Box>
                  <Typography variant="body2" color="text.secondary" mb={1}>
                    {getMimeTypeLabel(document.mime_type)} • {formatFileSize(document.file_size)}
                  </Typography>
                  {document.description && (
                    <Typography variant="body2" color="text.secondary" noWrap>
                      {document.description}
                    </Typography>
                  )}
                  {document.tags && document.tags.length > 0 && (
                    <Box mt={1}>
                      {document.tags.slice(0, 2).map((tag: string) => (
                        <Chip
                          key={tag}
                          label={tag}
                          size="small"
                          sx={{ mr: 0.5, mb: 0.5 }}
                        />
                      ))}
                      {document.tags.length > 2 && (
                        <Chip
                          label={`+${document.tags.length - 2}`}
                          size="small"
                          variant="outlined"
                        />
                      )}
                    </Box>
                  )}
                </CardContent>
                <CardActions>
                  <Tooltip title="Скачать">
                    <IconButton
                      size="small"
                      onClick={() => handleDownloadDocument(document)}
                    >
                      <Download />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Просмотр">
                    <IconButton
                      size="small"
                      onClick={() => {
                        // TODO: Implement view functionality
                        console.log('View document:', document)
                      }}
                    >
                      <Visibility />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Редактировать">
                    <IconButton
                      size="small"
                      onClick={() => {
                        // TODO: Implement edit functionality
                        console.log('Edit document:', document)
                      }}
                    >
                      <Edit />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Удалить">
                    <IconButton
                      size="small"
                      onClick={() => {
                        setSelectedItem(document)
                        setOpenDelete(true)
                      }}
                      color="error"
                    >
                      <Delete />
                    </IconButton>
                  </Tooltip>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Box>
    )
  }

  const renderBreadcrumbs = () => (
    <Breadcrumbs sx={{ mb: 2 }}>
      <Link
        component="button"
        variant="body1"
        onClick={() => {
          setCurrentFolder(null)
          setBreadcrumbs([])
        }}
        sx={{ display: 'flex', alignItems: 'center' }}
      >
        <Home sx={{ mr: 0.5 }} />
        Корневая папка
      </Link>
      {breadcrumbs.map((folder, index) => (
        <Link
          key={folder.id}
          component="button"
          variant="body1"
          onClick={() => handleBreadcrumbClick(index)}
          sx={{ display: 'flex', alignItems: 'center' }}
        >
          <Folder sx={{ mr: 0.5 }} />
          {folder.name}
        </Link>
      ))}
    </Breadcrumbs>
  )

  const renderDocumentsList = () => (
    <Box>
      {renderBreadcrumbs()}
      
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
        <Typography variant="h6">
          {currentFolder ? currentFolder.name : 'Корневая папка'}
        </Typography>
        <Box display="flex" gap={1}>
          <Button
            variant="outlined"
            startIcon={<Folder />}
            onClick={() => setOpenCreateFolder(true)}
          >
            Создать папку
          </Button>
          <Button
            variant="contained"
            startIcon={<Upload />}
            onClick={() => setOpenUpload(true)}
          >
            Загрузить файл
          </Button>
        </Box>
      </Box>

      {loading ? (
        <Box display="flex" justifyContent="center" p={4}>
          <CircularProgress />
        </Box>
      ) : (
        <>
          {/* Структурированные документы для корневой папки */}
          {!currentFolder && !currentModule && !currentCategory && structuredData && renderStructuredDocuments()}
          
          {/* Отображение модуля */}
          {currentModule && currentModuleData && !currentCategory && renderModuleDocuments(currentModule, currentModuleData)}
          
          {/* Отображение категории */}
          {currentModule && currentCategory && currentCategoryData && renderCategoryDocuments(currentModule, currentCategory, currentCategoryData)}
          
          {/* Обычные папки и документы только если нет структурированных данных */}
          {!structuredData && !currentModule && !currentCategory && (
            <Grid container spacing={2}>
            {/* Folders */}
            {folders.map((folder) => (
              <Grid item xs={12} sm={6} md={4} lg={3} key={folder.id}>
                <Card
                  sx={{
                    '&:hover': { boxShadow: 3 },
                  }}
                >
                  <CardContent
                    sx={{ cursor: 'pointer' }}
                    onClick={() => handleFolderClick(folder)}
                  >
                    <Box display="flex" alignItems="center" mb={1}>
                      <FolderOpen color="primary" sx={{ mr: 1 }} />
                      <Typography variant="h6" noWrap sx={{ flex: 1 }}>
                        {folder.name}
                      </Typography>
                    </Box>
                    {folder.description && (
                      <Typography variant="body2" color="text.secondary" noWrap>
                        {folder.description}
                      </Typography>
                    )}
                  </CardContent>
                  <CardActions>
                    <Tooltip title="Удалить папку">
                      <IconButton
                        size="small"
                        color="error"
                        onClick={(e) => {
                          e.stopPropagation()
                          setSelectedItem(folder)
                          setOpenDelete(true)
                        }}
                      >
                        <Delete />
                      </IconButton>
                    </Tooltip>
                  </CardActions>
                </Card>
              </Grid>
            ))}

            {/* Documents */}
            {documents.map((document) => (
            <Grid item xs={12} sm={6} md={4} lg={3} key={document.id}>
              <Card
                sx={{
                  '&:hover': { boxShadow: 3 },
                }}
              >
                <CardContent>
                  <Box display="flex" alignItems="center" mb={1}>
                    {getFileIconComponent(document.mime_type)}
                    <Typography variant="h6" noWrap sx={{ ml: 1, flex: 1 }}>
                      {document.name}
                    </Typography>
                  </Box>
                  <Typography variant="body2" color="text.secondary" mb={1}>
                    {getMimeTypeLabel(document.mime_type)} • {formatFileSize(document.file_size)}
                  </Typography>
                  {document.description && (
                    <Typography variant="body2" color="text.secondary" noWrap>
                      {document.description}
                    </Typography>
                  )}
                  {document.tags.length > 0 && (
                    <Box mt={1}>
                      {document.tags.slice(0, 2).map((tag) => (
                        <Chip
                          key={tag}
                          label={tag}
                          size="small"
                          sx={{ mr: 0.5, mb: 0.5 }}
                        />
                      ))}
                      {document.tags.length > 2 && (
                        <Chip
                          label={`+${document.tags.length - 2}`}
                          size="small"
                          variant="outlined"
                        />
                      )}
                    </Box>
                  )}
                </CardContent>
                <CardActions>
                  <Tooltip title="Скачать">
                    <IconButton
                      size="small"
                      onClick={() => handleDownloadDocument(document)}
                    >
                      <Download />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Просмотр">
                    <IconButton size="small">
                      <Visibility />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Редактировать">
                    <IconButton size="small">
                      <Edit />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Удалить">
                    <IconButton
                      size="small"
                      color="error"
                      onClick={() => {
                        setSelectedItem(document)
                        setOpenDelete(true)
                      }}
                    >
                      <Delete />
                    </IconButton>
                  </Tooltip>
                </CardActions>
              </Card>
            </Grid>
          ))}

            {folders.length === 0 && documents.length === 0 && !loading && (
              <Grid item xs={12}>
                <Box textAlign="center" p={4}>
                  <Typography variant="h6" color="text.secondary">
                    Папка пуста
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Создайте папку или загрузите файл
                  </Typography>
                </Box>
              </Grid>
            )}
          </Grid>
          )}
        </>
      )}
    </Box>
  )

  const renderSearch = () => (
    <Box>
      <Box display="flex" gap={2} mb={3}>
        <TextField
          fullWidth
          placeholder="Поиск документов..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          InputProps={{
            startAdornment: <Search sx={{ mr: 1, color: 'text.secondary' }} />
          }}
        />
        <Button
          variant="contained"
          onClick={handleSearch}
          disabled={!searchTerm.trim()}
        >
          Найти
        </Button>
      </Box>

      {loading ? (
        <Box display="flex" justifyContent="center" p={4}>
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={2}>
          {documents.map((document) => (
            <Grid item xs={12} sm={6} md={4} lg={3} key={document.id}>
              <Card>
                <CardContent>
                  <Box display="flex" alignItems="center" mb={1}>
                    {getFileIconComponent(document.mime_type)}
                    <Typography variant="h6" noWrap sx={{ ml: 1, flex: 1 }}>
                      {document.name}
                    </Typography>
                  </Box>
                  <Typography variant="body2" color="text.secondary" mb={1}>
                    {getMimeTypeLabel(document.mime_type)} • {formatFileSize(document.file_size)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {new Date(document.created_at).toLocaleDateString()}
                  </Typography>
                </CardContent>
                <CardActions>
                  <IconButton
                    size="small"
                    onClick={() => handleDownloadDocument(document)}
                  >
                    <Download />
                  </IconButton>
                  <IconButton size="small">
                    <Visibility />
                  </IconButton>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}
    </Box>
  )

  const renderStats = () => (
    <Box>
      {stats ? (
        <Grid container spacing={3}>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center">
                  <Description color="primary" sx={{ mr: 2 }} />
                  <Box>
                    <Typography variant="h4">{stats.total_documents}</Typography>
                    <Typography color="text.secondary">Документов</Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center">
                  <Folder color="primary" sx={{ mr: 2 }} />
                  <Box>
                    <Typography variant="h4">{stats.total_folders}</Typography>
                    <Typography color="text.secondary">Папок</Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center">
                  <CloudUpload color="primary" sx={{ mr: 2 }} />
                  <Box>
                    <Typography variant="h4">{formatFileSize(stats.total_size)}</Typography>
                    <Typography color="text.secondary">Общий размер</Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center">
                  <Description color="primary" sx={{ mr: 2 }} />
                  <Box>
                    <Typography variant="h4">{formatFileSize(stats.storage_usage)}</Typography>
                    <Typography color="text.secondary">Использовано</Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      ) : (
        <Box display="flex" justifyContent="center" p={4}>
          <CircularProgress />
        </Box>
      )}
    </Box>
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
        Файловое хранилище
      </Typography>

      <Paper sx={{ mb: 3 }}>
        <Tabs
          value={currentTab}
          onChange={(_, newValue) => setCurrentTab(newValue)}
          aria-label="file documents tabs"
        >
          <Tab label="Документы" />
          <Tab label="Поиск" />
          <Tab label="Статистика" />
        </Tabs>
      </Paper>

      <TabPanel value={currentTab} index={0}>
        {renderDocumentsList()}
      </TabPanel>

      <TabPanel value={currentTab} index={1}>
        {renderSearch()}
      </TabPanel>

      <TabPanel value={currentTab} index={2}>
        {renderStats()}
      </TabPanel>

      {/* Create Folder Dialog */}
      <Dialog open={openCreateFolder} onClose={() => setOpenCreateFolder(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Создать папку</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Название папки"
              value={folderName}
              onChange={(e) => setFolderName(e.target.value)}
              fullWidth
              required
            />
            <TextField
              label="Описание"
              value={folderDescription}
              onChange={(e) => setFolderDescription(e.target.value)}
              fullWidth
              multiline
              rows={3}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenCreateFolder(false)}>Отмена</Button>
          <Button
            onClick={handleCreateFolder}
            variant="contained"
            disabled={!folderName.trim()}
          >
            Создать
          </Button>
        </DialogActions>
      </Dialog>

      {/* Upload Document Dialog */}
      <Dialog open={openUpload} onClose={() => setOpenUpload(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Загрузить документ</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <Button
              variant="outlined"
              component="label"
              startIcon={<CloudUpload />}
              fullWidth
              sx={{ py: 2 }}
            >
              {uploadFile ? uploadFile.name : 'Выберите файл'}
              <input
                type="file"
                hidden
                onChange={(e) => {
                  const file = e.target.files?.[0]
                  if (file) {
                    setUploadFile(file)
                    setDocumentName(file.name)
                  }
                }}
              />
            </Button>
            <TextField
              label="Название документа"
              value={documentName}
              onChange={(e) => setDocumentName(e.target.value)}
              fullWidth
              required
            />
            <TextField
              label="Описание"
              value={documentDescription}
              onChange={(e) => setDocumentDescription(e.target.value)}
              fullWidth
              multiline
              rows={3}
            />
            <Box>
              <TextField
                label="Теги"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyPress={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault()
                    addTag()
                  }
                }}
                fullWidth
                placeholder="Введите тег и нажмите Enter"
              />
              <Box mt={1} display="flex" flexWrap="wrap" gap={0.5}>
                {documentTags.map((tag) => (
                  <Chip
                    key={tag}
                    label={tag}
                    onDelete={() => removeTag(tag)}
                    size="small"
                  />
                ))}
              </Box>
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenUpload(false)}>Отмена</Button>
          <Button
            onClick={handleUploadDocument}
            variant="contained"
            disabled={!uploadFile || !documentName.trim()}
          >
            Загрузить
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={openDelete} onClose={() => setOpenDelete(false)}>
        <DialogTitle>Удалить</DialogTitle>
        <DialogContent>
          <Typography>
            Вы уверены, что хотите удалить {selectedItem && 'file_path' in selectedItem ? 'документ' : 'папку'} "{selectedItem?.name}"?
            Это действие нельзя отменить.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDelete(false)}>Отмена</Button>
          <Button onClick={handleDeleteItem} variant="contained" color="error">
            Удалить
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  )
}

export default FileDocumentsPage
