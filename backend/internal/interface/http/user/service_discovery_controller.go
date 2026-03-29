package user

import (
	"errors"
	"net/http"
	"strconv"

	app "listen/backend/internal/application/service_discovery"
	domain "listen/backend/internal/domain/service_discovery"
)

type ServiceDiscoveryController struct {
	listCategoriesUC app.ListServiceCategoriesUseCase
	listProvidersUC  app.ListPublicProvidersUseCase
	getProviderUC    app.GetPublicProviderUseCase
	listItemsUC      app.ListProviderServiceItemsUseCase
}

func NewServiceDiscoveryController(
	listCategoriesUC app.ListServiceCategoriesUseCase,
	listProvidersUC app.ListPublicProvidersUseCase,
	getProviderUC app.GetPublicProviderUseCase,
	listItemsUC app.ListProviderServiceItemsUseCase,
) ServiceDiscoveryController {
	return ServiceDiscoveryController{
		listCategoriesUC: listCategoriesUC,
		listProvidersUC:  listProvidersUC,
		getProviderUC:    getProviderUC,
		listItemsUC:      listItemsUC,
	}
}

func (c ServiceDiscoveryController) HandleListCategories(w http.ResponseWriter, r *http.Request) {
	output, err := c.listCategoriesUC.Execute(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	items := make([]ServiceCategoryResponseDTO, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, ServiceCategoryResponseDTO{
			ID:   item.ID,
			Name: item.Name,
			Icon: item.Icon,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (c ServiceDiscoveryController) HandleListPublicProviders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	input := ListPublicProvidersRequestDTO{
		CategoryID: r.URL.Query().Get("category_id"),
		Keyword:    r.URL.Query().Get("keyword"),
		CityCode:   r.URL.Query().Get("city_code"),
		Page:       page,
		PageSize:   pageSize,
	}

	output, err := c.listProvidersUC.Execute(r.Context(), app.ListPublicProvidersInput{
		CategoryID: input.CategoryID,
		Keyword:    input.Keyword,
		CityCode:   input.CityCode,
		Page:       input.Page,
		PageSize:   input.PageSize,
	})
	if err != nil {
		writeServiceDiscoveryError(w, err)
		return
	}

	items := make([]ProviderPublicResponseDTO, 0, len(output.Items))
	for _, provider := range output.Items {
		items = append(items, buildProviderPublicDTO(provider))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": output.Total})
}

func (c ServiceDiscoveryController) HandleGetPublicProvider(w http.ResponseWriter, r *http.Request) {
	providerID := r.PathValue("id")
	if providerID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid provider id")
		return
	}

	output, err := c.getProviderUC.Execute(r.Context(), app.GetPublicProviderInput{ProviderID: providerID})
	if err != nil {
		writeServiceDiscoveryError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildProviderPublicDTO(output.Provider))
}

func (c ServiceDiscoveryController) HandleListProviderServiceItems(w http.ResponseWriter, r *http.Request) {
	providerID := r.PathValue("id")
	if providerID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid provider id")
		return
	}

	output, err := c.listItemsUC.Execute(r.Context(), app.ListProviderServiceItemsInput{ProviderID: providerID})
	if err != nil {
		writeServiceDiscoveryError(w, err)
		return
	}

	items := make([]ServiceItemResponseDTO, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, ServiceItemResponseDTO{
			ID:            item.ID,
			ProviderID:    item.ProviderID,
			CategoryID:    item.CategoryID,
			Title:         item.Title,
			Description:   item.Description,
			PriceAmount:   item.PriceAmount,
			PriceUnit:     item.PriceUnit,
			SupportOnline: item.SupportOnline,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func writeServiceDiscoveryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProviderNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}

func buildProviderPublicDTO(provider domain.ProviderPublicProfile) ProviderPublicResponseDTO {
	return ProviderPublicResponseDTO{
		ID:                provider.ID,
		DisplayName:       provider.DisplayName,
		AvatarURL:         provider.AvatarURL,
		CityCode:          provider.CityCode,
		Bio:               provider.Bio,
		RatingAvg:         provider.RatingAvg,
		CompletedOrders:   provider.CompletedOrders,
		Online:            provider.Online,
		VerificationLabel: provider.VerificationText,
		Tags:              provider.Tags,
		PriceFrom:         provider.PriceFrom,
		PriceUnit:         provider.PriceUnit,
	}
}
