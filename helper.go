package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"

	pokecache "github.com/OliverMackie/pokedex/internal"
)

type commands struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(c *config, param ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, param ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range getCommands() {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(c *config, param ...string) error {
	var locations LocationAreaList
	url := c.next
	if c.cache != nil {
		if data, exists := c.cache.Get(url); exists {
			if err := json.Unmarshal(data, &locations); err != nil {
				fmt.Printf("Error unmarshalling JSON from cache: %v\n", err)
				return err
			}
			c.prev = locations.Previous
			c.next = locations.Next
			for _, location := range locations.Results {
				fmt.Printf("%s\n", location.Name)
			}
			return nil
		}
	}
	req, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data from the API: %v\n", err)
		return err
	}

	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return err
	}

	if err := json.Unmarshal(data, &locations); err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return err
	}

	c.prev = locations.Previous
	c.next = locations.Next
	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}
	c.cache.Add(url, data)
	return nil
}

func commandMapb(c *config, param ...string) error {
	var locations LocationAreaList
	url := c.prev
	if url == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	if c.cache != nil {
		if data, exists := c.cache.Get(*url); exists {
			if err := json.Unmarshal(data, &locations); err != nil {
				fmt.Printf("Error unmarshalling JSON from cache: %v\n", err)
				return err
			}
			c.prev = locations.Previous
			c.next = locations.Next
			for _, location := range locations.Results {
				fmt.Printf("%s\n", location.Name)
			}
			return nil
		}
	}
	req, err := http.Get(*url)
	if err != nil {
		fmt.Printf("Error fetching data from the API: %v\n", err)
		return err
	}

	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return err
	}

	if err := json.Unmarshal(data, &locations); err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return err
	}

	c.prev = locations.Previous
	c.next = locations.Next
	for _, location := range locations.Results {
		fmt.Printf("%s\n", location.Name)
	}
	c.cache.Add(*url, data)
	return nil
}

func commandExplore(c *config, area ...string) error {
	var locationDetails LocationAreaDetails
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", area[0])
	fmt.Printf("Exploring %s...\n", area[0])
	fmt.Printf("Found Pokemon:\n")
	if c.cache != nil {
		if data, exists := c.cache.Get(url); exists {
			if err := json.Unmarshal(data, &locationDetails); err != nil {
				fmt.Printf("Error unmarshalling JSON from cache: %v\n", err)
				return err
			}
			fmt.Printf("Cached:\n")
			for _, encounter := range locationDetails.PokemonEncounters {
				fmt.Printf(" - %s\n", encounter.Pokemon.Name)
			}
			return nil
		}
	}
	req, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data from the API: %v\n", err)
		return err
	}

	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return err
	}

	if err := json.Unmarshal(data, &locationDetails); err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return err
	}
	for _, encounter := range locationDetails.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	c.cache.Add(url, data)
	return nil
}

func commandCatch(c *config, pokemon ...string) error {
	pokedex := *c.pokedex
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon[0])
	roll := rand.Intn(1000)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemon[0])
	var targetPokemon Pokemon

	req, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data from the API: %v\n", err)
		return err
	}

	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return err
	}

	if err := json.Unmarshal(data, &targetPokemon); err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return err
	}
	if roll < targetPokemon.BaseExperience {
		fmt.Printf("%s escaped!\n", pokemon[0])
		return nil
	} else {
		pokedex[pokemon[0]] = targetPokemon
		fmt.Printf("%s was caught!\n", pokemon[0])
		return nil
	}
}

func commandInspect(c *config, pokemon ...string) error {
	pokedex := *c.pokedex
	if p, exists := pokedex[pokemon[0]]; !exists {
		fmt.Printf("you have not caught that pokemon \n")
		return nil
	} else {
		fmt.Printf("Name: %s \n", p.Name)
		fmt.Printf("Height: %d \n", p.Height)
		fmt.Printf("Weight: %d \n", p.Weight)
		fmt.Printf("Stats: \n")
		for _, stat := range p.Stats {
			fmt.Printf("-%s: %d \n", stat.Stat.Name, (stat.BaseStat + stat.Effort))
		}
		fmt.Printf("Types: \n")
		for _, typeInfo := range p.Types {
			fmt.Printf("- %s \n", typeInfo.Type.Name)
		}
		return nil
	}
}

func commandPokedex(c *config, param ...string) error {
	pokedex := *c.pokedex
	if len(pokedex) == 0 {
		fmt.Printf("Your Pokedex is empty.\n")
		return nil
	}
	fmt.Printf("Your Pokedex:\n")
	for _, p := range pokedex {
		fmt.Printf(" - %s\n", p.Name)
	}
	return nil
}

func getCommands() map[string]commands {
	var CommandMap = map[string]commands{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},

		"help": {
			name:        "help",
			description: "displays a help message",
			callback:    commandHelp,
		},

		"map": {
			name:        "map",
			description: "lists next 20 locations in the Pokemon world",
			callback:    commandMap,
		},

		"mapb": {
			name:        "mapb",
			description: "lists previous 20 locations in the Pokemon world",
			callback:    commandMapb,
		},

		"explore": {
			name:        "explore",
			description: "lists the pokemon that spawn in a given location",
			callback:    commandExplore,
		},

		"catch": {
			name:        "catch",
			description: "attempts to catch the chosen pokemon",
			callback:    commandCatch,
		},

		"inspect": {
			name:        "inspect",
			description: "attempts to inspect the chosen pokemon",
			callback:    commandInspect,
		},

		"pokedex": {
			name:        "pokedex",
			description: "lists all the pokemon you have caught",
			callback:    commandPokedex,
		},
	}
	return CommandMap
}

type config struct {
	registry map[string]commands
	next     string
	prev     *string
	cache    *pokecache.Cache
	pokedex  *map[string]Pokemon
}
