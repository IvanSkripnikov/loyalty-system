package helpers

import (
	"net/http"
	"strings"

	"loyalty-system/models"

	"github.com/IvanSkripnikov/go-gormdb"
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
	var loyalty []models.Loyalty

	userID, err := getIDFromRequestString(strings.TrimSpace(r.URL.Path))
	if checkError(w, err, category) {
		return
	}

	db := gormdb.GetClient(models.ServiceDatabase)
	err = db.Where("user_id = ?", userID).Find(&loyalty).Error
	if checkError(w, err, category) {
		return
	}

	data := ResponseData{
		"response": loyalty,
	}
	SendResponse(w, data, category, http.StatusOK)
}
