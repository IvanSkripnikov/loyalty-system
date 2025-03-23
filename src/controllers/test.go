package controllers

import (
	"net/http"

	"loyalty-system/helpers"
)

func TestRunLoyalty(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		helpers.ApplyLoyalty()
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/test/run-loyalty")
	}
}

func TestRemoveLoyalty(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		helpers.CheckExpiredLoyalty()
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/test/remove-loyalty")
	}
}
