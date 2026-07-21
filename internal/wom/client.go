package wom

import(
	"fmt"
	"os"
	"bufio"
	"time"
	"encoding/json"
	"net/http"
)

var (
	WOM_API_KEY	string
	WOM_USER_AGENT	string
)

func init() {
	WOM_API_KEY = os.Getenv("WOM_API_KEY")
	WOM_USER_AGENT = os.Getenv("WOM_USER_AGENT")
}

func scanFile(path string) ([]string, error) {
        // Takes a filepath and scans its lines into a string array
        var lines []string

        file, err := os.Open(path)
        if err != nil {
                return nil, fmt.Errorf("Error occurred while trying to open file: %v", err)
        }
        defer file.Close()
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                if scanner.Text() != "" {
                        lines = append(lines, scanner.Text())
                }
        }
        if err := scanner.Err(); err != nil {
                return nil, fmt.Errorf("Error occurred when scanning the file: %v", err)
        }
        return lines, nil
}

func GetPlayerData(username string) (Player, error) {
    filename := fmt.Sprintf("../../data/%s/current.json", username)

    data, err := os.ReadFile(filename)
    if err != nil {
        return Player{}, err
    }

    var p Player
    if err := json.Unmarshal(data, &p); err != nil {
        return Player{}, err // Don't archive garbage, return error
    }

    if time.Since(p.UpdatedAt) > 2*time.Hour {
        archiveCurrentFile(username)
        return p, fmt.Errorf("stale")
    }

    return p, nil
}

func SaveData(username string, p Player) error {
	// 1. Ensure the directory exists
	dir := fmt.Sprintf("../../data/%s", username)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 2. Marshal with indent so it's readable if you ever open it
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	// 3. Write it out
	path := fmt.Sprintf("%s/current.json", dir)
	return os.WriteFile(path, data, 0644)
}

func archiveCurrentFile(username string) error {
    sourcePath := fmt.Sprintf("../../data/%s/current.json", username)

    data, err := os.ReadFile(sourcePath)
    if err != nil {
        return err // nothing to archive
    }

    var p Player
    if err := json.Unmarshal(data, &p); err != nil {
        return err // Corrupted, don't archive
    }

    destDir := fmt.Sprintf("../../data/%s/archive", username)
    os.MkdirAll(destDir, 0755)

    destPath := fmt.Sprintf("%s/%s.json", destDir, p.UpdatedAt.Format("20060102_150405"))
    return os.Rename(sourcePath, destPath)
}


func getPlayerFromAPI(username string) (Player, error) {
	var p Player
	var err error
	p, err = GetPlayerData(username)
	if err == nil {
		return p, nil
	}

	var req *http.Request
	var resp *http.Response
	url := fmt.Sprintf("https://api.wiseoldman.net/v2/players/%s", username)
	req, err = http.NewRequest("POST", url, nil)
	if err != nil {
		return Player{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", WOM_API_KEY)
	req.Header.Set("User-Agent", WOM_USER_AGENT)
	
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return Player{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return Player{}, fmt.Errorf("update failed with status: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return Player{}, err
	}

	// Save data to avoid spamming api
	SaveData(username, p)
	return p, nil

}
