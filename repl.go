package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/kartikey-tiwari/pokedex-go/internal/pokeapi"
)

const HIST_FILENAME = ".pokedex_history"
const HIST_SIZE = 5

var historyIndex int
var currentInput string
var modifications map[int]string

func enableRawMode() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func restoreNormalTTYSettings() {
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	exec.Command("stty", "-F", "/dev/tty", "-cbreak").Run()
}

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

func cleanInput(text string) []string {
	cleanedText := strings.Fields(strings.ToLower(text))
	return cleanedText
}

func commandExit(c *pokeapi.Config, arg string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	if histFile != nil {
		if len(history) > HIST_SIZE {
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
		histFile.Close()
	}
	restoreNormalTTYSettings()
	os.Exit(0)
	return nil
}

func commandHelp(c *pokeapi.Config, arg string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, v := range commandRegistry {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMapMain(c *pokeapi.Config, next bool) error {
	if !next && c.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	locations, err := pokeapi.GetLocationAreaNames(c, next)
	if err != nil {
		return err
	}

	for _, area := range locations {
		fmt.Println(area)
	}
	return nil
}

func commandMap(c *pokeapi.Config, arg string) error {
	return commandMapMain(c, true)
}

func commandMapBack(c *pokeapi.Config, arg string) error {
	return commandMapMain(c, false)
}

func commandExplore(c *pokeapi.Config, arg string) error {
	if arg == "" {
		fmt.Println("No location provided")
		return nil
	}
	pokemons, err := pokeapi.GetPokemonsInArea(arg)
	if err != nil {
		return err
	}

	for _, v := range pokemons {
		fmt.Println(v)
	}
	return nil
}

func commandCatch(c *pokeapi.Config, arg string) error {
	if arg == "" {
		fmt.Println("No pokemon provided")
		return nil
	}
	fmt.Println("Throwing a Pokeball at " + arg + "...")
	pokemon, err := pokeapi.GetPokemonInformation(arg)
	if err != nil {
		return err
	}
	baseXp := pokemon.BaseExperience
	caught := rand.IntN(650)+1 > baseXp
	if caught {
		pokedex[arg] = pokemon
		fmt.Println(arg + " was caught!")
		fmt.Println("You may now inspect it with the inspect command")
	} else {
		fmt.Println(arg + " escaped!")
	}
	return nil
}

func commandInspect(c *pokeapi.Config, arg string) error {
	pokemon, ok := pokedex[arg]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for _, v := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", v.Stat.Name, v.BaseStat)
	}
	fmt.Println("Types:")
	for _, v := range pokemon.Types {
		fmt.Println("  -", v.Type.Name)
	}
	return nil
}

func commandPokedex(c *pokeapi.Config, arg string) error {
	if len(pokedex) == 0 {
		fmt.Println("You haven't caught any pokemons yet")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, v := range pokedex {
		fmt.Println(" -", v.Name)
	}
	return nil
}

func commandHistory(c *pokeapi.Config, option string) error {
	width := len(fmt.Sprintf("%d", len(history)))

	for i, v := range history {
		fmt.Printf("%*d. %s\n", width, i+1, v)
	}
	return nil
}

type CliCommand struct {
	name        string
	description string
	callback    func(conf *pokeapi.Config, arg string) error
}

var commandRegistry map[string]CliCommand = make(map[string]CliCommand)
var config *pokeapi.Config = &pokeapi.Config{
	Next:     "https://pokeapi.co/api/v2/location-area?limit=20&offset=0",
	Previous: "",
}
var pokedex map[string]pokeapi.PokemonResponse = map[string]pokeapi.PokemonResponse{}
var history []string = []string{}
var histFile *os.File

func initCommands() {
	commandRegistry["help"] = CliCommand{
		name:        "help",
		description: "Dispays a help message",
		callback:    commandHelp,
	}
	commandRegistry["exit"] = CliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	commandRegistry["map"] = CliCommand{
		name:        "map",
		description: "Search for next location areas",
		callback:    commandMap,
	}
	commandRegistry["mapb"] = CliCommand{
		name:        "mapb",
		description: "Search for previous location areas",
		callback:    commandMapBack,
	}
	commandRegistry["explore"] = CliCommand{
		name:        "explore",
		description: "Explore a location area",
		callback:    commandExplore,
	}
	commandRegistry["catch"] = CliCommand{
		name:        "catch",
		description: "Try to catch a pokemon",
		callback:    commandCatch,
	}
	commandRegistry["inspect"] = CliCommand{
		name:        "inspect",
		description: "Display caught pokemon information",
		callback:    commandInspect,
	}
	commandRegistry["pokedex"] = CliCommand{
		name:        "pokedex",
		description: "Display names of caught pokemons",
		callback:    commandPokedex,
	}
	commandRegistry["history"] = CliCommand{
		name:        "history",
		description: "Displays previous commands",
		callback:    commandHistory,
	}
}

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

func startREPL() {
	defer restoreNormalTTYSettings()
	enableRawMode()
	initCommands()
	loadHistory()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		commandExit(config, "")
	}()

	for {
		input, ok := readInput("Pokedex > ")
		if !ok {
			commandExit(config, "")
		}

		trimmedInput := strings.TrimSpace(input)
		history = append(history, trimmedInput)
		if histFile != nil {
			histFile.WriteString(trimmedInput + "\n")
		}

		cleanedInput := cleanInput(input)
		if len(cleanedInput) == 0 {
			continue
		}

		command, ok := commandRegistry[cleanedInput[0]]
		var option string
		if len(cleanedInput) == 2 {
			option = cleanedInput[1]
		} else {
			option = ""
		}
		if !ok {
			fmt.Println("Unknown command. Type help to see list of commands.")
			continue
		} else {
			err := command.callback(config, option)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
