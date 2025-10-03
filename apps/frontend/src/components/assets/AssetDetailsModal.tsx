import React, { useState, useEffect } from 'react';
import { AssetWithDetails, AssetDocument, AssetSoftware, AssetHistory, DOCUMENT_TYPES } from '../../shared/api/assets';
import { assetsApi } from '../../shared/api/assets';
import AssetRelationsTab from './AssetRelationsTab';

interface AssetDetailsModalProps {
  assetId: string;
  onClose: () => void;
  onEdit: (asset: any) => void;
}

const AssetDetailsModal: React.FC<AssetDetailsModalProps> = ({ assetId, onClose, onEdit }) => {
  const [asset, setAsset] = useState<AssetWithDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'overview' | 'documents' | 'software' | 'history' | 'relations'>('overview');

  useEffect(() => {
    loadAssetDetails();
  }, [assetId]);

  const loadAssetDetails = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await assetsApi.getDetails(assetId);
      setAsset(response.data);
    } catch (err) {
      setError('Ошибка загрузки деталей актива');
      console.error('Error loading asset details:', err);
    } finally {
      setLoading(false);
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

  const getDocumentTypeLabel = (type: string) => {
    return DOCUMENT_TYPES.find(dt => dt.value === type)?.label || type;
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU');
  };

  if (loading) {
    return (
      <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
        <div className="relative top-20 mx-auto p-5 border w-11/12 md:w-4/5 lg:w-3/4 xl:w-2/3 shadow-lg rounded-md bg-white">
          <div className="flex justify-center items-center h-64">
            <div className="text-lg">Загрузка...</div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !asset) {
    return (
      <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
        <div className="relative top-20 mx-auto p-5 border w-11/12 md:w-4/5 lg:w-3/4 xl:w-2/3 shadow-lg rounded-md bg-white">
          <div className="text-center py-8">
            <div className="text-red-600 text-lg mb-4">{error || 'Актив не найден'}</div>
            <button
              onClick={onClose}
              className="px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600"
            >
              Закрыть
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-10 mx-auto p-5 border w-11/12 md:w-4/5 lg:w-3/4 xl:w-2/3 shadow-lg rounded-md bg-white">
        {/* Header */}
        <div className="flex justify-between items-start mb-6">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">{asset.name}</h2>
            <p className="text-sm text-gray-500">Инв. номер: {asset.inventory_number}</p>
          </div>
          <div className="flex space-x-2">
            <button
              onClick={() => onEdit(asset)}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              Редактировать
            </button>
            <button
              onClick={onClose}
              className="px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500"
            >
              Закрыть
            </button>
          </div>
        </div>

        {/* Tabs */}
        <div className="border-b border-gray-200 mb-6">
          <nav className="-mb-px flex space-x-8">
            {[
              { id: 'overview', label: 'Обзор' },
              { id: 'documents', label: `Документы (${asset.documents?.length || 0})` },
              { id: 'software', label: `ПО (${asset.software?.length || 0})` },
              { id: 'history', label: `История (${asset.history?.length || 0})` },
              { id: 'relations', label: 'Связи' }
            ].map(tab => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as any)}
                className={`py-2 px-1 border-b-2 font-medium text-sm ${
                  activeTab === tab.id
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </nav>
        </div>

        {/* Tab Content */}
        <div className="min-h-96">
          {activeTab === 'overview' && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* Основная информация */}
              <div className="space-y-4">
                <h3 className="text-lg font-medium text-gray-900">Основная информация</h3>
                <div className="bg-gray-50 p-4 rounded-lg space-y-3">
                  <div>
                    <label className="text-sm font-medium text-gray-500">Тип</label>
                    <p className="text-sm text-gray-900 capitalize">{asset.type}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Класс</label>
                    <p className="text-sm text-gray-900 capitalize">{asset.class}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Владелец</label>
                    <p className="text-sm text-gray-900">{asset.owner_name || 'Не назначен'}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Местоположение</label>
                    <p className="text-sm text-gray-900">{asset.location || 'Не указано'}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Статус</label>
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(asset.status)}`}>
                      {asset.status}
                    </span>
                  </div>
                </div>
              </div>

              {/* CIA Оценка */}
              <div className="space-y-4">
                <h3 className="text-lg font-medium text-gray-900">CIA Оценка</h3>
                <div className="bg-gray-50 p-4 rounded-lg space-y-3">
                  <div>
                    <label className="text-sm font-medium text-gray-500">Критичность</label>
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getCriticalityColor(asset.criticality)}`}>
                      {asset.criticality}
                    </span>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Конфиденциальность</label>
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getCriticalityColor(asset.confidentiality)}`}>
                      {asset.confidentiality}
                    </span>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Целостность</label>
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getCriticalityColor(asset.integrity)}`}>
                      {asset.integrity}
                    </span>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Доступность</label>
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getCriticalityColor(asset.availability)}`}>
                      {asset.availability}
                    </span>
                  </div>
                </div>
              </div>

              {/* Метаданные */}
              <div className="space-y-4 md:col-span-2">
                <h3 className="text-lg font-medium text-gray-900">Метаданные</h3>
                <div className="bg-gray-50 p-4 rounded-lg grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium text-gray-500">Создан</label>
                    <p className="text-sm text-gray-900">{formatDate(asset.created_at)}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Обновлен</label>
                    <p className="text-sm text-gray-900">{formatDate(asset.updated_at)}</p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'documents' && (
            <div className="space-y-4">
              <div className="flex justify-between items-center">
                <h3 className="text-lg font-medium text-gray-900">Документы актива</h3>
                <button className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
                  Добавить документ
                </button>
              </div>
              {asset.documents && asset.documents.length > 0 ? (
                <div className="bg-white shadow overflow-hidden sm:rounded-md">
                  <ul className="divide-y divide-gray-200">
                    {asset.documents.map((doc: AssetDocument) => (
                      <li key={doc.id} className="px-6 py-4">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center">
                            <div className="flex-shrink-0">
                              <div className="h-8 w-8 bg-gray-200 rounded-full flex items-center justify-center">
                                <span className="text-xs font-medium text-gray-600">DOC</span>
                              </div>
                            </div>
                            <div className="ml-4">
                              <div className="text-sm font-medium text-gray-900">
                                {getDocumentTypeLabel(doc.document_type)}
                              </div>
                              <div className="text-sm text-gray-500">
                                {formatDate(doc.created_at)}
                              </div>
                            </div>
                          </div>
                          <div className="flex space-x-2">
                            <button className="text-blue-600 hover:text-blue-900 text-sm">
                              Скачать
                            </button>
                            <button className="text-red-600 hover:text-red-900 text-sm">
                              Удалить
                            </button>
                          </div>
                        </div>
                      </li>
                    ))}
                  </ul>
                </div>
              ) : (
                <div className="text-center py-8 text-gray-500">
                  Документы не найдены
                </div>
              )}
            </div>
          )}

          {activeTab === 'software' && (
            <div className="space-y-4">
              <div className="flex justify-between items-center">
                <h3 className="text-lg font-medium text-gray-900">Установленное ПО</h3>
                <button className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
                  Добавить ПО
                </button>
              </div>
              {asset.software && asset.software.length > 0 ? (
                <div className="bg-white shadow overflow-hidden sm:rounded-md">
                  <ul className="divide-y divide-gray-200">
                    {asset.software.map((sw: AssetSoftware) => (
                      <li key={sw.id} className="px-6 py-4">
                        <div className="flex items-center justify-between">
                          <div>
                            <div className="text-sm font-medium text-gray-900">
                              {sw.software_name}
                            </div>
                            <div className="text-sm text-gray-500">
                              {sw.version && `Версия: ${sw.version}`}
                              {sw.installed_at && ` • Установлено: ${formatDate(sw.installed_at)}`}
                            </div>
                          </div>
                          <div className="flex space-x-2">
                            <button className="text-blue-600 hover:text-blue-900 text-sm">
                              Редактировать
                            </button>
                            <button className="text-red-600 hover:text-red-900 text-sm">
                              Удалить
                            </button>
                          </div>
                        </div>
                      </li>
                    ))}
                  </ul>
                </div>
              ) : (
                <div className="text-center py-8 text-gray-500">
                  ПО не найдено
                </div>
              )}
            </div>
          )}

          {activeTab === 'history' && (
            <div className="space-y-4">
              <h3 className="text-lg font-medium text-gray-900">История изменений</h3>
              {asset.history && asset.history.length > 0 ? (
                <div className="bg-white shadow overflow-hidden sm:rounded-md">
                  <ul className="divide-y divide-gray-200">
                    {asset.history.map((entry: AssetHistory) => (
                      <li key={entry.id} className="px-6 py-4">
                        <div className="flex items-start">
                          <div className="flex-shrink-0">
                            <div className="h-8 w-8 bg-gray-200 rounded-full flex items-center justify-center">
                              <span className="text-xs font-medium text-gray-600">H</span>
                            </div>
                          </div>
                          <div className="ml-4 flex-1">
                            <div className="text-sm font-medium text-gray-900">
                              {entry.field_changed}
                            </div>
                            <div className="text-sm text-gray-500">
                              {entry.old_value && `Было: ${entry.old_value}`}
                              {entry.old_value && entry.new_value && ' → '}
                              {entry.new_value && `Стало: ${entry.new_value}`}
                            </div>
                            <div className="text-xs text-gray-400 mt-1">
                              {formatDate(entry.changed_at)} • {entry.changed_by}
                            </div>
                          </div>
                        </div>
                      </li>
                    ))}
                  </ul>
                </div>
              ) : (
                <div className="text-center py-8 text-gray-500">
                  История изменений пуста
                </div>
              )}
            </div>
          )}

          {activeTab === 'relations' && (
            <AssetRelationsTab
              assetId={asset.id}
              assetName={asset.name}
            />
          )}
        </div>
      </div>
    </div>
  );
};

export default AssetDetailsModal;
