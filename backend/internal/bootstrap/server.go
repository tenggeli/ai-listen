package bootstrap

import (
	"log"
	"net/http"
	"strings"

	aiApp "listen/backend/internal/application/ai"
	audioApp "listen/backend/internal/application/audio"
	identityApp "listen/backend/internal/application/identity"
	providerApp "listen/backend/internal/application/provider"
	serviceDiscoveryApp "listen/backend/internal/application/service_discovery"
	domainAi "listen/backend/internal/domain/ai"
	infraAi "listen/backend/internal/infrastructure/ai"
	infraAudio "listen/backend/internal/infrastructure/audio"
	"listen/backend/internal/infrastructure/config"
	infraIdentity "listen/backend/internal/infrastructure/identity"
	memory "listen/backend/internal/infrastructure/persistence/memory"
	mysqlInfra "listen/backend/internal/infrastructure/persistence/mysql"
	adminHTTP "listen/backend/internal/interface/http/admin"
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
	authService := infraIdentity.NewMockAuthService()

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

			mux := http.NewServeMux()
			user.RegisterAIRoutes(mux, aiController)
			user.RegisterIdentityRoutes(mux, identityController)
			user.RegisterServiceDiscoveryRoutes(mux, serviceDiscoveryController)
			user.RegisterSoundRoutes(mux, soundController)
			adminHTTP.RegisterProviderRoutes(mux, adminController)
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

	mux := http.NewServeMux()
	user.RegisterAIRoutes(mux, aiController)
	user.RegisterIdentityRoutes(mux, identityController)
	user.RegisterServiceDiscoveryRoutes(mux, serviceDiscoveryController)
	user.RegisterSoundRoutes(mux, soundController)
	adminHTTP.RegisterProviderRoutes(mux, adminController)
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
