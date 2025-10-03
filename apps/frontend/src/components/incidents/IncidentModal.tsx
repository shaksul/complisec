import React, { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  Incident,
  CreateIncidentRequest,
  UpdateIncidentRequest,
  INCIDENT_CATEGORIES,
  INCIDENT_CRITICALITY,
  INCIDENT_STATUS,
  INCIDENT_SOURCE,
  getCategoryLabel,
  getCriticalityLabel,
  getStatusLabel,
  getSourceLabel,
} from '../../shared/api/incidents';
import { User } from '../../shared/api/users';

const incidentSchema = z.object({
  title: z.string().min(1, 'Название обязательно'),
  description: z.string().optional(),
  category: z.enum(Object.values(INCIDENT_CATEGORIES) as [string, ...string[]], {
    required_error: 'Категория обязательна',
  }),
  criticality: z.enum(Object.values(INCIDENT_CRITICALITY) as [string, ...string[]], {
    required_error: 'Критичность обязательна',
  }),
  source: z.enum(Object.values(INCIDENT_SOURCE) as [string, ...string[]], {
    required_error: 'Источник обязателен',
  }),
  assigned_to: z.string().optional().transform(val => val === '' ? undefined : val),
  detected_at: z.string().optional(),
  asset_ids: z.array(z.string()).optional(),
  risk_ids: z.array(z.string()).optional(),
});

type IncidentFormData = z.infer<typeof incidentSchema>;

interface IncidentModalProps {
  incident?: Incident;
  onClose: () => void;
  onSave: (data: CreateIncidentRequest | UpdateIncidentRequest) => Promise<void>;
  users: User[];
}

const IncidentModal: React.FC<IncidentModalProps> = ({
  incident,
  onClose,
  onSave,
  users,
}) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    watch,
  } = useForm<IncidentFormData>({
    resolver: zodResolver(incidentSchema),
    defaultValues: {
      title: incident?.title || '',
      description: incident?.description || '',
      category: incident?.category || INCIDENT_CATEGORIES.TECHNICAL_FAILURE,
      criticality: incident?.criticality || INCIDENT_CRITICALITY.MEDIUM,
      source: incident?.source || INCIDENT_SOURCE.USER_REPORT,
      assigned_to: incident?.assigned_to || '',
      detected_at: incident?.detected_at ? new Date(incident.detected_at).toISOString().slice(0, 16) : '',
      asset_ids: incident?.assets?.map(a => a.id) || [],
      risk_ids: incident?.risks?.map(r => r.id) || [],
    },
  });

  const isEditing = !!incident;

  const onSubmit = async (data: IncidentFormData) => {
    try {
      setLoading(true);
      setError(null);

      const submitData = {
        ...data,
        assigned_to: data.assigned_to && data.assigned_to.trim() !== '' ? data.assigned_to : undefined,
        detected_at: data.detected_at ? new Date(data.detected_at).toISOString() : undefined,
        // Remove empty arrays
        asset_ids: data.asset_ids && data.asset_ids.length > 0 ? data.asset_ids : undefined,
        risk_ids: data.risk_ids && data.risk_ids.length > 0 ? data.risk_ids : undefined,
      };

      await onSave(submitData);
    } catch (err: any) {
      setError(err.message || 'Произошла ошибка при сохранении');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    reset();
    setError(null);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-20 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-1/2 shadow-lg rounded-md bg-white">
        <div className="mt-3">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-medium text-gray-900">
              {isEditing ? 'Редактировать инцидент' : 'Создать инцидент'}
            </h3>
            <button
              onClick={handleClose}
              className="text-gray-400 hover:text-gray-600"
            >
              <span className="sr-only">Закрыть</span>
              <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {error && (
            <div className="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Название *
                </label>
                <input
                  {...register('title')}
                  type="text"
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Введите название инцидента"
                />
                {errors.title && (
                  <p className="mt-1 text-sm text-red-600">{errors.title.message}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Категория *
                </label>
                <select
                  {...register('category')}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  {Object.values(INCIDENT_CATEGORIES).map(category => (
                    <option key={category} value={category}>
                      {getCategoryLabel(category)}
                    </option>
                  ))}
                </select>
                {errors.category && (
                  <p className="mt-1 text-sm text-red-600">{errors.category.message}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Критичность *
                </label>
                <select
                  {...register('criticality')}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  {Object.values(INCIDENT_CRITICALITY).map(criticality => (
                    <option key={criticality} value={criticality}>
                      {getCriticalityLabel(criticality)}
                    </option>
                  ))}
                </select>
                {errors.criticality && (
                  <p className="mt-1 text-sm text-red-600">{errors.criticality.message}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Источник *
                </label>
                <select
                  {...register('source')}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  {Object.values(INCIDENT_SOURCE).map(source => (
                    <option key={source} value={source}>
                      {getSourceLabel(source)}
                    </option>
                  ))}
                </select>
                {errors.source && (
                  <p className="mt-1 text-sm text-red-600">{errors.source.message}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Ответственный
                </label>
                <select
                  {...register('assigned_to')}
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Не назначен</option>
                  {users.map(user => (
                    <option key={user.id} value={user.id}>
                      {user.first_name} {user.last_name}
                    </option>
                  ))}
                </select>
                {errors.assigned_to && (
                  <p className="mt-1 text-sm text-red-600">{errors.assigned_to.message}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Дата обнаружения
                </label>
                <input
                  {...register('detected_at')}
                  type="datetime-local"
                  className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                {errors.detected_at && (
                  <p className="mt-1 text-sm text-red-600">{errors.detected_at.message}</p>
                )}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Описание
              </label>
              <textarea
                {...register('description')}
                rows={4}
                className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Введите описание инцидента"
              />
              {errors.description && (
                <p className="mt-1 text-sm text-red-600">{errors.description.message}</p>
              )}
            </div>

            <div className="flex justify-end space-x-3 pt-4">
              <button
                type="button"
                onClick={handleClose}
                className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
              >
                Отмена
              </button>
              <button
                type="submit"
                disabled={loading}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
              >
                {loading ? 'Сохранение...' : (isEditing ? 'Сохранить' : 'Создать')}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default IncidentModal;