package repl

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

const HIST_FILENAME = ".pokedex_history"
const HIST_SIZE = 100

var historyIndex int
var currentInput string
var modifications map[int]string

func loadHistory() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	filePath := filepath.Join(home, HIST_FILENAME)

	histFile, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		histFile = nil
		return
	}

	histFile.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(histFile)
	for scanner.Scan() {
		text := scanner.Text()
		history = append(history, text)
	}
}

func updateAndTruncateHistory() {
	histFile.Seek(0, io.SeekStart)
	history = history[len(history)-HIST_SIZE:]
	var totalBytes int64 = 0
	didBreak := false
	for _, entry := range history {
		n, err := histFile.WriteString(entry + "\n")
		if err != nil {
			didBreak = true
		}
		totalBytes += int64(n)
	}
	if !didBreak {
		histFile.Truncate(totalBytes)
	}
}
