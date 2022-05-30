package v1

import (
	"errors"
	"strconv"
	"strings"

	"net/http"

	"github.com/coneno/logger"
	"github.com/gin-gonic/gin"
	"github.com/influenzanet/study-service/pkg/studyengine"
	"github.com/influenzanet/study-service/pkg/types"

	mw "github.com/tekenradar/content-service/pkg/http/middlewares"
	cstypes "github.com/tekenradar/content-service/pkg/types"
)

func (h *HttpEndpoints) AddContentManagementAPI(rg *gin.RouterGroup) {
	data := rg.Group("/data")
	data.Use(mw.HasValidAPIKey(h.apiKeys.readWrite))
	{
		data.POST("/tb-report", mw.RequirePayload(), h.addTBReportHandl)
		data.POST("/initialise-dbcollection", mw.RequirePayload(), h.loadTBMapDataHandl)
	}
}

func (h *HttpEndpoints) addTBReportHandl(c *gin.Context) {

	var req studyengine.ExternalEventPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	TBmapData, err := studyEventToTBMapData(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// save TBmapdata into DB
	_, err = h.contentDB.AddTickBiteMapData(req.InstanceID, TBmapData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to add data to data base"})
		return
	}

	// prepare response
	c.JSON(http.StatusOK, gin.H{
		"message": "Map Data successfully added to data base"})

}

func studyEventToTBMapData(event studyengine.ExternalEventPayload) (tickBiteMapData cstypes.TickBiteMapData, err error) {

	item, err := findResponseItem(event.Response.Responses, "TBLoc.Q2")
	if err != nil {
		return cstypes.TickBiteMapData{}, err
	}

	lat, err := parseResponseValueAsFloat(item, "lat")
	if err != nil {
		return cstypes.TickBiteMapData{}, err
	}

	lng, err := parseResponseValueAsFloat(item, "lng")
	if err != nil {
		return cstypes.TickBiteMapData{}, err
	}

	rtype, err := getReportType(event.Response.Key)
	if err != nil {
		return cstypes.TickBiteMapData{}, err
	}

	time := event.Response.SubmittedAt

	return cstypes.TickBiteMapData{
		Time: time,
		Lat:  lat,
		Lng:  lng,
		Type: rtype}, nil

}

func getReportType(key string) (Rtype string, err error) {
	if strings.Contains(key, "TB") {
		return "TB", nil
	} else if strings.Contains(key, "EM") {
		return "EM", nil
	} else if strings.Contains(key, "Fever") {
		return "FE", nil
	} else if strings.Contains(key, "LB") {
		return "Other", nil
	} else if strings.Contains(key, "CH") {
		return "Other", nil
	} else {
		return "", errors.New("Could not allocate type value")
	}
}

func parseResponseValueAsFloat(mapItem []types.ResponseItem, name string) (value float64, err error) {
	for _, mapItem := range mapItem {
		if mapItem.Key == name {
			val, err := strconv.ParseFloat(mapItem.Value, 64)
			if err != nil {
				return val, errors.New("Could not parse response value to float")
			}
			return val, nil
		}
	}
	return 0, errors.New("Could not find response value")
}

func findResponseItem(response []types.SurveyItemResponse, itemKey string) (item []types.ResponseItem, err error) {
	for i := range response {
		if strings.Contains(response[i].Key, itemKey) {

			for _, item := range response[i].Response.Items {

				if item.Key == "map" {
					return item.Items, nil
				}

			}
		}
	}
	return []types.ResponseItem{}, errors.New("Could not find response item")
}

func (h *HttpEndpoints) loadTBMapDataHandl(c *gin.Context) {
	instanceID := c.DefaultQuery("instanceID", "")
	// TODO: check if instanceID exists to prevent empty instance ids

	var TBMapData []cstypes.TickBiteMapData
	if err := c.ShouldBindJSON(&TBMapData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, d := range TBMapData {
		if _, err := h.contentDB.AddTickBiteMapData(instanceID, d); err != nil {
			logger.Error.Printf("Unable to add data to db: [%d]: %v", i, d)
		}
	}

	err := h.contentDB.CreateIndex(instanceID)
	if err != nil {
		logger.Error.Printf("Unexpected error: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Map Data loading finished"})
}
