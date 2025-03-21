package main

import (
	"context"
	"time"

	"loyalty-system/helpers"
	"loyalty-system/httphandler"
	"loyalty-system/models"

	"github.com/IvanSkripnikov/go-logger"
	"github.com/IvanSkripnikov/go-migrator"

	"github.com/go-co-op/gocron"
)

func main() {
	logger.Debug("Service starting")

	// регистрация общих метрик
	helpers.RegisterCommonMetrics()

	// настройка всех конфигов
	config, err := models.LoadConfig()
	if err != nil {
		logger.Fatalf("Config error: %v", err)
	}

	// настройка коннекта к БД
	helpers.InitDatabase(config.Database)

	// настройка коннекта к redis
	helpers.InitRedis(context.Background(), config.Redis)

	// выполнение миграций
	migrator.CreateTables(helpers.DB)

	// инициализация REST-api
	httphandler.InitHTTPServer()

	// запуск кронов
	InitTimer()

	logger.Info("Service started")
}

func InitTimer() {
	//interval, _ := strconv.Atoi(os.Getenv("UPDATE_INTERVAL"))
	interval := 60
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(interval).Seconds().Do(func() {
		logger.Info("Updating loyalty")
	})
	scheduler.StartAsync()
}
