package main

import (
	"fmt"
	"github.com/Quak1/pokedex-go/internal/pokeapi"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     *string
	Previous *string
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
	}
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println()

	return nil
}

func commandMap(config *Config) error {
	areas, err := pokeapi.GetLocationAreas(config.Next)
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

func commandMapb(config *Config) error {
	if config.Previous == nil {
		fmt.Println("You are on the first page")
		return nil
	}

	areas, err := pokeapi.GetLocationAreas(config.Previous)
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
