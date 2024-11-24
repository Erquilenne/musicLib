// @title           Music Library API
// @version         1.0
// @description     API Server for Music Library Application
// @host      localhost:5000
// @BasePath  /

package main

import (
	"log"
	"musiclib/config"
	"musiclib/internal/server"
	"musiclib/pkg/db/migrations"
	"musiclib/pkg/db/postgres"
	"musiclib/pkg/logger"
	_ "musiclib/docs" // Import swagger docs
	"os"
	"path/filepath"
)

func main() {
	log.Println("Starting api server")

	configPath := config.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %s", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %s", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s", "LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	// Подключение к базе данных
	db, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	appLogger.Infof("Postgres connected, Status: %#v", db.Stats())

	// Запуск миграций
	migrationsPath := filepath.Join("migrations")
	if err := migrations.RunMigrations(db.DB, migrationsPath); err != nil {
		appLogger.Fatalf("Could not run migrations: %v", err)
	}
	appLogger.Info("Migrations completed successfully")

	// Создание и запуск сервера
	srv := server.NewServer(cfg, db, appLogger)
	if err := srv.Run(); err != nil {
		appLogger.Fatalf("Error running server: %v", err)
	}

	defer db.Close()
}
