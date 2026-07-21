package main

import (
	"github.com/bja-dev/scry/internal/wom"
	"bytes"
	"net/http"
	"os"
	"time"
	"log"
	"github.com/joho/godotenv"
	"encoding/json"
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
			log.Printf("ERROR: reading local player %s: %v", player, er)
			continue
		}
		p, err := wom.GetPlayerFromAPI(player)
		if err != nil {
			log.Printf("ERROR: fetching API data for %s: %v", player, err)
			continue
		}

		if current.Username == "" {
			log.Printf("initialised baseline: %s",player)
			continue
		}

		log.Printf("DEBUG: Comparing player %s -> Local UpdatedAt: [%s] vs API UpdatedAt: [%s]", 
            player, current.UpdatedAt, p.UpdatedAt)
		if p.UpdatedAt == current.UpdatedAt {
			// curl update "no change", or no message
			log.Printf("DEBUG: No changes detected for %s. skipping", player)
			continue
		}
		log.Printf("DEBUG: Changes detected, sending webhook for %s", player)
		diff := p.GetDiff(current)
		payload := map[string]string{
    "content": diff.Format(),
}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("failed to marshal JSON: %v", err)
			continue
		}

		resp, err := http.Post(DISCORD_WEBHOOK_URL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("failed to send webhook request: %v", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			log.Printf("discord webhook returned non-success status: %d", resp.StatusCode)
		}
		time.Sleep(3 * time.Second) // see if this is necessary
	}
	return
}

