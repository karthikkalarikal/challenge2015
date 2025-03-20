package util

import (
	"context"
	"fmt"
	"net/http"
)

func Exists(url string) error {

	err := limiter.Wait(context.Background())
	if err != nil {
		return err
	}
	resp, err := http.Get("http://data.moviebuff.com/" + url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}
