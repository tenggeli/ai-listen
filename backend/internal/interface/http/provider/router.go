package provider

import "net/http"

func RegisterAuthRoutes(mux *http.ServeMux, controller AuthController) {
	mux.HandleFunc("POST /api/v1/provider/auth/login/mock", controller.HandleLoginMock)
	mux.HandleFunc("GET /api/v1/provider/profile", requireProviderAuth(controller.HandleGetMe))
}

func RegisterProfileRoutes(mux *http.ServeMux, controller ProfileController) {
	mux.HandleFunc("PUT /api/v1/provider/profile", requireProviderAuth(controller.HandleUpdateProfile))
	mux.HandleFunc("GET /api/v1/provider/services", requireProviderAuth(controller.HandleListServices))
}

func RegisterOrderRoutes(mux *http.ServeMux, controller OrderController) {
	mux.HandleFunc("GET /api/v1/provider/orders", requireProviderAuth(controller.HandleListOrders))
	mux.HandleFunc("GET /api/v1/provider/orders/{id}", requireProviderAuth(controller.HandleGetOrder))
	mux.HandleFunc("POST /api/v1/provider/orders/{id}/accept", requireProviderAuth(controller.HandleAccept))
	mux.HandleFunc("POST /api/v1/provider/orders/{id}/depart", requireProviderAuth(controller.HandleDepart))
	mux.HandleFunc("POST /api/v1/provider/orders/{id}/arrive", requireProviderAuth(controller.HandleArrive))
	mux.HandleFunc("POST /api/v1/provider/orders/{id}/start", requireProviderAuth(controller.HandleStart))
	mux.HandleFunc("POST /api/v1/provider/orders/{id}/complete", requireProviderAuth(controller.HandleComplete))
}
