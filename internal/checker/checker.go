package checker

import (
	"fmt"
	"net/http"
	"time"
)

type CheckResult struct {
	URL          string
	StatusCode   int
	ResponseTime time.Duration
	CheckedAt    time.Time
}

func CheckWebsite(client *http.Client, url string) (CheckResult, error) {

	start := time.Now()
	response, err := client.Get(url)
	duration := time.Since(start)

	result := CheckResult{
		URL:          url,
		ResponseTime: duration,
		CheckedAt:    start.UTC(),
	}

	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	result.StatusCode = response.StatusCode

	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf(
			"unexpected HTTP status: %d",
			response.StatusCode,
		)
	}

	return result, nil

}
