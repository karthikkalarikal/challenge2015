package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Every(1*time.Millisecond), 1000)

// fetch resources
func GetByURL(url string, target any) error {

	err := limiter.Wait(context.Background())
	if err != nil {
		return err
	}
	resp, err := http.Get("http://data.moviebuff.com/" + url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {

		return fmt.Errorf("access denied for URL: %s", url)

	}
	return json.NewDecoder(resp.Body).Decode(target)
}
