import React, { useState, useEffect } from 'react';
import {
  Box,
  Tabs,
  Tab,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  CardActions,
  Chip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Switch,
  FormControlLabel,
  Alert,
  CircularProgress,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  PlayArrow as PlayIcon,
  Assignment as AssignmentIcon,
  School as SchoolIcon,
  Quiz as QuizIcon,
  VideoLibrary as VideoIcon,
  Description as DocumentIcon,
  CheckCircle as AcknowledgmentIcon,
} from '@mui/icons-material';
import { materialsApi, coursesApi, MATERIAL_TYPES, MATERIAL_SOURCES, Material, TrainingCourse, CreateMaterialRequest, CreateCourseRequest } from '../shared/api/training';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`training-tabpanel-${index}`}
      aria-labelledby={`training-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

const TrainingPage: React.FC = () => {
  const [tabValue, setTabValue] = useState(0);
  const [materials, setMaterials] = useState<Material[]>([]);
  const [courses, setCourses] = useState<TrainingCourse[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Material dialog state
  const [materialDialogOpen, setMaterialDialogOpen] = useState(false);
  const [editingMaterial, setEditingMaterial] = useState<Material | null>(null);
  const [materialForm, setMaterialForm] = useState<CreateMaterialRequest>({
    title: '',
    description: '',
    uri: '',
    type: 'file',
    material_type: 'document',
    duration_minutes: undefined,
    tags: [],
    is_required: false,
    passing_score: 80,
    attempts_limit: undefined,
    metadata: {},
  });

  // Course dialog state
  const [courseDialogOpen, setCourseDialogOpen] = useState(false);
  const [editingCourse, setEditingCourse] = useState<TrainingCourse | null>(null);
  const [courseForm, setCourseForm] = useState<CreateCourseRequest>({
    title: '',
    description: '',
    is_active: true,
  });

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const loadMaterials = async () => {
    try {
      setLoading(true);
      const response = await materialsApi.list();
      setMaterials(response.items);
    } catch (err) {
      setError('Ошибка загрузки материалов');
    } finally {
      setLoading(false);
    }
  };

  const loadCourses = async () => {
    try {
      setLoading(true);
      const response = await coursesApi.list();
      setCourses(response.items);
    } catch (err) {
      setError('Ошибка загрузки курсов');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadMaterials();
    loadCourses();
  }, []);

  // Material handlers
  const handleCreateMaterial = () => {
    setEditingMaterial(null);
    setMaterialForm({
      title: '',
      description: '',
      uri: '',
      type: 'file',
      material_type: 'document',
      duration_minutes: undefined,
      tags: [],
      is_required: false,
      passing_score: 80,
      attempts_limit: undefined,
      metadata: {},
    });
    setMaterialDialogOpen(true);
  };

  const handleEditMaterial = (material: Material) => {
    setEditingMaterial(material);
    setMaterialForm({
      title: material.title,
      description: material.description || '',
      uri: material.uri,
      type: material.type,
      material_type: material.material_type,
      duration_minutes: material.duration_minutes,
      tags: material.tags,
      is_required: material.is_required,
      passing_score: material.passing_score,
      attempts_limit: material.attempts_limit,
      metadata: material.metadata,
    });
    setMaterialDialogOpen(true);
  };

  const handleSaveMaterial = async () => {
    try {
      if (editingMaterial) {
        await materialsApi.update(editingMaterial.id, materialForm);
      } else {
        await materialsApi.create(materialForm);
      }
      setMaterialDialogOpen(false);
      loadMaterials();
    } catch (err) {
      setError('Ошибка сохранения материала');
    }
  };

  const handleDeleteMaterial = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этот материал?')) {
      try {
        await materialsApi.delete(id);
        loadMaterials();
      } catch (err) {
        setError('Ошибка удаления материала');
      }
    }
  };

  // Course handlers
  const handleCreateCourse = () => {
    setEditingCourse(null);
    setCourseForm({
      title: '',
      description: '',
      is_active: true,
    });
    setCourseDialogOpen(true);
  };

  const handleEditCourse = (course: TrainingCourse) => {
    setEditingCourse(course);
    setCourseForm({
      title: course.title,
      description: course.description || '',
      is_active: course.is_active,
    });
    setCourseDialogOpen(true);
  };

  const handleSaveCourse = async () => {
    try {
      if (editingCourse) {
        await coursesApi.update(editingCourse.id, courseForm);
      } else {
        await coursesApi.create(courseForm);
      }
      setCourseDialogOpen(false);
      loadCourses();
    } catch (err) {
      setError('Ошибка сохранения курса');
    }
  };

  const handleDeleteCourse = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этот курс?')) {
      try {
        await coursesApi.delete(id);
        loadCourses();
      } catch (err) {
        setError('Ошибка удаления курса');
      }
    }
  };

  const getMaterialIcon = (materialType: string) => {
    switch (materialType) {
      case 'video':
        return <VideoIcon />;
      case 'quiz':
        return <QuizIcon />;
      case 'document':
        return <DocumentIcon />;
      case 'acknowledgment':
        return <AcknowledgmentIcon />;
      default:
        return <AssignmentIcon />;
    }
  };

  const getMaterialColor = (materialType: string) => {
    switch (materialType) {
      case 'video':
        return 'primary';
      case 'quiz':
        return 'secondary';
      case 'document':
        return 'success';
      case 'acknowledgment':
        return 'warning';
      default:
        return 'default';
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Управление обучением</Typography>
        <Box>
          {tabValue === 0 && (
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={handleCreateMaterial}
            >
              Добавить материал
            </Button>
          )}
          {tabValue === 1 && (
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={handleCreateCourse}
            >
              Создать курс
            </Button>
          )}
        </Box>
      </Box>

      {error && (
        <Alert severity="error" onClose={() => setError(null)} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={handleTabChange}>
          <Tab label="Материалы" icon={<AssignmentIcon />} iconPosition="start" />
          <Tab label="Курсы" icon={<SchoolIcon />} iconPosition="start" />
          <Tab label="Назначения" icon={<PlayIcon />} iconPosition="start" />
        </Tabs>
      </Box>

      <TabPanel value={tabValue} index={0}>
        <Grid container spacing={3}>
          {materials.map((material) => (
            <Grid item xs={12} sm={6} md={4} key={material.id}>
              <Card>
                <CardContent>
                  <Box display="flex" alignItems="center" mb={2}>
                    <Box mr={2} color={`${getMaterialColor(material.material_type)}.main`}>
                      {getMaterialIcon(material.material_type)}
                    </Box>
                    <Box flex={1}>
                      <Typography variant="h6" component="h3">
                        {material.title}
                      </Typography>
                      <Chip
                        label={MATERIAL_TYPES[material.material_type]}
                        size="small"
                        color={getMaterialColor(material.material_type)}
                      />
                    </Box>
                  </Box>
                  
                  {material.description && (
                    <Typography variant="body2" color="text.secondary" mb={2}>
                      {material.description}
                    </Typography>
                  )}
                  
                  <Box display="flex" gap={1} flexWrap="wrap">
                    <Chip label={MATERIAL_SOURCES[material.type]} size="small" />
                    {material.duration_minutes && (
                      <Chip label={`${material.duration_minutes} мин`} size="small" />
                    )}
                    {material.is_required && (
                      <Chip label="Обязательный" size="small" color="error" />
                    )}
                  </Box>
                </CardContent>
                
                <CardActions>
                  <IconButton onClick={() => handleEditMaterial(material)}>
                    <EditIcon />
                  </IconButton>
                  <IconButton onClick={() => handleDeleteMaterial(material.id)}>
                    <DeleteIcon />
                  </IconButton>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        <Grid container spacing={3}>
          {courses.map((course) => (
            <Grid item xs={12} sm={6} md={4} key={course.id}>
              <Card>
                <CardContent>
                  <Box display="flex" alignItems="center" mb={2}>
                    <SchoolIcon color="primary" sx={{ mr: 2 }} />
                    <Box flex={1}>
                      <Typography variant="h6" component="h3">
                        {course.title}
                      </Typography>
                      <Chip
                        label={course.is_active ? 'Активен' : 'Неактивен'}
                        size="small"
                        color={course.is_active ? 'success' : 'default'}
                      />
                    </Box>
                  </Box>
                  
                  {course.description && (
                    <Typography variant="body2" color="text.secondary" mb={2}>
                      {course.description}
                    </Typography>
                  )}
                  
                  {course.materials && (
                    <Typography variant="body2">
                      Материалов: {course.materials.length}
                    </Typography>
                  )}
                </CardContent>
                
                <CardActions>
                  <IconButton onClick={() => handleEditCourse(course)}>
                    <EditIcon />
                  </IconButton>
                  <IconButton onClick={() => handleDeleteCourse(course.id)}>
                    <DeleteIcon />
                  </IconButton>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      </TabPanel>

      <TabPanel value={tabValue} index={2}>
        <Typography variant="h6" color="text.secondary">
          Функционал назначений будет добавлен в следующих версиях
        </Typography>
      </TabPanel>

      {/* Material Dialog */}
      <Dialog open={materialDialogOpen} onClose={() => setMaterialDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          {editingMaterial ? 'Редактировать материал' : 'Создать материал'}
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <TextField
                label="Название"
                value={materialForm.title}
                onChange={(e) => setMaterialForm({ ...materialForm, title: e.target.value })}
                fullWidth
                required
              />
            </Grid>
            
            <Grid item xs={12}>
              <TextField
                label="Описание"
                value={materialForm.description}
                onChange={(e) => setMaterialForm({ ...materialForm, description: e.target.value })}
                fullWidth
                multiline
                rows={3}
              />
            </Grid>
            
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Тип материала</InputLabel>
                <Select
                  value={materialForm.material_type}
                  onChange={(e) => setMaterialForm({ ...materialForm, material_type: e.target.value as any })}
                >
                  <MenuItem value="document">Документ</MenuItem>
                  <MenuItem value="video">Видео</MenuItem>
                  <MenuItem value="quiz">Квиз</MenuItem>
                  <MenuItem value="simulation">Симуляция</MenuItem>
                  <MenuItem value="acknowledgment">Ознакомление</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Источник</InputLabel>
                <Select
                  value={materialForm.type}
                  onChange={(e) => setMaterialForm({ ...materialForm, type: e.target.value as any })}
                >
                  <MenuItem value="file">Файл</MenuItem>
                  <MenuItem value="link">Ссылка</MenuItem>
                  <MenuItem value="video">Видео</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            
            <Grid item xs={12}>
              <TextField
                label="URI/Путь"
                value={materialForm.uri}
                onChange={(e) => setMaterialForm({ ...materialForm, uri: e.target.value })}
                fullWidth
                required
                helperText="Путь к файлу или URL ссылка"
              />
            </Grid>
            
            <Grid item xs={12} sm={6}>
              <TextField
                label="Длительность (минуты)"
                type="number"
                value={materialForm.duration_minutes || ''}
                onChange={(e) => setMaterialForm({ 
                  ...materialForm, 
                  duration_minutes: e.target.value ? parseInt(e.target.value) : undefined 
                })}
                fullWidth
              />
            </Grid>
            
            <Grid item xs={12} sm={6}>
              <TextField
                label="Проходной балл"
                type="number"
                value={materialForm.passing_score}
                onChange={(e) => setMaterialForm({ 
                  ...materialForm, 
                  passing_score: parseInt(e.target.value) 
                })}
                fullWidth
                inputProps={{ min: 0, max: 100 }}
              />
            </Grid>
            
            <Grid item xs={12}>
              <FormControlLabel
                control={
                  <Switch
                    checked={materialForm.is_required}
                    onChange={(e) => setMaterialForm({ ...materialForm, is_required: e.target.checked })}
                  />
                }
                label="Обязательный материал"
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setMaterialDialogOpen(false)}>Отмена</Button>
          <Button onClick={handleSaveMaterial} variant="contained">
            {editingMaterial ? 'Обновить' : 'Создать'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Course Dialog */}
      <Dialog open={courseDialogOpen} onClose={() => setCourseDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {editingCourse ? 'Редактировать курс' : 'Создать курс'}
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <TextField
                label="Название курса"
                value={courseForm.title}
                onChange={(e) => setCourseForm({ ...courseForm, title: e.target.value })}
                fullWidth
                required
              />
            </Grid>
            
            <Grid item xs={12}>
              <TextField
                label="Описание"
                value={courseForm.description}
                onChange={(e) => setCourseForm({ ...courseForm, description: e.target.value })}
                fullWidth
                multiline
                rows={3}
              />
            </Grid>
            
            <Grid item xs={12}>
              <FormControlLabel
                control={
                  <Switch
                    checked={courseForm.is_active}
                    onChange={(e) => setCourseForm({ ...courseForm, is_active: e.target.checked })}
                  />
                }
                label="Активный курс"
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCourseDialogOpen(false)}>Отмена</Button>
          <Button onClick={handleSaveCourse} variant="contained">
            {editingCourse ? 'Обновить' : 'Создать'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default TrainingPage;