package main

import (
	"context"

	"loyalty-system/helpers"
	"loyalty-system/httphandler"
	"loyalty-system/models"

	"github.com/IvanSkripnikov/go-logger"
	"github.com/IvanSkripnikov/go-migrator"
)

func main() {
	logger.Debug("Service starting")

	// регистрация общих метрик
	helpers.RegisterCommonMetrics()
	logger.Debug("Metrics registered")

	// настройка всех конфигов
	config, err := models.LoadConfig()
	if err != nil {
		logger.Fatalf("Config error: %v", err)
	}
	// пробрасываем конфиг в helpers
	helpers.InitConfig(config)
	logger.Debug("Config initialized")

	// настройка коннекта к БД
	helpers.InitDatabase(config.Database)
	logger.Debug("Database initialized")

	// настройка коннекта к redis
	helpers.InitRedis(context.Background(), config.Redis)
	logger.Debug("Redis initialized")

	// выполнение миграций
	migrator.CreateTables(helpers.DB)
	logger.Debug("Migrations applied")

	// инициализация REST-api
	httphandler.InitHTTPServer()
	logger.Debug("Ape routes initialiazed")

	// запуск кронов
	helpers.InitTimer()
	logger.Debug("Cron commands initialized")

	// инициализация настроек системы лояльности
	helpers.LoadLoyaltyConfig()
	logger.Debug("Loyalty system configuration loaded")

	logger.Info("Service started")
}
