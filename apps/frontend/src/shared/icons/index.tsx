import React from 'react'
import type { SvgIconProps } from '@mui/material/SvgIcon'
import {
  SpaceDashboardRounded,
  PeopleRounded,
  ImportantDevicesRounded,
  WarningAmberRounded,
  DescriptionRounded,
  ReportProblemRounded,
  SchoolRounded,
  GavelRounded,
  PsychologyAltRounded,
  SmartToyRounded,
  ManageAccountsRounded,
  SecurityRounded,
  BusinessRounded,
  LogoutRounded,
} from '@mui/icons-material'

const withBrandIcon = (IconComponent: React.ComponentType<SvgIconProps>) => {
  const BrandedIcon: React.FC<SvgIconProps> = ({ sx, fontSize = 'medium', ...props }) => (
    <IconComponent
      fontSize={fontSize}
      {...props}
      sx={{
        color: 'inherit',
        transition: 'color 0.2s ease',
        ...sx,
      }}
    />
  )

  const name = IconComponent.displayName || IconComponent.name || 'Icon'
  BrandedIcon.displayName = `CorporateIcon(${name})`
  return BrandedIcon
}

export const DashboardIcon = withBrandIcon(SpaceDashboardRounded)
export const UsersIcon = withBrandIcon(PeopleRounded)
export const AssetsIcon = withBrandIcon(ImportantDevicesRounded)
export const RisksIcon = withBrandIcon(WarningAmberRounded)
export const DocumentsIcon = withBrandIcon(DescriptionRounded)
export const IncidentsIcon = withBrandIcon(ReportProblemRounded)
export const TrainingIcon = withBrandIcon(SchoolRounded)
export const ComplianceIcon = withBrandIcon(GavelRounded)
export const AIProvidersIcon = withBrandIcon(PsychologyAltRounded)
export const AIQueryIcon = withBrandIcon(SmartToyRounded)
export const AdminUsersIcon = withBrandIcon(ManageAccountsRounded)
export const RolesIcon = withBrandIcon(SecurityRounded)
export const OrganizationsIcon = withBrandIcon(BusinessRounded)
export const LogoutIcon = withBrandIcon(LogoutRounded)
