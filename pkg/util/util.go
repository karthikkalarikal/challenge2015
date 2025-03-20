package util

import (
	"encoding/json"
	"log"
	"net/http"
)

// fetch resources
func GetByURL(url string, target any) error {
	resp, err := http.Get("http://data.moviebuff.com/" + url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println(resp.Body)
	return json.NewDecoder(resp.Body).Decode(target)
}
