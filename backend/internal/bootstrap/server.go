package bootstrap

import (
	"log"
	"net/http"
	"strings"

	aiApp "listen/backend/internal/application/ai"
	providerApp "listen/backend/internal/application/provider"
	domainAi "listen/backend/internal/domain/ai"
	infraAi "listen/backend/internal/infrastructure/ai"
	"listen/backend/internal/infrastructure/config"
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
	matchService := infraAi.NewMockMatchService()
	quotaRepo := memory.NewMatchQuotaRepository()
	sessionRepo := memory.NewSessionRepository()
	providerRepo := memory.NewProviderRepository()

	if cfg.RepositoryDriver == "mysql" {
		db, err := mysqlInfra.NewDB(cfg.MySQLDSN)
		if err != nil {
			log.Printf("mysql init failed, fallback to memory: %v", err)
		} else {
			quotaRepo = nil
			sessionRepo = nil
			providerRepo = nil
			mysqlQuotaRepo := mysqlInfra.NewMatchQuotaRepository(db)
			mysqlSessionRepo := mysqlInfra.NewSessionRepository(db)
			mysqlProviderRepo := mysqlInfra.NewProviderRepository(db)

			aiController := user.NewAIController(
				aiApp.NewGetAiHomeUseCase(mysqlQuotaRepo, clock),
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

			mux := http.NewServeMux()
			user.RegisterAIRoutes(mux, aiController)
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
		aiApp.NewGetAiHomeUseCase(quotaRepo, clock),
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

	mux := http.NewServeMux()
	user.RegisterAIRoutes(mux, aiController)
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
