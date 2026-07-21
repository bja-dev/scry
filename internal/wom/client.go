package wom

import(
	"fmt"
	"os"
	"bufio"
	"time"
	"encoding/json"
	"net/http"
	"github.com/joho/godotenv"
	"log"
)

var (
	WOM_API_KEY	string
	WOM_USER_AGENT	string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	WOM_API_KEY = os.Getenv("WOM_API_KEY")
	WOM_USER_AGENT = os.Getenv("WOM_USER_AGENT")
}

func ScanFile(path string) ([]string, error) {
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

func ReadLocalPlayer(username string) (Player, error) {
	filename := fmt.Sprintf("data/%s/current.json", username)

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// find most recent archive if current doesn't exist (shouldn't happen but errors may occur while executing)
			archiveDir := fmt.Sprintf("data/%s/archive", username)
			files, readErr := os.ReadDir(archiveDir)
			if readErr == nil && len(files) > 0 {
				latestArchive := fmt.Sprintf("%s/%s", archiveDir, files[len(files)-1].Name())
				archiveData, readArchiveErr := os.ReadFile(latestArchive)
				if readArchiveErr == nil {
					var p Player
					if json.Unmarshal(archiveData, &p) == nil {
						_ = os.WriteFile(filename, archiveData, 0644)
						return p, nil
					}
				}
			}
			return Player{}, nil
		}
		return Player{}, err
	}

	var p Player
	if err := json.Unmarshal(data, &p); err != nil {
		return Player{}, err
	}

	return p, nil
}
func IsDataStale(username string) (bool, error) {
    p, err := ReadLocalPlayer(username)
    if err != nil {
        return true, err
    }
    
    // If the file doesn't exist yet, treat it as stale so it gets fetched and saved
    if p.Username == "" {
        return true, nil
    }

    if time.Since(p.UpdatedAt) > 2*time.Hour {
        if err := archiveCurrentFile(username); err != nil {
            return false, err
        }
        return true, nil
    }

    return false, nil
}

func SaveData(username string, p Player) error {
	// 1. Ensure the directory exists
	dir := fmt.Sprintf("data/%s", username)
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
    sourcePath := fmt.Sprintf("data/%s/current.json", username)

    data, err := os.ReadFile(sourcePath)
    if err != nil {
        return err // nothing to archive
    }

    var p Player
    if err := json.Unmarshal(data, &p); err != nil {
        return err // Corrupted, don't archive
    }

    destDir := fmt.Sprintf("data/%s/archive", username)
    os.MkdirAll(destDir, 0755)

    destPath := fmt.Sprintf("%s/%s.json", destDir, p.UpdatedAt.Format("20060102_150405"))
    return os.Rename(sourcePath, destPath)
}


// DRY ified
func fetchFromAPI(username string) (Player, error) {
	url := fmt.Sprintf("https://api.wiseoldman.net/v2/players/%s", username)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return Player{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", WOM_API_KEY)
	req.Header.Set("User-Agent", WOM_USER_AGENT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Player{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return Player{}, fmt.Errorf("update failed with status: %d", resp.StatusCode)
	}

	var p Player
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return Player{}, err
	}

	return p, nil
}

func GetPlayerFromAPI(username string) (Player, error) {
    p, err := ReadLocalPlayer(username)
    if err != nil {
        // No local file, fetch straight from API
        fetched, err := fetchFromAPI(username)
        if err != nil {
            return Player{}, err
        }
        err = SaveData(username, fetched)
		if err != nil {
		return Player{}, err
		}
        return fetched, nil
    }

    stale, err := IsDataStale(username)
    if err != nil {
        return Player{}, err
    }

    if !stale {
        return p, nil
    }

    // It's stale, fetch fresh data and save it
    fetched, err := fetchFromAPI(username)
    if err != nil {
        return Player{}, err
    }
    err = SaveData(username, fetched)
	    if err != nil {
		return Player{}, err
	    }
    return fetched, nil
}
