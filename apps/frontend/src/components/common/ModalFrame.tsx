import React from 'react'
import { Dialog, DialogTitle, DialogContent, DialogActions, type DialogProps, Typography, Box } from '@mui/material'

interface ModalFrameProps extends DialogProps {
  title: string
  description?: string
  actions?: React.ReactNode
}

export const ModalFrame: React.FC<ModalFrameProps> = ({ title, description, actions, children, ...props }) => (
  <Dialog fullWidth maxWidth="md" {...props}>
    <DialogTitle>
      <Box>
        <Typography variant="h6" component="h2">
          {title}
        </Typography>
        {description && (
          <Typography variant="body2" color="text.secondary" mt={0.5}>
            {description}
          </Typography>
        )}
      </Box>
    </DialogTitle>
    <DialogContent dividers sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
      {children}
    </DialogContent>
    {actions && <DialogActions>{actions}</DialogActions>}
  </Dialog>
)
