package appservices

import (
	"os"
	"testing"
	"time"

	"github.com/manso/gobackup/server/store"
)

func TestUploadShould(t *testing.T) {
	t.Run("upload file", func(t *testing.T) {
		saveCalled := false
		findBackedFileByPathArg := ""
		backedFileStorer := backedFileStoreFake{
			save: func(b *store.BackedFile) error {
				saveCalled = true
				return nil
			},

			findBackedFileByPath: func(path string) *store.BackedFile {
				findBackedFileByPathArg = path
				return nil
			},
		}
		storePath := os.TempDir()
		hash := "8afef86de3b0eff1ded7591aa5ff2769"
		path := "c:/path/test"
		modifiedOn := time.Now()
		size := int64(8)
		reader := ioReader{}
		err := UploadFile(&backedFileStorer, storePath, hash, path, modifiedOn, size, &reader)

		if err != nil {
			t.Errorf(err.Error())
		}

		if findBackedFileByPathArg != path {
			t.Errorf("Path argument not used to find backed file")
		}

		if !saveCalled {
			t.Errorf("Saved was not called")
		}

	})
}
