import { api } from "./client"

export interface AIProvider {
  id: string
  tenant_id: string
  name: string
  base_url: string
  api_key?: string
  roles: string[]
  models: string[]
  default_model: string
  prompt_template?: string
  is_active: boolean
}

export interface CreateAIProviderDTO {
  name: string
  base_url: string
  api_key?: string
  roles: string[]
  models?: string[]
  default_model?: string
  prompt_template?: string
}

export interface UpdateAIProviderDTO {
  name: string
  base_url: string
  api_key?: string
  roles: string[]
  models?: string[]
  default_model?: string
  prompt_template?: string
  is_active: boolean
}

export async function getProviders(): Promise<AIProvider[]> {
  const res = await api.get("/ai/providers")
  return res.data.data
}

export async function getProvider(id: string): Promise<AIProvider> {
  const res = await api.get(`/ai/providers/${id}`)
  return res.data.data
}

export async function addProvider(dto: CreateAIProviderDTO) {
  const res = await api.post("/ai/providers", dto)
  return res.data.data
}

export async function updateProvider(id: string, dto: UpdateAIProviderDTO) {
  const res = await api.put(`/ai/providers/${id}`, dto)
  return res.data.data
}

export async function deleteProvider(id: string) {
  const res = await api.delete(`/ai/providers/${id}`)
  return res.data.data
}

export interface QueryAIRequest {
  provider_id: string
  role: string
  input: string
  context?: any
  use_rag?: boolean
}

export interface RAGSource {
  document_id: string
  title: string
  content: string
  score: number
}

export interface QueryAIResponse {
  output?: string
  answer?: string
  sources?: RAGSource[]
}

export async function queryAI(dto: QueryAIRequest): Promise<QueryAIResponse> {
  const res = await api.post("/ai/query", dto)
  return res.data
}
