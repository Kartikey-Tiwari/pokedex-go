package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/kartikey-tiwari/pokedex-go/internal/pokeapi"
)

const HIST_FILENAME = ".pokedex_history"
const HIST_SIZE = 1000

func cleanInput(text string) []string {
	cleanedText := strings.Fields(strings.ToLower(text))
	return cleanedText
}

func commandExit(c *pokeapi.Config, arg string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	histFile.Close()
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

	histFile, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		histFile = nil
		return
	}

	scanner := bufio.NewScanner(histFile)
	for scanner.Scan() {
		text := scanner.Text()
		history = append(history, text)
	}
}

func startREPL() {
	initCommands()
	loadHistory()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			fmt.Println()
			commandExit(config, "")
		}
		input := scanner.Text()

		trimmedInput := strings.TrimSpace(input) + "\n"
		history = append(history, trimmedInput)
		if histFile != nil {
			histFile.WriteString(trimmedInput)
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
