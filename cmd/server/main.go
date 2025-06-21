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

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
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
// @name Authorization
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
			metrics.NewPrometheusMetrics,
			repository.NewUserRepository,
			service.NewUserService,
			middleware.NewLoggingMiddleware,
			middleware.NewMetricsMiddleware,
			middleware.NewCORSMiddleware,
			handler.NewUserHandler,
			newGinEngine,
			newHTTPServer,
		),
		// Invoke the server startup
		fx.Invoke(startServer),
		// Configure logging
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
	)

	// Start the application
	app.Run()
}

// newLogger creates a new Zap logger
func newLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	return logger
}

// newGinEngine creates a new Gin engine with middleware
func newGinEngine(
	loggingMiddleware middleware.LoggingMiddleware,
	metricsMiddleware middleware.MetricsMiddleware,
	corsMiddleware middleware.CORSMiddleware,
	userHandler *handler.UserHandler,
) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin engine
	engine := gin.New()

	// Add middleware
	engine.Use(loggingMiddleware.Handle())
	engine.Use(metricsMiddleware.Handle())
	engine.Use(corsMiddleware.Handle())
	engine.Use(gin.Recovery())

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

		// User routes
		users := api.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
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
