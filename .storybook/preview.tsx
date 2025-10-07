import React from 'react'
import type { Preview } from '@storybook/react'
import { ThemeProvider } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import Box from '@mui/material/Box'
import { corporateTheme } from '../apps/frontend/src/shared/theme'

const preview: Preview = {
  decorators: [
    (Story) => (
      <ThemeProvider theme={corporateTheme}>
        <CssBaseline />
        <Box sx={{ backgroundColor: 'background.default', minHeight: '100vh', p: 4 }}>
          <Story />
        </Box>
      </ThemeProvider>
    ),
  ],
  parameters: {
    backgrounds: {
      default: 'surface',
      values: [
        { name: 'surface', value: '#FFFFFF' },
        { name: 'workspace', value: '#F4F7FB' },
      ],
    },
    controls: { expanded: true },
  },
}

export default preview
