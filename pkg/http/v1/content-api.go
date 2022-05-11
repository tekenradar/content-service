package v1

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/coneno/logger"
	"github.com/gin-gonic/gin"
	"github.com/tekenradar/content-service/pkg/types"
)

func (h *HttpEndpoints) AddContentAPI(rg *gin.RouterGroup) {
	data := rg.Group("/data")
	data.GET("/tb-report", h.getTBReportMapDataHandl)

}

func (h *HttpEndpoints) getTBReportMapDataHandl(c *gin.Context) {

	//fetch data from DB
	nParam := c.DefaultQuery("weeks", "4")
	n, err := strconv.Atoi(nParam)
	if err != nil {
		logger.Error.Fatal("Could not read weeks parameter")
	}

	t := time.Now().AddDate(0, 0, -(n * 7)).Unix()
	//testInstanceID := strconv.FormatInt(1650359785, 10)
	InstanceID := c.DefaultQuery("InstanceID", "")
	points, err := h.contentDB.FindTickBiteMapDataByTime(InstanceID, t)

	//prepare response in TBReportMapData Format
	if err != nil {
		fmt.Println("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch data from db"})
		return
	}

	var tempReportMapData types.ReportMapData
	tempReportMapData.Slider.MinLabel = "-" + strconv.Itoa(n) + " weken"
	tempReportMapData.Slider.MaxLabel = "nu"
	tempReportMapData.Slider.Labels = make([]string, n)

	timeFormat := "02-01-2006"

	for i := 0; i < n; i++ {
		tempReportMapData.Slider.Labels[i] = time.Now().AddDate(0, 0, -(7*(n-i))).Format(timeFormat) + "  -  " + time.Now().AddDate(0, 0, -(7*(n-i-1)+1)).Format(timeFormat)
		//add todays date to current week
		if i == (n - 1) {
			tempReportMapData.Slider.Labels[i] = time.Now().AddDate(0, 0, -(7*(n-i))).Format(timeFormat) + "  -  " + time.Now().AddDate(0, 0, -(7*(n-1-i))).Format(timeFormat)
		}
	}

	tempReportMapData.Series = make([][]types.TickBiteMapData, n)

	for _, point := range points {

		TempMapData := types.TickBiteMapData{
			Lat:  point.Lat,
			Lng:  point.Lng,
			Type: point.Type,
		}
		tDays := time.Unix(point.Time, 0).Sub(time.Now().AddDate(0, 0, -n*7)).Hours() / 24
		index := int64(tDays / 7)
		tempReportMapData.Series[index] = append(tempReportMapData.Series[index], TempMapData)
	}

	c.JSON(http.StatusOK, tempReportMapData)

}
