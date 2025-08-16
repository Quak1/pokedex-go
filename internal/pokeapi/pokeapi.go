package pokeapi

import (
	"encoding/json"
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
