//go:build manual

package wom

import (
	"testing"
)

// go test -v -tags manual ./...
func TestFetchPlayerFromAPI(t *testing.T) {
	username := "bnl"

	t.Logf("Attempting to fetch data for: %s", username)

	p, err := getPlayerFromAPI(username)
	if err != nil {
		t.Fatalf("Failed to fetch player: %v", err)
	}

	if p.Username == "" {
		t.Error("Expected a username, got an empty string")
	}

	t.Logf("Successfully fetched player: %+v", p)
}
