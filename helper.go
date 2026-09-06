package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pokedex/internal/pokecache"
	"strings"
)

type commands struct {
	name        string
	description string
	callback    func(*config) error
}

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(*config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(*config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range getCommands() {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(c *config) error {
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
				fmt.Printf("  %s\n", location.Name)
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
		fmt.Printf("  %s\n", location.Name)
	}
	c.cache.Add(url, data)
	return nil
}

func commandMapb(c *config) error {
	var locations LocationAreaList
	url := c.prev
	if c.cache != nil {
		if data, exists := c.cache.Get(url); exists {
			if err := json.Unmarshal(data, &locations); err != nil {
				fmt.Printf("Error unmarshalling JSON from cache: %v\n", err)
				return err
			}
			c.prev = locations.Previous
			c.next = locations.Next
			for _, location := range locations.Results {
				fmt.Printf("  %s\n", location.Name)
			}
			return nil
		}
	}
	if url == nil {
		fmt.Println("you're on the first page")
		return nil
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
		fmt.Printf("  %s\n", location.Name)
	}
	c.cache.Add(*url, data)
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
	}
	return CommandMap
}

type config struct {
	registry map[string]commands
	next     string
	prev     *string
	cache    *pokecache.Cache
}
