package pokeapi

import (
	"encoding/json"
	"net/http"
)

const (
	baseUrl = "https://pokeapi.co/api/v2/"
)

func GetLocationAreas(page *string) (LocationAreaList, error) {
	url := baseUrl + "location-area/"
	if page != nil {
		url = *page
	}

	res, err := http.Get(url)
	if err != nil {
		return LocationAreaList{}, err
	}

	var areas LocationAreaList
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&areas); err != nil {
		return LocationAreaList{}, err
	}

	return areas, nil
}
