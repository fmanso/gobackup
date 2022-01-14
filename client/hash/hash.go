package hash

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"os"
)

const fileChunk = 8192

func CalculateHash(file string) (md5Hash string, err error) {
	f, e := os.Open(file)
	if e != nil {
		return "", e
	}

	defer f.Close()
	info, _ := f.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(fileChunk)))

	hash := md5.New()
	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(fileChunk, float64(filesize-int64(i*fileChunk))))
		buf := make([]byte, blocksize)
		f.Read(buf)
		io.WriteString(hash, string(buf))
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
