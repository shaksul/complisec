import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Asset, ASSET_STATUSES } from '../../shared/api/assets';
import { assetsApi } from '../../shared/api/assets';
import { assetInventorySchema, AssetInventoryFormData } from '../../shared/validation/assets';

interface BulkOperationsModalProps {
  selectedAssets: Asset[];
  onClose: () => void;
  onSuccess: () => void;
}

const BulkOperationsModal: React.FC<BulkOperationsModalProps> = ({ 
  selectedAssets, 
  onClose, 
  onSuccess 
}) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
    reset
  } = useForm<AssetInventoryFormData>({
    resolver: zodResolver(assetInventorySchema),
    defaultValues: {
      asset_ids: selectedAssets.map(asset => asset.id),
      action: 'verify',
      status: undefined,
      notes: ''
    }
  });

  const selectedAction = watch('action');

  const onSubmit = async (data: AssetInventoryFormData) => {
    try {
      setLoading(true);
      setError(null);
      
      await assetsApi.performInventory({
        action: data.action as "verify" | "update_status",
        asset_ids: data.asset_ids,
        status: data.status,
        notes: data.notes
      });
      onSuccess();
      onClose();
    } catch (err) {
      setError('Ошибка выполнения операции');
      console.error('Error performing bulk operation:', err);
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
          <h3 className="text-lg font-medium text-gray-900 mb-4">
            Массовые операции ({selectedAssets.length} активов)
          </h3>
          
          {error && (
            <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {/* Список выбранных активов */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Выбранные активы
              </label>
              <div className="max-h-32 overflow-y-auto border border-gray-300 rounded-md p-3 bg-gray-50">
                {selectedAssets.map(asset => (
                  <div key={asset.id} className="text-sm text-gray-700 py-1">
                    {asset.inventory_number} - {asset.name}
                  </div>
                ))}
              </div>
            </div>

            {/* Действие */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Действие *
              </label>
              <select
                {...register('action')}
                className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                  errors.action ? 'border-red-500' : 'border-gray-300'
                }`}
              >
                <option value="verify">Проверить</option>
                <option value="update_status">Обновить статус</option>
              </select>
              {errors.action && (
                <p className="mt-1 text-sm text-red-600">{errors.action.message}</p>
              )}
            </div>

            {/* Статус (только для update_status) */}
            {selectedAction === 'update_status' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Новый статус *
                </label>
                <select
                  {...register('status')}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors.status ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  <option value="">Выберите статус</option>
                  {ASSET_STATUSES.map(status => (
                    <option key={status.value} value={status.value}>
                      {status.label}
                    </option>
                  ))}
                </select>
                {errors.status && (
                  <p className="mt-1 text-sm text-red-600">{errors.status.message}</p>
                )}
              </div>
            )}

            {/* Примечания */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Примечания
              </label>
              <textarea
                {...register('notes')}
                rows={3}
                className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                  errors.notes ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Введите примечания к операции..."
              />
              {errors.notes && (
                <p className="mt-1 text-sm text-red-600">{errors.notes.message}</p>
              )}
            </div>

            {/* Кнопки */}
            <div className="flex justify-end space-x-2 pt-4">
              <button
                type="button"
                onClick={handleClose}
                disabled={loading}
                className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Отмена
              </button>
              <button
                type="submit"
                disabled={loading}
                className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? 'Выполнение...' : 'Выполнить'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default BulkOperationsModal;

