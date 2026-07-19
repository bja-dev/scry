package wom

import(
	"fmt"
	"os"
	"bufio"
	"time"
	"encoding/json"
)

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

func GetPlayerData(name string) (Player, error) {
	// checks cache to avoid abusing WOM api
	var p Player
	filename := fmt.Sprintf("data/%s/current.json", username)

	data, err := os.ReadFile(filename)
	if err != nil {
		return Player{}, err
	}
	json.Unmarshal(data, &p)

	if time.Since(p.UpdatedAt) > 2*time.Hour {
		archiveCurrentFile(username, data)
		return p, fmt.Errorf("stale")
	}
	return p, nil
}

func archiveCurrentFile(username string, data []byte) {
    // Create the directory if it doesn't exist
    os.MkdirAll(fmt.Sprintf("data/%s/archive", username), 0755)
    
    // Name it by timestamp
    timestamp := time.Now().Format("20060102_1504")
    archivePath := fmt.Sprintf("data/%s/archive/%s.json", username, timestamp)
    
    // Save the old data into the archive
    os.WriteFile(archivePath, data, 0644)
}

// TODO: request methods, use getplayerdata to check if err != nil then pull :D
