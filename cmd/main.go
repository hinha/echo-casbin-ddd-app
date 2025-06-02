// main entrypoint
package main

import (
	"context"
	"errors"
	"flag"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/hinha/echo-casbin-ddd-app/docs" // Import Swagger docs
	"github.com/hinha/echo-casbin-ddd-app/internal/application/usecase"
	"github.com/hinha/echo-casbin-ddd-app/internal/config"
	"github.com/hinha/echo-casbin-ddd-app/internal/infrastructure/auth"
	"github.com/hinha/echo-casbin-ddd-app/internal/infrastructure/persistence"
	"github.com/hinha/echo-casbin-ddd-app/internal/interfaces/api/handler"
	"github.com/hinha/echo-casbin-ddd-app/internal/interfaces/api/middleware"
	"github.com/hinha/echo-casbin-ddd-app/internal/interfaces/websocket"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// @title Echo Casbin DDD App
// @version 1.0
// @description Example API with JWT, Casbin, PostgreSQL, WebSocket
// @termsOfService http://swagger.io/terms/
// @contact.name Martinus Dawan
// @contact.email martinus@example.com
// @license.name MIT
// @BasePath /api/v1
func main() {
	// Define command line flags
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	// Load configuration
	cfg := config.NewConfig()

	// Initialize database
	db, err := persistence.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database models if --migrate flag is provided
	if *migrateFlag {
		log.Println("Running database migrations...")
		if err := db.AutoMigrate(); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("Database migrations completed successfully")
	}

	// Initialize repositories
	userRepo := persistence.NewUserRepository(db.DB)
	apiClientRepo := persistence.NewAPIClientRepository(db.DB)

	// Initialize and run seeder
	if *migrateFlag {
		seeder := persistence.NewSeeder(cfg, userRepo)
		if err := seeder.Seed(context.Background()); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		log.Println("Database seeding users completed successfully")
	}

	// Initialize auth services
	jwtService := auth.NewJWTService(cfg)
	apiKeyService := auth.NewAPIKeyService(cfg, apiClientRepo)
	casbinService, err := auth.NewCasbinService(db.DB, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin service: %v", err)
	}

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, jwtService, casbinService)
	apiClientUseCase := usecase.NewAPIClientUseCase(apiClientRepo, casbinService)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUseCase)
	apiClientHandler := handler.NewAPIClientHandler(apiClientUseCase)

	// Initialize WebSocket handler
	userWSHandler := websocket.NewUserWSHandler(userUseCase)
	userWSHandler.Start()

	// Register routes
	jwtMiddleware := middleware.JWTMiddleware(cfg)
	apiKeyMiddleware := middleware.APIKeyMiddleware(cfg, apiKeyService)
	casbinMiddleware := middleware.CasbinMiddleware(casbinService)

	// User routes with JWT authentication
	userHandler.RegisterRoutes(e, jwtMiddleware)

	// API client routes with API key authentication and admin authorization
	apiClientHandler.RegisterRoutes(e, apiKeyMiddleware, casbinMiddleware)

	// Serve static files
	e.Static("/", "web")

	// Serve Swagger documentation
	//e.GET("/swagger", func(c echo.Context) error {
	//	return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	//})
	e.File("/swagger/doc.json", "docs/swagger.json")
	e.File("/swagger/doc.yaml", "docs/swagger.yaml")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	go func() {
		if err := e.Start(":" + cfg.Server.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop WebSocket handler
	userWSHandler.Stop()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shut down server: %v", err)
	}

	// Close database connection
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}

	log.Println("Server stopped")
}
