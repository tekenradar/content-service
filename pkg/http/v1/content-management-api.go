package v1

import (
	"github.com/gin-gonic/gin"
	mw "github.com/tekenradar/content-service/pkg/http/middlewares"
)

func (h *HttpEndpoints) AddContentManagementAPI(rg *gin.RouterGroup) {
	data := rg.Group("/data")
	data.POST("/tb-report", mw.RequirePayload(), h.addTBReportHandl)

}

func (h *HttpEndpoints) addTBReportHandl(c *gin.Context) {
	// TODO: check if api key valid
	// TODO: for debugging, sve POST body as a JSON
	// TODO: POST body -> TBmapdata
	// TODO: save TBmapdata into DB
	// TODO: prepare response
}
