package admin

import (
	"encoding/json"
	"errors"
	"net/http"

	app "listen/backend/internal/application/admin_auth"
	domain "listen/backend/internal/domain/admin_auth"
)

type AuthController struct {
	loginMockUC app.LoginMockUseCase
	getMeUC     app.GetCurrentAdminUseCase
}

func NewAuthController(loginMockUC app.LoginMockUseCase, getMeUC app.GetCurrentAdminUseCase) AuthController {
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
		writeAdminAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, LoginMockResponseDTO{
		AccessToken: output.Identity.AccessToken,
		AdminID:     output.Identity.AdminID,
		Role:        output.Identity.Role,
	})
}

func (c AuthController) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	output, err := c.getMeUC.Execute(r.Context(), app.GetCurrentAdminInput{
		AdminID: currentAdminID(r),
	})
	if err != nil {
		writeAdminAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, AdminMeResponseDTO{
		AdminID:     output.Admin.AdminID,
		Account:     output.Admin.Account,
		Role:        output.Admin.Role,
		DisplayName: output.Admin.DisplayName,
		Status:      output.Admin.Status,
	})
}

func writeAdminAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidCredential):
		writeJSONError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		writeJSONError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrAdminNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
