package loaders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 20 * time.Second}

func getJSON(url string, out interface{}) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP response %d from %s", resp.StatusCode, url)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
