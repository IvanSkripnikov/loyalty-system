package controllers

import (
	"net/http"

	"loyalty-system/helpers"
)

func GetLoyaltyListV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		helpers.GetLoyaltyList(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/list")
	}
}

func GetLoyaltyV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		helpers.GetLoyalty(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/get")
	}
}

func GetLoyaltyForUserV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		helpers.GetLoyaltyForUser(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/get-for-user")
	}
}

func ApplyForOrderV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		helpers.ApplyForOrder(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/apply-for-order")
	}
}
