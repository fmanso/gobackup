package appservices

import (
	"testing"
	"time"

	"github.com/manso/gobackup/server/store"
)

func TestGetFilesShould(t *testing.T) {
	t.Run("get files", func(t *testing.T) {
		backedFileStorer := backedFileStoreFake{
			getBackedFiles: func() []store.BackedFile {
				return []store.BackedFile{
					{
						Versions: []store.BackedFileVersion{
							{
								ModifiedOn: time.Now(),
							},
						},
					},
					{
						Versions: []store.BackedFileVersion{
							{
								ModifiedOn: time.Now(),
							},
						},
					},
				}
			},
		}

		result := GetBackedFiles(&backedFileStorer)

		if len(result) != 2 {
			t.Error("Expected 2, but got", len(result))
		}
	})
}
