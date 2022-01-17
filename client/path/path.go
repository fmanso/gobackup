package path

import (
	"os"
	"path"
	"strings"
)

func Sanitize(p string) string {
	sanitized := strings.ReplaceAll(p, ":", "")
	sanitized = strings.ReplaceAll(sanitized, "\\", "/")
	return sanitized
}

func Split(p string) []string {
	return strings.Split(p, "/")
}

func GetPathWithoutFile(p string) string {
	split := Split(p)
	return strings.Join(split[:len(split)-1], "/")
}

func GetFileName(p string) string {
	split := Split(p)
	return split[len(split)-1]
}

func EnsureDirectoryPathExists(p string) error {
	sanitized := strings.ReplaceAll(p, "\\", "/")
	split := Split(sanitized)
	current := ""
	for _, s := range split {
		current = path.Join(current, s)
		if _, err := os.Stat(current); os.IsNotExist(err) {
			if err := os.Mkdir(current, os.ModePerm); err != nil {
				return err
			}
		}
	}

	return nil
}
