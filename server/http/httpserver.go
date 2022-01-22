package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manso/gobackup/server/appservices"
	"github.com/manso/gobackup/server/store"
)

var _bfStore *store.BackedFileStore
var _storePath string

func Start(bfStore *store.BackedFileStore, storePath string) {
	_bfStore = bfStore
	_storePath = storePath
	router := gin.Default()
	router.POST("/upload/:id", uploadController)
	router.GET("/files/", getFiles)
	router.GET("/file/:hash", downloadFileController)
	router.Run(":8080")
}

func getFiles(c *gin.Context) {
	result := appservices.GetBackedFiles(_bfStore)
	c.JSON(http.StatusOK, result)
}

func uploadController(c *gin.Context) {
	modifiedOn, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", c.Request.FormValue("ModifiedOn"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	formFile, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	bodyFile, err := formFile.Open()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	defer bodyFile.Close()

	err = appservices.UploadFile(_bfStore, _storePath, c.Param("id"), c.Request.FormValue("Path"), modifiedOn, formFile.Size, bodyFile)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.String(http.StatusOK, "")
}

func downloadFileController(c *gin.Context) {
	log.Println("Requested hash ", c.Param("hash"))
	filePath, err := appservices.DownloadBackedFileVersion(_bfStore, _storePath, c.Param("hash"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment")
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)

}
