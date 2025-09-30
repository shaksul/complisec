import React, { useState, useEffect } from 'react'
import { usersApi, User } from '../shared/api/users'
import { rolesApi, Role } from '../shared/api/roles'
import Pagination from '../components/Pagination'

const UsersManagementPage: React.FC = () => {
  const [users, setUsers] = useState<User[]>([])
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [showRolesModal, setShowRolesModal] = useState(false)

  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [hasNext, setHasNext] = useState(false)
  const [hasPrev, setHasPrev] = useState(false)
  const [total, setTotal] = useState(0)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async (page: number = currentPage) => {
    try {
      setLoading(true)
      const [usersResponse, rolesData] = await Promise.all([
        usersApi.getUsersPaginated(page, 20),
        rolesApi.getRoles()
      ])

      const normalizedUsers = usersResponse.data.map((user) => ({
        ...user,
        roles: user.roles || []
      }))

      setUsers(normalizedUsers)
      setRoles(rolesData)
      setCurrentPage(usersResponse.pagination.page)
      setTotalPages(usersResponse.pagination.total_pages)
      setHasNext(usersResponse.pagination.has_next)
      setHasPrev(usersResponse.pagination.has_prev)
      setTotal(usersResponse.pagination.total)
    } catch (error) {
      console.error('Ошибка загрузки данных:', error)
    } finally {
      setLoading(false)
    }
  }

  const handlePageChange = (page: number) => {
    setCurrentPage(page)
    loadData(page)
  }

  const handleCreateUser = async (userData: any) => {
    try {
      await usersApi.createUser(userData)
      await loadData(1)
      setCurrentPage(1)
      setShowCreateModal(false)
      alert('Пользователь успешно создан')
    } catch (error) {
      console.error('Ошибка создания пользователя:', error)
      alert('Ошибка создания пользователя: ' + (error as Error).message)
    }
  }

  const handleUpdateUser = async (id: string, userData: any) => {
    try {
      await usersApi.updateUser(id, userData)
      await loadData(currentPage)
      setShowEditModal(false)
      setSelectedUser(null)
    } catch (error) {
      console.error('Ошибка обновления пользователя:', error)
      alert('Ошибка обновления пользователя: ' + (error as Error).message)
    }
  }

  const handleDeleteUser = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить пользователя?')) {
      try {
        await usersApi.deleteUser(id)
        await loadData(currentPage)
      } catch (error) {
        console.error('Ошибка удаления пользователя:', error)
        alert('Ошибка удаления пользователя: ' + (error as Error).message)
      }
    }
  }

  if (loading) {
    return <div className="p-6">Загрузка...</div>
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Управление пользователями</h1>
        <button
          onClick={() => setShowCreateModal(true)}
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Добавить пользователя
        </button>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Пользователь
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Email
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Статус
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
            {users.map((user, index) => (
              <tr key={user.id || `user-${index}`}>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">
                    {[user.first_name, user.last_name].filter(Boolean).join(' ') || '—'}
                  </div>
                  {user.roles && user.roles.length > 0 && (
                    <div className="text-xs text-gray-500">
                      {user.roles.join(', ')}
                    </div>
                  )}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {user.email}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                      user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {user.is_active ? 'Активен' : 'Заблокирован'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {user.created_at ? new Date(user.created_at).toLocaleDateString() : '—'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div className="flex space-x-2">
                    <button
                      onClick={() => {
                        setSelectedUser(user)
                        setShowEditModal(true)
                      }}
                      className="text-indigo-600 hover:text-indigo-900"
                    >
                      Редактировать
                    </button>
                    <button
                      onClick={() => {
                        setSelectedUser(user)
                        setShowRolesModal(true)
                      }}
                      className="text-green-600 hover:text-green-900"
                    >
                      Роли
                    </button>
                    <button
                      onClick={() => handleDeleteUser(user.id)}
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

      <Pagination
        currentPage={currentPage}
        totalPages={totalPages}
        hasNext={hasNext}
        hasPrev={hasPrev}
        onPageChange={handlePageChange}
      />

      {showCreateModal && (
        <CreateUserModal
          onClose={() => setShowCreateModal(false)}
          onSubmit={handleCreateUser}
          roles={roles}
        />
      )}

      {showEditModal && selectedUser && (
        <EditUserModal
          user={selectedUser}
          onClose={() => {
            setShowEditModal(false)
            setSelectedUser(null)
          }}
          onSubmit={handleUpdateUser}
          roles={roles}
        />
      )}

      {showRolesModal && selectedUser && (
        <UserRolesModal
          user={selectedUser}
          onClose={() => {
            setShowRolesModal(false)
            setSelectedUser(null)
          }}
          roles={roles}
        />
      )}
    </div>
  )
}

const CreateUserModal: React.FC<{ onClose: () => void; onSubmit: (data: any) => void; roles: Role[] }> = ({ onClose, onSubmit, roles }) => {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    role_ids: [] as string[]
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.email || !formData.password || !formData.first_name || !formData.last_name) {
      alert('Пожалуйста, заполните все обязательные поля')
      return
    }

    if (formData.password.length < 6) {
      alert('Пароль должен содержать минимум 6 символов')
      return
    }

    // Фильтруем null значения из role_ids
    const filteredRoleIds = formData.role_ids.filter(id => id !== null && id !== undefined && id !== '')
    
    if (filteredRoleIds.length === 0) {
      alert('Выберите хотя бы одну роль')
      return
    }

    // Отправляем данные с отфильтрованными role_ids
    const submitData = {
      ...formData,
      role_ids: filteredRoleIds
    }
    console.log('Отправляемые данные:', submitData)
    onSubmit(submitData)
  }

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-96">
        <h2 className="text-xl font-bold mb-4">Создать пользователя</h2>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
            <input
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              className="w-full border border-gray-300 rounded px-3 py-2"
              autoComplete="email"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">Пароль</label>
            <input
              type="password"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              className="w-full border border-gray-300 rounded px-3 py-2"
              autoComplete="new-password"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">Имя</label>
            <input
              type="text"
              value={formData.first_name}
              onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
              className="w-full border border-gray-300 rounded px-3 py-2"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">Фамилия</label>
            <input
              type="text"
              value={formData.last_name}
              onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
              className="w-full border border-gray-300 rounded px-3 py-2"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">Роли</label>
            <div className="space-y-2">
              {roles.map((role, index) => (
                <label key={role.id || `role-${index}`} className="flex items-center">
                  <input
                    type="checkbox"
                    checked={formData.role_ids.includes(role.id)}
                    onChange={(e) => {
                      if (e.target.checked && role.id) {
                        setFormData({ ...formData, role_ids: [...formData.role_ids, role.id] })
                      } else if (role.id) {
                        setFormData({
                          ...formData,
                          role_ids: formData.role_ids.filter((id) => id !== role.id)
                        })
                      }
                    }}
                    className="mr-2"
                  />
                  {role.name}
                </label>
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
            <button type="submit" className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
              Создать
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

const EditUserModal: React.FC<{
  user: User;
  onClose: () => void;
  onSubmit: (id: string, data: any) => void;
  roles: Role[];
}> = ({ user, onClose, onSubmit, roles }) => {
  const [formData, setFormData] = useState({
    first_name: user.first_name || '',
    last_name: user.last_name || '',
    is_active: user.is_active,
    role_ids: [] as string[]
  });

  useEffect(() => {
    // Загружаем роли пользователя
    const loadUserRoles = async () => {
      try {
        const userRoles = await rolesApi.getUserRoles(user.id);
        setFormData(prev => ({ ...prev, role_ids: userRoles }));
      } catch (error) {
        console.error('Ошибка загрузки ролей пользователя:', error);
      }
    };
    loadUserRoles();
  }, [user.id]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Фильтруем null значения из role_ids
    const filteredRoleIds = formData.role_ids.filter(id => id !== null && id !== undefined && id !== '')
    
    onSubmit(user.id, {
      ...formData,
      role_ids: filteredRoleIds
    });
  };

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-96">
        <h2 className="text-xl font-bold mb-4">Редактировать пользователя</h2>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              type="email"
              value={user.email}
              disabled
              className="w-full border border-gray-300 rounded px-3 py-2 bg-gray-100"
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Имя
            </label>
            <input
              type="text"
              value={formData.first_name}
              onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
              className="w-full border border-gray-300 rounded px-3 py-2"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Фамилия
            </label>
            <input
              type="text"
              value={formData.last_name}
              onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
              className="w-full border border-gray-300 rounded px-3 py-2"
              required
            />
          </div>
          <div className="mb-4">
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={formData.is_active}
                onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
                className="mr-2"
              />
              Активен
            </label>
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Роли
            </label>
            <div className="space-y-2">
              {roles.map((role) => (
                <label key={role.id} className="flex items-center">
                  <input
                    type="checkbox"
                    checked={formData.role_ids.includes(role.id)}
                    onChange={(e) => {
                      if (e.target.checked && role.id) {
                        setFormData({
                          ...formData,
                          role_ids: [...formData.role_ids, role.id]
                        });
                      } else if (role.id) {
                        setFormData({
                          ...formData,
                          role_ids: formData.role_ids.filter(id => id !== role.id)
                        });
                      }
                    }}
                    className="mr-2"
                  />
                  {role.name}
                </label>
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

// Компонент управления ролями пользователя
const UserRolesModal: React.FC<{
  user: User;
  onClose: () => void;
  roles: Role[];
}> = ({ user, onClose, roles }) => {
  const [userRoles, setUserRoles] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadUserRoles = async () => {
      try {
        const roles = await rolesApi.getUserRoles(user.id);
        setUserRoles(roles);
      } catch (error) {
        console.error('Ошибка загрузки ролей пользователя:', error);
      } finally {
        setLoading(false);
      }
    };
    loadUserRoles();
  }, [user.id]);

  const handleRoleToggle = async (roleId: string, assigned: boolean) => {
    try {
      if (assigned) {
        await rolesApi.removeRoleFromUser(user.id, roleId);
        setUserRoles(prev => prev.filter(id => id !== roleId));
      } else {
        await rolesApi.assignRoleToUser(user.id, roleId);
        setUserRoles(prev => [...prev, roleId]);
      }
    } catch (error) {
      console.error('Ошибка изменения роли:', error);
    }
  };

  if (loading) {
    return (
      <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
        <div className="bg-white rounded-lg p-6">
          <div>Загрузка...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-96">
        <h2 className="text-xl font-bold mb-4">Роли пользователя: {user.first_name} {user.last_name}</h2>
        <div className="space-y-2">
          {roles.map((role) => {
            const isAssigned = userRoles.includes(role.id);
            return (
              <label key={role.id} className="flex items-center justify-between">
                <div>
                  <div className="font-medium">{role.name}</div>
                  {role.description && (
                    <div className="text-sm text-gray-500">{role.description}</div>
                  )}
                </div>
                <input
                  type="checkbox"
                  checked={isAssigned}
                  onChange={(e) => handleRoleToggle(role.id, isAssigned)}
                  className="ml-4"
                />
              </label>
            );
          })}
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

export default UsersManagementPage;
