package main

import (
	"time"

	"github.com/Quak1/pokedex-go/internal/pokeapi"
)

func main() {
	config := Config{
		Client:        pokeapi.NewClient(5*time.Second, 10*time.Second),
		CaughtPokemon: map[string]pokeapi.Pokemon{},
	}

	repl(&config)
}
