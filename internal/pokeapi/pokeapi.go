package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Quak1/pokedex-go/internal/pokecache"
)

const (
	baseUrl = "https://pokeapi.co/api/v2/"
)

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

func (c *Client) GetLocationAreas(page *string) (LocationAreaList, error) {
	url := baseUrl + "location-area/"
	if page != nil {
		url = *page
	}

	data, ok := c.cache.Get(url)
	if !ok {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return LocationAreaList{}, err
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return LocationAreaList{}, err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return LocationAreaList{}, err
		}

		c.cache.Add(url, data)
	}

	var areas LocationAreaList
	if err := json.Unmarshal(data, &areas); err != nil {
		return LocationAreaList{}, err
	}

	return areas, nil
}

func (c *Client) GetLocationDetails(name string) (LocationArea, error) {
	url := baseUrl + "location-area/" + name

	data, ok := c.cache.Get(url)
	if !ok {
		res, err := c.httpClient.Get(url)
		if err != nil {
			return LocationArea{}, err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			return LocationArea{}, fmt.Errorf("Area not found, check your spelling")
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return LocationArea{}, err
		}

		c.cache.Add(url, data)
	}

	var details LocationArea
	if err := json.Unmarshal(data, &details); err != nil {
		return LocationArea{}, err
	}

	return details, nil
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	url := baseUrl + "pokemon/" + name

	data, ok := c.cache.Get(url)
	if !ok {
		res, err := c.httpClient.Get(url)
		if err != nil {
			return Pokemon{}, err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			return Pokemon{}, fmt.Errorf("Pokemon: %s, not found", name)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return Pokemon{}, err
		}

		c.cache.Add(url, data)
	}

	var details Pokemon
	if err := json.Unmarshal(data, &details); err != nil {
		return Pokemon{}, err
	}

	return details, nil
}
