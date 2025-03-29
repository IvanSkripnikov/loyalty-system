package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"loyalty-system/models"

	"github.com/IvanSkripnikov/go-gormdb"
	"github.com/IvanSkripnikov/go-logger"
	"gorm.io/gorm"
)

func GetLoyaltyList(w http.ResponseWriter, _ *http.Request) {
	category := "/v1/loyalty/list"
	var loyalty []models.Loyalty

	db := gormdb.GetClient(models.ServiceDatabase)
	err := db.Find(&loyalty).Error
	if checkError(w, err, category) {
		return
	}

	data := ResponseData{
		"response": loyalty,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func GetLoyalty(w http.ResponseWriter, r *http.Request) {
	category := "/v1/loyalty/get"
	var loyalty models.Loyalty

	loyaltyID, err := getIDFromRequestString(strings.TrimSpace(r.URL.Path))
	if checkError(w, err, category) {
		return
	}

	db := gormdb.GetClient(models.ServiceDatabase)
	err = db.Where("id = ?", loyaltyID).First(&loyalty).Error
	if checkError(w, err, category) {
		return
	}

	data := ResponseData{
		"response": loyalty,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func GetLoyaltyForUser(w http.ResponseWriter, r *http.Request) {
	category := "/v1/loyalty/get-for-user"

	userID, err := getIDFromRequestString(strings.TrimSpace(r.URL.Path))
	if checkError(w, err, category) {
		return
	}

	var loyaltyUserList []models.LoyaltyUser
	err = GormDB.Where("user_id = ? AND active = ?", userID, 1).Find(&loyaltyUserList).Error
	if checkError(w, err, category) {
		return
	}

	var loyaltyIds []int
	for _, loyaltyUser := range loyaltyUserList {
		loyaltyIds = append(loyaltyIds, loyaltyUser.LoyaltyID)
	}

	var loyalty []models.Loyalty
	err = GormDB.Where("loyalty_id IN ?", loyaltyIds).Find(&loyalty).Error
	if checkError(w, err, category) {
		return
	}

	data := ResponseData{
		"response": loyalty,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func CreateLoyalty(w http.ResponseWriter, r *http.Request) {
	category := "/v1/loyalty/create"
	var loyalty models.Loyalty

	err := json.NewDecoder(r.Body).Decode(&loyalty)
	if checkError(w, err, category+":json_decode") {
		return
	}

	loyalty.Created = GetCurrentDate()
	loyalty.Expired = GetCurrentDate()
	loyalty.Active = 1

	err = GormDB.Create(&loyalty).Error
	if checkError(w, err, category+":create") {
		return
	}

	data := ResponseData{
		"response": models.Success,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func UpdateLoyalty(w http.ResponseWriter, r *http.Request) {
	category := "/v1/loyalty/update"
	var loyaltyRequest models.Loyalty

	err := json.NewDecoder(r.Body).Decode(&loyaltyRequest)
	if checkError(w, err, category) {
		return
	}

	var loyalty models.Loyalty
	err = GormDB.Where("id = ?", loyaltyRequest.ID).First(&loyalty).Error
	if checkError(w, err, category) {
		return
	}

	err = GormDB.Model(&loyalty).Updates(models.Loyalty{
		Title:   loyaltyRequest.Title,
		Expired: loyaltyRequest.Expired,
		Active:  loyaltyRequest.Active,
		Data:    loyaltyRequest.Data,
	}).Error
	if checkError(w, err, category) {
		return
	}

	data := ResponseData{
		"response": models.Success,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func DeleteLoyalty(w http.ResponseWriter, r *http.Request) {
	category := "/v1/loyalty/remove"

	loyaltyID, err := getIDFromRequestString(strings.TrimSpace(r.URL.Path))
	if checkError(w, err, category) {
		return
	}

	err = GormDB.Delete(&models.Loyalty{}, loyaltyID).Error
	if checkError(w, err, category) {
		return
	}
	data := ResponseData{
		"response": models.Success,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func ApplyForOrder(w http.ResponseWriter, r *http.Request) {
	category := "/v1/loyalty/apply-for-order"

	var order models.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	checkError(w, err, category)

	// здесь произвести рассчёт всех скидок и предолжений
	order = RecalculateForOrder(order)

	data := ResponseData{
		"response": order,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func RecalculateForOrder(order models.Order) models.Order {
	var err error
	var loyaltyUserList []models.LoyaltyUser
	err = GormDB.Where("user_id = ? AND active = ?", order.UserID, 1).Find(&loyaltyUserList).Error
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorf("Error getting loyalty for user: %v", err)
		}
		return order
	}

	var loyaltyIds []int
	for _, loyaltyUser := range loyaltyUserList {
		loyaltyIds = append(loyaltyIds, loyaltyUser.LoyaltyID)
	}

	var loyalty []models.Loyalty
	err = GormDB.Where("loyalty_id IN ?", loyaltyIds).Find(&loyalty).Error
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorf("Error getting loyalty for user: %v", err)
		}
		return order
	}

	hasPromocode := false
	hasFirstBuyDiscount := false
	hasRegularDiscount := false
	hasCertificate := false
	hasTemporaryDiscount := false
	var discountLoyalty models.Loyalty
	for _, loyaltyItem := range loyalty {
		if loyaltyItem.TypeID == models.LoyaltyTypePromocode && order.Promocode == loyaltyItem.Title {
			hasPromocode = true
		}
		if loyaltyItem.TypeID == models.LoyaltyTypeNoOrders {
			hasFirstBuyDiscount = true
		}
		if loyaltyItem.TypeID == models.LoyaltyTypeDiscount1 || loyaltyItem.TypeID == models.LoyaltyTypeDiscount2 || loyaltyItem.TypeID == models.LoyaltyTypeDiscount3 || loyaltyItem.TypeID == models.LoyaltyTypeDiscount4 {
			hasRegularDiscount = true
			discountLoyalty = loyaltyItem
		}
		if loyaltyItem.TypeID == models.LoyaltyTypeCertificate && order.Certificate == loyaltyItem.Title {
			hasCertificate = true
		}
		if loyaltyItem.TypeID == models.LoyaltyTypeTempDiscount {
			hasTemporaryDiscount = true
		}
	}

	if hasPromocode {
		promocodeLoyalty := getLoyaltyFromListByType(models.LoyaltyTypePromocode, loyalty)
		var promocode models.Promocode

		err := json.Unmarshal([]byte(promocodeLoyalty.Data), &promocode)
		if err != nil {
			logger.Errorf("Cant parse promocode value: %v", err)
		}

		if hasCertificate {
			order.Price = getPriceOfAppliedCertificate(order.Price, loyalty)
		}

		if promocode.Type == models.PromocodeTypeStatic {
			order.Price = order.Price - float32(promocode.Value)
		} else {
			percent := (order.Price / 100) * float32(promocode.Value)
			order.Price = order.Price - percent
		}
	} else if hasFirstBuyDiscount {
		firstBuyDiscountLoyalty := getLoyaltyFromListByType(models.LoyaltyTypeNoOrders, loyalty)
		var discount models.FirstDiscount
		err := json.Unmarshal([]byte(firstBuyDiscountLoyalty.Data), &discount)
		if err != nil {
			logger.Errorf("Cant parse first discount value: %v", err)
		}

		if hasCertificate {
			order.Price = getPriceOfAppliedCertificate(order.Price, loyalty)
		}

		if discount.Type == models.PromocodeTypeStatic {
			order.Price = order.Price - float32(discount.Value)
		} else {
			percent := (order.Price / 100) * float32(discount.Value)
			order.Price = order.Price - percent
		}
	} else if hasTemporaryDiscount {
		tempDiscountLoyalty := getLoyaltyFromListByType(models.LoyaltyTypeTempDiscount, loyalty)
		var discount models.TempDiscount
		err := json.Unmarshal([]byte(tempDiscountLoyalty.Data), &discount)
		if err != nil {
			logger.Errorf("Cant parse first discount value: %v", err)
		}

		t := time.Now()
		now := t.Format("2006-01-02")
		if now >= discount.FromDate && now <= discount.ToDate {
			if hasCertificate {
				order.Price = getPriceOfAppliedCertificate(order.Price, loyalty)
			}

			if discount.Type == models.PromocodeTypeStatic {
				order.Price = order.Price - float32(discount.Value)
			} else {
				percent := (order.Price / 100) * float32(discount.Value)
				order.Price = order.Price - percent
			}
		}

	} else if hasRegularDiscount {
		var discount models.FirstDiscount
		err := json.Unmarshal([]byte(discountLoyalty.Data), &discount)
		if err != nil {
			logger.Errorf("Cant parse first discount value: %v", err)
		}

		if hasCertificate {
			order.Price = getPriceOfAppliedCertificate(order.Price, loyalty)
		}

		if discount.Type == models.PromocodeTypeStatic {
			order.Price = order.Price - float32(discount.Value)
		} else {
			percent := (order.Price / 100) * float32(discount.Value)
			order.Price = order.Price - percent
		}
	} else if hasCertificate {
		order.Price = getPriceOfAppliedCertificate(order.Price, loyalty)
	}

	return order
}

func getLoyaltyFromListByType(typeID int, loyalty []models.Loyalty) models.Loyalty {
	var defaultValue models.Loyalty
	for _, loyaltyItem := range loyalty {
		if loyaltyItem.TypeID == typeID {
			defaultValue = loyaltyItem
			break
		}
	}

	return defaultValue
}

func getPriceOfAppliedCertificate(price float32, loyalty []models.Loyalty) float32 {
	certificateLoyalty := getLoyaltyFromListByType(models.LoyaltyTypeCertificate, loyalty)
	var certificate models.Certificate
	err := json.Unmarshal([]byte(certificateLoyalty.Data), &certificate)
	if err != nil {
		logger.Errorf("Cant parse certificate value: %v", err)
	}

	return price - float32(certificate.Value)
}
