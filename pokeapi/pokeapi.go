package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"pokecache"
)

func GetReq[T any](url string, response T, p *PokeApi, createErrorCb func(*http.Response) error) (T, error) {
	data, ok := p.cache.Get(url)
	if ok {
		fmt.Println("Cache hit!!!")
		err := json.Unmarshal(data, &response)
		if err != nil {
			log.Fatal(err)
		}

		return response, nil
	}

	time.Sleep(2 * time.Second)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if createErrorCb != nil && createErrorCb(res) != nil {
		return response, createErrorCb(res)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	p.cache.Add(url, body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	return response, nil
}

type RegionResult struct {
	Results []Region `json:"results"`
}

type Region struct {
	Name string `json:"name"`
}

type PokeApi struct {
	cache *pokecache.Cache
}

func (p *PokeApi) GetMap(limit, offset int) RegionResult {
	regionResult := RegionResult{}
	res, err := GetReq(
		fmt.Sprintf("https://pokeapi.co/api/v2/location/?limit=%d&offset=%d", limit, offset),
		regionResult,
		p,
		nil)

	if err != nil {
		log.Fatal(err)
	}

	return res
}

type Pokemon struct {
	Name           string `json:"name"`
	Id             int    `json:"id"`
	BaseExperience int    `json:"base_experience"`
	Weight         int    `json:"weight"`
	Height         int    `json:"height"`
	Types          []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

func (p Pokemon) PrintInfo() {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Id: %d\n", p.Id)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)

	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("\t-%s\n", t.Type.Name)
	}

	fmt.Println("Stats:")
	for _, stat := range p.Stats {
		fmt.Printf("\t-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("")
}

type ExploreResult struct {
	PokemonEncounters []struct {
		Pokemon Pokemon `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (p *PokeApi) ExploreLocation(location string) (ExploreResult, error) {
	exploreResult := ExploreResult{}
	res, err := GetReq(
		fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location),
		exploreResult,
		p,
		func(res *http.Response) error {
			if res.StatusCode == 404 {
				return errors.New("location not found")
			}
			return nil
		},
	)

	if err != nil {
		return exploreResult, err
	}

	return res, nil
}

func (p *PokeApi) GetPokemon(pokemon string) (Pokemon, error) {
	pokemonResult := Pokemon{}
	res, err := GetReq(
		fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon),
		pokemonResult,
		p,
		func(res *http.Response) error {
			if res.StatusCode == 404 {
				return errors.New("pokemon not found")
			}
			return nil
		},
	)

	if err != nil {
		log.Fatal(err)
		return pokemonResult, err
	}

	return res, nil
}

func NewPokeApi() PokeApi {
	return PokeApi{cache: pokecache.NewCache(60 * time.Second)}
}
