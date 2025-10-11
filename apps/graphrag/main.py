from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Dict, List, Any, Optional
import hashlib
import json
import os

app = FastAPI(title="GraphRAG Service", version="1.0.0")

# In-memory storage (в production использовать SQLite/Neo4j)
documents_store: Dict[str, Dict] = {}
vectors_store: Dict[str, List[float]] = {}

class IndexRequest(BaseModel):
    id: str
    title: str
    content: str
    metadata: Optional[Dict[str, Any]] = {}

class QueryRequest(BaseModel):
    query: str
    use_graph: bool = True
    top_k: int = 5
    filter: Optional[Dict[str, str]] = {}

class QueryResponse(BaseModel):
    answer: str
    sources: List[Dict[str, Any]]

@app.post("/index")
async def index_document(req: IndexRequest):
    """Индексирует документ (упрощенная версия)"""
    try:
        # Разбиваем на чанки (упрощенно - по 500 символов)
        content = req.content
        chunk_size = 500
        chunks = []
        
        for i in range(0, len(content), chunk_size):
            chunk_text = content[i:i + chunk_size]
            if chunk_text.strip():
                chunks.append({
                    "text": chunk_text,
                    "index": len(chunks)
                })
        
        # Сохраняем документ
        doc_data = {
            "id": req.id,
            "title": req.title,
            "content": req.content,
            "metadata": req.metadata,
            "chunks": chunks,
            "chunks_count": len(chunks),
            # Упрощенная имитация сущностей
            "entities_count": len(content.split()) // 10,  # ~10% слов как сущности
            "relationships_count": len(chunks) * 2  # ~2 связи на чанк
        }
        
        documents_store[req.id] = doc_data
        
        return {
            "status": "indexed",
            "doc_id": req.id,
            "chunks_count": doc_data["chunks_count"],
            "entities_count": doc_data["entities_count"],
            "relationships_count": doc_data["relationships_count"]
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/query")
async def query(req: QueryRequest):
    """Поиск по индексированным документам"""
    try:
        query_lower = req.query.lower()
        results = []
        
        # Упрощенный поиск - ищем совпадения в чанках
        for doc_id, doc_data in documents_store.items():
            # Применяем фильтры
            if req.filter:
                tenant_id = req.filter.get("tenant_id")
                if tenant_id and doc_data["metadata"].get("tenant_id") != tenant_id:
                    continue
            
            # Ищем релевантные чанки
            for chunk in doc_data["chunks"]:
                chunk_text = chunk["text"].lower()
                
                # Простой поиск по вхождению слов
                query_words = set(query_lower.split())
                chunk_words = set(chunk_text.split())
                common_words = query_words & chunk_words
                
                if common_words:
                    score = len(common_words) / len(query_words) if query_words else 0
                    results.append({
                        "document_id": doc_id,
                        "title": doc_data["title"],
                        "content": chunk["text"][:300] + "..." if len(chunk["text"]) > 300 else chunk["text"],
                        "score": score
                    })
        
        # Сортируем по релевантности
        results.sort(key=lambda x: x["score"], reverse=True)
        results = results[:req.top_k]
        
        # НЕ генерируем ответ - просто возвращаем источники
        # Ответ будет генерировать AI провайдер на основе этих источников
        return {
            "sources": results
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health():
    return {
        "status": "ok",
        "service": "graphrag",
        "documents_indexed": len(documents_store)
    }

@app.get("/stats")
async def stats():
    total_chunks = sum(doc["chunks_count"] for doc in documents_store.values())
    total_entities = sum(doc["entities_count"] for doc in documents_store.values())
    total_relationships = sum(doc["relationships_count"] for doc in documents_store.values())
    
    return {
        "documents": len(documents_store),
        "chunks": total_chunks,
        "entities": total_entities,
        "relationships": total_relationships
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)

