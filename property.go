package proprdr

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// PropertyFile interface for read property file
type PropertyFile interface {
	// Get return a string value
	Get(property string) (string, error)
	// GetAsInt convert string value to int
	GetAsInt(property string) (int, error)
	// GetAsFloat convert string value to float32 or float64 according to bitSize
	GetAsFloat(property string, bitSize int) (float64, error)
	// GetAsBool convert string bool representation to bool type
	GetAsBool(property string) bool
	// GetAll returns a key/value submap where the key starts with startWith
	GetAll(startWith string) map[string]string
	// Containts return true if key exist
	Contains(property string) (exist bool)
	// Size returns the number of properties
	Size() int
	// HasChanged verifies if the file has changed after its creation
	HasChanged() (bool, error)
}

type propFile struct {
	fileName   string
	created    time.Time
	properties map[string]string
}

// New return a implementation of PropertyFile
func New(fileName string) (PropertyFile, error) {
	dictionary, err := parseFile(fileName)
	if err != nil {
		return nil, err
	}

	return &propFile{fileName: fileName, properties: dictionary, created: time.Now()}, nil
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

func (p *propFile) GetAsBool(property string) bool {
	value, _ := p.Get(property)

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return boolValue
}

func (p *propFile) GetAll(startWith string) map[string]string {
	result := map[string]string{}

	for key, value := range p.properties {
		if strings.HasPrefix(key, startWith) {
			result[key] = value
		}
	}

	return result
}

func (p *propFile) HasChanged() (bool, error) {
	info, err := os.Lstat(p.fileName)
	if err != nil {
		return false, err
	}

	modTime := info.ModTime()
	if modTime.After(p.created) {
		return true, nil
	}

	return false, nil
}

func (p *propFile) Refresh() error {
	pfile, err := New(p.fileName)
	if err != nil {
		return err
	}

	p = pfile.(*propFile)
	return nil
}

func (p *propFile) UGet(property string) (string, error) {
	_, exist := p.properties[property]
	if !exist {
		return "", errors.New("Property not found")
	}

	value, err := findLine(p.fileName, property)
	if err != nil {
		value, gErr := p.Get(property)
		if gErr != nil {
			return "", err
		}

		// it was not possible to obtain the updated value
		// returns the current value loaded in memory, but
		// reporting error
		return value, err
	}

	// updates old property value
	p.properties[property] = value

	return value, nil
}
