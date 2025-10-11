import { z } from 'zod';

// Схема валидации для создания актива
export const createAssetSchema = z.object({
  name: z.string()
    .min(1, 'Название обязательно')
    .max(255, 'Название не должно превышать 255 символов')
    .trim(),
  inventory_number: z.string()
    .min(1, 'Инвентарный номер обязателен')
    .max(100, 'Инвентарный номер не должен превышать 100 символов')
    .trim()
    .optional()
    .or(z.literal('')),
  type: z.string()
    .min(1, 'Тип обязателен')
    .refine(val => ['server', 'workstation', 'computer', 'monitor', 'application', 'database', 'document', 'network_device', 'other'].includes(val), {
      message: 'Выберите корректный тип актива'
    }),
  class: z.string()
    .min(1, 'Класс обязателен')
    .refine(val => ['hardware', 'software', 'data', 'service'].includes(val), {
      message: 'Выберите корректный класс актива'
    }),
  owner_id: z.string()
    .optional()
    .refine(val => !val || /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(val), {
      message: 'Некорректный формат ID владельца'
    }),
  responsible_user_id: z.string()
    .optional()
    .refine(val => !val || /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(val), {
      message: 'Некорректный формат ID ответственного пользователя'
    }),
  location: z.string()
    .max(500, 'Местоположение не должно превышать 500 символов')
    .optional()
    .or(z.literal('')),
  criticality: z.string()
    .min(1, 'Критичность обязательна')
    .refine(val => ['low', 'medium', 'high'].includes(val), {
      message: 'Выберите корректный уровень критичности'
    }),
  confidentiality: z.string()
    .min(1, 'Конфиденциальность обязательна')
    .refine(val => ['low', 'medium', 'high'].includes(val), {
      message: 'Выберите корректный уровень конфиденциальности'
    }),
  integrity: z.string()
    .min(1, 'Целостность обязательна')
    .refine(val => ['low', 'medium', 'high'].includes(val), {
      message: 'Выберите корректный уровень целостности'
    }),
  availability: z.string()
    .min(1, 'Доступность обязательна')
    .refine(val => ['low', 'medium', 'high'].includes(val), {
      message: 'Выберите корректный уровень доступности'
    }),
  status: z.string()
    .optional()
    .refine(val => !val || ['active', 'in_repair', 'storage', 'decommissioned'].includes(val), {
      message: 'Выберите корректный статус'
    }),
  // Passport fields
  serial_number: z.string().max(255).optional().or(z.literal('')),
  pc_number: z.string().max(100).optional().or(z.literal('')),
  model: z.string().max(255).optional().or(z.literal('')),
  cpu: z.string().max(255).optional().or(z.literal('')),
  ram: z.string().max(100).optional().or(z.literal('')),
  hdd_info: z.string().optional().or(z.literal('')),
  network_card: z.string().max(255).optional().or(z.literal('')),
  optical_drive: z.string().max(255).optional().or(z.literal('')),
  ip_address: z.string()
    .optional()
    .refine(val => !val || val === '' || /^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/.test(val), {
      message: 'Некорректный формат IP адреса'
    })
    .or(z.literal('')),
  mac_address: z.string()
    .optional()
    .refine(val => !val || val === '' || /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/.test(val), {
      message: 'Некорректный формат MAC адреса'
    })
    .or(z.literal('')),
  manufacturer: z.string().max(255).optional().or(z.literal('')),
  purchase_year: z.number()
    .min(1900, 'Год не может быть раньше 1900')
    .max(2100, 'Год не может быть позже 2100')
    .optional()
    .nullable(),
  warranty_until: z.string().optional().or(z.literal('')),
  template_id: z.string().optional().or(z.literal('')),
});

// Схема валидации для обновления актива
export const updateAssetSchema = z.object({
  name: z.string()
    .min(1, 'Название обязательно')
    .max(255, 'Название не должно превышать 255 символов')
    .trim()
    .optional(),
  inventory_number: z.string()
    .max(100, 'Инвентарный номер не должен превышать 100 символов')
    .trim()
    .optional()
    .or(z.literal('')),
  type: z.string()
    .refine(val => !val || ['server', 'workstation', 'computer', 'monitor', 'application', 'database', 'document', 'network_device', 'other'].includes(val), {
      message: 'Выберите корректный тип актива'
    })
    .optional(),
  class: z.string()
    .refine(val => !val || ['hardware', 'software', 'data', 'service'].includes(val), {
      message: 'Выберите корректный класс актива'
    })
    .optional(),
  owner_id: z.string()
    .optional()
    .refine(val => !val || val === '' || /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(val), {
      message: 'Некорректный формат ID владельца'
    }),
  responsible_user_id: z.string()
    .optional()
    .refine(val => !val || val === '' || /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(val), {
      message: 'Некорректный формат ID ответственного пользователя'
    }),
  location: z.string()
    .max(500, 'Местоположение не должно превышать 500 символов')
    .optional(),
  criticality: z.string()
    .refine(val => !val || ['low', 'medium', 'high'].includes(val), {
      message: 'Выберите корректный уровень критичности'
    })
    .optional(),
  confidentiality: z.string()
    .refine(val => !val || ['low', 'medium', 'high'].includes(val), {
      message: 'Выберите корректный уровень конфиденциальности'
    })
    .optional(),
  integrity: z.string()
    .refine(val => !val || ['low', 'medium', 'high'].includes(val), {
      message: 'Выберите корректный уровень целостности'
    })
    .optional(),
  availability: z.string()
    .refine(val => !val || ['low', 'medium', 'high'].includes(val), {
      message: 'Выберите корректный уровень доступности'
    })
    .optional(),
  status: z.string()
    .refine(val => !val || ['active', 'in_repair', 'storage', 'decommissioned'].includes(val), {
      message: 'Выберите корректный статус'
    })
    .optional(),
  // Passport fields
  serial_number: z.string().max(255).optional().or(z.literal('')),
  pc_number: z.string().max(100).optional().or(z.literal('')),
  model: z.string().max(255).optional().or(z.literal('')),
  cpu: z.string().max(255).optional().or(z.literal('')),
  ram: z.string().max(100).optional().or(z.literal('')),
  hdd_info: z.string().optional().or(z.literal('')),
  network_card: z.string().max(255).optional().or(z.literal('')),
  optical_drive: z.string().max(255).optional().or(z.literal('')),
  ip_address: z.string()
    .optional()
    .refine(val => !val || val === '' || /^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/.test(val), {
      message: 'Некорректный формат IP адреса'
    })
    .or(z.literal('')),
  mac_address: z.string()
    .optional()
    .refine(val => !val || val === '' || /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/.test(val), {
      message: 'Некорректный формат MAC адреса'
    })
    .or(z.literal('')),
  manufacturer: z.string().max(255).optional().or(z.literal('')),
  purchase_year: z.number()
    .min(1900, 'Год не может быть раньше 1900')
    .max(2100, 'Год не может быть позже 2100')
    .optional()
    .nullable(),
  warranty_until: z.string().optional().or(z.literal('')),
  template_id: z.string().optional().or(z.literal('')),
});

// Схема валидации для документов актива
export const assetDocumentSchema = z.object({
  document_type: z.string()
    .min(1, 'Тип документа обязателен')
    .refine(val => ['passport', 'transfer_act', 'writeoff_act', 'repair_log', 'other'].includes(val), {
      message: 'Выберите корректный тип документа'
    }),
  file_path: z.string()
    .min(1, 'Путь к файлу обязателен')
    .max(1000, 'Путь к файлу не должен превышать 1000 символов')
});

// Схема валидации для ПО актива
export const assetSoftwareSchema = z.object({
  software_name: z.string()
    .min(1, 'Название ПО обязательно')
    .max(255, 'Название ПО не должно превышать 255 символов')
    .trim(),
  version: z.string()
    .max(100, 'Версия не должна превышать 100 символов')
    .optional()
    .or(z.literal('')),
  installed_at: z.string()
    .optional()
    .refine(val => !val || !isNaN(Date.parse(val)), {
      message: 'Некорректная дата установки'
    })
});

// Схема валидации для инвентаризации
export const assetInventorySchema = z.object({
  asset_ids: z.array(z.string().uuid('Некорректный ID актива'))
    .min(1, 'Выберите хотя бы один актив')
    .max(100, 'Максимум 100 активов за раз'),
  action: z.string()
    .min(1, 'Действие обязательно')
    .refine(val => ['verify', 'update_status'].includes(val), {
      message: 'Выберите корректное действие'
    }),
  status: z.string()
    .refine(val => !val || ['active', 'in_repair', 'storage', 'decommissioned'].includes(val), {
      message: 'Выберите корректный статус'
    })
    .optional(),
  notes: z.string()
    .max(1000, 'Примечания не должны превышать 1000 символов')
    .optional()
});

// Типы для TypeScript
export type CreateAssetFormData = z.infer<typeof createAssetSchema>;
export type UpdateAssetFormData = z.infer<typeof updateAssetSchema>;
export type AssetDocumentFormData = z.infer<typeof assetDocumentSchema>;
export type AssetSoftwareFormData = z.infer<typeof assetSoftwareSchema>;
export type AssetInventoryFormData = z.infer<typeof assetInventorySchema>;

// Вспомогательные функции для валидации
export const validateAssetName = (name: string): string | null => {
  if (!name || name.trim().length === 0) {
    return 'Название обязательно';
  }
  if (name.length > 255) {
    return 'Название не должно превышать 255 символов';
  }
  return null;
};

export const validateUUID = (id: string): boolean => {
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
  return uuidRegex.test(id);
};

export const validateCIAValue = (value: string): boolean => {
  return ['low', 'medium', 'high'].includes(value);
};

export const validateAssetType = (type: string): boolean => {
  return ['server', 'workstation', 'computer', 'monitor', 'application', 'database', 'document', 'network_device', 'other'].includes(type);
};

export const validateAssetClass = (classValue: string): boolean => {
  return ['hardware', 'software', 'data', 'service'].includes(classValue);
};

export const validateAssetStatus = (status: string): boolean => {
  return ['active', 'in_repair', 'storage', 'decommissioned'].includes(status);
};

