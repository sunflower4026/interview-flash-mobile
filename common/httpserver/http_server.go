package httpserver

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	"gitlab.com/sunflower4026/interview-flash-mobile/middleware"
	"gitlab.com/sunflower4026/interview-flash-mobile/routes"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/ginkit"
)

func RunHTTPService(ctx context.Context, cfg *viper.Viper, svc *httpservice.Service) {

	if cfg.GetString("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(ErrorHandler())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Token", "RefreshToken"}
	router.Use(cors.New(config))

	router.Use(middleware.CorrelationIDMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	runtimeCfg := ginkit.NewRuntimeConfig(cfg)
	runtimeCfg.HealthCheckFunc = svc.GetServiceHealth

	routes.Router(ctx, router, cfg, svc)

	ginkit.RunServerWithContext(ctx, router, runtimeCfg)
}
