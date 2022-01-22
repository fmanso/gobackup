package appservices

import (
	"sort"
	"time"

	"github.com/manso/gobackup/server/store"
)

type GetBackedFilesResponse struct {
	Path       string
	Hash       string
	Size       int64
	ModifiedOn time.Time
}

func GetBackedFiles(bfStore store.BackedFileStorer) []GetBackedFilesResponse {
	files := bfStore.GetBackedFiles()
	response := []GetBackedFilesResponse{}

	for _, f := range files {
		v := getLatesVersion(&f)
		response = append(response, GetBackedFilesResponse{
			Hash:       v.Hash,
			Size:       v.Size,
			Path:       f.Path,
			ModifiedOn: v.ModifiedOn,
		})
	}

	return response
}

func getLatesVersion(f *store.BackedFile) *store.BackedFileVersion {
	sort.SliceStable(f.Versions, func(i, j int) bool {
		return f.Versions[i].ModifiedOn.After(f.Versions[j].ModifiedOn)
	})

	return &f.Versions[0]
}
