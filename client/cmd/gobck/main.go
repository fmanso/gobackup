package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/devfacet/gocmd/v3"
	"github.com/manso/gobackup/client/api"
	"github.com/manso/gobackup/client/db"
	hashCalc "github.com/manso/gobackup/client/hash"
	"github.com/manso/gobackup/client/queue"
	"github.com/manso/gobackup/client/structs"
	"github.com/manso/gobackup/client/terminal"
	"github.com/pterm/pterm"
)

func main() {
	flags := struct {
		Help   bool `short:"h" long:"help" description:"Display usage" global:"true"`
		Backup struct {
			Path string `short:"p" long:"path" required:"true" description:"Root directory to be backed up"`
		} `command:"backup" description:"Perform backup"`
	}{}

	gocmd.HandleFlag("Backup", func(cmd *gocmd.Cmd, args []string) error {
		performBackup(flags.Backup.Path)
		return nil
	})

	gocmd.New(gocmd.Options{
		Name:        "gobck",
		Description: "Backup tool",
		Flags:       &flags,
		ConfigType:  gocmd.ConfigTypeAuto,
	})
}

func performBackup(path string) {
	rootPath := path
	start := time.Now()

	fileList := calculateBackupSize(rootPath)
	batch1, batch2 := hashFiles(fileList)
	files := append(batch1, batch2...)

	pterm.Println("Saving data to local database...")
	database := db.CreateDatabase("client.db")
	database.InsertBulk(files)
	pterm.Success.Println("Data saved")

	pterm.Println("Uploading files...")
	api.UploadFiles(files)
	pterm.Success.Println("Files uploaded")

	log.Println("Files found: ", len(fileList))
	log.Println("Uploaded files: ", len(files))
	elapsed := time.Since(start)
	log.Printf("Process took: %s", elapsed)
}

func calculateBackupSize(rootPath string) []string {
	queue := queue.Queue{}
	queue.Push(rootPath)

	fileList := make([]string, 0)
	backupSize := int64(0)
	calculatingBackupSize := terminal.CalculatingBackupSize{}
	calculatingBackupSize.Start()
	for !queue.Empty() {
		currentPath := queue.Pop()
		files, err := ioutil.ReadDir(currentPath)
		for _, file := range files {
			fileFullPath := path.Join(currentPath, file.Name())
			backupSize += file.Size()
			calculatingBackupSize.Update(backupSize)
			if file.IsDir() {
				queue.Push(fileFullPath)
			} else {
				fileList = append(fileList, fileFullPath)
			}
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	calculatingBackupSize.End(backupSize)

	return fileList
}

func hashFiles(fileList []string) (batch1 []structs.BackedFile, batch2 []structs.BackedFile) {
	c := make(chan []structs.BackedFile)
	p := make(chan int64)
	sizeHashed := int64(0)

	go hash(fileList[:len(fileList)/2], c, p)
	go hash(fileList[len(fileList)/2:], c, p)

	go func() {
		hashingFiles := terminal.HashingFiles{}
		hashingFiles.Start(len(fileList))
		for size := range p {
			sizeHashed = sizeHashed + size
			hashingFiles.Increment()
		}
	}()

	b1, b2 := <-c, <-c
	close(p)

	return b1, b2
}

func hash(files []string, c chan []structs.BackedFile, p chan int64) {
	hashedFiles := make([]structs.BackedFile, 0)

	for _, f := range files {
		fileInfo, err := os.Stat(f)
		if err != nil {
			log.Fatal(err)
			continue
		}

		hash, err := hashCalc.CalculateHash(f)
		if err != nil {
			log.Fatal(err)
			continue
		}

		hashedFile := structs.BackedFile{
			Path: f,
			Hash: hash,
		}

		hashedFiles = append(hashedFiles, hashedFile)
		p <- fileInfo.Size()
	}

	c <- hashedFiles
}
