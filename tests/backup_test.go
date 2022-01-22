package backuptests

import (
	"log"
	gpath "path"
	"path/filepath"
	"testing"

	"github.com/manso/gobackup/client/commands"
	"github.com/manso/gobackup/client/hash"
	"github.com/manso/gobackup/client/path"
	"github.com/manso/gobackup/server/http"
	"github.com/manso/gobackup/server/store"
)

func TestBackup(t *testing.T) {
	storePath, err := filepath.Abs("./store")
	if err != nil {
		panic(err)
	}

	bfStore := store.Open(gpath.Join(storePath, "gobackup.db"))
	log.Println("Backed up files: ", bfStore.GetBackedFileCount())
	go http.Start(bfStore, storePath)
	commands.PerformBackup("./data/")
	commands.PerformRestore("./restore")

	dataFile, err := filepath.Abs("./data/download.jpg")
	if err != nil {
		t.Error(err)
	}

	hash1, err := hash.CalculateHash(dataFile)
	if err != nil {
		t.Error(err)
	}

	sanitizedFilePath := path.Sanitize(dataFile)
	pathWithoutFile := path.GetPathWithoutFile(sanitizedFilePath)
	pathWithoutFile = gpath.Join("./restore/", pathWithoutFile)
	filePath := gpath.Join(pathWithoutFile, "download.jpg")
	hash2, err := hash.CalculateHash(filePath)
	if err != nil {
		t.Error(filePath)
		t.Error(err)
	}

	if hash1 != hash2 {
		t.Error("Restored file is different from backed up file")
	}

}
