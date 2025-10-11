from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import os
from converters import pdf_converter, docx_converter, xlsx_converter, ocr_handler

app = FastAPI(title="Document Processor Service", version="1.0.0")

class ConvertRequest(BaseModel):
    file_path: str
    use_ocr: bool = False

@app.post("/convert")
async def convert_document(req: ConvertRequest):
    """Конвертирует документ в текст"""
    if not os.path.exists(req.file_path):
        raise HTTPException(status_code=404, detail="File not found")
    
    ext = os.path.splitext(req.file_path)[1].lower()
    
    try:
        if ext == '.pdf':
            if req.use_ocr:
                text = ocr_handler.extract_pdf_with_ocr(req.file_path)
            else:
                text = pdf_converter.extract_text(req.file_path)
        elif ext == '.docx':
            text = docx_converter.extract_text(req.file_path)
        elif ext in ['.xlsx', '.xls']:
            text = xlsx_converter.extract_text(req.file_path)
        elif ext == '.txt':
            with open(req.file_path, 'r', encoding='utf-8') as f:
                text = f.read()
        else:
            raise HTTPException(status_code=400, detail=f"Unsupported format: {ext}")
        
        return {"text": text, "char_count": len(text)}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health():
    return {"status": "ok", "service": "doc-processor"}


