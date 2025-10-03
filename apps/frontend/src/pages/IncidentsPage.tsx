import React, { useState, useEffect } from 'react';
import { incidentsApi, Incident, CreateIncidentRequest, UpdateIncidentRequest } from '../shared/api/incidents';
import { usersApi, User } from '../shared/api/users';

const IncidentsPage: React.FC = () => {
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showViewModal, setShowViewModal] = useState(false);
  const [selectedIncident, setSelectedIncident] = useState<Incident | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [incidentsResponse, usersResponse] = await Promise.all([
        incidentsApi.list({ page: 1, page_size: 20 }),
        usersApi.list({ page: 1, page_size: 100 })
      ]);
      
      setIncidents(incidentsResponse.data || []);
      setUsers(usersResponse.data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateIncident = async (data: CreateIncidentRequest) => {
    try {
      await incidentsApi.create(data);
      setShowCreateModal(false);
      loadData();
    } catch (err: any) {
      throw new Error(err.response?.data?.error || 'Failed to create incident');
    }
  };

  const handleUpdateIncident = async (data: UpdateIncidentRequest) => {
    try {
      if (selectedIncident) {
        await incidentsApi.update(selectedIncident.id, data);
        setShowEditModal(false);
        setSelectedIncident(null);
        loadData();
      }
    } catch (err: any) {
      throw new Error(err.response?.data?.error || 'Failed to update incident');
    }
  };

  const handleDeleteIncident = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этот инцидент?')) {
      try {
        await incidentsApi.delete(id);
        loadData();
      } catch (err: any) {
        setError(err.response?.data?.error || 'Failed to delete incident');
      }
    }
  };

  const handleViewIncident = (incident: Incident) => {
    setSelectedIncident(incident);
    setShowViewModal(true);
  };

  const handleEditIncident = (incident: Incident) => {
    setSelectedIncident(incident);
    setShowEditModal(true);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'new': return 'bg-blue-100 text-blue-800';
      case 'assigned': return 'bg-yellow-100 text-yellow-800';
      case 'in_progress': return 'bg-orange-100 text-orange-800';
      case 'resolved': return 'bg-green-100 text-green-800';
      case 'closed': return 'bg-gray-100 text-gray-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getCriticalityColor = (criticality: string) => {
    switch (criticality) {
      case 'low': return 'bg-green-100 text-green-800';
      case 'medium': return 'bg-yellow-100 text-yellow-800';
      case 'high': return 'bg-orange-100 text-orange-800';
      case 'critical': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        <h1 className="text-2xl font-bold mb-6">Инциденты</h1>
        <div className="text-center py-8">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p className="mt-2 text-gray-600">Загрузка инцидентов...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6">
        <h1 className="text-2xl font-bold mb-6">Инциденты</h1>
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <p className="text-red-800">Ошибка: {error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Инциденты</h1>
        <button 
          onClick={() => setShowCreateModal(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
        >
          Создать инцидент
        </button>
      </div>

      <div className="bg-white shadow rounded-lg">
        {incidents.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-gray-500">Инциденты не найдены</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Название
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Статус
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Критичность
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Категория
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Ответственный
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Дата создания
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Действия
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
              {incidents.map((incident) => (
                  <tr key={incident.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-medium text-gray-900">
                      {incident.title}
                      </div>
                      {incident.description && (
                        <div className="text-sm text-gray-500 truncate max-w-xs">
                          {incident.description}
                        </div>
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(incident.status)}`}>
                        {incident.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getCriticalityColor(incident.criticality)}`}>
                        {incident.criticality}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {incident.category}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {incident.assigned_name || 'Не назначен'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(incident.created_at).toLocaleDateString('ru-RU')}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <div className="flex space-x-2">
                        <button
                          onClick={() => handleViewIncident(incident)}
                          className="text-blue-600 hover:text-blue-900"
                        >
                          Просмотр
                        </button>
                        <button
                          onClick={() => handleEditIncident(incident)}
                          className="text-indigo-600 hover:text-indigo-900"
                        >
                          Редактировать
                        </button>
                        <button
                          onClick={() => handleDeleteIncident(incident.id)}
                          className="text-red-600 hover:text-red-900"
                        >
                          Удалить
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Create Modal */}
      {showCreateModal && (
        <CreateIncidentModal
          onClose={() => setShowCreateModal(false)}
          onSave={handleCreateIncident}
          users={users}
        />
      )}

      {/* Edit Modal */}
      {showEditModal && selectedIncident && (
        <EditIncidentModal
          incident={selectedIncident}
          onClose={() => {
            setShowEditModal(false);
            setSelectedIncident(null);
          }}
          onSave={handleUpdateIncident}
          users={users}
        />
      )}

      {/* View Modal */}
      {showViewModal && selectedIncident && (
        <ViewIncidentModal
          incident={selectedIncident}
          onClose={() => {
            setShowViewModal(false);
            setSelectedIncident(null);
          }}
        />
      )}
    </div>
  );
};

// Create Incident Modal Component
const CreateIncidentModal: React.FC<{
  onClose: () => void;
  onSave: (data: CreateIncidentRequest) => Promise<void>;
  users: User[];
}> = ({ onClose, onSave, users }) => {
  const [formData, setFormData] = useState<CreateIncidentRequest>({
    title: '',
    description: '',
    category: 'technical_failure',
    criticality: 'medium',
    source: 'user_report',
    assigned_to: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await onSave(formData);
    } catch (error: any) {
      alert(error.message);
    }
  };

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div className="mt-3">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Создать инцидент</h3>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Название</label>
              <input
                type="text"
                required
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Описание</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                rows={3}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Категория</label>
              <select
                value={formData.category}
                onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="technical_failure">Технический сбой</option>
                <option value="data_breach">Утечка данных</option>
                <option value="unauthorized_access">Несанкционированный доступ</option>
                <option value="physical">Физический инцидент</option>
                <option value="malware">Вредоносное ПО</option>
                <option value="social_engineering">Социальная инженерия</option>
                <option value="other">Другое</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Критичность</label>
              <select
                value={formData.criticality}
                onChange={(e) => setFormData({ ...formData, criticality: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="low">Низкая</option>
                <option value="medium">Средняя</option>
                <option value="high">Высокая</option>
                <option value="critical">Критическая</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Источник</label>
              <select
                value={formData.source}
                onChange={(e) => setFormData({ ...formData, source: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="user_report">Сообщение пользователя</option>
                <option value="automatic_agent">Автоматический агент</option>
                <option value="admin_manual">Ручное создание админом</option>
                <option value="monitoring">Мониторинг</option>
                <option value="siem">SIEM</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Ответственный</label>
              <select
                value={formData.assigned_to}
                onChange={(e) => setFormData({ ...formData, assigned_to: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Не назначен</option>
                {users.map(user => (
                  <option key={user.id} value={user.id}>
                    {user.first_name} {user.last_name}
                  </option>
                ))}
              </select>
            </div>
            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={onClose}
                className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
              >
                Отмена
              </button>
              <button
                type="submit"
                className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
              >
                Создать
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

// Edit Incident Modal Component
const EditIncidentModal: React.FC<{
  incident: Incident;
  onClose: () => void;
  onSave: (data: UpdateIncidentRequest) => Promise<void>;
  users: User[];
}> = ({ incident, onClose, onSave, users }) => {
  const [formData, setFormData] = useState<UpdateIncidentRequest>({
    title: incident.title,
    description: incident.description || '',
    category: incident.category,
    criticality: incident.criticality,
    source: incident.source,
    assigned_to: incident.assigned_to || '',
    status: incident.status,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await onSave(formData);
    } catch (error: any) {
      alert(error.message);
    }
  };

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div className="mt-3">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Редактировать инцидент</h3>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Название</label>
              <input
                type="text"
                required
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Описание</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                rows={3}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Статус</label>
              <select
                value={formData.status}
                onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="new">Новый</option>
                <option value="assigned">Назначен</option>
                <option value="in_progress">В работе</option>
                <option value="resolved">Решен</option>
                <option value="closed">Закрыт</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Критичность</label>
              <select
                value={formData.criticality}
                onChange={(e) => setFormData({ ...formData, criticality: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="low">Низкая</option>
                <option value="medium">Средняя</option>
                <option value="high">Высокая</option>
                <option value="critical">Критическая</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Ответственный</label>
              <select
                value={formData.assigned_to}
                onChange={(e) => setFormData({ ...formData, assigned_to: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">Не назначен</option>
                {users.map(user => (
                  <option key={user.id} value={user.id}>
                    {user.first_name} {user.last_name}
                  </option>
                ))}
              </select>
            </div>
            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={onClose}
                className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
              >
                Отмена
              </button>
              <button
                type="submit"
                className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
              >
                Сохранить
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

// View Incident Modal Component
const ViewIncidentModal: React.FC<{
  incident: Incident;
  onClose: () => void;
}> = ({ incident, onClose }) => {
  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-20 mx-auto p-5 border w-2/3 shadow-lg rounded-md bg-white">
        <div className="mt-3">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Детали инцидента</h3>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Название</label>
              <p className="mt-1 text-sm text-gray-900">{incident.title}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Статус</label>
              <p className="mt-1 text-sm text-gray-900">{incident.status}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Критичность</label>
              <p className="mt-1 text-sm text-gray-900">{incident.criticality}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Категория</label>
              <p className="mt-1 text-sm text-gray-900">{incident.category}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Источник</label>
              <p className="mt-1 text-sm text-gray-900">{incident.source}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Ответственный</label>
              <p className="mt-1 text-sm text-gray-900">{incident.assigned_name || 'Не назначен'}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Дата обнаружения</label>
              <p className="mt-1 text-sm text-gray-900">
                {new Date(incident.detected_at).toLocaleString('ru-RU')}
              </p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Дата создания</label>
              <p className="mt-1 text-sm text-gray-900">
                {new Date(incident.created_at).toLocaleString('ru-RU')}
              </p>
            </div>
            {incident.resolved_at && (
              <div>
                <label className="block text-sm font-medium text-gray-700">Дата решения</label>
                <p className="mt-1 text-sm text-gray-900">
                  {new Date(incident.resolved_at).toLocaleString('ru-RU')}
                </p>
              </div>
            )}
            {incident.closed_at && (
              <div>
                <label className="block text-sm font-medium text-gray-700">Дата закрытия</label>
                <p className="mt-1 text-sm text-gray-900">
                  {new Date(incident.closed_at).toLocaleString('ru-RU')}
                </p>
              </div>
            )}
          </div>
          {incident.description && (
            <div className="mt-4">
              <label className="block text-sm font-medium text-gray-700">Описание</label>
              <p className="mt-1 text-sm text-gray-900">{incident.description}</p>
            </div>
          )}
          <div className="flex justify-end mt-6">
            <button
              onClick={onClose}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
            >
              Закрыть
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default IncidentsPage;