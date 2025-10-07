package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DocumentRepo - репозиторий для работы с документами
type DocumentRepo struct {
	db DBInterface
}

// NewDocumentRepo создает новый экземпляр DocumentRepo
func NewDocumentRepo(db DBInterface) *DocumentRepo {
	return &DocumentRepo{db: db}
}

// Folder структура папки
type Folder struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	ParentID    *string   `json:"parent_id"`
	OwnerID     string    `json:"owner_id"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
	Metadata    *string   `json:"metadata"`
}

// Document структура документа
type Document struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	Name         string    `json:"name"`
	OriginalName string    `json:"original_name"`
	Description  *string   `json:"description"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	FileHash     string    `json:"file_hash"`
	FolderID     *string   `json:"folder_id"`
	OwnerID      string    `json:"owner_id"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsActive     bool      `json:"is_active"`
	Version      string    `json:"version"`
	Metadata     *string   `json:"metadata"`
}

// DocumentTag структура тега документа
type DocumentTag struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Tag        string    `json:"tag"`
	CreatedAt  time.Time `json:"created_at"`
}

// DocumentLink структура связи документа с другими модулями
type DocumentLink struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Module     string    `json:"module"`
	EntityID   string    `json:"entity_id"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}

// OCRText структура OCR текста
type OCRText struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	Language   string    `json:"language"`
	Confidence *float64  `json:"confidence"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DocumentPermission структура разрешения документа
type DocumentPermission struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	SubjectType string     `json:"subject_type"`
	SubjectID   string     `json:"subject_id"`
	ObjectType  string     `json:"object_type"`
	ObjectID    string     `json:"object_id"`
	Permission  string     `json:"permission"`
	GrantedBy   string     `json:"granted_by"`
	GrantedAt   time.Time  `json:"granted_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsActive    bool       `json:"is_active"`
}

// DocumentVersion структура версии документа
type DocumentVersion struct {
	ID                string    `json:"id"`
	DocumentID        string    `json:"document_id"`
	VersionNumber     int       `json:"version_number"`
	FilePath          string    `json:"file_path"`
	FileSize          int64     `json:"file_size"`
	FileHash          string    `json:"file_hash"`
	CreatedBy         string    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	ChangeDescription *string   `json:"change_description"`
}

// DocumentAuditLog структура аудита документа
type DocumentAuditLog struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	DocumentID *string   `json:"document_id"`
	FolderID   *string   `json:"folder_id"`
	UserID     string    `json:"user_id"`
	Action     string    `json:"action"`
	Details    *string   `json:"details"`
	IPAddress  *string   `json:"ip_address"`
	UserAgent  *string   `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateFolder создает новую папку
func (r *DocumentRepo) CreateFolder(ctx context.Context, folder Folder) error {
	query := `
		INSERT INTO folders (id, tenant_id, name, description, parent_id, owner_id, created_by, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, query, folder.ID, folder.TenantID, folder.Name, folder.Description,
		folder.ParentID, folder.OwnerID, folder.CreatedBy, folder.Metadata)
	return err
}

// GetFolderByID получает папку по ID
func (r *DocumentRepo) GetFolderByID(ctx context.Context, id, tenantID string) (*Folder, error) {
	query := `
		SELECT id, tenant_id, name, description, parent_id, owner_id, created_by, 
		       created_at, updated_at, is_active, metadata
		FROM folders 
		WHERE id = $1 AND tenant_id = $2 AND is_active = true`

	var folder Folder
	err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&folder.ID, &folder.TenantID, &folder.Name, &folder.Description, &folder.ParentID,
		&folder.OwnerID, &folder.CreatedBy, &folder.CreatedAt, &folder.UpdatedAt,
		&folder.IsActive, &folder.Metadata)

	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// ListFolders получает список папок
func (r *DocumentRepo) ListFolders(ctx context.Context, tenantID string, parentID *string) ([]Folder, error) {
	query := `
		SELECT id, tenant_id, name, description, parent_id, owner_id, created_by, 
		       created_at, updated_at, is_active, metadata
		FROM folders 
		WHERE tenant_id = $1 AND is_active = true`

	args := []interface{}{tenantID}
	if parentID != nil {
		query += " AND parent_id = $2"
		args = append(args, *parentID)
	} else {
		query += " AND parent_id IS NULL"
	}

	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	folders := make([]Folder, 0)
	for rows.Next() {
		var folder Folder
		err := rows.Scan(&folder.ID, &folder.TenantID, &folder.Name, &folder.Description,
			&folder.ParentID, &folder.OwnerID, &folder.CreatedBy, &folder.CreatedAt,
			&folder.UpdatedAt, &folder.IsActive, &folder.Metadata)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

// UpdateFolder обновляет папку
func (r *DocumentRepo) UpdateFolder(ctx context.Context, folder Folder) error {
	query := `
		UPDATE folders 
		SET name = $1, description = $2, metadata = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND tenant_id = $5`

	_, err := r.db.ExecContext(ctx, query, folder.Name, folder.Description,
		folder.Metadata, folder.ID, folder.TenantID)
	return err
}

// DeleteFolder удаляет папку
func (r *DocumentRepo) DeleteFolder(ctx context.Context, id, tenantID string) error {
	query := `UPDATE folders SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, tenantID)
	return err
}

// CreateDocument создает новый документ
func (r *DocumentRepo) CreateDocument(ctx context.Context, document Document) error {
	query := `
		INSERT INTO documents (id, tenant_id, name, original_name, description, file_path, 
		                      file_size, mime_type, file_hash, folder_id, owner_id, created_by, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := r.db.ExecContext(ctx, query, document.ID, document.TenantID, document.Name,
		document.OriginalName, document.Description, document.FilePath, document.FileSize,
		document.MimeType, document.FileHash, document.FolderID, document.OwnerID,
		document.CreatedBy, document.Metadata)
	return err
}

// GetDocumentByID получает документ по ID
func (r *DocumentRepo) GetDocumentByID(ctx context.Context, id, tenantID string) (*Document, error) {
	query := `
		SELECT id, tenant_id, title as name, title as original_name, description, storage_uri as file_path, size_bytes as file_size, 
		       mime_type, checksum_sha256 as file_hash, NULL as folder_id, owner_id, created_by, created_at, 
		       updated_at, CASE WHEN deleted_at IS NULL THEN true ELSE false END as is_active, version, NULL as metadata
		FROM documents 
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`

	var document Document
	err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&document.ID, &document.TenantID, &document.Name, &document.OriginalName,
		&document.Description, &document.FilePath, &document.FileSize, &document.MimeType,
		&document.FileHash, &document.FolderID, &document.OwnerID, &document.CreatedBy,
		&document.CreatedAt, &document.UpdatedAt, &document.IsActive, &document.Version,
		&document.Metadata)

	if err != nil {
		return nil, err
	}
	return &document, nil
}

// ListDocuments получает список документов
func (r *DocumentRepo) ListDocuments(ctx context.Context, tenantID string, filters map[string]interface{}) ([]Document, error) {
	query := `
		SELECT id, tenant_id, title as name, title as original_name, description, storage_uri as file_path, size_bytes as file_size, 
		       mime_type, checksum_sha256 as file_hash, NULL as folder_id, owner_id, created_by, created_at, 
		       updated_at, CASE WHEN deleted_at IS NULL THEN true ELSE false END as is_active, version, NULL as metadata
		FROM documents 
		WHERE tenant_id = $1 AND deleted_at IS NULL`

	args := []interface{}{tenantID}
	argIndex := 2

	// Добавляем фильтры
	// folder_id фильтр убран, так как в таблице documents нет колонки folder_id

	if mimeType, ok := filters["mime_type"].(string); ok && mimeType != "" {
		query += fmt.Sprintf(" AND mime_type = $%d", argIndex)
		args = append(args, mimeType)
		argIndex++
	}

	if ownerID, ok := filters["owner_id"].(string); ok && ownerID != "" {
		query += fmt.Sprintf(" AND owner_id = $%d", argIndex)
		args = append(args, ownerID)
		argIndex++
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	// Сортировка
	sortBy := "created_at"
	if sb, ok := filters["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	sortOrder := "DESC"
	if so, ok := filters["sort_order"].(string); ok && so != "" {
		sortOrder = so
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Пагинация
	if page, ok := filters["page"].(int); ok && page > 0 {
		limit := 20
		if l, ok := filters["limit"].(int); ok && l > 0 {
			limit = l
		}
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	fmt.Printf("DEBUG: ListDocuments query: %s\n", query)
	fmt.Printf("DEBUG: ListDocuments args: %v\n", args)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		fmt.Printf("ERROR: ListDocuments query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	documents := make([]Document, 0)
	for rows.Next() {
		var document Document
		err := rows.Scan(&document.ID, &document.TenantID, &document.Name, &document.OriginalName,
			&document.Description, &document.FilePath, &document.FileSize, &document.MimeType,
			&document.FileHash, &document.FolderID, &document.OwnerID, &document.CreatedBy,
			&document.CreatedAt, &document.UpdatedAt, &document.IsActive, &document.Version,
			&document.Metadata)
		if err != nil {
			fmt.Printf("ERROR: ListDocuments scan failed: %v\n", err)
			return nil, err
		}
		documents = append(documents, document)
	}

	return documents, nil
}

// UpdateDocument обновляет документ
func (r *DocumentRepo) UpdateDocument(ctx context.Context, document Document) error {
	query := `
		UPDATE documents 
		SET name = $1, description = $2, folder_id = $3, metadata = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND tenant_id = $6`

	_, err := r.db.ExecContext(ctx, query, document.Name, document.Description,
		document.FolderID, document.Metadata, document.ID, document.TenantID)
	return err
}

// DeleteDocument удаляет документ
func (r *DocumentRepo) DeleteDocument(ctx context.Context, id, tenantID string) error {
	query := `UPDATE documents SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, tenantID)
	return err
}

// AddDocumentTag добавляет тег к документу
func (r *DocumentRepo) AddDocumentTag(ctx context.Context, documentID, tag string) error {
	query := `INSERT INTO document_tags (id, document_id, tag) VALUES ($1, $2, $3) ON CONFLICT (document_id, tag) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, uuid.New().String(), documentID, tag)
	return err
}

// RemoveDocumentTag удаляет тег у документа
func (r *DocumentRepo) RemoveDocumentTag(ctx context.Context, documentID, tag string) error {
	query := `DELETE FROM document_tags WHERE document_id = $1 AND tag = $2`
	_, err := r.db.ExecContext(ctx, query, documentID, tag)
	return err
}

// GetDocumentTags получает теги документа
func (r *DocumentRepo) GetDocumentTags(ctx context.Context, documentID string) ([]string, error) {
	query := `SELECT tag FROM document_tags WHERE document_id = $1 ORDER BY tag`
	rows, err := r.db.QueryContext(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// AddDocumentLink добавляет связь документа с другим модулем
func (r *DocumentRepo) AddDocumentLink(ctx context.Context, link DocumentLink) error {
	query := `
		INSERT INTO document_links (id, document_id, module, entity_id, created_by)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT (document_id, module, entity_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, link.ID, link.DocumentID, link.Module, link.EntityID, link.CreatedBy)
	return err
}

// GetDocumentLinks получает связи документа
func (r *DocumentRepo) GetDocumentLinks(ctx context.Context, documentID string) ([]DocumentLink, error) {
	query := `SELECT id, document_id, module, entity_id, created_by, created_at FROM document_links WHERE document_id = $1`
	rows, err := r.db.QueryContext(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]DocumentLink, 0)
	for rows.Next() {
		var link DocumentLink
		err := rows.Scan(&link.ID, &link.DocumentID, &link.Module, &link.EntityID, &link.CreatedBy, &link.CreatedAt)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

// CreateOCRText создает OCR текст для документа
func (r *DocumentRepo) CreateOCRText(ctx context.Context, ocrText OCRText) error {
	query := `
		INSERT INTO ocr_text (id, document_id, content, language, confidence)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT (document_id) DO UPDATE SET
		content = EXCLUDED.content, language = EXCLUDED.language, 
		confidence = EXCLUDED.confidence, updated_at = CURRENT_TIMESTAMP`
	_, err := r.db.ExecContext(ctx, query, ocrText.ID, ocrText.DocumentID, ocrText.Content, ocrText.Language, ocrText.Confidence)
	return err
}

// GetOCRText получает OCR текст документа
func (r *DocumentRepo) GetOCRText(ctx context.Context, documentID string) (*OCRText, error) {
	query := `SELECT id, document_id, content, language, confidence, created_at, updated_at FROM ocr_text WHERE document_id = $1`
	var ocrText OCRText
	err := r.db.QueryRowContext(ctx, query, documentID).Scan(
		&ocrText.ID, &ocrText.DocumentID, &ocrText.Content, &ocrText.Language,
		&ocrText.Confidence, &ocrText.CreatedAt, &ocrText.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &ocrText, nil
}

// CreateDocumentPermission создает разрешение для документа
func (r *DocumentRepo) CreateDocumentPermission(ctx context.Context, permission DocumentPermission) error {
	query := `
		INSERT INTO document_permissions (id, tenant_id, subject_type, subject_id, object_type, object_id, permission, granted_by, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, permission.ID, permission.TenantID, permission.SubjectType,
		permission.SubjectID, permission.ObjectType, permission.ObjectID, permission.Permission,
		permission.GrantedBy, permission.ExpiresAt)
	return err
}

// GetDocumentPermissions получает разрешения документа
func (r *DocumentRepo) GetDocumentPermissions(ctx context.Context, objectType, objectID, tenantID string) ([]DocumentPermission, error) {
	query := `
		SELECT id, tenant_id, subject_type, subject_id, object_type, object_id, permission, 
		       granted_by, granted_at, expires_at, is_active
		FROM document_permissions 
		WHERE object_type = $1 AND object_id = $2 AND tenant_id = $3 AND is_active = true`

	rows, err := r.db.QueryContext(ctx, query, objectType, objectID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]DocumentPermission, 0)
	for rows.Next() {
		var permission DocumentPermission
		err := rows.Scan(&permission.ID, &permission.TenantID, &permission.SubjectType,
			&permission.SubjectID, &permission.ObjectType, &permission.ObjectID,
			&permission.Permission, &permission.GrantedBy, &permission.GrantedAt,
			&permission.ExpiresAt, &permission.IsActive)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

// CreateDocumentVersion создает версию документа
func (r *DocumentRepo) CreateDocumentVersion(ctx context.Context, version DocumentVersion) error {
	query := `
		INSERT INTO document_versions (id, document_id, version_number, file_path, file_size, file_hash, created_by, change_description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, version.ID, version.DocumentID, version.VersionNumber,
		version.FilePath, version.FileSize, version.FileHash, version.CreatedBy, version.ChangeDescription)
	return err
}

// GetDocumentVersions получает версии документа
func (r *DocumentRepo) GetDocumentVersions(ctx context.Context, documentID string) ([]DocumentVersion, error) {
	query := `
		SELECT id, document_id, version_number, file_path, file_size, file_hash, created_by, created_at, change_description
		FROM document_versions 
		WHERE document_id = $1 
		ORDER BY version_number DESC`

	rows, err := r.db.QueryContext(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make([]DocumentVersion, 0)
	for rows.Next() {
		var version DocumentVersion
		err := rows.Scan(&version.ID, &version.DocumentID, &version.VersionNumber,
			&version.FilePath, &version.FileSize, &version.FileHash, &version.CreatedBy,
			&version.CreatedAt, &version.ChangeDescription)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

// CreateDocumentAuditLog создает запись аудита
func (r *DocumentRepo) CreateDocumentAuditLog(ctx context.Context, log DocumentAuditLog) error {
	query := `
		INSERT INTO document_audit_log (id, tenant_id, document_id, folder_id, user_id, action, details, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.TenantID, log.DocumentID, log.FolderID,
		log.UserID, log.Action, log.Details, log.IPAddress, log.UserAgent)
	return err
}

// GetDocumentAuditLog получает аудит документов
func (r *DocumentRepo) GetDocumentAuditLog(ctx context.Context, tenantID string, filters map[string]interface{}) ([]DocumentAuditLog, error) {
	query := `
		SELECT id, tenant_id, document_id, folder_id, user_id, action, details, ip_address, user_agent, created_at
		FROM document_audit_log 
		WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	argIndex := 2

	if documentID, ok := filters["document_id"].(string); ok && documentID != "" {
		query += fmt.Sprintf(" AND document_id = $%d", argIndex)
		args = append(args, documentID)
		argIndex++
	}

	if userID, ok := filters["user_id"].(string); ok && userID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}

	if action, ok := filters["action"].(string); ok && action != "" {
		query += fmt.Sprintf(" AND action = $%d", argIndex)
		args = append(args, action)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	// Пагинация
	if page, ok := filters["page"].(int); ok && page > 0 {
		limit := 50
		if l, ok := filters["limit"].(int); ok && l > 0 {
			limit = l
		}
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]DocumentAuditLog, 0)
	for rows.Next() {
		var log DocumentAuditLog
		err := rows.Scan(&log.ID, &log.TenantID, &log.DocumentID, &log.FolderID,
			&log.UserID, &log.Action, &log.Details, &log.IPAddress, &log.UserAgent, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// SearchDocuments выполняет поиск документов
func (r *DocumentRepo) SearchDocuments(ctx context.Context, tenantID, searchTerm string) ([]Document, error) {
	query := `
		SELECT DISTINCT d.id, d.tenant_id, d.title as name, d.title as original_name, d.description, d.storage_uri as file_path, 
		       d.size_bytes as file_size, d.mime_type, d.checksum_sha256 as file_hash, NULL as folder_id, d.owner_id, d.created_by, 
		       d.created_at, d.updated_at, CASE WHEN d.deleted_at IS NULL THEN true ELSE false END as is_active, d.version, NULL as metadata
		FROM documents d
		LEFT JOIN document_tags dt ON d.id = dt.document_id
		LEFT JOIN ocr_text ot ON d.id = ot.document_id
		WHERE d.tenant_id = $1 AND d.deleted_at IS NULL
		AND (d.title ILIKE $2 OR d.description ILIKE $2 OR dt.tag ILIKE $2 OR ot.content ILIKE $2)
		ORDER BY d.created_at DESC`

	searchPattern := "%" + searchTerm + "%"
	rows, err := r.db.QueryContext(ctx, query, tenantID, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := make([]Document, 0)
	for rows.Next() {
		var document Document
		err := rows.Scan(&document.ID, &document.TenantID, &document.Name, &document.OriginalName,
			&document.Description, &document.FilePath, &document.FileSize, &document.MimeType,
			&document.FileHash, &document.FolderID, &document.OwnerID, &document.CreatedBy,
			&document.CreatedAt, &document.UpdatedAt, &document.IsActive, &document.Version,
			&document.Metadata)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	return documents, nil
}

// DeleteDocumentLink удаляет связь документа с другим модулем
func (r *DocumentRepo) DeleteDocumentLink(ctx context.Context, documentID, module, entityID string) error {
	query := `DELETE FROM document_links WHERE document_id = $1 AND module = $2 AND entity_id = $3`
	_, err := r.db.ExecContext(ctx, query, documentID, module, entityID)
	return err
}
