package http

import (
	"github.com/gofiber/fiber/v2"
)

// UTF8Middleware - middleware для принудительной установки UTF-8 кодировки
func UTF8Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Устанавливаем заголовки для принудительной UTF-8 кодировки
		c.Set("Content-Type", "application/json; charset=utf-8")
		c.Set("Accept-Charset", "utf-8")
		
		// Продолжаем выполнение следующего middleware
		return c.Next()
	}
}

// UTF8TextMiddleware - middleware для текстовых ответов с UTF-8
func UTF8TextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Устанавливаем заголовки для текстовых файлов с UTF-8
		c.Set("Content-Type", "text/plain; charset=utf-8")
		c.Set("Accept-Charset", "utf-8")
		
		// Продолжаем выполнение следующего middleware
		return c.Next()
	}
}
