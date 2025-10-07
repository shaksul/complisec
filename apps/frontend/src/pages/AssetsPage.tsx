import React, { useState, useEffect } from 'react';
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
  LinearProgress,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
  IconButton,
  Tooltip,
  TableSortLabel,
} from '@mui/material';
import {
  Add,
  Computer,
  Edit,
  Search,
  FilterList,
  Download,
  Clear,
  Delete,
  Visibility
} from '@mui/icons-material';
import {
  assetsApi,
  Asset,
  AssetListParams,
  ASSET_TYPES,
  ASSET_CLASSES,
  CRITICALITY_LEVELS,
  ASSET_STATUSES,
  PaginationMeta,
} from '../shared/api/assets';
import { usersApi, UserCatalog } from '../shared/api/users';
import { useAuth } from '../contexts/AuthContext';
import Pagination from '../components/Pagination';
import AssetModal from '../components/assets/AssetModal';
import AssetDetailsModal from '../components/assets/AssetDetailsModal';
import BulkOperationsModal from '../components/assets/BulkOperationsModal';

interface AssetFilters {
  type: string;
  class: string;
  status: string;
  criticality: string;
  owner_id: string;
  search: string;
}

type SortField = 'name' | 'created_at' | 'type' | 'criticality' | 'status'
type SortDirection = 'asc' | 'desc'

export const AssetsPage: React.FC = () => {
  const [assets, setAssets] = useState<Asset[]>([]);
  const [users, setUsers] = useState<UserCatalog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState<PaginationMeta>({
    page: 1,
    page_size: 20,
    total: 0,
    total_pages: 0,
    has_next: false,
    has_prev: false,
  });
  const [filters, setFilters] = useState<AssetFilters>({
    type: '',
    class: '',
    status: '',
    criticality: '',
    owner_id: '',
    search: ''
  });
  const [sortField, setSortField] = useState<SortField>('created_at');
  const [sortDirection, setSortDirection] = useState<SortDirection>('desc');
  const [showFilters, setShowFilters] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingAsset, setEditingAsset] = useState<Asset | null>(null);
  const [viewingAssetId, setViewingAssetId] = useState<string | null>(null);
  const [selectedAssets, setSelectedAssets] = useState<Asset[]>([]);
  const [showBulkModal, setShowBulkModal] = useState(false);
  const { user } = useAuth();

  useEffect(() => {
    if (user) {
      void loadAssets();
    } else {
      setLoading(false);
    }
  }, [filters, pagination.page, pagination.page_size, user]);

  useEffect(() => {
    if (user) {
      void loadUsers();
    }
  }, [user]);

  useEffect(() => {
    setSelectedAssets([]);
  }, [filters, pagination.page]);

  const loadAssets = async () => {
    try {
      setLoading(true);
      setError(null);
      const params: AssetListParams = {
        page: pagination.page,
        page_size: pagination.page_size,
        ...filters
      };
      
      const response = await assetsApi.list(params);
      setAssets(response.data ?? []);
      if (response.pagination) {
        setPagination(response.pagination);
      }
    } catch (err) {
      setError('Ошибка загрузки активов');
      console.error('Error loading assets:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadUsers = async () => {
    try {
      const response = await usersApi.getUserCatalog({ page: 1, page_size: 100 });
      setUsers(response.data || []);
    } catch (err) {
      console.error('Error loading users:', err);
      setUsers([]);
    }
  };

  const handleFilterChange = (key: keyof AssetFilters, value: string) => {
    setFilters(prev => ({ ...prev, [key]: value }));
    setPagination(prev => ({ ...prev, page: 1 }));
  };

  const handlePageChange = (page: number) => {
    if (page < 1 || page === pagination.page || page > pagination.total_pages) {
      return;
    }
    setPagination(prev => ({ ...prev, page }));
    setSelectedAssets([]);
  };

  const handleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection(prev => prev === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortDirection('desc');
    }
  };

  const handleCreateAsset = async (assetData: any) => {
    try {
      await assetsApi.create(assetData);
      setShowCreateModal(false);
      loadAssets();
    } catch (err) {
      setError('Ошибка создания актива');
      console.error('Error creating asset:', err);
    }
  };

  const handleUpdateAsset = async (id: string, assetData: any) => {
    try {
      await assetsApi.update(id, assetData);
      setEditingAsset(null);
      loadAssets();
    } catch (err) {
      setError('Ошибка обновления актива');
      console.error('Error updating asset:', err);
    }
  };

  const handleDeleteAsset = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этот актив?')) {
      try {
        await assetsApi.delete(id);
        loadAssets();
      } catch (err) {
        setError('Ошибка удаления актива');
        console.error('Error deleting asset:', err);
      }
    }
  };

  const handleViewAsset = (id: string) => {
    setViewingAssetId(id);
  };

  const handleEditAsset = (asset: Asset) => {
    setEditingAsset(asset);
  };

  const handleSelectAsset = (asset: Asset, selected: boolean) => {
    if (selected) {
      setSelectedAssets(prev => [...prev, asset]);
    } else {
      setSelectedAssets(prev => prev.filter(a => a.id !== asset.id));
    }
  };

  const handleSelectAll = (selected: boolean) => {
    if (selected) {
      setSelectedAssets([...assets]);
    } else {
      setSelectedAssets([]);
    }
  };

  const handleBulkOperation = () => {
    if (selectedAssets.length === 0) {
      setError('Выберите активы для массовой операции');
      return;
    }
    setShowBulkModal(true);
  };

  const handleBulkSuccess = () => {
    setSelectedAssets([]);
    loadAssets();
  };

  const handleExport = async () => {
    try {
      const blob = await assetsApi.export(filters);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'assets.csv';
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (err) {
      setError('Ошибка экспорта');
      console.error('Error exporting assets:', err);
    }
  };

  const getCriticalityColor = (criticality: string) => {
    switch (criticality) {
      case 'high': return 'error'
      case 'medium': return 'warning'
      case 'low': return 'success'
      default: return 'default'
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'success'
      case 'in_repair': return 'warning'
      case 'storage': return 'info'
      case 'decommissioned': return 'error'
      default: return 'default'
    }
  };

  const clearFilters = () => {
    setFilters({
      type: '',
      class: '',
      status: '',
      criticality: '',
      owner_id: '',
      search: ''
    });
  };

  const getTypeLabel = (type: string) => {
    return ASSET_TYPES?.find(t => t.value === type)?.label || type;
  };

  const getStatusLabel = (status: string) => {
    return ASSET_STATUSES?.find(s => s.value === status)?.label || status;
  };

  const getCriticalityLabel = (criticality: string) => {
    return CRITICALITY_LEVELS?.find(c => c.value === criticality)?.label || criticality;
  };

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Активы</Typography>
        <Box display="flex" gap={1}>
          <Button 
            variant="outlined" 
            startIcon={<FilterList />}
            onClick={() => setShowFilters(!showFilters)}
          >
            Фильтры
          </Button>
          <Button 
            variant="outlined" 
            startIcon={<Download />}
            onClick={handleExport}
          >
            Экспорт
          </Button>
          <Button 
            variant="contained" 
            startIcon={<Add />} 
            onClick={() => setShowCreateModal(true)}
          >
            Добавить актив
          </Button>
        </Box>
      </Box>

      {error && (
        <Box mb={2} p={2} bgcolor="error.light" borderRadius={1}>
          <Typography color="error">{error}</Typography>
        </Box>
      )}

      {/* Search and Filters */}
      <Paper sx={{ mb: 2, p: 2 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={6}>
            <TextField
              fullWidth
              placeholder="Поиск по названию или инв. номеру..."
              value={filters.search}
              onChange={(e) => handleFilterChange('search', e.target.value)}
              InputProps={{
                startAdornment: <Search sx={{ mr: 1, color: 'text.secondary' }} />
              }}
            />
          </Grid>
          <Grid item xs={12} md={2}>
            <Button
              variant="outlined"
              startIcon={<Clear />}
              onClick={clearFilters}
              fullWidth
            >
              Очистить
            </Button>
          </Grid>
        </Grid>

        {showFilters && (
          <Grid container spacing={2} sx={{ mt: 2 }}>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Тип</InputLabel>
                <Select
                value={filters.type}
                onChange={(e) => handleFilterChange('type', e.target.value)}
                  label="Тип"
              >
                  <MenuItem value="">Все</MenuItem>
                {ASSET_TYPES?.map(type => (
                    <MenuItem key={type.value} value={type.value}>
                      {type.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Класс</InputLabel>
                <Select
                value={filters.class}
                onChange={(e) => handleFilterChange('class', e.target.value)}
                  label="Класс"
              >
                  <MenuItem value="">Все</MenuItem>
                {ASSET_CLASSES?.map(cls => (
                    <MenuItem key={cls.value} value={cls.value}>
                      {cls.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Статус</InputLabel>
                <Select
                value={filters.status}
                onChange={(e) => handleFilterChange('status', e.target.value)}
                  label="Статус"
              >
                  <MenuItem value="">Все</MenuItem>
                {ASSET_STATUSES?.map(status => (
                    <MenuItem key={status.value} value={status.value}>
                      {status.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Критичность</InputLabel>
                <Select
                value={filters.criticality}
                onChange={(e) => handleFilterChange('criticality', e.target.value)}
                  label="Критичность"
              >
                  <MenuItem value="">Все</MenuItem>
                {CRITICALITY_LEVELS?.map(level => (
                    <MenuItem key={level.value} value={level.value}>
                      {level.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        )}
      </Paper>

      {/* Bulk Actions */}
      {selectedAssets.length > 0 && (
        <Paper sx={{ mb: 2, p: 2 }}>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="body2">
              Выбрано активов: {selectedAssets.length}
            </Typography>
            <Button
              variant="contained"
              color="secondary"
              onClick={handleBulkOperation}
            >
              Массовые операции
            </Button>
          </Box>
        </Paper>
      )}

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell padding="checkbox">
                  <input
                    type="checkbox"
                    checked={selectedAssets.length === assets.length && assets.length > 0}
                    onChange={(e) => handleSelectAll(e.target.checked)}
                    style={{ transform: 'scale(1.2)' }}
                  />
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={sortField === 'name'}
                    direction={sortField === 'name' ? sortDirection : 'asc'}
                    onClick={() => handleSort('name')}
                  >
                  Название
                  </TableSortLabel>
                </TableCell>
                <TableCell>Тип</TableCell>
                <TableCell>Владелец</TableCell>
                <TableCell>Ответственный</TableCell>
                <TableCell>
                  <TableSortLabel
                    active={sortField === 'criticality'}
                    direction={sortField === 'criticality' ? sortDirection : 'asc'}
                    onClick={() => handleSort('criticality')}
                  >
                  Критичность
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={sortField === 'status'}
                    direction={sortField === 'status' ? sortDirection : 'asc'}
                    onClick={() => handleSort('status')}
                  >
                  Статус
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={sortField === 'created_at'}
                    direction={sortField === 'created_at' ? sortDirection : 'asc'}
                    onClick={() => handleSort('created_at')}
                  >
                    Дата создания
                  </TableSortLabel>
                </TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={9} align="center">
                    <LinearProgress />
                    <Typography sx={{ mt: 1 }}>Загрузка активов...</Typography>
                  </TableCell>
                </TableRow>
              ) : assets.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={9} align="center">
                    <Typography>Нет активов для отображения.</Typography>
                  </TableCell>
                </TableRow>
              ) : (
                assets.map((asset) => (
                  <TableRow key={asset.id} hover>
                    <TableCell padding="checkbox">
                    <input
                      type="checkbox"
                      checked={selectedAssets.some(a => a.id === asset.id)}
                      onChange={(e) => handleSelectAsset(asset, e.target.checked)}
                        style={{ transform: 'scale(1.2)' }}
                      />
                    </TableCell>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <Computer sx={{ mr: 1 }} />
                        <Box>
                          <Typography variant="body2" fontWeight="medium">
                      {asset.name}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            Инв. №: {asset.inventory_number}
                          </Typography>
                        </Box>
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {getTypeLabel(asset.type)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                    {asset.owner_name || 'Не назначен'}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                    {asset.responsible_user_name || 'Не назначен'}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={getCriticalityLabel(asset.criticality)}
                        color={getCriticalityColor(asset.criticality) as any}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={getStatusLabel(asset.status)}
                        color={getStatusColor(asset.status) as any}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {new Date(asset.created_at).toLocaleDateString('ru-RU')}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Box display="flex" gap={0.5}>
                        <Tooltip title="Просмотр деталей">
                          <IconButton 
                            size="small"
                            color="primary"
                        onClick={() => handleViewAsset(asset.id)}
                          >
                            <Visibility />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Редактировать актив">
                          <IconButton 
                            size="small"
                            color="primary"
                            onClick={() => handleEditAsset(asset)}
                          >
                            <Edit />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Удалить актив">
                          <IconButton 
                            size="small"
                            color="error"
                        onClick={() => handleDeleteAsset(asset.id)}
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

      {/* Pagination */}
      <Box mt={2}>
        <Pagination
          currentPage={pagination.page}
          totalPages={pagination.total_pages}
          hasNext={pagination.has_next}
          hasPrev={pagination.has_prev}
          onPageChange={handlePageChange}
        />
      </Box>

      {/* Create/Edit Modal */}
      <AssetModal
        open={showCreateModal || !!editingAsset}
        asset={editingAsset}
        users={users}
        onSave={editingAsset ? 
          (data) => handleUpdateAsset(editingAsset.id, data) : 
          handleCreateAsset
        }
        onClose={() => {
          setShowCreateModal(false);
          setEditingAsset(null);
        }}
      />

      {/* Asset Details Modal */}
      {viewingAssetId && (
        <AssetDetailsModal
          assetId={viewingAssetId}
          onClose={() => setViewingAssetId(null)}
        />
      )}

      {/* Bulk Operations Modal */}
      {showBulkModal && (
        <BulkOperationsModal
          selectedAssets={selectedAssets}
          onClose={() => setShowBulkModal(false)}
          onSuccess={handleBulkSuccess}
        />
      )}
    </Container>
  );
};
