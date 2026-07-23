package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bja-dev/scry/internal/wom"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	discordWebhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if discordWebhookURL == "" {
		log.Fatalf("Error: DISCORD_WEBHOOK_URL environment variable is not set")
	}

	players, e := wom.ScanFile("users.txt")
	if e != nil {
		log.Fatalf("Error: %v\n", e)
	}

	for _, player := range players {
		current, er := wom.ReadLocalPlayer(player)
		if er != nil {
			log.Printf("Error: reading local player %s: %v", player, er)
			continue
		}

		p, err := wom.FetchFromAPI(player)
		if err != nil {
			log.Printf("Error: fetching API data for %s: %v", player, err)
			continue
		}

		if current.Username == "" {
			log.Printf("Initialised baseline for: %s", player)
			if err := wom.SaveData(player, p); err != nil {
				log.Printf("Error: failed to save initial baseline for %s: %v", player, err)
			}
			continue
		}

		log.Printf("DEBUG: Comparing player %s -> Local UpdatedAt: [%s] vs API UpdatedAt: [%s]",
			player, current.UpdatedAt, p.UpdatedAt)

		if p.UpdatedAt == current.UpdatedAt {
			log.Printf("DEBUG: No changes detected for %s. Skipping.", player)
			continue
		}

		log.Printf("DEBUG: Changes detected, building diff and sending webhook for %s", player)

		diff := p.GetDiff(current)
		payload := map[string]string{
			"content": diff.Format(),
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("failed to marshal JSON for %s: %v", player, err)
			continue
		}

		resp, err := http.Post(discordWebhookURL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("failed to send webhook request for %s: %v", player, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			log.Printf("discord webhook returned non-success status for %s: %d", player, resp.StatusCode)
			continue
		}

		if err := wom.SaveData(player, p); err != nil {
			log.Printf("Error: failed to update local baseline for %s: %v", player, err)
		} else {
			log.Printf("Successfully updated baseline for %s", player)
		}

		time.Sleep(3 * time.Second)
	}
}
