import React from 'react'
import {
  Container,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  Box,
  Chip,
} from '@mui/material'
import { Add, Computer } from '@mui/icons-material'

export const AssetsPage: React.FC = () => {
  const assets = [
    { id: '1', name: 'Сервер Web-01', type: 'Server', status: 'active', location: 'Серверная 1' },
    { id: '2', name: 'ПК Администратора', type: 'Workstation', status: 'active', location: 'Офис 101' },
    { id: '3', name: 'Принтер HP LaserJet', type: 'Printer', status: 'maintenance', location: 'Офис 102' },
  ]

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'success'
      case 'maintenance': return 'warning'
      case 'inactive': return 'error'
      default: return 'default'
    }
  }

  return (
    <Container maxWidth="lg">
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Активы</Typography>
        <Button variant="contained" startIcon={<Add />}>
          Добавить актив
        </Button>
      </Box>

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Название</TableCell>
                <TableCell>Тип</TableCell>
                <TableCell>Статус</TableCell>
                <TableCell>Местоположение</TableCell>
                <TableCell>Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {assets.map((asset) => (
                <TableRow key={asset.id}>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      <Computer sx={{ mr: 1 }} />
                      {asset.name}
                    </Box>
                  </TableCell>
                  <TableCell>{asset.type}</TableCell>
                  <TableCell>
                    <Chip
                      label={asset.status}
                      color={getStatusColor(asset.status) as any}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>{asset.location}</TableCell>
                  <TableCell>
                    <Button size="small">Редактировать</Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
    </Container>
  )
}
