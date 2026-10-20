package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	checker "github.com/vitamaxoduol/Website-Uptime-Monitor/internal/checker"
	"github.com/vitamaxoduol/Website-Uptime-Monitor/internal/storage"
)

func main() {

	url := flag.String(
		"url",
		"https://go.dev",
		"https://www.example.com",
	)

	timeout := flag.Duration(
		"timeout",
		10*time.Second,
		"Maximum time to wait for a website reponse",
	)

	history := flag.Bool(
		"history",
		false,
		"Show the last 10 saved checks",
	)

	flag.Parse()

	if *timeout <= 0 {
		log.Print("Timeout must be greater than zero")
		return
	}

	db, err := storage.Open("uptime.db")
	if err != nil {
		log.Printf("Database setup failed: %v", err)
		return
	}
	defer db.Close()

	if *history {
		checks, err := storage.ListChecks(db, 10)
		if err != nil {
			log.Printf("No checks saved yet")
			return
		}

		if len(checks) == 0 {
			fmt.Println("No checks saved yet")
			return
		}

		for _, check := range checks {
			fmt.Printf(
				"URL: %s | Status: %d | Response time: %s | Checked at: %s\n",
				check.URL,
				check.StatusCode,
				time.Duration(check.ResponseTimeNS),
				check.CheckedAt,
			)

			if check.ErrorMessage != "" {
				fmt.Printf(" Error: %s\n", check.ErrorMessage)
			}
		}

		return
	}

	client := &http.Client{
		Timeout: *timeout,
	}

	result, err := checker.CheckWebsite(client, *url)

	saveErr := storage.SaveCheck(
		db,
		result.URL,
		result.StatusCode,
		result.ResponseTime,
		result.CheckedAt,
		err,
	)

	if saveErr != nil {
		log.Printf("could not save: %v", saveErr)
	}

	if err != nil {
		log.Printf("Check failed for %s: %v", *url, err)
		return
	}

	fmt.Printf(
		"URL: %s | Status: %d | Response time: %s | Checked at: %s\n",
		result.URL,
		result.StatusCode,
		result.ResponseTime,
		result.CheckedAt.Format(time.RFC3339),
	)
}
