# 🤖 AI Module Pack 

Этот пакет содержит полный набор файлов (с путями, названиями и кодом), чтобы интегрировать поддержку **ИИ-провайдеров** в систему. Все файлы нужно вставить строго в указанные директории.

---

## 📂 Backend (Go + Fiber)

### 1. DTO

**Файл:** `apps/backend/internal/dto/ai.go`

```go
package dto

type CreateAIProviderDTO struct {
  Name           string   `json:"name" validate:"required"`
  BaseURL        string   `json:"base_url" validate:"required,url"`
  APIKey         string   `json:"api_key"`
  Roles          []string `json:"roles" validate:"required,dive,required"`
  PromptTemplate string   `json:"prompt_template"`
}

type QueryAIRequest struct {
  ProviderID string      `json:"provider_id" validate:"required,uuid"`
  Role       string      `json:"role" validate:"required"`
  Input      string      `json:"input" validate:"required"`
  Context    interface{} `json:"context"`
}

type QueryAIResponse struct {
  Output string `json:"output"`
}
```

---

### 2. Repo

**Файл:** `apps/backend/internal/repo/ai_repo.go`

```go
package repo

import (
  "context"
  "errors"

  "github.com/jackc/pgx/v5"
)

type AIProvider struct {
  ID             string
  TenantID       string
  Name           string
  BaseURL        string
  APIKey         *string
  Roles          []string
  PromptTemplate *string
  IsActive       bool
}

type AIRepo struct { db *DB }

func NewAIRepo(db *DB) *AIRepo { return &AIRepo{db: db} }

func (r *AIRepo) List(ctx context.Context, tenantID string) ([]AIProvider, error) {
  rows, err := r.db.Query(ctx, `SELECT id, tenant_id, name, base_url, api_key, roles, prompt_template, is_active FROM ai_providers WHERE tenant_id=$1`, tenantID)
  if err != nil { return nil, err }
  defer rows.Close()
  var items []AIProvider
  for rows.Next() {
    var p AIProvider
    if err := rows.Scan(&p.ID,&p.TenantID,&p.Name,&p.BaseURL,&p.APIKey,&p.Roles,&p.PromptTemplate,&p.IsActive); err != nil { return nil, err }
    items = append(items, p)
  }
  return items, nil
}

func (r *AIRepo) Create(ctx context.Context, p AIProvider) error {
  _, err := r.db.Exec(ctx, `INSERT INTO ai_providers(id,tenant_id,name,base_url,api_key,roles,prompt_template,is_active) VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7)`, p.TenantID,p.Name,p.BaseURL,p.APIKey,p.Roles,p.PromptTemplate,p.IsActive)
  return err
}

func (r *AIRepo) Get(ctx context.Context, id string) (*AIProvider, error) {
  row := r.db.QueryRow(ctx, `SELECT id, tenant_id, name, base_url, api_key, roles, prompt_template, is_active FROM ai_providers WHERE id=$1`, id)
  var p AIProvider
  if err := row.Scan(&p.ID,&p.TenantID,&p.Name,&p.BaseURL,&p.APIKey,&p.Roles,&p.PromptTemplate,&p.IsActive); err != nil {
    if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
    return nil, err
  }
  return &p, nil
}
```

---

### 3. Service

**Файл:** `apps/backend/internal/domain/ai_service.go`

```go
package domain

import (
  "context"
  "encoding/json"
  "net/http"
  "bytes"
  "risknexus/backend/internal/repo"
)

type AIService struct { repo *repo.AIRepo }

func NewAIService(r *repo.AIRepo) *AIService { return &AIService{repo:r} }

func (s *AIService) List(ctx context.Context, tenantID string) ([]repo.AIProvider, error) {
  return s.repo.List(ctx, tenantID)
}

func (s *AIService) Create(ctx context.Context, p repo.AIProvider) error {
  return s.repo.Create(ctx, p)
}

func (s *AIService) Query(ctx context.Context, provider repo.AIProvider, role, input string, contextData any) (string, error) {
  payload := map[string]any{"role":role, "input":input, "context":contextData}
  body,_ := json.Marshal(payload)
  req, _ := http.NewRequestWithContext(ctx,"POST",provider.BaseURL,bytes.NewReader(body))
  req.Header.Set("Content-Type","application/json")
  if provider.APIKey!=nil { req.Header.Set("Authorization","Bearer "+*provider.APIKey) }

  resp, err := http.DefaultClient.Do(req)
  if err!=nil { return "", err }
  defer resp.Body.Close()
  var data map[string]any
  if err:=json.NewDecoder(resp.Body).Decode(&data); err!=nil { return "", err }
  if out,ok:=data["output"].(string); ok { return out,nil }
  return "(no output)",nil
}
```

---

### 4. HTTP Handler

**Файл:** `apps/backend/internal/http/ai_handler.go`

```go
package http

import (
  "context"
  "github.com/gofiber/fiber/v2"
  "risknexus/backend/internal/domain"
  "risknexus/backend/internal/dto"
  "risknexus/backend/internal/repo"
)

type AIHandler struct { service *domain.AIService }

func NewAIHandler(s *domain.AIService) *AIHandler { return &AIHandler{service:s} }

func (h *AIHandler) Register(r fiber.Router) {
  r.Get("/ai/providers", h.listProviders)
  r.Post("/ai/providers", h.createProvider)
  r.Post("/ai/query", h.query)
}

func (h *AIHandler) listProviders(c *fiber.Ctx) error {
  tenantID := "demo-tenant" // TODO: получить из JWT
  items, err := h.service.List(context.Background(), tenantID)
  if err!=nil { return c.Status(500).JSON(fiber.Map{"error":err.Error()}) }
  return c.JSON(fiber.Map{"data":items})
}

func (h *AIHandler) createProvider(c *fiber.Ctx) error {
  var dto dto.CreateAIProviderDTO
  if err:=c.BodyParser(&dto); err!=nil { return c.Status(400).JSON(fiber.Map{"error":"bad input"}) }
  p := repo.AIProvider{TenantID:"demo-tenant",Name:dto.Name,BaseURL:dto.BaseURL,APIKey:&dto.APIKey,Roles:dto.Roles,PromptTemplate:&dto.PromptTemplate,IsActive:true}
  if err:=h.service.Create(context.Background(), p); err!=nil { return c.Status(500).JSON(fiber.Map{"error":err.Error()}) }
  return c.JSON(fiber.Map{"data":"ok"})
}

func (h *AIHandler) query(c *fiber.Ctx) error {
  var req dto.QueryAIRequest
  if err:=c.BodyParser(&req); err!=nil { return c.Status(400).JSON(fiber.Map{"error":"bad input"}) }
  prov := repo.AIProvider{ID:req.ProviderID,BaseURL:"http://localhost:11434/api/chat"} // заглушка
  out, err := h.service.Query(context.Background(), prov, req.Role, req.Input, req.Context)
  if err!=nil { return c.Status(500).JSON(fiber.Map{"error":err.Error()}) }
  return c.JSON(dto.QueryAIResponse{Output: out})
}
```

---

## 📂 Frontend (React + Vite + TS)

### 1. API Client

**Файл:** `apps/frontend/src/shared/api/ai.ts`

```ts
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
```

---

### 2. ProvidersPage

**Файл:** `apps/frontend/src/features/ai/ProvidersPage.tsx`

```tsx
import { useEffect, useState } from "react"
import { getProviders, addProvider } from "@/shared/api/ai"

export default function ProvidersPage(){
  const [items,setItems]=useState<any[]>([])
  const [name,setName]=useState("")
  const [url,setUrl]=useState("")

  useEffect(()=>{ getProviders().then(setItems) },[])

  async function handleAdd(){
    await addProvider({name, base_url:url, roles:["docs"]})
    const data=await getProviders(); setItems(data)
  }

  return (
    <div style={{padding:20}}>
      <h2>AI Providers</h2>
      <ul>{items.map(p=>(<li key={p.id}>{p.name} - {p.base_url}</li>))}</ul>
      <input placeholder="name" value={name} onChange={e=>setName(e.target.value)} />
      <input placeholder="base url" value={url} onChange={e=>setUrl(e.target.value)} />
      <button onClick={handleAdd}>Add</button>
    </div>
  )
}
```

---

### 3. QueryPage

**Файл:** `apps/frontend/src/features/ai/QueryPage.tsx`

```tsx
import { useState } from "react"
import { queryAI } from "@/shared/api/ai"

export default function QueryPage(){
  const [input,setInput]=useState("")
  const [output,setOutput]=useState("")

  async function handleSend(){
    const res = await queryAI({provider_id:"demo", role:"docs", input})
    setOutput(res.output)
  }

  return (
    <div style={{padding:20}}>
      <h2>Query AI</h2>
      <textarea value={input} onChange={e=>setInput(e.target.value)} />
      <button onClick={handleSend}>Send</button>
      <pre>{output}</pre>
    </div>
  )
}
```

---

### 4. Routes

**Файл:** `apps/frontend/src/App.tsx` (добавить)

```tsx
import ProvidersPage from "./features/ai/ProvidersPage"
import QueryPage from "./features/ai/QueryPage"

// внутри <Routes>
<Route path="/ai/providers" element={<ProvidersPage />} />
<Route path="/ai/query" element={<QueryPage />} />
```

---

## ✅ Резюме

* Backend: 4 файла (`dto`, `repo`, `service`, `handler`).
* Frontend: 3 файла (`api.ts`, `ProvidersPage.tsx`, `QueryPage.tsx`) + изменения в `App.tsx`.
* Docs: используем `docs/specs/AI.md`.

Этого достаточно, чтобы:

1. Регистрировать ИИ-провайдеров.
2. Отправлять тестовые запросы.
3. Получать ответы от выбранного агента.

---

⚡ Дальше можно постепенно улучшать: добавлять JWT, привязку к tenant_id, полноценное хранение ключей и редактор промптов в UI.
