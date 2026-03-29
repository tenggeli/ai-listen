package user

import "net/http"

func RegisterAIRoutes(mux *http.ServeMux, controller AIController) {
	mux.HandleFunc("GET /api/v1/ai/home", controller.HandleGetHome)
	mux.HandleFunc("GET /api/v1/ai/match/remaining", controller.HandleGetRemaining)
	mux.HandleFunc("POST /api/v1/ai/match", controller.HandleMatch)
	mux.HandleFunc("POST /api/v1/ai/sessions", controller.HandleCreateSession)

	mux.HandleFunc("GET /api/v1/ai/sessions/", controller.HandleSessionDetail)
	mux.HandleFunc("POST /api/v1/ai/sessions/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) < len("/messages") || r.URL.Path[len(r.URL.Path)-9:] != "/messages" {
			writeJSONError(w, http.StatusNotFound, "route not found")
			return
		}
		controller.HandleAppendMessage(w, r)
	})
}

func RegisterIdentityRoutes(mux *http.ServeMux, controller IdentityController) {
	mux.HandleFunc("POST /api/v1/auth/login/sms", controller.HandleSMSLogin)
	mux.HandleFunc("POST /api/v1/auth/login/wechat/mock", controller.HandleWechatMockLogin)
	mux.HandleFunc("GET /api/v1/users/me", controller.HandleGetMe)
	mux.HandleFunc("PUT /api/v1/users/me/profile", controller.HandleSaveProfile)
	mux.HandleFunc("PUT /api/v1/users/me/personality", controller.HandleSavePersonality)
	mux.HandleFunc("POST /api/v1/users/me/personality/skip", controller.HandleSkipPersonality)
}

func RegisterServiceDiscoveryRoutes(mux *http.ServeMux, controller ServiceDiscoveryController) {
	mux.HandleFunc("GET /api/v1/services/categories", controller.HandleListCategories)
	mux.HandleFunc("GET /api/v1/providers/public", controller.HandleListPublicProviders)
	mux.HandleFunc("GET /api/v1/providers/public/{id}", controller.HandleGetPublicProvider)
	mux.HandleFunc("GET /api/v1/providers/public/{id}/service-items", controller.HandleListProviderServiceItems)
}
