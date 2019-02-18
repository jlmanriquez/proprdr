package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// PropertyFile interface for read property file
type PropertyFile interface {
	// Get return a string value
	Get(property string) string
	// GetAsInt convert string value to int
	GetAsInt(property string) (int, error)
	// Size returns the number of properties
	Size() int
}

type propFile struct {
	fileName   string
	properties map[string]string
}

// New return a implementation of PropertyFile
func New(fileName string) (PropertyFile, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	newPropertyFile := &propFile{fileName: fileName, properties: map[string]string{}}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		keyValue := strings.Split(line, "=")

		newPropertyFile.properties[keyValue[0]] = keyValue[1]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return newPropertyFile, nil
}

func (p *propFile) Get(property string) string {
	return p.properties[property]
}

func (p *propFile) GetAsInt(property string) (int, error) {
	value, err := strconv.Atoi(p.Get(property))
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (p *propFile) Size() int {
	return len(p.properties)
}
