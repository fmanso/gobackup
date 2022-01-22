package store

import (
	"os"
	"testing"
	"time"
)

func TestBackedFileStoreShould(t *testing.T) {
	t.Run("be empty when created", func(t *testing.T) {
		os.Remove("test.db")
		db := Open("test.db")
		defer db.Close()

		if db.GetBackedFileCount() != 0 {
			t.Error("BackedFileStore is not empty")
		}
	})

	t.Run("insert backed files", func(t *testing.T) {
		os.Remove("test.db")
		db := Open("test.db")
		defer db.Close()

		insertTestBackedFile(db, "c:/test/path")

		backedFiles := db.GetBackedFiles()
		if len(backedFiles) != 1 {
			t.Errorf("Expected 1 backed files, Got %d", len(backedFiles))
		}

		if len(backedFiles[0].Versions) != 1 {
			t.Errorf("Expected 1 file versions, Got %d", len(backedFiles[0].Versions))
		}
	})

	t.Run("find by path", func(t *testing.T) {
		os.Remove("test.db")
		db := Open("test.db")
		defer db.Close()

		insertTestBackedFile(db, "c:/test/path")
		backedFile := db.FindBackedFileByPath("c:/test/path")

		if backedFile == nil {
			t.Errorf("Backed file not found")
		}
	})

	t.Run("return nil when finding by path if not found", func(t *testing.T) {
		os.Remove("test.db")
		db := Open("test.db")
		defer db.Close()

		insertTestBackedFile(db, "c:/test/path")

		backedFile := db.FindBackedFileByPath("c:/test/path2")

		if backedFile != nil {
			t.Errorf("Backed file should return nil")
		}
	})

	t.Run("find by hash", func(t *testing.T) {
		os.Remove("test.db")
		db := Open("test.db")
		defer db.Close()

		insertTestBackedFile(db, "c:/test/path")
		backedFileVersion := db.FindBackedFileVersionByHash("hash1")

		if backedFileVersion == nil {
			t.Errorf("Backed file version not found")
		}
	})

	t.Run("return nil when finding by hash if not found", func(t *testing.T) {
		os.Remove("test.db")
		db := Open("test.db")
		defer db.Close()

		insertTestBackedFile(db, "c:/test/path")
		backedFileVersion := db.FindBackedFileVersionByHash("hash2")

		if backedFileVersion != nil {
			t.Errorf("Backed file should return nil")
		}
	})

	t.Run("return backed files count", func(t *testing.T) {
		os.Remove("test.db")
		db := Open("test.db")
		defer db.Close()

		insertTestBackedFile(db, "c:/test/path")
		count := db.GetBackedFileCount()

		if count != 1 {
			t.Errorf("Backed file count should be 1")
		}
	})
}

func insertTestBackedFile(db *BackedFileStore, path string) {
	backedFile := BackedFile{
		Path:     path,
		Versions: []BackedFileVersion{},
	}

	backedFile.Versions = append(backedFile.Versions, BackedFileVersion{
		Hash:       "hash1",
		ModifiedOn: time.Now(),
		Size:       int64(260),
	})

	db.Save(&backedFile)
}
