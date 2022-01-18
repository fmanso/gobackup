package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BackedFile struct {
	gorm.Model
	Path     string
	Versions []BackedFileVersion
}

type BackedFileVersion struct {
	gorm.Model
	BackedFileID uint
	Hash         string
	Size         int64
	ModifiedOn   time.Time
	UploadedOn   time.Time
}

type BackedFilesResponse struct {
	Path       string
	Hash       string
	Size       int64
	ModifiedOn time.Time
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
	gormdb.AutoMigrate(&BackedFileVersion{})

	router.POST("/upload/:id", func(c *gin.Context) {
		size, calcHash, err := receiveAndStoreFile(storePath, c.Param("id"), c)
		if err != nil {
			log.Fatal(err)
			c.AbortWithError(http.StatusBadRequest, err)
		}

		defer c.Request.Body.Close()

		modifiedOn, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", c.Request.FormValue("ModifiedOn"))

		path := c.Request.FormValue("Path")

		backedFile := BackedFile{}
		result := gormdb.Where(&BackedFile{Path: path}).First(&backedFile)
		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			backedFile.Path = path
		}

		backedFile.Versions = append(backedFile.Versions, BackedFileVersion{
			Hash:       calcHash,
			Size:       size,
			ModifiedOn: modifiedOn,
			UploadedOn: time.Now(),
		})

		gormdb.Save(&backedFile)
		c.String(http.StatusOK, "")
	})

	router.GET("/files/", func(c *gin.Context) {
		files := []BackedFile{}
		gormdb.Preload("Versions").Find(&files)

		response := []BackedFilesResponse{}

		for _, f := range files {
			v := getLatesVersion(f)
			response = append(response, BackedFilesResponse{
				Hash:       v.Hash,
				Size:       v.Size,
				Path:       f.Path,
				ModifiedOn: v.ModifiedOn,
			})
		}

		c.JSON(http.StatusOK, response)
	})

	router.GET("/file/:hash", func(c *gin.Context) {
		fileVersion := BackedFileVersion{}
		log.Println("Requested hash ", c.Param("hash"))
		result := gormdb.Where(BackedFileVersion{Hash: c.Param("hash")}).First(&fileVersion)
		if result.Error != nil {
			c.AbortWithError(http.StatusNotFound, nil)
		}

		filePath := path.Join(storePath, fileVersion.Hash)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment")
		c.Header("Content-Type", "application/octet-stream")
		c.File(filePath)

	})

	router.Run(":8080")
}

func getLatesVersion(f BackedFile) BackedFileVersion {
	sort.SliceStable(f.Versions, func(i, j int) bool {
		return f.Versions[i].ModifiedOn.After(f.Versions[j].ModifiedOn)
	})

	return f.Versions[0]
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
