import { api as apiClient } from './client';

export interface Role {
  id: string;
  name: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface RoleWithPermissions extends Role {
  permissions: string[];
}

export interface Permission {
  id: string;
  code: string;
  module: string;
  description?: string;
}

export interface CreateRoleRequest {
  name: string;
  description: string;
  permission_ids: string[];
}

export interface UpdateRoleRequest {
  name?: string;
  description?: string;
  permission_ids?: string[];
}

export interface UserRoleRequest {
  user_id: string;
  role_id: string;
}

export const rolesApi = {
  // Получить все роли
  getRoles: async (): Promise<Role[]> => {
    const response = await apiClient.get('/roles');
    return response.data.data;
  },

  // Получить роль по ID
  getRole: async (id: string): Promise<RoleWithPermissions> => {
    const response = await apiClient.get(`/roles/${id}`);
    return response.data.data;
  },

  // Создать роль
  createRole: async (data: CreateRoleRequest): Promise<Role> => {
    const response = await apiClient.post('/roles', data);
    return response.data.data;
  },

  // Обновить роль
  updateRole: async (id: string, data: UpdateRoleRequest): Promise<void> => {
    await apiClient.put(`/roles/${id}`, data);
  },

  // Удалить роль
  deleteRole: async (id: string): Promise<void> => {
    await apiClient.delete(`/roles/${id}`);
  },

  // Получить все права
  getPermissions: async (): Promise<Permission[]> => {
    const response = await apiClient.get('/permissions');
    return response.data.data;
  },

  // Получить пользователей роли
  getRoleUsers: async (roleId: string): Promise<any[]> => {
    const response = await apiClient.get(`/roles/${roleId}/users`);
    return response.data.data;
  },

  // Назначить роль пользователю
  assignRoleToUser: async (userId: string, roleId: string): Promise<void> => {
    await apiClient.post(`/users/${userId}/roles`, { user_id: userId, role_id: roleId });
  },

  // Убрать роль у пользователя
  removeRoleFromUser: async (userId: string, roleId: string): Promise<void> => {
    await apiClient.delete(`/users/${userId}/roles/${roleId}`);
  },

  // Получить роли пользователя
  getUserRoles: async (userId: string): Promise<string[]> => {
    const response = await apiClient.get(`/users/${userId}/roles`);
    return response.data.data;
  },
};
