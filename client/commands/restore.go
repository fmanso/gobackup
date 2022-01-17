package commands

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	gpath "path"
	"time"

	"github.com/manso/gobackup/client/path"
	"github.com/manso/gobackup/client/terminal"
	"github.com/pterm/pterm"
	"gorm.io/gorm"
)

type BackedFile struct {
	gorm.Model
	Path       string
	Hash       string
	Size       int64
	ModifiedOn time.Time
	UploadedOn time.Time
}

func PerformRestore(rootPath string) {
	resp, err := http.Get("http://localhost:8080/files")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	files := []BackedFile{}
	err = json.Unmarshal(body, &files)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	size := int64(0)
	for _, f := range files {
		size += f.Size
	}

	// TODO: Refactor function Calculate Size Progress
	pterm.Println("Files: ", len(files))
	pterm.Println("Size: ", terminal.CalculateSizeProgress(size))
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(files)).WithTitle("Downloading files").Start()
	for _, f := range files {
		sanitizedFilePath := path.Sanitize(f.Path)
		pathWithoutFile := path.GetPathWithoutFile(sanitizedFilePath)
		pathWithoutFile = gpath.Join(rootPath, pathWithoutFile)
		path.EnsureDirectoryPathExists(pathWithoutFile)

		resp, err := http.Get("http://localhost:8080/file/" + f.Hash)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		out, err := os.Create(gpath.Join(pathWithoutFile, path.GetFileName(f.Path)))
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		progressBar.Increment()
	}

	progressBar.Stop()
}
