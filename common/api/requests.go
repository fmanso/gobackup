package api

import "time"

type UploadFileRequest struct {
	Path       string
	CreatedOn  time.Time
	ModifiedOn time.Time
}
