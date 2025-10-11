import pytesseract
from pdf2image import convert_from_path
from PIL import Image
import os

def extract_pdf_with_ocr(file_path: str, lang: str = 'rus+eng') -> str:
    """Извлекает текст из PDF используя OCR"""
    try:
        # Конвертируем PDF страницы в изображения
        images = convert_from_path(file_path)
        
        text = []
        for i, image in enumerate(images):
            # Применяем OCR к каждой странице
            page_text = pytesseract.image_to_string(image, lang=lang)
            if page_text.strip():
                text.append(f"[Page {i+1}]\n{page_text}")
        
        return "\n\n".join(text)
    except Exception as e:
        raise Exception(f"Failed to extract PDF text with OCR: {str(e)}")

def extract_image_with_ocr(file_path: str, lang: str = 'rus+eng') -> str:
    """Извлекает текст из изображения используя OCR"""
    try:
        image = Image.open(file_path)
        text = pytesseract.image_to_string(image, lang=lang)
        return text.strip()
    except Exception as e:
        raise Exception(f"Failed to extract image text with OCR: {str(e)}")


