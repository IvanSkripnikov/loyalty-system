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

func CreateLoyaltyV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		helpers.CreateLoyalty(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/create")
	}
}

func UpdateLoyaltyV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		helpers.UpdateLoyalty(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/update")
	}
}

func DeleteLoyaltyV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		helpers.DeleteLoyalty(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/remove")
	}
}

func GetLoyaltyConfigurationListV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		helpers.GetLoyaltyConfigurationList(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/configuration/list")
	}
}

func UpdateLoyaltyConfigurationV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		helpers.UpdateLoyaltyConfiguration(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/configuration/update")
	}
}

func DeleteLoyaltyForUserV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		helpers.DeleteLoyaltyForUser(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/remove-for-user")
	}
}

func DeleteCertificateV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		helpers.DeleteCertificate(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/loyalty/remove-certificate")
	}
}
