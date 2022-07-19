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

type MapDataCache struct {
	LastUpdated int64
	Data        types.ReportMapData
}

var (
	Cache map[string]map[int]MapDataCache
)

func (h *HttpEndpoints) AddContentAPI(rg *gin.RouterGroup) {
	instanceGroup := rg.Group("/:instanceID")
	instanceGroup.Use((mw.HasValidInstanceID(h.instanceIDs)))
	{
		data := instanceGroup.Group("/data")
		data.Use(mw.HasValidAPIKey(h.apiKeys.readOnly))
		{
			data.GET("/tb-report", h.getTBReportMapDataHandl)
		}
		newsitems := instanceGroup.Group("/news-items")
		newsitems.GET("", h.getPublishedNewsItemsHandl)
		newsitem := instanceGroup.Group("/news-item")
		newsitem.GET("", h.getNewsItemHandl)

	}
	files := rg.Group("/files")
	files.Static("/assets", h.assetsDir)
}

func (h *HttpEndpoints) getTBReportMapDataHandl(c *gin.Context) {
	nParam := c.DefaultQuery("weeks", "4")
	n, err := strconv.Atoi(nParam)
	if err != nil {
		logger.Error.Println("Could not read weeks parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read weeks parameter"})
	}

	t := time.Now().AddDate(0, 0, -(n * 7)).Unix()
	instanceID := c.Param("instanceID")

	if Cache == nil {
		Cache = make(map[string]map[int]MapDataCache)
	} else {
		if mdcache, ok := Cache[instanceID][n]; ok {
			if time.Since(time.Unix(mdcache.LastUpdated, 0)).Seconds() < float64(h.mapDataStoringDuration) {
				c.JSON(http.StatusOK, mdcache.Data)
				return
			}
		}
	}

	//fetch data from DB
	points, err := h.contentDB.FindTickBiteMapDataNewerThan(instanceID, t)
	if err != nil {
		logger.Error.Printf("unexpected error: %v", err)
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
			logger.Error.Printf("error while allocating map data to time intervals, check time of map data: %v", point.Time)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "could not allocate map data to time interval"})
			return
		}

		rmd.Series[index] = append(rmd.Series[index], md)
	}

	Cache[instanceID] = make(map[int]MapDataCache)
	Cache[instanceID][n] = MapDataCache{
		LastUpdated: time.Now().Unix(),
		Data:        rmd,
	}

	c.JSON(http.StatusOK, rmd)
}

func (h *HttpEndpoints) getNewsItemHandl(c *gin.Context) {
	instanceID := c.Param("instanceID")
	newsItemID := c.DefaultQuery("news-item-ID", "")
	if newsItemID == "" {
		logger.Error.Printf("error: ID of news item is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID of news item is empty"})
		return
	}
	newsItem, err := h.contentDB.FindNewsItem(instanceID, newsItemID)
	if err != nil {
		logger.Error.Printf("unexpected error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "news item not found"})
		return
	}
	if newsItem.Status != "published" {
		logger.Error.Printf("error: news item is not published")
		c.JSON(http.StatusBadRequest, gin.H{"error": "news item not found"})
		return
	}
	c.JSON(http.StatusOK, newsItem)
}

func (h *HttpEndpoints) getPublishedNewsItemsHandl(c *gin.Context) {
	instanceID := c.Param("instanceID")
	fromString := c.DefaultQuery("from", "0")
	from, err := strconv.ParseInt(fromString, 10, 64)
	if err != nil {
		logger.Error.Println("Could not read start date parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read start date parameter"})
	}
	untilString := c.DefaultQuery("until", "0")
	until, _ := strconv.ParseInt(untilString, 10, 64)
	if err != nil {
		logger.Error.Println("Could not read end date parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read end date parameter"})
	}
	if until == 0 {
		until = time.Now().Unix()
	}
	nString := c.DefaultQuery("n", "-1")
	n, err := strconv.Atoi(nString)
	if err != nil {
		logger.Error.Println("Could not read number of news items parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read number of news items parameter"})
	}

	if from > until {
		logger.Error.Println("error: end date of news items interval is older than start date")
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date of news items interval is older than start date"})
	}

	newsItemList, err := h.contentDB.FindNewsItemsInTimeInterval(instanceID, from, until, true)
	if err != nil {
		logger.Error.Printf("unexpected error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch data from db"})
		return
	}
	if n >= 0 && n < len(newsItemList) {
		newsItemList = newsItemList[:n]
	}

	c.JSON(http.StatusOK, newsItemList)
}
