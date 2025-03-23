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
	// loyalty
	newRoute(http.MethodGet, "/v1/loyalty/list", controllers.GetLoyaltyListV1),
	newRoute(http.MethodGet, "/v1/loyalty/get/([0-9]+)", controllers.GetLoyaltyV1),
	newRoute(http.MethodGet, "/v1/loyalty/get-for-user/([0-9]+)", controllers.GetLoyaltyForUserV1),
	newRoute(http.MethodPut, "/v1/loyalty/apply-for-order", controllers.ApplyForOrderV1),
}
