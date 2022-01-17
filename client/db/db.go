package db

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/manso/gobackup/client/structs"
)

type BackupDb struct {
	db *gorm.DB
}

func Open(file string) *BackupDb {
	gormdb, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	gormdb.AutoMigrate(&structs.BackedFile{})
	return &BackupDb{gormdb}
}

func (bdb *BackupDb) GetBackedFiles() []structs.BackedFile {
	backedFiles := []structs.BackedFile{}
	bdb.db.Find(&backedFiles)
	return backedFiles
}

func (bdb *BackupDb) Insert(hashedFile *structs.BackedFile) {
	bdb.db.Save(hashedFile)
}

func (bdb *BackupDb) InsertBulk(hashedFiles []structs.BackedFile) {
	bdb.db.Transaction(func(tx *gorm.DB) error {
		for _, f := range hashedFiles {
			result := bdb.db.Where(map[string]interface{}{"Path": f.Path, "Hash": f.Hash}).First(&structs.BackedFile{})
			if result != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
				bdb.Insert(&f)
			}
		}

		return nil
	})
}

func (bdb *BackupDb) SetUploaded(hash string, size int64) {
	hashedFile := structs.BackedFile{}
	bdb.db.Where(&structs.BackedFile{Hash: hash}).First(&hashedFile)
	hashedFile.UploadedOn = time.Now()
	hashedFile.Size = size
	bdb.db.Save(&hashedFile)
}
