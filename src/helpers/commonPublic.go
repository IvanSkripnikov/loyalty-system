package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"loyalty-system/models"
	"net/http"
	"time"

	"github.com/IvanSkripnikov/go-logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var Config *models.Config
var ConfigMap map[string]string

func InitConfig(cfg *models.Config) {
	Config = cfg
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func FormatResponse(w http.ResponseWriter, httpStatus int, category string) {
	w.WriteHeader(httpStatus)

	data := ResponseData{
		"error": "Unsuccessfull request",
	}
	SendResponse(w, data, category, httpStatus)
}

func CreateQueryWithResponse(method, url string, data any) (any, error) {
	var err error
	var response any

	jsonData, err := json.Marshal(data)
	if err != nil {
		return response, err
	}
	logger.Infof("json data: %v", string(jsonData))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return response, err
	}

	resp, err := client.Do(req)
	logger.Infof("response for request %v: %v", url, resp)
	if err != nil {
		return response, err
	}

	var result map[string]any
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &result)

	logger.Infof("Data from response %v", result)

	// Преобразуем JSON-строку в map
	if err != nil {
		return response, err
	}

	response, ok := result["response"]
	if !ok {
		return "", errors.New("failed to get response")
	}

	return response, nil
}

func SendNotification(message map[string]interface{}) {
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf("Error connection to Redis: %v", err)
	}

	_, err = redisClient.XAdd(context.Background(), &redis.XAddArgs{
		Stream: Config.Redis.Stream,
		Values: message,
	}).Result()
	if err != nil {
		logger.Fatalf("Error sending to redis stream: %v", err)
	} else {
		logger.Info("Succsessfuly send to stream")
	}
}

func LoadLoyaltyConfig() {
	var loyaltyConfiguration []models.LoyaltyConfiguration
	err := GormDB.Find(&loyaltyConfiguration).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Fatalf("Cant load loyalty configuration: %v", err)
		return
	}

	ConfigMap = make(map[string]string)
	for _, cfg := range loyaltyConfiguration {
		ConfigMap[cfg.Property] = cfg.Value
	}
}
