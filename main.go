package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	cleanedText := strings.Fields(strings.ToLower(text))
	return cleanedText
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println(`Welcome to the Pokedex!
Usage:
`)
	for _, v := range commandRegistry {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commandRegistry map[string]cliCommand = make(map[string]cliCommand)

func initCommands() {
	commandRegistry["help"] = cliCommand{
		name:        "help",
		description: "Dispays a help message",
		callback:    commandHelp,
	}
	commandRegistry["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
}

func main() {
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
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
