package main

import (
	"log"
	"os"
	"path"

	"github.com/manso/gobackup/server/http"
	"github.com/manso/gobackup/server/store"
)

func main() {
	storePath := "./store"
	if len(os.Args) > 1 {
		storePath = os.Args[1]
	}

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		if err := os.Mkdir(storePath, os.ModePerm); err != nil {
			panic(err)
		}
	}

	bfStore := store.Open(path.Join(storePath, "gobackup.go"))

	log.Println("Backed up files: ", bfStore.GetBackedFileCount())

	http.Start(bfStore, storePath)
}
