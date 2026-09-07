package main

import (
	"bufio"
	"fmt"
	"os"

	pokecache "github.com/OliverMackie/pokedex/internal"
)

func main() {
	pokedex := make(map[string]Pokemon)
	pokecache := pokecache.NewCache(5 * 60)
	url := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	config := &config{getCommands(), url, nil, pokecache, &pokedex}
	CommandMap := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			fmt.Printf("Error reading input. %v", scanner.Err())
			break
		}
		input := scanner.Text()
		CleanedInput := CleanInput(input)
		if len(CleanedInput) == 0 {
			CommandMap["exit"].callback(config)
			break
		}
		switch CleanedInput[0] {
		case "exit":
			CommandMap["exit"].callback(config)
		case "help":
			CommandMap["help"].callback(config)
		case "map":
			CommandMap["map"].callback(config)
		case "mapb":
			CommandMap["mapb"].callback(config)
		case "explore":
			CommandMap["explore"].callback(config, CleanedInput[1])
		case "catch":
			CommandMap["catch"].callback(config, CleanedInput[1])
		case "inspect":
			CommandMap["inspect"].callback(config, CleanedInput[1])
		case "pokedex":
			CommandMap["pokedex"].callback(config)
		default:
			fmt.Printf("Unknown command\n")
		}
	}
}
