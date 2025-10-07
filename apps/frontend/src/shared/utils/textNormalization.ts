/**
 * Утилиты для нормализации текста и корректного отображения UTF-8 символов
 */

/**
 * Нормализует текст для корректного отображения
 * @param text - исходный текст
 * @returns нормализованный текст
 */
export function normalizeText(text: string | null | undefined): string {
  if (!text) return ''
  
  // Убираем лишние пробелы и переносы строк
  let normalized = text.trim()
  
  // Нормализуем Unicode символы (NFD -> NFC)
  normalized = normalized.normalize('NFC')
  
  // Исправляем возможные проблемы с кодировкой
  normalized = normalized.replace(/\u00A0/g, ' ') // Заменяем неразрывные пробелы
  normalized = normalized.replace(/\u2013/g, '-') // Заменяем длинные тире
  normalized = normalized.replace(/\u2014/g, '--') // Заменяем очень длинные тире
  normalized = normalized.replace(/\u2018/g, "'") // Заменяем левые одинарные кавычки
  normalized = normalized.replace(/\u2019/g, "'") // Заменяем правые одинарные кавычки
  normalized = normalized.replace(/\u201C/g, '"') // Заменяем левые двойные кавычки
  normalized = normalized.replace(/\u201D/g, '"') // Заменяем правые двойные кавычки
  
  return normalized
}

/**
 * Проверяет, содержит ли текст корректные UTF-8 символы
 * @param text - текст для проверки
 * @returns true, если текст содержит корректные UTF-8 символы
 */
export function isValidUTF8(text: string): boolean {
  try {
    // Пытаемся декодировать текст как UTF-8
    const encoder = new TextEncoder()
    const decoder = new TextDecoder('utf-8', { fatal: true })
    const encoded = encoder.encode(text)
    decoder.decode(encoded)
    return true
  } catch {
    return false
  }
}

/**
 * Исправляет возможные проблемы с кодировкой в тексте
 * @param text - текст для исправления
 * @returns исправленный текст
 */
export function fixEncoding(text: string): string {
  if (!text) return ''
  
  try {
    // Если текст уже корректный UTF-8, возвращаем его
    if (isValidUTF8(text)) {
      return normalizeText(text)
    }
    
    // Пытаемся исправить кодировку
    const bytes = new Uint8Array(text.length)
    for (let i = 0; i < text.length; i++) {
      bytes[i] = text.charCodeAt(i)
    }
    
    const decoder = new TextDecoder('utf-8', { fatal: false })
    const fixed = decoder.decode(bytes)
    
    return normalizeText(fixed)
  } catch {
    // Если ничего не получилось, возвращаем исходный текст
    return text
  }
}

/**
 * Форматирует текст для отображения в UI
 * @param text - исходный текст
 * @param maxLength - максимальная длина (опционально)
 * @returns отформатированный текст
 */
export function formatTextForDisplay(text: string | null | undefined, maxLength?: number): string {
  let formatted = normalizeText(text)
  
  if (maxLength && formatted.length > maxLength) {
    formatted = formatted.substring(0, maxLength) + '...'
  }
  
  return formatted
}
