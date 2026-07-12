package main

import (
	"testing"
	"encoding/json"
	"os"
	"fmt"
)

func TestPlayerStruct(t *testing.T) {
    file, _ := os.ReadFile("test_bnl.json")

    var player Player
    err := json.Unmarshal(file, &player)
    if err != nil {
        t.Fatalf("Unmarshal Error: %v", err)
    }

	fmt.Println("Available activity keys:")
	for k := range player.LatestSnapshot.Data.Activities {
	    fmt.Printf(" - %s\n", k)
	}

    if player.Username != "bnl" {
        t.Errorf("Username mismatch: want %s, got %s", "bnl", player.Username)
    }

    if player.LatestSnapshot.Data.Skills["attack"].Level != 99 {
        t.Errorf("Attack Level mismatch: want 99, got %d", player.LatestSnapshot.Data.Skills["attack"].Level)
    }

    if player.LatestSnapshot.Data.Skills["crafting"].Rank != 75273 {
        t.Errorf("Crafting Rank mismatch: want 75273, got %d", player.LatestSnapshot.Data.Skills["crafting"].Rank)
    }

    // NOTE: For floats like EHB/EHP, use a small delta or direct compare
    if player.LatestSnapshot.Data.Bosses["zulrah"].Ehb != 25.52174 {
        t.Errorf("Zulrah EHB mismatch: want 25.52174, got %f", player.LatestSnapshot.Data.Bosses["zulrah"].Ehb)
    }

    if player.LatestSnapshot.Data.Activities["clue_scrolls_master"].Score != 29 {
        t.Errorf("Master Score mismatch: want 29, got %d", player.LatestSnapshot.Data.Activities["master_clue_scroll"].Score)
    }

    if player.LatestSnapshot.Data.Skills["overall"].Ehp != 782.1558 {
        t.Errorf("EHP mismatch: want 782.1558, got %f", player.LatestSnapshot.Data.Skills["overall"].Ehp)
    }
}
