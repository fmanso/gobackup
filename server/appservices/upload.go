package appservices

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"time"

	"github.com/manso/gobackup/server/store"
)

func UploadFile(
	bfStore store.BackedFileStorer,
	storePath string,
	hash string,
	path string,
	modifiedOn time.Time,
	size int64,
	reader io.Reader) error {

	size, calcHash, err := receiveAndStoreFile(storePath, hash, size, reader)
	if err != nil {
		log.Fatal(err)
		return err
	}

	backedFile := bfStore.FindBackedFileByPath(path)
	if backedFile == nil {
		backedFile = &store.BackedFile{
			Path: path,
		}
	}

	backedFile.Versions = append(backedFile.Versions, store.BackedFileVersion{
		Hash:       calcHash,
		Size:       size,
		ModifiedOn: modifiedOn,
		UploadedOn: time.Now(),
	})

	bfStore.Save(backedFile)

	return nil
}

const fileChunk = 8192

func receiveAndStoreFile(storePath string, id string, size int64, reader io.Reader) (int64, string, error) {
	blocks := uint64(math.Ceil(float64(size) / float64(fileChunk)))
	savedFilePath := path.Join(storePath, id)
	savedFile, err := os.Create(savedFilePath)

	if err != nil {
		return 0, "", err
	}

	log.Println("Saving file ", savedFilePath)

	defer savedFile.Close()
	calculatedSize := int64(0)
	calcHash := md5.New()
	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(fileChunk, float64(size-int64(i*fileChunk))))
		buf := make([]byte, blocksize)
		bytesRead, err := reader.Read(buf)

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
