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
	// Get return a string value. If property does not exist, return a errors
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
	// Refresh reread the property file and replace it in memory
	Refresh() error
	// UGet reread the property file to get the property and update it in memory
	UGet(property string) (string, error)
}

type properties struct {
	fileName  string
	created   time.Time
	keyValues map[string]string
}

// New return a implementation of PropertyFile
func New(fileName string) (PropertyFile, error) {
	dictionary, err := parseFile(fileName)
	if err != nil {
		return nil, err
	}

	return &properties{fileName: fileName, keyValues: dictionary, created: time.Now()}, nil
}

func (p *properties) Get(property string) (string, error) {
	value, exist := p.keyValues[property]
	if !exist {
		return "", errors.New("Property not found")
	}
	return value, nil
}

func (p *properties) GetAsInt(property string) (int, error) {
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

func (p *properties) GetAsFloat(property string, bitSize int) (float64, error) {
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

func (p *properties) Size() int {
	return len(p.keyValues)
}

func (p *properties) Contains(property string) (exist bool) {
	_, exist = p.keyValues[property]
	return
}

func (p *properties) GetAsBool(property string) bool {
	value, _ := p.Get(property)

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return boolValue
}

func (p *properties) GetAll(startWith string) map[string]string {
	result := map[string]string{}

	for key, value := range p.keyValues {
		if strings.HasPrefix(key, startWith) {
			result[key] = value
		}
	}

	return result
}

func (p *properties) HasChanged() (bool, error) {
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

func (p *properties) Refresh() error {
	pfile, err := New(p.fileName)
	if err != nil {
		return err
	}

	newProperties := pfile.(*properties)

	p.created = newProperties.created
	p.keyValues = newProperties.keyValues
	return nil
}

func (p *properties) UGet(property string) (string, error) {
	_, exist := p.keyValues[property]
	if !exist {
		return "", errors.New("Property not found")
	}

	value, err := findLine(p.fileName, property)
	if err != nil {
		value, gErr := p.Get(property)
		if gErr != nil {
			return "", gErr
		}

		// it was not possible to obtain the updated value
		// returns the current value loaded in memory, but
		// reporting error
		return value, err
	}

	// updates old property value
	p.keyValues[property] = value

	return value, nil
}
