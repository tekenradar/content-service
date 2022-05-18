package v1

import (
	"strconv"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/study-service/pkg/studyengine"

	mw "github.com/tekenradar/content-service/pkg/http/middlewares"
	cstypes "github.com/tekenradar/content-service/pkg/types"
)

func (h *HttpEndpoints) AddContentManagementAPI(rg *gin.RouterGroup) {
	data := rg.Group("/data")
	data.Use(mw.HasValidAPIKey(h.apiKeys.readWrite))
	{
		data.POST("/tb-report", mw.RequirePayload(), h.addTBReportHandl)
	}
}

func (h *HttpEndpoints) addTBReportHandl(c *gin.Context) {

	var load studyengine.ExternalEventPayload

	if err := c.ShouldBindJSON(&load); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// POST body -> TBmapdata
	var TBdata cstypes.TickBiteMapData
	TBdata.Time = load.Response.SubmittedAt

	if strings.Contains(load.Response.Key, "TB") {
		TBdata.Type = "TB"
	} else if strings.Contains(load.Response.Key, "EM") {
		TBdata.Type = "EM"
	} else if strings.Contains(load.Response.Key, "Fever") {
		TBdata.Type = "FE"
	} else {
		TBdata.Type = "Other"
	}

	//find map data in responses
	//TODO: give error if element is not found
	for i := range load.Response.Responses {
		if strings.Contains(load.Response.Responses[i].Key, "TBLoc.Q2") {

			for _, item := range load.Response.Responses[i].Response.Items {

				if item.Key == "map" {
					for _, mapItem := range item.Items {
						if mapItem.Key == "lat" {
							lat, err := strconv.ParseFloat(mapItem.Value, 64)
							if err != nil {
								//Correct handling of error here?
								c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
							}
							TBdata.Lat = lat

						}
						if mapItem.Key == "lng" {
							lng, err := strconv.ParseFloat(mapItem.Value, 64)
							if err != nil {
								//Correct handling of error here?
								c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
							}
							TBdata.Lng = lng

						}
					}
					break
				}
			}
			break
		}
	}

	// save TBmapdata into DB
	_, err := h.contentDB.AddTickBiteMapData(load.InstanceID, TBdata)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to add data to data base"})
		return
	}

	// prepare response
	c.JSON(http.StatusOK, gin.H{
		"message": "Map Data successfully added to data base"})

}
