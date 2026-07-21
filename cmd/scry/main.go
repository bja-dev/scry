package main

import (
	"github.com/bja-dev/scry/internal/wom"
	"fmt"
	"bytes"
	"net/http"
	"os"
	"time"
	"log"
	"github.com/joho/godotenv"
)



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var DISCORD_WEBHOOK_URL = os.Getenv("DISCORD_WEBHOOK_URL")
	if DISCORD_WEBHOOK_URL == "" {
		log.Fatalf("Error: DISCORD_WEBHOOK_URL environment variable is not set")
	}
	players, e := wom.ScanFile("users.txt")
	if e != nil {
		log.Fatalf("Error: %v\n", e)
		os.Exit(1)
	}
	for _, player := range players {
		current, er := wom.ReadLocalPlayer(player)
		if er != nil {
			log.Fatalf("Error: %v\n", er)
			os.Exit(1)
		}
		p, err := wom.GetPlayerFromAPI(player)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
			os.Exit(1)
		}

		if current.Username == "" {
			log.Printf("initialised baseline: %s",player)
			continue
		}

		if p.UpdatedAt == current.UpdatedAt {
			// curl update "no change", or no message
			continue
		}
		diff := p.GetDiff(current)
		content := fmt.Sprintf("{ \"content\": \"%s\"}", diff.Format())
		payloadBytes := []byte(content)
		resp, err := http.Post(DISCORD_WEBHOOK_URL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Fatalf("failed to send webhook request: %v", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			log.Fatalf("discord webhook returned non-success status: %d", resp.StatusCode)
			os.Exit(1)
		}
		time.Sleep(3 * time.Second) // see if this is necessary
	}
	return
}

