package domain

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

//go:embed templates/*.html
var templatesFS embed.FS

// TemplateService handles document template operations
type TemplateService struct {
	templateRepo    *repo.TemplateRepo
	assetRepo       *repo.AssetRepo
	documentService *DocumentService
}

// NewTemplateService creates a new template service
func NewTemplateService(templateRepo *repo.TemplateRepo, assetRepo *repo.AssetRepo, documentService *DocumentService) *TemplateService {
	return &TemplateService{
		templateRepo:    templateRepo,
		assetRepo:       assetRepo,
		documentService: documentService,
	}
}

// ListTemplates returns all templates for a tenant
func (s *TemplateService) ListTemplates(ctx context.Context, tenantID string, filters map[string]interface{}) ([]dto.DocumentTemplateDTO, error) {
	templates, err := s.templateRepo.ListTemplates(ctx, tenantID, filters)
	if err != nil {
		return nil, err
	}

	result := make([]dto.DocumentTemplateDTO, len(templates))
	for i, t := range templates {
		result[i] = t.ToDTO()
	}

	return result, nil
}

// GetTemplate retrieves a template by ID
func (s *TemplateService) GetTemplate(ctx context.Context, id, tenantID string) (*dto.DocumentTemplateDTO, error) {
	template, err := s.templateRepo.GetTemplateByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("template not found")
	}

	result := template.ToDTO()
	return &result, nil
}

// CreateTemplate creates a new template
func (s *TemplateService) CreateTemplate(ctx context.Context, tenantID, userID string, req dto.CreateTemplateRequest) (*dto.DocumentTemplateDTO, error) {
	// Validate template content
	if err := s.ValidateTemplate(req.Content); err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}

	template := &repo.DocumentTemplate{
		ID:           uuid.New().String(),
		TenantID:     tenantID,
		Name:         req.Name,
		Description:  req.Description,
		TemplateType: req.TemplateType,
		Content:      req.Content,
		IsSystem:     false, // User templates are never system templates
		IsActive:     true,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.templateRepo.CreateTemplate(ctx, template); err != nil {
		return nil, err
	}

	result := template.ToDTO()
	return &result, nil
}

// UpdateTemplate updates an existing template
func (s *TemplateService) UpdateTemplate(ctx context.Context, id, tenantID string, req dto.UpdateTemplateRequest) error {
	existing, err := s.templateRepo.GetTemplateByID(ctx, id, tenantID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("template not found")
	}

	// System templates cannot be modified
	if existing.IsSystem {
		return fmt.Errorf("cannot modify system templates")
	}

	// Update fields
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.TemplateType != nil {
		existing.TemplateType = *req.TemplateType
	}
	if req.Content != nil {
		if err := s.ValidateTemplate(*req.Content); err != nil {
			return fmt.Errorf("invalid template: %w", err)
		}
		existing.Content = *req.Content
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	return s.templateRepo.UpdateTemplate(ctx, existing)
}

// DeleteTemplate deletes a template
func (s *TemplateService) DeleteTemplate(ctx context.Context, id, tenantID string) error {
	return s.templateRepo.DeleteTemplate(ctx, id, tenantID)
}

// FillTemplate fills a template with asset data
func (s *TemplateService) FillTemplate(ctx context.Context, tenantID, userID string, req dto.FillTemplateRequest) (*dto.FillTemplateResponse, error) {
	log.Printf("DEBUG: FillTemplate called - SaveAsDocument=%v, GeneratePDF=%v, DocumentTitle=%s",
		req.SaveAsDocument, req.GeneratePDF, req.DocumentTitle)

	// Get template
	template, err := s.templateRepo.GetTemplateByID(ctx, req.TemplateID, tenantID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("template not found")
	}

	// Get asset
	asset, err := s.assetRepo.GetByID(ctx, req.AssetID)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, fmt.Errorf("asset not found")
	}

	// Prepare data for template rendering
	data := s.prepareTemplateData(asset, req.AdditionalData)

	// Render template
	html := s.RenderTemplate(template.Content, data)

	response := &dto.FillTemplateResponse{
		HTML: html,
	}

	// Generate PDF if requested
	if req.GeneratePDF {
		pdfBytes, err := s.GeneratePDFFromHTML(ctx, html)
		if err != nil {
			log.Printf("ERROR: Failed to generate PDF: %v", err)
			return nil, fmt.Errorf("failed to generate PDF: %w", err)
		}
		// Encode to base64 for JSON transport
		pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)
		response.PDFBase64 = &pdfBase64
	}

	// Save as document if requested
	if req.SaveAsDocument {
		title := fmt.Sprintf("Паспорт %s", asset.Name)
		if req.DocumentTitle != "" {
			title = req.DocumentTitle
		}

		// Determine content, file name, and MIME type based on format
		var content []byte
		var mimeType string
		var fileName string

		if req.GeneratePDF {
			// Save as PDF
			pdfBytes, err := s.GeneratePDFFromHTML(ctx, html)
			if err != nil {
				log.Printf("ERROR: Failed to generate PDF for saving: %v", err)
				return nil, fmt.Errorf("failed to generate PDF: %w", err)
			}
			content = pdfBytes
			mimeType = "application/pdf"
			fileName = fmt.Sprintf("%s.pdf", uuid.New().String())
		} else {
			// Save as HTML (if PDF not requested)
			content = []byte(html)
			mimeType = "text/html"
			fileName = fmt.Sprintf("%s.html", uuid.New().String())
		}

		// Save document to storage
		documentID, err := s.saveDocumentAsBytes(ctx, tenantID, userID, req.AssetID, title, content, fileName, mimeType, template.TemplateType)
		if err != nil {
			log.Printf("ERROR: Failed to save document: %v", err)
			return nil, fmt.Errorf("failed to save document: %w", err)
		}

		response.DocumentID = &documentID
		log.Printf("DEBUG: Document saved successfully: %s (format: %s)", documentID, mimeType)
	}

	return response, nil
}

// prepareTemplateData prepares data for template rendering
func (s *TemplateService) prepareTemplateData(asset *repo.Asset, additionalData map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})

	// Asset basic fields
	data["asset_name"] = asset.Name
	data["inventory_number"] = asset.InventoryNumber
	data["asset_type"] = asset.Type
	data["asset_class"] = asset.Class
	data["criticality"] = asset.Criticality
	data["confidentiality"] = asset.Confidentiality
	data["integrity"] = asset.Integrity
	data["availability"] = asset.Availability
	data["status"] = asset.Status

	// Optional fields
	if asset.OwnerName != nil {
		data["owner_name"] = *asset.OwnerName
	}
	if asset.ResponsibleUserName != nil {
		data["responsible_user_name"] = *asset.ResponsibleUserName
	}
	if asset.Location != nil {
		data["location"] = *asset.Location
	}

	// Passport-specific fields (will be empty if not set)
	data["serial_number"] = safeString(asset.SerialNumber)
	data["pc_number"] = safeString(asset.PCNumber)
	data["model"] = safeString(asset.Model)
	data["cpu"] = safeString(asset.CPU)
	data["ram"] = safeString(asset.RAM)
	data["hdd_info"] = safeString(asset.HDDInfo)
	data["network_card"] = safeString(asset.NetworkCard)
	data["optical_drive"] = safeString(asset.OpticalDrive)
	data["ip_address"] = safeString(asset.IPAddress)
	data["mac_address"] = safeString(asset.MACAddress)
	data["manufacturer"] = safeString(asset.Manufacturer)

	if asset.PurchaseYear != nil {
		data["purchase_year"] = strconv.Itoa(*asset.PurchaseYear)
	} else {
		data["purchase_year"] = ""
	}

	if asset.WarrantyUntil != nil {
		data["warranty_until"] = asset.WarrantyUntil.Format("02.01.2006")
	} else {
		data["warranty_until"] = ""
	}

	// Date fields
	data["current_date"] = time.Now().Format("02.01.2006")
	data["current_datetime"] = time.Now().Format("02.01.2006 15:04")
	data["created_at"] = asset.CreatedAt.Format("02.01.2006")
	data["updated_at"] = asset.UpdatedAt.Format("02.01.2006")

	// Merge additional data
	for k, v := range additionalData {
		data[k] = v
	}

	return data
}

// safeString returns empty string if pointer is nil
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// RenderTemplate replaces placeholders in template with actual values
func (s *TemplateService) RenderTemplate(template string, data map[string]interface{}) string {
	result := template

	// Replace all {{variable}} placeholders
	re := regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)
	result = re.ReplaceAllStringFunc(result, func(match string) string {
		// Extract variable name
		varName := strings.Trim(match, "{}")
		varName = strings.TrimSpace(varName)

		// Look up value in data
		if val, ok := data[varName]; ok {
			return fmt.Sprintf("%v", val)
		}

		// Return empty string if variable not found
		return ""
	})

	return result
}

// ValidateTemplate validates template syntax
func (s *TemplateService) ValidateTemplate(content string) error {
	// Check for balanced {{ }}
	openCount := strings.Count(content, "{{")
	closeCount := strings.Count(content, "}}")

	if openCount != closeCount {
		return fmt.Errorf("unbalanced template placeholders")
	}

	// Check for valid variable names
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		varName := strings.TrimSpace(match[1])
		if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(varName) {
			return fmt.Errorf("invalid variable name: %s", varName)
		}
	}

	return nil
}

// GetAvailableVariables returns list of available template variables
func (s *TemplateService) GetAvailableVariables() dto.TemplateVariablesResponse {
	variables := []dto.TemplateVariable{
		// Asset basic fields
		{Name: "asset_name", Placeholder: "{{asset_name}}", Description: "Название актива", Example: "Рабочая станция №1", Category: "asset"},
		{Name: "inventory_number", Placeholder: "{{inventory_number}}", Description: "Инвентарный номер", Example: "РСП-2025-0001", Category: "asset"},
		{Name: "asset_type", Placeholder: "{{asset_type}}", Description: "Тип актива", Example: "hardware", Category: "asset"},
		{Name: "asset_class", Placeholder: "{{asset_class}}", Description: "Класс актива", Example: "workstation", Category: "asset"},

		// CIA fields
		{Name: "criticality", Placeholder: "{{criticality}}", Description: "Критичность", Example: "high", Category: "asset"},
		{Name: "confidentiality", Placeholder: "{{confidentiality}}", Description: "Конфиденциальность", Example: "medium", Category: "asset"},
		{Name: "integrity", Placeholder: "{{integrity}}", Description: "Целостность", Example: "high", Category: "asset"},
		{Name: "availability", Placeholder: "{{availability}}", Description: "Доступность", Example: "high", Category: "asset"},
		{Name: "status", Placeholder: "{{status}}", Description: "Статус", Example: "active", Category: "asset"},

		// Passport fields
		{Name: "serial_number", Placeholder: "{{serial_number}}", Description: "Серийный номер (S/N)", Example: "ABC123456", Category: "passport"},
		{Name: "pc_number", Placeholder: "{{pc_number}}", Description: "Номер ПК", Example: "PC-101", Category: "passport"},
		{Name: "model", Placeholder: "{{model}}", Description: "Модель", Example: "Dell OptiPlex 7090", Category: "passport"},
		{Name: "cpu", Placeholder: "{{cpu}}", Description: "Процессор", Example: "Intel Core i7-11700", Category: "passport"},
		{Name: "ram", Placeholder: "{{ram}}", Description: "Оперативная память", Example: "16 GB DDR4", Category: "passport"},
		{Name: "hdd_info", Placeholder: "{{hdd_info}}", Description: "Информация о HDD", Example: "SSD 512GB", Category: "passport"},
		{Name: "network_card", Placeholder: "{{network_card}}", Description: "Сетевая карта", Example: "Intel I219-V", Category: "passport"},
		{Name: "optical_drive", Placeholder: "{{optical_drive}}", Description: "Оптический привод", Example: "DVD-RW", Category: "passport"},
		{Name: "ip_address", Placeholder: "{{ip_address}}", Description: "IP адрес", Example: "192.168.1.100", Category: "passport"},
		{Name: "mac_address", Placeholder: "{{mac_address}}", Description: "MAC адрес", Example: "00:1A:2B:3C:4D:5E", Category: "passport"},
		{Name: "manufacturer", Placeholder: "{{manufacturer}}", Description: "Производитель", Example: "Dell Inc.", Category: "passport"},
		{Name: "purchase_year", Placeholder: "{{purchase_year}}", Description: "Год приобретения", Example: "2023", Category: "passport"},
		{Name: "warranty_until", Placeholder: "{{warranty_until}}", Description: "Гарантия до", Example: "31.12.2026", Category: "passport"},

		// User fields
		{Name: "owner_name", Placeholder: "{{owner_name}}", Description: "Владелец", Example: "Иванов И.И.", Category: "user"},
		{Name: "responsible_user_name", Placeholder: "{{responsible_user_name}}", Description: "Ответственный пользователь", Example: "Петров П.П.", Category: "user"},
		{Name: "location", Placeholder: "{{location}}", Description: "Местоположение", Example: "Кабинет 101", Category: "asset"},

		// Date fields
		{Name: "current_date", Placeholder: "{{current_date}}", Description: "Текущая дата", Example: "08.10.2025", Category: "date"},
		{Name: "current_datetime", Placeholder: "{{current_datetime}}", Description: "Текущая дата и время", Example: "08.10.2025 14:30", Category: "date"},
		{Name: "created_at", Placeholder: "{{created_at}}", Description: "Дата создания актива", Example: "01.01.2025", Category: "date"},
		{Name: "updated_at", Placeholder: "{{updated_at}}", Description: "Дата обновления актива", Example: "08.10.2025", Category: "date"},
	}

	return dto.TemplateVariablesResponse{
		Variables: variables,
	}
}

// InitializeDefaultTemplates creates system templates from embedded files
func (s *TemplateService) InitializeDefaultTemplates(ctx context.Context, tenantID, userID string) error {
	log.Printf("Initializing default templates for tenant %s", tenantID)

	templates := []struct {
		filename     string
		name         string
		description  string
		templateType string
	}{
		{"passport_pc.html", "Паспорт персонального компьютера", "Стандартный паспорт для ПК", "passport_pc"},
		{"passport_monitor.html", "Паспорт монитора", "Стандартный паспорт для мониторов", "passport_monitor"},
		{"passport_device.html", "Паспорт устройства", "Универсальный паспорт для любого устройства", "passport_device"},
		{"passport_printer.html", "Паспорт принтера/МФУ", "Паспорт для печатающих устройств", "passport_device"},
		{"passport_network.html", "Паспорт сетевого оборудования", "Паспорт для роутеров, коммутаторов и другого сетевого оборудования", "passport_device"},
		{"passport_storage.html", "Паспорт съемного носителя", "Паспорт для USB-флешек, внешних HDD и других носителей", "passport_device"},
	}

	for _, tmpl := range templates {
		// Read template file from embedded FS
		content, err := templatesFS.ReadFile("templates/" + tmpl.filename)
		if err != nil {
			log.Printf("WARN: Could not read template file %s: %v", tmpl.filename, err)
			continue
		}

		// Check if template already exists
		existing, _ := s.templateRepo.ListTemplates(ctx, tenantID, map[string]interface{}{
			"template_type": tmpl.templateType,
			"is_system":     true,
		})

		if len(existing) > 0 {
			log.Printf("System template %s already exists, skipping", tmpl.name)
			continue
		}

		// Create template
		template := &repo.DocumentTemplate{
			ID:           uuid.New().String(),
			TenantID:     tenantID,
			Name:         tmpl.name,
			Description:  &tmpl.description,
			TemplateType: tmpl.templateType,
			Content:      string(content),
			IsSystem:     true,
			IsActive:     true,
			CreatedBy:    userID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.templateRepo.CreateTemplate(ctx, template); err != nil {
			log.Printf("ERROR: Failed to create system template %s: %v", tmpl.name, err)
			return err
		}

		log.Printf("Created system template: %s", tmpl.name)
	}

	return nil
}

// ListInventoryRules returns all inventory rules for a tenant
func (s *TemplateService) ListInventoryRules(ctx context.Context, tenantID string) ([]dto.InventoryNumberRuleDTO, error) {
	rules, err := s.templateRepo.ListInventoryRules(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.InventoryNumberRuleDTO, len(rules))
	for i, r := range rules {
		result[i] = r.ToDTO()
	}

	return result, nil
}

// CreateInventoryRule creates a new inventory rule
func (s *TemplateService) CreateInventoryRule(ctx context.Context, tenantID string, req dto.CreateInventoryRuleRequest) (*dto.InventoryNumberRuleDTO, error) {
	rule := &repo.InventoryNumberRule{
		ID:              uuid.New().String(),
		TenantID:        tenantID,
		AssetType:       req.AssetType,
		AssetClass:      req.AssetClass,
		Pattern:         req.Pattern,
		CurrentSequence: 0,
		Prefix:          req.Prefix,
		Description:     req.Description,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.templateRepo.CreateInventoryRule(ctx, rule); err != nil {
		return nil, err
	}

	result := rule.ToDTO()
	return &result, nil
}

// UpdateInventoryRule updates an existing rule
func (s *TemplateService) UpdateInventoryRule(ctx context.Context, id, tenantID string, req dto.UpdateInventoryRuleRequest) error {
	// Get existing rule
	existing, err := s.templateRepo.GetInventoryRuleByType(ctx, tenantID, "", nil)
	if err != nil {
		return err
	}
	if existing == nil || existing.ID != id {
		// Fallback: try to find by ID (would need to add this method to repo)
		return fmt.Errorf("inventory rule not found")
	}

	// Update fields
	if req.Pattern != nil {
		existing.Pattern = *req.Pattern
	}
	if req.CurrentSequence != nil {
		existing.CurrentSequence = *req.CurrentSequence
	}
	if req.Prefix != nil {
		existing.Prefix = req.Prefix
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	return s.templateRepo.UpdateInventoryRule(ctx, existing)
}

// GenerateInventoryNumber generates a new inventory number based on rules
func (s *TemplateService) GenerateInventoryNumber(ctx context.Context, tenantID, assetType string, assetClass *string) (*dto.GenerateInventoryNumberResponse, error) {
	// Get rule for this asset type
	rule, err := s.templateRepo.GetInventoryRuleByType(ctx, tenantID, assetType, assetClass)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, fmt.Errorf("no inventory number rule found for asset type: %s", assetType)
	}

	// Increment sequence
	newSequence, err := s.templateRepo.IncrementSequence(ctx, rule.ID)
	if err != nil {
		return nil, err
	}

	// Generate number from pattern
	inventoryNumber := s.renderInventoryPattern(rule.Pattern, assetType, assetClass, newSequence)

	return &dto.GenerateInventoryNumberResponse{
		InventoryNumber: inventoryNumber,
		Pattern:         rule.Pattern,
		Sequence:        newSequence,
	}, nil
}

// renderInventoryPattern renders pattern with variables
func (s *TemplateService) renderInventoryPattern(pattern, assetType string, assetClass *string, sequence int) string {
	result := pattern
	now := time.Now()

	// Replace variables
	result = strings.ReplaceAll(result, "{{type_code}}", assetType)
	result = strings.ReplaceAll(result, "{{year}}", fmt.Sprintf("%d", now.Year()))
	result = strings.ReplaceAll(result, "{{month}}", fmt.Sprintf("%02d", now.Month()))

	if assetClass != nil {
		result = strings.ReplaceAll(result, "{{class}}", *assetClass)
	}

	// Handle sequence with format (e.g., {{sequence:0000}})
	re := regexp.MustCompile(`\{\{sequence(?::(\d+))?\}\}`)
	result = re.ReplaceAllStringFunc(result, func(match string) string {
		matches := re.FindStringSubmatch(match)
		if len(matches) > 1 && matches[1] != "" {
			// Has format specifier
			width, _ := strconv.Atoi(matches[1])
			format := fmt.Sprintf("%%0%dd", width)
			return fmt.Sprintf(format, sequence)
		}
		// No format, just return number
		return fmt.Sprintf("%d", sequence)
	})

	return result
}

// GeneratePDFFromHTML generates a PDF document from HTML content using chromedp
func (s *TemplateService) GeneratePDFFromHTML(ctx context.Context, html string) ([]byte, error) {
	// Create chromedp context with proper options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Create a buffer to store PDF
	var pdfBuffer []byte

	// Encode HTML to base64 to avoid URL encoding issues
	htmlBase64 := base64.StdEncoding.EncodeToString([]byte(html))
	dataURL := "data:text/html;base64," + htmlBase64

	// Navigate to data URL with HTML content and print to PDF
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(dataURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond), // Give time for rendering
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuffer, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).   // A4 width in inches
				WithPaperHeight(11.69). // A4 height in inches
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				WithDisplayHeaderFooter(false).
				WithPreferCSSPageSize(false).
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		log.Printf("ERROR: GeneratePDFFromHTML chromedp failed: %v", err)
		return nil, fmt.Errorf("chromedp failed: %w", err)
	}

	log.Printf("DEBUG: GeneratePDFFromHTML generated PDF of size %d bytes", len(pdfBuffer))
	return pdfBuffer, nil
}

// saveDocumentAsBytes saves document content (PDF or HTML) as a document using DocumentService
func (s *TemplateService) saveDocumentAsBytes(ctx context.Context, tenantID, userID, assetID, title string, content []byte, fileName, mimeType, documentType string) (string, error) {
	log.Printf("DEBUG: saveDocumentAsBytes called - title=%s, assetID=%s, documentType=%s, mimeType=%s", title, assetID, documentType, mimeType)

	// Get asset info for metadata
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return "", fmt.Errorf("failed to get asset: %w", err)
	}
	if asset == nil {
		return "", fmt.Errorf("asset not found")
	}

	// Prepare upload request with correct tags for proper categorization
	uploadReq := dto.UploadDocumentDTO{
		Name:        title,
		Description: templateStringPtr(fmt.Sprintf("Сгенерированный паспорт для актива %s (%s)", asset.Name, asset.InventoryNumber)),
		FolderID:    nil,
		Tags:        []string{"#passport", "#активы"}, // #passport first for correct category detection
		LinkedTo: &dto.DocumentLinkDTO{
			Module:   "assets",
			EntityID: assetID,
		},
		Metadata: templateStringPtr(fmt.Sprintf(`{"asset_id": "%s", "asset_name": "%s", "document_type": "%s", "generated": true}`, assetID, asset.Name, documentType)),
	}

	// Save using DocumentService (handles all file operations, DB records, links, auditing)
	document, err := s.documentService.SaveGeneratedDocument(
		ctx,
		tenantID,
		content,
		fileName,
		mimeType,
		uploadReq,
		userID,
	)
	if err != nil {
		log.Printf("ERROR: Failed to save generated document: %v", err)
		return "", fmt.Errorf("failed to save document: %w", err)
	}

	log.Printf("DEBUG: Document saved successfully: ID=%s, Title=%s, AssetID=%s, Format=%s", document.ID, document.Title, assetID, mimeType)
	return document.ID, nil
}

// templateStringPtr returns pointer to string
func templateStringPtr(s string) *string {
	return &s
}
