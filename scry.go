package scry

import(
	"fmt"
	"bufio"
	"os"
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
