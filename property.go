package proprdr

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

// PropertyFile interface for read property file
type PropertyFile interface {
	// Get return a string value
	Get(property string) (string, error)
	// GetAsInt convert string value to int
	GetAsInt(property string) (int, error)
	// GetAsFloat convert string value to float32 or float64 according to bitSize
	GetAsFloat(property string, bitSize int) (float64, error)
	// Containts return true if key exist
	Contains(property string) (exist bool)
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

func (p *propFile) Get(property string) (string, error) {
	value, exist := p.properties[property]
	if !exist {
		return "", errors.New("Property not found")
	}
	return value, nil
}

func (p *propFile) GetAsInt(property string) (int, error) {
	strValue, err := p.Get(property)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (p *propFile) GetAsFloat(property string, bitSize int) (float64, error) {
	strValue, err := p.Get(property)
	if err != nil {
		return 0.00, err
	}

	value, err := strconv.ParseFloat(strValue, bitSize)
	if err != nil {
		return 0.0, nil
	}

	return value, nil
}

func (p *propFile) Size() int {
	return len(p.properties)
}

func (p *propFile) Contains(property string) (exist bool) {
	_, exist = p.properties[property]
	return
}
