package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ukharwa/pokedex_api/internal"
)

var commands map[string]cliCommand
var cache *internal.Cache
var api_path string = "https://pokeapi.co/api/v2/"
var pokedex = make(map[string]pokemon)

func init() {
	commands = map[string]cliCommand{
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
			description: "Displays next 20 location-areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 location-areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "get a list of pokemon in the area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon you caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all the pokemon you have caught",
			callback:    commandPokedex,
		},
	}

	var err error
	cache, err = internal.NewCache(1 * time.Minute)
	if err != nil {
		fmt.Println("Error creating cache")
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := config{next: api_path + "location-area/",
		prev: "",
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) > 0 {
			if command, exists := commands[input[0]]; exists {
				if len(input) > 1 {
					err := command.callback(&config, input[1])
					if err != nil {
						fmt.Println(err)
					}
				} else {
					err := command.callback(&config, "")
					if err != nil {
						fmt.Println(err)
					}
				}
			} else {
				fmt.Println("Unknown Command")
			}
		} else {
			fmt.Println("Please enter a command.")
		}
	}
}

func cleanInput(text string) []string {
	output := strings.Fields(strings.ToLower(text))
	return output
}

func commandExit(url *config, input string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(url *config, input string) error {
	fmt.Println("Welcome to the Pokedex\nUsage:")
	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}

func commandMap(url *config, input string) error {
	val, exists, err := cache.Get(url.next)
	if err != nil {
		return err
	}

	var data locationAreaResponse

	if !exists {
		res, err := http.Get(url.next)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d", res.StatusCode)
		}

		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return err
		}

		bodyBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		cache.Add(url.next, bodyBytes)

	} else {
		err = json.Unmarshal(val, &data)
	}
	url.next = data.Next
	url.prev = data.Previous
	for _, location := range data.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(url *config, input string) error {
	val, exists, err := cache.Get(url.prev)
	if err != nil {
		return err
	}

	var data locationAreaResponse

	if !exists {
		res, err := http.Get(url.prev)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d", res.StatusCode)
		}

		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return err
		}

		bodyBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		cache.Add(url.prev, bodyBytes)

	} else {
		err = json.Unmarshal(val, &data)
	}
	url.next = data.Next
	url.prev = data.Previous
	for _, location := range data.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(c *config, input string) error {
	url := api_path + "location-area/" + input + "/"
	val, exists, err := cache.Get(url)
	if err != nil {
		return err
	}

	var data locationEncounters

	if !exists {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d", res.StatusCode)
		}

		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return err
		}

		bodyBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		cache.Add(url, bodyBytes)

	} else {
		err = json.Unmarshal(val, &data)
		if err != nil {
			return err
		}
	}

	for _, encounter := range data.Encounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(c *config, input string) error {
	url := api_path + "pokemon/" + input
	fmt.Printf("Throwing a pokeball at %s\n", input)

	val, exists, err := cache.Get(url)
	if err != nil {
		return err
	}

	var data pokemon

	if !exists {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d", res.StatusCode)
		}

		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return err
		}

		bodyBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cache.Add(url, bodyBytes)
	} else {
		err = json.Unmarshal(val, &data)
		if err != nil {
			return err
		}
	}

	base := data.BaseXp
	ratio := float64(base) / float64(650)
	probability := 1.0 - (ratio * 0.8)
	caught := rand.Float64() < probability
	if caught {
		fmt.Printf("%s was caught!\n", data.Name)
		pokedex[data.Name] = data
	} else {
		fmt.Printf("%s escaped!\n", data.Name)
	}

	return nil
}

func commandInspect(c *config, input string) error {
	poke, exists := pokedex[input]

	if exists {
		fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", poke.Name, poke.Height, poke.Weight)
		fmt.Print("Stats:\n")
		for _, stat := range poke.Stats {
			fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Print("Types:\n")
		for _, t := range poke.Types {
			fmt.Printf(" -%s\n", t.Type.Name)
		}
	} else {
		return fmt.Errorf("You have not caught that pokemon")
	}
	return nil
}

func commandPokedex(c *config, input string) error {
	fmt.Println("Your Pokedex: ")
	for key := range pokedex {
		fmt.Printf(" - %s\n", key)
	}
	return nil
}
