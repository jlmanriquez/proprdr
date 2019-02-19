package proprdr

import (
	"bufio"
	"os"
	"strings"
)

const (
	commentLine = "#"
)

func parseFile(fileName string) (map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	results := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Is a comment line
		if strings.HasPrefix(line, commentLine) {
			continue
		}

		keyValue := strings.Split(line, "=")
		results[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
