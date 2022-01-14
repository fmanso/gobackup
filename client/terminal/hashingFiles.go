package terminal

import "github.com/pterm/pterm"

type HashingFiles struct {
	progressBar *pterm.ProgressbarPrinter
}

func (h *HashingFiles) Start(total int) {
	progressBar, _ := pterm.DefaultProgressbar.WithTotal(total).WithTitle("Hashing files").Start()
	h.progressBar = progressBar
}

func (h *HashingFiles) Increment() {
	h.progressBar.Increment()
}
