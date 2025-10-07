import React from 'react'
import { usePermissions } from '../hooks/usePermissions'

interface PermissionGuardProps {
  children: React.ReactNode
  permission?: string
  permissions?: string[]
  requireAll?: boolean
  role?: string
  roles?: string[]
  requireAllRoles?: boolean
  fallback?: React.ReactNode
}

/**
 * Компонент для условного рендеринга на основе прав доступа
 * 
 * @param children - Компоненты для рендеринга при наличии прав
 * @param permission - Конкретное право для проверки
 * @param permissions - Массив прав для проверки
 * @param requireAll - Требовать все права (по умолчанию false - достаточно одного)
 * @param role - Конкретная роль для проверки
 * @param roles - Массив ролей для проверки
 * @param requireAllRoles - Требовать все роли (по умолчанию false - достаточно одной)
 * @param fallback - Компонент для рендеринга при отсутствии прав
 */
export const PermissionGuard: React.FC<PermissionGuardProps> = ({
  children,
  permission,
  permissions,
  requireAll = false,
  role,
  roles,
  requireAllRoles = false,
  fallback = null
}) => {
  const { 
    hasPermission, 
    hasAnyPermission, 
    hasAllPermissions,
    hasRole,
    hasAnyRole,
    hasAllRoles: hasAllRolesCheck
  } = usePermissions()

  // Проверка прав
  let hasRequiredPermissions = true

  if (permission) {
    hasRequiredPermissions = hasPermission(permission)
  } else if (permissions && permissions.length > 0) {
    hasRequiredPermissions = requireAll 
      ? hasAllPermissions(permissions)
      : hasAnyPermission(permissions)
  }

  // Проверка ролей
  let hasRequiredRoles = true

  if (role) {
    hasRequiredRoles = hasRole(role)
  } else if (roles && roles.length > 0) {
    hasRequiredRoles = requireAllRoles
      ? hasAllRolesCheck(roles)
      : hasAnyRole(roles)
  }

  // Если указаны и права, и роли, то должны выполняться оба условия
  const hasAccess = hasRequiredPermissions && hasRequiredRoles

  return hasAccess ? <>{children}</> : <>{fallback}</>
}

// Экспорт по умолчанию для удобства
export default PermissionGuard
