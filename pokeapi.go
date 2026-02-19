package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yelkcubwerdna/pokedex/internal/pokecache"
)

func (c *Client) pokeapi_location_area(url string) location_area {
	location := location_area{}

	// See if data is in cache
	data, ok := c.cache.Get(url)
	if ok {
		err := json.Unmarshal(data, &location)
		if err != nil {
			fmt.Printf("Error unmarshal: %v\n", err)
		}
	} else {
		res, err := c.httpClient.Get(url)
		if err != nil {
			fmt.Printf("Error with GET request: %v\n", err)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			fmt.Printf("Response failed with status code: %d\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			fmt.Printf("Error with reading body: %v\n", err)
		}

		//fmt.Printf("Fetched data from url: %s\n", body)

		err = json.Unmarshal(body, &location)
		if err != nil {
			fmt.Printf("Error unmarshal: %v\n", err)
		}
		c.cache.Add(url, body)
	}

	return location
}

func (c *Client) pokeapi_pokemon(url string) pokemon_full {
	pokemon := pokemon_full{}

	data, ok := c.cache.Get(url)

	if ok {
		err := json.Unmarshal(data, &pokemon)
		if err != nil {
			fmt.Printf("Unmarshal Error: %v", err)
		}
	} else {
		res, err := c.httpClient.Get(url)

		if err != nil {
			fmt.Printf("Error with http.Get: %v\n", err)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			fmt.Printf("Response failed with status code: %d\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			fmt.Printf("Error reading body: %v", err)
		}

		err = json.Unmarshal(body, &pokemon)
		if err != nil {
			fmt.Printf("Unmarshal Error: %v", err)
		}

		c.cache.Add(url, body)
	}

	return pokemon
}

type pokemon_full struct {
	Name     string  `json:"name"`
	Base_exp int     `json:"base_experience"`
	Height   int     `json:"height"`
	Weight   int     `json:"weight"`
	Stats    []stat  `json:"stats"`
	Types    []types `json:"types"`
}

type types struct {
	Slot int       `json:"slot"`
	Type type_info `json:"type"`
}

type type_info struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type stat struct {
	Base_stat int       `json:"base_stat"`
	Effort    int       `json:"effort"`
	Stat      stat_info `json:"stat"`
}

type stat_info struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type pokemon_encounter struct {
	Pokemon         pokemon    `json:"pokemon"`
	Version_details []struct{} `json:"version_details"`
}

type location_area struct {
	Encounter_method_rates []struct{}          `json:"encounter_method_rates"`
	Game_index             int                 `json:"game_index"`
	Id                     int                 `json:"id"`
	Location               struct{}            `json:"location"`
	Name                   string              `json:"name"`
	Names                  []struct{}          `json:"names"`
	Pokemon_encounters     []pokemon_encounter `json:"pokemon_encounters"`
}

type Client struct {
	cache      pokecache.Cache
	httpClient http.Client
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (p *pokemon_full) Inspect() {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.Base_stat)
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
}
