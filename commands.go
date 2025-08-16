package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/Quak1/pokedex-go/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, args ...string) error
}

type Config struct {
	Next          *string
	Previous      *string
	Client        pokeapi.Client
	CaughtPokemon map[string]pokeapi.Pokemon
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Get the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map",
			description: "Get the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "Get pokemon encounters for a given area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Catch a Pokemon!",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "View a Pokemon's details",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all your caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func commandExit(config *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, args ...string) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println()

	return nil
}

func commandMap(config *Config, args ...string) error {
	areas, err := config.Client.GetLocationAreas(config.Next)
	if err != nil {
		return err
	}

	config.Next = areas.Next
	config.Previous = areas.Previous

	for _, area := range areas.Results {
		println(area.Name)
	}

	return nil
}

func commandMapb(config *Config, args ...string) error {
	if config.Previous == nil {
		fmt.Println("You are on the first page")
		return nil
	}

	areas, err := config.Client.GetLocationAreas(config.Previous)
	if err != nil {
		return err
	}

	config.Next = areas.Next
	config.Previous = areas.Previous

	for _, area := range areas.Results {
		println(area.Name)
	}

	return nil
}

func commandExplore(config *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("Missing area argument")
	}

	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)

	area, err := config.Client.GetLocationDetails(areaName)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("Missing pokemon name")
	}

	pokemonName := strings.ToLower(args[0])
	pokemon, ok := config.CaughtPokemon[pokemonName]
	if ok {
		fmt.Printf("%s has already been caught\n", pokemon.Name)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	pokemon, err := config.Client.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	isCaught := rand.Intn(pokemon.BaseExperience) <= 60
	if !isCaught {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemon.Name)
	config.CaughtPokemon[pokemonName] = pokemon

	return nil
}

func commandInspect(config *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("Missing pokemon name")
	}

	pokemonName := strings.ToLower(args[0])
	pokemon, ok := config.CaughtPokemon[pokemonName]
	if !ok {
		fmt.Println("You have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, ptype := range pokemon.Types {
		fmt.Printf("  - %s\n", ptype.Type.Name)
	}

	return nil
}

func commandPokedex(config *Config, args ...string) error {
	if len(config.CaughtPokemon) == 0 {
		fmt.Println("Your Pokedex is empty!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.CaughtPokemon {
		fmt.Printf("  - %s\n", pokemon.Name)
	}

	return nil
}
