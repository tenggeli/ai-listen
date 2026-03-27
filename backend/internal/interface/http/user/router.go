package user

import "net/http"

func RegisterAIRoutes(mux *http.ServeMux, controller AIController) {
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
