package wom

import (
	"encoding/json"
	"os"
	"testing"
)

func TestPlayerStruct(t *testing.T) {
	file, _ := os.ReadFile("../../data/test/test_bnl.json")

	var player Player
	err := json.Unmarshal(file, &player)
	if err != nil {
		t.Fatalf("Unmarshal Error: %v", err)
	}

	if player.Username != "bnl" {
		t.Errorf("Username mismatch: want %s, got %s", "bnl", player.Username)
	}

	if player.LatestSnapshot.Data.Skills["attack"].Level != 99 {
		t.Errorf("Attack Level mismatch: want 99, got %d", player.LatestSnapshot.Data.Skills["attack"].Level)
	}

	if player.LatestSnapshot.Data.Bosses["zulrah"].Ehb != 25.52174 {
		t.Errorf("Zulrah EHB mismatch: want 25.52174, got %f", player.LatestSnapshot.Data.Bosses["zulrah"].Ehb)
	}

	// Fixed key to match standard WOM naming ("master_clue_scroll") TODO: is this wrong? wtf
	if player.LatestSnapshot.Data.Activities["clue_scrolls_master"].Score != 29 {
		t.Errorf("Master Score mismatch: want 29, got %d", player.LatestSnapshot.Data.Activities["clue_scrolls_master"].Score)
	}

	if player.LatestSnapshot.Data.Skills["overall"].Ehp != 782.1558 {
		t.Errorf("EHP mismatch: want 782.1558, got %f", player.LatestSnapshot.Data.Skills["overall"].Ehp)
	}
}

func TestEqual(t *testing.T) {
	file, _ := os.ReadFile("../../data/test/test_bnl.json")

	var p1 Player
	err := json.Unmarshal(file, &p1)
	if err != nil {
		t.Fatalf("Unmarshal Error: %v", err)
	}

	file2, _ := os.ReadFile("../../data/test/test_dfg.json")
	var p2 Player
	err = json.Unmarshal(file2, &p2)
	if err != nil {
		t.Fatalf("Unmarshal Error: %v", err)
	}
	if p1.Equal(p2) {
		t.Errorf("Difference mismatch: want %t, got %t", false, p1.Equal(p2))
	}
	if !p1.Equal(p1) {
		t.Errorf("Difference mismatch: want %t, got %t", true, p1.Equal(p1))
	}
}

func TestDiff(t *testing.T) {
	file, _ := os.ReadFile("../../data/test/test_bnl.json")
	var p1 Player
	err := json.Unmarshal(file, &p1)
	if err != nil {
		t.Fatalf("Unmarshal Error: %v", err)
	}

	file2, _ := os.ReadFile("../../data/test/test_bnl2.json")
	var p2 Player
	err = json.Unmarshal(file2, &p2)
	if err != nil {
		t.Fatalf("Unmarshal Error: %v", err)
	}

	got := p1.GetDiff(p2)

	want := SnapshotDiff{
		Exp:        20000000,
		Skills:     make(map[string]MetricData),
		Bosses:     make(map[string]MetricData),
		Activities: make(map[string]MetricData),
	}
	want.Skills["magic"] = MetricData{Metric: "magic", Experience: 1000000}
	want.Skills["fletching"] = MetricData{Metric: "fletching", Experience: 2000000}
	want.Skills["crafting"] = MetricData{Metric: "crafting", Experience: 20000}
	want.Skills["herblore"] = MetricData{Metric: "herblore", Experience: 3}
	want.Skills["slayer"] = MetricData{Metric: "slayer", Level: 1, Experience: 1000000}
	want.Bosses["abyssal_sire"] = MetricData{Metric: "abyssal_sire", Kills: 60}
	want.Activities["clue_scrolls_all"] = MetricData{Metric: "clue_scrolls_all", Score: 1}
	want.Activities["clue_scrolls_master"] = MetricData{Metric: "clue_scrolls_master", Score: 1}

	if got.Exp != want.Exp {
		t.Errorf("Error: experience was not accurate: %d != %d\n", got.Exp, want.Exp)
	}
	if !want.Skills["magic"].Equal(got.Skills["magic"]) {
		t.Errorf("Error: magic was not accurate\n")
	}
	if !want.Skills["fletching"].Equal(got.Skills["fletching"]) {
		t.Errorf("Error: fletching was not accurate\n")
	}
	if !want.Skills["crafting"].Equal(got.Skills["crafting"]) {
		t.Errorf("Error: crafting was not accurate\n")
	}
	if !want.Skills["herblore"].Equal(got.Skills["herblore"]) {
		t.Errorf("Error: herblore was not accurate\n")
	}
	if !want.Skills["slayer"].Equal(got.Skills["slayer"]) {
		t.Errorf("Error: slayer was not accurate\n")
	}
	if !want.Bosses["abyssal_sire"].Equal(got.Bosses["abyssal_sire"]) {
		t.Errorf("Error: abyssal_sire was not accurate\n")
	}
	if !want.Activities["clue_scrolls_all"].Equal(got.Activities["clue_scrolls_all"]) {
		t.Errorf("Error: clue_scrolls_all was not accurate\n")
	}
	if !want.Activities["clue_scrolls_master"].Equal(got.Activities["clue_scrolls_master"]) {
		t.Errorf("Error: clue_scrolls_master was not accurate\n")
	}
}

func TestPrintDiffNoAssertion(t *testing.T) {
	file, _ := os.ReadFile("../../data/test/test_bnl.json")
	var p3 Player
	err := json.Unmarshal(file, &p3)
	if err != nil {
		t.Fatalf("Unmarshal Error: %v", err)
	}

	file2, _ := os.ReadFile("../../data/test/test_bnl-updated.json")
	var p4 Player
	err = json.Unmarshal(file2, &p4)
	if err != nil {
		t.Fatalf("Unmarshal Error: %v", err)
	}

	got := p3.GetDiff(p4)
	print(got.Format())
}
