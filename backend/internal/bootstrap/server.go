package bootstrap

import (
	"log"
	"net/http"
	"strings"

	adminAuthApp "listen/backend/internal/application/admin_auth"
	adminOrderApp "listen/backend/internal/application/admin_order"
	aiApp "listen/backend/internal/application/ai"
	audioApp "listen/backend/internal/application/audio"
	feedbackApp "listen/backend/internal/application/feedback"
	identityApp "listen/backend/internal/application/identity"
	orderApp "listen/backend/internal/application/order"
	providerApp "listen/backend/internal/application/provider"
	providerAuthApp "listen/backend/internal/application/provider_auth"
	serviceDiscoveryApp "listen/backend/internal/application/service_discovery"
	adminServiceItemApp "listen/backend/internal/application/service_item_admin"
	userSettingsApp "listen/backend/internal/application/user_settings"
	adminAuthDomain "listen/backend/internal/domain/admin_auth"
	domainAi "listen/backend/internal/domain/ai"
	providerAuthDomain "listen/backend/internal/domain/provider_auth"
	serviceItemDomain "listen/backend/internal/domain/service_item_admin"
	userSettingsDomain "listen/backend/internal/domain/user_settings"
	infraAi "listen/backend/internal/infrastructure/ai"
	infraAudio "listen/backend/internal/infrastructure/audio"
	"listen/backend/internal/infrastructure/config"
	infraIdentity "listen/backend/internal/infrastructure/identity"
	memory "listen/backend/internal/infrastructure/persistence/memory"
	mysqlInfra "listen/backend/internal/infrastructure/persistence/mysql"
	adminHTTP "listen/backend/internal/interface/http/admin"
	providerHTTP "listen/backend/internal/interface/http/provider"
	"listen/backend/internal/interface/http/user"
)

type Server struct {
	Port string
	http *http.Server
}

func NewServer() Server {
	cfg := config.LoadServerConfig()

	clock := aiApp.SystemClock{}
	idGenerator := aiApp.NewTimestampIDGenerator(clock)
	homeOverviewService := infraAi.NewMockHomeOverviewService()
	matchService := infraAi.NewMockMatchService()
	soundPageService := infraAudio.NewMockSoundPageService()
	quotaRepo := memory.NewMatchQuotaRepository()
	sessionRepo := memory.NewSessionRepository()
	providerRepo := memory.NewProviderRepository()
	serviceDiscoveryRepo := memory.NewServiceDiscoveryRepository()
	identityRepo := memory.NewIdentityRepository()
	userSettingsRepo := userSettingsDomain.Repository(memory.NewUserSettingsRepository())
	orderRepo := memory.NewOrderRepository()
	feedbackRepo := memory.NewFeedbackRepository()
	serviceItemRepo := serviceItemDomain.Repository(memory.NewServiceItemAdminRepository())
	authService := infraIdentity.NewMockAuthService()
	adminAuthRepo := adminAuthApp.NewInMemoryRepository([]adminAuthDomain.AdminAccount{
		{
			AdminID:     "admin_001",
			Account:     "admin",
			Password:    "admin123",
			Role:        "super_admin",
			DisplayName: "平台管理员",
			Status:      "active",
		},
		{
			AdminID:     "admin_002",
			Account:     "reviewer",
			Password:    "reviewer123",
			Role:        "reviewer",
			DisplayName: "审核专员",
			Status:      "active",
		},
	})
	providerAuthRepo := providerAuthApp.NewInMemoryRepository([]providerAuthDomain.ProviderAccount{
		{
			ProviderID:  "p_pub_001",
			Account:     "provider",
			Password:    "provider123",
			DisplayName: "暖心倾听师 · 小林",
			Status:      "active",
			CityCode:    "310000",
		},
		{
			ProviderID:  "p_pub_002",
			Account:     "listener2",
			Password:    "listener123",
			DisplayName: "倾听师 · 安安",
			Status:      "active",
			CityCode:    "110000",
		},
	})

	if cfg.RepositoryDriver == "mysql" {
		db, err := mysqlInfra.NewDB(cfg.MySQLDSN)
		if err != nil {
			log.Printf("mysql init failed, fallback to memory: %v", err)
		} else {
			quotaRepo = nil
			sessionRepo = nil
			providerRepo = nil
			serviceDiscoveryRepo = nil
			mysqlQuotaRepo := mysqlInfra.NewMatchQuotaRepository(db)
			mysqlSessionRepo := mysqlInfra.NewSessionRepository(db)
			mysqlProviderRepo := mysqlInfra.NewProviderRepository(db)
			mysqlServiceDiscoveryRepo := mysqlInfra.NewServiceDiscoveryRepository(db)
			mysqlIdentityRepo := mysqlInfra.NewIdentityRepository(db)
			mysqlUserSettingsRepo := mysqlInfra.NewUserSettingsRepository(db)
			mysqlOrderRepo := mysqlInfra.NewOrderRepository(db)
			mysqlFeedbackRepo := mysqlInfra.NewFeedbackRepository(db)
			mysqlServiceItemRepo := mysqlInfra.NewServiceItemAdminRepository(db)
			mysqlOrderActionRepo := mysqlInfra.NewOrderAdminActionRepository(db)

			aiController := user.NewAIController(
				aiApp.NewGetAiHomeUseCase(mysqlQuotaRepo, homeOverviewService, clock),
				aiApp.NewGetRemainingMatchUseCase(mysqlQuotaRepo, clock),
				aiApp.NewSubmitMatchUseCase(mysqlQuotaRepo, matchService, clock),
				aiApp.NewCreateAiSessionUseCase(mysqlSessionRepo, idGenerator),
				aiApp.NewGetAiSessionUseCase(mysqlSessionRepo),
				aiApp.NewAppendAiMessageUseCase(mysqlSessionRepo, clock),
			)
			adminController := adminHTTP.NewProviderController(
				providerApp.NewListReviewProvidersUseCase(mysqlProviderRepo),
				providerApp.NewGetProviderDetailUseCase(mysqlProviderRepo),
				providerApp.NewReviewProviderUseCase(mysqlProviderRepo),
			)
			adminAuthController := adminHTTP.NewAuthController(
				adminAuthApp.NewLoginMockUseCase(adminAuthRepo, clock),
				adminAuthApp.NewGetCurrentAdminUseCase(adminAuthRepo),
			)
			serviceDiscoveryController := user.NewServiceDiscoveryController(
				serviceDiscoveryApp.NewListServiceCategoriesUseCase(mysqlServiceDiscoveryRepo),
				serviceDiscoveryApp.NewListPublicProvidersUseCase(mysqlServiceDiscoveryRepo),
				serviceDiscoveryApp.NewGetPublicProviderUseCase(mysqlServiceDiscoveryRepo),
				serviceDiscoveryApp.NewListProviderServiceItemsUseCase(mysqlServiceDiscoveryRepo),
			)
			soundController := user.NewSoundController(
				audioApp.NewGetSoundPageUseCase(soundPageService),
			)
			identityController := user.NewIdentityController(
				identityApp.NewLoginBySMSUseCase(mysqlIdentityRepo, authService, clock, idGenerator),
				identityApp.NewLoginByWechatUseCase(mysqlIdentityRepo, authService, clock, idGenerator),
				identityApp.NewGetUserProfileUseCase(mysqlIdentityRepo),
				identityApp.NewSaveUserProfileUseCase(mysqlIdentityRepo),
				identityApp.NewSaveUserPersonalityUseCase(mysqlIdentityRepo),
				identityApp.NewSkipUserPersonalityUseCase(mysqlIdentityRepo),
			)
			orderController := user.NewOrderController(
				orderApp.NewCreateOrderUseCase(mysqlOrderRepo, idGenerator, clock),
				orderApp.NewListOrdersUseCase(mysqlOrderRepo),
				orderApp.NewGetOrderUseCase(mysqlOrderRepo),
				orderApp.NewPayOrderMockSuccessUseCase(mysqlOrderRepo, clock),
			)
			feedbackController := user.NewFeedbackController(
				feedbackApp.NewSubmitOrderFeedbackUseCase(mysqlFeedbackRepo, mysqlOrderRepo, idGenerator, clock),
				feedbackApp.NewGetOrderFeedbackUseCase(mysqlFeedbackRepo, mysqlOrderRepo),
			)
			settingsController := user.NewSettingsController(
				userSettingsApp.NewGetSettingsUseCase(mysqlUserSettingsRepo, mysqlIdentityRepo),
				userSettingsApp.NewSaveSettingsUseCase(mysqlUserSettingsRepo, mysqlIdentityRepo),
			)
			serviceItemController := adminHTTP.NewServiceItemController(
				adminServiceItemApp.NewListServiceItemsUseCase(mysqlServiceItemRepo),
				adminServiceItemApp.NewGetServiceItemDetailUseCase(mysqlServiceItemRepo),
				adminServiceItemApp.NewUpdateServiceItemStatusUseCase(mysqlServiceItemRepo),
			)
			orderAdminController := adminHTTP.NewOrderController(
				adminOrderApp.NewUseCase(mysqlOrderRepo, mysqlFeedbackRepo, mysqlOrderActionRepo, idGenerator, clock),
			)
			providerAuthController := providerHTTP.NewAuthController(
				providerAuthApp.NewLoginMockUseCase(providerAuthRepo, clock),
				providerAuthApp.NewGetCurrentProviderUseCase(providerAuthRepo),
			)
			providerOrderController := providerHTTP.NewOrderController(
				orderApp.NewProviderListOrdersUseCase(mysqlOrderRepo),
				orderApp.NewProviderGetOrderUseCase(mysqlOrderRepo),
				orderApp.NewProviderOperateOrderUseCase(mysqlOrderRepo),
			)

			mux := http.NewServeMux()
			user.RegisterAIRoutes(mux, aiController)
			user.RegisterIdentityRoutes(mux, identityController)
			user.RegisterServiceDiscoveryRoutes(mux, serviceDiscoveryController)
			user.RegisterSoundRoutes(mux, soundController)
			user.RegisterOrderRoutes(mux, orderController)
			user.RegisterFeedbackRoutes(mux, feedbackController)
			user.RegisterSettingsRoutes(mux, settingsController)
			adminHTTP.RegisterAuthRoutes(mux, adminAuthController)
			adminHTTP.RegisterProviderRoutes(mux, adminController)
			adminHTTP.RegisterServiceItemRoutes(mux, serviceItemController)
			adminHTTP.RegisterOrderRoutes(mux, orderAdminController)
			adminHTTP.RegisterComplaintRoutes(mux, orderAdminController)
			providerHTTP.RegisterAuthRoutes(mux, providerAuthController)
			providerHTTP.RegisterOrderRoutes(mux, providerOrderController)
			mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
			})
			return Server{
				Port: cfg.Port,
				http: &http.Server{
					Addr:    ":" + cfg.Port,
					Handler: corsMiddleware(mux),
				},
			}
		}
	}

	aiController := user.NewAIController(
		aiApp.NewGetAiHomeUseCase(quotaRepo, homeOverviewService, clock),
		aiApp.NewGetRemainingMatchUseCase(quotaRepo, clock),
		aiApp.NewSubmitMatchUseCase(quotaRepo, matchService, clock),
		aiApp.NewCreateAiSessionUseCase(sessionRepo, idGenerator),
		aiApp.NewGetAiSessionUseCase(sessionRepo),
		aiApp.NewAppendAiMessageUseCase(sessionRepo, clock),
	)
	adminController := adminHTTP.NewProviderController(
		providerApp.NewListReviewProvidersUseCase(providerRepo),
		providerApp.NewGetProviderDetailUseCase(providerRepo),
		providerApp.NewReviewProviderUseCase(providerRepo),
	)
	adminAuthController := adminHTTP.NewAuthController(
		adminAuthApp.NewLoginMockUseCase(adminAuthRepo, clock),
		adminAuthApp.NewGetCurrentAdminUseCase(adminAuthRepo),
	)
	serviceDiscoveryController := user.NewServiceDiscoveryController(
		serviceDiscoveryApp.NewListServiceCategoriesUseCase(serviceDiscoveryRepo),
		serviceDiscoveryApp.NewListPublicProvidersUseCase(serviceDiscoveryRepo),
		serviceDiscoveryApp.NewGetPublicProviderUseCase(serviceDiscoveryRepo),
		serviceDiscoveryApp.NewListProviderServiceItemsUseCase(serviceDiscoveryRepo),
	)
	soundController := user.NewSoundController(
		audioApp.NewGetSoundPageUseCase(soundPageService),
	)
	identityController := user.NewIdentityController(
		identityApp.NewLoginBySMSUseCase(identityRepo, authService, clock, idGenerator),
		identityApp.NewLoginByWechatUseCase(identityRepo, authService, clock, idGenerator),
		identityApp.NewGetUserProfileUseCase(identityRepo),
		identityApp.NewSaveUserProfileUseCase(identityRepo),
		identityApp.NewSaveUserPersonalityUseCase(identityRepo),
		identityApp.NewSkipUserPersonalityUseCase(identityRepo),
	)
	orderController := user.NewOrderController(
		orderApp.NewCreateOrderUseCase(orderRepo, idGenerator, clock),
		orderApp.NewListOrdersUseCase(orderRepo),
		orderApp.NewGetOrderUseCase(orderRepo),
		orderApp.NewPayOrderMockSuccessUseCase(orderRepo, clock),
	)
	feedbackController := user.NewFeedbackController(
		feedbackApp.NewSubmitOrderFeedbackUseCase(feedbackRepo, orderRepo, idGenerator, clock),
		feedbackApp.NewGetOrderFeedbackUseCase(feedbackRepo, orderRepo),
	)
	settingsController := user.NewSettingsController(
		userSettingsApp.NewGetSettingsUseCase(userSettingsRepo, identityRepo),
		userSettingsApp.NewSaveSettingsUseCase(userSettingsRepo, identityRepo),
	)
	serviceItemController := adminHTTP.NewServiceItemController(
		adminServiceItemApp.NewListServiceItemsUseCase(serviceItemRepo),
		adminServiceItemApp.NewGetServiceItemDetailUseCase(serviceItemRepo),
		adminServiceItemApp.NewUpdateServiceItemStatusUseCase(serviceItemRepo),
	)
	orderActionRepo := memory.NewOrderAdminActionRepository()
	orderAdminController := adminHTTP.NewOrderController(
		adminOrderApp.NewUseCase(orderRepo, feedbackRepo, orderActionRepo, idGenerator, clock),
	)
	providerAuthController := providerHTTP.NewAuthController(
		providerAuthApp.NewLoginMockUseCase(providerAuthRepo, clock),
		providerAuthApp.NewGetCurrentProviderUseCase(providerAuthRepo),
	)
	providerOrderController := providerHTTP.NewOrderController(
		orderApp.NewProviderListOrdersUseCase(orderRepo),
		orderApp.NewProviderGetOrderUseCase(orderRepo),
		orderApp.NewProviderOperateOrderUseCase(orderRepo),
	)

	mux := http.NewServeMux()
	user.RegisterAIRoutes(mux, aiController)
	user.RegisterIdentityRoutes(mux, identityController)
	user.RegisterServiceDiscoveryRoutes(mux, serviceDiscoveryController)
	user.RegisterSoundRoutes(mux, soundController)
	user.RegisterOrderRoutes(mux, orderController)
	user.RegisterFeedbackRoutes(mux, feedbackController)
	user.RegisterSettingsRoutes(mux, settingsController)
	adminHTTP.RegisterAuthRoutes(mux, adminAuthController)
	adminHTTP.RegisterProviderRoutes(mux, adminController)
	adminHTTP.RegisterServiceItemRoutes(mux, serviceItemController)
	adminHTTP.RegisterOrderRoutes(mux, orderAdminController)
	adminHTTP.RegisterComplaintRoutes(mux, orderAdminController)
	providerHTTP.RegisterAuthRoutes(mux, providerAuthController)
	providerHTTP.RegisterOrderRoutes(mux, providerOrderController)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return Server{
		Port: cfg.Port,
		http: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: corsMiddleware(mux),
		},
	}
}

func (s Server) Run() error {
	return s.http.ListenAndServe()
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

var _ domainAi.MatchService = infraAi.MockMatchService{}
