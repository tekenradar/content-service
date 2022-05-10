package v1

import (
	"fmt"
	"io/ioutil"

	"net/http"

	"github.com/gin-gonic/gin"
	mw "github.com/tekenradar/content-service/pkg/http/middlewares"
)

func (h *HttpEndpoints) AddContentManagementAPI(rg *gin.RouterGroup) {
	data := rg.Group("/data")
	data.Use(mw.HasValidAPIKey(h.apiKeys.readWrite))
	{
		data.POST("/tb-report", mw.RequirePayload(), h.addTBReportHandl)
	}
}

func (h *HttpEndpoints) addTBReportHandl(c *gin.Context) {

	// TODO: for debugging, save POST body as a JSON
	resp, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to read request body",
		})
		return
	}

	err1 := ioutil.WriteFile("test.json", resp, 0644)

	if err1 != nil {
		fmt.Println("error:", err1)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	// File saved successfully. Return proper result
	c.JSON(http.StatusOK, gin.H{
		"message": "Your file has been successfully saved."})

	// TODO: POST body -> TBmapdata
	/*req TBMapData
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("error:", err)
		return
	}*/
	//id, err := h.contentDB.AddTickBiteMapData(testInstanceID, testMapData)
	// TODO: save TBmapdata into DB
	// TODO: prepare response
}
