import React, { useState, useEffect } from 'react';
import {
  assetsApi, 
  Asset, 
  AssetListParams, 
  ASSET_TYPES, 
  ASSET_CLASSES, 
  CRITICALITY_LEVELS, 
  ASSET_STATUSES 
} from '../shared/api/assets';
import { usersApi, User, UserCatalog } from '../shared/api/users';
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

export const AssetsPage: React.FC = () => {
  const [assets, setAssets] = useState<Asset[]>([]);
  const [users, setUsers] = useState<UserCatalog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState({
    page: 1,
    page_size: 20,
    total: 0,
    total_pages: 0
  });
  const [filters, setFilters] = useState<AssetFilters>({
    type: '',
    class: '',
    status: '',
    criticality: '',
    owner_id: '',
    search: ''
  });
  const [searchTerm, setSearchTerm] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingAsset, setEditingAsset] = useState<Asset | null>(null);
  const [viewingAssetId, setViewingAssetId] = useState<string | null>(null);
  const [selectedAssets, setSelectedAssets] = useState<Asset[]>([]);
  const [showBulkModal, setShowBulkModal] = useState(false);

  useEffect(() => {
    loadAssets();
    loadUsers();
  }, [pagination.page, filters]);

      // Debounce search
      useEffect(() => {
        console.log('Search term changed:', searchTerm);
        const timer = setTimeout(() => {
          console.log('Setting search filter:', searchTerm);
          setFilters(prev => {
            const newFilters = { 
              type: prev.type || '',
              class: prev.class || '',
              status: prev.status || '',
              criticality: prev.criticality || '',
              owner_id: prev.owner_id || '',
              search: searchTerm
            };
            console.log('New filters with search:', newFilters);
            return newFilters;
          });
        }, searchTerm === '' ? 0 : 500); // No delay for empty search

        return () => clearTimeout(timer);
      }, [searchTerm]);

  const loadAssets = async () => {
    try {
      setLoading(true);
      const params: AssetListParams = {
        page: pagination.page,
        page_size: pagination.page_size,
        ...filters
      };
      
      console.log('Loading assets with params:', params);
      console.log('Filters state:', filters);
      console.log('Search term in filters:', filters.search);
      const response = await assetsApi.list(params);
      console.log('Assets API response:', response);
      console.log('Assets data:', response.data);
      console.log('Assets pagination:', response.pagination);
      if (response.data && response.data.length > 0) {
        console.log('First asset:', response.data[0]);
      }
      setAssets(response.data);
      setPagination(response.pagination);
    } catch (err) {
      setError('Ошибка загрузки активов');
      console.error('Error loading assets:', err);
    } finally {
      console.log('Setting loading to false');
      setLoading(false);
    }
  };

  const loadUsers = async () => {
    try {
      const response = await usersApi.getUserCatalog();
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
    setPagination(prev => ({ ...prev, page }));
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

  const handleSelectAsset = (asset: Asset, selected: boolean) => {
    if (selected) {
      setSelectedAssets(prev => [...prev, asset]);
    } else {
      setSelectedAssets(prev => prev.filter(a => a.id !== asset.id));
    }
  };

  const handleSelectAll = (selected: boolean) => {
    if (selected) {
      setSelectedAssets(assets);
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
      case 'high': return 'text-red-600 bg-red-100';
      case 'medium': return 'text-yellow-600 bg-yellow-100';
      case 'low': return 'text-green-600 bg-green-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'text-green-600 bg-green-100';
      case 'in_repair': return 'text-yellow-600 bg-yellow-100';
      case 'storage': return 'text-blue-600 bg-blue-100';
      case 'decommissioned': return 'text-red-600 bg-red-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg">Загрузка...</div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">Управление активами</h1>
        
        {/* Filters */}
        <div className="bg-white p-4 rounded-lg shadow mb-4">
          <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Поиск</label>
              <input
                type="text"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                placeholder="Название или инв. номер"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Тип</label>
              <select
                value={filters.type}
                onChange={(e) => handleFilterChange('type', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Все типы</option>
                {ASSET_TYPES?.map(type => (
                  <option key={type.value} value={type.value}>{type.label}</option>
                )) || []}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Класс</label>
              <select
                value={filters.class}
                onChange={(e) => handleFilterChange('class', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Все классы</option>
                {ASSET_CLASSES?.map(cls => (
                  <option key={cls.value} value={cls.value}>{cls.label}</option>
                )) || []}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Статус</label>
              <select
                value={filters.status}
                onChange={(e) => handleFilterChange('status', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Все статусы</option>
                {ASSET_STATUSES?.map(status => (
                  <option key={status.value} value={status.value}>{status.label}</option>
                )) || []}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Критичность</label>
              <select
                value={filters.criticality}
                onChange={(e) => handleFilterChange('criticality', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Все уровни</option>
                {CRITICALITY_LEVELS?.map(level => (
                  <option key={level.value} value={level.value}>{level.label}</option>
                )) || []}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Владелец</label>
              <select
                value={filters.owner_id}
                onChange={(e) => handleFilterChange('owner_id', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Все владельцы</option>
                {users?.map(user => (
                  <option key={user.id} value={user.id}>
                    {user.first_name} {user.last_name}
                  </option>
                )) || []}
              </select>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex justify-between items-center">
          <div className="flex space-x-2">
            <button
              onClick={() => setShowCreateModal(true)}
              className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              Добавить актив
            </button>
            <button
              onClick={handleBulkOperation}
              disabled={selectedAssets.length === 0}
              className="bg-purple-600 text-white px-4 py-2 rounded-md hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Массовые операции ({selectedAssets.length})
            </button>
            <button
              onClick={handleExport}
              className="bg-green-600 text-white px-4 py-2 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500"
            >
              Экспорт
            </button>
          </div>
          
          <div className="text-sm text-gray-500">
            Всего: {pagination.total} активов
            {selectedAssets.length > 0 && ` • Выбрано: ${selectedAssets.length}`}
          </div>
        </div>
      </div>

      {/* Error message */}
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      {/* Assets table */}
      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  <input
                    type="checkbox"
                    checked={selectedAssets.length === assets.length && assets.length > 0}
                    onChange={(e) => handleSelectAll(e.target.checked)}
                    className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Инв. номер
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Название
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Тип
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Владелец
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Ответственный
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Критичность
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Статус
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Действия
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {assets?.map((asset) => (
                <tr key={asset.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <input
                      type="checkbox"
                      checked={selectedAssets.some(a => a.id === asset.id)}
                      onChange={(e) => handleSelectAsset(asset, e.target.checked)}
                      className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {asset.inventory_number}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {asset.name}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {ASSET_TYPES?.find(t => t.value === asset.type)?.label || asset.type}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {asset.owner_name || 'Организация'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {asset.responsible_user_name || 'Не назначен'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getCriticalityColor(asset.criticality)}`}>
                      {CRITICALITY_LEVELS?.find(c => c.value === asset.criticality)?.label || asset.criticality}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(asset.status)}`}>
                      {ASSET_STATUSES?.find(s => s.value === asset.status)?.label || asset.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <div className="flex space-x-2">
                      <button
                        onClick={() => handleViewAsset(asset.id)}
                        className="text-green-600 hover:text-green-900"
                        title="Просмотр деталей"
                      >
                        Просмотр
                      </button>
                      <button
                        onClick={() => setEditingAsset(asset)}
                        className="text-blue-600 hover:text-blue-900"
                        title="Редактировать"
                      >
                        Редактировать
                      </button>
                      <button
                        onClick={() => handleDeleteAsset(asset.id)}
                        className="text-red-600 hover:text-red-900"
                        title="Удалить"
                      >
                        Удалить
                      </button>
                    </div>
                  </td>
                </tr>
              )) || []}
            </tbody>
          </table>
        </div>
      </div>

      {/* Pagination */}
      <div className="mt-4">
        <Pagination
          currentPage={pagination.page}
          totalPages={pagination.total_pages}
          hasNext={pagination.page < pagination.total_pages}
          hasPrev={pagination.page > 1}
          onPageChange={handlePageChange}
        />
      </div>

      {/* Create/Edit Modal */}
      {(showCreateModal || editingAsset) && (
        <AssetModal
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
      )}

      {/* Asset Details Modal */}
      {viewingAssetId && (
        <AssetDetailsModal
          assetId={viewingAssetId}
          onClose={() => setViewingAssetId(null)}
          onEdit={(asset) => {
            setViewingAssetId(null);
            setEditingAsset(asset);
          }}
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
    </div>
  );
};
