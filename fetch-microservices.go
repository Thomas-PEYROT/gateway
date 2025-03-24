package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MicroserviceInstance struct {
	Port uint32 `json:"port"`
}

func FetchMicroservices(apiURL string) (map[string]map[string]MicroserviceInstance, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]map[string]MicroserviceInstance
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
