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
	// актуализируем конфиг
	LoadLoyaltyConfig()

	// Получаем всех активных пользователей
	var response any
	var err error

	response, err = CreateQueryWithResponse(http.MethodGet, Config.ShopServiceUrl+"/v1/users/get-active", nil)
	if err != nil {
		logger.Errorf("Cant get users list: %v", err)
	}

	users := getUsersListFromResponse(response)

	// получить все новые промокоды для пользователей
	var newPromocodes []models.Loyalty
	err = GormDB.Where("type_id = ? AND active = ?", models.LoyaltyTypePromocode, 1).Find(&newPromocodes).Error
	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Infof("Cant get new promocodes list: %v", err)
	}

	// получить текущие акции для пользователей
	var newTempDiscounts []models.Loyalty
	err = GormDB.Where("type_id = ? AND active = ?", models.LoyaltyTypeTempDiscount, 1).Find(&newTempDiscounts).Error
	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Infof("Cant get new temp discounts list: %v", err)
	}

	// начинаем просмотр всех пользователей
	for _, user := range users {
		logger.Debugf("User info: %v", user)

		// 1. выставление скидок по анализу всех покупок
		// получить все заказы пользователя
		var commonPrice float32
		response, err = CreateQueryWithResponse(http.MethodGet, Config.OrdersServiceUrl+"/v1/orders/get-by-user/"+strconv.Itoa(user.ID), nil)
		if err != nil {
			logger.Infof("Cant get orders list: %v", err)
		}

		orders := getOrdersListFromResponse(response)
		if len(orders) == 0 {
			var noOrdersLoyalty models.Loyalty
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
			// проверка на возможное проставление постоянной скидки или её повышения
			SetDiscountForUser(user.ID, commonPrice)
		}

		// 2. сделать доступными пользователю новые промокоды
		for _, promocode := range newPromocodes {
			SetLoyalty(user.ID, promocode.ID)
		}

		// 3. сделать доступными пользователю временные скидки
		for _, tempDiscount := range newTempDiscounts {
			SetLoyalty(user.ID, tempDiscount.ID)
		}

		// 4. поменять группу пользователя анализируя платежи
		response, err = CreateQueryWithResponse(http.MethodGet, Config.ShopServiceUrl+"/v1/user-category/get-by-user/"+strconv.Itoa(user.ID), nil)
		if err != nil {
			logger.Errorf("Cant get user category: %v", err)
		}
		category := getUserCategoryFromResponse(response)
		if category.ID == models.UserCategoryVIP {
			continue
		} else {
			CheckForVIPCategory(user.ID)
		}
	}
}

// производятся проверки на то, можно ли перевести пользователя в статус VIP
func CheckForVIPCategory(userID int) {
	var response any

	vipAmount, err := strconv.ParseFloat(ConfigMap[models.TriggerSwitchVIPUserCategory], 32)
	if err != nil {
		logger.Errorf("Cant get config value TriggerSwitchVIPUserCategory: %v", err)
	}

	response, err = CreateQueryWithResponse(http.MethodGet, Config.PaymentServiceUrl+"/v1/payment/get-deposits-by-user/"+strconv.Itoa(userID), nil)
	if err != nil {
		logger.Infof("Cant get deposits list: %v", err)
	}
	deposits := getPaymentsListFromResponse(response)

	vipFlag := true
	for _, deposit := range deposits {
		if deposit.Amount >= float32(vipAmount) {
			changeCategoryParams := models.UserCategoryParams{UserID: userID, CategoryID: models.UserCategoryVIP}
			_, err = CreateQueryWithResponse(http.MethodPut, Config.ShopServiceUrl+"/v1/users/category-update", changeCategoryParams)
			if err != nil {
				logger.Errorf("Cant change user category: %v", err)
			} else {
				messageData := map[string]interface{}{
					"title":       "Congratulation! You get new privelegy: VIP",
					"description": "Since this time you have access to VIP items!",
					"user":        userID,
					"category":    "loyalty",
				}
				SendNotification(messageData)
				vipFlag = false
			}
			break
		}
	}

	// 5. поменять группу пользователя анализируя счета
	response, err = CreateQueryWithResponse(http.MethodGet, Config.BillingServiceUrl+"/v1/account/get-balance/"+strconv.Itoa(userID), nil)
	if err != nil {
		logger.Errorf("Cant get account balance: %v", err)
	}
	balance := response.(float64)
	if err != nil {
		logger.Errorf("Cant parse account balance: %v", err)
	}
	if balance >= vipAmount {
		changeCategoryParams := models.UserCategoryParams{UserID: userID, CategoryID: models.UserCategoryVIP}
		_, err = CreateQueryWithResponse(http.MethodPut, Config.ShopServiceUrl+"/v1/users/category-update", changeCategoryParams)
		if err != nil {
			logger.Errorf("Cant change user category: %v", err)
		} else {
			if vipFlag {
				messageData := map[string]interface{}{
					"title":       "Congratulation! You get new privelegy: VIP",
					"description": "Since this time you have access to VIP items!",
					"user":        userID,
					"category":    "loyalty",
				}
				SendNotification(messageData)
			}
		}
	}
}

// здесь деактивируются истёкшие лояльности
func CheckExpiredLoyalty() {
	t := time.Now()
	now := t.Format("2006-01-02")
	var loyalty []models.Loyalty

	err := GormDB.Where("expired < ? AND type_id NOT IN (2, 3, 4, 5)", now).Find(&loyalty).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error getting loyalty for deactivate: %v", err)
	}

	loyaltyIDs := make([]int, 0)
	for _, l := range loyalty {
		loyaltyIDs = append(loyaltyIDs, l.ID)
	}

	GormDB.Model(&models.LoyaltyUser{}).Where("loyalty_id IN ?", loyaltyIDs).Update("active", 0)
	GormDB.Model(&models.Loyalty{}).Where("id IN ?", loyaltyIDs).Update("active", 0)
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

// проанализировать уровень скидки и выставить её пользователю
func SetDiscountForUser(userID int, price float32) {
	// инициализируем объекты для сверки
	firstLevelSum, err := strconv.ParseFloat(ConfigMap[models.TriggerMinimalOrdersSum], 32)
	if err != nil {
		logger.Errorf("Cant get config value 1: %v", err)
		return
	}

	// если покупок не набралось ни на какую скидку - выходим
	if price < float32(firstLevelSum) {
		return
	}

	advancedLevelSum, err := strconv.ParseFloat(ConfigMap[models.TriggerFirstLevelOrdersSum], 32)
	if err != nil {
		logger.Errorf("Cant get config value 2: %v", err)
		return
	}
	profiLevelSum, err := strconv.ParseFloat(ConfigMap[models.TriggerSecondLevelOrdersSum], 32)
	if err != nil {
		logger.Errorf("Cant get config value 3: %v", err)
		return
	}
	lastLevelSum, err := strconv.ParseFloat(ConfigMap[models.TriggerThirdLevelOrdersSum], 32)
	if err != nil {
		logger.Errorf("Cant get config value 4: %v", err)
		return
	}

	var maxDiscountLoyalty models.Loyalty
	err = GormDB.Where("type_id = ?", models.LoyaltyTypeDiscount4).First(&maxDiscountLoyalty).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error during get discount loyalty 4: %v", err)
		return
	}
	var profiDiscountLoyalty models.Loyalty
	err = GormDB.Where("type_id = ?", models.LoyaltyTypeDiscount3).First(&profiDiscountLoyalty).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error during get discount loyalty 3: %v", err)
		return
	}
	var advancedDiscountLoyalty models.Loyalty
	err = GormDB.Where("type_id = ?", models.LoyaltyTypeDiscount2).First(&advancedDiscountLoyalty).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error during get discount loyalty 2: %v", err)
		return
	}
	var firstDiscountLoyalty models.Loyalty
	err = GormDB.Where("type_id = ?", models.LoyaltyTypeDiscount1).First(&firstDiscountLoyalty).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Error during get discount loyalty 4: %v", err)
		return
	}

	// проверяем, достиг ли пользователь последнего уровня скидок
	if price >= float32(lastLevelSum) {
		SetLoyalty(userID, maxDiscountLoyalty.ID)
		RemoveLoyalty(userID, profiDiscountLoyalty.ID)

		return
	}

	// проверяем, достиг ли пользователь предпоследнего уровня скидок
	if price >= float32(profiLevelSum) {
		SetLoyalty(userID, profiDiscountLoyalty.ID)
		RemoveLoyalty(userID, advancedDiscountLoyalty.ID)

		return
	}

	// проверяем, достиг ли пользователь предпоследнего уровня скидок
	if price >= float32(advancedLevelSum) {
		SetLoyalty(userID, advancedDiscountLoyalty.ID)
		RemoveLoyalty(userID, firstDiscountLoyalty.ID)

		return
	}

	// проверяем, достиг ли пользователь предпоследнего уровня скидок
	if price >= float32(firstLevelSum) {
		SetLoyalty(userID, firstDiscountLoyalty.ID)

		return
	}
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
	userLoyalty.Active = 1
	err := GormDB.Create(&userLoyalty).Error
	if err != nil {
		logger.Errorf("Cant create new loyalty %v", err)
	}

	// повышаем счётчик лояльностей
	LoyaltyTotal.Inc()

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

func RemoveLoyalty(userID, loyaltyID int) {
	err := GormDB.Model(models.LoyaltyUser{}).Where("user_id = ? AND loyalty_id = ?", userID, loyaltyID).Updates(models.LoyaltyUser{Active: 0})
	if err != nil {
		logger.Errorf("Cant delete loyalty for user: %v", err)
	}
}

func getUsersListFromResponse(response any) []models.User {
	var users []models.User
	responseArray := response.([]any)
	for _, item := range responseArray {
		userMap, ok := item.(map[string]any)
		if !ok {
			logger.Errorf("Error asserting item to map[string]interface{}")
		}

		// Создаем экземпляр User и заполняем его данными
		user := models.User{
			ID:         int(userMap["id"].(float64)),
			UserName:   userMap["username"].(string),
			Password:   userMap["password"].(string),
			FirstName:  userMap["first_name"].(string),
			LastName:   userMap["last_name"].(string),
			Email:      userMap["email"].(string),
			Phone:      userMap["phone"].(string),
			CategoryID: int(userMap["category_id"].(float64)),
			Created:    userMap["created"].(string),
			Updated:    userMap["updated"].(string),
			Active:     int(userMap["active"].(float64)),
		}

		users = append(users, user)
	}

	return users
}

func getOrdersListFromResponse(response any) []models.Order {
	var orders []models.Order
	responseArray := response.([]any)
	for _, item := range responseArray {
		userMap, ok := item.(map[string]any)
		if !ok {
			logger.Errorf("Error asserting item to map[string]interface{}")
		}

		// Создаем экземпляр Order и заполняем его данными
		order := models.Order{
			ID:          int(userMap["id"].(float64)),
			UserID:      int(userMap["userId"].(float64)),
			ItemID:      int(userMap["itemId"].(float64)),
			Volume:      int(userMap["volume"].(float64)),
			Price:       float32(userMap["price"].(float64)),
			Created:     userMap["created"].(string),
			Updated:     userMap["updated"].(string),
			Status:      int(userMap["id"].(float64)),
			RequestID:   userMap["requestId"].(string),
			Promocode:   "",
			Certificate: "",
		}

		orders = append(orders, order)
	}

	return orders
}

func getPaymentsListFromResponse(response any) []models.Payment {
	var payments []models.Payment
	responseArray := response.([]any)
	for _, item := range responseArray {
		userMap, ok := item.(map[string]any)
		if !ok {
			logger.Errorf("Error asserting item to map[string]interface{}")
		}

		// Создаем экземпляр Payment и заполняем его данными
		payment := models.Payment{
			ID:        int(userMap["id"].(float64)),
			UserID:    int(userMap["userId"].(float64)),
			Type:      userMap["type"].(string),
			Amount:    float32(userMap["amount"].(float64)),
			Created:   userMap["created"].(string),
			Status:    int(userMap["id"].(float64)),
			RequestID: userMap["requestId"].(string),
		}

		payments = append(payments, payment)
	}

	return payments
}

func getUserCategoryFromResponse(response any) models.UserCategory {
	userMap, ok := response.(map[string]any)
	if !ok {
		logger.Errorf("Error asserting item to map[string]interface{}")
	}

	userCategory := models.UserCategory{
		ID:      int(userMap["id"].(float64)),
		Title:   userMap["title"].(string),
		Created: userMap["created"].(string),
		Active:  int(userMap["active"].(float64)),
	}

	return userCategory
}
