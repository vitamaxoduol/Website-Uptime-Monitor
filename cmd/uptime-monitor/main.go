package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	checker "github.com/vitamaxoduol/Website-Uptime-Monitor/internal/checker"
	storage "github.com/vitamaxoduol/Website-Uptime-Monitor/internal/storage"
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

	addMonitor := flag.String(
		"add-monitor",
		"",
		"Website URL to save as a monitor",
	)
	interval := flag.Int(
		"interval",
		60,
		"Monitor check interval in seconds",
	)

	listMonitors := flag.Bool(
		"list-monitors",
		false,
		"List saved monitors",
	)

	deleteMonitor := flag.Int64(
		"delete-monitor",
		0,
		"ID of the monitor to delete",
	)

	flag.Parse()

	// this block rejects conflicting flags
	actionCount := 0

	flag.Visit(func(currentFlag *flag.Flag) {
		switch currentFlag.Name {
		case "add-monitor", "list-monitors", "delete-monitor", "history":
			actionCount++
		}
	})

	if actionCount > 1 {
		log.Print("Choose only action: -add-monitor, -list-monitors, -delete-monitor, or -history")
		return
	}
	if *timeout <= 0 {
		log.Print("Timeout must be greater than zero")
		return
	}

	invalidAction := false

	flag.Visit(func(currentFlag *flag.Flag) {
		switch currentFlag.Name {
		case "add-monitor":
			if *addMonitor == "" {
				log.Print("-add-monitor requires a website URL")
				invalidAction = true
			}
		case "delete-monitor":
			if *deleteMonitor <= 0 {
				log.Print("-delete-monitor requires a positive monitor ID")
				invalidAction = true
			}
		}
	})

	if invalidAction {
		return
	}

	db, err := storage.Open("uptime.db")
	if err != nil {
		log.Printf("Database setup failed: %v", err)
		return
	}
	defer db.Close()

	if *addMonitor != "" {
		id, err := storage.CreateMonitor(db, *addMonitor, *interval)
		if err != nil {
			log.Printf("Could not create monitor: %v", err)
			return
		}

		fmt.Printf(
			"Monitor created  | ID: %d | URL: %s | Interval: %ds\n",
			id,
			*addMonitor,
			*interval,
		)
		return
	}

	if *listMonitors {
		monitors, err := storage.ListMonitors(db)
		if err != nil {
			log.Fatalf("Could not list monitors: %v", err)
			return
		}

		if len(monitors) == 0 {
			fmt.Println("No monitors found yet")
		}

		for _, monitor := range monitors {
			fmt.Printf(
				"ID: %d | URL: %s | Interval: %ds | Created at: %s\n",
				monitor.ID,
				monitor.URL,
				monitor.IntervalSeconds,
				monitor.CreatedAt,
			)

		}

		return
	}

	if *deleteMonitor != 0 {
		deleted, err := storage.DeleteMonitor(db, *deleteMonitor)
		if err != nil {
			log.Printf("Could not delete monitor: %v", err)
			return
		}

		if !deleted {
			fmt.Printf("No monitor found with ID %d\n", *deleteMonitor)
			return
		}

		fmt.Printf("Monitor deleted | ID: %d\n", *deleteMonitor)
		return
	}

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
