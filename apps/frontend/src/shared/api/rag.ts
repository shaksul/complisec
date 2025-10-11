import { api } from './client'

export interface IndexedDocument {
  id: string
  document_id: string
  status: string
  error_message?: string
  retry_count: number
  max_retries: number
  graphrag_doc_id?: string
  chunks_count: number
  entities_count: number
  relationships_count: number
  indexed_at: string | null
  last_retry_at: string | null
  created_at: string
  updated_at: string
}

export interface RAGSource {
  document_id: string
  title: string
  content: string
  score: number
}

export interface RAGQueryResult {
  answer: string
  sources: RAGSource[]
}

export async function indexDocument(documentId: string) {
  const res = await api.post(`/rag/index/${documentId}`)
  return res.data
}

export async function getIndexedDocuments() {
  const res = await api.get('/rag/indexed')
  return res.data.data as IndexedDocument[]
}

export async function queryWithRAG(query: string, useGraph: boolean = true, topK: number = 5) {
  const res = await api.post('/rag/query', { 
    query, 
    use_graph: useGraph, 
    top_k: topK 
  })
  return res.data as RAGQueryResult
}

export async function indexAllDocuments() {
  const res = await api.post('/rag/index-all')
  return res.data
}

