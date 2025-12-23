package repl

import (
	"fmt"
	"os"
	"strings"
)

func readInput(prompt string) (string, bool) {
	fmt.Print(prompt)

	var input []byte
	historyIndex = len(history)
	currentInput = ""
	if modifications == nil {
		modifications = make(map[int]string)
	}

	for {
		char := make([]byte, 1)
		n, err := os.Stdin.Read(char)
		if n == 0 || err != nil { // EOF
			fmt.Println()
			return "", false
		}

		switch char[0] {
		case 10, 13: // Enter
			{
				fmt.Println()
				modifications = make(map[int]string)
				return string(input), true
			}
		case 27: // Escape
			{
				seq := make([]byte, 2)
				os.Stdin.Read(seq)
				if seq[0] == 91 {
					switch seq[1] {
					case 65: // Up arrow
						{
							if historyIndex == len(history) {
								currentInput = string(input)
							} else {
								modifications[historyIndex] = string(input)
							}

							if historyIndex > 0 {
								historyIndex--
								if modifiedCommand, exists := modifications[historyIndex]; exists {
									input = []byte(modifiedCommand)
								} else {
									histEntry := strings.TrimSpace(history[historyIndex])
									input = []byte(histEntry)
								}
								fmt.Print("\r\033[K" + prompt + string(input))
							}
						}
					case 66: // down arrow
						{
							if historyIndex < len(history)-1 {
								modifications[historyIndex] = string(input)
								historyIndex++

								if modifiedCmd, exists := modifications[historyIndex]; exists {
									input = []byte(modifiedCmd)
								} else {
									histEntry := strings.TrimSpace(history[historyIndex])
									input = []byte(histEntry)
								}
								fmt.Print("\r\033[K" + prompt + string(input))
							} else if historyIndex == len(history)-1 {
								modifications[historyIndex] = string(input)
								historyIndex = len(history)
								input = []byte(currentInput)
								fmt.Print("\r\033[K" + prompt + string(input))
							}
						}
					}
				}
			}
		case 127: // Backspace
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Print("\b \b")
			}
		case 3, 4: // SIGINT, EOF
			commandExit(config, "")

		default:
			if char[0] >= 32 && char[0] <= 126 {
				input = append(input, char[0])
				fmt.Print(string(char[0]))
			}
		}
	}
}
