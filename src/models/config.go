package models

import (
	"os"
	"strconv"

	"github.com/IvanSkripnikov/go-gormdb"
)

type Config struct {
	Database          gormdb.Database
	Redis             Redis
	ShopServiceUrl    string
	OrdersServiceUrl  string
	PaymentServiceUrl string
	BillingServiceUrl string
}

func LoadConfig() (*Config, error) {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Database: gormdb.Database{
			Address:  os.Getenv("DB_ADDRESS"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DB:       os.Getenv("DB_NAME"),
		},
		Redis: Redis{
			Address:  os.Getenv("REDIS_ADDRESS"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
			Stream:   os.Getenv("REDIS_STREAM"),
		},
		ShopServiceUrl:    os.Getenv("SHOP_SERVICE_URL"),
		OrdersServiceUrl:  os.Getenv("ORDERS_SERViCE_URL"),
		PaymentServiceUrl: os.Getenv("PAYMENT_SERVICE_URL"),
		BillingServiceUrl: os.Getenv("BILLING_SERViCE_URL"),
	}, nil
}

func GetRequiredVariables() []string {
	return []string{
		// Обязательные переменные окружения для подключения к БД сервиса
		"DB_ADDRESS", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",

		// Обязательные переменные окружения для подключения к redis
		"REDIS_ADDRESS", "REDIS_PORT", "REDIS_PASSWORD", "REDIS_DB_NUMBER", "REDIS_STREAM",

		// Обязательные переменные для обновления крон задач
		"UPDATE_INTERVAL",

		// Url сервиса магазина
		"SHOP_SERVICE_URL",

		// Url сервиса заказов
		"ORDERS_SERViCE_URL",

		// Url сервиса платежей
		"PAYMENT_SERVICE_URL",

		// Url сервиса счетов
		"BILLING_SERViCE_URL",
	}
}
