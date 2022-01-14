package terminal

import (
	"fmt"
	"time"

	"github.com/pterm/pterm"
)

type CalculatingBackupSize struct {
	lastUpdate time.Time
}

func (b *CalculatingBackupSize) Start() {
	b.lastUpdate = time.Now()
	pterm.Println("Calculating backup size...")
}

func (b *CalculatingBackupSize) Update(backupSize int64) {
	if time.Since(b.lastUpdate).Milliseconds() > 10 {
		text := fmt.Sprintf("Backup Size: %s", calculateSizeProgress(backupSize))
		pterm.Printo(text)
		b.lastUpdate = time.Now()
	}
}

func (b *CalculatingBackupSize) End(backupSize int64) {
	text := fmt.Sprintf("Backup Size: %s", calculateSizeProgress(backupSize))
	pterm.Printo(text)
	pterm.Println()
}

func calculateSizeProgress(backupSize int64) string {
	mb := backupSize >> 20
	if mb < 1000 {
		return fmt.Sprintf("%d MB", mb)
	}

	gb := mb >> 10
	return fmt.Sprintf("%d GB", gb)
}
