# Pokedex Go

A command-line Pokedex application written in Go that interacts with the [PokeAPI](https://pokeapi.co/). This tool allows you to explore the Pokemon world, discover location areas, and catch Pokemon directly from your terminal.

## Features

- **Interactive REPL**: A robust Read-Eval-Print Loop (REPL) environment to interact with the application.
- **Caching**: Implements a custom caching system to store API responses and reduce network calls, improving performance.
- **Command History**: Persists command history to `.pokedex_history` in your home directory, allowing you to navigate previous commands using Up/Down arrow keys.
- **Game Mechanics**: Includes a catching mechanic where success is based on the Pokemon's base experience level.

## Installation

Ensure you have Go installed (version 1.24.4 or higher is recommended).

1.  **Clone the repository:**

    ```bash
    git clone [https://github.com/kartikey-tiwari/pokedex-go.git](https://github.com/kartikey-tiwari/pokedex-go.git)
    cd pokedex-go
    ```

2.  **Build and run the application:**

    ```bash
    go build -o pokedex
    ./pokedex
    ```

    Or run directly without building:

    ```bash
    go run .
    ```

## Usage

Once the application is running, you will see the `Pokedex >` prompt. You can interact with the Pokedex using the commands listed below.

### Commands

- `help`: Displays a help message describing available commands.
- `exit`: Exits the Pokedex application.
- `map`: Displays the next 20 location areas in the Pokemon world.
- `mapb`: Displays the previous 20 location areas.
- `explore <area_name>`: Lists all Pokemon found in a specific location area.
  - _Example:_ `explore pastoria-city-area`
- `catch <pokemon_name>`: Attempts to catch a specific Pokemon. Catching is probabilistic; harder Pokemon are more difficult to catch.
- `inspect <pokemon_name>`: View details (height, weight, stats, types) of a Pokemon you have successfully caught.
- `pokedex`: Lists the names of all Pokemon you have caught so far.
- `history`: Displays a list of your previously executed commands.

## Development

### Project Structure

- `main.go`: Entry point of the application.
- `repl.go`: Handles the REPL loop, command parsing, history management, and raw terminal mode.
- `internal/pokeapi/`: Client logic for interacting with the PokeAPI, including data types and fetching functions.
- `internal/pokecache/`: A custom implementation of a cache with time-to-live (TTL) eviction.
