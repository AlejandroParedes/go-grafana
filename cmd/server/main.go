package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-grafana/docs"
	"go-grafana/internal/config"
	"go-grafana/internal/domain/repository"
	"go-grafana/internal/handler"
	"go-grafana/internal/middleware"
	"go-grafana/internal/service"
	"go-grafana/pkg/database"
	"go-grafana/pkg/metrics"
	"go-grafana/pkg/sentry"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/TheZeroSlave/zapsentry"
	sentrysdk "github.com/getsentry/sentry-go"
)

// @title Go Grafana Web API
// @version 1.0
// @description A Go web application with Grafana monitoring
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	// Initialize Swagger info
	docs.SwaggerInfo.Title = "Go Grafana Web API"
	docs.SwaggerInfo.Description = "A Go web application with Grafana monitoring"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	app := fx.New(
		// Provide all dependencies
		fx.Provide(
			config.NewConfig,
			newLogger,
			database.NewPostgresDB,
			func() prometheus.Registerer { return prometheus.DefaultRegisterer },
			metrics.NewPrometheusMetrics,
			repository.NewUserRepository,
			repository.NewAPIKeyRepository,
			service.NewUserService,
			service.NewAPIKeyService,
			middleware.NewLoggingMiddleware,
			middleware.NewMetricsMiddleware,
			middleware.NewCORSMiddleware,
			handler.NewUserHandler,
			handler.NewAPIKeyHandler,
			newGinEngine,
			newHTTPServer,
		),
		// Invoke the server startup
		fx.Invoke(startServer),
		fx.Invoke(sentry.InitSentry),
		// Configure logging
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
	)

	// Start the application
	app.Run()
}

// newLogger creates a new Zap logger
func newLogger(cfg *config.Config) *zap.Logger {
	var logger *zap.Logger
	var err error

	// Default Zap logger configuration
	zapConfig := zap.NewProductionConfig()

	// Set log level
	switch cfg.Logging.Level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Build the logger
	logger, err = zapConfig.Build()
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}

	// Add Sentry core if DSN is configured
	if cfg.Sentry.DSN != "" {
		sentryCfg := zapsentry.Configuration{
			Level:             zapcore.ErrorLevel, //when to send message to sentry
			EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
			BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
		}
		sentryCore, err := zapsentry.NewCore(sentryCfg, zapsentry.NewSentryClientFromClient(sentrysdk.CurrentHub().Client()))
		if err != nil {
			logger.Error("Failed to create Sentry core for Zap", zap.Error(err))
		} else {
			logger = zapsentry.AttachCoreToLogger(sentryCore, logger)
		}
	}

	return logger
}

// newGinEngine creates a new Gin engine with middleware
func newGinEngine(
	loggingMiddleware middleware.LoggingMiddleware,
	metricsMiddleware middleware.MetricsMiddleware,
	corsMiddleware middleware.CORSMiddleware,
	userHandler *handler.UserHandler,
	apiKeyHandler *handler.APIKeyHandler,
	apiKeyService service.APIKeyService,
	logger *zap.Logger,
) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin engine
	engine := gin.New()

	// Add middleware
	engine.Use(loggingMiddleware.Handle())
	engine.Use(metricsMiddleware.Handle())
	engine.Use(corsMiddleware.Handle())
	engine.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	// Create API key authentication middleware
	apiKeyAuthMiddleware := middleware.APIKeyAuthMiddleware(apiKeyService, logger)

	// API routes
	api := engine.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Service is healthy",
				"time":    time.Now().UTC(),
			})
		})

		// Metrics endpoint for Prometheus
		api.GET("/metrics", metricsMiddleware.MetricsHandler())
		api.GET("/foo", func(ctx *gin.Context) {
			panic("test")
		})

		// User routes
		users := api.Group("/users")
		{
			// Public endpoints (no API key required)
			users.GET("/", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUserByID)

			// Protected endpoints (API key required)
			users.POST("/", apiKeyAuthMiddleware, userHandler.CreateUser)
			users.PUT("/:id", apiKeyAuthMiddleware, userHandler.UpdateUser)
			users.DELETE("/:id", apiKeyAuthMiddleware, userHandler.DeleteUser)
		}

		// API Key management routes (protected by API key)
		apiKeys := api.Group("/api-keys")
		{
			apiKeys.POST("/", apiKeyAuthMiddleware, apiKeyHandler.CreateAPIKey)
			apiKeys.GET("/", apiKeyAuthMiddleware, apiKeyHandler.GetAPIKeys)
			apiKeys.GET("/:id", apiKeyAuthMiddleware, apiKeyHandler.GetAPIKeyByID)
			apiKeys.PUT("/:id", apiKeyAuthMiddleware, apiKeyHandler.UpdateAPIKey)
			apiKeys.DELETE("/:id", apiKeyAuthMiddleware, apiKeyHandler.DeleteAPIKey)
		}
	}

	// Swagger documentation
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return engine
}

// newHTTPServer creates a new HTTP server
func newHTTPServer(engine *gin.Engine, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: engine,
	}
}

// startServer starts the HTTP server with graceful shutdown
func startServer(lifecycle fx.Lifecycle, server *http.Server, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting HTTP server", zap.String("addr", server.Addr))
			// Start server in a goroutine
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("Failed to start server", zap.Error(err))
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server")
			// Create a channel to listen for interrupt signals
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			// Create a context with timeout for graceful shutdown
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				logger.Error("Server forced to shutdown", zap.Error(err))
				return err
			}
			logger.Info("Server exited")
			return nil
		},
	})
}
