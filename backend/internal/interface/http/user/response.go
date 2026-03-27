package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Envelope struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	RequestID string `json:"request_id"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Code: 0, Message: "ok", Data: data, RequestID: requestID()})
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Code: status, Message: message, Data: map[string]any{}, RequestID: requestID()})
}

func requestID() string {
	return "req_" + strconv.FormatInt(time.Now().UnixNano(), 10)
}
