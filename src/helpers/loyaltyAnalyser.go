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

	// получить все новые промокоды для пользователей
	var newPromocodes []models.Loyalty
	err = GormDB.Where("type_id = ? AND active = ?", models.LoyaltyTypePromocode, 1).Find(&newPromocodes).Error
	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Fatalf("Cant get new promocodes list: %v", err)
	}

	// получить текущие акции для пользователей
	var newTempDiscounts []models.Loyalty
	err = GormDB.Where("type_id = ? AND active = ?", models.LoyaltyTypeTempDiscount, 1).Find(&newTempDiscounts).Error
	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Fatalf("Cant get new temp discounts list: %v", err)
	}

	// начинаем просмотр всех пользователей
	for _, user := range users {
		logger.Debugf("User info: %v", user)

		// 1. выставление скидок по анализу всех покупок
		// получить все заказы пользователя
		var orders []models.Order
		var commonPrice float32
		response, err = CreateQueryWithScalarResponse(http.MethodGet, Config.OrdersServiceUrl+"/v1/orders/get-by-user/"+strconv.Itoa(user.ID), nil)
		if err != nil {
			logger.Fatalf("Cant get orders list: %v", err)
		}
		orders = response.([]models.Order)

		if len(orders) != 0 {
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
		// TODO реализовать в сервисе магазина получение категории пользователя
		response, err = CreateQueryWithScalarResponse(http.MethodGet, Config.ShopServiceUrl+"/v1/user-category/get-by-user"+strconv.Itoa(user.ID), nil)
		if err != nil {
			logger.Errorf("Cant get user category: %v", err)
		}
		category := response.(models.UserCategory)
		if category.CategoryID == models.UserCategoryVIP {
			continue
		} else {
			CheckForVIPCategory(user.ID)
		}
	}
}

// производятся проверки на то, можно ли перевести пользователя в статус VIP
func CheckForVIPCategory(userID int) {
	var deposits []models.Payment
	var response any

	vipAmount, err := strconv.ParseFloat(ConfigMap[models.TriggerSwitchVIPUserCategory], 32)
	if err != nil {
		logger.Errorf("Cant get config value TriggerSwitchVIPUserCategory: %v", err)
	}
	// TODO реализовать в сервисе платежей получение списка депозитов
	response, err = CreateQueryWithScalarResponse(http.MethodGet, Config.PaymentServiceUrl+"/v1/payment/get-deposits-by-user/"+strconv.Itoa(userID), nil)
	if err != nil {
		logger.Fatalf("Cant get deposits list: %v", err)
	}
	deposits = response.([]models.Payment)

	for _, deposit := range deposits {
		if deposit.Amount >= float32(vipAmount) {
			// TODO реализовать в сервисе магазина смену категории пользователю
			_, err = CreateQueryWithScalarResponse(http.MethodPut, Config.ShopServiceUrl+"/v1/user-category/update", nil)
			if err != nil {
				logger.Fatalf("Cant change user category: %v", err)
			}
			break
		}
	}

	// 5. поменять группу пользователя анализируя счета
	response, err = CreateQueryWithScalarResponse(http.MethodGet, Config.BillingServiceUrl+"/v1/account/get-balance/"+strconv.Itoa(userID), nil)
	if err != nil {
		logger.Errorf("Cant get account balance: %v", err)
	}
	balance, err := strconv.ParseFloat(response.(string), 32)
	if err != nil {
		logger.Errorf("Cant parse account balance: %v", err)
	}
	if balance >= vipAmount {
		// TODO реализовать в сервисе магазина смену категории пользователю
		_, err = CreateQueryWithScalarResponse(http.MethodPut, Config.ShopServiceUrl+"/v1/user-category/update", nil)
		if err != nil {
			logger.Fatalf("Cant change user category: %v", err)
		}
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

func RemoveLoyalty(userID, loyaltyID int) {
	err := GormDB.Model(models.LoyaltyUser{}).Where("user_id = ? AND loyalty_id = ?", userID, loyaltyID).Updates(models.LoyaltyUser{Active: 0})
	if err != nil {
		logger.Errorf("Cant delete loyalty for user: %v", err)
	}
}
