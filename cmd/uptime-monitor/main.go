package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vitamaxoduol/Website-Uptime-Monitor/internal"
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

	flag.Parse()

	if *timeout <= 0 {
		log.Print("Timeout must be greater than zero")
		return
	}

	client := &http.Client{
		Timeout: *timeout,
	}

	result, err := internal.CheckWebsite(client, *url)
	if err != nil {
		log.Printf("Check failed for %s: %v", *url, err)
		return
	}

	fmt.Printf(
		// "Website is up | Status: %d | Response time: %s\n",
		"URL: %s | Status: %d | Response time: %s\n",
		result.URL,
		result.StatusCode,
		result.ResponseTime,
	)
}
