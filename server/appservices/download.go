package appservices

import (
	"errors"
	"path"

	"github.com/manso/gobackup/server/store"
)

func DownloadBackedFileVersion(bfStore store.BackedFileStorer, storePath string, hash string) (string, error) {
	fileVersion := bfStore.FindBackedFileVersionByHash(hash)
	if fileVersion == nil {
		return "", errors.New("backed file version not found")
	}

	filePath := path.Join(storePath, fileVersion.Hash)
	return filePath, nil
}
