package appservices

import (
	"testing"

	"github.com/manso/gobackup/server/store"
)

func TestDownloadShould(t *testing.T) {
	t.Run("download file version", func(t *testing.T) {
		backedFileStorer := backedFileStoreFake{
			findBackedFileVersionByHash: func(hash string) *store.BackedFileVersion {
				return &store.BackedFileVersion{
					Hash: "hash1",
				}
			},
		}

		filePath, err := DownloadBackedFileVersion(&backedFileStorer, "c:/path/test/", "hash1")

		if err != nil {
			t.Error("Error found ", err)
		}

		if filePath != "c:/path/test/hash1" {
			t.Errorf("Expected %s, got %s", "c:/path/test/hash1", filePath)
		}
	})
}
