import { useAuth } from '../contexts/AuthContext'

/**
 * Хук для проверки прав доступа пользователя
 */
export function usePermissions() {
  const { user } = useAuth()

  /**
   * Проверяет, есть ли у пользователя конкретное право
   */
  const hasPermission = (permission: string): boolean => {
    if (!user || !user.permissions) {
      return false
    }
    return user.permissions.includes(permission)
  }

  /**
   * Проверяет, есть ли у пользователя хотя бы одно из указанных прав
   */
  const hasAnyPermission = (permissions: string[]): boolean => {
    if (!user || !user.permissions || permissions.length === 0) {
      return false
    }
    return permissions.some(permission => user.permissions.includes(permission))
  }

  /**
   * Проверяет, есть ли у пользователя все указанные права
   */
  const hasAllPermissions = (permissions: string[]): boolean => {
    if (!user || !user.permissions || permissions.length === 0) {
      return false
    }
    return permissions.every(permission => user.permissions.includes(permission))
  }

  /**
   * Проверяет, есть ли у пользователя конкретная роль
   */
  const hasRole = (role: string): boolean => {
    if (!user || !user.roles) {
      return false
    }
    return user.roles.includes(role)
  }

  /**
   * Проверяет, есть ли у пользователя хотя бы одна из указанных ролей
   */
  const hasAnyRole = (roles: string[]): boolean => {
    if (!user || !user.roles || roles.length === 0) {
      return false
    }
    return roles.some(role => user.roles.includes(role))
  }

  /**
   * Проверяет, есть ли у пользователя все указанные роли
   */
  const hasAllRoles = (roles: string[]): boolean => {
    if (!user || !user.roles || roles.length === 0) {
      return false
    }
    return roles.every(role => user.roles.includes(role))
  }

  /**
   * Проверяет, является ли пользователь администратором
   */
  const isAdmin = (): boolean => {
    return hasRole('Admin') || hasRole('Администратор')
  }

  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    hasRole,
    hasAnyRole,
    hasAllRoles,
    isAdmin,
    permissions: user?.permissions || [],
    roles: user?.roles || []
  }
}
