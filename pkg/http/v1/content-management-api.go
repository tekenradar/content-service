package v1

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/coneno/logger"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"github.com/influenzanet/study-service/pkg/studyengine"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tekenradar/content-service/pkg/http/helpers"
	mw "github.com/tekenradar/content-service/pkg/http/middlewares"
	"github.com/tekenradar/content-service/pkg/types"
)

func (h *HttpEndpoints) AddContentManagementAPI(rg *gin.RouterGroup) {
	studyevents := rg.Group("/study-event-handlers")
	studyevents.Use(mw.HasValidAPIKey(h.apiKeys.readWrite))
	{
		studyevents.POST("/tb-map-point-aggregator", mw.RequirePayload(), h.addTBReportHandl)
		studyevents.POST("/lpp-submission", mw.RequirePayload(), h.LPPSubmissionHandl)
	}

	instanceGroup := rg.Group("/:instanceID")
	instanceGroup.Use((mw.HasValidInstanceID(h.instanceIDs)))
	{
		data := instanceGroup.Group("/data")
		data.Use(mw.HasValidAPIKey(h.apiKeys.readWrite))
		{
			data.POST("/tb-map-data", mw.RequirePayload(), h.loadTBMapDataHandl)
		}

		filesW := instanceGroup.Group("/files")
		filesW.Use(mw.HasValidAPIKey(h.apiKeys.readWrite))
		{
			filesW.POST("", mw.RequirePayload(), h.uploadFileHandl)
			filesW.DELETE("", h.deleteFileHandl)
		}
		filesR := instanceGroup.Group("/files")
		filesR.Use(mw.HasValidAPIKey(h.apiKeys.readOnly))
		{
			filesR.GET("", h.getFileInfosHandl)
		}
		newsitemsR := instanceGroup.Group("/news-items")
		newsitemsR.Use(mw.HasValidAPIKey(h.apiKeys.readOnly))
		{
			newsitemsR.GET("", h.getNewsItemsHandl)
		}
		newsitemsW := instanceGroup.Group("/news-items")
		newsitemsW.Use(mw.HasValidAPIKey(h.apiKeys.readWrite))
		{
			newsitemsW.POST("", h.addNewsItemHandl)
		}
	}
}

func (h *HttpEndpoints) addTBReportHandl(c *gin.Context) {
	var req studyengine.ExternalEventPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error.Printf("error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	InstanceID := req.InstanceID
	err := helpers.CheckInstanceID(h.instanceIDs, InstanceID)
	if err != nil {
		logger.Error.Printf("error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	TBmapData, err := helpers.StudyEventToTBMapData(req)
	if err != nil {
		logger.Error.Printf("error while processing study event to TBMapData: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// save TBmapdata into DB
	_, err = h.contentDB.AddTickBiteMapData(InstanceID, TBmapData)
	if err != nil {
		logger.Error.Printf("unexpected error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add data to data base"})
		return
	}

	// prepare response
	c.JSON(http.StatusOK, gin.H{
		"message": "Map Data successfully added to data base"})
}

func (h *HttpEndpoints) LPPSubmissionHandl(c *gin.Context) {
	var req studyengine.ExternalEventPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error.Printf("error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	InstanceID := req.InstanceID
	err := helpers.CheckInstanceID(h.instanceIDs, InstanceID)
	if err != nil {
		logger.Error.Printf("error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := h.contentDB.GetLPPParticipant(InstanceID, "todo")
	if err != nil {
		logger.Error.Printf("error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSubmissions := p.Submissions
	newSubmissions[req.Response.Key] = time.Now()

	err = h.contentDB.UpdateLPPParticipantSubmissions(InstanceID, p.PID, newSubmissions, &types.TempParticipantInfo{
		ID:        req.ParticipantState.ParticipantID,
		EnteredAt: req.ParticipantState.EnteredAt,
	})
	if err != nil {
		logger.Error.Printf("error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info.Printf("LPP submission received for instance %s", InstanceID)
	c.JSON(http.StatusOK, gin.H{
		"message": "LPP submission registered successfully"})
}

func (h *HttpEndpoints) loadTBMapDataHandl(c *gin.Context) {
	instanceID := c.Param("instanceID")

	var TBMapData []types.TickBiteMapData
	if err := c.ShouldBindJSON(&TBMapData); err != nil {
		logger.Error.Printf("error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, d := range TBMapData {
		if _, err := h.contentDB.AddTickBiteMapData(instanceID, d); err != nil {
			logger.Error.Printf("Unable to add data to db: [%d]: %v, %v", i, d, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add data to data base"})
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Map Data loading finished"})
}

func (h *HttpEndpoints) uploadFileHandl(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Debug.Printf("error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	extension := filepath.Ext(file.Filename)

	//get file type
	fileContent, _ := file.Open()
	buffer := make([]byte, 512)
	_, err = fileContent.Read(buffer)
	if err != nil {
		logger.Error.Printf("unexpected error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error reading file",
		})
		return
	}
	kind, _ := filetype.Match(buffer)
	if kind == filetype.Unknown {
		logger.Error.Printf("unexpected error: file type is unknown")
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"message": "Unknown file type",
		})
		return
	}

	// Create file reference entry in DB
	instanceID := c.Param("instanceID")
	db_ID := primitive.NewObjectID()
	newFileName := db_ID.Hex() + extension
	dst := path.Join(h.assetsDir, newFileName)

	fi, err := h.contentDB.SaveFileInfo(instanceID, types.FileInfo{
		ID:         db_ID,
		Path:       dst,
		UploadedAt: time.Now().Unix(),
		FileType:   kind.MIME.Value,
		Name:       newFileName,
		Size:       int32(file.Size),
	})
	if err != nil {
		logger.Error.Printf("error saving file in db: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Unable to save file in db": err.Error()})
		return
	}

	err = os.MkdirAll(h.assetsDir, os.ModePerm)
	if err != nil {
		logger.Error.Printf("error uploading file: err at target path mkdir %v", err.Error())

		//if error delete db object
		_, err = h.contentDB.DeleteFileInfo(instanceID, fi.ID.String())
		if err != nil {
			logger.Error.Printf("unexpected error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"unexpected error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to upload file",
		})
		return
	}

	if err := c.SaveUploadedFile(file, dst); err != nil {
		logger.Error.Println("error while saving file at ", dst)

		//if error delete db object
		_, err = h.contentDB.DeleteFileInfo(instanceID, fi.ID.String())
		if err != nil {
			logger.Error.Printf("unexpected error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"unexpected error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File has been successfully uploaded."})
}

func (h *HttpEndpoints) deleteFileHandl(c *gin.Context) {
	instanceID := c.Param("instanceID")
	fileID := c.DefaultQuery("fileID", "")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fileID is empty"})
		return
	}

	fileIDs := strings.Split(fileID, ",")

	for _, id := range fileIDs {
		fileInfo, err := h.contentDB.FindFileInfo(instanceID, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"unexpected error": "file info not found"})
			return
		}

		// delete file
		err = os.Remove(fileInfo.Path)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"unexpected error": err.Error()})
			continue
		}

		// remove file info
		count, err := h.contentDB.DeleteFileInfo(instanceID, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"unexpected error": err.Error()})
			continue
		}
		logger.Debug.Printf("%d file info removed", count)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "File has been successfully removed."})
}

func (h *HttpEndpoints) getFileInfosHandl(c *gin.Context) {
	instanceID := c.Param("instanceID")
	//fetch data from DB
	fileInfoList, err := h.contentDB.GetFileInfoList(instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch data from db"})
		return
	}

	c.JSON(http.StatusOK, fileInfoList)
}

func (h *HttpEndpoints) getNewsItemsHandl(c *gin.Context) {
	instanceID := c.Param("instanceID")
	//fetch data from DB
	newsItemList, err := h.contentDB.GetNewsItemsList(instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch data from db"})
		return
	}

	c.JSON(http.StatusOK, newsItemList)
}

func (h *HttpEndpoints) addNewsItemHandl(c *gin.Context) {
	var req types.NewsItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	instanceID := c.Param("instanceID")

	// save TBmapdata into DB
	_, err := h.contentDB.AddNewsItem(instanceID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to add news item to data base"})
		return
	}

	// prepare response
	c.JSON(http.StatusOK, gin.H{
		"message": "News Item successfully added to data base"})

}
