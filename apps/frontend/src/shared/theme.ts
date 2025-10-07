import { createTheme, type ThemeOptions } from '@mui/material/styles'

export const CORPORATE_COLORS = {
  primary: '#1F4B8E',
  primaryDark: '#122B53',
  primaryLight: '#3F6EB8',
  secondary: '#F47C3C',
  accent: '#2FB2A2',
  neutral900: '#101B2E',
  neutral700: '#3A4A64',
  neutral500: '#64748B',
  neutral300: '#CBD5E1',
  neutral100: '#F4F7FB',
  surface: '#FFFFFF',
  border: '#D5DCE5',
}

const shadows: ThemeOptions['shadows'] = [
  'none',
  '0px 2px 4px rgba(16, 27, 46, 0.08)',
  '0px 4px 10px rgba(16, 27, 46, 0.10)',
  '0px 10px 30px rgba(16, 27, 46, 0.12)',
  ...Array(21).fill('none')
] as ThemeOptions['shadows']

const themeOptions: ThemeOptions = {
  palette: {
    mode: 'light',
    primary: {
      main: CORPORATE_COLORS.primary,
      light: CORPORATE_COLORS.primaryLight,
      dark: CORPORATE_COLORS.primaryDark,
      contrastText: '#FFFFFF',
    },
    secondary: {
      main: CORPORATE_COLORS.secondary,
      contrastText: '#FFFFFF',
    },
    success: {
      main: '#2FB37F',
    },
    warning: {
      main: '#F0A202',
    },
    error: {
      main: '#D64545',
    },
    info: {
      main: CORPORATE_COLORS.accent,
    },
    background: {
      default: CORPORATE_COLORS.neutral100,
      paper: CORPORATE_COLORS.surface,
    },
    text: {
      primary: CORPORATE_COLORS.neutral900,
      secondary: CORPORATE_COLORS.neutral500,
    },
    divider: CORPORATE_COLORS.border,
  },
  typography: {
    fontFamily: "'Inter', 'Roboto', 'Segoe UI', 'Helvetica Neue', Arial, sans-serif",
    h1: {
      fontWeight: 600,
      fontSize: '2.75rem',
      letterSpacing: '-0.02em',
    },
    h2: {
      fontWeight: 600,
      fontSize: '2.125rem',
      letterSpacing: '-0.01em',
    },
    h3: {
      fontWeight: 600,
      fontSize: '1.75rem',
    },
    h4: {
      fontWeight: 600,
      fontSize: '1.5rem',
    },
    h5: {
      fontWeight: 600,
      fontSize: '1.25rem',
    },
    h6: {
      fontWeight: 600,
      fontSize: '1.125rem',
    },
    subtitle1: {
      fontWeight: 500,
    },
    subtitle2: {
      fontWeight: 500,
    },
    body1: {
      fontSize: '1rem',
      lineHeight: 1.6,
    },
    body2: {
      fontSize: '0.875rem',
      lineHeight: 1.5,
    },
    button: {
      fontWeight: 600,
      letterSpacing: '0.02em',
      textTransform: 'none',
    },
  },
  shape: {
    borderRadius: 14,
  },
  spacing: 8,
  shadows,
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        ':root': {
          colorScheme: 'light',
          backgroundColor: CORPORATE_COLORS.neutral100,
        },
        body: {
          backgroundColor: CORPORATE_COLORS.neutral100,
          color: CORPORATE_COLORS.neutral900,
          fontFeatureSettings: '"cv11", "ss01"',
        },
        a: {
          color: CORPORATE_COLORS.primary,
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: CORPORATE_COLORS.surface,
          color: CORPORATE_COLORS.neutral900,
          boxShadow: '0 8px 24px rgba(16, 27, 46, 0.08)',
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          backgroundColor: CORPORATE_COLORS.surface,
          backgroundImage: 'none',
          borderRight: `1px solid ${CORPORATE_COLORS.border}`,
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 16,
          boxShadow: '0 12px 36px rgba(16, 27, 46, 0.08)',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          paddingInline: 18,
          paddingBlock: 10,
        },
        containedPrimary: {
          boxShadow: '0 10px 24px rgba(31, 75, 142, 0.25)',
          '&:hover': {
            boxShadow: '0 12px 28px rgba(31, 75, 142, 0.35)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 20,
          boxShadow: '0 28px 56px rgba(16, 27, 46, 0.10)',
        },
      },
    },
    MuiListItemButton: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          marginInline: 8,
          color: CORPORATE_COLORS.neutral700,
          '& .MuiListItemIcon-root': {
            color: 'inherit',
          },
          '&:hover': {
            backgroundColor: 'rgba(31, 75, 142, 0.08)',
            color: CORPORATE_COLORS.primary,
            '& .MuiListItemIcon-root': {
              color: CORPORATE_COLORS.primary,
            },
          },
          '&.Mui-selected': {
            backgroundColor: 'rgba(31, 75, 142, 0.12)',
            color: CORPORATE_COLORS.primary,
            '& .MuiListItemIcon-root': {
              color: CORPORATE_COLORS.primary,
            },
            '&:hover': {
              backgroundColor: 'rgba(31, 75, 142, 0.18)',
            },
          },
        },
      },
    },
    MuiAvatar: {
      styleOverrides: {
        root: {
          backgroundColor: CORPORATE_COLORS.primary,
          color: '#FFFFFF',
        },
      },
    },
  },
}

export const corporateTheme = createTheme(themeOptions)
export type CorporateTheme = typeof corporateTheme




