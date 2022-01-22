package store

import (
	"errors"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BackedFile struct {
	gorm.Model
	Path     string
	Versions []BackedFileVersion
}

type BackedFileVersion struct {
	gorm.Model
	BackedFileID uint
	Hash         string
	Size         int64
	ModifiedOn   time.Time
	UploadedOn   time.Time
}

type BackedFileStore struct {
	db *gorm.DB
}

type BackedFileStorer interface {
	Save(backedFile *BackedFile) error
	FindBackedFileByPath(path string) *BackedFile
	GetBackedFiles() []BackedFile
	FindBackedFileVersionByHash(hash string) *BackedFileVersion
}

func Open(path string) *BackedFileStore {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&BackedFile{})
	db.AutoMigrate(&BackedFileVersion{})

	store := BackedFileStore{}
	store.db = db
	return &store
}

func (s *BackedFileStore) GetBackedFiles() []BackedFile {
	backedFiles := []BackedFile{}
	result := s.db.Preload("Versions").Find(&backedFiles)
	if result.Error != nil {
		panic(result.Error)
	}

	return backedFiles
}

func (s *BackedFileStore) FindBackedFileByPath(path string) *BackedFile {
	backedFile := BackedFile{}
	result := s.db.Where(&BackedFile{Path: path}).First(&backedFile)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &backedFile
}

func (s *BackedFileStore) FindBackedFileVersionByHash(hash string) *BackedFileVersion {
	backedFileVersion := BackedFileVersion{}
	result := s.db.Where(&BackedFileVersion{Hash: hash}).First(&backedFileVersion)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &backedFileVersion
}

func (s *BackedFileStore) GetBackedFileCount() int64 {
	var count int64
	s.db.Model(&BackedFile{}).Count(&count)
	return count
}

func (s *BackedFileStore) Save(backedFile *BackedFile) error {
	result := s.db.Save(backedFile)
	return result.Error
}

func (s *BackedFileStore) Close() {
	db, err := s.db.DB()
	if err != nil {
		panic(err)
	}

	db.Close()
}
