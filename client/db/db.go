package db

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/manso/gobackup/client/structs"
)

type BackupDb struct {
	db *gorm.DB
}

func CreateDatabase(file string) *BackupDb {
	gormdb, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	gormdb.AutoMigrate(&structs.BackedFile{})
	return &BackupDb{gormdb}
}

func (bdb *BackupDb) Insert(hashedFile *structs.BackedFile) {
	bdb.db.Save(hashedFile)
}

func (bdb *BackupDb) InsertBulk(hashedFiles []structs.BackedFile) {
	bdb.db.Transaction(func(tx *gorm.DB) error {
		for _, f := range hashedFiles {
			bdb.Insert(&f)
		}

		return nil
	})
}

func (bdb *BackupDb) GetFilesNotUploaded() []structs.BackedFile {
	hashedFiles := []structs.BackedFile{}
	bdb.db.Where(map[string]interface{}{"uploaded": 0}).Find(&hashedFiles)
	return hashedFiles
}

func (bdb *BackupDb) Clear() {
	bdb.db.Exec("DELETE FROM hashed_files")
}

func (bdb *BackupDb) SetUploaded(hash string, size int64) {
	hashedFile := structs.BackedFile{}
	bdb.db.Where(&structs.BackedFile{Hash: hash}).First(&hashedFile)
	hashedFile.UploadedOn = time.Now()
	hashedFile.Size = size
	bdb.db.Save(&hashedFile)
}
