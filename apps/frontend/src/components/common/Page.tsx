import React from 'react'
import { Container, type ContainerProps, Box, Typography, type TypographyProps, Paper, type PaperProps, Stack } from '@mui/material'

export const PageContainer: React.FC<ContainerProps> = ({ children, maxWidth = 'xl', sx, ...props }) => (
  <Container
    maxWidth={maxWidth}
    {...props}
    sx={{
      py: 4,
      display: 'flex',
      flexDirection: 'column',
      gap: 3,
      ...sx,
    }}
  >
    {children}
  </Container>
)

interface PageHeaderProps {
  title: string
  subtitle?: string
  actions?: React.ReactNode
  titleProps?: TypographyProps
}

export const PageHeader: React.FC<PageHeaderProps> = ({ title, subtitle, actions, titleProps }) => (
  <Box display="flex" alignItems={{ xs: 'stretch', md: 'center' }} flexDirection={{ xs: 'column', md: 'row' }} gap={2}>
    <Box flexGrow={1}>
      <Typography variant="h4" {...titleProps}>
        {title}
      </Typography>
      {subtitle && (
        <Typography variant="body1" color="text.secondary" mt={0.5}>
          {subtitle}
        </Typography>
      )}
    </Box>
    {actions && (
      <Box display="flex" alignItems="center" gap={1}>
        {actions}
      </Box>
    )}
  </Box>
)

interface SectionCardProps extends PaperProps {
  title?: string
  description?: string
  action?: React.ReactNode
}

export const SectionCard: React.FC<SectionCardProps> = ({ title, description, action, children, sx, ...props }) => (
  <Paper
    elevation={0}
    {...props}
    sx={{
      p: 3,
      borderRadius: 3,
      display: 'flex',
      flexDirection: 'column',
      gap: 2.5,
      ...sx,
    }}
  >
    {(title || description || action) && (
      <Box display="flex" flexDirection={{ xs: 'column', sm: 'row' }} gap={1.5} alignItems={{ xs: 'flex-start', sm: 'center' }}>
        <Box flexGrow={1}>
          {title && (
            <Typography variant="h6" component="h3">
              {title}
            </Typography>
          )}
          {description && (
            <Typography variant="body2" color="text.secondary" mt={0.5}>
              {description}
            </Typography>
          )}
        </Box>
        {action && <Box>{action}</Box>}
      </Box>
    )}
    <Stack spacing={2}>{children}</Stack>
  </Paper>
)
