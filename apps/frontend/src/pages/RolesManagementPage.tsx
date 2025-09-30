import React, { useState, useEffect } from 'react';
import { rolesApi, Role, RoleWithPermissions, Permission } from '../shared/api/roles';

// Функция для получения иконки модуля
const getModuleIcon = (module: string) => {
  const icons: Record<string, string> = {
    'Документы': '📄',
    'Инциденты': '⚠️',
    'Риски': '📊',
    'Активы': '🛠️',
    'Обучение': '🎓',
    'ИИ': '🤖',
    'Пользователи': '👤',
    'Роли': '🛡️',
    'Аудит': '📜',
    'Дашборд': '📈',
    'Соответствие': '✅',
    'Отчеты': '📋'
  };
  return icons[module] || '📁';
};

const RolesManagementPage: React.FC = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedRole, setSelectedRole] = useState<RoleWithPermissions | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showPermissionsModal, setShowPermissionsModal] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [rolesRaw, permissionsRaw] = await Promise.all([
        rolesApi.getRoles(),
        rolesApi.getPermissions()
      ]);
      // Нормализация ключей в camelCase на случай, если пришли поля в PascalCase
      const rolesData: Role[] = (rolesRaw as any[]).map((r) => ({
        id: r.id ?? r.ID,
        name: r.name ?? r.Name,
        description: r.description ?? r.Description ?? undefined,
        created_at: r.created_at ?? r.CreatedAt,
        updated_at: r.updated_at ?? r.UpdatedAt,
      }));
      const permissionsData: Permission[] = (permissionsRaw as any[]).map((p) => ({
        id: p.id ?? p.ID,
        code: p.code ?? p.Code,
        module: p.module ?? p.Module,
        description: p.description ?? p.Description ?? undefined,
      }));
      setRoles(rolesData);
      setPermissions(permissionsData);
      
      // Debug: проверим структуру данных
      console.log('Загруженные роли (нормализовано):', rolesData);
      console.log('Пример роли:', rolesData[0]);
      console.log('Загруженные права (нормализовано):', permissionsData);
      console.log('Пример права:', permissionsData[0]);
    } catch (error) {
      console.error('Ошибка загрузки данных:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateRole = async (roleData: any) => {
    try {
      await rolesApi.createRole(roleData);
      await loadData();
      setShowCreateModal(false);
    } catch (error) {
      console.error('Ошибка создания роли:', error);
    }
  };

  const handleUpdateRole = async (id: string, roleData: any) => {
    try {
      await rolesApi.updateRole(id, roleData);
      await loadData();
      setShowEditModal(false);
      setSelectedRole(null);
    } catch (error) {
      console.error('Ошибка обновления роли:', error);
    }
  };

  const handleDeleteRole = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить роль?')) {
      try {
        await rolesApi.deleteRole(id);
        await loadData();
      } catch (error) {
        console.error('Ошибка удаления роли:', error);
      }
    }
  };

  const handleViewPermissions = async (role: Role) => {
    try {
      const roleWithPermissions = await rolesApi.getRole(role.id);
      setSelectedRole(roleWithPermissions);
      setShowPermissionsModal(true);
    } catch (error) {
      console.error('Ошибка загрузки роли:', error);
    }
  };

  if (loading) {
    return <div className="p-6">Загрузка...</div>;
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Управление ролями</h1>
        <button
          onClick={() => setShowCreateModal(true)}
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Создать роль
        </button>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Название
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Описание
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Пользователи
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
            {roles.map((role, index) => (
              <tr key={role.id || `role-${index}`}>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">
                    {role.name || '-'}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {role.description || '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                    0
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {role.created_at ? new Date(role.created_at).toLocaleDateString() : '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div className="flex space-x-2">
                    <button
                      onClick={() => {
                        setSelectedRole(role as RoleWithPermissions);
                        setShowEditModal(true);
                      }}
                      className="text-indigo-600 hover:text-indigo-900"
                    >
                      Редактировать
                    </button>
                    <button
                      onClick={() => handleViewPermissions(role)}
                      className="text-green-600 hover:text-green-900"
                    >
                      Права
                    </button>
                    <button
                      onClick={() => handleDeleteRole(role.id)}
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

      {showCreateModal && (
        <CreateRoleModal
          onClose={() => setShowCreateModal(false)}
          onSubmit={handleCreateRole}
          permissions={permissions}
        />
      )}

      {showEditModal && selectedRole && (
        <EditRoleModal
          role={selectedRole}
          onClose={() => {
            setShowEditModal(false);
            setSelectedRole(null);
          }}
          onSubmit={handleUpdateRole}
          permissions={permissions}
        />
      )}

      {showPermissionsModal && selectedRole && (
        <RolePermissionsModal
          role={selectedRole}
          onClose={() => {
            setShowPermissionsModal(false);
            setSelectedRole(null);
          }}
          permissions={permissions}
        />
      )}
    </div>
  );
};

// Компонент создания роли
const CreateRoleModal: React.FC<{
  onClose: () => void;
  onSubmit: (data: any) => void;
  permissions: Permission[];
}> = ({ onClose, onSubmit, permissions }) => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    permission_ids: [] as string[],
    permissionSearch: ''
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Фильтруем null/undefined значения из permission_ids
    const filteredPermissionIds = formData.permission_ids.filter(id => id !== null && id !== undefined && id !== '');
    
    const submitData = {
      ...formData,
      permission_ids: filteredPermissionIds
    };
    
    console.log('Отправляемые данные роли:', submitData);
    onSubmit(submitData);
  };

  const groupedPermissions = permissions.reduce((acc, perm) => {
    const module = perm.module || 'Общие';
    if (!acc[module]) {
      acc[module] = [];
    }
    acc[module].push(perm);
    return acc;
  }, {} as Record<string, Permission[]>);

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-96 max-h-96 overflow-y-auto">
        <h2 className="text-xl font-bold mb-4">Создать роль</h2>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Название
            </label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full border border-gray-300 rounded px-3 py-2"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Описание
            </label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="w-full border border-gray-300 rounded px-3 py-2"
              rows={3}
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Права
            </label>
            <input
              type="text"
              placeholder="Поиск по правам..."
              className="w-full border border-gray-300 rounded px-3 py-2 mb-3 text-sm"
              onChange={() => {}}
            />
            <div className="space-y-4 max-h-48 overflow-y-auto">
              {Object.entries(groupedPermissions).map(([module, perms]) => (
                <div key={module}>
                  <h4 className="font-medium text-gray-900 mb-2 flex items-center">
                    <span className="mr-2">{getModuleIcon(module)}</span>
                    {module}
                    <label className="ml-auto flex items-center text-sm text-blue-600">
                      <input
                        type="checkbox"
                        checked={perms.every(perm => formData.permission_ids.includes(perm.id))}
                        onChange={(e) => {
                          if (e.target.checked) {
                            // Выбрать все права в модуле
                            const modulePermissionIds = perms.filter(perm => perm.id).map(perm => perm.id);
                            const newPermissionIds = [...new Set([...formData.permission_ids, ...modulePermissionIds])];
                            setFormData({
                              ...formData,
                              permission_ids: newPermissionIds
                            });
                          } else {
                            // Снять все права в модуле
                            const modulePermissionIds = perms.filter(perm => perm.id).map(perm => perm.id);
                            const newPermissionIds = formData.permission_ids.filter(id => !modulePermissionIds.includes(id));
                            setFormData({
                              ...formData,
                              permission_ids: newPermissionIds
                            });
                          }
                        }}
                        className="mr-1"
                      />
                      Выбрать всё
                    </label>
                  </h4>
                  <div className="space-y-1 ml-4">
                    {perms.map((perm, index) => {
                      const pid = perm.id;
                      const pcode = perm.code;
                      const pdesc = perm.description;
                      return (
                        <label key={pid || `perm-${index}`} className="flex items-center">
                          <input
                            type="checkbox"
                            checked={pid ? formData.permission_ids.includes(pid) : false}
                            onChange={(e) => {
                              if (pid) {
                                if (e.target.checked) {
                                  setFormData({
                                    ...formData,
                                    permission_ids: [...new Set([...formData.permission_ids, pid])]
                                  });
                                } else {
                                  setFormData({
                                    ...formData,
                                    permission_ids: formData.permission_ids.filter(id => id !== pid)
                                  });
                                }
                              }
                            }}
                            className="mr-2"
                          />
                          <span className="text-sm">{pcode}</span>
                          {pdesc && (
                            <span className="text-xs text-gray-500 ml-2">
                              - {pdesc}
                            </span>
                          )}
                        </label>
                      );
                    })}
                  </div>
                </div>
              ))}
            </div>
          </div>
          <div className="flex justify-end space-x-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-gray-600 border border-gray-300 rounded hover:bg-gray-50"
            >
              Отмена
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Создать
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Компонент редактирования роли
const EditRoleModal: React.FC<{
  role: RoleWithPermissions;
  onClose: () => void;
  onSubmit: (id: string, data: any) => void;
  permissions: Permission[];
}> = ({ role, onClose, onSubmit, permissions }) => {
  const [formData, setFormData] = useState({
    name: role.name,
    description: role.description || '',
    permission_ids: role.permissions || []
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(role.id, formData);
  };

  const groupedPermissions = permissions.reduce((acc, perm) => {
    const module = perm.module || (perm as any).Module || 'Общие';
    if (!acc[module]) {
      acc[module] = [];
    }
    acc[module].push(perm);
    return acc;
  }, {} as Record<string, Permission[]>);

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-96 max-h-96 overflow-y-auto">
        <h2 className="text-xl font-bold mb-4">Редактировать роль</h2>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Название
            </label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full border border-gray-300 rounded px-3 py-2"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Описание
            </label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="w-full border border-gray-300 rounded px-3 py-2"
              rows={3}
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Права
            </label>
            <input
              type="text"
              placeholder="Поиск по правам..."
              className="w-full border border-gray-300 rounded px-3 py-2 mb-3 text-sm"
              onChange={() => {}}
            />
            <div className="space-y-4 max-h-48 overflow-y-auto">
              {Object.entries(groupedPermissions).map(([module, perms]) => (
                <div key={module}>
                  <h4 className="font-medium text-gray-900 mb-2 flex items-center">
                    <span className="mr-2">{getModuleIcon(module)}</span>
                    {module}
                    <label className="ml-auto flex items-center text-sm text-blue-600">
                      <input
                        type="checkbox"
                        checked={perms.every(perm => formData.permission_ids.includes(perm.id))}
                        onChange={(e) => {
                          if (e.target.checked) {
                            // Выбрать все права в модуле
                            const modulePermissionIds = perms.filter(perm => perm.id).map(perm => perm.id);
                            const newPermissionIds = [...new Set([...formData.permission_ids, ...modulePermissionIds])];
                            setFormData({
                              ...formData,
                              permission_ids: newPermissionIds
                            });
                          } else {
                            // Снять все права в модуле
                            const modulePermissionIds = perms.filter(perm => perm.id).map(perm => perm.id);
                            const newPermissionIds = formData.permission_ids.filter(id => !modulePermissionIds.includes(id));
                            setFormData({
                              ...formData,
                              permission_ids: newPermissionIds
                            });
                          }
                        }}
                        className="mr-1"
                      />
                      Выбрать всё
                    </label>
                  </h4>
                  <div className="space-y-1 ml-4">
                    {perms.map((perm, index) => {
                      const pid = perm.id;
                      const pcode = perm.code;
                      const pdesc = perm.description;
                      return (
                        <label key={pid || `perm-${index}`} className="flex items-center">
                          <input
                            type="checkbox"
                            checked={pid ? formData.permission_ids.includes(pid) : false}
                            onChange={(e) => {
                              if (pid) {
                                if (e.target.checked) {
                                  setFormData({
                                    ...formData,
                                    permission_ids: [...new Set([...formData.permission_ids, pid])]
                                  });
                                } else {
                                  setFormData({
                                    ...formData,
                                    permission_ids: formData.permission_ids.filter(id => id !== pid)
                                  });
                                }
                              }
                            }}
                            className="mr-2"
                          />
                          <span className="text-sm">{pcode}</span>
                          {pdesc && (
                            <span className="text-xs text-gray-500 ml-2">
                              - {pdesc}
                            </span>
                          )}
                        </label>
                      );
                    })}
                  </div>
                </div>
              ))}
            </div>
          </div>
          <div className="flex justify-end space-x-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-gray-600 border border-gray-300 rounded hover:bg-gray-50"
            >
              Отмена
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Сохранить
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Компонент просмотра прав роли
const RolePermissionsModal: React.FC<{
  role: RoleWithPermissions;
  onClose: () => void;
  permissions: Permission[];
}> = ({ role, onClose, permissions }) => {
  const groupedPermissions = permissions.reduce((acc, perm) => {
    const module = perm.module || 'Общие';
    if (!acc[module]) {
      acc[module] = [];
    }
    acc[module].push(perm);
    return acc;
  }, {} as Record<string, Permission[]>);

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-96 max-h-96 overflow-y-auto">
        <h2 className="text-xl font-bold mb-4">Права роли: {role.name}</h2>
        <div className="space-y-4">
          {Object.entries(groupedPermissions).map(([module, perms]) => (
            <div key={module}>
              <h4 className="font-medium text-gray-900 mb-2">{module}</h4>
              <div className="space-y-1 ml-4">
                {perms.map((perm, index) => {
                  const hasPermission = role.permissions.includes(perm.code);
                  return (
                    <div
                      key={perm.id || `perm-${index}`}
                      className={`flex items-center p-2 rounded ${
                        hasPermission ? 'bg-green-50 text-green-800' : 'bg-gray-50 text-gray-500'
                      }`}
                    >
                      <span className="text-sm font-medium">{perm.code}</span>
                      {perm.description && (
                        <span className="text-xs ml-2">
                          - {perm.description}
                        </span>
                      )}
                      {hasPermission && (
                        <span className="ml-auto text-green-600">✓</span>
                      )}
                    </div>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
        <div className="flex justify-end mt-6">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            Закрыть
          </button>
        </div>
      </div>
    </div>
  );
};

export default RolesManagementPage;
