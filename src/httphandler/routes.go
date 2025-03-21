package httphandler

import (
	"net/http"
	"regexp"

	"loyalty-system/controllers"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

var routes = []route{
	// system
	newRoute(http.MethodGet, "/health", controllers.HealthCheck),
	// notifications
	newRoute(http.MethodGet, "/v1/loyalty/list", controllers.GetLoyaltyListV1),
	newRoute(http.MethodGet, "/v1/loyalty/get/([0-9]+)", controllers.GetLoyaltyV1),
	newRoute(http.MethodGet, "/v1/loyalty/get-for-user/([0-9]+)", controllers.GetLoyaltyForUserV1),
}
