package main

import (
	"fmt"
	"strconv"

	//"time"

	"github.com/coneno/logger"
	"github.com/tekenradar/content-service/pkg/dbs/contentdb"
	"github.com/tekenradar/content-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {

	conf := InitConfig()

	logger.SetLevel(conf.LogLevel)

	contentDBService := contentdb.NewContentDBService(conf.ContentDBConfig)

	logger.Debug.Println(contentDBService)

	testMapData := types.MapData{
		Time: 1650530469, 
		Lng: 21.262332,
		Lat: 6.34534,
		Type: "TB",
	}
	
	//var (
	//testInstanceID = strconv.FormatInt(time.Now().Unix(), 10)
	//)

	testInstanceID := strconv.FormatInt(1650359785,10)

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
	fmt.Printf("Time: %v\tLatitude: %v\tLongitude: %v\n",el.Time, el.Lat, el.Lng)
	}

}
