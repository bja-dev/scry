package wom

import (
	"testing"
	"encoding/json"
	"os"
//	"fmt"
)

func TestPlayerStruct(t *testing.T) {
    file, _ := os.ReadFile("../../data/test/test_bnl.json")

    var player Player
    err := json.Unmarshal(file, &player)
    if err != nil {
        t.Fatalf("Unmarshal Error: %v", err)
    }

	/*fmt.Println("Available activity keys:")
	for k := range player.LatestSnapshot.Data.Activities {
	    fmt.Printf(" - %s\n", k)
	}
	fmt.Println("Available skill keys:")
	for k := range player.LatestSnapshot.Data.Skills{
	    fmt.Printf(" - %s\n", k)
	}

	fmt.Println("Available boss keys:")
	for k := range player.LatestSnapshot.Data.Bosses{
	    fmt.Printf(" - %s\n", k)
	}*/


    if player.Username != "bnl" {
        t.Errorf("Username mismatch: want %s, got %s", "bnl", player.Username)
    }

    if player.LatestSnapshot.Data.Skills["attack"].Level != 99 {
        t.Errorf("Attack Level mismatch: want 99, got %d", player.LatestSnapshot.Data.Skills["attack"].Level)
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
	//fmt.Printf("DEBUG: Skills count: %d, Bosses count: %d, Activities count: %d\n", len(got.Skills), len(got.Bosses), len(got.Activities))
	/*
	for _, v := range got.Skills {
		v.Print()
	}
	for _, v := range got.Bosses {
		v.Print()
	}
	for _, v := range got.Activities {
		v.Print()
	}
	*/
	want := SnapshotDiff{
		Exp:	20000000,
		Skills:		make(map[string]MetricData),
		Bosses:		make(map[string]MetricData),
		Activities:	make(map[string]MetricData),
	}
	want.Skills["magic"] = MetricData{Experience: 1000000}
	want.Skills["fletching"] = MetricData{Experience: 2000000}
	want.Skills["crafting"] = MetricData{Experience: 20000}
	want.Skills["herblore"] = MetricData{Experience: 3}
	want.Skills["slayer"] = MetricData{Level: 1, Experience: 1000000}
	want.Bosses["abyssal_sire"] = MetricData{Kills: 60}
	want.Activities["clue_scrolls_all"] = MetricData{Score: 1}
	want.Activities["clue_scrolls_master"] = MetricData{Score: 1}
	if got.Exp != want.Exp {
		t.Errorf("Error: experience was not accurate: %d != %d\n", got.Exp, want.Exp)
	}
	if !want.Skills["magic"].Equal(got.Skills["magic"]) {
		t.Errorf("Error: magic was not accurate: %d != %d\n", got.Skills["magic"].Experience, want.Skills["magic"].Experience)
	}
	if !want.Skills["fletching"].Equal(got.Skills["fletching"]) {
		t.Errorf("Error: fletching was not accurate: %d != %d\n", got.Skills["fletching"].Experience, want.Skills["fletching"].Experience)
	}
	if !want.Skills["crafting"].Equal(got.Skills["crafting"]) {
		t.Errorf("Error: crafting was not accurate: %d != %d\n", got.Skills["crafting"].Experience, want.Skills["crafting"].Experience)
	}
	if !want.Skills["herblore"].Equal(got.Skills["herblore"]) {
		t.Errorf("Error: herblore was not accurate: %d != %d\n", got.Skills["herblore"].Experience, want.Skills["herblore"].Experience)
	}
	if !want.Skills["slayer"].Equal(got.Skills["slayer"]) {
		t.Errorf("Error: slayer was not accurate: %d != %d\n", got.Skills["slayer"].Level, want.Skills["slayer"].Level)
	}
	if !want.Bosses["abyssal_sire"].Equal(got.Bosses["abyssal_sire"]) {
		t.Errorf("Error: abyssal_sire was not accurate: %d != %d\n", got.Bosses["abyssal_sire"].Kills, want.Bosses["abyssal_sire"].Kills)
	}
	if !want.Activities["clue_scrolls_all"].Equal(got.Activities["clue_scrolls_all"]) {
		t.Errorf("Error: clue_scrolls_all was not accurate: %d != %d\n", got.Activities["clue_scrolls_all"].Score, want.Activities["clue_scrolls_all"].Kills)
	}
	if !want.Activities["clue_scrolls_master"].Equal(got.Activities["clue_scrolls_master"]) {
		t.Errorf("Error: clue_scrolls_master was not accurate: %d != %d\n", got.Activities["clue_scrolls_master"].Score, want.Activities["clue_scrolls_master"].Score)
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
	got.Print()
}
