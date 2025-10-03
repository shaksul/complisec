import React, { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Asset, ASSET_TYPES, ASSET_CLASSES, CRITICALITY_LEVELS, ASSET_STATUSES } from '../../shared/api/assets';
import { User, UserCatalog } from '../../shared/api/users';
import { createAssetSchema, updateAssetSchema, CreateAssetFormData, UpdateAssetFormData } from '../../shared/validation/assets';

interface AssetModalProps {
  asset?: Asset | null;
  users: UserCatalog[];
  onSave: (data: any) => void;
  onClose: () => void;
}

const AssetModal: React.FC<AssetModalProps> = ({ asset, users, onSave, onClose }) => {
  const isEdit = !!asset;
  const schema = isEdit ? updateAssetSchema : createAssetSchema;
  
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset
  } = useForm<CreateAssetFormData | UpdateAssetFormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: asset?.name || '',
      type: asset?.type || '',
      class: asset?.class || '',
      owner_id: asset?.owner_id || '',
      responsible_user_id: asset?.responsible_user_id || '',
      location: asset?.location || '',
      criticality: asset?.criticality || '',
      confidentiality: asset?.confidentiality || '',
      integrity: asset?.integrity || '',
      availability: asset?.availability || '',
      status: asset?.status || 'active'
    }
  });

  // Сброс формы при изменении актива
  useEffect(() => {
    if (asset) {
      reset({
        name: asset.name,
        type: asset.type,
        class: asset.class,
        owner_id: asset.owner_id || '',
        responsible_user_id: asset.responsible_user_id || '',
        location: asset.location || '',
        criticality: asset.criticality,
        confidentiality: asset.confidentiality,
        integrity: asset.integrity,
        availability: asset.availability,
        status: asset.status
      });
    } else {
      reset({
        name: '',
        type: '',
        class: '',
        owner_id: '',
        responsible_user_id: '',
        location: '',
        criticality: '',
        confidentiality: '',
        integrity: '',
        availability: '',
        status: 'active'
      });
    }
  }, [asset, reset]);

  const onSubmit = (data: CreateAssetFormData | UpdateAssetFormData) => {
    // Очистка пустых строковых полей для обновления
    if (isEdit) {
      const updateData = { ...data };
      Object.keys(updateData).forEach(key => {
        const value = updateData[key as keyof typeof updateData];
        if (value === '' || value === null || value === undefined) {
          delete updateData[key as keyof typeof updateData];
        }
      });
      onSave(updateData);
    } else {
      onSave(data);
    }
  };

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-20 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-1/2 shadow-lg rounded-md bg-white">
        <div className="mt-3">
          <h3 className="text-lg font-medium text-gray-900 mb-4">
            {asset ? 'Редактировать актив' : 'Создать актив'}
          </h3>
          
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Название {!isEdit && '*'}
                </label>
                <input
                  type="text"
                  {...register('name')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.name ? 'border-red-500' : 'border-gray-300'
                  }`}
                  placeholder="Введите название актива"
                />
                {errors.name && (
                  <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Тип {!isEdit && '*'}
                </label>
                <select
                  {...register('type')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.type ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  <option value="">Выберите тип</option>
                  {ASSET_TYPES.map(type => (
                    <option key={type.value} value={type.value}>{type.label}</option>
                  ))}
                </select>
                {errors.type && (
                  <p className="mt-1 text-sm text-red-600">{errors.type.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Класс {!isEdit && '*'}
                </label>
                <select
                  {...register('class')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.class ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  <option value="">Выберите класс</option>
                  {ASSET_CLASSES.map(cls => (
                    <option key={cls.value} value={cls.value}>{cls.label}</option>
                  ))}
                </select>
                {errors.class && (
                  <p className="mt-1 text-sm text-red-600">{errors.class.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Владелец (организация)</label>
                <select
                  {...register('owner_id')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.owner_id ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  <option value="">Организация</option>
                  {users.map(user => (
                    <option key={user.id} value={user.id}>
                      {user.first_name && user.last_name 
                        ? `${user.first_name} ${user.last_name}` 
                        : user.email}
                    </option>
                  ))}
                </select>
                {errors.owner_id && (
                  <p className="mt-1 text-sm text-red-600">{errors.owner_id.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Ответственный пользователь</label>
                <select
                  {...register('responsible_user_id')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.responsible_user_id ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  <option value="">Не назначен</option>
                  {users.map(user => (
                    <option key={user.id} value={user.id}>
                      {user.first_name && user.last_name 
                        ? `${user.first_name} ${user.last_name}` 
                        : user.email}
                    </option>
                  ))}
                </select>
                {errors.responsible_user_id && (
                  <p className="mt-1 text-sm text-red-600">{errors.responsible_user_id.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Местоположение</label>
                <input
                  type="text"
                  {...register('location')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.location ? 'border-red-500' : 'border-gray-300'
                  }`}
                  placeholder="Введите местоположение"
                />
                {errors.location && (
                  <p className="mt-1 text-sm text-red-600">{errors.location.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Критичность {!isEdit && '*'}
                </label>
                <select
                  {...register('criticality')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.criticality ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  <option value="">Выберите критичность</option>
                  {CRITICALITY_LEVELS.map(level => (
                    <option key={level.value} value={level.value}>{level.label}</option>
                  ))}
                </select>
                {errors.criticality && (
                  <p className="mt-1 text-sm text-red-600">{errors.criticality.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Конфиденциальность {!isEdit && '*'}
                </label>
                <select
                  {...register('confidentiality')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.confidentiality ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  <option value="">Выберите уровень</option>
                  {CRITICALITY_LEVELS.map(level => (
                    <option key={level.value} value={level.value}>{level.label}</option>
                  ))}
                </select>
                {errors.confidentiality && (
                  <p className="mt-1 text-sm text-red-600">{errors.confidentiality.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Целостность {!isEdit && '*'}
                </label>
                <select
                  {...register('integrity')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.integrity ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  <option value="">Выберите уровень</option>
                  {CRITICALITY_LEVELS.map(level => (
                    <option key={level.value} value={level.value}>{level.label}</option>
                  ))}
                </select>
                {errors.integrity && (
                  <p className="mt-1 text-sm text-red-600">{errors.integrity.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Доступность {!isEdit && '*'}
                </label>
                <select
                  {...register('availability')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.availability ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  <option value="">Выберите уровень</option>
                  {CRITICALITY_LEVELS.map(level => (
                    <option key={level.value} value={level.value}>{level.label}</option>
                  ))}
                </select>
                {errors.availability && (
                  <p className="mt-1 text-sm text-red-600">{errors.availability.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Статус</label>
                <select
                  {...register('status')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.status ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  {ASSET_STATUSES.map(status => (
                    <option key={status.value} value={status.value}>{status.label}</option>
                  ))}
                </select>
                {errors.status && (
                  <p className="mt-1 text-sm text-red-600">{errors.status.message}</p>
                )}
              </div>
            </div>
            
            <div className="flex justify-end space-x-2 pt-4">
              <button
                type="button"
                onClick={onClose}
                disabled={isSubmitting}
                className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Отмена
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? 'Сохранение...' : (asset ? 'Сохранить' : 'Создать')}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default AssetModal;
