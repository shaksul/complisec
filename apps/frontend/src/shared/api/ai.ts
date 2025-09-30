import { api } from "./client"

export async function getProviders() {
  const res = await api.get("/ai/providers")
  return res.data.data
}

export async function addProvider(dto: any) {
  const res = await api.post("/ai/providers", dto)
  return res.data.data
}

export async function queryAI(dto: any) {
  const res = await api.post("/ai/query", dto)
  return res.data
}
