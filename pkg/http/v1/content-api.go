package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/coneno/logger"
	"github.com/gin-gonic/gin"
	mw "github.com/tekenradar/content-service/pkg/http/middlewares"
	"github.com/tekenradar/content-service/pkg/types"
)

func (h *HttpEndpoints) AddContentAPI(rg *gin.RouterGroup) {
	data := rg.Group("/:instanceID/data")
	data.Use(mw.HasValidAPIKey(h.apiKeys.readOnly))
	{
		data.GET("/tb-report", h.getTBReportMapDataHandl)
	}
	files := rg.Group("/files")
	files.Static("/assets", h.assetsDir) 
}

func (h *HttpEndpoints) getTBReportMapDataHandl(c *gin.Context) {
	nParam := c.DefaultQuery("weeks", "4")
	n, err := strconv.Atoi(nParam)
	if err != nil {
		logger.Error.Println("Could not read weeks parameter")
	}

	t := time.Now().AddDate(0, 0, -(n * 7)).Unix()
	instanceID := c.Param("instanceID")
	if instanceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "instanceID is empty"})
		return
	}

	//fetch data from DB
	points, err := h.contentDB.FindTickBiteMapDataNewerThan(instanceID, t)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch data from db"})
		return
	}

	//prepare response in TBReportMapData Format
	var rmd types.ReportMapData
	rmd.Slider.MinLabel = fmt.Sprintf("-%d weken", n)
	rmd.Slider.MaxLabel = "nu"
	rmd.Slider.Labels = make([]string, n)

	timeFormat := "02-01-2006"

	for i := 0; i < n; i++ {
		start_date_days := -(7 * (n - i))
		end_date_days := -(7*(n-i-1) + 1)
		if i == (n - 1) {
			end_date_days = -(7 * (n - 1 - i)) //add todays date to current week
		}
		rmd.Slider.Labels[i] = time.Now().AddDate(0, 0, start_date_days).Format(timeFormat) + "  -  " + time.Now().AddDate(0, 0, end_date_days).Format(timeFormat)
	}

	rmd.Series = make([][]types.TickBiteMapData, n)
	for i := range rmd.Series {
		rmd.Series[i] = make([]types.TickBiteMapData, 0)
	}

	for _, point := range points {
		end_date_days := -(7*n + 1)
		if time.Unix(point.Time, 0).Format("02-01-2006") == time.Now().Format("02-01-2006") {
			end_date_days = -(7 * n) //handle different if date is today
		}
		md := types.TickBiteMapData{
			Lat:  point.Lat,
			Lng:  point.Lng,
			Type: point.Type,
		}
		t := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 999999999, time.Local)
		tDays := time.Unix(point.Time, 0).Sub(t.AddDate(0, 0, end_date_days)).Hours() / 24
		index := int64(tDays / 7)
		if index >= int64(n) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not allocate map data to time interval"})
			return
		}

		rmd.Series[index] = append(rmd.Series[index], md)
	}

	c.JSON(http.StatusOK, rmd)

}
