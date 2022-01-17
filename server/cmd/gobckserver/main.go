package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BackedFile struct {
	gorm.Model
	Path       string
	Hash       string
	Size       int64
	ModifiedOn time.Time
	UploadedOn time.Time
}

func main() {
	storePath := "./store"
	if len(os.Args) > 1 {
		storePath = os.Args[1]
	}

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		if err := os.Mkdir(storePath, os.ModePerm); err != nil {
			panic(err)
		}
	}

	router := gin.Default()
	gormdb, err := gorm.Open(sqlite.Open(path.Join(storePath, "backup.db")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	gormdb.AutoMigrate(&BackedFile{})

	router.POST("/upload/:id", func(c *gin.Context) {
		size, calcHash, err := receiveAndStoreFile(storePath, c.Param("id"), c)
		if err != nil {
			log.Fatal(err)
			c.AbortWithError(http.StatusBadRequest, err)
		}

		defer c.Request.Body.Close()

		modifiedOn, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", c.Request.FormValue("ModifiedOn"))
		gormdb.Save(&BackedFile{
			Path:       c.Request.FormValue("Path"),
			Hash:       calcHash,
			Size:       size,
			ModifiedOn: modifiedOn,
			UploadedOn: time.Now(),
		})

		c.String(http.StatusOK, "")
	})

	router.GET("/files/", func(c *gin.Context) {
		files := []BackedFile{}
		gormdb.Find(&files)
		c.JSON(http.StatusOK, files)
	})

	router.GET("/file/:id", func(c *gin.Context) {
		file := BackedFile{}
		log.Println("Requested file ", c.Param("id"))
		gormdb.Debug().Where(BackedFile{Hash: c.Param("id")}).First(&file)
		if file.Hash == "" {
			c.AbortWithError(http.StatusNotFound, nil)
		}

		filePath := path.Join(storePath, file.Hash)

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment")
		c.Header("Content-Type", "application/octet-stream")
		c.File(filePath)

	})

	router.Run(":8080")
}

const fileChunk = 8192

func receiveAndStoreFile(storePath string, id string, c *gin.Context) (int64, string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return 0, "", err
	}

	bodyFile, err := file.Open()
	if err != nil {
		log.Fatal(err)
		return 0, "", err
	}

	defer bodyFile.Close()

	blocks := uint64(math.Ceil(float64(file.Size) / float64(fileChunk)))
	savedFilePath := path.Join(storePath, id)
	savedFile, err := os.Create(savedFilePath)

	if err != nil {
		log.Fatal(err)
		return 0, "", err
	}

	log.Println("Saving file ", savedFilePath)

	defer savedFile.Close()
	calculatedSize := int64(0)
	calcHash := md5.New()
	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(fileChunk, float64(file.Size-int64(i*fileChunk))))
		buf := make([]byte, blocksize)
		bytesRead, err := bodyFile.Read(buf)

		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}

			break
		}

		calculatedSize += int64(bytesRead)
		io.WriteString(calcHash, string(buf))
		savedFile.Write(buf)
	}

	h := fmt.Sprintf("%x", calcHash.Sum(nil))
	if id != h {
		log.Fatal("Calculated hash: ", h, " Request hash: ", id)
	}

	return calculatedSize, fmt.Sprintf("%x", calcHash.Sum(nil)), nil
}
