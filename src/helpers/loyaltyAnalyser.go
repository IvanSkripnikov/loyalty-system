package helpers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"loyalty-system/models"

	"github.com/IvanSkripnikov/go-logger"
	"gorm.io/gorm"
)

// здесь производится проверка на необхоимость присовения пользователям новой лояльности
func ApplyLoyalty() {
	// Получаем всех активных пользователей
	var users []models.User
	var response any
	var err error

	// TODO реализовать в сервисе магазина получение списка ативных пользователей
	response, err = CreateQueryWithScalarResponse(http.MethodGet, Config.ShopServiceUrl+"/v1/users/get-active", nil)
	if err != nil {
		logger.Fatalf("Cant get users list: %v", err)
	}
	users = response.([]models.User)

	// начинаем просмотр всех пользователей
	for _, user := range users {
		logger.Debugf("User info: %v", user)

		// получить все заказы пользователя
		var orders []models.Order
		var commonPrice float32
		var noOrdersLoyalty models.Loyalty
		response, err = CreateQueryWithScalarResponse(http.MethodGet, Config.OrdersServiceUrl+"/v1/orders/get-by-user/"+strconv.Itoa(user.ID), nil)
		if err != nil {
			logger.Fatalf("Cant get users list: %v", err)
		}
		orders = response.([]models.Order)

		if len(orders) != 0 {
			// если пользователь не совершал заказов - делаем скидку на первый заказ
			err := GormDB.Where("type_id = ?", models.LoyaltyTypeNoOrders).First(&noOrdersLoyalty).Error
			if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
				SetLoyalty(user.ID, noOrdersLoyalty.ID)
			}
		} else {
			for _, order := range orders {
				logger.Debugf("Order info: %v", order)
				commonPrice = commonPrice + order.Price
			}

			firstLevelSum, err := strconv.ParseFloat(ConfigMap[models.TriggerFirstLevelOrdersSum], 32)
			if err != nil {
				logger.Errorf("Cant get config value: %v", err)
				continue
			}
			if commonPrice > float32(firstLevelSum) {
				// TODO продолжить реализацию
			}
		}

		// получить все платежи пользователя
	}
}

// здесь деактивируются истёкшие лояльности
func CheckExpiredLoyalty() {
	t := time.Now()
	now := t.Format("2006-01-02")
	GormDB.Model(models.Loyalty{}).Where("expired < ?", now).Updates(models.Loyalty{Active: 0})
}

// проверка, есть ли у пользователя данная лояльность
func isExistsLoyalty(userID, loyaltyID int) bool {
	var relation models.LoyaltyUser
	err := GormDB.Where("user_id = ? AND loyalty_id = ?", userID, loyaltyID).First(&relation).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	if err != nil {
		return true
	}

	return true
}

func SetLoyalty(userID, loyaltyID int) {
	// проверяем, есть ли уже такая лояльность у пользователя
	if isExistsLoyalty(userID, loyaltyID) {
		return
	}

	// создаём лояльность для пользователя
	var userLoyalty models.LoyaltyUser
	userLoyalty.UserID = userID
	userLoyalty.LoyaltyID = loyaltyID
	err := GormDB.Create(&userLoyalty).Error
	if err != nil {
		logger.Errorf("Cant create new loyalty %v", err)
	}

	// уведомляем пользователя
	SendNewLoyaltyNotification(userID, loyaltyID)
}

func SendNewLoyaltyNotification(userID, loyaltyID int) {
	var loyalty models.Loyalty
	var err error
	err = GormDB.Where("id = ?", loyaltyID).First(&loyalty).Error
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		if err != nil {
			logger.Errorf("Error during get loyalty from database: %v", err)
		}
		return
	}
	var loyaltyType models.LoyaltyType
	err = GormDB.Where("id = ?", loyalty.TypeID).First(&loyaltyType).Error
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		if err != nil {
			logger.Errorf("Error during get loyalty type from database: %v", err)
		}
		return
	}
	messageData := map[string]interface{}{
		"title":       "Congratulation! You get new privelegy: " + loyalty.Title,
		"description": "Since this time you have access to: " + loyaltyType.Description,
		"user":        userID,
		"category":    "loyalty",
	}
	SendNotification(messageData)
}
