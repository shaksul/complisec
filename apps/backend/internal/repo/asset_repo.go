package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Asset struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenant_id"`
	InventoryNumber string     `json:"inventory_number"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	Class           string     `json:"class"`
	OwnerID         *string    `json:"owner_id"`
	Location        *string    `json:"location"`
	Criticality     string     `json:"criticality"`
	Confidentiality string     `json:"confidentiality"`
	Integrity       string     `json:"integrity"`
	Availability    string     `json:"availability"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type AssetDocument struct {
	ID           string    `json:"id"`
	AssetID      string    `json:"asset_id"`
	DocumentType string    `json:"document_type"`
	FilePath     string    `json:"file_path"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type AssetHistory struct {
	ID           string    `json:"id"`
	AssetID      string    `json:"asset_id"`
	FieldChanged string    `json:"field_changed"`
	OldValue     *string   `json:"old_value,omitempty"`
	NewValue     string    `json:"new_value"`
	ChangedBy    string    `json:"changed_by"`
	ChangedAt    time.Time `json:"changed_at"`
}

type AssetSoftware struct {
	ID           string     `json:"id"`
	AssetID      string     `json:"asset_id"`
	SoftwareName string     `json:"software_name"`
	Version      *string    `json:"version,omitempty"`
	InstalledAt  *time.Time `json:"installed_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type AssetWithDetails struct {
	Asset
	OwnerName *string
	Documents []AssetDocument
	Software  []AssetSoftware
	History   []AssetHistory
}

type AssetRepo struct {
	db *DB
}

func NewAssetRepo(db *DB) *AssetRepo {
	return &AssetRepo{db: db}
}

func (r *AssetRepo) Create(ctx context.Context, asset Asset) error {
	// Generate inventory number if not provided
	if asset.InventoryNumber == "" {
		asset.InventoryNumber = r.generateInventoryNumber(ctx, asset.TenantID)
	}

	log.Printf("DEBUG: asset_repo.Create inserting asset tenant=%s name=%s", asset.TenantID, asset.Name)
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO assets (id, tenant_id, inventory_number, name, type, class, owner_id, location, 
                           criticality, confidentiality, integrity, availability, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `, asset.ID, asset.TenantID, asset.InventoryNumber, asset.Name, asset.Type, asset.Class,
		asset.OwnerID, asset.Location, asset.Criticality, asset.Confidentiality,
		asset.Integrity, asset.Availability, asset.Status)
	if err != nil {
		log.Printf("ERROR: asset_repo.Create insert failed: %v", err)
	}
	return err
}

func (r *AssetRepo) GetByID(ctx context.Context, id string) (*Asset, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, inventory_number, name, type, class, owner_id, location,
		       criticality, confidentiality, integrity, availability, status, created_at, updated_at, deleted_at
		FROM assets WHERE id = $1 AND deleted_at IS NULL
	`, id)

	var asset Asset
	err := row.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
		&asset.Type, &asset.Class, &asset.OwnerID, &asset.Location,
		&asset.Criticality, &asset.Confidentiality, &asset.Integrity,
		&asset.Availability, &asset.Status, &asset.CreatedAt, &asset.UpdatedAt, &asset.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepo) GetByInventoryNumber(ctx context.Context, tenantID, inventoryNumber string) (*Asset, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, inventory_number, name, type, class, owner_id, location,
		       criticality, confidentiality, integrity, availability, status, created_at, updated_at, deleted_at
		FROM assets WHERE tenant_id = $1 AND inventory_number = $2 AND deleted_at IS NULL
	`, tenantID, inventoryNumber)

	var asset Asset
	err := row.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
		&asset.Type, &asset.Class, &asset.OwnerID, &asset.Location,
		&asset.Criticality, &asset.Confidentiality, &asset.Integrity,
		&asset.Availability, &asset.Status, &asset.CreatedAt, &asset.UpdatedAt, &asset.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepo) List(ctx context.Context, tenantID string, filters map[string]interface{}) ([]Asset, error) {
	query := `
		SELECT id, tenant_id, inventory_number, name, type, class, owner_id, location,
		       criticality, confidentiality, integrity, availability, status, created_at, updated_at, deleted_at
		FROM assets WHERE tenant_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{tenantID}
	argIndex := 2

	// Apply filters
	if assetType, ok := filters["type"].(string); ok && assetType != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, assetType)
		argIndex++
	}
	if class, ok := filters["class"].(string); ok && class != "" {
		query += fmt.Sprintf(" AND class = $%d", argIndex)
		args = append(args, class)
		argIndex++
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if criticality, ok := filters["criticality"].(string); ok && criticality != "" {
		query += fmt.Sprintf(" AND criticality = $%d", argIndex)
		args = append(args, criticality)
		argIndex++
	}
	if ownerID, ok := filters["owner_id"].(string); ok && ownerID != "" {
		query += fmt.Sprintf(" AND owner_id = $%d", argIndex)
		args = append(args, ownerID)
		argIndex++
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR inventory_number ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		err := rows.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
			&asset.Type, &asset.Class, &asset.OwnerID, &asset.Location,
			&asset.Criticality, &asset.Confidentiality, &asset.Integrity,
			&asset.Availability, &asset.Status, &asset.CreatedAt, &asset.UpdatedAt, &asset.DeletedAt)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

func (r *AssetRepo) ListPaginated(ctx context.Context, tenantID string, page, pageSize int, filters map[string]interface{}) ([]Asset, int64, error) {
	offset := (page - 1) * pageSize

	// Build base query for counting
	countQuery := "SELECT COUNT(*) FROM assets WHERE tenant_id = $1 AND deleted_at IS NULL"
	args := []interface{}{tenantID}
	argIndex := 2

	// Apply filters for count
	if assetType, ok := filters["type"].(string); ok && assetType != "" {
		countQuery += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, assetType)
		argIndex++
	}
	if class, ok := filters["class"].(string); ok && class != "" {
		countQuery += fmt.Sprintf(" AND class = $%d", argIndex)
		args = append(args, class)
		argIndex++
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if criticality, ok := filters["criticality"].(string); ok && criticality != "" {
		countQuery += fmt.Sprintf(" AND criticality = $%d", argIndex)
		args = append(args, criticality)
		argIndex++
	}
	if ownerID, ok := filters["owner_id"].(string); ok && ownerID != "" {
		countQuery += fmt.Sprintf(" AND owner_id = $%d", argIndex)
		args = append(args, ownerID)
		argIndex++
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR inventory_number ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	// Get total count
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build query for data
	query := `
		SELECT id, tenant_id, inventory_number, name, type, class, owner_id, location,
		       criticality, confidentiality, integrity, availability, status, created_at, updated_at, deleted_at
		FROM assets WHERE tenant_id = $1 AND deleted_at IS NULL
	`
	dataArgs := []interface{}{tenantID}
	dataArgIndex := 2

	// Apply same filters for data query
	if assetType, ok := filters["type"].(string); ok && assetType != "" {
		query += fmt.Sprintf(" AND type = $%d", dataArgIndex)
		dataArgs = append(dataArgs, assetType)
		dataArgIndex++
	}
	if class, ok := filters["class"].(string); ok && class != "" {
		query += fmt.Sprintf(" AND class = $%d", dataArgIndex)
		dataArgs = append(dataArgs, class)
		dataArgIndex++
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query += fmt.Sprintf(" AND status = $%d", dataArgIndex)
		dataArgs = append(dataArgs, status)
		dataArgIndex++
	}
	if criticality, ok := filters["criticality"].(string); ok && criticality != "" {
		query += fmt.Sprintf(" AND criticality = $%d", dataArgIndex)
		dataArgs = append(dataArgs, criticality)
		dataArgIndex++
	}
	if ownerID, ok := filters["owner_id"].(string); ok && ownerID != "" {
		query += fmt.Sprintf(" AND owner_id = $%d", dataArgIndex)
		dataArgs = append(dataArgs, ownerID)
		dataArgIndex++
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR inventory_number ILIKE $%d)", dataArgIndex, dataArgIndex)
		dataArgs = append(dataArgs, "%"+search+"%")
		dataArgIndex++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", dataArgIndex, dataArgIndex+1)
	dataArgs = append(dataArgs, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		err := rows.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
			&asset.Type, &asset.Class, &asset.OwnerID, &asset.Location,
			&asset.Criticality, &asset.Confidentiality, &asset.Integrity,
			&asset.Availability, &asset.Status, &asset.CreatedAt, &asset.UpdatedAt, &asset.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		assets = append(assets, asset)
	}

	return assets, total, nil
}

func (r *AssetRepo) Update(ctx context.Context, asset Asset) error {
	_, err := r.db.Exec(`
		UPDATE assets SET name = $1, type = $2, class = $3, owner_id = $4, location = $5,
		                  criticality = $6, confidentiality = $7, integrity = $8, availability = $9,
		                  status = $10, updated_at = CURRENT_TIMESTAMP
		WHERE id = $11
	`, asset.Name, asset.Type, asset.Class, asset.OwnerID, asset.Location,
		asset.Criticality, asset.Confidentiality, asset.Integrity,
		asset.Availability, asset.Status, asset.ID)
	return err
}

func (r *AssetRepo) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(`
		UPDATE assets SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id)
	return err
}

func (r *AssetRepo) GetWithDetails(ctx context.Context, id string) (*AssetWithDetails, error) {
	asset, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, nil
	}

	// Get owner name
	var ownerName *string
	if asset.OwnerID != nil {
		var firstName, lastName sql.NullString
		err := r.db.QueryRow(`
			SELECT first_name, last_name FROM users WHERE id = $1
		`, *asset.OwnerID).Scan(&firstName, &lastName)
		if err == nil {
			if firstName.Valid && lastName.Valid {
				fullName := firstName.String + " " + lastName.String
				ownerName = &fullName
			} else if firstName.Valid {
				ownerName = &firstName.String
			}
		}
	}

	// Get documents
	documents, err := r.GetAssetDocuments(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get software
	software, err := r.GetAssetSoftware(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get history
	history, err := r.GetAssetHistory(ctx, id)
	if err != nil {
		return nil, err
	}

	return &AssetWithDetails{
		Asset:     *asset,
		OwnerName: ownerName,
		Documents: documents,
		Software:  software,
		History:   history,
	}, nil
}

func (r *AssetRepo) AddDocument(ctx context.Context, assetID, documentType, filePath, createdBy string) error {
	_, err := r.db.Exec(`
		INSERT INTO asset_documents (id, asset_id, document_type, file_path, created_by)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.New().String(), assetID, documentType, filePath, createdBy)
	return err
}

func (r *AssetRepo) GetAssetDocuments(ctx context.Context, assetID string) ([]AssetDocument, error) {
	rows, err := r.db.Query(`
		SELECT id, asset_id, document_type, file_path, created_by, created_at
		FROM asset_documents WHERE asset_id = $1 ORDER BY created_at DESC
	`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []AssetDocument
	for rows.Next() {
		var doc AssetDocument
		err := rows.Scan(&doc.ID, &doc.AssetID, &doc.DocumentType, &doc.FilePath, &doc.CreatedBy, &doc.CreatedAt)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func (r *AssetRepo) AddSoftware(ctx context.Context, assetID, softwareName, version string, installedAt *time.Time) error {
	_, err := r.db.Exec(`
		INSERT INTO asset_software (id, asset_id, software_name, version, installed_at)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.New().String(), assetID, softwareName, version, installedAt)
	return err
}

func (r *AssetRepo) GetAssetSoftware(ctx context.Context, assetID string) ([]AssetSoftware, error) {
	rows, err := r.db.Query(`
		SELECT id, asset_id, software_name, version, installed_at, updated_at
		FROM asset_software WHERE asset_id = $1 ORDER BY updated_at DESC
	`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var software []AssetSoftware
	for rows.Next() {
		var sw AssetSoftware
		err := rows.Scan(&sw.ID, &sw.AssetID, &sw.SoftwareName, &sw.Version, &sw.InstalledAt, &sw.UpdatedAt)
		if err != nil {
			return nil, err
		}
		software = append(software, sw)
	}
	return software, nil
}

func (r *AssetRepo) AddHistory(ctx context.Context, assetID, fieldChanged, oldValue, newValue, changedBy string) error {
	_, err := r.db.Exec(`
		INSERT INTO asset_history (id, asset_id, field_changed, old_value, new_value, changed_by)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, uuid.New().String(), assetID, fieldChanged, oldValue, newValue, changedBy)
	return err
}

func (r *AssetRepo) GetAssetHistory(ctx context.Context, assetID string) ([]AssetHistory, error) {
	rows, err := r.db.Query(`
		SELECT id, asset_id, field_changed, old_value, new_value, changed_by, changed_at
		FROM asset_history WHERE asset_id = $1 ORDER BY changed_at DESC
	`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []AssetHistory
	for rows.Next() {
		var h AssetHistory
		err := rows.Scan(&h.ID, &h.AssetID, &h.FieldChanged, &h.OldValue, &h.NewValue, &h.ChangedBy, &h.ChangedAt)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func (r *AssetRepo) generateInventoryNumber(ctx context.Context, tenantID string) string {
	// Simple implementation - in production, this should be more sophisticated
	prefix := "AST"
	timestamp := time.Now().Format("20060102")
	random := strings.ToUpper(uuid.New().String()[:8])
	return fmt.Sprintf("%s-%s-%s", prefix, timestamp, random)
}
