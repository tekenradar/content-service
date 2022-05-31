package v1

import (
	"net/http"

	"github.com/coneno/logger"
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
