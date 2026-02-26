package main

import (
	cfgPkg "github.com/ArtemNeGopher/url-shortener/pkg/config"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/config"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/grpc"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/handler"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// Конфиг
	config := &config.Config{}
	cfgPkg.MustInit("config/config.yaml", config)

	// gRPC клиенты
	cacheClient := grpc.NewCacheClient(&config.CacheClient)
	defer cacheClient.Close()
	urlClient := grpc.NewURLClient(&config.URLClient)
	defer urlClient.Close()
	analyticsClient := grpc.NewAnalyticsClient(&config.AnalyticsClient)
	defer analyticsClient.Close()

	// Сервисы
	urlService := service.NewURLService(cacheClient, urlClient)
	analyticsService := service.NewAnalyticsService(analyticsClient)

	h := handler.NewHandler(urlService, analyticsService)
	r := gin.Default()
	handler.RegisterRoutes(r, h)

	r.Run(config.Host + ":" + config.Port)
}
