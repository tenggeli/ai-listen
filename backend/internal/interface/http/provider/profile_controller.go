package provider

import (
	"encoding/json"
	"errors"
	"net/http"

	providerApp "listen/backend/internal/application/provider"
	providerAuthDomain "listen/backend/internal/domain/provider_auth"
	serviceDiscoveryDomain "listen/backend/internal/domain/service_discovery"
)

type ProfileController struct {
	updateProfileUC providerApp.UpdateCurrentProfileUseCase
	listServicesUC  providerApp.ListCurrentProviderServicesUseCase
}

func NewProfileController(updateProfileUC providerApp.UpdateCurrentProfileUseCase, listServicesUC providerApp.ListCurrentProviderServicesUseCase) ProfileController {
	return ProfileController{updateProfileUC: updateProfileUC, listServicesUC: listServicesUC}
}

func (c ProfileController) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	var body ProviderProfileUpdateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.updateProfileUC.Execute(r.Context(), providerApp.UpdateCurrentProfileInput{
		ProviderID:  currentProviderID(r),
		DisplayName: body.DisplayName,
		CityCode:    body.CityCode,
	})
	if err != nil {
		writeProviderProfileError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, ProviderProfileResponseDTO{
		ProviderID:  output.Provider.ProviderID,
		Account:     output.Provider.Account,
		DisplayName: output.Provider.DisplayName,
		Status:      output.Provider.Status,
		CityCode:    output.Provider.CityCode,
	})
}

func (c ProfileController) HandleListServices(w http.ResponseWriter, r *http.Request) {
	output, err := c.listServicesUC.Execute(r.Context(), providerApp.ListCurrentProviderServicesInput{ProviderID: currentProviderID(r)})
	if err != nil {
		writeProviderProfileError(w, err)
		return
	}

	items := make([]ProviderServiceItemDTO, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, ProviderServiceItemDTO{
			ItemID:        item.ID,
			ProviderID:    item.ProviderID,
			CategoryID:    item.CategoryID,
			Title:         item.Title,
			Description:   item.Description,
			PriceAmount:   item.PriceAmount,
			PriceUnit:     item.PriceUnit,
			SupportOnline: item.SupportOnline,
			SortOrder:     item.SortOrder,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func writeProviderProfileError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, providerApp.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, providerApp.ErrUnauthorized):
		writeJSONError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, providerAuthDomain.ErrProviderNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, serviceDiscoveryDomain.ErrProviderNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
