import pdfplumber
import PyPDF2

def extract_text(file_path: str) -> str:
    """Извлекает текст из PDF файла"""
    text = ""
    
    try:
        # Попытка с pdfplumber (лучше сохраняет структуру)
        with pdfplumber.open(file_path) as pdf:
            for page in pdf.pages:
                page_text = page.extract_text()
                if page_text:
                    text += page_text + "\n\n"
    except Exception as e:
        # Fallback на PyPDF2
        try:
            with open(file_path, 'rb') as file:
                pdf_reader = PyPDF2.PdfReader(file)
                for page in pdf_reader.pages:
                    page_text = page.extract_text()
                    if page_text:
                        text += page_text + "\n\n"
        except Exception as e2:
            raise Exception(f"Failed to extract PDF text: {str(e)}, {str(e2)}")
    
    return text.strip()


