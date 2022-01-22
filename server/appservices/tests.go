package appservices

import "github.com/manso/gobackup/server/store"

type backedFileStoreFake struct {
	save                        func(backedFile *store.BackedFile) error
	findBackedFileByPath        func(path string) *store.BackedFile
	getBackedFiles              func() []store.BackedFile
	findBackedFileVersionByHash func(hash string) *store.BackedFileVersion
}

func (b *backedFileStoreFake) Save(backedFile *store.BackedFile) error {
	return b.save(backedFile)
}

func (b *backedFileStoreFake) FindBackedFileByPath(path string) *store.BackedFile {
	return b.findBackedFileByPath(path)
}

func (b *backedFileStoreFake) GetBackedFiles() []store.BackedFile {
	return b.getBackedFiles()
}

func (b *backedFileStoreFake) FindBackedFileVersionByHash(hash string) *store.BackedFileVersion {
	return b.findBackedFileVersionByHash(hash)
}

type ioReader struct {
}

func (io *ioReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 1
	}

	return 10, nil
}
