package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"autostrike/internal/application"
	"autostrike/internal/domain/service"
	"autostrike/internal/infrastructure/api/rest"
	"autostrike/internal/infrastructure/persistence/sqlite"
	"autostrike/internal/infrastructure/websocket"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	// Register SQLite3 driver for database/sql
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load .env file (optional - won't fail if not found)
	_ = godotenv.Load()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	// Load configuration
	if err := loadConfig(); err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize database
	db, err := sql.Open("sqlite3", viper.GetString("database.path"))
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize schema
	if err := sqlite.InitSchema(db); err != nil {
		logger.Fatal("Failed to initialize schema", zap.Error(err))
	}

	// Initialize repositories
	agentRepo := sqlite.NewAgentRepository(db)
	scenarioRepo := sqlite.NewScenarioRepository(db)
	techniqueRepo := sqlite.NewTechniqueRepository(db)
	resultRepo := sqlite.NewResultRepository(db)

	// Initialize domain services
	validator := service.NewTechniqueValidator()
	orchestrator := service.NewAttackOrchestrator(agentRepo, techniqueRepo, validator, logger)
	calculator := service.NewScoreCalculator()

	// Initialize application services
	agentService := application.NewAgentService(agentRepo)
	scenarioService := application.NewScenarioService(scenarioRepo, techniqueRepo, validator)
	executionService := application.NewExecutionService(
		resultRepo,
		scenarioRepo,
		techniqueRepo,
		agentRepo,
		orchestrator,
		calculator,
	)
	techniqueService := application.NewTechniqueService(techniqueRepo)

	// Auto-import techniques from configs directory at startup
	autoImportTechniques(techniqueService, logger)

	// Initialize WebSocket hub
	hub := websocket.NewHub(logger)
	go hub.Run()

	// Initialize HTTP server
	server := rest.NewServer(
		agentService,
		scenarioService,
		executionService,
		techniqueService,
		hub,
		logger,
	)

	// Start server
	go func() {
		addr := viper.GetString("server.address")
		logger.Info("Starting AutoStrike server", zap.String("address", addr))
		if err := server.Run(addr); err != nil {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
}

func autoImportTechniques(service *application.TechniqueService, logger *zap.Logger) {
	// Import techniques from all YAML files in configs/techniques
	paths := []string{
		"./configs/techniques/discovery.yaml",
		"./configs/techniques/execution.yaml",
		"./configs/techniques/persistence.yaml",
		"./configs/techniques/defense-evasion.yaml",
	}

	imported := 0
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			if err := service.ImportTechniques(context.Background(), path); err != nil {
				logger.Warn("Failed to import techniques", zap.String("path", path), zap.Error(err))
			} else {
				imported++
				logger.Info("Imported techniques", zap.String("path", path))
			}
		}
	}

	if imported > 0 {
		logger.Info("Auto-imported technique files", zap.Int("count", imported))
	}
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Defaults
	viper.SetDefault("server.address", ":8443")
	viper.SetDefault("database.path", "./data/autostrike.db")
	viper.SetDefault("agent.beacon_interval", 30)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}
