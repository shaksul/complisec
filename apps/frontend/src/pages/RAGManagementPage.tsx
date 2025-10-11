import { useEffect, useState } from 'react'
import {
  Box,
  Button,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Paper,
  CircularProgress,
} from '@mui/material'
import {
  Refresh,
  CheckCircle,
  Error as ErrorIcon,
  HourglassEmpty,
} from '@mui/icons-material'
import { PageContainer, PageHeader, SectionCard } from '@/components/common/Page'
import { getIndexedDocuments, indexAllDocuments, type IndexedDocument } from '@/shared/api/rag'

export default function RAGManagementPage() {
  const [documents, setDocuments] = useState<IndexedDocument[]>([])
  const [loading, setLoading] = useState(true)
  const [indexing, setIndexing] = useState(false)

  useEffect(() => {
    loadDocuments()
  }, [])

  async function loadDocuments() {
    setLoading(true)
    try {
      const data = await getIndexedDocuments()
      setDocuments(data || [])
    } catch (error) {
      console.error('Failed to load indexed documents:', error)
      setDocuments([])
    } finally {
      setLoading(false)
    }
  }

  async function handleIndexAll() {
    setIndexing(true)
    try {
      await indexAllDocuments()
      alert('Массовая индексация запущена в фоне. Обновите страницу через минуту.')
      loadDocuments()
    } catch (error) {
      console.error('Failed to start indexing:', error)
      alert('Ошибка запуска индексации')
    } finally {
      setIndexing(false)
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'indexed':
        return <CheckCircle color="success" fontSize="small" />
      case 'failed':
        return <ErrorIcon color="error" fontSize="small" />
      case 'processing':
      case 'retrying':
        return <HourglassEmpty color="warning" fontSize="small" />
      default:
        return null
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'indexed':
        return 'success'
      case 'failed':
        return 'error'
      case 'processing':
      case 'retrying':
        return 'warning'
      case 'pending':
        return 'info'
      default:
        return 'default'
    }
  }

  const statusLabels: Record<string, string> = {
    pending: 'Ожидает',
    processing: 'Обработка',
    indexed: 'Индексирован',
    failed: 'Ошибка',
    retrying: 'Повтор',
  }

  return (
    <PageContainer>
      <PageHeader
        title="Управление RAG индексацией"
        subtitle="Статус индексированных документов для GraphRAG"
        actions={
          <Box display="flex" gap={2}>
            <Button
              variant="contained"
              color="primary"
              onClick={handleIndexAll}
              disabled={indexing}
            >
              {indexing ? 'Индексация...' : 'Переиндексировать все'}
            </Button>
            <Button
              variant="outlined"
              startIcon={<Refresh />}
              onClick={loadDocuments}
              disabled={loading}
            >
              Обновить
            </Button>
          </Box>
        }
      />

      <SectionCard
        title="Индексированные документы"
        description="Документы, обработанные GraphRAG для семантического поиска"
      >
        {loading ? (
          <Box display="flex" justifyContent="center" py={4}>
            <CircularProgress />
          </Box>
        ) : !documents || documents.length === 0 ? (
          <Box textAlign="center" py={4}>
            <Typography variant="body1" color="text.secondary">
              Нет индексированных документов
            </Typography>
            <Typography variant="body2" color="text.secondary" mt={1}>
              Документы будут отображаться здесь после индексации
            </Typography>
          </Box>
        ) : (
          <TableContainer component={Paper} variant="outlined">
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Статус</TableCell>
                  <TableCell>ID документа</TableCell>
                  <TableCell align="center">Чанки</TableCell>
                  <TableCell align="center">Сущности</TableCell>
                  <TableCell align="center">Связи</TableCell>
                  <TableCell align="center">Попытки</TableCell>
                  <TableCell>Индексирован</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {documents && documents.map((doc) => (
                  <TableRow key={doc.id} hover>
                    <TableCell>
                      <Box display="flex" alignItems="center" gap={1}>
                        {getStatusIcon(doc.status)}
                        <Chip
                          label={statusLabels[doc.status] || doc.status}
                          color={getStatusColor(doc.status) as any}
                          size="small"
                        />
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" fontFamily="monospace">
                        {doc.document_id.substring(0, 8)}...
                      </Typography>
                    </TableCell>
                    <TableCell align="center">
                      <Typography variant="body2">{doc.chunks_count}</Typography>
                    </TableCell>
                    <TableCell align="center">
                      <Typography variant="body2">{doc.entities_count}</Typography>
                    </TableCell>
                    <TableCell align="center">
                      <Typography variant="body2">{doc.relationships_count}</Typography>
                    </TableCell>
                    <TableCell align="center">
                      <Typography variant="body2">
                        {doc.retry_count} / {doc.max_retries}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" color="text.secondary">
                        {doc.indexed_at
                          ? new Date(doc.indexed_at).toLocaleString('ru-RU')
                          : '—'}
                      </Typography>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}

        {documents && documents.length > 0 && (
          <Box mt={2}>
            <Typography variant="body2" color="text.secondary">
              Всего индексированных документов: {documents.length}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Успешно: {documents.filter((d) => d.status === 'indexed').length} | В
              обработке: {documents.filter((d) => d.status === 'processing').length} |
              Ошибок: {documents.filter((d) => d.status === 'failed').length}
            </Typography>
          </Box>
        )}
      </SectionCard>
    </PageContainer>
  )
}

