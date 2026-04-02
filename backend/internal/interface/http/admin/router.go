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

func RegisterOrderRoutes(mux *http.ServeMux, controller OrderController) {
	mux.HandleFunc("GET /api/v1/admin/orders", requireAdminAuth(controller.HandleListOrders))
	mux.HandleFunc("GET /api/v1/admin/orders/{id}", requireAdminAuth(controller.HandleGetOrderDetail))
	mux.HandleFunc("POST /api/v1/admin/orders/{id}/intervene", requireAdminAuth(controller.HandleInterveneOrder))
	mux.HandleFunc("POST /api/v1/admin/orders/{id}/close", requireAdminAuth(controller.HandleCloseOrder))
}

func RegisterComplaintRoutes(mux *http.ServeMux, controller OrderController) {
	mux.HandleFunc("GET /api/v1/admin/complaints", requireAdminAuth(controller.HandleListComplaints))
	mux.HandleFunc("GET /api/v1/admin/complaints/{id}", requireAdminAuth(controller.HandleGetComplaintDetail))
	mux.HandleFunc("POST /api/v1/admin/complaints/{id}/intervene", requireAdminAuth(controller.HandleInterveneComplaint))
	mux.HandleFunc("POST /api/v1/admin/complaints/{id}/resolve", requireAdminAuth(controller.HandleResolveComplaint))
}
