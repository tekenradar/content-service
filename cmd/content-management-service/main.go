package main

import (
	"net/http"
	"time"

	"github.com/coneno/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tekenradar/content-service/pkg/dbs/contentdb"
	v1 "github.com/tekenradar/content-service/pkg/http/v1"
)

func healthCheckHandle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {

	conf := InitConfig()

	logger.SetLevel(conf.LogLevel)
	if !conf.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	contentDBService := contentdb.NewContentDBService(conf.ContentDBConfig, conf.InstanceIDs)

	// Start webserver
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		// AllowAllOrigins: true,
		AllowOrigins:     conf.AllowOrigins,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Content-Length", "Api-Key"},
		ExposeHeaders:    []string{"Authorization", "Content-Type", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/", healthCheckHandle)
	v1Root := router.Group("/v1")

	v1APIHandlers := v1.NewHTTPHandler(contentDBService, conf.APIKeyForReadOnly, conf.APIKeyForRW, conf.InstanceIDs, 0, conf.AssetsDir)
	v1APIHandlers.AddContentManagementAPI(v1Root)

	logger.Info.Printf("gateway listening on port %s", conf.Port)
	logger.Error.Fatal(router.Run(":" + conf.Port))
}
