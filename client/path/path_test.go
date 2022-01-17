package path

import (
	"os"
	"testing"
)

func TestSanitize(t *testing.T) {
	toSanitize := "f:\\users\\broker\\documentos/Finanzas/~$Crypto.xlsx"

	sanitized := Sanitize(toSanitize)

	good := "f/users/broker/documentos/Finanzas/~$Crypto.xlsx"
	if sanitized != good {
		t.Errorf("Expected %s but got %s", good, sanitized)
	}
}

func TestSplit(t *testing.T) {
	toSplit := "f/users/broker/documentos/Finanzas/~$Crypto.xlsx"

	good := []string{
		"f",
		"users",
		"broker",
		"documentos",
		"Finanzas",
		"~$Crypto.xlsx",
	}

	split := Split(toSplit)

	for i, s := range split {
		if good[i] != s {
			t.Errorf("Expected %s, but got %s", good[i], s)
		}
	}
}

func TestGetPathWithoutFile(t *testing.T) {
	toPathWithoutFile := "f/users/broker/documentos/Finanzas/~$Crypto.xlsx"
	good := "f/users/broker/documentos/Finanzas"
	pathWithoutFile := GetPathWithoutFile(toPathWithoutFile)

	if good != pathWithoutFile {
		t.Errorf("Expected %s, but got %s", good, pathWithoutFile)
	}
}

func TestEnsureDirectoryPathExists(t *testing.T) {
	path := "f/users/broker/common/"

	err := EnsureDirectoryPathExists(path)
	if err != nil {
		t.Errorf(err.Error())
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Path is not created")
	}

	path = "f/users/broker/common/newone"
	err = EnsureDirectoryPathExists(path)
	if err != nil {
		t.Errorf(err.Error())
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Path is not created")
	}
}

func TestGetFileName(t *testing.T) {
	path := "f:/asdas/dasd/filename.pdf"

	fileName := GetFileName(path)

	good := "filename.pdf"
	if good != fileName {
		t.Errorf("Expected %s, but got %s", good, fileName)
	}

}
