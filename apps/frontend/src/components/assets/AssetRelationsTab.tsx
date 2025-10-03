import React, { useState, useEffect } from 'react';
import { risksApi, Risk, RISK_LEVELS, RISK_STATUSES } from '../../shared/api/risks';
import { incidentsApi, Incident, INCIDENT_CRITICALITY, INCIDENT_STATUS } from '../../shared/api/incidents';

interface AssetRelationsTabProps {
  assetId: string;
  assetName: string;
}

const AssetRelationsTab: React.FC<AssetRelationsTabProps> = ({ assetId, assetName }) => {
  const [risks, setRisks] = useState<Risk[]>([]);
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'risks' | 'incidents'>('risks');

  useEffect(() => {
    loadRelations();
  }, [assetId]);

  const loadRelations = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const [risksResponse, incidentsResponse] = await Promise.all([
        risksApi.getByAsset(assetId),
        incidentsApi.getByAsset(assetId)
      ]);
      
      setRisks(risksResponse.data || []);
      setIncidents(incidentsResponse.data || []);
    } catch (err) {
      setError('Ошибка загрузки связанных данных');
      console.error('Error loading relations:', err);
    } finally {
      setLoading(false);
    }
  };

  const getRiskLevelColor = (level: string) => {
    switch (level) {
      case 'critical': return 'text-red-600 bg-red-100';
      case 'high': return 'text-orange-600 bg-orange-100';
      case 'medium': return 'text-yellow-600 bg-yellow-100';
      case 'low': return 'text-green-600 bg-green-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'text-red-600 bg-red-100';
      case 'high': return 'text-orange-600 bg-orange-100';
      case 'medium': return 'text-yellow-600 bg-yellow-100';
      case 'low': return 'text-green-600 bg-green-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getStatusColor = (status: string, isRisk: boolean = true) => {
    if (isRisk) {
      switch (status) {
        case 'open': return 'text-red-600 bg-red-100';
        case 'mitigated': return 'text-yellow-600 bg-yellow-100';
        case 'accepted': return 'text-blue-600 bg-blue-100';
        case 'closed': return 'text-green-600 bg-green-100';
        default: return 'text-gray-600 bg-gray-100';
      }
    } else {
      switch (status) {
        case 'open': return 'text-red-600 bg-red-100';
        case 'in_progress': return 'text-yellow-600 bg-yellow-100';
        case 'resolved': return 'text-blue-600 bg-blue-100';
        case 'closed': return 'text-green-600 bg-green-100';
        default: return 'text-gray-600 bg-gray-100';
      }
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU');
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg">Загрузка связанных данных...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8 text-red-600">
        {error}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Tabs */}
      <div className="border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          <button
            onClick={() => setActiveTab('risks')}
            className={`py-2 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'risks'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            Риски ({risks.length})
          </button>
          <button
            onClick={() => setActiveTab('incidents')}
            className={`py-2 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'incidents'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            Инциденты ({incidents.length})
          </button>
        </nav>
      </div>

      {/* Risks Tab */}
      {activeTab === 'risks' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-medium text-gray-900">
              Риски для актива "{assetName}"
            </h3>
            <button className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
              Добавить риск
            </button>
          </div>
          
          {risks.length > 0 ? (
            <div className="bg-white shadow overflow-hidden sm:rounded-md">
              <ul className="divide-y divide-gray-200">
                {risks.map((risk) => (
                  <li key={risk.id} className="px-6 py-4">
                    <div className="flex items-center justify-between">
                      <div className="flex-1">
                        <div className="flex items-center">
                          <h4 className="text-sm font-medium text-gray-900 mr-2">
                            {risk.name}
                          </h4>
                          <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getRiskLevelColor(risk.risk_level)}`}>
                            {RISK_LEVELS.find(r => r.value === risk.risk_level)?.label || risk.risk_level}
                          </span>
                          <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ml-2 ${getStatusColor(risk.status, true)}`}>
                            {RISK_STATUSES.find(s => s.value === risk.status)?.label || risk.status}
                          </span>
                        </div>
                        {risk.description && (
                          <p className="text-sm text-gray-500 mt-1">{risk.description}</p>
                        )}
                        <div className="text-xs text-gray-400 mt-1">
                          Создан: {formatDate(risk.created_at)}
                          {risk.owner_name && ` • Владелец: ${risk.owner_name}`}
                        </div>
                      </div>
                      <div className="flex space-x-2">
                        <button className="text-blue-600 hover:text-blue-900 text-sm">
                          Просмотр
                        </button>
                        <button className="text-green-600 hover:text-green-900 text-sm">
                          Редактировать
                        </button>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">
              Риски не найдены
            </div>
          )}
        </div>
      )}

      {/* Incidents Tab */}
      {activeTab === 'incidents' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-medium text-gray-900">
              Инциденты для актива "{assetName}"
            </h3>
            <button className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
              Добавить инцидент
            </button>
          </div>
          
          {incidents.length > 0 ? (
            <div className="bg-white shadow overflow-hidden sm:rounded-md">
              <ul className="divide-y divide-gray-200">
                {incidents.map((incident) => (
                  <li key={incident.id} className="px-6 py-4">
                    <div className="flex items-center justify-between">
                      <div className="flex-1">
                        <div className="flex items-center">
                          <h4 className="text-sm font-medium text-gray-900 mr-2">
                            {incident.title}
                          </h4>
                          <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getSeverityColor(incident.criticality)}`}>
                            {Object.values(INCIDENT_CRITICALITY).find(s => s === incident.criticality) || incident.criticality}
                          </span>
                          <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ml-2 ${getStatusColor(incident.status, false)}`}>
                            {Object.values(INCIDENT_STATUS).find(s => s === incident.status) || incident.status}
                          </span>
                        </div>
                        {incident.description && (
                          <p className="text-sm text-gray-500 mt-1">{incident.description}</p>
                        )}
                        <div className="text-xs text-gray-400 mt-1">
                          Создан: {formatDate(incident.created_at)}
                          {incident.assigned_to_name && ` • Назначен: ${incident.assigned_to_name}`}
                          {incident.resolved_at && ` • Решен: ${formatDate(incident.resolved_at)}`}
                        </div>
                      </div>
                      <div className="flex space-x-2">
                        <button className="text-blue-600 hover:text-blue-900 text-sm">
                          Просмотр
                        </button>
                        <button className="text-green-600 hover:text-green-900 text-sm">
                          Редактировать
                        </button>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">
              Инциденты не найдены
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default AssetRelationsTab;

