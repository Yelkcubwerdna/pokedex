package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	client := NewClient(10*time.Second, 5*time.Minute)
	config := config{}
	config.Next = 1
	config.Previous = 1
	config.Client = client

	config.Pokedex = make(map[string]pokemon_full)

	commands := getCommands()

	scanner := bufio.NewScanner(os.Stdin)
	var input string

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input = scanner.Text()
		cInput := cleanInput(input)

		cmd, ok := commands[cInput[0]]

		if !ok {
			fmt.Println("Unknown command")
		} else {
			args := []string{}
			if len(cInput) > 1 {
				args = cInput[1:]
			}
			cmd.callback(&config, args)
		}
	}
}

type config struct {
	Next     int
	Previous int
	Client   Client
	Pokedex  map[string]pokemon_full
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *config, args []string) error
}

func getCommands() map[string]cliCommand {

	commands := map[string]cliCommand{
		"catch": {
			name:        "catch",
			description: "Attenpt to catch a pokemon",
			callback:    commandCatch,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "List all Pokemon found in a location",
			callback:    commandExplore,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"inspect": {
			name:        "inspect",
			description: "Display info on caught pokemon",
			callback:    commandInspect,
		},
		"map": {
			name:        "map",
			description: "List next 20 map locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "List previous 20 map locations",
			callback:    commandMapB,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all the pokemon registered in your pokedex",
			callback:    commandPokedex,
		},
	}

	return commands
}

func commandExit(config *config, args []string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, args []string) error {
	commands := getCommands()

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range commands {
		n := cmd.name
		d := cmd.description

		fmt.Printf("%s: %s\n", n, d)
	}

	return nil
}

func commandMap(config *config, args []string) error {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"

	start := config.Next
	end := start + 20

	config.Previous = max(start-20, 1)
	config.Next = end

	for i := start; i < end; i++ {
		fullUrl := baseUrl + fmt.Sprintf("%d", i)
		l := config.Client.pokeapi_location_area(fullUrl)
		fmt.Println(l.Name)
	}

	return nil
}

func commandMapB(config *config, args []string) error {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"

	start := config.Previous
	end := start + 20

	config.Previous = max(start-20, 1)
	config.Next = end

	for i := start; i < end; i++ {
		fullUrl := baseUrl + fmt.Sprintf("%d", i)
		l := config.Client.pokeapi_location_area(fullUrl)
		fmt.Println(l.Name)
	}

	return nil
}

func commandExplore(config *config, args []string) error {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	fullUrl := baseUrl + args[0]
	location := config.Client.pokeapi_location_area(fullUrl)

	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Pokemon found:")

	for _, e := range location.Pokemon_encounters {
		fmt.Printf(" - %s\n", e.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *config, args []string) error {
	baseUrl := "https://pokeapi.co/api/v2/pokemon/"
	fullUrl := baseUrl + args[0]

	pokemon := config.Client.pokeapi_pokemon(fullUrl)

	if pokemon.Name == "" {
		fmt.Println("Invalid pokemon name")
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	roll := rand.Intn(100)

	catchRate := max(100-((pokemon.Base_exp/100)*25), 25)

	if roll <= catchRate {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		config.Pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(config *config, args []string) error {
	pokemon, ok := config.Pokedex[args[0]]

	if !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		pokemon.Inspect()
	}

	return nil
}

func commandPokedex(config *config, args []string) error {
	fmt.Println("Your Pokedex:")
	for _, p := range config.Pokedex {
		fmt.Printf(" - %s\n", p.Name)
	}

	return nil
}
