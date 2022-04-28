package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	//"time"

	"github.com/coneno/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tekenradar/content-service/pkg/dbs/contentdb"
	v1 "github.com/tekenradar/content-service/pkg/http/v1"
	"github.com/tekenradar/content-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func healthCheckHandle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {

	conf := InitConfig()

	logger.SetLevel(conf.LogLevel)

	contentDBService := contentdb.NewContentDBService(conf.ContentDBConfig)

	// ---> TEST CODE
	testMapData := types.MapData{
		Time: 1650530469,
		Lng:  21.262332,
		Lat:  6.34534,
		Type: "TB",
	}

	//var (
	//testInstanceID = strconv.FormatInt(time.Now().Unix(), 10)
	//)

	testInstanceID := strconv.FormatInt(1650359785, 10)

	id, err := contentDBService.AddMapData(testInstanceID, testMapData)

	if err != nil {
		logger.Error.Printf("unexpected error: %s", err.Error())
	}
	if len(id) < 2 || id == primitive.NilObjectID.Hex() {
		logger.Error.Printf("unexpected id: %s", id)
	}

	mapData, err := contentDBService.FindMapDataByTime(testInstanceID, 12323100000)

	if err != nil {
		logger.Error.Printf("unexpected error: %s", err.Error())
	}

	for _, el := range mapData {
		fmt.Printf("Time: %v\tLatitude: %v\tLongitude: %v\n", el.Time, el.Lat, el.Lng)
	}
	// <--- END TEST CODE

	// Start webserver
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		// AllowAllOrigins: true,
		AllowOrigins:     conf.AllowOrigins,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Authorization", "Content-Type", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/", healthCheckHandle)
	v1Root := router.Group("/v1")

	v1APIHandlers := v1.NewHTTPHandler(contentDBService)
	v1APIHandlers.AddContentManagementAPI(v1Root)

	logger.Info.Printf("gateway listening on port %s", conf.Port)
	logger.Error.Fatal(router.Run(":" + conf.Port))

}
