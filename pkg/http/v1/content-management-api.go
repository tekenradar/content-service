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

	mw "github.com/tekenradar/content-service/pkg/http/middlewares"
	"github.com/tekenradar/content-service/pkg/types"
	cstypes "github.com/tekenradar/content-service/pkg/types"
)

func (h *HttpEndpoints) AddContentManagementAPI(rg *gin.RouterGroup) {
	studyevents := rg.Group("/study-event-handlers")
	studyevents.Use(mw.HasValidAPIKey(h.apiKeys.readWrite))
	{
		studyevents.POST("/tb-map-point-aggregator", mw.RequirePayload(), h.addTBReportHandl)
	}
	data := rg.Group("/:instanceID/data")
	data.Use(mw.HasValidAPIKey(h.apiKeys.readWrite))
	{
		data.POST("/tb-map-data", mw.RequirePayload(), h.loadTBMapDataHandl)
	}
	files := rg.Group("/:instanceID/files")
	files.Use(mw.HasValidAPIKey(h.apiKeys.readWrite))
	{
		files.POST("", mw.RequirePayload(), h.uploadFileHandl)
		files.DELETE("", mw.RequirePayload(), h.deleteFileHandl)
	}
	files.Use(mw.HasValidAPIKey(h.apiKeys.readOnly))
	{
		files.GET("", mw.RequirePayload(), h.getFileInfosHandl)
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
	instanceID := c.Param("instanceID")
	if instanceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "instanceID is empty"})
		return
	}

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

func (h *HttpEndpoints) uploadFileHandl(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error reading file",
		})
		return
	}
	kind, _ := filetype.Match(buffer)
	if kind == filetype.Unknown {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unknown file type",
		})
		return
	}

	// Create file reference entry in DB
	instanceID := c.Param("instanceID")
	if instanceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "instanceID is empty"})
		return
	}
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
		logger.Error.Printf("Error UploadFile: %v", err.Error())
		return
	}

	err = os.MkdirAll(h.assetsDir, os.ModePerm)
	if err != nil {
		logger.Info.Printf("Error uploading file: err at target path mkdir %v", err.Error())
		//if error delete db object
		_, err = h.contentDB.DeleteFileInfo(instanceID, fi.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"unexpected error": err.Error()})
		}
		return
	}

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to save the file",
		})
		logger.Info.Println(dst)
		//if error delete db object
		_, err = h.contentDB.DeleteFileInfo(instanceID, fi.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"unexpected error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File has been successfully uploaded."})
}

func (h *HttpEndpoints) deleteFileHandl(c *gin.Context) {
	instanceID := c.Param("instanceID")
	if instanceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "instanceID is empty"})
		return
	}
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
	if instanceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "instanceID is empty"})
		return
	}

	//fetch data from DB
	fileInfoList, err := h.contentDB.GetFileInfoList(instanceID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch data from db"})
		return
	}

	c.JSON(http.StatusOK, fileInfoList)
}
