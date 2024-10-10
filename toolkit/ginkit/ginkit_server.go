package ginkit

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"

	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/runtimekit"

	stdLog "log"

	gin_prometheus "github.com/zsais/go-gin-prometheus"
)

const (
	defaultInfoPath   = "/actuator/info"
	defaultHealthPath = "/actuator/health"
	defaultReqTimeout = 7 * time.Second
	defaultPort       = 8088
)

// RuntimeConfig defines echo REST API runtime config with healthcheck.
type RuntimeConfig struct {
	Port                    int           `json:"port,omitempty"`
	Name                    string        `json:"name,omitempty"`
	BuildInfo               string        `json:"build_info,omitempty"`
	ShutdownWaitDuration    time.Duration `json:"shutdown_wait_duration,omitempty"`
	ShutdownTimeoutDuration time.Duration `json:"shutdown_timeout_duration,omitempty"`
	HealthCheckPath         string        `json:"health_check_path,omitempty"`
	InfoCheckPath           string        `json:"info_check_path,omitempty"`
	HealthCheckFunc         `json:"-"`
}

func (cfg *RuntimeConfig) validate() {
	// port
	if cfg.Port == 0 {
		cfg.Port = defaultPort
	}

	// healthcheck
	if cfg.HealthCheckPath == "" {
		cfg.HealthCheckPath = defaultHealthPath
	}
}

type healthStatus struct {
	serving bool
	Status  string `json:"status"`
}

// HealthCheckFunc is healthcheck interface func.
type HealthCheckFunc func(ctx context.Context) error

// RunServer run graceful restapi server.
func RunServer(e *gin.Engine, cfg *RuntimeConfig) {
	appCtx, done := runtimekit.NewRuntimeContext()
	defer done()

	RunServerWithContext(appCtx, e, cfg)
}

// RunServerWithContext run graceful restapi server with existing background context
// provides default '/actuator/health' as healthcheck endpoint
// provides '/metrics' as prometheus metrics endpoint.
// set echo.Validator using `web.Validator` from `web` package.
func RunServerWithContext(appCtx context.Context, e *gin.Engine, cfg *RuntimeConfig) {
	cfg.Name = strcase.ToSnake(cfg.Name)

	logger := log.FromCtx(appCtx)

	// validator := validator.New()

	cfg.validate()

	// request validator setup
	// e.Use(ValidatorTranslatorMiddleware(validator), TimeoutMiddleware(cfg.RequestTimeoutConfig))

	// healthcheck
	hs := &healthStatus{
		serving: true,
	}

	e.GET(cfg.HealthCheckPath, func(c *gin.Context) {
		if !hs.serving {
			hs.Status = "OUT_OF_SERVICE"

			c.JSON(http.StatusOK, hs)
			return
		}

		if cfg.HealthCheckFunc == nil {
			hs.Status = "UP"

			c.JSON(http.StatusOK, hs)
			return
		}

		err := cfg.HealthCheckFunc(c.Request.Context())
		if err != nil {
			hs.Status = "OUT_OF_SERVICE"

			c.JSON(http.StatusOK, hs)
			return
		}

		hs.Status = "UP"

		c.JSON(http.StatusOK, hs)
	})

	if cfg.InfoCheckPath == "" {
		cfg.InfoCheckPath = defaultInfoPath
	}

	e.GET(cfg.InfoCheckPath, func(c *gin.Context) {
		var v struct {
			Version string `json:"version"`
		}
		v.Version = cfg.BuildInfo

		c.JSON(http.StatusOK, v)
	})

	// prometheus
	p := gin_prometheus.NewPrometheus(cfg.Name, nil)
	p.Use(e)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: e,
	}

	go func() {
		<-appCtx.Done()

		hs.serving = false

		logger.Info(fmt.Sprintf("shutting down REST HTTP server in %d seconds", cfg.ShutdownWaitDuration))
		<-time.After(cfg.ShutdownWaitDuration)

		// stop the server
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeoutDuration)
		defer cancel()

		if err := s.Shutdown(shutdownCtx); err != nil {
			logger.Error(err, "shutdown http server")
		}
	}()

	// error fallback handler
	e.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	})

	PrintRoutes(e)

	// start server
	logger.Info("serving REST HTTP server", "config", cfg)

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err, "starting http server")
	}
}

// PrintRoutes prints *gin.Engine routes.
func PrintRoutes(e *gin.Engine) {
	stdLog.Println("=====> initializing http routes")

	for _, r := range e.Routes() {
		stdLog.Printf("=====> %s %s %s", r.Method, r.Path, r.Handler)
	}
}
