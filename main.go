package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex/internal/pokecache"
)

func main() {
	pokecache := pokecache.NewCache(60 * 1000) // 5 minutes
	url := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	config := &config{getCommands(), url, nil, pokecache}
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
		default:
			fmt.Printf("Unknown command\n")
		}
	}
}
