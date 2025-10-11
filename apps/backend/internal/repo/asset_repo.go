package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"risknexus/backend/internal/dto"

	"github.com/google/uuid"
)

type Asset struct {
	ID                  string     `json:"id"`
	TenantID            string     `json:"tenant_id"`
	InventoryNumber     string     `json:"inventory_number"`
	Name                string     `json:"name"`
	Type                string     `json:"type"`
	Class               string     `json:"class"`
	OwnerID             *string    `json:"owner_id"`
	OwnerName           *string    `json:"owner_name,omitempty"`
	ResponsibleUserID   *string    `json:"responsible_user_id"`
	ResponsibleUserName *string    `json:"responsible_user_name,omitempty"`
	Location            *string    `json:"location"`
	Criticality         string     `json:"criticality"`
	Confidentiality     string     `json:"confidentiality"`
	Integrity           string     `json:"integrity"`
	Availability        string     `json:"availability"`
	Status              string     `json:"status"`
	SerialNumber        *string    `json:"serial_number,omitempty"`
	PCNumber            *string    `json:"pc_number,omitempty"`
	Model               *string    `json:"model,omitempty"`
	CPU                 *string    `json:"cpu,omitempty"`
	RAM                 *string    `json:"ram,omitempty"`
	HDDInfo             *string    `json:"hdd_info,omitempty"`
	NetworkCard         *string    `json:"network_card,omitempty"`
	OpticalDrive        *string    `json:"optical_drive,omitempty"`
	IPAddress           *string    `json:"ip_address,omitempty"`
	MACAddress          *string    `json:"mac_address,omitempty"`
	Manufacturer        *string    `json:"manufacturer,omitempty"`
	PurchaseYear        *int       `json:"purchase_year,omitempty"`
	WarrantyUntil       *time.Time `json:"warranty_until,omitempty"`
	Metadata            *string    `json:"metadata,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

type AssetDocument struct {
	ID             string    `json:"id"`
	AssetID        string    `json:"asset_id"`
	DocumentType   string    `json:"document_type"`
	FilePath       string    `json:"file_path"`
	Title          string    `json:"title"`
	Mime           string    `json:"mime"`
	SizeBytes      int64     `json:"size_bytes"`
	CreatedBy      string    `json:"created_by"`
	CreatedByName  *string   `json:"created_by_name,omitempty"`  // ФИО пользователя
	CreatedByEmail *string   `json:"created_by_email,omitempty"` // Email пользователя
	CreatedAt      time.Time `json:"created_at"`
}

type AssetHistory struct {
	ID             string    `json:"id"`
	AssetID        string    `json:"asset_id"`
	FieldChanged   string    `json:"field_changed"`
	OldValue       *string   `json:"old_value,omitempty"`
	NewValue       string    `json:"new_value"`
	ChangedBy      string    `json:"changed_by"`
	ChangedByName  *string   `json:"changed_by_name,omitempty"`  // ФИО пользователя
	ChangedByEmail *string   `json:"changed_by_email,omitempty"` // Email пользователя
	ChangedAt      time.Time `json:"changed_at"`
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
	OwnerName *string         `json:"owner_name,omitempty"`
	Documents []AssetDocument `json:"documents"`
	Software  []AssetSoftware `json:"software"`
	History   []AssetHistory  `json:"history"`
}

// AssetRisk represents a risk associated with an asset
type AssetRisk struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Category    *string   `json:"category"`
	Likelihood  *int      `json:"likelihood"`
	Impact      *int      `json:"impact"`
	Level       *int      `json:"level"`
	Status      string    `json:"status"`
	OwnerID     *string   `json:"owner_id"`
	AssetID     *string   `json:"asset_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AssetIncident represents an incident associated with an asset
type AssetIncident struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Category    string     `json:"category"`
	Status      string     `json:"status"`
	Severity    string     `json:"severity"`
	Source      string     `json:"source"`
	ReportedBy  string     `json:"reported_by"`
	AssignedTo  *string    `json:"assigned_to"`
	AssetID     *string    `json:"asset_id"`
	RiskID      *string    `json:"risk_id"`
	CreatedBy   string     `json:"created_by"`
	DetectedAt  time.Time  `json:"detected_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	ClosedAt    *time.Time `json:"closed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type AssetRepo struct {
	db DBInterface
}

func NewAssetRepo(db DBInterface) *AssetRepo {
	return &AssetRepo{db: db}
}

func (r *AssetRepo) Create(ctx context.Context, asset Asset) error {
	// Generate inventory number if not provided
	if asset.InventoryNumber == "" {
		asset.InventoryNumber = r.generateInventoryNumber(ctx, asset.TenantID)
	}

	log.Printf("DEBUG: asset_repo.Create inserting asset tenant=%s name=%s", asset.TenantID, asset.Name)
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO assets (id, tenant_id, inventory_number, name, type, class, owner_id, responsible_user_id, location, 
                           criticality, confidentiality, integrity, availability, status,
                           serial_number, pc_number, model, cpu, ram, hdd_info, network_card, optical_drive,
                           ip_address, mac_address, manufacturer, purchase_year, warranty_until)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)
    `, asset.ID, asset.TenantID, asset.InventoryNumber, asset.Name, asset.Type, asset.Class,
		asset.OwnerID, asset.ResponsibleUserID, asset.Location, asset.Criticality, asset.Confidentiality,
		asset.Integrity, asset.Availability, asset.Status,
		asset.SerialNumber, asset.PCNumber, asset.Model, asset.CPU, asset.RAM, asset.HDDInfo,
		asset.NetworkCard, asset.OpticalDrive, asset.IPAddress, asset.MACAddress, asset.Manufacturer,
		asset.PurchaseYear, asset.WarrantyUntil)
	if err != nil {
		log.Printf("ERROR: asset_repo.Create insert failed: %v", err)
	}
	return err
}

func (r *AssetRepo) GetByID(ctx context.Context, id string) (*Asset, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, inventory_number, name, type, class, owner_id, responsible_user_id, location,
		       criticality, confidentiality, integrity, availability, status, 
		       serial_number, pc_number, model, cpu, ram, hdd_info, network_card, optical_drive,
		       ip_address, mac_address, manufacturer, purchase_year, warranty_until, metadata,
		       created_at, updated_at, deleted_at
		FROM assets WHERE id = $1 AND deleted_at IS NULL
	`, id)

	var asset Asset
	var serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer sql.NullString
	var purchaseYear sql.NullInt32
	var warrantyUntil sql.NullTime
	var metadata sql.NullString

	err := row.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
		&asset.Type, &asset.Class, &asset.OwnerID, &asset.ResponsibleUserID, &asset.Location,
		&asset.Criticality, &asset.Confidentiality, &asset.Integrity,
		&asset.Availability, &asset.Status,
		&serialNumber, &pcNumber, &model, &cpu, &ram, &hddInfo, &networkCard, &opticalDrive,
		&ipAddress, &macAddress, &manufacturer, &purchaseYear, &warrantyUntil, &metadata,
		&asset.CreatedAt, &asset.UpdatedAt, &asset.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Populate passport fields
	populatePassportFieldsFromRow(&asset, serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer, purchaseYear, warrantyUntil)

	if metadata.Valid {
		asset.Metadata = &metadata.String
	}

	return &asset, nil
}

func (r *AssetRepo) GetByInventoryNumber(ctx context.Context, tenantID, inventoryNumber string) (*Asset, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, inventory_number, name, type, class, owner_id, responsible_user_id, location,
		       criticality, confidentiality, integrity, availability, status,
		       serial_number, pc_number, model, cpu, ram, hdd_info, network_card, optical_drive,
		       ip_address, mac_address, manufacturer, purchase_year, warranty_until, metadata,
		       created_at, updated_at, deleted_at
		FROM assets WHERE tenant_id = $1 AND inventory_number = $2 AND deleted_at IS NULL
	`, tenantID, inventoryNumber)

	var asset Asset
	var serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer sql.NullString
	var purchaseYear sql.NullInt32
	var warrantyUntil sql.NullTime
	var metadata sql.NullString

	err := row.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
		&asset.Type, &asset.Class, &asset.OwnerID, &asset.ResponsibleUserID, &asset.Location,
		&asset.Criticality, &asset.Confidentiality, &asset.Integrity,
		&asset.Availability, &asset.Status,
		&serialNumber, &pcNumber, &model, &cpu, &ram, &hddInfo, &networkCard, &opticalDrive,
		&ipAddress, &macAddress, &manufacturer, &purchaseYear, &warrantyUntil, &metadata,
		&asset.CreatedAt, &asset.UpdatedAt, &asset.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Populate passport fields
	populatePassportFieldsFromRow(&asset, serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer, purchaseYear, warrantyUntil)

	if metadata.Valid {
		asset.Metadata = &metadata.String
	}

	return &asset, nil
}

func (r *AssetRepo) List(ctx context.Context, tenantID string, filters map[string]interface{}) ([]Asset, error) {
	query := `
		SELECT a.id, a.tenant_id, a.inventory_number, a.name, a.type, a.class, a.owner_id, 
		       COALESCE(u_owner.first_name || ' ' || u_owner.last_name, u_owner.email) as owner_name,
		       a.responsible_user_id,
		       COALESCE(u_resp.first_name || ' ' || u_resp.last_name, u_resp.email) as responsible_user_name,
		       a.location, a.criticality, a.confidentiality, a.integrity, a.availability, 
		       a.status, a.serial_number, a.pc_number, a.model, a.cpu, a.ram, a.hdd_info, a.network_card, a.optical_drive,
		       a.ip_address, a.mac_address, a.manufacturer, a.purchase_year, a.warranty_until,
		       a.created_at, a.updated_at, a.deleted_at
		FROM assets a
		LEFT JOIN users u_owner ON a.owner_id = u_owner.id
		LEFT JOIN users u_resp ON a.responsible_user_id = u_resp.id
		WHERE a.tenant_id = $1 AND a.deleted_at IS NULL
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
		var ownerName, responsibleUserName sql.NullString
		var serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer sql.NullString
		var purchaseYear sql.NullInt32
		var warrantyUntil sql.NullTime
		err := rows.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
			&asset.Type, &asset.Class, &asset.OwnerID, &ownerName, &asset.ResponsibleUserID,
			&responsibleUserName, &asset.Location, &asset.Criticality, &asset.Confidentiality,
			&asset.Integrity, &asset.Availability, &asset.Status,
			&serialNumber, &pcNumber, &model, &cpu, &ram, &hddInfo, &networkCard, &opticalDrive,
			&ipAddress, &macAddress, &manufacturer, &purchaseYear, &warrantyUntil,
			&asset.CreatedAt,
			&asset.UpdatedAt, &asset.DeletedAt)
		if err != nil {
			return nil, err
		}
		if ownerName.Valid {
			asset.OwnerName = &ownerName.String
		}
		if responsibleUserName.Valid {
			asset.ResponsibleUserName = &responsibleUserName.String
		}
		populatePassportFieldsFromRow(&asset, serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer, purchaseYear, warrantyUntil)
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

	// Build query for data with owner name and responsible user name
	query := `
		SELECT a.id, a.tenant_id, a.inventory_number, a.name, a.type, a.class, a.owner_id, 
		       COALESCE(u_owner.first_name || ' ' || u_owner.last_name, u_owner.email) as owner_name,
		       a.responsible_user_id,
		       COALESCE(u_resp.first_name || ' ' || u_resp.last_name, u_resp.email) as responsible_user_name,
		       a.location, a.criticality, a.confidentiality, a.integrity, a.availability, 
		       a.status, a.serial_number, a.pc_number, a.model, a.cpu, a.ram, a.hdd_info, a.network_card, a.optical_drive,
		       a.ip_address, a.mac_address, a.manufacturer, a.purchase_year, a.warranty_until,
		       a.created_at, a.updated_at, a.deleted_at
		FROM assets a
		LEFT JOIN users u_owner ON a.owner_id = u_owner.id
		LEFT JOIN users u_resp ON a.responsible_user_id = u_resp.id
		WHERE a.tenant_id = $1 AND a.deleted_at IS NULL
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
		var ownerName, responsibleUserName sql.NullString
		var serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer sql.NullString
		var purchaseYear sql.NullInt32
		var warrantyUntil sql.NullTime
		err := rows.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
			&asset.Type, &asset.Class, &asset.OwnerID, &ownerName, &asset.ResponsibleUserID,
			&responsibleUserName, &asset.Location, &asset.Criticality, &asset.Confidentiality,
			&asset.Integrity, &asset.Availability, &asset.Status,
			&serialNumber, &pcNumber, &model, &cpu, &ram, &hddInfo, &networkCard, &opticalDrive,
			&ipAddress, &macAddress, &manufacturer, &purchaseYear, &warrantyUntil,
			&asset.CreatedAt,
			&asset.UpdatedAt, &asset.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		if ownerName.Valid {
			asset.OwnerName = &ownerName.String
		}
		if responsibleUserName.Valid {
			asset.ResponsibleUserName = &responsibleUserName.String
		}
		populatePassportFieldsFromRow(&asset, serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer, purchaseYear, warrantyUntil)
		assets = append(assets, asset)
	}

	return assets, total, nil
}

func (r *AssetRepo) Update(ctx context.Context, asset Asset) error {
	_, err := r.db.Exec(`
		UPDATE assets SET name = $1, type = $2, class = $3, owner_id = $4, responsible_user_id = $5, location = $6,
		                  criticality = $7, confidentiality = $8, integrity = $9, availability = $10,
		                  status = $11,
		                  serial_number = $12, pc_number = $13, model = $14, cpu = $15, ram = $16,
		                  hdd_info = $17, network_card = $18, optical_drive = $19,
		                  ip_address = $20, mac_address = $21, manufacturer = $22,
		                  purchase_year = $23, warranty_until = $24,
		                  updated_at = CURRENT_TIMESTAMP
		WHERE id = $25
	`, asset.Name, asset.Type, asset.Class, asset.OwnerID, asset.ResponsibleUserID, asset.Location,
		asset.Criticality, asset.Confidentiality, asset.Integrity,
		asset.Availability, asset.Status,
		asset.SerialNumber, asset.PCNumber, asset.Model, asset.CPU, asset.RAM,
		asset.HDDInfo, asset.NetworkCard, asset.OpticalDrive,
		asset.IPAddress, asset.MACAddress, asset.Manufacturer,
		asset.PurchaseYear, asset.WarrantyUntil,
		asset.ID)
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
	// Убеждаемся, что documents не nil
	if documents == nil {
		documents = []AssetDocument{}
	}

	// Get software
	software, err := r.GetAssetSoftware(ctx, id)
	if err != nil {
		return nil, err
	}
	// Убеждаемся, что software не nil
	if software == nil {
		software = []AssetSoftware{}
	}

	// Get history
	history, err := r.GetAssetHistory(ctx, id)
	if err != nil {
		return nil, err
	}
	// Убеждаемся, что history не nil
	if history == nil {
		history = []AssetHistory{}
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
	// Используем document_links для получения документов, связанных с активом
	// JOIN с users для получения имени и email создателя
	rows, err := r.db.QueryContext(ctx, `
		SELECT d.id, dl.entity_id as asset_id, 
		       COALESCE(d.category, 'other') as document_type,
		       d.storage_uri as file_path,
		       d.title,
		       d.mime_type as mime,
		       d.size_bytes,
		       d.created_by,
		       COALESCE(u.first_name || ' ' || u.last_name, u.email) as created_by_name,
		       u.email as created_by_email,
		       d.created_at
		FROM documents d
		INNER JOIN document_links dl ON d.id = dl.document_id
		LEFT JOIN users u ON d.created_by = u.id
		WHERE d.deleted_at IS NULL 
		  AND dl.module = 'assets'
		  AND dl.entity_id = $1
		ORDER BY d.created_at DESC
	`, assetID)
	if err != nil {
		log.Printf("ERROR: GetAssetDocuments query failed for assetID=%s: %v", assetID, err)
		return nil, err
	}
	defer rows.Close()

	var documents []AssetDocument
	for rows.Next() {
		var doc AssetDocument
		var createdByName, createdByEmail sql.NullString
		err := rows.Scan(&doc.ID, &doc.AssetID, &doc.DocumentType, &doc.FilePath, &doc.Title, &doc.Mime, &doc.SizeBytes, &doc.CreatedBy, &createdByName, &createdByEmail, &doc.CreatedAt)
		if err != nil {
			log.Printf("ERROR: GetAssetDocuments scan failed: %v", err)
			return nil, err
		}
		if createdByName.Valid {
			doc.CreatedByName = &createdByName.String
		}
		if createdByEmail.Valid {
			doc.CreatedByEmail = &createdByEmail.String
		}
		documents = append(documents, doc)
	}

	log.Printf("DEBUG: GetAssetDocuments assetID=%s returned %d documents", assetID, len(documents))
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
	// JOIN с users для получения имени и email пользователя
	rows, err := r.db.Query(`
		SELECT h.id, h.asset_id, h.field_changed, h.old_value, h.new_value, h.changed_by, 
		       COALESCE(u.first_name || ' ' || u.last_name, u.email) as changed_by_name,
		       u.email as changed_by_email,
		       h.changed_at
		FROM asset_history h
		LEFT JOIN users u ON h.changed_by = u.id
		WHERE h.asset_id = $1 
		ORDER BY h.changed_at DESC
	`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []AssetHistory
	for rows.Next() {
		var h AssetHistory
		var changedByName, changedByEmail sql.NullString
		err := rows.Scan(&h.ID, &h.AssetID, &h.FieldChanged, &h.OldValue, &h.NewValue, &h.ChangedBy, &changedByName, &changedByEmail, &h.ChangedAt)
		if err != nil {
			return nil, err
		}
		if changedByName.Valid {
			h.ChangedByName = &changedByName.String
		}
		if changedByEmail.Valid {
			h.ChangedByEmail = &changedByEmail.String
		}
		history = append(history, h)
	}
	return history, nil
}

// GetAssetHistoryWithFilters returns asset history with optional filters
func (r *AssetRepo) GetAssetHistoryWithFilters(ctx context.Context, assetID string, filters map[string]interface{}) ([]AssetHistory, error) {
	// JOIN с users для получения имени и email пользователя
	query := `
		SELECT h.id, h.asset_id, h.field_changed, h.old_value, h.new_value, h.changed_by,
		       COALESCE(u.first_name || ' ' || u.last_name, u.email) as changed_by_name,
		       u.email as changed_by_email,
		       h.changed_at
		FROM asset_history h
		LEFT JOIN users u ON h.changed_by = u.id
		WHERE h.asset_id = $1
	`
	args := []interface{}{assetID}
	argIndex := 2

	// Apply filters
	if changedBy, ok := filters["changed_by"].(string); ok && changedBy != "" {
		query += fmt.Sprintf(" AND h.changed_by = $%d", argIndex)
		args = append(args, changedBy)
		argIndex++
	}
	if fromDate, ok := filters["from_date"].(string); ok && fromDate != "" {
		query += fmt.Sprintf(" AND h.changed_at >= $%d", argIndex)
		args = append(args, fromDate)
		argIndex++
	}
	if toDate, ok := filters["to_date"].(string); ok && toDate != "" {
		query += fmt.Sprintf(" AND h.changed_at <= $%d", argIndex)
		args = append(args, toDate)
		argIndex++
	}

	query += " ORDER BY h.changed_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []AssetHistory
	for rows.Next() {
		var h AssetHistory
		var changedByName, changedByEmail sql.NullString
		err := rows.Scan(&h.ID, &h.AssetID, &h.FieldChanged, &h.OldValue, &h.NewValue, &h.ChangedBy, &changedByName, &changedByEmail, &h.ChangedAt)
		if err != nil {
			return nil, err
		}
		if changedByName.Valid {
			h.ChangedByName = &changedByName.String
		}
		if changedByEmail.Valid {
			h.ChangedByEmail = &changedByEmail.String
		}
		history = append(history, h)
	}
	return history, nil
}

// DeleteDocument removes a document from an asset
func (r *AssetRepo) DeleteDocument(ctx context.Context, documentID string) error {
	_, err := r.db.Exec(`
		DELETE FROM asset_documents WHERE id = $1
	`, documentID)
	return err
}

// GetDocumentByID returns a specific document
func (r *AssetRepo) GetDocumentByID(ctx context.Context, documentID string) (*AssetDocument, error) {
	var doc AssetDocument
	err := r.db.QueryRow(`
		SELECT id, asset_id, document_type, file_path, title, mime, size_bytes, created_by, created_at
		FROM asset_documents WHERE id = $1
	`, documentID).Scan(&doc.ID, &doc.AssetID, &doc.DocumentType, &doc.FilePath, &doc.Title, &doc.Mime, &doc.SizeBytes, &doc.CreatedBy, &doc.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// GetAssetRisks returns risks associated with an asset
func (r *AssetRepo) GetAssetRisks(ctx context.Context, assetID string) ([]Risk, error) {
	rows, err := r.db.Query(`
		SELECT r.id, r.tenant_id, r.title, r.description, r.category, r.likelihood, r.impact, r.level, r.status, r.owner_user_id, r.asset_id, r.created_at, r.updated_at
		FROM risks r
		WHERE r.asset_id = $1
		ORDER BY r.created_at DESC
	`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var risks []Risk
	for rows.Next() {
		var risk Risk
		err := rows.Scan(&risk.ID, &risk.TenantID, &risk.Title, &risk.Description, &risk.Category, &risk.Likelihood, &risk.Impact, &risk.Level, &risk.Status, &risk.OwnerUserID, &risk.AssetID, &risk.CreatedAt, &risk.UpdatedAt)
		if err != nil {
			return nil, err
		}
		risks = append(risks, risk)
	}
	return risks, nil
}

// GetAssetIncidents returns incidents associated with an asset
func (r *AssetRepo) GetAssetIncidents(ctx context.Context, assetID string) ([]Incident, error) {
	rows, err := r.db.Query(`
		SELECT i.id, i.tenant_id, i.title, i.description, i.category, i.status, i.severity, i.source, i.reported_by, i.assigned_to, i.asset_id, i.risk_id, i.created_by, i.detected_at, i.resolved_at, i.closed_at, i.created_at, i.updated_at, i.deleted_at
		FROM incidents i
		WHERE i.asset_id = $1
		ORDER BY i.created_at DESC
	`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []Incident
	for rows.Next() {
		var incident Incident
		err := rows.Scan(&incident.ID, &incident.TenantID, &incident.Title, &incident.Description, &incident.Category, &incident.Status, &incident.Severity, &incident.Source, &incident.ReportedBy, &incident.AssignedTo, &incident.AssetID, &incident.RiskID, &incident.CreatedBy, &incident.DetectedAt, &incident.ResolvedAt, &incident.ClosedAt, &incident.CreatedAt, &incident.UpdatedAt, &incident.DeletedAt)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}
	return incidents, nil
}

// GetAssetsWithoutOwner returns assets without owner
func (r *AssetRepo) GetAssetsWithoutOwner(ctx context.Context, tenantID string) ([]Asset, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.tenant_id, a.inventory_number, a.name, a.type, a.class, a.owner_id, 
		       COALESCE(u_owner.first_name || ' ' || u_owner.last_name, u_owner.email) as owner_name,
		       a.responsible_user_id,
		       COALESCE(u_resp.first_name || ' ' || u_resp.last_name, u_resp.email) as responsible_user_name,
		       a.location, a.criticality, a.confidentiality, a.integrity, a.availability, 
		       a.status, a.serial_number, a.pc_number, a.model, a.cpu, a.ram, a.hdd_info, a.network_card, a.optical_drive,
		       a.ip_address, a.mac_address, a.manufacturer, a.purchase_year, a.warranty_until,
		       a.created_at, a.updated_at, a.deleted_at
		FROM assets a
		LEFT JOIN users u_owner ON a.owner_id = u_owner.id
		LEFT JOIN users u_resp ON a.responsible_user_id = u_resp.id
		WHERE a.tenant_id = $1 AND a.deleted_at IS NULL AND a.owner_id IS NULL
		ORDER BY a.created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		var ownerName, responsibleUserName sql.NullString
		var serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer sql.NullString
		var purchaseYear sql.NullInt32
		var warrantyUntil sql.NullTime
		err := rows.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
			&asset.Type, &asset.Class, &asset.OwnerID, &ownerName, &asset.ResponsibleUserID,
			&responsibleUserName, &asset.Location, &asset.Criticality, &asset.Confidentiality,
			&asset.Integrity, &asset.Availability, &asset.Status,
			&serialNumber, &pcNumber, &model, &cpu, &ram, &hddInfo, &networkCard, &opticalDrive,
			&ipAddress, &macAddress, &manufacturer, &purchaseYear, &warrantyUntil,
			&asset.CreatedAt,
			&asset.UpdatedAt, &asset.DeletedAt)
		if err != nil {
			return nil, err
		}
		if ownerName.Valid {
			asset.OwnerName = &ownerName.String
		}
		if responsibleUserName.Valid {
			asset.ResponsibleUserName = &responsibleUserName.String
		}
		populatePassportFieldsFromRow(&asset, serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer, purchaseYear, warrantyUntil)
		assets = append(assets, asset)
	}
	return assets, nil
}

// GetAssetsWithoutPassport returns assets without passport document
func (r *AssetRepo) GetAssetsWithoutPassport(ctx context.Context, tenantID string) ([]Asset, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.tenant_id, a.inventory_number, a.name, a.type, a.class, a.owner_id, 
		       COALESCE(u_owner.first_name || ' ' || u_owner.last_name, u_owner.email) as owner_name,
		       a.responsible_user_id,
		       COALESCE(u_resp.first_name || ' ' || u_resp.last_name, u_resp.email) as responsible_user_name,
		       a.location, a.criticality, a.confidentiality, a.integrity, a.availability, 
		       a.status, a.serial_number, a.pc_number, a.model, a.cpu, a.ram, a.hdd_info, a.network_card, a.optical_drive,
		       a.ip_address, a.mac_address, a.manufacturer, a.purchase_year, a.warranty_until,
		       a.created_at, a.updated_at, a.deleted_at
		FROM assets a
		LEFT JOIN users u_owner ON a.owner_id = u_owner.id
		LEFT JOIN users u_resp ON a.responsible_user_id = u_resp.id
		WHERE a.tenant_id = $1 AND a.deleted_at IS NULL 
		  AND a.class = 'hardware'
		  AND (a.serial_number IS NULL OR a.serial_number = '' 
		       OR a.model IS NULL OR a.model = '' 
		       OR a.manufacturer IS NULL OR a.manufacturer = '')
		ORDER BY a.created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		var ownerName, responsibleUserName sql.NullString
		var serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer sql.NullString
		var purchaseYear sql.NullInt32
		var warrantyUntil sql.NullTime
		err := rows.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
			&asset.Type, &asset.Class, &asset.OwnerID, &ownerName, &asset.ResponsibleUserID,
			&responsibleUserName, &asset.Location, &asset.Criticality, &asset.Confidentiality,
			&asset.Integrity, &asset.Availability, &asset.Status,
			&serialNumber, &pcNumber, &model, &cpu, &ram, &hddInfo, &networkCard, &opticalDrive,
			&ipAddress, &macAddress, &manufacturer, &purchaseYear, &warrantyUntil,
			&asset.CreatedAt,
			&asset.UpdatedAt, &asset.DeletedAt)
		if err != nil {
			return nil, err
		}
		if ownerName.Valid {
			asset.OwnerName = &ownerName.String
		}
		if responsibleUserName.Valid {
			asset.ResponsibleUserName = &responsibleUserName.String
		}
		populatePassportFieldsFromRow(&asset, serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer, purchaseYear, warrantyUntil)
		assets = append(assets, asset)
	}
	return assets, nil
}

// GetAssetsWithoutCriticality returns assets without criticality assessment
func (r *AssetRepo) GetAssetsWithoutCriticality(ctx context.Context, tenantID string) ([]Asset, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.tenant_id, a.inventory_number, a.name, a.type, a.class, a.owner_id, 
		       COALESCE(u_owner.first_name || ' ' || u_owner.last_name, u_owner.email) as owner_name,
		       a.responsible_user_id,
		       COALESCE(u_resp.first_name || ' ' || u_resp.last_name, u_resp.email) as responsible_user_name,
		       a.location, a.criticality, a.confidentiality, a.integrity, a.availability, 
		       a.status, a.serial_number, a.pc_number, a.model, a.cpu, a.ram, a.hdd_info, a.network_card, a.optical_drive,
		       a.ip_address, a.mac_address, a.manufacturer, a.purchase_year, a.warranty_until,
		       a.created_at, a.updated_at, a.deleted_at
		FROM assets a
		LEFT JOIN users u_owner ON a.owner_id = u_owner.id
		LEFT JOIN users u_resp ON a.responsible_user_id = u_resp.id
		WHERE a.tenant_id = $1 AND a.deleted_at IS NULL 
		AND (a.criticality IS NULL OR a.criticality = '' OR a.confidentiality IS NULL OR a.confidentiality = '' 
		     OR a.integrity IS NULL OR a.integrity = '' OR a.availability IS NULL OR a.availability = '')
		ORDER BY a.created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		var ownerName, responsibleUserName sql.NullString
		var serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer sql.NullString
		var purchaseYear sql.NullInt32
		var warrantyUntil sql.NullTime
		err := rows.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
			&asset.Type, &asset.Class, &asset.OwnerID, &ownerName, &asset.ResponsibleUserID,
			&responsibleUserName, &asset.Location, &asset.Criticality, &asset.Confidentiality,
			&asset.Integrity, &asset.Availability, &asset.Status,
			&serialNumber, &pcNumber, &model, &cpu, &ram, &hddInfo, &networkCard, &opticalDrive,
			&ipAddress, &macAddress, &manufacturer, &purchaseYear, &warrantyUntil,
			&asset.CreatedAt,
			&asset.UpdatedAt, &asset.DeletedAt)
		if err != nil {
			return nil, err
		}
		if ownerName.Valid {
			asset.OwnerName = &ownerName.String
		}
		if responsibleUserName.Valid {
			asset.ResponsibleUserName = &responsibleUserName.String
		}
		populatePassportFieldsFromRow(&asset, serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer, purchaseYear, warrantyUntil)
		assets = append(assets, asset)
	}
	return assets, nil
}

func (r *AssetRepo) generateInventoryNumber(ctx context.Context, tenantID string) string {
	// Simple implementation - in production, this should be more sophisticated
	prefix := "AST"
	timestamp := time.Now().Format("20060102")
	random := strings.ToUpper(uuid.New().String()[:8])
	return fmt.Sprintf("%s-%s-%s", prefix, timestamp, random)
}

// AddDocumentWithFile adds a document with file details to an asset
func (r *AssetRepo) AddDocumentWithFile(ctx context.Context, assetID, documentID, documentType, filePath, fileName, mimeType string, fileSize int64, createdBy string) error {
	_, err := r.db.Exec(`
		INSERT INTO asset_documents (id, asset_id, document_type, file_path, title, mime, size_bytes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, documentID, assetID, documentType, filePath, fileName, mimeType, fileSize, createdBy)
	return err
}

// GetDocumentFromStorage retrieves a document from storage
func (r *AssetRepo) GetDocumentFromStorage(ctx context.Context, documentID string) (*AssetDocument, error) {
	var doc AssetDocument
	err := r.db.QueryRow(`
		SELECT id, asset_id, document_type, file_path, title, mime, size_bytes, created_by, created_at
		FROM asset_documents WHERE id = $1
	`, documentID).Scan(&doc.ID, &doc.AssetID, &doc.DocumentType, &doc.FilePath, &doc.Title, &doc.Mime, &doc.SizeBytes, &doc.CreatedBy, &doc.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// LinkDocumentToAsset links an existing document to an asset
func (r *AssetRepo) LinkDocumentToAsset(ctx context.Context, assetID, documentID, storageDocumentID, documentType, createdBy string) error {
	_, err := r.db.Exec(`
		INSERT INTO asset_documents (id, asset_id, document_type, file_path, title, mime, size_bytes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, documentID, assetID, documentType, storageDocumentID, "", "", 0, createdBy)
	return err
}

// GetUserResponsibleAssets returns assets for which the user is responsible
func (r *AssetRepo) GetUserResponsibleAssets(ctx context.Context, tenantID, userID string) ([]Asset, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.tenant_id, a.inventory_number, a.name, a.type, a.class, a.owner_id, 
		       COALESCE(u_owner.first_name || ' ' || u_owner.last_name, u_owner.email) as owner_name,
		       a.responsible_user_id,
		       COALESCE(u_resp.first_name || ' ' || u_resp.last_name, u_resp.email) as responsible_user_name,
		       a.location, a.criticality, a.confidentiality, a.integrity, a.availability, 
		       a.status, a.serial_number, a.pc_number, a.model, a.cpu, a.ram, a.hdd_info, a.network_card, a.optical_drive,
		       a.ip_address, a.mac_address, a.manufacturer, a.purchase_year, a.warranty_until,
		       a.created_at, a.updated_at, a.deleted_at
		FROM assets a
		LEFT JOIN users u_owner ON a.owner_id = u_owner.id
		LEFT JOIN users u_resp ON a.responsible_user_id = u_resp.id
		WHERE a.tenant_id = $1 AND a.responsible_user_id = $2 AND a.deleted_at IS NULL
		ORDER BY a.created_at DESC
	`, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		var ownerName, responsibleUserName sql.NullString
		var serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer sql.NullString
		var purchaseYear sql.NullInt32
		var warrantyUntil sql.NullTime
		err := rows.Scan(&asset.ID, &asset.TenantID, &asset.InventoryNumber, &asset.Name,
			&asset.Type, &asset.Class, &asset.OwnerID, &ownerName, &asset.ResponsibleUserID,
			&responsibleUserName, &asset.Location, &asset.Criticality, &asset.Confidentiality,
			&asset.Integrity, &asset.Availability, &asset.Status,
			&serialNumber, &pcNumber, &model, &cpu, &ram, &hddInfo, &networkCard, &opticalDrive,
			&ipAddress, &macAddress, &manufacturer, &purchaseYear, &warrantyUntil,
			&asset.CreatedAt,
			&asset.UpdatedAt, &asset.DeletedAt)
		if err != nil {
			return nil, err
		}
		if ownerName.Valid {
			asset.OwnerName = &ownerName.String
		}
		if responsibleUserName.Valid {
			asset.ResponsibleUserName = &responsibleUserName.String
		}
		populatePassportFieldsFromRow(&asset, serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer, purchaseYear, warrantyUntil)
		assets = append(assets, asset)
	}
	return assets, nil
}

// GetDocumentStorage returns documents from storage with pagination and filtering
func (r *AssetRepo) GetDocumentStorage(ctx context.Context, tenantID string, req dto.DocumentStorageRequest) ([]dto.DocumentStorageResponse, int64, error) {
	// Базовый запрос для получения документов из централизованного хранилища
	query := `
		SELECT d.id, COALESCE(d.title, '') as title, d.category as document_type, 
		       d.version, d.size_bytes, d.mime_type as mime, d.created_by, d.created_at
		FROM documents d
		WHERE d.tenant_id = $1 AND d.deleted_at IS NULL`

	args := []interface{}{tenantID}
	argIndex := 2

	// Фильтр по типу
	if req.Type != "" {
		query += fmt.Sprintf(" AND d.type = $%d", argIndex)
		args = append(args, req.Type)
		argIndex++
	}

	// Фильтр по поисковому запросу
	if req.Query != "" {
		query += fmt.Sprintf(" AND (d.title ILIKE $%d OR d.description ILIKE $%d)", argIndex, argIndex)
		searchTerm := "%" + req.Query + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	// Подсчет общего количества
	countQuery := "SELECT COUNT(*) FROM (" + query + ") as count_query"
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Пагинация
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 25
	}
	offset := (req.Page - 1) * req.PageSize
	query += fmt.Sprintf(" ORDER BY d.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, req.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var documents []dto.DocumentStorageResponse
	for rows.Next() {
		var doc dto.DocumentStorageResponse
		err := rows.Scan(&doc.ID, &doc.Title, &doc.DocumentType,
			&doc.Version, &doc.SizeBytes, &doc.Mime, &doc.CreatedBy, &doc.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		documents = append(documents, doc)
	}

	return documents, total, nil
}

func populatePassportFieldsFromRow(asset *Asset, serialNumber, pcNumber, model, cpu, ram, hddInfo, networkCard, opticalDrive, ipAddress, macAddress, manufacturer sql.NullString, purchaseYear sql.NullInt32, warrantyUntil sql.NullTime) {
	if serialNumber.Valid {
		value := serialNumber.String
		asset.SerialNumber = &value
	} else {
		asset.SerialNumber = nil
	}

	if pcNumber.Valid {
		value := pcNumber.String
		asset.PCNumber = &value
	} else {
		asset.PCNumber = nil
	}

	if model.Valid {
		value := model.String
		asset.Model = &value
	} else {
		asset.Model = nil
	}

	if cpu.Valid {
		value := cpu.String
		asset.CPU = &value
	} else {
		asset.CPU = nil
	}

	if ram.Valid {
		value := ram.String
		asset.RAM = &value
	} else {
		asset.RAM = nil
	}

	if hddInfo.Valid {
		value := hddInfo.String
		asset.HDDInfo = &value
	} else {
		asset.HDDInfo = nil
	}

	if networkCard.Valid {
		value := networkCard.String
		asset.NetworkCard = &value
	} else {
		asset.NetworkCard = nil
	}

	if opticalDrive.Valid {
		value := opticalDrive.String
		asset.OpticalDrive = &value
	} else {
		asset.OpticalDrive = nil
	}

	if ipAddress.Valid {
		value := ipAddress.String
		asset.IPAddress = &value
	} else {
		asset.IPAddress = nil
	}

	if macAddress.Valid {
		value := macAddress.String
		asset.MACAddress = &value
	} else {
		asset.MACAddress = nil
	}

	if manufacturer.Valid {
		value := manufacturer.String
		asset.Manufacturer = &value
	} else {
		asset.Manufacturer = nil
	}

	if purchaseYear.Valid {
		year := int(purchaseYear.Int32)
		asset.PurchaseYear = &year
	} else {
		asset.PurchaseYear = nil
	}

	if warrantyUntil.Valid {
		value := warrantyUntil.Time
		asset.WarrantyUntil = &value
	} else {
		asset.WarrantyUntil = nil
	}
}
