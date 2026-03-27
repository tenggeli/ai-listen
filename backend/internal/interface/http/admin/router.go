package admin

import "net/http"

func RegisterProviderRoutes(mux *http.ServeMux, controller ProviderController) {
	mux.HandleFunc("GET /api/v1/admin/providers", controller.HandleList)
	mux.HandleFunc("GET /api/v1/admin/providers/{id}", controller.HandleDetail)
	mux.HandleFunc("POST /api/v1/admin/providers/{id}/approve", controller.HandleApprove)
	mux.HandleFunc("POST /api/v1/admin/providers/{id}/reject", controller.HandleReject)
	mux.HandleFunc("POST /api/v1/admin/providers/{id}/require-supplement", controller.HandleRequireSupplement)
}
