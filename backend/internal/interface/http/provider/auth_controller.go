package provider

import (
	"encoding/json"
	"errors"
	"net/http"

	app "listen/backend/internal/application/provider_auth"
	domain "listen/backend/internal/domain/provider_auth"
)

type AuthController struct {
	loginMockUC app.LoginMockUseCase
	getMeUC     app.GetCurrentProviderUseCase
}

func NewAuthController(loginMockUC app.LoginMockUseCase, getMeUC app.GetCurrentProviderUseCase) AuthController {
	return AuthController{loginMockUC: loginMockUC, getMeUC: getMeUC}
}

func (c AuthController) HandleLoginMock(w http.ResponseWriter, r *http.Request) {
	var body LoginMockRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.loginMockUC.Execute(r.Context(), app.LoginMockInput{
		Account:  body.Account,
		Password: body.Password,
	})
	if err != nil {
		writeProviderAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, LoginMockResponseDTO{
		AccessToken: output.Identity.AccessToken,
		ProviderID:  output.Identity.ProviderID,
	})
}

func (c AuthController) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	output, err := c.getMeUC.Execute(r.Context(), app.GetCurrentProviderInput{ProviderID: currentProviderID(r)})
	if err != nil {
		writeProviderAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, ProviderMeResponseDTO{
		ProviderID:  output.Provider.ProviderID,
		Account:     output.Provider.Account,
		DisplayName: output.Provider.DisplayName,
		Status:      output.Provider.Status,
		CityCode:    output.Provider.CityCode,
	})
}

func writeProviderAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidCredential):
		writeJSONError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		writeJSONError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrProviderNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
