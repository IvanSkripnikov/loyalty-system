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
	newRoute(http.MethodPost, "/v1/loyalty/create", controllers.CreateLoyaltyV1),
	newRoute(http.MethodPut, "/v1/loyalty/update", controllers.UpdateLoyaltyV1),
	newRoute(http.MethodDelete, "/v1/loyalty/remove/([0-9]+)", controllers.DeleteLoyaltyV1),
	newRoute(http.MethodDelete, "/v1/loyalty/remove-for-user/([0-9]+)", controllers.DeleteLoyaltyForUserV1),
	newRoute(http.MethodDelete, "/v1/loyalty/remove-certificate", controllers.DeleteCertificateV1),
	newRoute(http.MethodGet, "/v1/loyalty/configuration/list", controllers.GetLoyaltyConfigurationListV1),
	newRoute(http.MethodPut, "/v1/loyalty/configuration/update", controllers.UpdateLoyaltyConfigurationV1),
	// test
	newRoute(http.MethodGet, "/test/run-loyalty", controllers.TestRunLoyalty),
	newRoute(http.MethodGet, "/test/remove-loyalty", controllers.TestRemoveLoyalty),
}
