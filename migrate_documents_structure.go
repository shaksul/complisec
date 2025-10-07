package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

// DocumentInfo структура для информации о документе
type DocumentInfo struct {
	ID           string
	TenantID     string
	Title        string
	OriginalName string
	FilePath     string
	FileSize     int64
	MimeType     string
	FileHash     string
	Tags         []string
	Links        []DocumentLink
}

// DocumentLink структура для связей документов
type DocumentLink struct {
	Module   string
	EntityID string
}

func main() {
	// Подключение к базе данных
	db, err := sql.Open("postgres", "postgres://complisec:complisec123@localhost:5432/complisec?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Starting document structure migration...")

	// Получаем все документы из старой структуры
	documents, err := getDocumentsFromOldStructure(db)
	if err != nil {
		log.Fatal("Failed to get documents:", err)
	}

	fmt.Printf("Found %d documents to migrate\n", len(documents))

	// Мигрируем каждый документ
	migratedCount := 0
	for _, doc := range documents {
		if err := migrateDocument(doc); err != nil {
			log.Printf("Failed to migrate document %s: %v", doc.ID, err)
			continue
		}
		migratedCount++
		fmt.Printf("Migrated document: %s (%s)\n", doc.Title, doc.ID)
	}

	fmt.Printf("Migration completed. Migrated %d out of %d documents\n", migratedCount, len(documents))
}

// getDocumentsFromOldStructure получает все документы из старой структуры
func getDocumentsFromOldStructure(db *sql.DB) ([]DocumentInfo, error) {
	query := `
		SELECT d.id, d.tenant_id, d.title, d.title as original_name, 
		       d.storage_uri as file_path, d.size_bytes as file_size, 
		       d.mime_type, d.checksum_sha256 as file_hash
		FROM documents d
		WHERE d.deleted_at IS NULL AND d.storage_uri IS NOT NULL
		ORDER BY d.created_at`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []DocumentInfo
	for rows.Next() {
		var doc DocumentInfo
		err := rows.Scan(&doc.ID, &doc.TenantID, &doc.Title, &doc.OriginalName,
			&doc.FilePath, &doc.FileSize, &doc.MimeType, &doc.FileHash)
		if err != nil {
			return nil, err
		}

		// Получаем теги для документа
		tags, err := getDocumentTags(db, doc.ID)
		if err != nil {
			log.Printf("Failed to get tags for document %s: %v", doc.ID, err)
		} else {
			doc.Tags = tags
		}

		// Получаем связи для документа
		links, err := getDocumentLinks(db, doc.ID)
		if err != nil {
			log.Printf("Failed to get links for document %s: %v", doc.ID, err)
		} else {
			doc.Links = links
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

// getDocumentTags получает теги документа
func getDocumentTags(db *sql.DB, documentID string) ([]string, error) {
	query := `SELECT tag FROM document_tags WHERE document_id = $1 ORDER BY tag`
	rows, err := db.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// getDocumentLinks получает связи документа
func getDocumentLinks(db *sql.DB, documentID string) ([]DocumentLink, error) {
	query := `SELECT module, entity_id FROM document_links WHERE document_id = $1`
	rows, err := db.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []DocumentLink
	for rows.Next() {
		var link DocumentLink
		if err := rows.Scan(&link.Module, &link.EntityID); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

// migrateDocument мигрирует один документ в новую структуру
func migrateDocument(doc DocumentInfo) error {
	// Определяем модуль и категорию
	module := detectModule(doc)
	category := detectCategory(doc)

	// Создаем новый путь
	oldPath := doc.FilePath
	fileExt := filepath.Ext(oldPath)
	fileName := fmt.Sprintf("%s%s", doc.ID, fileExt)

	// Новая структура: storage/documents/{tenant_id}/modules/{module}/categories/{category}/{filename}
	newPath := filepath.Join("storage", "documents", doc.TenantID, "modules", module, "categories", category, fileName)

	// Создаем директорию если не существует
	if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Проверяем существование старого файла
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("old file does not exist: %s", oldPath)
	}

	// Копируем файл
	if err := copyFile(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Обновляем путь в базе данных
	if err := updateDocumentPath(doc.ID, newPath); err != nil {
		// Если не удалось обновить БД, удаляем новый файл
		os.Remove(newPath)
		return fmt.Errorf("failed to update database: %w", err)
	}

	// Удаляем старый файл только после успешного обновления БД
	if err := os.Remove(oldPath); err != nil {
		log.Printf("Warning: failed to remove old file %s: %v", oldPath, err)
	}

	return nil
}

// detectModule определяет модуль документа
func detectModule(doc DocumentInfo) string {
	// Если есть связи с модулями, используем их
	if len(doc.Links) > 0 {
		return doc.Links[0].Module
	}

	// Если есть теги, пытаемся определить модуль по тегам
	for _, tag := range doc.Tags {
		switch strings.ToLower(tag) {
		case "#активы", "#assets":
			return "assets"
		case "#риски", "#risks":
			return "risks"
		case "#инциденты", "#incidents":
			return "incidents"
		case "#обучение", "#training":
			return "training"
		case "#соответствие", "#compliance":
			return "compliance"
		}
	}

	// По умолчанию - общие документы
	return "general"
}

// detectCategory определяет категорию документа
func detectCategory(doc DocumentInfo) string {
	// Если есть теги, используем первый как категорию
	for _, tag := range doc.Tags {
		if strings.HasPrefix(tag, "#") {
			category := strings.TrimPrefix(tag, "#")
			if category != "" {
				return category
			}
		}
	}

	// По умолчанию - uncategorized
	return "uncategorized"
}

// copyFile копирует файл из старого пути в новый
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// updateDocumentPath обновляет путь к файлу в базе данных
func updateDocumentPath(documentID, newPath string) error {
	db, err := sql.Open("postgres", "postgres://complisec:complisec123@localhost:5432/complisec?sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE documents SET storage_uri = $1 WHERE id = $2`
	_, err = db.Exec(query, newPath, documentID)
	return err
}

