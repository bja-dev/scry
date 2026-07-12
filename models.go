package main

import (
	"time"
)
type Player struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	DisplayName    string    `json:"displayName"`
	Type           string    `json:"type"`
	Build          string    `json:"build"`
	Status         string    `json:"status"`
	Patron         bool      `json:"patron"`
	Exp            int       `json:"exp"`
	Ehp            float64   `json:"ehp"`
	Ehb            float64   `json:"ehb"`
	UpdatedAt      time.Time `json:"updatedAt"`
	LastChangedAt  time.Time `json:"lastChangedAt"`
	CombatLevel    int       `json:"combatLevel"`
	LatestSnapshot struct {
		PlayerID   int       `json:"playerId"`
		CreatedAt  time.Time `json:"createdAt"`
		ImportedAt any       `json:"importedAt"`
		Data       struct {
			Skills map[string]MetricData `json:"skills"`
			Bosses map[string]MetricData `json:"bosses"`
			Activities map[string]MetricData `json:"activities"`
		} `json:"data"`
	} `json:"latestSnapshot"`
}

type MetricData struct {
Metric		string	`json:"metric"`
Experience	int	`json:"experience,omitempty"`
Rank 		int	`json:"rank"`
Level		int	`json:"level,omitempty"`
Ehp		float64	`json:"ehp,omitempty"`
Kills		int	`json:"kills,omitempty"`
Ehb		float64	`json:"ehb,omitempty"`
Score		int	`json:"score,omitempty"`
}
