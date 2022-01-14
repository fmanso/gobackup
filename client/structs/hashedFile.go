package structs

import (
	"time"

	"gorm.io/gorm"
)

type BackedFile struct {
	gorm.Model
	Hash       string `gorm:uniqueIndex`
	Path       string
	Size       int64
	ModifiedOn time.Time
	CreatedOn  time.Time
	UploadedOn time.Time
}
