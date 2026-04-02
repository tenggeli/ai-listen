package admin

import "net/http"

func RegisterAuthRoutes(mux *http.ServeMux, controller AuthController) {
	mux.HandleFunc("POST /api/v1/admin/auth/login/mock", controller.HandleLoginMock)
	mux.HandleFunc("GET /api/v1/admin/auth/me", requireAdminAuth(controller.HandleGetMe))
}

func RegisterProviderRoutes(mux *http.ServeMux, controller ProviderController) {
	mux.HandleFunc("GET /api/v1/admin/providers", requireAdminAuth(controller.HandleList))
	mux.HandleFunc("GET /api/v1/admin/providers/{id}", requireAdminAuth(controller.HandleDetail))
	mux.HandleFunc("POST /api/v1/admin/providers/{id}/approve", requireAdminAuth(controller.HandleApprove))
	mux.HandleFunc("POST /api/v1/admin/providers/{id}/reject", requireAdminAuth(controller.HandleReject))
	mux.HandleFunc("POST /api/v1/admin/providers/{id}/require-supplement", requireAdminAuth(controller.HandleRequireSupplement))
}

func RegisterServiceItemRoutes(mux *http.ServeMux, controller ServiceItemController) {
	mux.HandleFunc("GET /api/v1/admin/service-items", requireAdminAuth(controller.HandleList))
	mux.HandleFunc("GET /api/v1/admin/service-items/{id}", requireAdminAuth(controller.HandleDetail))
	mux.HandleFunc("POST /api/v1/admin/service-items/{id}/activate", requireAdminAuth(controller.HandleActivate))
	mux.HandleFunc("POST /api/v1/admin/service-items/{id}/deactivate", requireAdminAuth(controller.HandleDeactivate))
}
