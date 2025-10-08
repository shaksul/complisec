import type { ReactNode } from 'react'
import { DashboardIcon, UsersIcon, AssetsIcon, RisksIcon, DocumentsIcon, IncidentsIcon, TrainingIcon, ComplianceIcon, AIProvidersIcon, AIQueryIcon, RolesIcon, OrganizationsIcon } from './icons'
import { Folder } from '@mui/icons-material'

export interface NavigationItem {
  label: string
  to: string
  icon: ReactNode
  exact?: boolean
  permission?: string
  permissions?: string[]
  requireAll?: boolean
}

export const PRIMARY_NAVIGATION: NavigationItem[] = [
  { label: 'Панель мониторинга', to: '/dashboard', icon: <DashboardIcon /> },
  { label: 'Пользователи', to: '/users', icon: <UsersIcon />, permission: 'users.view' },
  { label: 'Активы', to: '/assets', icon: <AssetsIcon />, permission: 'asset.view' },
  { label: 'Риски', to: '/risks', icon: <RisksIcon />, permission: 'risk.view' },
  { label: 'Документы', to: '/documents', icon: <DocumentsIcon />, permission: 'document.read' },
  { label: 'Файловое хранилище', to: '/file-documents', icon: <Folder />, permission: 'document.read' },
  { label: 'Инциденты', to: '/incidents', icon: <IncidentsIcon />, permission: 'incidents.view' },
  { label: 'Обучение', to: '/training', icon: <TrainingIcon />, permission: 'training.view' },
  { label: 'Комплаенс', to: '/compliance', icon: <ComplianceIcon />, permission: 'compliance.view' },
  { label: 'AI-провайдеры', to: '/ai/providers', icon: <AIProvidersIcon />, permission: 'ai.providers.view' },
  { label: 'AI-аналитика', to: '/ai/query', icon: <AIQueryIcon />, permission: 'ai.query.view' },
]

export const ADMIN_NAVIGATION: NavigationItem[] = [
  { label: 'Роли и полномочия', to: '/admin/roles', icon: <RolesIcon />, permission: 'roles.view' },
  { label: 'Организации', to: '/admin/organizations', icon: <OrganizationsIcon />, permission: 'organizations.manage' },
]

