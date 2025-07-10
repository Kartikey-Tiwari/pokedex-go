package main

import (
	"bufio"
	"fmt"
	"github.com/kartikey-tiwari/pokedex-go/internal/pokeapi"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	cleanedText := strings.Fields(strings.ToLower(text))
	return cleanedText
}

func commandExit(c *pokeapi.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *pokeapi.Config) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:
`)
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

func commandMap(c *pokeapi.Config) error {
	return commandMapMain(c, true)
}

func commandMapBack(c *pokeapi.Config) error {
	return commandMapMain(c, false)
}

type CliCommand struct {
	name        string
	description string
	callback    func(conf *pokeapi.Config) error
}

var commandRegistry map[string]CliCommand = make(map[string]CliCommand)
var config *pokeapi.Config = &pokeapi.Config{
	Next:     "https://pokeapi.co/api/v2/location-area?limit=20&offset=0",
	Previous: "",
}

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
}

func startREPL() {
	initCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()

		cleanedInput := cleanInput(input)
		if len(cleanedInput) == 0 {
			continue
		}

		command, ok := commandRegistry[cleanedInput[0]]
		if !ok {
			fmt.Println("Unknown command. Type help to see usage")
			continue
		} else {
			err := command.callback(config)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
