package api

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/manso/gobackup/client/structs"
	"github.com/pterm/pterm"
)

func UploadFiles(hashedFiles []structs.BackedFile) {
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(hashedFiles)).WithTitle("Uploading files").Start()

	for _, hashedFile := range hashedFiles {
		uploadFile(hashedFile)
		progressBar.Increment()
	}

	progressBar.Stop()
}

func uploadFile(hashedFile structs.BackedFile) {
	file, _ := os.Open(hashedFile.Path)
	fileInfo, _ := file.Stat()
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("Path", hashedFile.Path)
	writer.WriteField("ModifiedOn", fileInfo.ModTime().String())
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	r, err := http.NewRequest("POST", "http://localhost:8080/upload/"+hashedFile.Hash, body)

	if err != nil {
		log.Fatal(err)
		return
	}

	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("StatusCode ", resp.StatusCode)
		return
	}
}
