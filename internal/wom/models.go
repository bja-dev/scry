package wom

import (
	"time"
	"fmt"
	"github.com/dustin/go-humanize"
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
// Rank 		int	`json:"rank"`			// Removed Rank because always updates
Level		int	`json:"level,omitempty"`
Ehp		float64	`json:"ehp,omitempty"`
Kills		int	`json:"kills,omitempty"`
Ehb		float64	`json:"ehb,omitempty"`
Score		int	`json:"score,omitempty"`
}

func (md MetricData) Print() {
	// NOTE: this method currently infers that difference is always positive. I'm unsure if EHP/EHB methods updating will lower those numbers, which could mean that those numbers can decrease. But, the checks should mean that won't print if its the case.
	if md.Metric == "" {
		return
	}
	fmt.Println("Stats for", md.Metric)
	if md.Experience > 0 {
		fmt.Printf("-> Experience: +%s\n",humanize.Comma(int64(md.Experience)))
	}

	if md.Level > 0 {
		fmt.Printf("-> Level: +%s\n",humanize.Comma(int64(md.Level)))
	}

	if md.Ehp > 0.0 {
		fmt.Printf("-> Ehp: +%f\n", md.Ehp)
	}

	if md.Kills > 0 {
		fmt.Printf("-> Kills: +%s\n", humanize.Comma(int64(md.Kills)))
	}

	if md.Ehb > 0.0 {
		fmt.Printf("-> Ehb: +%f\n", md.Ehb)
	}

	if md.Score > 0 {
		fmt.Printf("-> Score: +%s\n", humanize.Comma(int64(md.Score)))
	}
	fmt.Println("")
}

func (md MetricData) Equal(md2 MetricData) bool{
	return  md.Experience == md2.Experience &&
		md.Level == md2.Level &&
		md.Kills == md2.Kills &&
		md.Ehb == md2.Ehb &&
		md.Ehp == md2.Ehp &&
		md.Score == md2.Score
}

type SnapshotDiff struct {
	Exp		int
	Ehp		float64
	Ehb		float64
	Skills		map[string]MetricData `json:"skills"`
	Bosses		map[string]MetricData `json:"bosses"`
	Activities	map[string]MetricData `json:"activities"`
}

func (sd SnapshotDiff) Print() {
	fmt.Println("Snapshot Diff:")
	fmt.Printf("Exp: %d    -    Ehp: %f    -    Ehb: %f\n", sd.Exp, sd.Ehp, sd.Ehb)
	if (len(sd.Skills) >0) {
		fmt.Println("\nSkills:")
		for _, v := range sd.Skills { // NOTE: any way to ensure overall is printed first?
			v.Print()
		}
	}

	if (len(sd.Bosses) >0) {
		fmt.Println("\nBosses:")
		for _, v := range sd.Bosses {
			v.Print()
		}
	}

	if (len(sd.Activities) >0) {
		fmt.Println("\nActivities:")
		for _, v := range sd.Activities{
			v.Print()
		}
	}
}

func (p Player) Equal(p2 Player) bool {
	// Manual Shallow comparison using ID, Username, EXP, EHP and EHB. NOTE: this could bug out if u increase metrics without gaining exp. Clue Caskets, etc.
	return p.ID == p2.ID && p.Username == p2.Username && p.Exp == p2.Exp && p.Ehp == p2.Ehp && p.Ehb == p2.Ehb 
}

func diffMap(oldMap, newMap map[string]MetricData) map[string]MetricData {
	//NOTE: something is wrong here regarding EHP/EHB.
    diffs := make(map[string]MetricData)
    for name, newM := range newMap {
        oldM, exists := oldMap[name]
        if !exists { continue }

        // Create the delta struct
        delta := MetricData{
		Metric: name,
		Experience: newM.Experience - oldM.Experience,
		Level:      newM.Level - oldM.Level,
		Kills:      newM.Kills - oldM.Kills,
		Score:      newM.Score - oldM.Score,
        }

        // Only keep if something changed
	if delta.Experience != 0 || delta.Level != 0 || delta.Kills != 0 || delta.Score != 0 {
		diffs[name] = delta
		//delta.Print()
        }
    }
    return diffs
}

func (p1 Player) GetDiff(p2 Player) SnapshotDiff {
	return SnapshotDiff{
		Exp:		p2.Exp - p1.Exp,
		Ehp:		p2.Ehp - p1.Ehp,
		Ehb:		p2.Ehb - p1.Ehb,
        	Skills:		diffMap(p1.LatestSnapshot.Data.Skills, p2.LatestSnapshot.Data.Skills),
        	Bosses:		diffMap(p1.LatestSnapshot.Data.Bosses, p2.LatestSnapshot.Data.Bosses),
        	Activities: 	diffMap(p1.LatestSnapshot.Data.Activities, p2.LatestSnapshot.Data.Activities),
    }
}

//TODO: setup discord webhook
