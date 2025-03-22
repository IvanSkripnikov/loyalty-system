package helpers

import (
	"os"
	"strconv"
	"time"

	"github.com/IvanSkripnikov/go-logger"
	"github.com/go-co-op/gocron"
)

func InitTimer() {
	interval, err := strconv.Atoi(os.Getenv("UPDATE_INTERVAL"))
	if err != nil {
		logger.Errorf("Cant get update interval value: %v", err)
	}

	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(interval).Seconds().Do(ApplyLoyalty)
	scheduler.Every(interval).Seconds().Do(CheckExpiredLoyalty)

	scheduler.StartAsync()
}
