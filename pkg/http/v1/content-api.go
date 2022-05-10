package v1

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"


	"github.com/gin-gonic/gin"
	"github.com/tekenradar/content-service/pkg/types"
)

func (h *HttpEndpoints) AddContentAPI(rg *gin.RouterGroup) {
	data := rg.Group("/data")
	data.GET("/tb-report", h.getTBReportMapDataHandl)

}

func (h *HttpEndpoints) getTBReportMapDataHandl(c *gin.Context) {
	//fetch data from DB
	//date 4 weeks ago
	t :=time.Now().AddDate(0,0,-28).Unix()
	testInstanceID := strconv.FormatInt(1650359785, 10)
	points, err := h.contentDB.FindTickBiteMapDataByTime(testInstanceID,t)


	//prepare response in TBReportMapData Format
	if err != nil {
		fmt.Println("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch data from db"})
		return
	}

	var tempReportMapData types.ReportMapData;
	tempReportMapData.Slider.MinLabel = "-4 weken";
	tempReportMapData.Slider.MaxLabel = "nu";
	tempReportMapData.Slider.Labels = make([]string, 4)
	
	timeFormat := "02-01-2006"
	tempReportMapData.Slider.Labels[3] = time.Now().AddDate(0,0,-7).Format(timeFormat) + "  -  " + time.Now().Format(timeFormat)
	tempReportMapData.Slider.Labels[2] = time.Now().AddDate(0,0,-14).Format(timeFormat) + "  -  " + time.Now().AddDate(0,0,-8).Format(timeFormat)
	tempReportMapData.Slider.Labels[1] = time.Now().AddDate(0,0,-21).Format(timeFormat) + "  -  " + time.Now().AddDate(0,0,-15).Format(timeFormat)
	tempReportMapData.Slider.Labels[0] = time.Now().AddDate(0,0,-28).Format(timeFormat) + "  -  " + time.Now().AddDate(0,0,-22).Format(timeFormat)

	tempReportMapData.Series = make([][]types.TickBiteMapData, 4)
	for i := range points {

		TempMapData := types.TickBiteMapData{
			Lat : points[i].Lat,
			Lng : points[i].Lng,
			Type : points[i].Type,
		}
		switch temp_t := points[i].Time; {
		case temp_t < time.Now().AddDate(0,0,-21).Unix():
			tempReportMapData.Series[0] = append(tempReportMapData.Series[0],TempMapData)
		case temp_t < time.Now().AddDate(0,0,-14).Unix():
			tempReportMapData.Series[1] = append(tempReportMapData.Series[1],TempMapData)
		case temp_t < time.Now().AddDate(0,0,-7).Unix():
			tempReportMapData.Series[2] = append(tempReportMapData.Series[2],TempMapData)
		default:
			tempReportMapData.Series[3] = append(tempReportMapData.Series[3],TempMapData)
		}
	}

	c.JSON(http.StatusOK, tempReportMapData)

}
