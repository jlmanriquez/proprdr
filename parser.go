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

		// Is a comment line or a empty line
		if strings.HasPrefix(line, commentLine) || len(strings.TrimSpace(line)) == 0 {
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

func findLine(fileName, startWith string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Is a comment line
		if strings.HasPrefix(line, commentLine) {
			continue
		}

		if !strings.HasPrefix(line, startWith) {
			continue
		}

		keyValue := strings.Split(line, "=")
		return keyValue[1], nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}
