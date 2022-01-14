package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manso/gobackup/common/api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BackedFile struct {
	gorm.Model
	Path       string
	Hash       string
	Size       int64
	CreatedOn  time.Time
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

	router.POST("/upload/:id", func(c *gin.Context) {
		var calcHash *string
		var size *int64
		err := receiveAndStoreFile(storePath, c.Param("id"), calcHash, size, c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		body, _ := ioutil.ReadAll(c.Request.Body)
		requestPayload := api.UploadFileRequest{}
		json.Unmarshal(body, &requestPayload)
		gormdb.Save(BackedFile{
			Path:       requestPayload.Path,
			Hash:       *calcHash,
			Size:       *size,
			CreatedOn:  requestPayload.CreatedOn,
			ModifiedOn: requestPayload.ModifiedOn,
			UploadedOn: time.Now(),
		})
	})

	router.Run(":8080")
}

func receiveAndStoreFile(storePath string, id string, hash *string, size *int64, c *gin.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	bodyFile, err := file.Open()
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024*16)
	savedFile, err := os.Create(path.Join(storePath, id))

	if err != nil {
		return err
	}

	defer savedFile.Close()
	*size = 0
	calcHash := md5.New()
	for {
		bytesRead, err := bodyFile.Read(buffer)

		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
				*hash = fmt.Sprintf("%x", calcHash.Sum(nil))
			}

			break
		}

		*size += int64(bytesRead)
		io.WriteString(calcHash, string(buffer))
		savedFile.Write(buffer[:bytesRead])
	}

	return nil
}
