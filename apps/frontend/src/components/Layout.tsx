import React, { useState } from 'react'
import {
  AppBar,
  Avatar,
  Box,
  CssBaseline,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  ListSubheader,
  Menu,
  MenuItem,
  Toolbar,
  Typography,
} from '@mui/material'
import { MenuRounded } from '@mui/icons-material'
import { useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { PRIMARY_NAVIGATION, ADMIN_NAVIGATION } from '../shared/navigation'
import { usePermissions } from '../hooks/usePermissions'
import { LogoutIcon } from '../shared/icons'

const drawerWidth = 264

interface LayoutProps {
  children: React.ReactNode
}

const NavigationList: React.FC<{
  items: typeof PRIMARY_NAVIGATION
  currentPath: string
  onNavigate: (url: string) => void
}> = ({ items, currentPath, onNavigate }) => {
  const { hasPermission, hasAnyPermission, hasAllPermissions } = usePermissions()

  // Фильтруем элементы навигации по правам доступа
  const filteredItems = items.filter(item => {
    // Если нет требований к правам, показываем элемент
    if (!item.permission && !item.permissions) {
      return true
    }

    // Проверяем конкретное право
    if (item.permission) {
      return hasPermission(item.permission)
    }

    // Проверяем массив прав
    if (item.permissions && item.permissions.length > 0) {
      return item.requireAll 
        ? hasAllPermissions(item.permissions)
        : hasAnyPermission(item.permissions)
    }

    return true
  })

  return (
    <List disablePadding sx={{ px: 1 }}>
      {filteredItems.map((item) => (
        <ListItem key={item.to} disablePadding sx={{ mb: 0.5 }}>
          <ListItemButton
            selected={currentPath.startsWith(item.to)}
            onClick={() => onNavigate(item.to)}
            sx={{ px: 2, py: 1.1 }}
          >
            <ListItemIcon sx={{ minWidth: 40, color: 'inherit' }}>{item.icon}</ListItemIcon>
            <ListItemText
              primary={item.label}
              primaryTypographyProps={{ variant: 'body2', fontWeight: 600 }}
            />
          </ListItemButton>
        </ListItem>
      ))}
    </List>
  )
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [mobileOpen, setMobileOpen] = useState(false)
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuth()

  const drawerContent = (
    <Box height="100%" display="flex" flexDirection="column">
      <Toolbar sx={{ px: 3, py: 2.5 }}>
        <Box>
          <Typography variant="h6" fontWeight={700} letterSpacing="0.08em">
            RISKNEXUS
          </Typography>
          <Typography variant="caption" color="text.secondary">
            Центр управления безопасностью
          </Typography>
        </Box>
      </Toolbar>

      <Divider sx={{ mx: 2, mb: 1 }} />

      <Box flexGrow={1} overflow="auto" sx={{ pb: 2 }}>
        <NavigationList
          items={PRIMARY_NAVIGATION}
          currentPath={location.pathname}
          onNavigate={(url) => {
            navigate(url)
            setMobileOpen(false)
          }}
        />

        <Divider sx={{ mx: 2, my: 2 }} />
        <ListSubheader
          inset
          sx={{
            px: 3,
            py: 1.5,
            fontSize: 12,
            fontWeight: 700,
            textTransform: 'uppercase',
            color: 'text.secondary',
          }}
        >
          Администрирование
        </ListSubheader>
        <NavigationList
          items={ADMIN_NAVIGATION}
          currentPath={location.pathname}
          onNavigate={(url) => {
            navigate(url)
            setMobileOpen(false)
          }}
        />
      </Box>

      <Box
        px={3}
        py={2}
        sx={{ borderTop: (theme) => `1px solid ${theme.palette.divider}` }}
      >
        <Typography variant="caption" color="text.secondary">
          © {new Date().getFullYear()} RiskNexus Platform
        </Typography>
      </Box>
    </Box>
  )

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar
        position="fixed"
        color="default"
        elevation={0}
        sx={{
          backgroundColor: 'background.paper',
          color: 'text.primary',
          borderBottom: (theme) => `1px solid ${theme.palette.divider}`,
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          ml: { sm: `${drawerWidth}px` },
        }}
      >
        <Toolbar sx={{ gap: 2 }}>
          <IconButton
            color="inherit"
            edge="start"
            onClick={() => setMobileOpen((prev) => !prev)}
            sx={{ display: { sm: 'none' } }}
            aria-label="Открыть меню"
          >
            <MenuRounded />
          </IconButton>

          <Box flexGrow={1} minWidth={0}>
            <Typography variant="h6" fontWeight={600} noWrap>
              Консоль управления кибербезопасностью
            </Typography>
            <Typography variant="body2" color="text.secondary" noWrap>
              Единый мониторинг инцидентов, активов, рисков и соблюдения требований
            </Typography>
          </Box>

          <IconButton
            size="large"
            onClick={(event) => setAnchorEl(event.currentTarget)}
            aria-haspopup="true"
            aria-controls="user-menu"
            color="inherit"
          >
            <Avatar sx={{ width: 36, height: 36, fontWeight: 600 }}>
              {(user?.firstName?.[0] ?? '').toUpperCase()}
              {(user?.lastName?.[0] ?? '').toUpperCase()}
            </Avatar>
          </IconButton>

          <Menu
            id="user-menu"
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={() => setAnchorEl(null)}
            anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            transformOrigin={{ vertical: 'top', horizontal: 'right' }}
          >
            <MenuItem
              onClick={() => {
                logout()
                setAnchorEl(null)
              }}
            >
              <ListItemIcon>
                <LogoutIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText primary="Выйти" />
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      <Box component="nav" sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}>
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={() => setMobileOpen(false)}
          ModalProps={{ keepMounted: true }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          {drawerContent}
        </Drawer>
        <Drawer
          variant="permanent"
          open
          sx={{
            display: { xs: 'none', sm: 'block' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          {drawerContent}
        </Drawer>
      </Box>

      <Box
        component="main"
        sx={{
          flexGrow: 1,
          minHeight: '100vh',
          backgroundColor: 'background.default',
          width: { sm: `calc(100% - ${drawerWidth}px)` },
        }}
      >
        <Toolbar />
        <Box component="section" sx={{ px: { xs: 2, sm: 4 }, py: 4 }}>
          {children}
        </Box>
      </Box>
    </Box>
  )
}
