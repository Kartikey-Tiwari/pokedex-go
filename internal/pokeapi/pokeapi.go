package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/kartikey-tiwari/pokedex-go/internal/pokecache"
)

var cache *pokecache.Cache = pokecache.NewCache(5 * time.Second)

type LocationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type AreaResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Config struct {
	Next     string
	Previous string
}

func decodeJson[T any](data []byte) (T, error) {
	var res T

	if err := json.Unmarshal(data, &res); err != nil {
		return res, err
	}
	return res, nil
}

func fetchLocationData[T any](url string) (T, error) {
	var fetchedData T
	res, err := http.Get(url)
	if err != nil {
		return fetchedData, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fetchedData, err
	}

	return decodeJson[T](data)
}

func getLocationsFromResponse(locationRes LocationResponse) []string {
	locations := []string{}
	if locationRes.Results == nil {
		return locations
	}
	for _, v := range locationRes.Results {
		locations = append(locations, v.Name)
	}

	return locations
}

func setConfig(c *Config, locationRes LocationResponse) {
	c.Next = locationRes.Next
	if locationRes.Previous != nil {
		prev, ok := locationRes.Previous.(string)
		if ok {
			c.Previous = prev
		}
	}
}

func getData[T any](url string) (T, error) {
	var res T
	var err error
	data, ok := cache.Get(url)
	if ok {
		res, err = decodeJson[T](data)
		if err != nil {
			return res, err
		}
	} else {

		res, err = fetchLocationData[T](url)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

func GetLocationAreaNames(c *Config, next bool) ([]string, error) {
	var url string
	if next {
		url = c.Next
	} else {
		url = c.Previous
	}

	var locationRes LocationResponse
	locationRes, err := getData[LocationResponse](url)
	if err != nil {
		return nil, err
	}
	setConfig(c, locationRes)
	return getLocationsFromResponse(locationRes), nil
}

func GetPokemonsInArea(area string) ([]string, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + area

	var areaRes AreaResponse
	areaRes, err := getData[AreaResponse](url)
	if err != nil {
		return nil, err
	}

	pokemons := []string{}
	for _, pokemon := range areaRes.PokemonEncounters {
		pokemons = append(pokemons, pokemon.Pokemon.Name)
	}

	return pokemons, nil
}
