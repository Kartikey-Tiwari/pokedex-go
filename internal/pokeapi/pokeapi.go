package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next     string
	Previous string
}

func fetchLocationData(url string) (LocationResponse, error) {
	var locationRes LocationResponse

	res, err := http.Get(url)
	if err != nil {
		return LocationResponse{}, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&locationRes); err != nil {
		return LocationResponse{}, err
	}

	return locationRes, nil
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
	fmt.Println("next:", c.Next)
	fmt.Println("previous:", c.Previous)
}

func GetLocationAreaNames(c *Config, next bool) ([]string, error) {
	var url string
	if next {
		url = c.Next
	} else {
		url = c.Previous
	}
	locationRes, err := fetchLocationData(url)
	if err != nil {
		return nil, err
	}

	setConfig(c, locationRes)
	return getLocationsFromResponse(locationRes), nil
}
