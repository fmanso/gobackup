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
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	r, _ := http.NewRequest("POST", "http://localhost:8080/upload/"+hashedFile.Hash, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()
}
