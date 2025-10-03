import React, { useState, useEffect } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  Tabs,
  Tab,
  Box,
  Typography,
  IconButton,
  CircularProgress,
} from '@mui/material'
import { Close } from '@mui/icons-material'
import { Risk } from '../../shared/api/risks'
import { RiskGeneralTab } from './tabs/RiskGeneralTab'
import { RiskControlsTab } from './tabs/RiskControlsTab'
import { RiskCommentsTab } from './tabs/RiskCommentsTab'
import { RiskHistoryTab } from './tabs/RiskHistoryTab'
import { RiskAttachmentsTab } from './tabs/RiskAttachmentsTab'

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`risk-tabpanel-${index}`}
      aria-labelledby={`risk-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          {children}
        </Box>
      )}
    </div>
  )
}

interface RiskDetailsModalProps {
  open: boolean
  onClose: () => void
  risk: Risk | null
}

export const RiskDetailsModal: React.FC<RiskDetailsModalProps> = ({
  open,
  onClose,
  risk,
}) => {
  const [tabValue, setTabValue] = useState(0)
  const [loading, setLoading] = useState(false)

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue)
  }

  const handleClose = () => {
    setTabValue(0)
    onClose()
  }

  if (!risk) {
    return null
  }

  return (
    <Dialog 
      open={open} 
      onClose={handleClose} 
      maxWidth="lg" 
      fullWidth
      fullScreen
    >
      <DialogTitle>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h5" component="div">
            {risk.title}
          </Typography>
          <IconButton onClick={handleClose}>
            <Close />
          </IconButton>
        </Box>
      </DialogTitle>
      
      <DialogContent sx={{ p: 0 }}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs 
            value={tabValue} 
            onChange={handleTabChange} 
            aria-label="risk details tabs"
            variant="scrollable"
            scrollButtons="auto"
          >
            <Tab label="Общие данные" />
            <Tab label="Контроли" />
            <Tab label="Комментарии" />
            <Tab label="История изменений" />
            <Tab label="Вложения" />
          </Tabs>
        </Box>

        <TabPanel value={tabValue} index={0}>
          <RiskGeneralTab risk={risk} onUpdate={onClose} />
        </TabPanel>

        <TabPanel value={tabValue} index={1}>
          <RiskControlsTab riskId={risk.id} />
        </TabPanel>

        <TabPanel value={tabValue} index={2}>
          <RiskCommentsTab riskId={risk.id} />
        </TabPanel>

        <TabPanel value={tabValue} index={3}>
          <RiskHistoryTab riskId={risk.id} />
        </TabPanel>

        <TabPanel value={tabValue} index={4}>
          <RiskAttachmentsTab riskId={risk.id} />
        </TabPanel>
      </DialogContent>
    </Dialog>
  )
}

