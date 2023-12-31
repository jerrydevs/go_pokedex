package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"pokeapi"
	"strings"
)

const (
	BASE_ESCAPE_RATE    = 0.2
	MAX_BASE_EXPERIENCE = 700
)

func printPrompt() {
	fmt.Print("Pokedex> ")
}

func printSpaces() {
	fmt.Println()
	fmt.Println()
	fmt.Println()
}

type commandConfig struct {
	offset   int `default:"0"`
	pageSize int `default:"20"`
	client   pokeapi.PokeApi
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *commandConfig, arg string, state *State)
}

func commandHelp(c *commandConfig, arg string, state *State) {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: <command> [arguments]")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the next 20 locations in the Pokedex, starting from the current page")
	fmt.Println("mapb: Go back to the previous 20 locations in the Pokedex")
	fmt.Println("explore: List all the pokemon in a given location")
	fmt.Println()
}

func commandExit(c *commandConfig, arg string, state *State) {
	os.Exit(0)
}

func commandMap(c *commandConfig, arg string, state *State) {
	regionResult := c.client.GetMap(c.pageSize, c.offset)
	for _, region := range regionResult.Results {
		fmt.Println(region.Name)
	}

	c.offset += c.pageSize
}

func commandMapBack(c *commandConfig, arg string, state *State) {
	if c.offset == 0 {
		fmt.Println("You are already at the beginning of the list")
		return
	}

	c.offset -= c.pageSize
	regionResult := c.client.GetMap(c.pageSize, c.offset)
	for _, region := range regionResult.Results {
		fmt.Println(region.Name)
	}
}

func commandExplore(c *commandConfig, arg string, state *State) {
	if arg == "" {
		fmt.Println("Provide a location to explore")
		return
	}

	exploreResult, err := c.client.ExploreLocation(arg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, pokemonObj := range exploreResult.PokemonEncounters {
		fmt.Println(pokemonObj.Pokemon.Name)
	}
}

func commandCatch(c *commandConfig, arg string, state *State) {
	if arg == "" {
		fmt.Println("Provide a pokemon to catch")
		return
	}

	pokemon, err := c.client.GetPokemon(arg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	captureThreshold := (float64(pokemon.BaseExperience) / float64(MAX_BASE_EXPERIENCE))
	captureAttempt := rand.Float64()
	fmt.Printf("Pokemon BaseExp: %d, Capture Threshold: %.2f, captureAttempt: %.2f\n", pokemon.BaseExperience, captureThreshold, captureAttempt)

	if captureAttempt >= captureThreshold {
		fmt.Printf("You caught %s! Added to Pokedex.\n", pokemon.Name)
		state.CaughtPokemon = append(state.CaughtPokemon, pokemon)
	} else {
		fmt.Printf("Oh no! The %s escaped!\n", pokemon.Name)
	}
}

func commandInspect(c *commandConfig, arg string, state *State) {
	if arg == "" {
		fmt.Printf("You have %d pokemon in your pokedex\n", len(state.CaughtPokemon))
		return
	}

	hasPokemon := false
	for _, pokemon := range state.CaughtPokemon {
		if pokemon.Name == arg {
			hasPokemon = true
			pokemon.PrintInfo()
		}
	}

	if !hasPokemon {
		fmt.Println("You haven't caught that pokemon yet!")
		return
	}

}

func commandPokedex(c *commandConfig, arg string, state *State) {
	if len(state.CaughtPokemon) == 0 {
		fmt.Println("You haven't caught any pokemon yet!")
		return
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range state.CaughtPokemon {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays a list of all regions in the pokedex",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous list of regions",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "List all the pokemon in a given location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon in your pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the pokemon you've caught",
			callback:    commandPokedex,
		},
	}
}

type State struct {
	CaughtPokemon []pokeapi.Pokemon
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	pageSize := 20
	pokeApi := pokeapi.NewPokeApi()
	config := &commandConfig{pageSize: pageSize, offset: 0, client: pokeApi}
	state := &State{}

	for {
		printPrompt()
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err)
		}

		input := scanner.Text()
		command := strings.Split(input, " ")[0]
		arg := ""
		if len(strings.Split(input, " ")) > 1 {
			arg = strings.Split(input, " ")[1]
		}

		fmt.Printf("command: %s, arg: %s\n", command, arg)

		cliObj, ok := commands[command]

		if !ok {
			fmt.Println("Command not found")
			continue
		}

		cliObj.callback(config, arg, state)
		printSpaces()
	}

}
