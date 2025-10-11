import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  CircularProgress,
  Alert,
  Chip,
  Button,
} from '@mui/material';
import {
  Warning as WarningIcon,
  Assignment as AssignmentIcon,
  Person as PersonIcon,
  Download as DownloadIcon,
} from '@mui/icons-material';
import { assetsApi, Asset } from '../shared/api/assets';

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
      id={`inventory-tabpanel-${index}`}
      aria-labelledby={`inventory-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

const AssetsInventoryPage: React.FC = () => {
  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [assetsWithoutOwner, setAssetsWithoutOwner] = useState<Asset[]>([]);
  const [assetsWithoutPassport, setAssetsWithoutPassport] = useState<Asset[]>([]);
  const [assetsWithoutCriticality, setAssetsWithoutCriticality] = useState<Asset[]>([]);

  useEffect(() => {
    loadInventoryData();
  }, []);

  const loadInventoryData = async () => {
    try {
      setLoading(true);
      setError(null);

      const [withoutOwner, withoutPassport, withoutCriticality] = await Promise.all([
        assetsApi.getAssetsWithoutOwner(),
        assetsApi.getAssetsWithoutPassport(),
        assetsApi.getAssetsWithoutCriticality(),
      ]);

      setAssetsWithoutOwner(withoutOwner);
      setAssetsWithoutPassport(withoutPassport);
      setAssetsWithoutCriticality(withoutCriticality);
    } catch (err: any) {
      setError(err.message || 'Ошибка загрузки данных инвентаризации');
      console.error('Error loading inventory data:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const exportToCSV = (data: Asset[], filename: string) => {
    const headers = ['Инв. номер', 'Название', 'Тип', 'Класс', 'Статус', 'Критичность'];
    const rows = data.map(asset => [
      asset.inventory_number,
      asset.name,
      asset.type,
      asset.class,
      asset.status,
      asset.criticality,
    ]);

    const csvContent = [
      headers.join(','),
      ...rows.map(row => row.join(',')),
    ].join('\n');

    const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', `${filename}.csv`);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'success';
      case 'in_repair':
        return 'warning';
      case 'storage':
        return 'info';
      case 'decommissioned':
        return 'error';
      default:
        return 'default';
    }
  };

  const getCriticalityColor = (criticality: string) => {
    switch (criticality) {
      case 'high':
        return 'error';
      case 'medium':
        return 'warning';
      case 'low':
        return 'success';
      default:
        return 'default';
    }
  };

  const renderAssetTable = (assets: Asset[], emptyMessage: string, exportFilename: string) => {
    if (loading) {
      return (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
          <CircularProgress />
        </Box>
      );
    }

    if (assets.length === 0) {
      return (
        <Box textAlign="center" py={4}>
          <Typography color="text.secondary">{emptyMessage}</Typography>
        </Box>
      );
    }

    return (
      <>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h6">Найдено: {assets.length} активов</Typography>
          <Button
            startIcon={<DownloadIcon />}
            onClick={() => exportToCSV(assets, exportFilename)}
            variant="outlined"
            size="small"
          >
            Экспорт в CSV
          </Button>
        </Box>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Инв. номер</TableCell>
                <TableCell>Название</TableCell>
                <TableCell>Тип</TableCell>
                <TableCell>Класс</TableCell>
                <TableCell>Владелец</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Критичность</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {assets.map((asset) => (
                <TableRow key={asset.id} hover>
                  <TableCell>{asset.inventory_number}</TableCell>
                  <TableCell>{asset.name}</TableCell>
                  <TableCell sx={{ textTransform: 'capitalize' }}>{asset.type}</TableCell>
                  <TableCell sx={{ textTransform: 'capitalize' }}>{asset.class}</TableCell>
                  <TableCell>{asset.owner_name || '-'}</TableCell>
                  <TableCell>
                    <Chip
                      label={asset.status}
                      color={getStatusColor(asset.status) as any}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={asset.criticality}
                      color={getCriticalityColor(asset.criticality) as any}
                      size="small"
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </>
    );
  };

  return (
    <Box sx={{ p: 3 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Инвентаризация активов</Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={loadInventoryData}
          disabled={loading}
        >
          Обновить данные
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <Paper>
        <Tabs value={tabValue} onChange={handleTabChange} aria-label="inventory tabs">
          <Tab
            icon={<PersonIcon />}
            label={`Без владельца (${assetsWithoutOwner.length})`}
            id="inventory-tab-0"
          />
          <Tab
            icon={<AssignmentIcon />}
            label={`Без паспорта (${assetsWithoutPassport.length})`}
            id="inventory-tab-1"
          />
          <Tab
            icon={<WarningIcon />}
            label={`Без критичности (${assetsWithoutCriticality.length})`}
            id="inventory-tab-2"
          />
        </Tabs>

        <TabPanel value={tabValue} index={0}>
          <Alert severity="warning" sx={{ mb: 2 }}>
            Активы без назначенного владельца. Необходимо назначить ответственное лицо для
            каждого актива.
          </Alert>
          {renderAssetTable(
            assetsWithoutOwner,
            'Все активы имеют назначенного владельца',
            'активы_без_владельца'
          )}
        </TabPanel>

        <TabPanel value={tabValue} index={1}>
          <Alert severity="warning" sx={{ mb: 2 }}>
            Активы класса "hardware" без заполненного паспорта. Необходимо заполнить
            паспортные данные (серийный номер, модель, производитель).
          </Alert>
          {renderAssetTable(
            assetsWithoutPassport,
            'Все активы hardware имеют заполненный паспорт',
            'активы_без_паспорта'
          )}
        </TabPanel>

        <TabPanel value={tabValue} index={2}>
          <Alert severity="info" sx={{ mb: 2 }}>
            Активы без оценки критичности. Рекомендуется провести оценку CIA (Confidentiality,
            Integrity, Availability) для всех активов.
          </Alert>
          {renderAssetTable(
            assetsWithoutCriticality,
            'Все активы имеют оценку критичности',
            'активы_без_критичности'
          )}
        </TabPanel>
      </Paper>
    </Box>
  );
};

export default AssetsInventoryPage;


