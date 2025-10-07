package main

import (
	"log"
	"os"
	"time"

	"risknexus/backend/internal/cache"
	"risknexus/backend/internal/config"
	"risknexus/backend/internal/database"
	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/http"
	"risknexus/backend/internal/migrate"
	"risknexus/backend/internal/repo"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Check if we should run migrations
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Println("Running database migrations...")
		if err := migrate.RunMigrations(db.DB, "./migrations"); err != nil {
			log.Fatal("Migration failed:", err)
		}
		log.Println("Migrations completed successfully")
		return
	}

	// Auto-run migrations on startup
	log.Println("Running database migrations on startup...")
	if err := migrate.RunMigrations(db.DB, "./migrations"); err != nil {
		log.Printf("Warning: Migration failed: %v", err)
		// Don't fail startup, just log the warning
	} else {
		log.Println("Migrations completed successfully")
	}

	// Initialize cache
	memoryCache := cache.NewMemoryCache()

	// Initialize repositories
	userRepo := repo.NewUserRepo(db)
	baseRoleRepo := repo.NewRoleRepo(db)
	roleRepo := repo.NewCachedRoleRepo(baseRoleRepo, memoryCache, 5*time.Minute)
	permissionRepo := repo.NewPermissionRepo(db)
	tenantRepo := repo.NewTenantRepo(db)
	assetRepo := repo.NewAssetRepo(db)
	riskRepo := repo.NewRiskRepo(db)
	documentRepo := repo.NewDocumentRepo(db)
	incidentRepo := repo.NewIncidentRepository(db.DB)
	trainingRepo := repo.NewTrainingRepo(db)
	auditRepo := repo.NewAuditRepo(db)
	aiRepo := repo.NewAIRepo(db)
	complianceRepo := repo.NewComplianceRepo(db)
	emailChangeRepo := repo.NewEmailChangeRepo(db)

	// Initialize services
	authService := domain.NewAuthService(userRepo, baseRoleRepo, permissionRepo, cfg.JWTSecret)
	userService := domain.NewUserService(userRepo, baseRoleRepo)
	roleService := domain.NewRoleService(roleRepo, userRepo, auditRepo)
	tenantService := domain.NewTenantService(tenantRepo, auditRepo)
	assetService := domain.NewAssetService(assetRepo, userRepo)
	documentService := domain.NewDocumentService(documentRepo, "./storage/documents")
	riskService := domain.NewRiskService(riskRepo, auditRepo, documentService)
	incidentService := domain.NewIncidentService(incidentRepo, userRepo, assetRepo, riskRepo)
	trainingService := domain.NewTrainingService(trainingRepo)
	aiService := domain.NewAIService(aiRepo)
	complianceService := domain.NewComplianceService(complianceRepo)
	emailChangeService := domain.NewEmailChangeService(emailChangeRepo, userRepo)

	// Initialize handlers
	authHandler := http.NewAuthHandler(authService)
	userHandler := http.NewUserHandler(userService, roleService)
	roleHandler := http.NewRoleHandler(roleService)
	log.Printf("DEBUG: main.go roleHandler created: %+v", roleHandler)
	tenantHandler := http.NewTenantHandler(tenantService)
	auditHandler := http.NewAuditHandler(auditRepo)

	// Set global permission checker
	http.SetPermissionChecker(userService)
	assetHandler := http.NewAssetHandler(assetService)
	riskHandler := http.NewRiskHandler(riskService)
	documentHandler := http.NewDocumentHandler(documentService)
	incidentHandler := http.NewIncidentHandler(incidentService)
	trainingHandler := http.NewTrainingHandler(trainingService)
	aiHandler := http.NewAIHandler(aiService)
	complianceHandler := http.NewComplianceHandler(complianceService)
	emailChangeHandler := http.NewEmailChangeHandler(emailChangeService, validator.New())

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Simple test endpoint (no auth)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK", "message": "Backend is healthy"})
	})

	// Routes
	api := app.Group("/api")

	// Auth routes
	authHandler.Register(api)

	// Test endpoint (this should work without auth)
	api.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK", "message": "Backend is working"})
	})

	// Protected routes
	protected := api.Group("", http.AuthMiddleware(authService))
	// Register protected auth routes
	authHandler.RegisterProtected(protected)
	userHandler.Register(protected)
	log.Printf("DEBUG: main.go registering roleHandler")
	roleHandler.Register(protected)
	log.Printf("DEBUG: main.go roleHandler registered")
	tenantHandler.Register(protected)
	auditHandler.Register(protected)
	assetHandler.Register(protected)
	riskHandler.Register(protected)
	documentHandler.RegisterRoutes(protected)
	incidentHandler.Register(protected)
	trainingHandler.Register(protected)
	aiHandler.Register(protected)
	complianceHandler.Register(protected)
	emailChangeHandler.Register(protected)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
