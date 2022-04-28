package v1

import (
	"github.com/gin-gonic/gin"
)

func (h *HttpEndpoints) AddContentAPI(rg *gin.RouterGroup) {
	data := rg.Group("/data")
	data.GET("/tb-report", h.getTBReportMapDataHandl)

}

func (h *HttpEndpoints) getTBReportMapDataHandl(c *gin.Context) {
	// TODO: fetch data from DB
	// TODO: prepare response
}
